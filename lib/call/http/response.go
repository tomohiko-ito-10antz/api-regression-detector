package http

import (
	"net/http"

	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
)

type Response struct {
	Header map[string][]string
	Body   *wrap.JsonValue
	Code   int
}

func NewResponse() *Response {
	return &Response{
		Header: map[string][]string{},
		Body:   wrap.Null(),
		Code:   http.StatusOK,
	}
}
