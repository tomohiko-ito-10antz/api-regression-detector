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
		ColumnTypes: columnTypes,
		PrimaryKeys: primaryKeys,
		References:  referencedTables,
	}, nil
}

func getReferencedTables(ctx context.Context, tx db.Tx, tableName string) (referencedTables []string, err error) {
	stmt := `SELECT
	DISTINCT ccu.table_name AS referenced_table_name
  FROM 
	information_schema.TABLE_CONSTRAINTS AS tc
	JOIN information_schema.CONSTRAINT_COLUMN_USAGE AS ccu
	  ON tc.constraint_name = ccu.constraint_name
  WHERE tc.constraint_type = 'FOREIGN KEY'
	AND tc.table_name = $1`
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

func getColumnTypes(ctx context.Context, tx db.Tx, tableName string) (db.ColumnTypes, error) {
	stmt := `SELECT column_name, data_type FROM information_schema.columns WHERE table_name = $1`
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

		dataType, ok := row["data_type"]
		if !ok {
			return nil, errors.BadKeyAccess.New(
				errInfo.AppendTo("key data_type not found"))
		}

		dataTypeString, err := dataType.AsString()
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.With("dataType", dataType).AppendTo("fail to parse data_type to string"))
		}

		typ := dataTypeString.String
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

		col, err := columnName.AsBytes()
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.With("columnName", columnName).AppendTo("fail to parse column_name to string"))
		}

		primaryKeys = append(primaryKeys, string(col.Bytes))
	}

	return primaryKeys, nil
}
