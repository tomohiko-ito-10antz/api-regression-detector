package main

import (
	"fmt"
	"log"

	"github.com/Jumpaku/api-regression-detector/lib/call/http"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
)

func callSayHello() {
	b, _ := wrap.FromAny(map[string]any{"name": "My-Name", "title": "Dr."})
	req := &http.Request{Body: b}
	res := &http.Response{Body: req.Body}
	res, err := http.CallHTTP("http://api:80/say/hello/[name]", http.MethodGet, req)
	if err != nil {
		log.Fatalf("fail to call HTTP, %+v", err)
	}

	ab := wrap.ToAny(res.Body)
	ac := res.Code
	fmt.Printf("header %#v\n", res.Header)
	fmt.Printf("body   %#v\n", ab)
	fmt.Printf("code   %#v\n", ac)
}

func callGetError() {
	b, _ := wrap.FromAny(map[string]any{"name": "My-Name", "title": "Dr."})
	req := &http.Request{Body: b}
	res := &http.Response{Body: req.Body}
	res, err := http.CallHTTP("http://api:80/error", http.MethodGet, req)
	if err != nil {
		log.Fatalf("fail to call HTTP, %+v", err)
	}

	ab := wrap.ToAny(res.Body)
	ac := res.Code
	fmt.Printf("header %#v\n", res.Header)
	fmt.Printf("body   %#v\n", ab)
	fmt.Printf("code   %#v\n", ac)
}

func main() {
	callSayHello()
	//callGetError()
}
