package mysql

import (
	"context"
	"strings"

	"github.com/Jumpaku/api-regression-detector/db"
)

func getColumns(rows db.Rows) (columns []string) {
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

func getColumnTypes(ctx context.Context, tx db.Exec, table string) (columnTypes db.ColumnTypes, err error) {
	rows, err := tx.Read(ctx, `SELECT column_name, column_type FROM information_schema.columns WHERE table_name = ?`, []any{table})
	if err != nil {
		return nil, err
	}
	columnTypes = db.ColumnTypes{}
	for _, row := range rows {
		colAny, ok := row["column_name"]
		if !ok || colAny == nil {
			return nil, err
		}
		colBytes, ok := colAny.([]byte)
		if !ok || colAny == nil {
			return nil, err
		}
		col := strings.ToLower(string(colBytes))
		typAny, ok := row["column_type"]
		if !ok || typAny == nil {
			return nil, err
		}
		typ, ok := typAny.([]byte)
		if !ok || typAny == nil {
			return nil, err
		}
		s := strings.ToUpper(string(typ))
		startsWithAny := func(prefixes ...string) bool {
			for _, prefix := range prefixes {
				if strings.HasPrefix(s, prefix) {
					return true
				}
			}
			return false
		}
		switch {
		case startsWithAny("BOOL", "TINYINT(1)"):
			columnTypes[col] = db.ColumnTypeBoolean
		case startsWithAny("INT", "TINYINT", "SMALLINT", "MEDIUMINT", "BIGINT"):
			columnTypes[col] = db.ColumnTypeInteger
		case startsWithAny("FLOAT", "DOUBLE", "REAL"):
			columnTypes[col] = db.ColumnTypeFloat
		default:
			columnTypes[col] = db.ColumnTypeString
		}
	}
	return columnTypes, nil
}
