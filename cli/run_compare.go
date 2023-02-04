package cli

import (
	"fmt"
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"go.uber.org/multierr"
)

func RunCompare(expectedJson string, actualJson string, verbose bool, strict bool) (code int, err error) {
	expectedJsonFile, err := os.Open(expectedJson)
	if err != nil {
		return 1, err
	}

	defer func() {
		err = multierr.Combine(err, expectedJsonFile.Close())
		if err != nil {
			code = 1
		}
	}()

	actualJsonFile, err := os.Open(actualJson)
	defer func() {
		err = multierr.Combine(err, actualJsonFile.Close())
		if err != nil {
			code = 1
		}
	}()

	match, detail, err := cmd.Compare(expectedJsonFile, actualJsonFile)
	if err != nil {
		return 1, err
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
