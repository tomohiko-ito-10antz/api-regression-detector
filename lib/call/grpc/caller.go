package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/Jumpaku/api-regression-detector/lib/call"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func CallGRPC(endpoint string, fullMethod string, req *call.Request) (*call.Response, error) {
	_, err := call.ToReader(req.Body)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to read JsonValue: %#v", req.Body)
	}

	_, err = resolveMethodDescriptor(endpoint, fullMethod)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.GRPCFailure),
			"fail to resolve GRPC reflection method: %s %s", endpoint, fullMethod)
	}
	/*
		in := dynamic.NewMessage(md.GetInputType())
		if err := jsonpb.Unmarshal(reader, in); err != nil {
			return nil, errors.Wrap(
				errors.Join(err, errors.BadConversion),
				"fail to resolve GRPC reflection method: %s %s", endpoint, fullMethod)
		}
	*/
	return &call.Response{Body: req.Body}, nil
}

func resolveMethodDescriptor(endpoint string, fullMethod string) (*desc.MethodDescriptor, error) {
	splitFullMethod := strings.Split(strings.TrimLeft(fullMethod, "/"), "/")
	fmt.Println("fm", fullMethod)
	fmt.Println("split", splitFullMethod)
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

	refCtx := context.Background()
	refClient, err := grpc_reflection_v1alpha.NewServerReflectionClient(cc).ServerReflectionInfo(refCtx)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.GRPCFailure),
			"fail to create GRPC reflection client: %s", endpoint)
	}
	fmt.Println(fullMethod)
	fmt.Println(service)
	err = refClient.Send(&grpc_reflection_v1alpha.ServerReflectionRequest{
		Host:           endpoint,
		MessageRequest: &grpc_reflection_v1alpha.ServerReflectionRequest_FileContainingSymbol{FileContainingSymbol: strings.ReplaceAll(fullMethod, "/", ".")}})
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
	files, err := protodesc.NewFiles(&fds)
	fmt.Printf("%#v\n", err)
	fmt.Printf("%#v\n", files.NumFiles())
	found, err := files.FindDescriptorByName(protoreflect.FullName(service))
	fmt.Printf("%#v\n", err)
	fmt.Printf("%#v\n", found)
	sd, ok := found.(protoreflect.ServiceDescriptor)
	fmt.Printf("%#v\n", ok)
	fmt.Printf("%#v\n", sd)
	//refClient := grpcreflect.NewClientV1Alpha(refCtx, grpc_reflection_v1alpha.NewServerReflectionClient(cc))
	/*
		sd, err := refClient.ResolveService(service)
		if err != nil {
			return nil, errors.Wrap(
				errors.Join(err, errors.GRPCFailure),
				"fail to get GRPC reflection files: %s %s", endpoint, fullMethod)
		}

		md := sd.FindMethodByName(method)
		if md == nil {
			return nil, errors.Wrap(errors.GRPCFailure,
				"fail to get GRPC reflection method: %s %s", endpoint, fullMethod)
		}
		return md, nil
	*/
	return nil, nil
}
