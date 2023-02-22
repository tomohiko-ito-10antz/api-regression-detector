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
	errInfo := errors.Info{"tableName": tableName}

	columnTypes, err := getColumnTypes(ctx, tx, tableName)
	if err != nil {
		return db.Schema{}, errors.Wrap(
			errors.DBFailure.Err(err),
			errInfo.AppendTo("fail to get column types in table"))
	}

	primaryKeys, err := getPrimaryKeys(ctx, tx, tableName)
	if err != nil {
		return db.Schema{}, errors.Wrap(
			errors.DBFailure.Err(err),
			errInfo.AppendTo("fail to get primary keys in table"))
	}

	return db.Schema{
		ColumnTypes: columnTypes,
		PrimaryKeys: primaryKeys,
	}, nil
}

func getColumnTypes(ctx context.Context, tx db.Tx, tableName string) (db.ColumnTypes, error) {
	stmt := `SELECT name, type FROM pragma_table_info(?)`
	params := []any{tableName}
	errInfo := errors.Info{"stmt": stmt, "params": params}

	rows, err := tx.Read(ctx, stmt, params)
	if err != nil {
		return nil, errors.Wrap(
			errors.DBFailure.Err(err),
			errInfo.AppendTo("fail to select column types in table"))
	}

	columnTypes := db.ColumnTypes{}
	for _, row := range rows {
		name, ok := row["name"]
		if !ok {
			return nil, errors.BadKeyAccess.New(
				errInfo.AppendTo("key name not found"))
		}

		nameString, err := name.AsString()
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.With("name", name).AppendTo("fail to parse name to string"))
		}

		col := nameString.String

		columnType, ok := row["type"]
		if !ok {
			return nil, errors.BadKeyAccess.New(
				errInfo.AppendTo("key type not found"))
		}

		columnTypeString, err := columnType.AsString()
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.With("columnType", columnType).AppendTo("fail to parse type to string"))
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
	errInfo := errors.Info{"stmt": stmt, "params": params}

	rows, err := tx.Read(ctx, stmt, params)
	if err != nil {
		return nil, errors.Wrap(
			errors.DBFailure.Err(err),
			errInfo.AppendTo("fail to select primary keys in table"))
	}

	primaryKeys := []string{}
	for _, row := range rows {
		name, ok := row["name"]
		if !ok {
			return nil, errors.BadKeyAccess.New(
				errInfo.AppendTo("key name not found"))
		}

		nameString, err := name.AsString()
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.With("name", name).AppendTo("fail to parse name to string"))
		}

		primaryKeys = append(primaryKeys, nameString.String)
	}

	return primaryKeys, nil
}
