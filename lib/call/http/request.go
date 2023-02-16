package http

import (
	"fmt"
	nethttp "net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/Jumpaku/api-regression-detector/lib/call"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
)

type Request struct {
	Header map[string][]string
	Body   *wrap.JsonValue
}

func ToHTTPRequest(rawURL string, method Method, body *wrap.JsonValue) (*nethttp.Request, error) {
	reader, err := call.ToReader(body)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to read JsonValue: %#v", body)
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to parse url: %s", rawURL)
	}

	primitives, err := call.EnumeratePrimitives(body)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to parse url: %s", rawURL)
	}

	fmt.Printf("%#v\n", primitives)

	isPathParam := regexp.MustCompile(`^\[.+\]$`)
	pathElms := strings.Split(parsed.Path, "/")
	for i, pathElm := range pathElms {
		fmt.Printf("%d: %#v\n", i, pathElm)

		if !isPathParam.MatchString(pathElm) {
			continue
		}

		pathParamName := strings.TrimPrefix(strings.TrimSuffix(pathElm, "]"), "[")
		fmt.Printf("%d: %#v\n", i, pathParamName)

		if primitive, ok := primitives[pathParamName]; ok {
			pathElms[i] = primitive.String()
			continue
		}

		for jsonPath, primitive := range primitives {
			if strings.HasSuffix(jsonPath, "."+pathParamName) {
				pathElms[i] = primitive.String()
				break
			}
		}
	}
	fmt.Printf("%#v\n", pathElms)
	parsed.Path = ""
	parsed = parsed.JoinPath(pathElms...)

	fmt.Printf("%#v\n", parsed)
	fmt.Printf("%#v\n", parsed.String())

	request, err := nethttp.NewRequest(string(method), parsed.String(), reader)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to create request: %s %v %#v", rawURL, method, body)
	}

	return request, nil
}
