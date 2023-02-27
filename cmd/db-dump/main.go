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
	code := cmd.RunDump(
		cmd.Stdio,
		args["<database-driver>"].(string),
		args["<connection-string>"].(string))

	os.Exit(code)
}
