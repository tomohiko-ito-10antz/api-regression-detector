package cli

import (
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/call/http"
	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
)

func RunCallHTTP(endpointURL string, method http.Method /*, configJson string*/) (code int, err error) {
	//configJsonFile, err := os.Open(configJson)
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

	reqBodyAny, err := jsonio.LoadJson[any](os.Stdin)
	if err != nil {
		return 1, errors.Wrap(err, "fail RunCallHTTP")
	}

	reqBody, err := wrap.FromAny(reqBodyAny)
	if err != nil {
		return 1, errors.Wrap(err, "fail RunCallHTTP")
	}

	res, err := cmd.CallHTTP(endpointURL, method, &http.Request{Body: reqBody})
	if err != nil {
		return 1, errors.Wrap(err, "fail RunCallHTTP")
	}

	if err := jsonio.SaveJson(wrap.ToAny(res.Body), os.Stdout); err != nil {
		return 1, errors.Wrap(err, "fail RunCallHTTP")
	}

	return 0, nil
}
