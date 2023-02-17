package cmd

import (
	"github.com/Jumpaku/api-regression-detector/lib/call/grpc"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

func CallGRPC(endpoint string, fullMethod string, req *grpc.Request) (*grpc.Response, error) {
	registry, err := grpc.InvokeServerReflection(endpoint, fullMethod)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.GRPCFailure),
			"fail to resolve reflection registry: %s %s", endpoint, fullMethod)
	}

	methodDescriptor, err := registry.FindMethodDescriptor(fullMethod)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.GRPCFailure),
			"fail to resolve reflection method: %s %s", endpoint, fullMethod)
	}

	res, err := grpc.InvokeRPC(endpoint, methodDescriptor, *req)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to invoke GRPC call: %s %s", endpoint, fullMethod)
	}

	return res, nil
}
