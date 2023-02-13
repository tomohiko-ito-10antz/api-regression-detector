package curl

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
	"github.com/Jumpaku/api-regression-detector/lib/rpc"
)

type caller struct {
}

func Caller() caller { return caller{} }

var _ cmd.HTTPCaller = caller{}

func (caller) CallHTTP(ctx context.Context, endpoint string, method cmd.HTTPMethod, req rpc.Request) (rpc.Response, error) {
	request, err := http.NewRequest(string(method), endpoint, bytes.NewReader([]byte("")))
	if err != nil {
		return rpc.Response{}, errors.Wrap(
			errors.Join(errors.HTTPFailure, err),
			"fail to create HTTP request %s %v %#v", endpoint, method, req)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return rpc.Response{}, errors.Wrap(
			errors.Join(errors.HTTPFailure, err),
			"fail to call HTTP API %s %v %#v", endpoint, method, req)
	}

	decoder := json.NewDecoder(response.Body)
	decoder.UseNumber()

	var resBody any
	if err := decoder.Decode(&resBody); err != nil {
		return rpc.Response{}, errors.Wrap(
			errors.Join(errors.HTTPFailure, err),
			"fail to read HTTP response %#v", endpoint, method, response)
	}

	res := rpc.Response{Header: response.Header}
	res.Body, err = wrap.FromAny(resBody)
	if err != nil {
		return rpc.Response{}, errors.Wrap(
			errors.Join(errors.HTTPFailure, err),
			"fail to read HTTP response %#v", endpoint, method, response)
	}

	return res, nil
}
