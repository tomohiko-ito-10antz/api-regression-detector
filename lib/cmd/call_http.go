package cmd

import (
	"io"
	"net/http"

	libhttp "github.com/Jumpaku/api-regression-detector/lib/call/http"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
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

	resBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to read response body as JSON: %#v", request.Body)
	}

	res := libhttp.NewResponse()
	res.Code = response.StatusCode
	res.Header = response.Header

	if err = res.Body.UnmarshalJSON(resBodyBytes); err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to read response body as JSON: %#v", request.Body)
	}

	return res, nil
}
