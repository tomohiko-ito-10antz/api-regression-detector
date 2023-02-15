package grpc

import (
	"context"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

func InvokeRPC(endpoint string, fullMethod string, inputMessage proto.Message, outputMessage proto.Message) error {
	cc, err := grpc.Dial(
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return errors.Wrap(
			errors.Join(err, errors.GRPCFailure),
			"fail to dial: %s", endpoint)
	}

	defer cc.Close()

	clientStream, err := grpc.NewClientStream(context.Background(), &grpc.StreamDesc{
		StreamName:    "ClientStream",
		Handler:       nil, //func(srv any, stream grpc.ServerStream) error { return nil },
		ServerStreams: true,
		ClientStreams: true,
	}, cc, fullMethod)
	if err != nil {
		return errors.Wrap(
			errors.Join(err, errors.GRPCFailure),
			"fail to create client stream: %s %s", endpoint, fullMethod)
	}

	if err = clientStream.SendMsg(inputMessage); err != nil {
		return errors.Wrap(
			errors.Join(err, errors.GRPCFailure),
			"fail to send message: %s %s", endpoint, fullMethod)
	}

	if err = clientStream.RecvMsg(outputMessage); err != nil {
		return errors.Wrap(
			errors.Join(err, errors.GRPCFailure),
			"fail to receive message: %s %s", endpoint, fullMethod)
	}

	return nil
}
