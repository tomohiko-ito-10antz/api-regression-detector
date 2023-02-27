package cmd

import (
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/call/grpc"
	libcmd "github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/cmd/cli"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
	"google.golang.org/grpc/codes"
)

func RunCallGRPC(stdio *cli.Stdio, endpoint string, fullMethod string, metadata []string) (code int) {
	errorInfo := errors.Info{"endpoint": endpoint, "fullMethod": fullMethod, "metadata": metadata}

	metadataMap := map[string][]string{}
	for _, md := range metadata {
		key, val, err := grpc.ParseMetadata(md)
		if err != nil {
			PrintError(os.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunCallGRPC")))
			return 1
		}

		metadataMap[key] = append(metadataMap[key], val)
	}

	reqBodyAny, err := jsonio.LoadJson[any](stdio.Stdin)
	if err != nil {
		PrintError(os.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunCallGRPC")))
		return 1
	}

	reqBody, err := wrap.FromAny(reqBodyAny)
	if err != nil {
		PrintError(os.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunCallGRPC")))
		return 1
	}

	res, err := libcmd.CallGRPC(endpoint, fullMethod, &grpc.Request{Metadata: metadataMap, Body: reqBody})
	if err != nil {
		PrintError(os.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunCallGRPC")))
		return 1
	}

	if res.Status.Code() == codes.OK {
		if err := jsonio.SaveJson(wrap.ToAny(res.Body), os.Stdout); err != nil {
			PrintError(os.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunCallGRPC")))
			return 1
		}
	} else {
		if err := jsonio.SaveJson(wrap.ToAny(res.Error), os.Stdout); err != nil {
			PrintError(os.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunCallGRPC")))
			return 1
		}
	}

	return 0
}
