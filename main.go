package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Jumpaku/api-regression-detector/cmd"
	"github.com/Jumpaku/api-regression-detector/impl/mysql"
	"github.com/Jumpaku/api-regression-detector/impl/postgres"
	"github.com/Jumpaku/api-regression-detector/impl/spanner"
	"github.com/Jumpaku/api-regression-detector/impl/sqlite"
	"github.com/Jumpaku/api-regression-detector/io"
	"github.com/Jumpaku/api-regression-detector/log"
	"github.com/docopt/docopt-go"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/googleapis/go-sql-spanner"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/multierr"
)

func fail(err error) {
	log.Stderr("Error\n%v", err)
	panic(err)
}

func Connect(name string, connectionString string) (*cmd.Driver, error) {
	var driver *cmd.Driver
	switch name {
	case "mysql":
		driver = &cmd.Driver{
			Select:   mysql.Select(),
			Truncate: mysql.ClearRows(),
			Insert:   mysql.Insert(),
		}
	case "postgres":
		driver = &cmd.Driver{
			Select:   postgres.Select(),
			Truncate: postgres.Truncate(),
			Insert:   postgres.Insert(),
		}
	case "sqlite3":
		driver = &cmd.Driver{
			Select:   sqlite.Select(),
			Truncate: sqlite.Truncate(),
			Insert:   sqlite.Insert(),
		}
	case "spanner":
		driver = &cmd.Driver{
			Select:   spanner.Select(),
			Truncate: spanner.Truncate(),
			Insert:   spanner.Insert(),
		}
	default:
		return nil, fmt.Errorf("invalid driver name")
	}
	err := driver.Open(name, connectionString)
	if err != nil {
		return nil, err
	}
	return driver, nil
}
func main() {
	usage := `Regression detector.
The following commands are available:
* init: It initializes database according to json provided by stdin.
* dump: It outputs database according to json provided by stdin.
* compare: It compares two JSON files and outputs the comparison result to stdout.

Usage:
  program init <database-driver> <connection-string>
  program dump <database-driver> <connection-string>
  program compare [--verbose] [--strict] <expected-json> <actual-json>
  program -h | --help
  program --version

Options:
  -h --help          Show this screen.
  --version          Show version.
  --verbose          Show verbose difference. [default: false]
  --strict           Disallow superset match. [default: false]`

	args, _ := docopt.ParseArgs(usage, os.Args[1:], "1.0.0")
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
	case args["init"]:
		code, err = RunInit(
			args["<database-driver>"].(string),
			args["<connection-string>"].(string))
	case args["dump"]:
		code, err = RunDump(
			args["<database-driver>"].(string),
			args["<connection-string>"].(string))
	default:
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
	defer func() {
		err = multierr.Combine(err, expectedJsonFile.Close())
		if err != nil {
			code = 1
		}
	}()
	actualJsonFile, err := os.Open(actualJson)
	defer func() {
		err = multierr.Combine(err, actualJsonFile.Close())
		if err != nil {
			code = 1
		}
	}()
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

func RunInit(databaseDriver string, connectionString string) (code int, err error) {
	driver, err := Connect(databaseDriver, connectionString)
	if err != nil {
		return 1, err
	}
	defer func() {
		err = multierr.Combine(err, driver.Close())
		if err != nil {
			code = 1
		}
	}()
	tables, err := io.Load(os.Stdin)
	if err != nil {
		return 1, err
	}
	err = cmd.Init(context.Background(), driver.DB, tables, driver.Truncate, driver.Insert)
	if err != nil {
		return 1, err
	}
	return 0, nil
}

func RunDump(databaseDriver string, connectionString string) (code int, err error) {
	driver, err := Connect(databaseDriver, connectionString)
	if err != nil {
		return 1, err
	}
	defer func() {
		err = multierr.Combine(err, driver.Close())
		if err != nil {
			code = 1
		}
	}()
	tables, err := io.Load(os.Stdin)
	if err != nil {
		return 1, err
	}
	tableNames := []string{}
	for tableName := range tables {
		tableNames = append(tableNames, tableName)
	}
	dump, err := cmd.Dump(context.Background(), driver.DB, tableNames, driver.Select)
	if err != nil {
		return 1, err
	}
	err = io.Save(dump, os.Stdout)
	if err != nil {
		return 1, err
	}
	return 0, nil
}
