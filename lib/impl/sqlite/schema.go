package sqlite

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

func (o schemaGetter) GetSchema(ctx context.Context, tx db.Tx, tableName string) (db.Schema, error) {
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

func getColumnTypes(ctx context.Context, tx db.Tx, table string) (db.ColumnTypes, error) {
	rows, err := tx.Read(ctx, `SELECT name, type FROM pragma_table_info(?)`, []any{table})
	if err != nil {
		return nil, err
	}

	columnTypes := db.ColumnTypes{}
	for _, row := range rows {
		columnName, ok := row["name"]
		if !ok {
			return nil, fmt.Errorf("column %s not found", "name")
		}
		columnNameString, err := columnName.AsString()
		if err != nil {
			return nil, err
		}
		col := columnNameString.String
		columnType, ok := row["type"]
		if !ok {
			return nil, fmt.Errorf("column %s not found", "type")
		}
		columnTypeString, err := columnType.AsString()
		if err != nil {
			return nil, err
		}
		typ := columnTypeString.String
		startsWithAny := func(prefixes ...string) bool {
			for _, prefix := range prefixes {
				if strings.HasPrefix(strings.ToLower(typ), strings.ToLower(prefix)) {
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
		case startsWithAny("DATE", "DATETIME"):
			columnTypes[col] = db.ColumnTypeTime
		default:
			columnTypes[col] = db.ColumnTypeString
		}
	}

	return columnTypes, nil
}

func getPrimaryKeys(ctx context.Context, tx db.Tx, table string) ([]string, error) {
	rows, err := tx.Read(ctx, `SELECT name FROM pragma_table_info(?) WHERE pk > 0 ORDER BY pk`, []any{table})
	if err != nil {
		return nil, err
	}

	primaryKeys := []string{}
	for _, row := range rows {
		columnName, ok := row["name"]
		if !ok {
			return nil, fmt.Errorf("column %s not found", "name")
		}

		columnNameString, err := columnName.AsString()
		if err != nil {
			return nil, err
		}

		primaryKeys = append(primaryKeys, columnNameString.String)
	}

	return primaryKeys, nil
}
