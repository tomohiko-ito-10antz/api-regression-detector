package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/db"
)

type schemaGetter struct{}

func GetSchema() schemaGetter { return schemaGetter{} }

var _ cmd.SchemaGetter = schemaGetter{}

func (o schemaGetter) GetSchema(ctx context.Context, tx db.Tx, tableName string) (schema db.Schema, err error) {
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

func getColumnTypes(ctx context.Context, tx db.Tx, table string) (columnTypes db.ColumnTypes, err error) {
	rows, err := tx.Read(ctx, `SELECT column_name, data_type FROM information_schema.columns WHERE table_name = $1`, []any{table})
	fmt.Println("1")
	if err != nil {
		return nil, err
	}
	fmt.Println("2")
	columnTypes = db.ColumnTypes{}
	for _, row := range rows {
		fmt.Printf("row %v", row)
		columnName, ok := row.GetColumnValue("column_name")
		if !ok {
			return nil, fmt.Errorf("column %s not found", "column_name")
		}
		columnNameBytes, err := columnName.AsBytes()
		if err != nil {
			return nil, err
		}
		col := string(columnNameBytes.Bytes)

		columnType, ok := row.GetColumnValue("data_type")
		if !ok {
			return nil, fmt.Errorf("column %s not found", "data_type")
		}
		columnTypeString, err := columnType.AsString()
		if err != nil {
			return nil, err
		}
		typ := columnTypeString.String
		lowerTyp := strings.ToLower(typ)
		startsWithAny := func(prefixes ...string) bool {
			for _, prefix := range prefixes {
				if strings.HasPrefix(lowerTyp, strings.ToLower(prefix)) {
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

func getPrimaryKeys(ctx context.Context, tx db.Tx, table string) (primaryKeys []string, err error) {
	rows, err := tx.Read(ctx, `
SELECT 
    column_name
FROM 
    information_schema.key_column_usage AS keys 
    JOIN information_schema.table_constraints AS constraints 
    ON keys.constraint_name = constraints.constraint_name 
        AND keys.table_name = constraints.table_name
WHERE
    keys.table_name = $1
    AND constraint_type = 'PRIMARY KEY'
ORDER BY
    ordinal_position
`, []any{table})
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		columnName, ok := row.GetColumnValue("column_name")
		if !ok {
			return nil, fmt.Errorf("column %s not found", "column_name")
		}
		col, _ := columnName.AsString()
		primaryKeys = append(primaryKeys, col.String)
	}
	return primaryKeys, nil
}
