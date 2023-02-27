package main

import (
	"os"

	"github.com/Jumpaku/api-regression-detector/cmd"
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
	code := cmd.RunInit(
		cmd.Stdio,
		args["<database-driver>"].(string),
		args["<connection-string>"].(string))

	os.Exit(code)
}
