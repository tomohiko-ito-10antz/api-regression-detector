package main

import (
	"fmt"
	"log"

	"github.com/Jumpaku/api-regression-detector/lib/call"
	"github.com/Jumpaku/api-regression-detector/lib/call/http"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
)

func main() {
	b, _ := wrap.FromAny(nil)
	req := &call.Request{Body: b}
	res := &call.Response{Body: req.Body}
	res, err := http.CallHTTP("api:80/say/hello/my-name", http.MethodGet, req)
	if err != nil {
		log.Fatalf("fail to call HTTP, %+v", err)
	}

	a, _ := call.ToAny(res.Body)
	fmt.Printf("%#v", a)
}
