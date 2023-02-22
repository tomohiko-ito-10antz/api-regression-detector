package cmd

import (
	"io"
	nethttp "net/http"

	libhttp "github.com/Jumpaku/api-regression-detector/lib/call/http"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

func CallHTTP(endpointURL string, method libhttp.Method, req *libhttp.Request) (*libhttp.Response, error) {
	request, err := req.ToHTTPRequest(endpointURL, method)
	if err != nil {
		return nil, errors.Wrap(
			errors.HTTPFailure.Err(err),
			errors.Info{"req": req}.AppendTo("fail to create request"))
	}

	errInfo := errors.Info{"endpointURL": endpointURL, "method": method, "request": request}

	response, err := nethttp.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(
			errors.HTTPFailure.Err(err),
			errInfo.AppendTo("fail to perform HTTP call"))
	}

	resBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(
			errors.HTTPFailure.Err(err),
			errInfo.AppendTo("fail to read response body as JSON"))
	}

	res := libhttp.NewResponse()
	res.Code = response.StatusCode
	res.Header = response.Header

	if err = res.Body.UnmarshalJSON(resBodyBytes); err != nil {
		return nil, errors.Wrap(
			errors.HTTPFailure.Err(err),
			errors.Info{"resBodyBytes": string(resBodyBytes)}.AppendTo("fail to parse JSON response body"))
	}

	return res, nil
}
