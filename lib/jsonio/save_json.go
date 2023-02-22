package jsonio

import (
	"encoding/json"
	"io"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/log"
)

type NamedWriter interface {
	io.Writer
	Name() string
}

func SaveJson[T any](jsonValue T, file NamedWriter) (err error) {
	log.Stderr("OUTPUT JSON TO %s", file.Name())
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(jsonValue); err != nil {
		return errors.Wrap(errors.BadJSON.Err(err), "fail to encode JSON")

	}

	return nil
}
