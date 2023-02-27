package grpc

import (
	"strings"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
)

type Request struct {
	Metadata map[string][]string
	Body     *wrap.JsonValue
}

func ParseMetadata(metadata string) (string, string, error) {
	key, val, ok := strings.Cut(metadata, ":")
	if !ok {
		return "", "", errors.BadArgs.New(errors.Info{"header": metadata}.AppendTo("header must be in the form 'Key: value'"))
	}

	return strings.ToLower(key), strings.Trim(val, " \t"), nil
}
