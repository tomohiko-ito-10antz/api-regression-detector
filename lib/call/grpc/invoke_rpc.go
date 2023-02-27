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
			errors.BadConversion.Err(err),
			errors.Info{"req.Body": req.Body}.AppendTo("fail to convert request body to JSON"))
	}

	inputMessage := dynamicpb.NewMessage(methodDescriptor.Input())
	if err := (&protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqBodyBytes, inputMessage); err != nil {
		errorInfo := errors.Info{"reqBodyBytes": string(reqBodyBytes), "methodInputFullname": methodDescriptor.Input().FullName()}

		return nil, errors.Wrap(
			errors.BadConversion.Err(err),
			errorInfo.AppendTo("fail to convert request body from JSON to protobuf message"))
	}

	outputMessage := dynamicpb.NewMessage(methodDescriptor.Output())

	res := NewResponse()
	fullMethod := fmt.Sprintf(`%s/%s`, methodDescriptor.Parent().FullName(), methodDescriptor.Name())

	err = invokeRPCImpl(
		endpoint,
		fullMethod,
		metadata.MD(req.Metadata),
		inputMessage,
		(*metadata.MD)(&res.Header),
		outputMessage)
	if err != nil {
		errorStatus, ok := status.FromError(err)
		if !ok {
			errorInfo := errors.Info{"endpoint": endpoint, "fullMethod": fullMethod}

			return nil, errors.Wrap(
				errors.GRPCFailure.Err(err),
				errorInfo.AppendTo("fail to invoke grpc"))
		}

		errorMessageBytes, err := protojson.MarshalOptions{EmitUnpopulated: true, AllowPartial: true}.Marshal(errorStatus.Proto())
		if err != nil {
			errorInfo := errors.Info{"errorStatus": errorStatus.Proto()}

			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errorInfo.AppendTo("fail to convert error response body from protobuf message to JSON"))
		}

		if err = res.Error.UnmarshalJSON(errorMessageBytes); err != nil {
			errorInfo := errors.Info{"errorMessageBytes": string(errorMessageBytes)}

			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errorInfo.AppendTo("fail to parse error response body from JSON"))
		}

		res.Status = errorStatus
	}

	bodyMessageBytes, err := protojson.MarshalOptions{EmitUnpopulated: true, AllowPartial: true}.Marshal(outputMessage)
	if err != nil {
		errorInfo := errors.Info{"methodOutputFullname": methodDescriptor.Output().FullName()}
		return nil, errors.Wrap(
			errors.BadConversion.Err(err),
			errorInfo.AppendTo("fail to convert response body from protobuf message to JSON"))
	}

	if err = res.Body.UnmarshalJSON(bodyMessageBytes); err != nil {
		errorInfo := errors.Info{"errorMessageBytes": string(bodyMessageBytes)}

		return nil, errors.Wrap(
			errors.BadConversion.Err(err),
			errorInfo.AppendTo("fail to parse response body from JSON"))
	}

	return res, nil
}

func invokeRPCImpl(endpoint string, fullMethod string, inputMetadata metadata.MD, inputMessage proto.Message, outputMetadata *metadata.MD, outputMessage proto.Message) error {
	errInfo := errors.Info{"endpoint": endpoint}
	cc, err := grpc.Dial(
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return errors.Wrap(
			errors.GRPCFailure.Err(err),
			errInfo.AppendTo("fail to dial"))
	}

	defer cc.Close()

	ctx := metadata.NewOutgoingContext(context.Background(), inputMetadata)
	errInfo = errInfo.With("fullMethod", fullMethod)

	clientStream, err := grpc.NewClientStream(ctx, &grpc.StreamDesc{
		StreamName:    "ClientStream",
		ServerStreams: true,
		ClientStreams: true,
	}, cc, fullMethod)
	if err != nil {
		return errors.Wrap(
			errors.GRPCFailure.Err(err),
			errInfo.AppendTo("fail to create client stream"))
	}

	if err = clientStream.SendMsg(inputMessage); err != nil {
		return errors.Wrap(
			errors.GRPCFailure.Err(err),
			errInfo.With("inputMessage", inputMessage).AppendTo("fail to send message"))
	}

	*outputMetadata = metadata.MD{}
	headerMetadata, err := clientStream.Header()
	if err != nil {
		return errors.Wrap(
			errors.GRPCFailure.Err(err),
			errInfo.AppendTo("fail to receive header metadata"))
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
				errors.GRPCFailure.Err(err),
				errInfo.AppendTo("fail to receive trailer metadata"))
		}

		return err
	}

	return nil
}
