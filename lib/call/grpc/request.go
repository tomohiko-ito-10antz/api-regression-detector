package grpc

import (
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
)

type Request struct {
	Header map[string][]string
	Body   *wrap.JsonValue
}
