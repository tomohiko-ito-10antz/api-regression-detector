package main

import (
	"fmt"
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/log"
	"github.com/docopt/docopt-go"
)

const doc = `Regression detector compare.
compare compares two JSON files.

Usage:
	program compare [--verbose] [--strict] <expected-json> <actual-json>
	program -h | --help
	program --version

Options:
	-h --help          Show this screen.
	--version          Show version.
	--verbose          Show verbose difference. [default: false]
	--strict           Disallow superset match. [default: false]`

func main() {
	args, _ := docopt.ParseArgs(doc, os.Args[1:], "1.0.0")
	code, err := RunCompare(
		args["<expected-json>"].(string),
		args["<actual-json>"].(string),
		args["--verbose"].(bool),
		args["--strict"].(bool))
	if err != nil {
		log.Stderr("Error\n%+v", err)
	}
	os.Exit(code)
}

func RunCompare(expectedJson string, actualJson string, verbose bool, strict bool) (code int, err error) {
	expectedJsonFile, err := os.Open(expectedJson)
	if err != nil {
		return 1, errors.Wrap(errors.Join(err, errors.IOFailure), "fail to open %s", expectedJson)
	}

	defer func() {
		if errs := errors.Join(err, expectedJsonFile.Close()); err != nil {
			err = errors.Wrap(errors.Join(errs, errors.IOFailure), "fail RunCompare")
			code = 1
		}
	}()

	actualJsonFile, err := os.Open(actualJson)
	if err != nil {
		return 1, errors.Wrap(errors.Join(err, errors.IOFailure), "fail to open %s", actualJson)
	}

	defer func() {
		if errs := errors.Join(err, actualJsonFile.Close()); err != nil {
			err = errors.Wrap(errors.Join(errs, errors.IOFailure), "fail RunCompare")
			code = 1
		}
	}()

	match, detail, err := cmd.Compare(expectedJsonFile, actualJsonFile)
	if err != nil {
		return 1, errors.Wrap(err, "fail RunCompare %s", detail)
	}

	fmt.Println(match)

	if verbose {
		fmt.Println(detail)
	}

	switch match {
	case cmd.CompareResultFullMatch:
		return 0, nil
	case cmd.CompareResultSupersetMatch:
		if strict {
			return 1, nil
		}
		return 0, nil
	default:
		return 1, nil
	}
}
