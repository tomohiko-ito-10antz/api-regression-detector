package http

import (
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

func (r *Request) ToHTTPRequest(endpointURL string, method Method) (*nethttp.Request, error) {
	reader, err := call.ToReader(r.Body)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to read JsonValue: %#v", r.Body)
	}

	parsed, err := url.Parse(endpointURL)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to parse url: %s", endpointURL)
	}

	primitives, err := call.EnumeratePrimitives(r.Body)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to parse url: %s", endpointURL)
	}

	isPathParam := regexp.MustCompile(`^\[.+\]$`)
	pathElms := strings.Split(parsed.Path, "/")

	// Add all primitive values in JSON body to queryParams
	queryParams := url.Values{}
	for jsonPath, primitive := range primitives {
		queryParams.Add(strings.TrimPrefix(jsonPath, "."), primitive.String())
	}

	// Add primitive values in JSON body and remove them from queryParams if required as path params
	for i, pathElm := range pathElms {
		if !isPathParam.MatchString(pathElm) {
			continue
		}

		pathParamName := strings.TrimPrefix(strings.TrimSuffix(pathElm, "]"), "[")

		if primitive, ok := primitives[pathParamName]; ok {
			pathElms[i] = primitive.String()
			queryParams.Del(strings.TrimPrefix(pathParamName, "."))
			continue
		}

		for jsonPath, primitive := range primitives {
			if strings.HasSuffix(jsonPath, "."+pathParamName) {
				pathElms[i] = primitive.String()
				queryParams.Del(strings.TrimPrefix(jsonPath, "."))
				break
			}
		}
	}

	parsed.Path = ""
	parsed = parsed.JoinPath(pathElms...)
	parsed.RawQuery = queryParams.Encode()

	request, err := nethttp.NewRequest(string(method), parsed.String(), reader)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to create request: %s %v %#v", endpointURL, method, r)
	}

	request.Header = r.Header

	return request, nil
}
