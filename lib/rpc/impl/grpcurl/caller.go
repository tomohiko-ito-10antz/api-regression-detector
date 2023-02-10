package grpcurl

import (
	"context"
	"net/rpc"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
)

type caller struct {
}

func Caller() caller { return caller{} }

var _ cmd.Caller = caller{}

func (caller) Call(ctx context.Context, endpoint string, method string, req rpc.Request) (rpc.Response, error) {
	return rpc.Response{}, nil
}
