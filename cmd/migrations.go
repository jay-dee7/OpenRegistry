package cmd

import (
	"crypto/tls"
	"database/sql"
	"errors"
	"net"
	"os"
	"strings"
	"time"

	"github.com/containerish/OpenRegistry/store/v2/types"
	"github.com/fatih/color"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/urfave/cli/v2"
)

func NewMigrationsCommand() *cli.Command {
	return &cli.Command{
		Name:    "migrations",
		Aliases: []string{"m"},
		Usage:   "Run database migrations for OpenRegistry data store",
		Subcommands: []*cli.Command{
			newDatabaseInitCommand(),
			newMigrationsRunCommand(),
			newMigrationsRollbackCommand(),
			newMigrationsGenrateCommand(),
			newDatabaseResetCommand(),
		},
		Action: func(ctx *cli.Context) error {
			return nil
		},
	}
}

func getOpenRegistryDB(connector *pgdriver.Connector) *bun.DB {
	sqlDB := sql.OpenDB(connector)
	bunWrappedDB := bun.NewDB(sqlDB, pgdialect.New())
	if err := bunWrappedDB.Ping(); err != nil {
		color.Red("error connecting to database: %s", err)
		os.Exit(1100)
	}
	bunWrappedDB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return bunWrappedDB
}

func getAdminBunDB(connector *pgdriver.Connector) *bun.DB {
	sqlDB := sql.OpenDB(connector)
	bunWrappedDB := bun.NewDB(sqlDB, pgdialect.New())
	if err := bunWrappedDB.Ping(); err != nil {
		color.Red("error connecting to database: %s", err)
		os.Exit(1100)
	}
	bunWrappedDB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return bunWrappedDB
}

func createOpenRegistryTables(ctx *cli.Context, db *bun.DB) error {
	_, err := db.NewCreateTable().Model(&types.User{}).Table().IfNotExists().Exec(ctx.Context)
	if err != nil {
		return errors.New(
			color.RedString("Table=users Created=❌ Error=%s", err),
		)
	}
	color.Green(`Table "users" created ✔︎`)

	_, err = db.NewCreateTable().Model(&types.ContainerImageLayer{}).Table().IfNotExists().Exec(ctx.Context)
	if err != nil {
		return errors.New(
			color.RedString("Table=layers Created=❌ Error=%s", err),
		)
	}
	color.Green(`Table "layers" created ✔︎`)

	_, err = db.
		NewCreateTable().
		Model(&types.ContainerImageRepository{}).
		Table().
		IfNotExists().
		Exec(ctx.Context)
	if err != nil {
		return errors.New(
			color.RedString("Table=repositories Created=❌ Error=%s", err),
		)
	}
	color.Green(`Table "repositories" created ✔︎`)

	_, err = db.NewCreateTable().Model(&types.ImageManifest{}).Table().IfNotExists().Exec(ctx.Context)
	if err != nil {
		return errors.New(
			color.RedString("Table=image_manifests Created=❌ Error=%s", err),
		)
	}
	color.Green(`Table "image_manifests" created ✔︎`)

	_, err = db.NewCreateTable().Model(&types.Session{}).Table().IfNotExists().Exec(ctx.Context)
	if err != nil {
		return errors.New(
			color.RedString("Table=sessions Created=❌ Error=%s", err),
		)
	}
	color.Green(`Table "sessions" created ✔︎`)

	_, err = db.NewCreateTable().Model(&types.WebauthnSession{}).Table().IfNotExists().Exec(ctx.Context)
	if err != nil {
		return errors.New(
			color.RedString("Table=webauthn_sessions Created=❌ Error=%s", err),
		)
	}
	color.Green(`Table "webauthn_sessions" created ✔︎`)

	_, err = db.NewCreateTable().Model(&types.WebauthnCredential{}).Table().IfNotExists().Exec(ctx.Context)
	if err != nil {
		return errors.New(
			color.RedString("Table=webauthn_credentials Created=❌ Error=%s", err),
		)
	}
	color.Green(`Table "webauthn_credentials" created ✔︎`)

	_, err = db.NewCreateTable().Model(&types.Email{}).Table().IfNotExists().Exec(ctx.Context)
	if err != nil {
		return errors.New(
			color.RedString("Table=emails Created=❌ Error=%s", err),
		)
	}
	color.Green(`Table "emails" created ✔︎`)

	return nil
}

