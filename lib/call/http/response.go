package http

import (
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
)

type Response struct {
	Header map[string][]string
	Body   *wrap.JsonValue
	Code   int
}
