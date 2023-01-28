package mysql

import (
	"context"
	"strings"

	"github.com/Jumpaku/api-regression-detector/db"
)

func getColumnTypes(ctx context.Context, tx db.Transaction, table string) (columnTypes db.ColumnTypes, err error) {
	rows, err := tx.Read(ctx, `SELECT column_name, column_type FROM information_schema.columns WHERE table_name = ?`, []any{table})
	if err != nil {
		return nil, err
	}
	columnTypes = db.ColumnTypes{}
	for _, row := range rows {
		col := ""
		{
			columnName, err := row.GetColumnValue("column_name")
			if err != nil {
				return nil, err
			}
			columnNameBytes, err := columnName.AsBytes()
			if err != nil {
				return nil, err
			}
			col = string(columnNameBytes.Bytes)
		}
		typ := ""
		{
			columnType, err := row.GetColumnValue("column_type")
			if err != nil {
				return nil, err
			}
			columnTypeBytes, err := columnType.AsBytes()
			if err != nil {
				return nil, err
			}
			typ = string(columnTypeBytes.Bytes)
		}
		lowerTyp := strings.ToLower(typ)
		startsWithAny := func(prefixes ...string) bool {
			for _, prefix := range prefixes {
				if strings.HasPrefix(lowerTyp, strings.ToUpper(prefix)) {
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
		case startsWithAny("DATE", "DATETIME", "TIMESTAMP"):
			columnTypes[col] = db.ColumnTypeTime
		default:
			columnTypes[col] = db.ColumnTypeString
		}
	}
	return columnTypes, nil
}

func getPrimaryKeys(ctx context.Context, tx db.Transaction, tableName string) (primaryKeys []string, err error) {
	table, err := tx.Read(ctx, `
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
`, []any{tableName})
	if err != nil {
		return nil, err
	}
	for _, row := range table {
		columnName, err := row.GetColumnValue("column_name")
		if err != nil {
			return nil, err
		}
		columnNameBytes, err := columnName.AsBytes()
		if err != nil {
			return nil, err
		}
		primaryKeys = append(primaryKeys, string(columnNameBytes.Bytes))
	}
	return primaryKeys, nil
}
