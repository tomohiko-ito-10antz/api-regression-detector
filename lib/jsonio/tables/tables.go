package tables

import (
	"fmt"
	"sort"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
	"golang.org/x/exp/maps"
)

type Table struct {
	Name string
	Rows []Row
}

type InitTables []Table

type DumpTables map[string][]Row

func (tables InitTables) GetTableNames() []string {
	tableNameSet := map[string]any{}
	for _, table := range tables {
		tableNameSet[table.Name] = nil
	}

	tableNames := maps.Keys(tableNameSet)

	sort.Slice(tableNames, func(i, j int) bool {
		return tableNames[i] < tableNames[j]
	})

	return tableNames
}

func (tables DumpTables) GetTableNames() []string {
	tableNames := maps.Keys(tables)

	sort.Slice(tableNames, func(i, j int) bool {
		return tableNames[i] < tableNames[j]
	})

	return tableNames
}

func LoadInitTables(file jsonio.NamedReader) (InitTables, error) {
	json, err := jsonio.LoadJson[[]struct {
		Name string           `json:"name"`
		Rows []map[string]any `json:"rows"`
	}](file)
	if err != nil {
		return nil, errors.Wrap(err, "fail to load tables in %s", file.Name())
	}
	fmt.Printf("%#v\n", json)
	tables := []Table{}
	for _, jsonTable := range json {
		rows := []Row{}
		for _, jsonRow := range jsonTable.Rows {
			row := Row{}
			for columnName, jsonColumnValue := range jsonRow {
				row[columnName], _ = wrap.FromAny(jsonColumnValue)
			}

			rows = append(rows, row)
		}

		tables = append(tables, Table{Name: jsonTable.Name, Rows: rows})
	}

	return tables, nil
}

func SaveDumpTables(tables DumpTables, file jsonio.NamedWriter) (err error) {
	json := map[string][]map[string]any{}
	for tableName, rows := range tables {
		rowArr := []map[string]any{}
		for _, row := range rows {
			rowObj := map[string]any{}
			for columnName, columnValue := range row {
				if columnValue == nil {
					rowObj[columnName] = nil
					continue
				}
				switch columnValue.Type {
				case wrap.JsonTypeNull:
					rowObj[columnName] = nil
				case wrap.JsonTypeBoolean:
					rowObj[columnName] = columnValue.Bool()
				case wrap.JsonTypeNumber:
					var ok bool

					rowObj[columnName], ok = columnValue.Int64()
					if !ok {
						rowObj[columnName], ok = columnValue.Float64()
						if !ok {
							return errors.Wrap(
								errors.BadConversion,
								"fail to parse %v column value %v:%T of %s as JSON primitive", columnValue.Type, columnValue, columnValue, columnName)
						}
					}
				case wrap.JsonTypeString:
					rowObj[columnName] = columnValue.String()
				default:
					return errors.Wrap(
						errors.Join(err, errors.BadConversion, errors.Unsupported),
						"unsupported conversion for column %s of type %v to JSON primitive", columnValue, columnValue.Type)
				}
			}

			rowArr = append(rowArr, rowObj)
		}

		json[tableName] = rowArr
	}

	if err := jsonio.SaveJson(json, file); err != nil {
		return errors.Wrap(
			errors.Join(err, errors.BadArgs),
			"fail to save tables %v:%T as JSON to %s", json, json, file.Name())
	}

	return nil
}
