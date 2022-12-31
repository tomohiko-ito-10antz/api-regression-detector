package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/Jumpaku/api-regression-detector/cmd"
	"github.com/Jumpaku/api-regression-detector/log"
	"github.com/Jumpaku/api-regression-detector/mysql"
	"github.com/Jumpaku/api-regression-detector/postgres"
	"github.com/Jumpaku/api-regression-detector/sqlite"
)

func fail(err error) {
	log.Stderr("Error\n%v", err)
	panic(err)
}

type Driver struct {
	Name     string
	DB       *sql.DB
	Select   cmd.Select
	Truncate cmd.Truncate
	Insert   cmd.Insert
}

func (d *Driver) Close() error {
	return d.DB.Close()
}
func Connect(name string, connectionString string) (*Driver, error) {
	switch name {
	case "mysql":
		db, err := sql.Open(name, connectionString)
		if err != nil {
			return nil, err
		}
		return &Driver{
			Name:     name,
			DB:       db,
			Select:   mysql.Select(),
			Truncate: mysql.Truncate(),
			Insert:   mysql.Insert(),
		}, nil
	case "postgres":
		db, err := sql.Open(name, connectionString)
		if err != nil {
			return nil, err
		}
		return &Driver{
			Name:     name,
			DB:       db,
			Select:   postgres.Select(),
			Truncate: postgres.Truncate(),
			Insert:   postgres.Insert(),
		}, nil
	case "sqlite3":
		db, err := sql.Open(name, connectionString)
		if err != nil {
			return nil, err
		}
		return &Driver{
			Name:     name,
			DB:       db,
			Select:   sqlite.Select(),
			Truncate: sqlite.Truncate(),
			Insert:   sqlite.Insert(),
		}, nil
	default:
		return nil, fmt.Errorf("invalid driver name")
	}
}
func main() {
	usage := `Regression detector.
Examples of prepare
./program prepare mysql "root:password@(mysql)/main" <examples/prepare.json
./program prepare postgres "user=postgres password=password host=postgres dbname=main sslmode=disable" <examples/prepare.json
./program prepare sqlite3 "file:examples/sqlite/sqlite.db" <examples/prepare.json

Examples of dump
./program dump mysql "root:password@(mysql)/main" >examples/dump.json
./program dump postgres "user=postgres password=password host=postgres dbname=main sslmode=disable" >examples/dump.json
./program dump sqlite3 "file:examples/sqlite/sqlite.db" >examples/dump.json

Examples of compare
./program compare examples/expected.json examples/actual.json
./program compare --verbose examples/expected.json examples/actual.json


Usage:
  program prepare <database-driver> <connection-string>
  program dump <database-driver> <connection-string>
  program compare [--verbose] [--strict] <expected-json> <actual-json>
  program -h | --help
  program --version

Options:
  -h --help       Show this screen.
  --version       Show version.
  --verbose       Show verbose difference. [default: false]
  --strict        Disallow superset match. [default: false]`

	args, _ := docopt.ParseArgs(usage, os.Args[1:], "1.0.0")
	fmt.Println(args)
	var (
		code int
		err  error
	)
	switch {
	case args["compare"]:
		code, err = RunCompare(
			args["<expected-json>"].(string),
			args["<actual-json>"].(string),
			args["--verbose"].(bool),
			args["--strict"].(bool))
	case args["prepare"]:
		fmt.Println("prepare", args["<database-driver>"].(string), args["<connection-string>"].(string))
	case args["dump"]:
		fmt.Println("dump", args["<database-driver>"].(string), args["<connection-string>"].(string))
	}

	if err != nil {
		fail(err)
	}
	os.Exit(code)
}

func RunCompare(expectedJson string, actualJson string, verbose bool, strict bool) (code int, err error) {
	expectedJsonFile, err := os.Open(expectedJson)
	if err != nil {
		return 1, err
	}
	actualJsonFile, err := os.Open(actualJson)
	if err != nil {
		return 1, err
	}
	match, detail, err := cmd.Compare(expectedJsonFile, actualJsonFile)
	if err != nil {
		return 1, err
	}
	fmt.Println(match)
	if verbose {
		fmt.Println(detail)
	}
	switch match {
	case cmd.CompareResult_FullMatch:
		return 0, nil
	case cmd.CompareResult_SupersetMatch:
		if strict {
			return 1, nil
		}
		return 0, nil
	default:
		return 1, nil
	}
}
