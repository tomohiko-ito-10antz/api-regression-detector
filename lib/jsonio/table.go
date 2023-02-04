package jsonio

import (
	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

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

func TableFromJson(json map[string][]map[string]any) (Tables, error) {
	tables := Tables{}
	for tableName, rowsArr := range json {
		rows := []Row{}
		for _, rowObj := range rowsArr {
			row := Row{}
			for columnName, columnValue := range rowObj {
				jsonVal, err := NewJson(columnValue)
				if err != nil {
					return nil, errors.Wrap(
						errors.Join(err, errors.BadConversion),
						"fail to parse JSON primitive %v:%T of column %s", columnName, columnValue, columnValue)
				}

				row[columnName] = jsonVal
			}

			rows = append(rows, row)
		}

		tables[tableName] = Table{Rows: rows}
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
						return nil, errors.Wrap(
							errors.Join(err, errors.BadConversion),
							"fail to parse %v column value %v:%T of %s as JSON primitive", columnValue.Type, columnValue, columnValue, columnName)
					}
				case JsonTypeNumber:
					rowObj[columnName], err = columnValue.ToInt64()
					if err != nil {
						rowObj[columnName], err = columnValue.ToFloat64()
						if err != nil {
							return nil, errors.Wrap(
								errors.Join(err, errors.BadConversion),
								"fail to parse %v column value %v:%T of %s as JSON primitive", columnValue.Type, columnValue, columnValue, columnName)
						}
					}
				case JsonTypeString:
					rowObj[columnName], err = columnValue.ToString()
					if err != nil {
						return nil, errors.Wrap(
							errors.Join(err, errors.BadConversion),
							"fail to parse %v column value %v:%T of %s as JSON primitive", columnValue.Type, columnValue, columnValue, columnName)
					}
				default:
					return nil, errors.Wrap(
						errors.Join(err, errors.BadConversion, errors.Unsupported),
						"unsupported conversion for column %s of type %v to JSON primitive", columnValue, columnValue.Type)
				}
			}

			rowArr = append(rowArr, rowObj)
		}

		json[tableName] = rowArr
	}

	return json, nil
}
