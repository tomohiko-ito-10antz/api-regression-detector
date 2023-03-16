package mysql

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

	referencedTables, err := getReferencedTables(ctx, tx, tableName)
	if err != nil {
		return db.Schema{}, errors.Wrap(
			errors.DBFailure.Err(err),
			errInfo.AppendTo("fail to get referenced table of table"))
	}

	return db.Schema{
		PrimaryKeys: primaryKeys,
		ColumnTypes: columnTypes,
		References:  referencedTables,
	}, nil
}

func getReferencedTables(ctx context.Context, tx db.Tx, tableName string) (referencedTables []string, err error) {
	stmt := `SELECT
    referenced_table_name AS referenced_table_name
FROM information_schema.KEY_COLUMN_USAGE
WHERE table_name = ?
    AND referenced_table_name IS NOT NULL`
	params := []any{tableName}
	errInfo := errors.Info{"stmt": stmt, "params": params}

	rows, err := tx.Read(ctx, stmt, params)
	if err != nil {
		return nil, errors.Wrap(
			errors.DBFailure.Err(err),
			errInfo.AppendTo("fail to select referenced tables of table"))
	}

	referencedTables = []string{}
	for _, row := range rows {
		referencedTable, ok := row["referenced_table_name"]
		if !ok {
			return nil, errors.BadKeyAccess.New(
				errInfo.AppendTo("key referenced_table_name not found"))
		}

		referencedTableBytes, err := referencedTable.AsBytes()
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.With("referencedTable", referencedTable).AppendTo("fail to referenced_table_name to []byte"))
		}

		ref := string(referencedTableBytes.Bytes)

		referencedTables = append(referencedTables, ref)
	}

	return referencedTables, nil
}

func getColumnTypes(ctx context.Context, tx db.Tx, tableName string) (columnTypes db.ColumnTypes, err error) {
	stmt := `SELECT
	column_name AS column_name,
	column_type AS column_type
FROM information_schema.columns
WHERE table_name = ?`
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

		columnNameBytes, err := columnName.AsBytes()
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.With("columnName", columnName).AppendTo("fail to parse column_name to []byte"))
		}

		col := string(columnNameBytes.Bytes)

		columnType, ok := row["column_type"]
		if !ok {
			return nil, errors.BadKeyAccess.New(
				errInfo.AppendTo("key column_type not found"))
		}

		columnTypeBytes, err := columnType.AsBytes()
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.With("columnType", columnType).AppendTo("fail to parse column_type to []byte"))
		}

		typ := string(columnTypeBytes.Bytes)
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

func getPrimaryKeys(ctx context.Context, tx db.Tx, tableName string) ([]string, error) {
	stmt := `SELECT 
    column_name AS column_name
FROM 
    information_schema.key_column_usage AS key_columns 
    JOIN information_schema.table_constraints AS constraints 
    ON key_columns.constraint_name = constraints.constraint_name 
        AND key_columns.table_name = constraints.table_name
WHERE
    key_columns.table_name = ?
    AND constraint_type = 'PRIMARY KEY'
ORDER BY
    ordinal_position`
	params := []any{tableName}
	errInfo := errors.Info{"stmt": stmt, "params": params}

	table, err := tx.Read(ctx, stmt, params)
	if err != nil {
		return nil, errors.Wrap(
			errors.DBFailure.Err(err),
			errInfo.AppendTo("fail to select primary keys in table"))
	}

	primaryKeys := []string{}
	for _, row := range table {
		columnName, ok := row["column_name"]
		if !ok {
			return nil, errors.BadKeyAccess.New(
				errInfo.AppendTo("key column_name not found"))
		}

		columnNameBytes, err := columnName.AsBytes()
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.With("columnName", columnName).AppendTo("fail to parse column_name to []byte"))
		}

		primaryKeys = append(primaryKeys, string(columnNameBytes.Bytes))
	}

	return primaryKeys, nil
}
