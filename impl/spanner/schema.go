package spanner

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
				columns = append(columns, strings.ToLower(column))
			}
		}
	}
	return columns
}

func getColumnTypes(ctx context.Context, tx db.Exec, table string) (columnTypes db.ColumnTypes, err error) {
	rows, err := tx.Read(ctx, `SELECT column_name, spanner_type FROM information_schema.columns WHERE table_name = ?`, []any{table})
	if err != nil {
		return nil, err
	}
	columnTypes = db.ColumnTypes{}
	for _, row := range rows {
		colAny, ok := row["column_name"]
		if !ok || colAny == nil {
			return nil, err
		}
		col, ok := colAny.(string)
		if !ok || colAny == nil {
			return nil, err
		}
		typAny, ok := row["spanner_type"]
		if !ok || typAny == nil {
			return nil, err
		}
		typ, ok := typAny.(string)
		if !ok || typAny == nil {
			return nil, err
		}
		s := strings.ToUpper(typ)
		startsWith := func(prefix string) bool {
			return strings.HasPrefix(s, prefix)
		}
		switch {
		case startsWith("BOOL"):
			columnTypes[strings.ToLower(col)] = db.ColumnTypeBoolean
		case startsWith("INT64"):
			columnTypes[strings.ToLower(col)] = db.ColumnTypeInteger
		case startsWith("FLOAT64"):
			columnTypes[strings.ToLower(col)] = db.ColumnTypeFloat
		default:
			columnTypes[strings.ToLower(col)] = db.ColumnTypeString
		}
	}
	return columnTypes, nil
}
