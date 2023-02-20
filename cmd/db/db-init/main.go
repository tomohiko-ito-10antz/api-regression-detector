package main

import (
	"context"
	"os"

	"github.com/Jumpaku/api-regression-detector/cmd/db"
	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/tables"
	"github.com/Jumpaku/api-regression-detector/lib/log"
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
	code, err := RunInit(
		args["<database-driver>"].(string),
		args["<connection-string>"].(string))
	if err != nil {
		log.Stderr("Error\n%+v", err)
	}
	os.Exit(code)
}

func RunInit(databaseDriver string, connectionString string) (code int, err error) {
	driver, err := db.NewDriver(databaseDriver)
	if err != nil {
		return 1, errors.Wrap(err, "fail to new %s", databaseDriver)
	}

	err = driver.Open(connectionString)
	if err != nil {
		return 1, errors.Wrap(err, "fail to connect %s", connectionString)
	}

	defer func() {
		err = errors.Wrap(errors.Join(err, driver.Close()), "fail RunInit")
		if err != nil {
			code = 1
		}
	}()

	initTables, err := tables.LoadInitTables(os.Stdin)
	if err != nil {
		return 1, errors.Wrap(err, "fail to load JSON from stdin")
	}

	err = cmd.Init(context.Background(), driver.DB, initTables, driver.SchemaGetter, driver.RowClearer, driver.RowCreator)
	if err != nil {
		return 1, errors.Wrap(err, "fail Init")
	}

	return 0, nil
}
