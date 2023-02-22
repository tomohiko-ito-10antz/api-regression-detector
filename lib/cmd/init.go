package cmd

import (
	"context"

	libdb "github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/tables"
)

func Init(ctx context.Context,
	db libdb.DB,
	jsonTables tables.InitTables,
	schemaGetter SchemaGetter,
	rowClearer RowClearer,
	rowCreator RowCreator,
) error {
	err := db.RunTransaction(ctx, func(ctx context.Context, tx libdb.Tx) error {
		tableSchema := map[string]libdb.Schema{}
		for _, table := range jsonTables {
			if _, ok := tableSchema[table.Name]; ok {
				continue
			}

			schema, err := schemaGetter.GetSchema(ctx, tx, table.Name)
			if err != nil {
				return errors.Wrap(
					errors.DBFailure.Err(err),
					errors.Info{"tableName": table.Name}.AppendTo("fail to get schema of table"))
			}

			tableSchema[table.Name] = schema
		}

		for tableName := range tableSchema {
			err := rowClearer.ClearRows(ctx, tx, tableName)
			if err != nil {
				return errors.Wrap(
					errors.DBFailure.Err(err),
					errors.Info{"tableName": tableName}.AppendTo("fail to clear rows in table"))
			}
		}

		for _, table := range jsonTables {
			err := rowCreator.CreateRows(ctx, tx, table.Name, tableSchema[table.Name], table.Rows)
			if err != nil {
				return errors.Wrap(
					errors.DBFailure.Err(err),
					errors.Info{"tableName": table.Name}.AppendTo("fail to create rows in table"))
			}
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(errors.DBFailure.Err(err), "fail to run transaction for Init")
	}

	return nil
}
