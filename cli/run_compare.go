package cli

import (
	"fmt"
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

func RunCompare(expectedJson string, actualJson string, verbose bool, strict bool) (code int, err error) {
	expectedJsonFile, err := os.Open(expectedJson)
	if err != nil {
		return 1, errors.Wrap(errors.Join(err, errors.IOFailure), "fail to open %s", expectedJson)
	}

	defer func() {
		err = errors.Wrap(errors.Join(err, expectedJsonFile.Close()), "fail RunCompare")
		if err != nil {
			code = 1
		}
	}()

	actualJsonFile, err := os.Open(actualJson)
	if err != nil {
		return 1, errors.Wrap(errors.Join(err, errors.IOFailure), "fail to open %s", actualJson)
	}

	defer func() {
		err = errors.Wrap(errors.Join(err, actualJsonFile.Close()), "fail RunCompare")
		if err != nil {
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
