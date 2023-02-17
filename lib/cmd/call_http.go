package cmd

import (
	"io"
	"net/http"

	libhttp "github.com/Jumpaku/api-regression-detector/lib/call/http"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
)

func CallHTTP(endpointURL string, method libhttp.Method, req *libhttp.Request) (*libhttp.Response, error) {

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

	res := &libhttp.Response{Header: response.Header, Code: response.StatusCode}

	resBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to read response body as JSON: %#v", request.Body)
	}

	if res.Body, err = wrap.Decode(resBodyBytes); err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to read response body as JSON: %#v", request.Body)
	}

	return res, nil
}
