package sqlite

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
	stmt := `SELECT name, type FROM pragma_table_info(?)`
	params := []any{tableName}

	rows, err := tx.Read(ctx, stmt, params)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.DBFailure),
			"fail to select column types in table %s (stmt=%v)", tableName, stmt, params)
	}

	columnTypes := db.ColumnTypes{}
	for _, row := range rows {
		columnName, ok := row["name"]
		if !ok {
			return nil, errors.Wrap(
				errors.BadKeyAccess,
				"key %s not find in table %s", "name", tableName)
		}

		columnNameString, err := columnName.AsString()
		if err != nil {
			return nil, errors.Wrap(
				err,
				"fail to parse value of column %s as string", columnName)
		}

		col := columnNameString.String

		columnType, ok := row["type"]
		if !ok {
			return nil, errors.Wrap(
				errors.BadKeyAccess,
				"key %s not find in table %s", "type", tableName)
		}

		columnTypeString, err := columnType.AsString()
		if err != nil {
			return nil, errors.Wrap(
				err,
				"fail to parse value of column %s as string", columnName)
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

func getPrimaryKeys(ctx context.Context, tx db.Tx, tableName string) ([]string, error) {
	stmt := `SELECT name FROM pragma_table_info(?) WHERE pk > 0 ORDER BY pk`
	params := []any{tableName}

	rows, err := tx.Read(ctx, stmt, params)
	if err != nil {
		return nil, errors.Wrap(
			errors.Join(err, errors.DBFailure),
			"fail to select primary keys in table %s (stmt=%v,params=%v)", tableName, stmt, params)
	}

	primaryKeys := []string{}
	for _, row := range rows {
		columnName, ok := row["name"]
		if !ok {
			return nil, errors.Wrap(
				errors.BadKeyAccess,
				"key %s not find in table %s", "name", tableName)
		}

		columnNameString, err := columnName.AsString()
		if err != nil {
			return nil, errors.Wrap(
				err,
				"fail to parse value of column %s as string", columnName)
		}

		primaryKeys = append(primaryKeys, columnNameString.String)
	}

	return primaryKeys, nil
}
