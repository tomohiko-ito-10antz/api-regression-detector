package cmd

import (
	"context"

	"github.com/Jumpaku/api-regression-detector/db"
	"github.com/Jumpaku/api-regression-detector/io"
)

type RowLister interface {
	ListRows(ctx context.Context, exec db.Transaction, tableName string) (db.Table, error)
}

type RowClearer interface {
	ClearRows(ctx context.Context, exec db.Transaction, tableName string) error
}

type RowCreator interface {
	CreateRows(ctx context.Context, exec db.Transaction, tableName string, table io.Table) error
}

type SchemaGetter interface {
	GetSchema(ctx context.Context, exec db.Transaction, tableName string) (db.Schema, error)
}
