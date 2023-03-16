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
	referencedTables, err := getReferencedTables(ctx, tx, tableName)
	if err != nil {
		return db.Schema{}, errors.Wrap(
			errors.DBFailure.Err(err),
			errInfo.AppendTo("fail to get referenced table of table"))
	}
	interleavedTables, err := getParentTables(ctx, tx, tableName)
	if err != nil {
		return db.Schema{}, errors.Wrap(
			errors.DBFailure.Err(err),
			errInfo.AppendTo("fail to get interleaved table of table"))
	}

	return db.Schema{
		ColumnTypes:  columnTypes,
		PrimaryKeys:  primaryKeys,
		Dependencies: append(referencedTables, interleavedTables...),
	}, nil
}

func getReferencedTables(ctx context.Context, tx db.Tx, tableName string) (referencedTables []string, err error) {
	stmt := `SELECT
	DISTINCT ctu.table_name AS referenced_table_name
  FROM 
	information_schema.TABLE_CONSTRAINTS AS tc
	JOIN information_schema.CONSTRAINT_TABLE_USAGE AS ctu
	  ON tc.constraint_name = ctu.constraint_name
  WHERE tc.constraint_type = 'FOREIGN KEY'
	AND tc.table_name = ?`
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

		referencedTableString, err := referencedTable.AsString()
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.With("referencedTable", referencedTable).AppendTo("fail to referenced_table_name to []byte"))
		}

		ref := string(referencedTableString.String)

		referencedTables = append(referencedTables, ref)
	}

	return referencedTables, nil
}

func getParentTables(ctx context.Context, tx db.Tx, tableName string) (parentTables []string, err error) {
	stmt := `SELECT
    parent_table_name AS parent_table_name
FROM information_schema.TABLES
WHERE parent_table_name IS NOT NULL
    AND table_name = ?`
	params := []any{tableName}
	errInfo := errors.Info{"stmt": stmt, "params": params}

	rows, err := tx.Read(ctx, stmt, params)
	if err != nil {
		return nil, errors.Wrap(
			errors.DBFailure.Err(err),
			errInfo.AppendTo("fail to select referenced tables of table"))
	}

	parentTables = []string{}
	for _, row := range rows {
		parentTable, ok := row["parent_table_name"]
		if !ok {
			return nil, errors.BadKeyAccess.New(
				errInfo.AppendTo("key parent_table_name not found"))
		}

		parentTablesString, err := parentTable.AsString()
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.With("parentTables", parentTables).AppendTo("fail to parent_table_name to []byte"))
		}

		ref := string(parentTablesString.String)

		parentTables = append(parentTables, ref)
	}

	return parentTables, nil
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
