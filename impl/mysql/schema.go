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
		colBytes, err := row.GetByteArray("column_name")
		if err != nil {
			return nil, err
		}
		col := string(colBytes)
		typBytes, err := row.GetByteArray("column_type")
		if err != nil {
			return nil, err
		}
		typ := strings.ToUpper(string(typBytes))
		startsWithAny := func(prefixes ...string) bool {
			for _, prefix := range prefixes {
				if strings.HasPrefix(typ, prefix) {
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

func getPrimaryKeys(ctx context.Context, tx db.Exec, table string) (primaryKeys []string, err error) {
	rows, err := tx.Read(ctx, `
SELECT 
    column_name
FROM 
    information_schema.key_column_usage AS key_columns 
    JOIN information_schema.table_constraints AS constraints 
    ON key_columns.constraint_name = constraints.constraint_name 
        AND key_columns.table_name = constraints.table_name
WHERE
    key_columns.table_name = ?
    AND constraint_type = 'PRIMARY KEY'
ORDER BY
    ordinal_position
`, []any{table})
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		col, err := row.GetByteArray("column_name")
		if err != nil {
			return nil, err
		}
		primaryKeys = append(primaryKeys, strings.ToLower(string(col)))
	}
	return primaryKeys, nil
}
