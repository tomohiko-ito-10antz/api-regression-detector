package cmd

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Jumpaku/api-regression-detector/db"
	"github.com/Jumpaku/api-regression-detector/io"
)

type RowLister interface {
	ListRows(ctx context.Context, exec db.Transaction, tableName string, schema db.Schema) ([]db.Row, error)
}

type RowClearer interface {
	ClearRows(ctx context.Context, exec db.Transaction, tableName string) error
}

type RowCreator interface {
	CreateRows(ctx context.Context, exec db.Transaction, tableName string, rows []io.Row) error
}

type SchemaGetter interface {
	GetSchema(ctx context.Context, exec db.Transaction, tableName string) (db.Schema, error)
}

type Driver struct {
	Name         string
	DB           *sql.DB
	ListRows     RowLister
	ClearRows    RowClearer
	CreateRows   RowCreator
	SchemaGetter SchemaGetter
}

func (d *Driver) Close() error {
	return d.DB.Close()
}

func (d *Driver) Open(name string, connectionString string) (err error) {
	switch name {
	default:
		return fmt.Errorf("invalid driver name")
	case "mysql":
	case "postgres":
	case "sqlite3":
	case "spanner":
	}
	d.DB, err = sql.Open(name, connectionString)
	if err != nil {
		return err
	}
	d.Name = name
	return nil
}
