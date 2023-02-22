package main

import (
	"fmt"
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/call/http"
	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
	"github.com/docopt/docopt-go"
)

const doc = `Regression detector call-http.
call-http calls HTTP API: sending JSON request and receiving JSON response.

Usage:
	call-http <endpoint-url> <http-method>
	call-http -h | --help
	call-http --version

Options:
	<endpoint-url>     The URL of the HTTP endpoint which may has path parameters enclosed in '[' and ']'.
	<http-method>      One of GET, HEAD, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, or PATCH.
	-h --help          Show this screen.
	--version          Show version.`

func main() {
	args, _ := docopt.ParseArgs(doc, os.Args[1:], "1.0.0")
	code, err := RunCallHTTP(
		args["<endpoint-url>"].(string),
		http.Method(args["<http-method>"].(string)),
	)
	if err != nil {
		fmt.Printf("Error\n%s\n%+v", err, err)
	}
	os.Exit(code)
}

func RunCallHTTP(endpointURL string, method http.Method /*, configJson string*/) (code int, err error) {
	// configJsonFile, err := os.Open(configJson)
	//if err != nil {
	//	return 1, errors.Wrap(errors.Join(err, errors.IOFailure), "fail to open %s", configJson)
	//}
	//
	//defer func() {
	//	if errs := errors.Join(err, configJsonFile.Close()); err != nil {
	//		err = errors.Wrap(errors.Join(errs, errors.IOFailure), "fail RunCompare")
	//		code = 1
	//	}
	//}()
	errorInfo := errors.Info{"endpointURL": endpointURL, "method": method}

	reqBodyAny, err := jsonio.LoadJson[any](os.Stdin)
	if err != nil {
		return 1, errors.Wrap(err, errorInfo.AppendTo("fail RunCallHTTP"))
	}

	reqBody, err := wrap.FromAny(reqBodyAny)
	if err != nil {
		return 1, errors.Wrap(err, errorInfo.AppendTo("fail RunCallHTTP"))
	}

	res, err := cmd.CallHTTP(endpointURL, method, &http.Request{Body: reqBody})
	if err != nil {
		return 1, errors.Wrap(err, errorInfo.AppendTo("fail RunCallHTTP"))
	}

	if err := jsonio.SaveJson(wrap.ToAny(res.Body), os.Stdout); err != nil {
		return 1, errors.Wrap(err, errorInfo.AppendTo("fail RunCallHTTP"))
	}

	return 0, nil
}
