package io

import (
	"encoding/json"
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/log"
)

func LoadJson[T any](file *os.File) (jsonValue T, err error) {
	log.Stderr("INPUT JSON FROM %s", file.Name())
	decoder := json.NewDecoder(file)
	decoder.UseNumber()
	if err := decoder.Decode(&jsonValue); err != nil {
		return jsonValue, err
	}
	return jsonValue, nil
}
