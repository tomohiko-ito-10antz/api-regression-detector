package cmd

import (
	"fmt"
	"os"

	libcmd "github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/cmd/cli"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

func RunCompare(stdio *cli.Stdio, expectedJson string, actualJson string, showDiff bool, noSuperset bool) (code int) {
	errorInfo := errors.Info{"expectedJson": expectedJson, "actualJson": actualJson, "showDiff": showDiff, "noSuperset": noSuperset}

	expectedJsonFile, err := os.Open(expectedJson)
	if err != nil {
		PrintError(os.Stderr, errors.Wrap(err, errors.Info{"expectedJson": expectedJson}.AppendTo("fail to open expected JSON file")))
		return 1
	}

	defer func() {
		if errs := errors.Join(err, errors.IOFailure.Err(expectedJsonFile.Close())); errs != nil {
			PrintError(os.Stderr, errors.Wrap(errs, errorInfo.AppendTo("fail RunCompare")))
			code = 1
		}
	}()

	actualJsonFile, err := os.Open(actualJson)
	if err != nil {
		PrintError(os.Stderr, errors.Wrap(err, errors.Info{"actualJson": actualJson}.AppendTo("fail to open actual JSON file")))
		return 1
	}

	defer func() {
		if err := errors.Join(err, errors.IOFailure.Err(actualJsonFile.Close())); err != nil {
			PrintError(os.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunCompare")))
			code = 1
		}
	}()

	match, detail, err := libcmd.Compare(expectedJsonFile, actualJsonFile)
	if err != nil {
		PrintError(os.Stderr, errors.Wrap(err, errorInfo.With("detail", detail).AppendTo("fail RunCompare")))
		return 1
	}

	fmt.Println(match)

	if showDiff {
		fmt.Println(detail)
	}

	switch match {
	case libcmd.CompareResultFullMatch:
		return 0
	case libcmd.CompareResultSupersetMatch:
		if noSuperset {
			return 1
		}
		return 0
	default:
		return 1
	}
}
