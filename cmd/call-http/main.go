package main

import (
	"os"

	"github.com/Jumpaku/api-regression-detector/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/call/http"
	"github.com/docopt/docopt-go"
)

const doc = `Regression detector call-http.
call-http calls HTTP API: sending JSON request and receiving JSON response.

Usage:
	call-http [--header=<http-header>]... <endpoint-url> <http-method>
	call-http -h | --help
	call-http --version

Options:
	<http-header>      Header in the form 'Key: value'.
	<endpoint-url>     The URL of the HTTP endpoint which may has path parameters enclosed in '[' and ']'.
	<http-method>      One of GET, HEAD, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, or PATCH.
	-h --help          Show this screen.
	--version          Show version.`

func main() {
	args, _ := docopt.ParseArgs(doc, os.Args[1:], "1.0.0")
	code := cmd.RunCallHTTP(
		cmd.Stdio,
		args["<endpoint-url>"].(string),
		http.Method(args["<http-method>"].(string)),
		args["--header"].([]string),
	)

	os.Exit(code)
}
