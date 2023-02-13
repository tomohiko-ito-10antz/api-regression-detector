package main

import (
	"context"
	"os"

	"github.com/Jumpaku/api-regression-detector/cli"
	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/log"
	"github.com/Jumpaku/api-regression-detector/lib/rpc"
	http "github.com/Jumpaku/api-regression-detector/lib/rpc/impl/http"
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
		_, _ = cmd.CallHTTP(context.Background(), rpc.Request{}, "https://api:80", "GET", http.Caller())
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
