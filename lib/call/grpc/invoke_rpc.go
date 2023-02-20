package grpc

import (
	"context"
	"fmt"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

func InvokeRPC(endpoint string, methodDescriptor protoreflect.MethodDescriptor, req Request) (*Response, error) {
	reqBodyBytes, err := req.Body.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to read JsonValue: %#v", req.Body)
	}

	inputMessage := dynamicpb.NewMessage(methodDescriptor.Input())
	if err := (&protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqBodyBytes, inputMessage); err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to convert json body to protobuf message: %s", methodDescriptor.Input().FullName())
	}

	outputMessage := dynamicpb.NewMessage(methodDescriptor.Output())

	res := NewResponse()

	err = invokeRPCImpl(
		endpoint,
		fmt.Sprintf(`%s/%s`, methodDescriptor.Parent().FullName(), methodDescriptor.Name()),
		metadata.MD(req.Header),
		inputMessage,
		(*metadata.MD)(&res.Header),
		outputMessage)
	if err != nil {
		errorStatus, ok := status.FromError(err)
		if !ok {
			return nil, errors.Wrap(
				errors.Join(err, errors.BadConversion),
				"fail to invoke grpc: %s", methodDescriptor.Input().FullName())
		}

		errorMessageBytes, err := protojson.MarshalOptions{EmitUnpopulated: true, AllowPartial: true}.Marshal(errorStatus.Proto())
		if err != nil {
			return nil, errors.Wrap(
				errors.Join(err, errors.BadConversion),
				"fail to parse response as JSON: %#v", outputMessage)
		}

		if err = res.Error.UnmarshalJSON(errorMessageBytes); err != nil {
			return nil, errors.Wrap(
				errors.Join(err, errors.BadConversion),
				"fail to parse error as JSON: %s", string(errorMessageBytes))
		}

		res.Status = errorStatus
	}

	bodyMessageBytes, err := protojson.MarshalOptions{EmitUnpopulated: true, AllowPartial: true}.Marshal(outputMessage)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to parse response as JSON: %#v", outputMessage)
	}

	if err = res.Body.UnmarshalJSON(bodyMessageBytes); err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to parse response as JSON: %s", string(bodyMessageBytes))
	}

	return res, nil
}

func invokeRPCImpl(endpoint string, fullMethod string, inputMetadata metadata.MD, inputMessage proto.Message, outputMetadata *metadata.MD, outputMessage proto.Message) error {
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

	ctx := metadata.NewOutgoingContext(context.Background(), inputMetadata)

	clientStream, err := grpc.NewClientStream(ctx, &grpc.StreamDesc{
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

	*outputMetadata = metadata.MD{}
	headerMetadata, err := clientStream.Header()
	if err != nil {
		return errors.Wrap(
			errors.Join(err, errors.GRPCFailure),
			"fail to receive header metadata: %s %s", endpoint, fullMethod)
	}

	for k, v := range headerMetadata {
		(*outputMetadata)[k] = append((*outputMetadata)[k], v...)
	}

	if err = clientStream.RecvMsg(outputMessage); err != nil {
		for k, v := range clientStream.Trailer() {
			(*outputMetadata)[k] = append((*outputMetadata)[k], v...)
		}

		_, ok := status.FromError(err)
		if !ok {
			return errors.Wrap(
				errors.Join(err, errors.GRPCFailure),
				"fail to receive message: %s %s", endpoint, fullMethod)
		}

		return err
	}

	return nil
}
