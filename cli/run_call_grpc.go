package cli

import (
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/call/grpc"
	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
	"google.golang.org/grpc/codes"
)

func RunCallGRPC(endpoint string, fullMethod string /*, configJson string*/) (code int, err error) {
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
		return 1, errors.Wrap(err, "fail RunCallGRPC")
	}
	reqBody, err := wrap.FromAny(reqBodyAny)
	if err != nil {
		return 1, errors.Wrap(err, "fail RunCallGRPC")
	}

	res, err := cmd.CallGRPC(endpoint, fullMethod, &grpc.Request{Body: reqBody})
	if err != nil {
		return 1, errors.Wrap(err, "fail RunCallGRPC")
	}

	if res.Status.Code() == codes.OK {
		if err := jsonio.SaveJson(wrap.ToAny(res.Body), os.Stdout); err != nil {
			return 1, errors.Wrap(err, "fail RunCallGRPC")
		}
	} else {
		if err := jsonio.SaveJson(wrap.ToAny(res.Error), os.Stdout); err != nil {
			return 1, errors.Wrap(err, "fail RunCallGRPC")
		}
	}

	return 0, nil
}
