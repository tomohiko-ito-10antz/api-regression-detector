package cmd_test

import (
	"strings"
	"testing"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/test/assert"
)

func TestCompare_FullMatch(t *testing.T) {
	expectedJson := strings.NewReader(`{
	"x": {
		"a": "abc",
		"b": 123,
		"c": -123.45,
		"d": true,
		"e": false,
		"f": null,
		"g": [],
		"h": {}
	},
	"y": [
		"abc",
		123,
		-123.45,
		true,
		false,
		null,
		[],
		{}
	]
}`)
	actualJson := strings.NewReader(`{
	"x": {
		"a": "abc",
		"b": 123,
		"c": -123.45,
		"d": true,
		"e": false,
		"f": null,
		"g": [],
		"h": {}
	},
	"y": [
		"abc",
		123,
		-123.45,
		true,
		false,
		null,
		[],
		{}
	]
}`)
	aResult, _, err := cmd.Compare(expectedJson, actualJson)
	assert.Equal(t, err, nil)
	assert.Equal(t, aResult, cmd.CompareResult_FullMatch)
}

func TestCompare_SupersetMatch(t *testing.T) {
	expectedJson := strings.NewReader(`{
	"x": {
		"a": "abc",
		"b": 123,
		"c": -123.45,
		"d": true,
		"e": false,
		"f": null,
		"g": [],
		"h": {}
	},
	"y": [
		"abc",
		123,
		-123.45,
		true,
		false,
		null,
		[],
		{}
	]
}`)
	actualJson := strings.NewReader(`{
	"x": {
		"a": "abc",
		"b": 123,
		"c": -123.45,
		"d": true,
		"e": false,
		"f": null,
		"g": [],
		"h": {},
		"i": "extended"
	},
	"y": [
		"abc",
		123,
		-123.45,
		true,
		false,
		null,
		[],
		{},
		"extended"
	],
	"z": "extended"
}`)
	aResult, _, err := cmd.Compare(expectedJson, actualJson)
	assert.Equal(t, err, nil)
	assert.Equal(t, aResult, cmd.CompareResult_SupersetMatch)
}

func TestCompare_NoMatch(t *testing.T) {
	expectedJson := strings.NewReader(`{
	"x": {
		"a": "abc",
		"b": 123,
		"c": -123.45,
		"d": true,
		"e": false,
		"f": null,
		"g": [],
		"h": {}
	},
	"y": [
		"abc",
		123,
		-123.45,
		true,
		false,
		null,
		[],
		{}
	]
}`)
	actualJson := strings.NewReader(`{
	"x": {
		"a": "change",
		"b": 123,
		"c": -123.45,
		"d": true,
		"e": false,
		"f": null,
		"g": [],
		"h": {}
	},
	"y": [
		"abc",
		123,
		-123.45,
		"change",
		true,
		false,
		null,
		[],
		{}
	]
}`)
	aResult, _, err := cmd.Compare(expectedJson, actualJson)
	assert.Equal(t, err, nil)
	assert.Equal(t, aResult, cmd.CompareResult_NoMatch)
}

func TestCompare_Error(t *testing.T) {
	expectedJson := strings.NewReader(`{
	"x": {
		"a": "abc",
		"b": 123,
		"c": -123.45,
		"d": true,
		"e": false,
		"f": null,
		"g": [],
		"h": {}
	},
	unexpected
	"y": [
		"abc",
		123,
		-123.45,
		true,
		false,
		null,
		[],
		{}
	]
}`)
	actualJson := strings.NewReader(`{
	"x": {
		"a": "change",
		"b": 123,
		"c": -123.45,
		"d": true,
		"e": false,
		"f": null,
		"g": [],
		"h": {}
	},
	unexpected
	"y": [
		"abc",
		123,
		-123.45,
		"change",
		true,
		false,
		null,
		[],
		{}
	]
}`)
	aResult, _, err := cmd.Compare(expectedJson, actualJson)
	assert.Equal(t, err, nil)
	assert.Equal(t, aResult, cmd.CompareResult_Error)
}
