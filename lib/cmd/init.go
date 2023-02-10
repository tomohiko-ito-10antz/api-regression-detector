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
		for _, table := range jsonTables {
			schema, err := schemaGetter.GetSchema(ctx, tx, table.Name)
			if err != nil {
				return errors.Wrap(err, "fail to get schema of table %s", table.Name)
			}

			err = rowClearer.ClearRows(ctx, tx, table.Name)
			if err != nil {
				return errors.Wrap(err, "fail to clear rows in table %s", table.Name)
			}

			err = rowCreator.CreateRows(ctx, tx, table.Name, schema, table.Rows)
			if err != nil {
				return errors.Wrap(err, "fail to create rows in table %s", table.Name)
			}
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "transaction for Init failed")
	}

	return nil
}
