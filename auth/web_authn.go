package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/fatih/color"
	"github.com/google/uuid"

	"github.com/containerish/OpenRegistry/types"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)

func (a *auth) BeginRegistration(ctx echo.Context) error {
	ctx.Set(types.HandlerStartTime, time.Now())
	var user types.User

	if err := json.NewDecoder(ctx.Request().Body).Decode(&user); err != nil {
		echoErr := ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "invalid JSON object",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}
	_ = ctx.Request().Body.Close()

	err := user.Validate()
	if err != nil {
		echoErr := ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "invalid data provided for user login",
			"code":    "INVALID_CREDENTIALS",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}

	key := user.Email
	if user.Username != "" {
		key = user.Username
	}

	userFromDb, err := a.pgStore.GetUser(ctx.Request().Context(), key, true)
	if err != nil {
		if errors.Unwrap(err) == pgx.ErrNoRows {
			//user does not exist, create new user
			if err = a.pgStore.AddUser(ctx.Request().Context(), &user); err != nil {
				echoErr := ctx.JSON(http.StatusBadRequest, echo.Map{
					"error":   err.Error(),
					"message": "database error, failed to add user",
				})
				a.logger.Log(ctx, err)
				return echoErr
			}
			// user successfully created
			options, sessionData, wErr := a.webAuthN.BeginRegistration(&user)
			if wErr != nil {
				echoErr := ctx.JSON(http.StatusBadRequest, echo.Map{
					"error":   err.Error(),
					"message": "error begin registration",
				})
				a.logger.Log(ctx, err)
				return echoErr
			}
			// store session data in DB
			if err := a.pgStore.AddWebAuthSessionData(ctx.Request().Context(), user.Id, sessionData, "registration"); err != nil {
				echoErr := ctx.JSON(http.StatusInternalServerError, echo.Map{
					"error":   err.Error(),
					"message": "database error, failed to add web authn session data for new user",
				})
				a.logger.Log(ctx, err)
				return echoErr
			}
			//return response
			echoErr := ctx.JSON(http.StatusOK, echo.Map{
				"message": "registration successful",
				"options": &options,
			})
			a.logger.Log(ctx, echoErr)
			return echoErr

		}
		echoErr := ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "database error, failed to get user",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}

	options, sessionData, err := a.webAuthN.BeginRegistration(userFromDb)
	if err != nil {
		echoErr := ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "error begin registration",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}

	// store session data in DB
	if err := a.pgStore.AddWebAuthSessionData(ctx.Request().Context(), userFromDb.Id, sessionData, "registration"); err != nil {
		echoErr := ctx.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"message": "database error, failed to add web authn session data for existing user",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}

	echoErr := ctx.JSON(http.StatusOK, echo.Map{
		"options": &options,
	})
	a.logger.Log(ctx, echoErr)
	return echoErr
}

func (a *auth) FinishRegistration(ctx echo.Context) error {
	ctx.Set(types.HandlerStartTime, time.Now())

	username := ctx.QueryParam("username")
	userFromDB, err := a.pgStore.GetUser(ctx.Request().Context(), username, false)
	if err != nil {
		echoErr := ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "database error, user not found",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}

	sessionData, err := a.pgStore.GetWebAuthNSessionData(ctx.Request().Context(), userFromDB.Id, "registration")
	if err != nil {
		echoErr := ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "database error, session data not found",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}

	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(ctx.Request().Body)
	if err != nil {
		echoErr := ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "error parsing credential creation response body",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}
	defer ctx.Request().Body.Close()
	color.Red("sessionData: %+v", sessionData)
	color.Yellow("userFromDB: %+v", userFromDB)
	color.Green("parsedResponse: %+v", parsedResponse)

	credentials, err := a.webAuthN.CreateCredential(userFromDB, *sessionData, parsedResponse)
	if err != nil {
		echoErr := ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "error creating webauthn credentials",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}

	userFromDB.AddWebAuthNCredential(credentials)
	if err := a.pgStore.AddWebAuthNCredentials(ctx.Request().Context(), userFromDB.Id, credentials); err != nil {
		echoErr := ctx.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"message": "database error storing webauthn credentials",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}

	echoErr := ctx.JSON(http.StatusOK, echo.Map{
		"message": "registration successful",
	})
	a.logger.Log(ctx, echoErr)
	return echoErr
}

