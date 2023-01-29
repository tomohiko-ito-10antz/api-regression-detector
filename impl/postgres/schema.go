package postgres

import (
	"context"
	"strings"

	"github.com/Jumpaku/api-regression-detector/cmd"
	"github.com/Jumpaku/api-regression-detector/db"
)

type schemaGetter struct{}

func GetSchema() schemaGetter { return schemaGetter{} }

var _ cmd.SchemaGetter = schemaGetter{}

func (o schemaGetter) GetSchema(ctx context.Context, tx db.Transaction, tableName string) (schema db.Schema, err error) {
	columnTypes, err := getColumnTypes(ctx, tx, tableName)
	if err != nil {
		return db.Schema{}, err
	}
	primaryKeys, err := getPrimaryKeys(ctx, tx, tableName)
	if err != nil {
		return db.Schema{}, err
	}
	return db.Schema{
		ColumnTypes: columnTypes,
		PrimaryKeys: primaryKeys,
	}, nil
}

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
			columnNameString, err := columnName.AsString()
			if err != nil {
				return nil, err
			}
			col = columnNameString.String
		}
		typ := ""
		{
			columnType, err := row.GetColumnValue("data_type")
			if err != nil {
				return nil, err
			}
			columnTypeString, err := columnType.AsString()
			if err != nil {
				return nil, err
			}
			typ = columnTypeString.String
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
		case startsWithAny("bool"):
			columnTypes[col] = db.ColumnTypeBoolean
		case startsWithAny("int", "smallint", "bigint", "smallserial", "serial", "bigserial"):
			columnTypes[col] = db.ColumnTypeInteger
		case startsWithAny("float", "double", "real", "numeric"):
			columnTypes[col] = db.ColumnTypeFloat
		case startsWithAny("date", "timestamp"):
			columnTypes[col] = db.ColumnTypeTime
		default:
			columnTypes[col] = db.ColumnTypeString
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
