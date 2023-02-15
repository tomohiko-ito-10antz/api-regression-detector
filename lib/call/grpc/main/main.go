package main

import (
	"fmt"
	"log"

	"github.com/Jumpaku/api-regression-detector/lib/call"
	"github.com/Jumpaku/api-regression-detector/lib/call/grpc"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
)

func main() {
	b, _ := wrap.FromAny(map[string]any{"name": "My-Name", "title": "Dr."})
	req := &call.Request{Body: b}
	res, err := grpc.CallGRPC("localhost:50051", "api.GreetingService/SayHello", req)
	if err != nil {
		log.Fatalf("fail to call GRPC, %+v", err)
	}

	a, _ := call.ToAny(res.Body)
	fmt.Printf("%#v", a)
}
