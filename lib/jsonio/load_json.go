package jsonio

import (
	"encoding/json"
	"io"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/log"
)

type NamedReader interface {
	io.Reader
	Name() string
}

func LoadJson[T any](file NamedReader) (T, error) {
	log.Stderr("INPUT JSON FROM %s", file.Name())
	decoder := json.NewDecoder(file)
	decoder.UseNumber()

	var jsonValue T
	if err := decoder.Decode(&jsonValue); err != nil {
		return jsonValue, errors.Wrap(errors.Join(err, errors.BadJSON), "fail to decode JSON from %s", file.Name())
	}

	return jsonValue, nil
}
