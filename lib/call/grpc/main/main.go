package main

import (
	"fmt"
	"log"

	"github.com/Jumpaku/api-regression-detector/lib/call"
	"github.com/Jumpaku/api-regression-detector/lib/call/grpc"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
)

func callSayHello() {
	b, _ := wrap.FromAny(map[string]any{"name": "My-Name", "title": "Dr."})
	req := &grpc.Request{Body: b}
	res, err := grpc.CallGRPC("localhost:50051", "api.GreetingService/SayHello", req)
	if err != nil {
		log.Fatalf("fail to call GRPC, %+v", err)
	}

	ab, _ := call.ToAny(res.Body)
	ae, _ := call.ToAny(res.Error)
	fmt.Printf("header %#v\n", res.Header)
	fmt.Printf("body   %#v\n", ab)
	fmt.Printf("error  %#v\n", ae)
}

func callGetError() {
	b, _ := wrap.FromAny(map[string]any{})
	req := &grpc.Request{Body: b}
	res, err := grpc.CallGRPC("localhost:50051", "api.GreetingService/GetError", req)
	if err != nil {
		log.Fatalf("fail to call GRPC, %+v", err)
	}

	ab, _ := call.ToAny(res.Body)
	ae, _ := call.ToAny(res.Error)
	fmt.Printf("header %#v\n", res.Header)
	fmt.Printf("body   %#v\n", ab)
	fmt.Printf("error  %#v\n", ae)
}

func main() {
	//callSayHello()
	callGetError()
}
