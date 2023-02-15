package http

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/Jumpaku/api-regression-detector/lib/call"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
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

func CallHTTP(url string, method Method, req *call.Request) (*call.Response, error) {

	request, err := CreateRequest(url, method, req.Body)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to create request: %s %v %#v", url, method, req)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to do request: %#v", request)
	}

	res := &call.Response{Header: response.Header}
	if res.Body, err = call.FromReader(response.Body); err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to read response body as JSON: %#v", request)
	}

	return res, nil
}

func CreateRequest(rawURL string, method Method, body *wrap.JsonValue) (*http.Request, error) {
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

	primitives, err := EnumeratePrimitives(body)
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

	request, err := http.NewRequest(string(method), parsed.String(), reader)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.HTTPFailure),
			"fail to create request: %s %v %#v", rawURL, method, body)
	}

	return request, nil
}

func EnumeratePrimitives(v *wrap.JsonValue) (map[string]*wrap.JsonValue, error) {
	m := map[string]*wrap.JsonValue{}
	if err := enumeratePrimitivesImpl(v, "", m); err != nil {
		return nil, errors.Wrap(errors.BadState, "unexpected JsoType %v", v.Type)
	}

	return m, nil
}

func enumeratePrimitivesImpl(v *wrap.JsonValue, path string, m map[string]*wrap.JsonValue) error {
	switch v.Type {
	case wrap.JsonTypeNull, wrap.JsonTypeBoolean, wrap.JsonTypeNumber, wrap.JsonTypeString:
		if path == "" {
			path = "."
		}
		m[path] = v
	case wrap.JsonTypeObject:
		for k, vk := range v.Object() {
			enumeratePrimitivesImpl(vk, fmt.Sprintf(`%s.%s`, path, k), m)
		}
	case wrap.JsonTypeArray:
		for i, vi := range v.Array() {
			enumeratePrimitivesImpl(vi, fmt.Sprintf(`%s.%d`, path, i), m)
		}
	default:
		return errors.Wrap(errors.BadState, "unexpected JsoType %v", v.Type)
	}

	return nil
}