func (a *auth) BeginLogin(ctx echo.Context) error {
	ctx.Set(types.HandlerStartTime, time.Now())

	username := ctx.QueryParam("username")
	userFromDB, err := a.pgStore.GetUser(ctx.Request().Context(), username, false)
	if err != nil {
		echoErr := ctx.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"message": "database error: user not found",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}

	creds, err := a.pgStore.GetWebAuthNCredentials(ctx.Request().Context(), userFromDB.Id)
	if err != nil {
		echoErr := ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "error getting credentials for user",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}

	userFromDB.AddWebAuthNCredential(creds)
	options, sessionData, err := a.webAuthN.BeginLogin(userFromDB)
	if err != nil {
		echoErr := ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "error begin login",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}

	if err := a.pgStore.AddWebAuthSessionData(ctx.Request().Context(), userFromDB.Id, sessionData, "authentication"); err != nil {
		echoErr := ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "database error: storing session data while web authn begin login",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}

	echoErr := ctx.JSON(http.StatusOK, echo.Map{
		"options": &options,
	})
	a.logger.Log(ctx, echoErr)
	return echoErr

}
func (a *auth) FinishLogin(ctx echo.Context) error {
	ctx.Set(types.HandlerStartTime, time.Now())

	username := ctx.QueryParam("username")
	userFromDb, err := a.pgStore.GetUser(ctx.Request().Context(), username, false)
	if err != nil {
		echoErr := ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "database error: user not found",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}

	sessionData, err := a.pgStore.GetWebAuthNSessionData(ctx.Request().Context(), userFromDb.Id, "authentication")
	if err != nil {
		echoErr := ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "database error: session data for user not found in finish login",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}

	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(ctx.Request().Body)
	if err != nil {
		echoErr := ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "parsing error: could not parse credential request body in finish login",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}
	defer ctx.Request().Body.Close()
	color.Red("parsed Response: %+v", parsedResponse)
	color.Red("session data: %+v", *sessionData)
	creds, err := a.pgStore.GetWebAuthNCredentials(ctx.Request().Context(), userFromDb.Id)
	if err != nil {
		echoErr := ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "error getting credentials for user",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}

	userFromDb.AddWebAuthNCredential(creds)

	//Validate login gives back credential
	_, err = a.webAuthN.ValidateLogin(userFromDb, *sessionData, parsedResponse)
	if err != nil {
		echoErr := ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "could not validate user login",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}

	access, err := a.newWebLoginToken(userFromDb.Id, userFromDb.Username, "access")
	if err != nil {
		echoErr := ctx.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"message": "error creating web login token",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}

	refresh, err := a.newWebLoginToken(userFromDb.Id, userFromDb.Username, "refresh")
	if err != nil {
		echoErr := ctx.JSON(http.StatusInternalServerError, echo.Map{
			"error":   err.Error(),
			"message": "error creating refresh token",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}
	id := uuid.NewString()
	sessionId := fmt.Sprintf("%s:%s", id, userFromDb.Id)

	if err = a.pgStore.AddSession(ctx.Request().Context(), id, refresh, userFromDb.Username); err != nil {
		echoErr := ctx.JSON(http.StatusBadRequest, echo.Map{
			"error":   err.Error(),
			"message": "error creating session",
		})
		a.logger.Log(ctx, err)
		return echoErr
	}

	sessionCookie := a.createCookie("session_id", sessionId, false, time.Now().Add(time.Hour*750))
	accessCookie := a.createCookie("access", access, true, time.Now().Add(time.Hour*750))
	refreshCookie := a.createCookie("refresh", refresh, true, time.Now().Add(time.Hour*750))
	ctx.SetCookie(accessCookie)
	ctx.SetCookie(refreshCookie)
	ctx.SetCookie(sessionCookie)

	echoErr := ctx.JSON(http.StatusOK, echo.Map{
		"message": "Login Success",
	})

	a.logger.Log(ctx, echoErr)
	return echoErr
}
