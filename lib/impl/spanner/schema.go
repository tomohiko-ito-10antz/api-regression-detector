package spanner

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
	rows, err := tx.Read(ctx, `SELECT column_name, spanner_type FROM information_schema.columns WHERE table_name = ?`, []any{table})
	if err != nil {
		return nil, err
	}
	columnTypes = db.ColumnTypes{}
	for _, row := range rows {
		columnName, ok := row["column_name"]
		if !ok {
			return nil, fmt.Errorf("column %s not found", "column_name")
		}
		col, _ := columnName.AsString()
		spannerType, ok := row["spanner_type"]
		if !ok {
			return nil, fmt.Errorf("column %s not found", "spanner_type")
		}
		typ, _ := spannerType.AsString()
		startsWithAny := func(prefixes ...string) bool {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(typ.String), strings.ToLower(prefix)) {
					return true
				}
			}
			return false
		}
		switch {
		case startsWithAny("BOOL"):
			columnTypes[col.String] = db.ColumnTypeBoolean
		case startsWithAny("INT64"):
			columnTypes[col.String] = db.ColumnTypeInteger
		case startsWithAny("FLOAT64"):
			columnTypes[col.String] = db.ColumnTypeFloat
		case startsWithAny("DATE", "TIMESTAMP"):
			columnTypes[col.String] = db.ColumnTypeTime
		default:
			columnTypes[col.String] = db.ColumnTypeString
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
    keys.table_name = ?
    AND constraint_type = 'PRIMARY KEY'
ORDER BY
    ordinal_position
`, []any{table})
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		columnName, ok := row["column_name"]
		if !ok {
			return nil, fmt.Errorf("column %s not found", "column_name")
		}
		col, _ := columnName.AsString()
		primaryKeys = append(primaryKeys, col.String)
	}
	return primaryKeys, nil
}
