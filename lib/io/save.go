package io

import (
	"encoding/json"
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/log"
)

func SaveJson(jsonValue any, file *os.File) (err error) {
	log.Stderr("OUTPUT JSON TO %s", file.Name())
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(jsonValue); err != nil {
		return err
	}
	return nil
}
