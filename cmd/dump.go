package cmd

import (
	"context"
	"database/sql"

	"github.com/Jumpaku/api-regression-detector/db"
	"github.com/Jumpaku/api-regression-detector/io"
)

type Select interface {
	Select(ctx context.Context, exec db.Exec, table string) (db.Rows, error)
}

func Dump(ctx context.Context, database *sql.DB, tableNames []string, s Select) (tables io.JsonTables, err error) {
	tables = io.JsonTables{}
	err = db.Transaction(ctx, database, func(ctx context.Context, exec db.Exec) error {
		dbTables := db.Tables{}
		for _, tableName := range tableNames {
			rows, err := s.Select(ctx, exec, tableName)
			if err != nil {
				return err
			}
			dbTables[tableName] = rows
		}
		tables = convertTablesDBToJson(dbTables)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tables, nil
}

func convertTablesDBToJson(dbTables db.Tables) (jsonTables io.JsonTables) {
	jsonTables = io.JsonTables{}
	for dbTableName, dbRows := range dbTables {
		jsonRows := io.JsonRows{}
		for _, dbRow := range dbRows {
			jsonRow := io.JsonRow{}
			for dbColumnName, dbColumnValue := range dbRow {
				jsonRow[dbColumnName] = dbColumnValue
			}
			jsonRows = append(jsonRows, jsonRow)
		}
		jsonTables[dbTableName] = jsonRows
	}
	return jsonTables
}
