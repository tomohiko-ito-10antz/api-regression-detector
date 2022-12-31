package io

import (
	"encoding/json"
	"os"

	"github.com/Jumpaku/api-regression-detector/db"
)

func Save(tables db.Tables, file *os.File) (err error) {
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(tables); err != nil {
		return err
	}
	return nil
}
