package io

import (
	"encoding/json"
	"io"
)

func Load(file io.Reader) (tables JsonTables, err error) {
	decoder := json.NewDecoder(file)
	decoder.UseNumber()
	if err := decoder.Decode(&tables); err != nil {
		return nil, err
	}
	return tables, nil
}
