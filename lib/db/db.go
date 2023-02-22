package db

import (
	"context"
	"database/sql"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

type DB interface {
	RunTransaction(ctx context.Context, handler func(ctx context.Context, tx Tx) error) error
	Open() error
	Close() error
}

func NewDB(driver string, connection string) *database {
	return &database{connection: connection, driver: driver}
}

type database struct {
	driver     string
	connection string
	db         *sql.DB
}

func (d *database) RunTransaction(ctx context.Context, handler func(ctx context.Context, tx Tx) error) error {
	err := runTransaction(ctx, d.db, handler)
	if err != nil {
		return errors.Wrap(err, "transaction failed")
	}

	return nil
}

func (d *database) Open() error {
	errInfo := errors.Info{"driverName": d.driver, "connectionString": d.connection}

	db, err := sql.Open(d.driver, d.connection)
	if err != nil {
		return errors.Wrap(errors.DBFailure.Err(err), errInfo.AppendTo("fail to open database"))
	}

	d.db = db

	return nil
}

func (d *database) Close() error {
	errInfo := errors.Info{"driverName": d.driver, "connectionString": d.connection}

	if err := d.db.Close(); err != nil {
		return errors.Wrap(errors.DBFailure.Err(err), errInfo.AppendTo("fail to close database"))
	}

	return nil
}
