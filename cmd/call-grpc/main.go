package main

import (
	"os"

	"github.com/Jumpaku/api-regression-detector/cmd"
	"github.com/docopt/docopt-go"
)

const doc = `Regression detector call-grpc.
call-grpc calls GRPC API: sending JSON request and receiving JSON response.

Usage:
	call-grpc [--metadata=<grpc-metadata>]... <grpc-endpoint> <grpc-full-method>
	call-grpc -h | --help
	call-grpc --version

Options:
	<grpc-metadata>	   metadata in the form 'key: value'.
	<grpc-endpoint>    host and port joined by ':'.
	<grpc-full-method> full method in the form 'package.name.ServiceName/MethodName'.
	-h --help          Show this screen.
	--version          Show version.`

func main() {
	args, _ := docopt.ParseArgs(doc, os.Args[1:], "1.0.0")
	code := cmd.RunCallGRPC(
		cmd.Stdio,
		args["<grpc-endpoint>"].(string),
		args["<grpc-full-method>"].(string),
	)

	os.Exit(code)
}
