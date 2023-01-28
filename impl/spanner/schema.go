package spanner

import (
	"context"
	"strings"

	"github.com/Jumpaku/api-regression-detector/db"
)

func getColumnTypes(ctx context.Context, tx db.Transaction, table string) (columnTypes db.ColumnTypes, err error) {
	rows, err := tx.Read(ctx, `SELECT column_name, spanner_type FROM information_schema.columns WHERE table_name = ?`, []any{table})
	if err != nil {
		return nil, err
	}
	columnTypes = db.ColumnTypes{}
	for _, row := range rows {
		columnName, err := row.GetColumnValue("column_name")
		if err != nil {
			return nil, err
		}
		col, _ := columnName.AsString()
		spannerType, err := row.GetColumnValue("spanner_type")
		if err != nil {
			return nil, err
		}
		typ, _ := spannerType.AsString()
		startsWith := func(prefix string) bool {
			return strings.HasPrefix(strings.ToLower(typ.String), strings.ToLower(prefix))
		}
		switch {
		case startsWith("BOOL"):
			columnTypes[col.String] = db.ColumnTypeBoolean
		case startsWith("INT64"):
			columnTypes[col.String] = db.ColumnTypeInteger
		case startsWith("FLOAT64"):
			columnTypes[col.String] = db.ColumnTypeFloat
		default:
			columnTypes[col.String] = db.ColumnTypeString
		}
	}
	return columnTypes, nil
}

func getPrimaryKeys(ctx context.Context, tx db.Transaction, table string) (primaryKeys []string, err error) {
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
		columnName, err := row.GetColumnValue("column_name")
		if err != nil {
			return nil, err
		}
		col, _ := columnName.AsString()
		primaryKeys = append(primaryKeys, col.String)
	}
	return primaryKeys, nil
}
