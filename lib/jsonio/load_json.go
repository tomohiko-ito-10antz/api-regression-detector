package jsonio

import (
	"encoding/json"

	"github.com/Jumpaku/api-regression-detector/lib/log"
)

func LoadJson[T any](file NamedReader) (T, error) {
	log.Stderr("INPUT JSON FROM %s", file.Name())
	decoder := json.NewDecoder(file)
	decoder.UseNumber()

	var jsonValue T
	if err := decoder.Decode(&jsonValue); err != nil {
		return jsonValue, err
	}

	return jsonValue, nil
}
