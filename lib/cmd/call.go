package cmd

import (
	"context"
	"net/rpc"
)

type Caller interface {
	Call(ctx context.Context, endpoint string, method string, req rpc.Request) (rpc.Response, error)
}

func Call(
	ctx context.Context,
	req rpc.Request,
	endpoint string,
	method string,
	caller Caller,
) (rpc.Response, error) {
	return caller.Call(ctx, endpoint, method, req)
}
