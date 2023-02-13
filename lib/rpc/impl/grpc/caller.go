package grpcurl

import (
	"context"

	"github.com/Jumpaku/api-regression-detector/lib/rpc"
)

type caller struct {
}

func GRPCCaller() caller { return caller{} }

func (caller) CallGRPC(ctx context.Context, endpoint string, method string, req rpc.Request) (rpc.Response, error) {
	return rpc.Response{}, nil
}
