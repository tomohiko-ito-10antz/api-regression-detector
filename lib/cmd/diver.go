package cmd

import (
	"context"

	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/tables"
)

type RowLister interface {
	ListRows(ctx context.Context, tx db.Tx, tableName string, schema db.Schema) ([]db.Row, error)
}

type RowClearer interface {
	ClearRows(ctx context.Context, tx db.Tx, tableName string) error
}

type RowCreator interface {
	CreateRows(ctx context.Context, tx db.Tx, tableName string, schema db.Schema, rows []tables.Row) error
}

type SchemaGetter interface {
	GetSchema(ctx context.Context, tx db.Tx, tableName string) (db.Schema, error)
}

type Driver struct {
	Name         string
	DB           db.DB
	RowLister    RowLister
	RowClearer   RowClearer
	RowCreator   RowCreator
	SchemaGetter SchemaGetter
}

func (d *Driver) Close() error {
	if err := d.DB.Close(); err != nil {
		return errors.Wrap(errors.Join(err, errors.IOFailure), "fail to close database")
	}

	return nil
}

func (d *Driver) Open(connectionString string) error {
	switch d.Name {
	default:
		return errors.Wrap(errors.BadArgs, "unsupported driver name")
	case "mysql", "postgres", "sqlite3", "spanner":
	}

	d.DB = db.NewDB(d.Name, connectionString)

	if err := d.DB.Open(); err != nil {
		return errors.Wrap(err, "fail to open database")
	}

	return nil
}
