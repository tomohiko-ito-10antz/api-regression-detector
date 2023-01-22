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
		col, err := row.GetString("column_name")
		if err != nil {
			return nil, err
		}
		typ, err := row.GetString("spanner_type")
		if err != nil {
			return nil, err
		}
		startsWith := func(prefix string) bool {
			return strings.HasPrefix(strings.ToUpper(typ), prefix)
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

func getPrimaryKeys(ctx context.Context, tx db.Exec, table string) (primaryKeys []string, err error) {
	rows, err := tx.Read(ctx, `
SELECT 
    column_name
FROM 
    information_schema.key_column_usage AS keys 
    JOIN information_schema.table_constraints AS constraints 
    ON keys.constraint_name = constraints.constraint_name 
        AND keys.table_name = constraints.table_name
WHERE
    keys.table_name = ?
    AND constraint_type = 'PRIMARY KEY'
ORDER BY
    ordinal_position
`, []any{table})
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		col, err := row.GetString("column_name")
		if err != nil {
			return nil, err
		}
		primaryKeys = append(primaryKeys, strings.ToLower(col))
	}
	return primaryKeys, nil
}
