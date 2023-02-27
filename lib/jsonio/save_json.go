package jsonio

import (
	"encoding/json"
	"io"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

type NamedWriter interface {
	io.Writer
	Name() string
}

func SaveJson[T any](jsonValue T, file io.Writer) (err error) {
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(jsonValue); err != nil {
		return errors.Wrap(errors.BadJSON.Err(err), "fail to encode JSON")

	}

	return nil
}
