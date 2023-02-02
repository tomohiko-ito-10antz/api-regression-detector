package cmd

import (
	"context"
	"fmt"

	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/io"
)

type RowLister interface {
	ListRows(ctx context.Context, tx db.Tx, tableName string, schema db.Schema) ([]db.Row, error)
}

type RowClearer interface {
	ClearRows(ctx context.Context, tx db.Tx, tableName string) error
}

type RowCreator interface {
	CreateRows(ctx context.Context, tx db.Tx, tableName string, schema db.Schema, rows []io.Row) error
}

type SchemaGetter interface {
	GetSchema(ctx context.Context, tx db.Tx, tableName string) (db.Schema, error)
}

type Driver struct {
	Name         string
	DB           db.DB
	ListRows     RowLister
	ClearRows    RowClearer
	CreateRows   RowCreator
	SchemaGetter SchemaGetter
}

func (d *Driver) Close() error {
	return d.DB.Close()
}

func (d *Driver) Open(connectionString string) (err error) {
	switch d.Name {
	default:
		return fmt.Errorf("invalid driver name")
	case "mysql", "postgres", "sqlite3", "spanner":
	}
	d.DB = db.NewDB(d.Name, connectionString)
	return d.DB.Open()
}
