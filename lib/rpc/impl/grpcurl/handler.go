package grpcurl

import (
	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Handler struct {
	Formatter grpcurl.Formatter
	Response  Response
	Error     error
}

var _ grpcurl.InvocationEventHandler = (*Handler)(nil)

func (handler *Handler) OnResolveMethod(desc *desc.MethodDescriptor) {}

func (handler *Handler) OnSendHeaders(meta metadata.MD) {}

func (handler *Handler) OnReceiveHeaders(meta metadata.MD) {
	handler.Response.Headers = meta
}

func (handler *Handler) OnReceiveResponse(msg proto.Message) {
	handler.Response.Body, handler.Error = handler.Formatter(msg)
}

func (handler *Handler) OnReceiveTrailers(status *status.Status, md metadata.MD) {}
