package grpc

import (
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
	"google.golang.org/grpc/status"
)

type Response struct {
	Header map[string][]string
	Body   *wrap.JsonValue
	Error  *wrap.JsonValue
	Status *status.Status
}

func NewResponse() *Response {
	return &Response{
		Header: map[string][]string{},
		Body:   wrap.Null(),
		Error:  wrap.Null(),
		Status: nil,
	}
}
