package cmd

import (
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/call/http"
	libcmd "github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/cmd/cli"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
)

func RunCallHTTP(stdio *cli.Stdio, endpointURL string, method http.Method /*, configJson string*/) (code int) {
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
		PrintError(os.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunCallHTTP")))
		return 1
	}

	reqBody, err := wrap.FromAny(reqBodyAny)
	if err != nil {
		PrintError(os.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunCallHTTP")))
		return 1
	}

	res, err := libcmd.CallHTTP(endpointURL, method, &http.Request{Body: reqBody})
	if err != nil {
		PrintError(os.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunCallHTTP")))
		return 1
	}

	if err := jsonio.SaveJson(wrap.ToAny(res.Body), os.Stdout); err != nil {
		PrintError(os.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunCallHTTP")))
		return 1
	}

	return 0
}
