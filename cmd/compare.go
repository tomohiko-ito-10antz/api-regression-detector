package cmd

import (
	"io/ioutil"
	"os"

	"github.com/nsf/jsondiff"
)

func Compare(expect *os.File, actual *os.File) (jsondiff.Difference, string, error) {
	e, err := ioutil.ReadAll(expect)
	if err != nil {
		return 0, "", err
	}
	a, err := ioutil.ReadAll(actual)
	if err != nil {
		return 0, "", err
	}
	opt := jsondiff.DefaultConsoleOptions()
	label, desc := jsondiff.Compare(a, e, &opt)
	return label, desc, nil

}
