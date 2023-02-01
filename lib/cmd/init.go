package cmd

import (
	"context"
	"database/sql"

	lib_db "github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/io"
)

func Init(ctx context.Context,
	db *sql.DB,
	jsonTables io.Tables,
	schemaGetter SchemaGetter,
	clearer RowClearer,
	creator RowCreator,
) (err error) {
	return lib_db.RunTransaction(ctx, db, func(ctx context.Context, tx lib_db.Tx) error {
		for tableName, table := range jsonTables {
			schema, err := schemaGetter.GetSchema(ctx, tx, tableName)
			if err != nil {
				return err
			}
			err = clearer.ClearRows(ctx, tx, tableName)
			if err != nil {
				return err
			}
			err = creator.CreateRows(ctx, tx, tableName, schema, table.Rows)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

/*
func convertTablesJsonToDB(jsonTables io.JsonTables) (dbTables db.Tables) {
	dbTables = db.Tables{}
	for jsonTableName, jsonRows := range jsonTables {
		dbRows := db.Table{}
		for _, jsonRow := range jsonRows {
			dbRow := db.Row{}
			for jsonColumnName, jsonColumnValue := range jsonRow {
				dbRow[jsonColumnName] = jsonColumnValue
			}
			dbRows = append(dbRows, dbRow)
		}
		dbTables[jsonTableName] = dbRows
	}
	return dbTables
}
*/
