package main

import (
	"context"
	"os"

	"github.com/Jumpaku/api-regression-detector/cmd"
	libcmd "github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/cmd/cli"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/tables"
	"github.com/docopt/docopt-go"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/googleapis/go-sql-spanner"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const doc = `Regression detector db-dump.
db-dump outputs data within tables in JSON format.

Usage:
	db-dump <database-driver> <connection-string>
	db-dump -h | --help
	db-dump --version

Options:
	<database-driver>   Supported database driver name which is one of mysql, spanner, sqlite3, or postgres
	<connection-string> Connection string corresponding to the database driver.
	-h --help           Show this screen.
	--version           Show version.`

func main() {
	args, _ := docopt.ParseArgs(doc, os.Args[1:], "1.0.0")
	code := RunDump(
		cmd.Stdio,
		args["<database-driver>"].(string),
		args["<connection-string>"].(string))

	os.Exit(code)
}

func RunDump(stdio *cli.Stdio, databaseDriver string, connectionString string) (code int) {
	errorInfo := errors.Info{"databaseDriver": databaseDriver}

	driver, err := cmd.NewDriver(databaseDriver)
	if err != nil {
		cmd.PrintError(stdio.Stderr, errors.Wrap(errors.BadArgs.Err(err), errorInfo.AppendTo("fail RunDump")))
		return 1
	}

	errorInfo = errorInfo.With("connectionString", connectionString)

	if err := driver.Open(connectionString); err != nil {
		cmd.PrintError(stdio.Stderr, errors.Wrap(errors.IOFailure.Err(err), errorInfo.AppendTo("fail to open database")))
		return 1
	}
	defer func() {
		if err := errors.Join(err, driver.Close()); err != nil {
			cmd.PrintError(os.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunDump")))
			code = 1
		}
	}()

	tableNames, err := jsonio.LoadJson[[]string](os.Stdin)
	if err != nil {
		cmd.PrintError(os.Stderr, errors.Wrap(err, "fail to load init tables as JSON from stdin"))
		return 1
	}

	dump, err := libcmd.Dump(context.Background(), driver.DB, tableNames, driver.SchemaGetter, driver.RowLister)
	if err != nil {
		cmd.PrintError(os.Stderr, errors.Wrap(err, errorInfo.With("tableNames", tableNames).AppendTo("fail to get tables from database")))
		return 1
	}

	if err := tables.SaveDumpTables(dump, os.Stdout); err != nil {
		cmd.PrintError(os.Stderr, errors.Wrap(err, "fail to dump tables as JSON to stdout"))
		return 1
	}

	return 0
}
