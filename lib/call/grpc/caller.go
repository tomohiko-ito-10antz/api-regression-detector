package grpc

import (
	"bytes"

	"github.com/Jumpaku/api-regression-detector/lib/call"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/dynamicpb"
)

func CallGRPC(endpoint string, fullMethod string, req *call.Request) (*call.Response, error) {
	registry, err := InvokeServerReflection(endpoint, fullMethod)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.GRPCFailure),
			"fail to resolve reflection registry: %s %s", endpoint, fullMethod)
	}

	methodDescriptor, err := registry.FindMethodDescriptor(fullMethod)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.GRPCFailure),
			"fail to resolve reflection registry: %s %s", endpoint, fullMethod)
	}

	reader, err := call.ToReader(req.Body)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to read JsonValue: %#v", req.Body)
	}

	inputMessage := dynamicpb.NewMessage(methodDescriptor.Input())
	if err := protojson.Unmarshal(reader.Bytes(), inputMessage); err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to convert json body to protobuf message: %s %s", endpoint, fullMethod)
	}

	outputMessage := dynamicpb.NewMessage(methodDescriptor.Output())

	if err := InvokeRPC(endpoint, fullMethod, inputMessage, outputMessage); err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to invoke GRPC call: %s %s", endpoint, fullMethod)
	}

	b, err := protojson.Marshal(outputMessage)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to parse response as JSON: %#v", outputMessage)
	}

	res := call.Response{}

	res.Body, err = call.FromReader(bytes.NewBuffer(b))
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to parse response as JSON: %s", string(b))
	}

	return &res, nil
}