func getDBConnectorFromCtx(isAdmin bool, opts *databaseOptions) *pgdriver.Connector {
	if opts.openRegistryDSN != "" && opts.adminDSN != "" {
		panic("found DSN for both, openregistry and admin db, but only one should be present")
	}

	dsn := opts.openRegistryDSN
	if dsn == "" {
		dsn = opts.adminDSN
	}

	if dsn != "" {
		return pgdriver.NewConnector(pgdriver.WithDSN(dsn), pgdriver.WithInsecure(opts.insecure))
	}

	if isAdmin {
		return pgdriver.NewConnector(
			pgdriver.WithNetwork("tcp"),
			pgdriver.WithAddr(opts.address),
			//nolint
			pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: opts.insecure}),
			pgdriver.WithInsecure(opts.insecure),
			pgdriver.WithUser(opts.adminUsername),
			pgdriver.WithTimeout(opts.timeout),
			pgdriver.WithDatabase(opts.adminDB),
			pgdriver.WithPassword(opts.adminPassword),
			pgdriver.WithApplicationName("OpenRegistry"),
		)
	}

	return pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(opts.address),
		//nolint
		pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: opts.insecure}),
		pgdriver.WithInsecure(opts.insecure),
		pgdriver.WithUser(opts.username),
		pgdriver.WithTimeout(opts.timeout),
		pgdriver.WithDatabase(opts.database),
		pgdriver.WithPassword(opts.password),
		pgdriver.WithApplicationName("OpenRegistry"),
	)
}

func createOpenRegistryDatabase(ctx *cli.Context, opts *databaseOptions) (*bun.DB, error) {
	adminConnector := getDBConnectorFromCtx(true, opts)
	adminDB := getAdminBunDB(adminConnector)

	_, err := adminDB.ExecContext(
		ctx.Context,
		"CREATE USER ? WITH ENCRYPTED PASSWORD ?",
		bun.Ident(opts.username),
		opts.password,
	)
	if err != nil && !strings.Contains(err.Error(), "SQLSTATE=42710") {
		return nil, errors.New(
			color.RedString("Action=CreateUser Created=❌ Error=%s", err),
		)
	}

	_, err = adminDB.Exec("create database ? with owner = ?", bun.Ident(opts.database), opts.username)
	if err != nil && !strings.Contains(err.Error(), "SQLSTATE=42P04") {
		return nil, errors.New(
			color.RedString("Action=CreateDatabase Created=❌ Error=%s", err),
		)
	}

	_, err = adminDB.
		ExecContext(
			ctx.Context,
			"GRANT ALL PRIVILEGES ON DATABASE ?0 to ?1",
			// "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ?",
			bun.Ident(opts.database),
			bun.Ident(opts.username),
		)
	if err != nil {
		return nil, errors.New(
			color.RedString("Action=GrantDBPrivleges Created=❌ Error=%s", err),
		)
	}
	color.Green(`Action "GrantDBPrivleges" succeeded ✔︎`)

	openregistryDB := getOpenRegistryDB(getDBConnectorFromCtx(false, opts))
	_, err = adminDB.
		ExecContext(
			ctx.Context,
			// "GRANT ALL ON SCHEMA public to ?",
			"GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ?",
			bun.Ident(opts.username),
		)
	if err != nil {
		return nil, errors.New(
			color.RedString("Action=GrantAll Created=❌ Error=%s", err),
		)
	}
	color.Green(`Action "GrantAllPrivleges" succeeded ✔︎`)
	return openregistryDB, nil
}

type databaseOptions struct {
	address         string
	database        string
	username        string
	password        string
	adminDSN        string
	openRegistryDSN string
	adminDB         string
	adminUsername   string
	adminPassword   string
	timeout         time.Duration
	insecure        bool
}

func parseDatabaseFlags(ctx *cli.Context) *databaseOptions {
	address := net.JoinHostPort(ctx.String("host"), ctx.String("port"))
	database := ctx.String("database")
	username := ctx.String("username")
	password := ctx.String("password")
	timeout := ctx.Duration("timeout")
	insecure := ctx.Bool("insecure")
	openRegistryDSN := ctx.String("openregistry-db-dsn")
	adminDsn := ctx.String("admin-db-dsn")
	adminDB := ctx.String("admin-db")
	adminUsername := ctx.String("admin-db-username")
	adminPassword := ctx.String("admin-db-password")

	return &databaseOptions{
		address:         address,
		database:        database,
		username:        username,
		password:        password,
		timeout:         timeout,
		insecure:        insecure,
		openRegistryDSN: openRegistryDSN,
		adminDSN:        adminDsn,
		adminDB:         adminDB,
		adminUsername:   adminUsername,
		adminPassword:   adminPassword,
	}
}
