package http

import (
	"bytes"
	nethttp "net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
)

type Request struct {
	Header map[string][]string
	Body   *wrap.JsonValue
}

func ParseHeader(header string) (string, string, error) {
	key, val, ok := strings.Cut(header, ":")
	if !ok {
		return "", "", errors.BadArgs.New(errors.Info{"header": header}.AppendTo("header must be in the form 'Key: value'"))
	}
	return key, strings.Trim(val, " \t"), nil
}

func (r *Request) ToHTTPRequest(endpointURL string, method Method) (*nethttp.Request, error) {
	errInfo := errors.Info{"requestBody": r.Body}
	reqBodyBytes, err := r.Body.MarshalJSON()
	if err != nil {
		return nil, errors.Wrap(
			errors.HTTPFailure.Err(err),
			errInfo.AppendTo("fail to convert request body to JSON"))
	}

	errInfo = errInfo.With("endpointURL", endpointURL).With("method", method)

	urlWithParams, err := AssignParamsToURL(endpointURL, r)
	if err != nil {
		return nil, errors.Wrap(
			errors.HTTPFailure.Err(err),
			errInfo.AppendTo("fail to assign request body parameters to URL"))
	}

	if method == MethodGet {
		endpointURL = urlWithParams.String()
		reqBodyBytes = nil
	}

	request, err := nethttp.NewRequest(string(method), endpointURL, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return nil, errors.Wrap(
			errors.HTTPFailure.Err(err),
			errInfo.AppendTo("fail to create request"))
	}

	request.Header = r.Header

	return request, nil
}

func AssignParamsToURL(templateURL string, req *Request) (*url.URL, error) {
	parsed, err := url.Parse(templateURL)
	if err != nil {
		return nil, errors.Wrap(
			errors.HTTPFailure.Err(err),
			errors.Info{"templateURL": templateURL}.AppendTo("fail to parse url"))
	}

	// Add all non-null primitive values in JSON body to queryParams
	queryParams := url.Values{}
	_ = req.Body.Walk(func(key []wrap.JsonKey, val *wrap.JsonValue) error {
		if len(key) == 0 || val.Type == wrap.JsonTypeNull {
			return nil
		}
		keyStrings := []string{}
		for _, v := range key {
			keyStrings = append(keyStrings, v.String())
		}
		keyPath := strings.Join(keyStrings, ".")
		switch val.Type {
		case wrap.JsonTypeBoolean:
			queryParams.Add(keyPath, strconv.FormatBool(val.MustBool()))
		case wrap.JsonTypeNumber:
			queryParams.Add(keyPath, string(val.MustNumber()))
		case wrap.JsonTypeString:
			queryParams.Add(keyPath, val.MustString())
		}
		return nil
	})

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
	parsed.Path = strings.Join(pathElms, "/")
	parsed.RawQuery = queryParams.Encode()

	return parsed, nil
}
