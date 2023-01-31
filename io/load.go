package io

import (
	"encoding/json"
	"os"

	"github.com/Jumpaku/api-regression-detector/log"
)

func Load(file *os.File) (tables Tables, err error) {
	decoder := json.NewDecoder(file)
	decoder.UseNumber()
	var jsonTables map[string][]map[string]any
	if err := decoder.Decode(&jsonTables); err != nil {
		return nil, err
	}
	log.Stderr("INPUT FROM %s\n\tcontent: %v", file.Name(), jsonTables)
	tables = Tables{}
	for tableName, jsonTable := range jsonTables {
		table := Table{}
		for _, jsonRow := range jsonTable {
			row := Row{}
			for column, jsonColumnValue := range jsonRow {
				jsonValue, err := NewJson(jsonColumnValue)
				if err != nil {
					return nil, err
				}
				row[column] = jsonValue
			}
			table.Rows = append(table.Rows, row)
		}
		tables[tableName] = table
	}
	return tables, nil
}
