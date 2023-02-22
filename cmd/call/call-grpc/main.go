package main

import (
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/call/grpc"
	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
	"github.com/Jumpaku/api-regression-detector/lib/log"
	"github.com/docopt/docopt-go"
	"google.golang.org/grpc/codes"
)

const doc = `Regression detector call-grpc.
call-grpc calls GRPC API: sending JSON request and receiving JSON response.

Usage:
	call-grpc <grpc-endpoint> <grpc-full-method>
	call-grpc -h | --help
	call-grpc --version

Options:
	<grpc-endpoint>    host and port joined by ':'.
	<grpc-full-method> full method in the form 'package.name.ServiceName/MethodName'.
	-h --help          Show this screen.
	--version          Show version.`

func main() {
	args, _ := docopt.ParseArgs(doc, os.Args[1:], "1.0.0")
	code, err := RunCallGRPC(
		args["<grpc-endpoint>"].(string),
		args["<grpc-full-method>"].(string),
	)
	if err != nil {
		log.Stderr("Error\n%s\n%+v", err, err)
	}
	os.Exit(code)
}

func RunCallGRPC(endpoint string, fullMethod string /*, configJson string*/) (code int, err error) {
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
	errorInfo := errors.Info{"endpoint": endpoint, "fullMethod": fullMethod}

	reqBodyAny, err := jsonio.LoadJson[any](os.Stdin)
	if err != nil {
		return 1, errors.Wrap(err, errorInfo.AppendTo("fail RunCallGRPC"))
	}

	reqBody, err := wrap.FromAny(reqBodyAny)
	if err != nil {
		return 1, errors.Wrap(err, errorInfo.AppendTo("fail RunCallGRPC"))
	}

	res, err := cmd.CallGRPC(endpoint, fullMethod, &grpc.Request{Body: reqBody})
	if err != nil {
		return 1, errors.Wrap(err, errorInfo.AppendTo("fail RunCallGRPC"))
	}

	if res.Status.Code() == codes.OK {
		if err := jsonio.SaveJson(wrap.ToAny(res.Body), os.Stdout); err != nil {
			return 1, errors.Wrap(err, errorInfo.AppendTo("fail RunCallGRPC"))
		}
	} else {
		if err := jsonio.SaveJson(wrap.ToAny(res.Error), os.Stdout); err != nil {
			return 1, errors.Wrap(err, errorInfo.AppendTo("fail RunCallGRPC"))
		}
	}

	return 0, nil
}
