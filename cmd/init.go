package cmd

import (
	"context"
	"database/sql"

	"github.com/Jumpaku/api-regression-detector/db"
	"github.com/Jumpaku/api-regression-detector/io"
)

type RowClearer interface {
	ClearRows(ctx context.Context, exec db.Transaction, table string) error
}
type RowCreator interface {
	CreateRows(ctx context.Context, exec db.Transaction, tableName string, table io.JsonTable) error
}

func Init(ctx context.Context, database *sql.DB, jsonTables io.JsonTables, clearer RowClearer, creator RowCreator) (err error) {
	return db.ExecuteTransaction(ctx, database, func(ctx context.Context, exec db.Transaction) error {
		for tableName, table := range jsonTables {
			err = clearer.ClearRows(ctx, exec, tableName)
			if err != nil {
				return err
			}
			err = creator.CreateRows(ctx, exec, tableName, table)
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
