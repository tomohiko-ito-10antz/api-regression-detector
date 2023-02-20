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

	cc, err := grpc.Dial(
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.GRPCFailure),
			"fail to dial: %s", endpoint)
	}

	defer cc.Close()

	refCtx := context.Background()
	refClient, err := grpc_reflection_v1alpha.NewServerReflectionClient(cc).ServerReflectionInfo(refCtx)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.GRPCFailure),
			"fail to create GRPC reflection client: %s", endpoint)
	}

	err = refClient.Send(&grpc_reflection_v1alpha.ServerReflectionRequest{
		Host:           endpoint,
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_FileContainingSymbol{FileContainingSymbol: service},
	})
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.GRPCFailure),
			"fail to send GRPC reflection request: %s", endpoint)
	}

	refRes, err := refClient.Recv()
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.GRPCFailure),
			"fail to receive GRPC reflection response: %s", endpoint)
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
			errors.Join(err, errors.GRPCFailure),
			"fail to create protobuf registry files: %s", endpoint)
	}

	return &ReflectionRegistry{registryFiles}, nil
}

func (r *ReflectionRegistry) FindMethodDescriptor(fullMethod string) (protoreflect.MethodDescriptor, error) {
	serviceName := getServiceName(fullMethod)
	desc, err := r.RegistryFiles.FindDescriptorByName(serviceName)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.GRPCFailure),
			"fail to resolve service: %s", serviceName)
	}

	serviceDescriptor, ok := desc.(protoreflect.ServiceDescriptor)
	if !ok {
		return nil, errors.Wrap(
			errors.Join(err, errors.GRPCFailure),
			"fail to resolve service: %s", serviceName)
	}

	methodName := getMethodName(fullMethod)
	methodDescriptor := serviceDescriptor.Methods().ByName(methodName)
	if methodDescriptor == nil {
		return nil, errors.Wrap(
			errors.GRPCFailure,
			"fail to resolve method: %s", methodName)
	}

	return methodDescriptor, nil
}
