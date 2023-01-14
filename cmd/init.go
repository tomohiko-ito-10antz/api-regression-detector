package cmd

import (
	"context"
	"database/sql"

	"github.com/Jumpaku/api-regression-detector/db"
)

type Truncate interface {
	Truncate(ctx context.Context, exec db.Exec, table string) error
}
type Insert interface {
	Insert(ctx context.Context, exec db.Exec, table string, rows db.Rows) error
}

func Init(ctx context.Context, database *sql.DB, tables db.Tables, truncate Truncate, insert Insert) (err error) {
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
