package io

import (
	"encoding/json"
	"io"
)

func Save(tables JsonTables, file io.Writer) (err error) {
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(tables); err != nil {
		return err
	}
	return nil
}
