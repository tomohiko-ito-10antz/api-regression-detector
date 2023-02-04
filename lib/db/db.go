package db

import (
	"context"
	"database/sql"
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
	return runTransaction(ctx, d.db, handler)
}

func (d *database) Open() error {
	db, err := sql.Open(d.driver, d.connection)
	if err != nil {
		return err
	}
	d.db = db

	return nil
}

func (d *database) Close() error {
	return d.db.Close()
}
