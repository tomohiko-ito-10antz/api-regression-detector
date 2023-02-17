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

	jsonPrimitiveKeys := r.Body.EnumeratePrimitiveKeys()

	// Add all primitive values in JSON body to queryParams
	queryParams := url.Values{}
	for _, key := range jsonPrimitiveKeys {
		val, ok := r.Body.Find(key...)
		if !ok || val.Type == wrap.JsonTypeNull {
			continue
		}
		keyStrings := []string{}
		for _, v := range key {
			keyStrings = append(keyStrings, v.String())
		}
		b, err := wrap.Encode(val)
		if err != nil {
			return nil, errors.Wrap(errors.Join(err, errors.BadConversion), "fail to encode JsonValue")
		}
		queryParams.Add(strings.Join(keyStrings, "."), string(b))
	}

	isPathParam := regexp.MustCompile(`^\[.+\]$`)
	pathElms := strings.Split(parsed.Path, "/")
	// Add primitive values in JSON body and remove them from queryParams if required as path params
	for i, pathElm := range pathElms {
		if !isPathParam.MatchString(pathElm) {
			continue
		}

		pathParamName := strings.TrimPrefix(strings.TrimSuffix(pathElm, "]"), "[")

		if queryParams.Has(pathParamName) {
			pathElms[i] = queryParams.Get(pathParamName)
			queryParams.Del(pathParamName)
			continue
		}

		for queryParamName := range queryParams {
			if strings.HasSuffix(queryParamName, "."+pathParamName) {
				pathElms[i] = queryParams.Get(queryParamName)
				queryParams.Del(queryParamName)
				continue
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
