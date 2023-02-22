package spanner

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

func getColumnTypes(ctx context.Context, tx db.Tx, tableName string) (columnTypes db.ColumnTypes, err error) {
	stmt := `SELECT column_name, spanner_type FROM information_schema.columns WHERE table_name = ?`
	params := []any{tableName}
	errInfo := errors.Info{"stmt": stmt, "params": params}

	rows, err := tx.Read(ctx, stmt, params)
	if err != nil {
		return nil, errors.Wrap(
			errors.DBFailure.Err(err),
			errInfo.AppendTo("fail to select column types in table"))
	}

	columnTypes = db.ColumnTypes{}
	for _, row := range rows {
		columnName, ok := row["column_name"]
		if !ok {
			return nil, errors.BadKeyAccess.New(
				errInfo.AppendTo("key column_name not found"))
		}

		col, err := columnName.AsString()
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.With("columnName", columnName).AppendTo("fail to parse column_name to string"))
		}

		spannerType, ok := row["spanner_type"]
		if !ok {
			return nil, errors.BadKeyAccess.New(
				errInfo.AppendTo("key spanner_type not found"))
		}

		typ, err := spannerType.AsString()
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.With("spannerType", spannerType).AppendTo("fail to parse spanner_type to string"))
		}

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

func getPrimaryKeys(ctx context.Context, tx db.Tx, tableName string) ([]string, error) {
	stmt := `SELECT 
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
    ordinal_position`
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
		columnName, ok := row["column_name"]
		if !ok {
			return nil, errors.BadKeyAccess.New(
				errInfo.AppendTo("key column_name not found"))
		}

		col, err := columnName.AsString()
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.With("columnName", columnName).AppendTo("fail to parse column_name to string"))
		}

		primaryKeys = append(primaryKeys, col.String)
	}
	return primaryKeys, nil
}
