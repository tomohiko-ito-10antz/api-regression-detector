package http

import (
	"net/http"

	"github.com/Jumpaku/api-regression-detector/lib/call"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

type Method string

const (
	MethodGet     Method = http.MethodGet
	MethodHead    Method = http.MethodHead
	MethodPost    Method = http.MethodPost
	MethodPut     Method = http.MethodPut
	MethodPatch   Method = http.MethodPatch
	MethodDelete  Method = http.MethodDelete
	MethodConnect Method = http.MethodConnect
	MethodOptions Method = http.MethodOptions
	MethodTrace   Method = http.MethodTrace
)

func CallHTTP(endpointURL string, method Method, req *Request) (*Response, error) {

	request, err := req.ToHTTPRequest(endpointURL, method)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to create request: %s %v %#v", endpointURL, method, req)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to do request: %#v", request)
	}

	res := &Response{Header: response.Header, Code: response.StatusCode}
	if res.Body, err = call.FromReader(response.Body); err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to read response body as JSON: %#v", request)
	}

	return res, nil
}
