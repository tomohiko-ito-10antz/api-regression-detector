package jsonio

import (
	"encoding/json"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/log"
)

func SaveJson[T any](jsonValue T, file NamedWriter) (err error) {
	log.Stderr("OUTPUT JSON TO %s", file.Name())
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	if err := encoder.Encode(jsonValue); err != nil {
		return errors.Wrap(errors.Join(err, errors.BadArgs), "fail to encode %v:%T as JSON to %s", jsonValue, jsonValue, file.Name())
	}

	return nil
}
