package jsonio

import "fmt"

type Table struct {
	Rows []Row
}
type Tables map[string]Table

func (tables Tables) GetTableNames() []string {
	tableNames := []string{}
	for tableName := range tables {
		tableNames = append(tableNames, tableName)
	}
	return tableNames
}

func TableFromJson(json map[string][]map[string]any) (tables Tables, err error) {
	tables = Tables{}
	for tableName, rows := range json {
		table := Table{}
		for _, rowObj := range rows {
			row := Row{}
			for columnName, columnValue := range rowObj {
				jsonVal, err := NewJson(columnValue)
				if err != nil {
					return nil, err
				}
				row[columnName] = jsonVal
			}
			table.Rows = append(table.Rows, row)
		}
		tables[tableName] = table
	}
	return tables, nil
}

func TableToJson(tables Tables) (json map[string][]map[string]any, err error) {
	json = map[string][]map[string]any{}
	for tableName, table := range tables {
		rowArr := []map[string]any{}
		for _, row := range table.Rows {
			rowObj := map[string]any{}
			for columnName, columnValue := range row {
				switch columnValue.Type {
				case JsonTypeNull:
					rowObj[columnName] = nil
				case JsonTypeBoolean:
					rowObj[columnName], err = columnValue.ToBool()
					if err != nil {
						return nil, err
					}
				case JsonTypeNumber:
					rowObj[columnName], err = columnValue.ToInt64()
					if err != nil {
						rowObj[columnName], err = columnValue.ToFloat64()
						if err != nil {
							return nil, err
						}
					}
				case JsonTypeString:
					rowObj[columnName], err = columnValue.ToString()
					if err != nil {
						return nil, err
					}
				default:
					return nil, fmt.Errorf("unsupported value %v of type %v", columnValue, columnValue.Type)
				}
			}
			rowArr = append(rowArr, rowObj)
		}
		json[tableName] = rowArr
	}
	return json, nil
}
