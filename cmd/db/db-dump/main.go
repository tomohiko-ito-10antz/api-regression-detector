package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Jumpaku/api-regression-detector/cmd/db"
	"github.com/Jumpaku/api-regression-detector/lib/cmd"
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
	code, err := RunDump(
		args["<database-driver>"].(string),
		args["<connection-string>"].(string))
	if err != nil {
		fmt.Printf("Error\n%s\n%+v", err, err)
	}
	os.Exit(code)
}

func RunDump(databaseDriver string, connectionString string) (code int, err error) {
	errorInfo := errors.Info{"databaseDriver": databaseDriver}

	driver, err := db.NewDriver(databaseDriver)
	if err != nil {
		return 1, errors.Wrap(errors.BadArgs.Err(err), errorInfo.AppendTo("fail RunDump"))
	}

	errorInfo = errorInfo.With("connectionString", connectionString)

	if err := driver.Open(connectionString); err != nil {
		return 1, errors.Wrap(errors.IOFailure.Err(err), errorInfo.AppendTo("fail to open database"))
	}
	defer func() {
		if errs := errors.Wrap(errors.Join(err, driver.Close()), errorInfo.AppendTo("fail RunDump")); errs != nil {
			err = errs
			code = 1
		}
	}()

	tableNames, err := jsonio.LoadJson[[]string](os.Stdin)
	if err != nil {
		return 1, errors.Wrap(err, "fail to load init tables as JSON from stdin")
	}

	dump, err := cmd.Dump(context.Background(), driver.DB, tableNames, driver.SchemaGetter, driver.RowLister)
	if err != nil {
		return 1, errors.Wrap(err, errorInfo.With("tableNames", tableNames).AppendTo("fail to get tables from database"))
	}

	if err := tables.SaveDumpTables(dump, os.Stdout); err != nil {
		return 1, errors.Wrap(err, "fail to dump tables as JSON to stdout")
	}

	return 0, nil
}
