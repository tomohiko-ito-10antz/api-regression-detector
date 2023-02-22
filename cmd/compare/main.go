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
	compare [--show-diff] [--no-superset] <expected-json> <actual-json>
	compare -h | --help
	compare --version

Options:
	<expected-json>    JSON file path of expected value.
	<actual-json>      JSON file path of actual value.
	--show-diff        Show difference. [default: false]
	--no-superset      Disallow superset match. [default: false]
	-h --help          Show this screen.
	--version          Show version.`

func main() {
	args, _ := docopt.ParseArgs(doc, os.Args[1:], "1.0.0")
	code, err := RunCompare(
		args["<expected-json>"].(string),
		args["<actual-json>"].(string),
		args["--show-diff"].(bool),
		args["--no-superset"].(bool))
	if err != nil {
		log.Stderr("Error\n%+v", err)
	}
	os.Exit(code)
}

func RunCompare(expectedJson string, actualJson string, showDiff bool, noSuperset bool) (code int, err error) {
	errorInfo := errors.Info{"expectedJson": expectedJson, "actualJson": actualJson, "showDiff": showDiff, "noSuperset": noSuperset}

	expectedJsonFile, err := os.Open(expectedJson)
	if err != nil {
		return 1, errors.Wrap(errors.IOFailure.Err(err), errors.Info{"expectedJson": expectedJson}.AppendTo("fail to open expected JSON file"))
	}

	defer func() {
		if errs := errors.Join(err, errors.IOFailure.Err(expectedJsonFile.Close())); errs != nil {
			err = errors.Wrap(errs, errorInfo.AppendTo("fail RunCompare"))
			code = 1
		}
	}()

	actualJsonFile, err := os.Open(actualJson)
	if err != nil {
		return 1, errors.Wrap(errors.IOFailure.Err(err), errors.Info{"actualJson": actualJson}.AppendTo("fail to open actual JSON file"))
	}

	defer func() {
		if errs := errors.Join(err, errors.IOFailure.Err(actualJsonFile.Close())); errs != nil {
			err = errors.Wrap(errs, errorInfo.AppendTo("fail RunCompare"))
			code = 1
		}
	}()

	match, detail, err := cmd.Compare(expectedJsonFile, actualJsonFile)
	if err != nil {
		return 1, errors.Wrap(err, errorInfo.With("detail", detail).AppendTo("fail RunCompare"))
	}

	fmt.Println(match)

	if showDiff {
		fmt.Println(detail)
	}

	switch match {
	case cmd.CompareResultFullMatch:
		return 0, nil
	case cmd.CompareResultSupersetMatch:
		if noSuperset {
			return 1, nil
		}
		return 0, nil
	default:
		return 1, nil
	}
}
