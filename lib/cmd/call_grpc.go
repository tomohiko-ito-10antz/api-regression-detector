package cmd

import (
	"github.com/Jumpaku/api-regression-detector/lib/call/grpc"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

func CallGRPC(endpoint string, fullMethod string, req *grpc.Request) (*grpc.Response, error) {
	errInfo := errors.Info{"endpoint": endpoint, "fullMethod": fullMethod}

	registry, err := grpc.InvokeServerReflection(endpoint, fullMethod)
	if err != nil {
		return nil, errors.Wrap(
			errors.GRPCFailure.Err(err),
			errInfo.AppendTo("fail to resolve GRPC reflection registry"))
	}

	methodDescriptor, err := registry.FindMethodDescriptor(fullMethod)
	if err != nil {
		return nil, errors.Wrap(
			errors.GRPCFailure.Err(err),
			errInfo.AppendTo("fail to resolve GRPC method by reflection"))
	}

	res, err := grpc.InvokeRPC(endpoint, methodDescriptor, *req)
	if err != nil {
		return nil, errors.Wrap(
			errors.GRPCFailure.Err(err),
			errInfo.With("request", req).AppendTo("fail to invoke GRPC call"))
	}

	return res, nil
}
