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
	return RunTransaction(ctx, d.db, handler)
}
func (d *database) Open() (err error) {
	if d == nil || d.db == nil {
		return nil
	}
	d.db, err = sql.Open(d.driver, d.connection)
	return err
}
func (d *database) Close() error {
	if d == nil || d.db == nil {
		return nil
	}
	return d.db.Close()
}
