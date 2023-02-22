package grpc

import (
	"context"
	"strings"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

type ReflectionRegistry struct {
	RegistryFiles *protoregistry.Files
}

func InvokeServerReflection(endpoint string, fullMethod string) (*ReflectionRegistry, error) {
	splitFullMethod := strings.Split(strings.TrimLeft(fullMethod, "/"), "/")
	service := splitFullMethod[0]
	errorInfo := errors.Info{"endpoint": endpoint}

	cc, err := grpc.Dial(
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, errors.Wrap(
			errors.GRPCFailure.Err(err),
			errorInfo.AppendTo("fail to dial"))
	}

	defer cc.Close()

	refCtx := context.Background()
	refClient, err := grpc_reflection_v1alpha.NewServerReflectionClient(cc).ServerReflectionInfo(refCtx)
	if err != nil {
		return nil, errors.Wrap(
			errors.GRPCFailure.Err(err),
			errorInfo.AppendTo("fail to create GRPC reflection client"))
	}

	errorInfo = errorInfo.With("service", service)
	err = refClient.Send(&grpc_reflection_v1alpha.ServerReflectionRequest{
		Host:           endpoint,
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_FileContainingSymbol{FileContainingSymbol: service},
	})
	if err != nil {
		return nil, errors.Wrap(
			errors.GRPCFailure.Err(err),
			errorInfo.AppendTo("fail to send GRPC reflection request"))
	}

	refRes, err := refClient.Recv()
	if err != nil {
		return nil, errors.Wrap(
			errors.GRPCFailure.Err(err),
			errorInfo.AppendTo("fail to receive GRPC reflection response"))
	}

	fds := descriptorpb.FileDescriptorSet{}
	for _, b := range refRes.GetFileDescriptorResponse().GetFileDescriptorProto() {
		fdp := descriptorpb.FileDescriptorProto{}
		proto.Unmarshal(b, &fdp)
		fds.File = append(fds.File, &fdp)
	}

	registryFiles, err := protodesc.NewFiles(&fds)
	if err != nil {
		return nil, errors.Wrap(
			errors.GRPCFailure.Err(err),
			errorInfo.AppendTo("fail to create protobuf registry files"))
	}

	return &ReflectionRegistry{registryFiles}, nil
}

func (r *ReflectionRegistry) FindMethodDescriptor(fullMethod string) (protoreflect.MethodDescriptor, error) {
	serviceName := getServiceName(fullMethod)

	errorInfo := errors.Info{"serviceName": serviceName}

	desc, err := r.RegistryFiles.FindDescriptorByName(serviceName)
	if err != nil {
		return nil, errors.Wrap(
			errors.GRPCFailure.Err(err),
			errorInfo.AppendTo("fail to resolve service"))
	}

	serviceDescriptor, ok := desc.(protoreflect.ServiceDescriptor)
	if !ok {
		return nil, errors.Wrap(
			errors.GRPCFailure.Err(err),
			errorInfo.AppendTo("fail to resolve service"))
	}

	methodName := getMethodName(fullMethod)
	errorInfo = errorInfo.With("methodName", methodName)

	methodDescriptor := serviceDescriptor.Methods().ByName(methodName)
	if methodDescriptor == nil {
		return nil, errors.GRPCFailure.New(errorInfo.AppendTo("fail to resolve method"))
	}

	return methodDescriptor, nil
}
