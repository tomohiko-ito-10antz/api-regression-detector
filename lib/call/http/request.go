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

func (r *Request) ToHTTPRequest(endpointURL string, method Method) (*nethttp.Request, error) {
	reqBodyBytes, err := wrap.Encode(r.Body)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to read JsonValue: %#v", r.Body)
	}

	if method == MethodGet {
		urlWithParams, err := AssignParamsToURL(endpointURL, r)
		if err != nil {
			return nil, errors.Wrap(
				errors.Join(err, errors.HTTPFailure),
				"fail to assign JsonValue: %#v", r.Body)
		}
		endpointURL = urlWithParams.String()
		reqBodyBytes = nil
	}

	request, err := nethttp.NewRequest(string(method), endpointURL, bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to create request: %s %v %#v", endpointURL, method, r)
	}

	request.Header = r.Header

	return request, nil
}

func AssignParamsToURL(templateURL string, req *Request) (*url.URL, error) {
	parsed, err := url.Parse(templateURL)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to parse url: %s", templateURL)
	}

	jsonPrimitiveKeys := req.Body.EnumeratePrimitiveKeys()

	// Add all primitive values in JSON body to queryParams
	queryParams := url.Values{}
	for _, key := range jsonPrimitiveKeys {
		val, ok := req.Body.Find(key...)
		if !ok || val.Type == wrap.JsonTypeNull {
			continue
		}
		keyStrings := []string{}
		for _, v := range key {
			keyStrings = append(keyStrings, v.String())
		}
		switch val.Type {
		case wrap.JsonTypeBoolean:
			queryParams.Add(strings.Join(keyStrings, "."), strconv.FormatBool(val.MustBool()))
		case wrap.JsonTypeNumber:
			queryParams.Add(strings.Join(keyStrings, "."), string(val.MustNumber()))
		case wrap.JsonTypeString:
			queryParams.Add(strings.Join(keyStrings, "."), val.MustString())
		}
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

	return parsed, nil
}
