package postgres

import (
	"context"
	"strings"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

type schemaGetter struct{}

func GetSchema() schemaGetter { return schemaGetter{} }

var _ cmd.SchemaGetter = schemaGetter{}

func (o schemaGetter) GetSchema(ctx context.Context, tx db.Tx, tableName string) (db.Schema, error) {
	columnTypes, err := getColumnTypes(ctx, tx, tableName)
	if err != nil {
		return db.Schema{}, errors.Wrap(
			errors.DBFailure,
			"fail to get column types in table %s", tableName)
	}

	primaryKeys, err := getPrimaryKeys(ctx, tx, tableName)
	if err != nil {
		return db.Schema{}, errors.Wrap(
			errors.DBFailure,
			"fail to get primary keys in table %s", tableName)
	}

	return db.Schema{
		ColumnTypes: columnTypes,
		PrimaryKeys: primaryKeys,
	}, nil
}

func getColumnTypes(ctx context.Context, tx db.Tx, tableName string) (db.ColumnTypes, error) {
	stmt := `SELECT column_name, data_type FROM information_schema.columns WHERE table_name = $1`
	params := []any{tableName}

	rows, err := tx.Read(ctx, stmt, params)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.DBFailure),
			"fail to select column types in table %s (stmt=%v)", tableName, stmt, params)
	}

	columnTypes := db.ColumnTypes{}
	for _, row := range rows {
		columnName, ok := row["column_name"]
		if !ok {
			return nil, errors.Wrap(
				errors.BadKeyAccess,
				"key %s not find in table %s", "column_name", tableName)
		}

		columnNameBytes, err := columnName.AsBytes()
		if err != nil {
			return nil, errors.Wrap(
				err,
				"fail to parse value of column %s as []byte", columnName)
		}

		col := string(columnNameBytes.Bytes)

		columnType, ok := row["data_type"]
		if !ok {
			return nil, errors.Wrap(
				errors.BadKeyAccess,
				"key %s not find in table %s", "data_type", tableName)
		}

		columnTypeString, err := columnType.AsString()
		if err != nil {
			return nil, errors.Wrap(
				err,
				"fail to parse value of column %s as string", columnName)
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

func getPrimaryKeys(ctx context.Context, tx db.Tx, tableName string) ([]string, error) {
	stmt := `SELECT 
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
	ordinal_position`
	params := []any{tableName}

	rows, err := tx.Read(ctx, stmt, params)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.DBFailure),
			"fail to select primary keys in table %s (stmt=%v,params=%v)", tableName, stmt, params)
	}

	primaryKeys := []string{}
	for _, row := range rows {
		columnName, ok := row["column_name"]
		if !ok {
			return nil, errors.Wrap(
				errors.BadKeyAccess,
				"key %s not find in table %s", "column_name", tableName)
		}

		col, err := columnName.AsString()
		if err != nil {
			return nil, errors.Wrap(
				err,
				"fail to parse value of column %s as string", columnName)
		}

		primaryKeys = append(primaryKeys, col.String)
	}

	return primaryKeys, nil
}
