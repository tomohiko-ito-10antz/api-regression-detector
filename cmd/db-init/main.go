package main

import (
	"context"
	"os"

	"github.com/Jumpaku/api-regression-detector/cmd"
	libcmd "github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/cmd/cli"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/tables"
	"github.com/docopt/docopt-go"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/googleapis/go-sql-spanner"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const doc = `Regression detector db-init.
db-init initializes tables according to JSON data.

Usage:
	db-init <database-driver> <connection-string>
	db-init -h | --help
	db-init --version

Options:
	<database-driver>   Supported database driver name which is one of mysql, spanner, sqlite3, or postgres
	<connection-string> Connection string corresponding to the database driver.
	-h --help          Show this screen.
	--version          Show version.`

func main() {
	args, _ := docopt.ParseArgs(doc, os.Args[1:], "1.0.0")
	code := RunInit(
		cmd.Stdio,
		args["<database-driver>"].(string),
		args["<connection-string>"].(string))

	os.Exit(code)
}

func RunInit(stdio *cli.Stdio, databaseDriver string, connectionString string) (code int) {
	errorInfo := errors.Info{"databaseDriver": databaseDriver}

	driver, err := cmd.NewDriver(databaseDriver)
	if err != nil {
		cmd.PrintError(stdio.Stderr, errors.Wrap(errors.BadArgs.Err(err), errorInfo.AppendTo("fail RunInit")))
		return 1
	}

	errorInfo = errorInfo.With("connectionString", connectionString)

	if err := driver.Open(connectionString); err != nil {
		cmd.PrintError(stdio.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail to open database")))
		return 1
	}
	defer func() {
		if err := errors.Join(err, driver.Close()); err != nil {
			cmd.PrintError(stdio.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunInit")))
			code = 1
		}
	}()

	initTables, err := tables.LoadInitTables(os.Stdin)
	if err != nil {
		cmd.PrintError(stdio.Stderr, errors.Wrap(err, "fail to load init tables from stdin"))
		return 1
	}

	err = libcmd.Init(context.Background(), driver.DB, initTables, driver.SchemaGetter, driver.RowClearer, driver.RowCreator)
	if err != nil {
		cmd.PrintError(stdio.Stderr, errors.Wrap(err, errorInfo.With("initTables", initTables).AppendTo("fail to init tables in database")))
		return 1
	}

	return 0
}
