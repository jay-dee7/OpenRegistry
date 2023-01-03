package auth

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/containerish/OpenRegistry/types"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	AccessCookieKey = "access"
	RefreshCookKey  = "refresh"
	QueryToken      = "token"
)

// JWT basically uses the default JWT middleware by echo, but has a slightly different skipper func
func (a *auth) JWT() echo.MiddlewareFunc {
	privBz, err := os.ReadFile(a.c.Registry.TLS.PrivateKey)
	if err != nil {
		panic(err)
	}
	privkey, err := jwt.ParseRSAPrivateKeyFromPEM(privBz)
	if err != nil {
		panic(err)
	}

	pubBz, err := os.ReadFile(a.c.Registry.TLS.PubKey)
	if err != nil {
		panic(err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBz)
	if err != nil {
		panic(err)
	}

	return middleware.JWTWithConfig(middleware.JWTConfig{
		Skipper: func(ctx echo.Context) bool {
			if strings.HasPrefix(ctx.Request().RequestURI, "/auth") {
				return false
			}

			// if JWT_AUTH is not set, we don't need to perform JWT authentication
			jwtAuth, ok := ctx.Get(JWT_AUTH_KEY).(bool)
			if !ok {
				return true
			}

			if jwtAuth {
				return false
			}

			return true
		},
		BeforeFunc:     middleware.DefaultJWTConfig.BeforeFunc,
		SuccessHandler: middleware.DefaultJWTConfig.SuccessHandler,
		ErrorHandler:   nil,
		ErrorHandlerWithContext: func(err error, ctx echo.Context) error {
			// ErrorHandlerWithContext only logs the failing requtest
			ctx.Set(types.HandlerStartTime, time.Now())
			a.logger.Log(ctx, err)
			return ctx.JSON(http.StatusUnauthorized, echo.Map{
				"error":   err.Error(),
				"message": "missing authentication information",
			})
		},
		KeyFunc: func(t *jwt.Token) (interface{}, error) {
			return pubKey, nil
		},
		ParseTokenFunc: middleware.DefaultJWTConfig.ParseTokenFunc,
		SigningKey:     privkey,
		SigningKeys:    map[string]interface{}{},
		SigningMethod:  jwt.SigningMethodRS256.Name,
		Claims:         &Claims{},
		TokenLookup:    fmt.Sprintf("cookie:%s,header:%s", AccessCookieKey, echo.HeaderAuthorization),
	})
}

// ACL implies a basic Access Control List on protected resources
func (a *auth) ACL() echo.MiddlewareFunc {
	return func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set(types.HandlerStartTime, time.Now())

			m := ctx.Request().Method
			if m == http.MethodGet || m == http.MethodHead {
				return hf(ctx)
			}

			token, ok := ctx.Get("user").(*jwt.Token)
			if !ok {
				a.logger.Log(ctx, fmt.Errorf("ACL: unauthorized"))
				return ctx.NoContent(http.StatusUnauthorized)
			}

			claims, ok := token.Claims.(*Claims)
			if !ok {
				a.logger.Log(ctx, fmt.Errorf("ACL: invalid claims"))
				return ctx.NoContent(http.StatusUnauthorized)
			}

			username := ctx.Param("username")

			user, err := a.pgStore.GetUserById(ctx.Request().Context(), claims.Id, false)
			if err != nil {
				a.logger.Log(ctx, err)
				return ctx.NoContent(http.StatusUnauthorized)
			}
			if user.Username == username {
				return hf(ctx)
			}

			return ctx.NoContent(http.StatusUnauthorized)

		}
	}
}

// JWT basically uses the default JWT middleware by echo, but has a slightly different skipper func
func (a *auth) JWTRest() echo.MiddlewareFunc {
	privBz, err := os.ReadFile(a.c.Registry.TLS.PrivateKey)
	if err != nil {
		panic(err)
	}
	privkey, err := jwt.ParseRSAPrivateKeyFromPEM(privBz)
	if err != nil {
		panic(err)
	}

	pubBz, err := os.ReadFile(a.c.Registry.TLS.PubKey)
	if err != nil {
		panic(err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBz)
	if err != nil {
		panic(err)
	}

	return middleware.JWTWithConfig(middleware.JWTConfig{
		BeforeFunc:     middleware.DefaultJWTConfig.BeforeFunc,
		SuccessHandler: middleware.DefaultJWTConfig.SuccessHandler,
		ErrorHandler:   nil,
		ErrorHandlerWithContext: func(err error, ctx echo.Context) error {
			// ErrorHandlerWithContext only logs the failing requtest
			ctx.Set(types.HandlerStartTime, time.Now())
			a.logger.Log(ctx, err)
			return ctx.JSON(http.StatusUnauthorized, echo.Map{
				"error":   err.Error(),
				"message": "missing authentication information",
			})
		},
		KeyFunc: func(t *jwt.Token) (interface{}, error) {
			return pubKey, nil
		},
		ParseTokenFunc: middleware.DefaultJWTConfig.ParseTokenFunc,
		SigningKey:     privkey,
		SigningMethod:  jwt.SigningMethodRS256.Name,
		Claims:         &Claims{},
		TokenLookup:    fmt.Sprintf("cookie:%s,header:%s", AccessCookieKey, echo.HeaderAuthorization),
	})
}
