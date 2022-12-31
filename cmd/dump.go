package cmd

import (
	"context"
	"database/sql"

	"github.com/Jumpaku/api-regression-detector/db"
)

type Select interface {
	Select(ctx context.Context, exec db.Exec, table string) (db.Rows, error)
}

func Dump(ctx context.Context, database *sql.DB, tableNames []string, s Select) (tables db.Tables, err error) {
	tables = db.Tables{}
	err = db.Transaction(ctx, database, func(ctx context.Context, exec db.Exec) error {
		for _, tableName := range tableNames {
			rows, err := s.Select(ctx, exec, tableName)
			if err != nil {
				return err
			}
			tables[tableName] = rows
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tables, nil
}
