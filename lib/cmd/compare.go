package cmd

import (
	"io"
	"strings"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/nsf/jsondiff"
)

type CompareResult string

const (
	CompareResultSupersetMatch CompareResult = "SupersetMatch"
	CompareResultFullMatch     CompareResult = "FullMatch"
	CompareResultNoMatch       CompareResult = "NoMatch"
	CompareResultError         CompareResult = "Error"
)

func Compare(expectedJson io.Reader, actualJson io.Reader) (CompareResult, string, error) {
	expected, err := io.ReadAll(expectedJson)
	if err != nil {
		return CompareResultError, "", errors.Wrap(errors.IOFailure, "fail to read first JSON")
	}

	actual, err := io.ReadAll(actualJson)
	if err != nil {
		return CompareResultError, "", errors.Wrap(errors.IOFailure, "fail to read second JSON")
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
	_, description := jsondiff.Compare(expected, actual, &opt)

	switch match {
	case jsondiff.FullMatch:
		return CompareResultFullMatch, describeDiff(description), nil
	case jsondiff.SupersetMatch:
		return CompareResultSupersetMatch, describeDiff(description), nil
	case jsondiff.NoMatch:
		return CompareResultNoMatch, describeDiff(description), nil
	case jsondiff.BothArgsAreInvalidJson, jsondiff.FirstArgIsInvalidJson, jsondiff.SecondArgIsInvalidJson:
		return CompareResultError, "", errors.Wrap(errors.BadJSON, description)
	default:
		return CompareResultError, "", errors.Wrap(errors.BadState, "unexpected case %v", match)
	}
}

func describeDiff(diff string) string {
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

		if prefix := ((trim + "  ")[:2]); prefix == "@+" {
			line = addBegin + "+|" + strings.Replace(line, prefix, "", 1) + addEnd
		} else if prefix := ((trim + "  ")[:2]); prefix == "@-" {
			line = removeBegin + "-|" + strings.Replace(line, prefix, "", 1) + removeEnd
		} else if suffix2 := (("  " + trim)[len("  "+trim)-2:]); suffix2 == "@~" {
			line = changeBegin + "~|" + line[:len(line)-2] + changeEnd
		} else if suffix3 := (("   " + trim)[len("   "+trim)-3:]); suffix3 == "@~," {
			line = changeBegin + "~|" + line[:len(line)-3] + "," + changeEnd
		} else {
			line = " |" + line
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
