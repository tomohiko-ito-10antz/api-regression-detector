package cmd

import (
	"context"

	libdb "github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
)

func Init(ctx context.Context,
	db libdb.DB,
	jsonTables jsonio.Tables,
	schemaGetter SchemaGetter,
	clearer RowClearer,
	creator RowCreator,
) error {
	err := db.RunTransaction(ctx, func(ctx context.Context, tx libdb.Tx) error {
		for tableName, table := range jsonTables {
			schema, err := schemaGetter.GetSchema(ctx, tx, tableName)
			if err != nil {
				return errors.Wrap(err, "fail to get schema of table %s", tableName)
			}

			err = clearer.ClearRows(ctx, tx, tableName)
			if err != nil {
				return errors.Wrap(err, "fail to clear rows in table %s", tableName)
			}

			err = creator.CreateRows(ctx, tx, tableName, schema, table.Rows)
			if err != nil {
				return errors.Wrap(err, "fail to create rows in table %s", tableName)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
