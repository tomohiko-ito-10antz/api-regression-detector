package cmd

import (
	"strings"
	"testing"

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
	aResult, _, err := Compare(expectedJson, actualJson)
	assert.Equal(t, err, nil)
	assert.Equal(t, aResult, CompareResult_FullMatch)
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
	aResult, _, err := Compare(expectedJson, actualJson)
	assert.Equal(t, err, nil)
	assert.Equal(t, aResult, CompareResult_SupersetMatch)
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
	aResult, _, err := Compare(expectedJson, actualJson)
	assert.Equal(t, err, nil)
	assert.Equal(t, aResult, CompareResult_NoMatch)
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
	aResult, _, err := Compare(expectedJson, actualJson)
	assert.Equal(t, err, nil)
	assert.Equal(t, aResult, CompareResult_Error)
}

/*
type CompareResult string

const (
	CompareResult_SupersetMatch = "SupersetMatch"
	CompareResult_FullMatch     = "FullMatch"
	CompareResult_NoMatch       = "NoMatch"
	CompareResult_Error         = "Error"
)

func Compare(expectedJson io.Reader, actualJson io.Reader) (CompareResult, string, error) {
	expected, err := io.ReadAll(expectedJson)
	if err != nil {
		return CompareResult_Error, "", err
	}
	actual, err := io.ReadAll(actualJson)
	if err != nil {
		return CompareResult_Error, "", err
	}
	opt := jsondiff.DefaultConsoleOptions()
	opt.SkipMatches = true
	opt.Removed.Begin = "@-"
	opt.Added.Begin = "@+"
	opt.Changed.Begin = ""
	opt.Changed.End = "@~"
	// Check actual value matches or is a superset of expected value
	match, _ := jsondiff.Compare(actual, expected, &opt)
	// Describe how actual value is different from expected value
	_, diff := jsondiff.Compare(expected, actual, &opt)
	switch match {
	case jsondiff.FullMatch:
		return CompareResult_FullMatch, describe(diff), nil
	case jsondiff.SupersetMatch:
		return CompareResult_SupersetMatch, describe(diff), nil
	case jsondiff.NoMatch:
		return CompareResult_NoMatch, describe(diff), nil
	default:
		return CompareResult_Error, "", nil
	}
}

func describe(diff string) string {
	var (
		addBegin    = "\033[0;32m"
		addEnd      = "\033[0m"
		removeBegin = "\033[0;31m"
		removeEnd   = "\033[0m"
		changeBegin = "\033[0;33m"
		changeEnd   = "\033[0m"
	)
	lines := []string{}
	for _, line := range strings.Split(diff, "\n") {
		trim := strings.Trim(line, " \t\n")
		if prefix := string((trim + "  ")[:2]); prefix == "@+" {
			line = addBegin + "+|" + strings.Replace(line, prefix, "", 1) + addEnd
		} else if prefix := string((trim + "  ")[:2]); prefix == "@-" {
			line = removeBegin + "-|" + strings.Replace(line, prefix, "", 1) + removeEnd
		} else if suffix2 := string(("  " + trim)[len("  "+trim)-2:]); suffix2 == "@~" {
			line = changeBegin + "~|" + line[:len(line)-2] + changeEnd
		} else if suffix3 := string(("   " + trim)[len("   "+trim)-3:]); suffix3 == "@~," {
			line = changeBegin + "~|" + line[:len(line)-3] + "," + changeEnd
		} else {
			line = " |" + line
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}
*/
