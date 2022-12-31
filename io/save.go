package io

import (
	"encoding/json"
	"io"

	"github.com/Jumpaku/api-regression-detector/db"
)

func Save(tables db.Tables, file io.Writer) (err error) {
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(tables); err != nil {
		return err
	}
	return nil
}
