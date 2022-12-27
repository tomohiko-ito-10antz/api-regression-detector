package mysql

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Jumpaku/api-regression-detector/prepare"
)

func Build(tables prepare.Tables) (sql string) {
	for table, rows := range tables {
		sql += makeTruncateStatement(table)
		sql += makeInsertStatement(table, rows)
	}
	return sql
}

func getColumns(rows prepare.Rows) (columns []string) {
	columnAdded := map[string]bool{}
	for _, row := range rows {
		for column := range row {
			if _, added := columnAdded[column]; !added {
				columnAdded[column] = true
				columns = append(columns, column)
			}
		}
	}
	return columns
}
func makeTruncateStatement(table string) (sql string) {
	sql += fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY;\n", table)
	sql += fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = 1;\n", table)
	return sql
}
func makeInsertValues(columns []string, rows prepare.Rows) (sql string) {
	values := []string{}
	for _, row := range rows {
		cells := []string{}
		for _, column := range columns {
			cell, ok := row[column]
			if !ok || cell == nil {
				cells = append(cells, "NULL")
			} else {
				switch cell := cell.(type) {
				case nil:
					cells = append(cells, "NULL")
				case json.Number, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
					cells = append(cells, fmt.Sprint(cell))
				case string:
					cells = append(cells, fmt.Sprintf(`'%s'`, strings.ReplaceAll(cell, `'`, `''`)))
				case bool:
					if cell {
						cells = append(cells, "TRUE")
					} else {
						cells = append(cells, "FALSE")
					}
				default:
				}
			}
		}
		values = append(values, fmt.Sprintf("  (%s)", strings.Join(cells, ", ")))
	}
	return strings.Join(values, ",\n")
}
func makeInsertStatement(table string, rows prepare.Rows) (sql string) {
	columns := getColumns(rows)
	if len(columns) == 0 {
		return sql
	}
	sql += fmt.Sprintf("INSERT INTO %s (%s) VALUES\n", table, strings.Join(columns, ", "))
	sql += makeInsertValues(columns, rows)
	return sql + ";\n"
}
