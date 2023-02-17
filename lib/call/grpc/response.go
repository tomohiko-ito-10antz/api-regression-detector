package grpc

import "github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"

type Response struct {
	Header map[string][]string
	Body   *wrap.JsonValue
	Error  *wrap.JsonValue
}

func NewResponse() *Response {
	return &Response{
		Header: map[string][]string{},
		Body:   wrap.Null(),
		Error:  wrap.Null(),
	}
}
