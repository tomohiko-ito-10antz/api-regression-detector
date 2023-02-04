package io_json

import (
	"bytes"
	"encoding/json"
	"regexp"
	"testing"

	"github.com/Jumpaku/api-regression-detector/lib/io_json/mock"
	"github.com/Jumpaku/api-regression-detector/test/assert"
)

func TestSaveJson_Tables(t *testing.T) {
	e := `{
		"t1": [
			{
				"a": null,
				"b": 123,
				"c": -123.45,
				"d": "abc",
				"e": false,
				"f": true
			},
			{
				"a": null,
				"b": 123,
				"c": -123.45,
				"d": "abc",
				"e": false,
				"f": true
			}
		],
		"t2": [
			{
				"x": null,
				"y": 123,
				"z": -123.45
			},
			{
				"x": null,
				"y": 123,
				"z": -123.45
			}
		]
	}`
	v := map[string][]map[string]any{
		"t1": {
			{
				"a": nil,
				"b": int64(123),
				"c": float64(-123.45),
				"d": "abc",
				"e": false,
				"f": true,
			},
			{
				"a": nil,
				"b": int64(123),
				"c": float64(-123.45),
				"d": "abc",
				"e": false,
				"f": true,
			},
		},
		"t2": {
			{
				"x": nil,
				"y": json.Number("123"),
				"z": json.Number("-123.45"),
			},
			{
				"x": nil,
				"y": json.Number("123"),
				"z": json.Number("-123.45"),
			},
		},
	}
	writer := mock.MockNamedBuffer{Buffer: bytes.NewBuffer(nil)}
	err := SaveJson(v, writer)
	assert.Equal(t, err, nil)
	a := writer.Buffer.String()
	eCompact := regexp.MustCompile(`\s`).ReplaceAllString(e, "")
	aCompact := regexp.MustCompile(`\s`).ReplaceAllString(a, "")
	assert.Equal(t, eCompact, aCompact)
}

/*
func SaveJson(jsonValue any, file NamedWriter) (err error) {
	log.Stderr("OUTPUT JSON TO %s", file.Name())
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(jsonValue); err != nil {
		return err
	}
	return nil
}
*/
