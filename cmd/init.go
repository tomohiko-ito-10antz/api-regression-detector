package cmd

import (
	"context"
	"database/sql"

	"github.com/Jumpaku/api-regression-detector/db"
	"github.com/Jumpaku/api-regression-detector/io"
)

type Truncate interface {
	Truncate(ctx context.Context, exec db.Exec, table string) error
}
type Insert interface {
	Insert(ctx context.Context, exec db.Exec, table string, rows db.Rows) error
}

func Init(ctx context.Context, database *sql.DB, jsonTables io.JsonTables, truncate Truncate, insert Insert) (err error) {
	tables := convertTablesJsonToDB(jsonTables)
	return db.Transaction(ctx, database, func(ctx context.Context, exec db.Exec) error {
		for table, rows := range tables {
			err = truncate.Truncate(ctx, exec, table)
			if err != nil {
				return err
			}
			err = insert.Insert(ctx, exec, table, rows)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func convertTablesJsonToDB(jsonTables io.JsonTables) (dbTables db.Tables) {
	dbTables = db.Tables{}
	for jsonTableName, jsonRows := range jsonTables {
		dbRows := db.Rows{}
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
