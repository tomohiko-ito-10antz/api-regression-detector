package sqlite

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
	rows, err := tx.Read(ctx, `SELECT name, type FROM pragma_table_info(?)`, []any{table})
	if err != nil {
		return nil, err
	}
	columnTypes = db.ColumnTypes{}
	for _, row := range rows {
		col := ""
		{
			columnName, err := row.GetColumnValue("name")
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
			columnType, err := row.GetColumnValue("type")
			if err != nil {
				return nil, err
			}
			columnTypeString, err := columnType.AsString()
			if err != nil {
				return nil, err
			}
			typ = columnTypeString.String
		}
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

func getPrimaryKeys(ctx context.Context, tx db.Transaction, table string) (primaryKeys []string, err error) {
	rows, err := tx.Read(ctx, `SELECT name FROM pragma_table_info(?) WHERE pk > 0 ORDER BY pk`, []any{table})
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		columnName, err := row.GetColumnValue("name")
		if err != nil {
			return nil, err
		}
		columnNameString, err := columnName.AsString()
		if err != nil {
			return nil, err
		}
		primaryKeys = append(primaryKeys, columnNameString.String)
	}
	return primaryKeys, nil
}
