package grpc

import (
	"bytes"
	"context"

	"github.com/Jumpaku/api-regression-detector/lib/call"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"

	statuspb "google.golang.org/genproto/googleapis/rpc/status"
)

func InvokeRPC(endpoint string, methodDescriptor protoreflect.MethodDescriptor, req Request) (*Response, error) {
	reader, err := call.ToReader(req.Body)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to read JsonValue: %#v", req.Body)
	}

	inputMessage := dynamicpb.NewMessage(methodDescriptor.Input())
	if err := (&protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reader.Bytes(), inputMessage); err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to convert json body to protobuf message: %s", methodDescriptor.Input().FullName())
	}

	outputMessage := dynamicpb.NewMessage(methodDescriptor.Output())

	errorMessage := statuspb.Status{}

	res := NewResponse()

	if err := invokeRPCImpl(endpoint, string(methodDescriptor.Parent().FullName()+"/"+protoreflect.FullName(methodDescriptor.Name())), metadata.MD(req.Header), inputMessage, (*metadata.MD)(&res.Header), outputMessage, &errorMessage); err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to invoke grpc: %s", methodDescriptor.Input().FullName())
	}

	b, err := protojson.MarshalOptions{EmitUnpopulated: true, AllowPartial: true}.Marshal(outputMessage)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to parse response as JSON: %#v", outputMessage)
	}

	res.Body, err = call.FromReader(bytes.NewBuffer(b))
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to parse response as JSON: %s", string(b))
	}

	e, err := protojson.MarshalOptions{EmitUnpopulated: true, AllowPartial: true}.Marshal(&errorMessage)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to parse response as JSON: %#v", outputMessage)
	}

	res.Error, err = call.FromReader(bytes.NewBuffer(e))
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.BadConversion),
			"fail to parse response as JSON: %s", string(b))
	}

	return res, nil
}

func invokeRPCImpl(endpoint string, fullMethod string, inputMetadata metadata.MD, inputMessage proto.Message, outputMetadata *metadata.MD, outputMessage proto.Message, errorMessage proto.Message) error {
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

		status, ok := status.FromError(err)
		if !ok {
			return errors.Wrap(
				errors.Join(err, errors.GRPCFailure),
				"fail to receive message: %s %s", endpoint, fullMethod)
		}
		fields := errorMessage.ProtoReflect().Descriptor().Fields()
		for i := 0; i < fields.Len(); i++ {
			field := fields.ByNumber(protowire.Number(1 + i))
			if status.Proto().ProtoReflect().Has(field) {
				errorMessage.ProtoReflect().Set(field, status.Proto().ProtoReflect().Get(field))
			}
		}
	}

	return nil
}
