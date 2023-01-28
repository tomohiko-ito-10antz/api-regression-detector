package io

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Jumpaku/api-regression-detector/log"
)

func Save(tables Tables, file *os.File) (err error) {
	log.Stderr("OUTPUT TO %s\n\tcontent: %v", file.Name(), tables)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	jsonTables := map[string][]map[string]any{}
	for tableName, table := range tables {
		jsonTable := []map[string]any{}
		for _, row := range table {
			jsonRow := map[string]any{}
			for columnName, jsonValue := range row {
				switch jsonValue.Type {
				case JsonTypeNull:
					jsonRow[columnName] = nil
				case JsonTypeBoolean:
					jsonRow[columnName], err = jsonValue.ToBool()
					if err != nil {
						return err
					}
				case JsonTypeNumber:
					jsonRow[columnName], err = jsonValue.ToInt64()
					if err != nil {
						jsonRow[columnName], err = jsonValue.ToFloat64()
						if err != nil {
							return err
						}
					}
				case JsonTypeString:
					jsonRow[columnName], err = jsonValue.ToString()
					if err != nil {
						return err
					}
				default:
					return fmt.Errorf("cannot convert value of type %v to json value", jsonValue.Type)
				}
			}
		}
		jsonTables[tableName] = jsonTable
	}
	if err := encoder.Encode(jsonTables); err != nil {
		return err
	}
	return nil
}
