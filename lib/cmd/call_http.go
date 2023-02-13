package cmd

import (
	"context"
	"net/http"

	"github.com/Jumpaku/api-regression-detector/lib/rpc"
)

type HTTPCaller interface {
	CallHTTP(ctx context.Context, endpoint string, method HTTPMethod, req rpc.Request) (rpc.Response, error)
}

type HTTPMethod string

const (
	HTTPMethodGet     HTTPMethod = http.MethodGet
	HTTPMethodPost    HTTPMethod = http.MethodPost
	HTTPMethodPut     HTTPMethod = http.MethodPut
	HTTPMethodPatch   HTTPMethod = http.MethodPatch
	HTTPMethodDelete  HTTPMethod = http.MethodDelete
	HTTPMethodHead    HTTPMethod = http.MethodHead
	HTTPMethodConnect HTTPMethod = http.MethodConnect
	HTTPMethodOptions HTTPMethod = http.MethodOptions
	HTTPMethodTrace   HTTPMethod = http.MethodTrace
)

func CallHTTP(
	ctx context.Context,
	req rpc.Request,
	endpoint string,
	method HTTPMethod,
	caller HTTPCaller,
) (rpc.Response, error) {
	return caller.CallHTTP(ctx, endpoint, method, req)
}
