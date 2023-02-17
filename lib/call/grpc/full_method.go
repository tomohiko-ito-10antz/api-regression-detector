package grpc

import (
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func getServiceName(fullMethod string) protoreflect.FullName {
	splitFullMethod := strings.Split(strings.TrimLeft(fullMethod, "/"), "/")
	return protoreflect.FullName(splitFullMethod[0])
}

func getMethodName(fullMethod string) protoreflect.Name {
	splitFullMethod := strings.Split(strings.TrimLeft(fullMethod, "/"), "/")
	return protoreflect.Name(splitFullMethod[1])
}
