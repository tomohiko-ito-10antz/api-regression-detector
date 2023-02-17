package main

import (
	"os"

	"github.com/Jumpaku/api-regression-detector/cli"
	"github.com/Jumpaku/api-regression-detector/lib/call/http"
	"github.com/Jumpaku/api-regression-detector/lib/log"
	"github.com/docopt/docopt-go"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/googleapis/go-sql-spanner"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var (
		code int
		err  error
	)

	args, _ := docopt.ParseArgs(cli.GetDoc(), os.Args[1:], "1.0.0")

	switch {
	case args["compare"]:
		code, err = cli.RunCompare(
			args["<expected-json>"].(string),
			args["<actual-json>"].(string),
			args["--verbose"].(bool),
			args["--strict"].(bool))
	case args["init"]:
		code, err = cli.RunInit(
			args["<database-driver>"].(string),
			args["<connection-string>"].(string))
	case args["dump"]:
		code, err = cli.RunDump(
			args["<database-driver>"].(string),
			args["<connection-string>"].(string))
	case args["call"]:
		switch {
		case args["http"]:
			code, err = cli.RunCallHTTP(
				args["<endpoint-url>"].(string),
				http.Method(args["<http-method>"].(string)))
		case args["grpc"]:
			code, err = cli.RunCallGRPC(
				args["<grpc-endpoint>"].(string),
				args["<grpc-full-method>"].(string))
		default:
		}
	default:
	}

	if err != nil {
		fail(err)
	}

	os.Exit(code)
}

func fail(err error) {
	log.Stderr("Error\n%+v", err)
	panic(err)
}
