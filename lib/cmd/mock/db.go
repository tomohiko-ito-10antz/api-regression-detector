package mock

import (
	"context"

	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/test"
)

type DB struct{}

func (DB) Open() error {
	return nil
}

func (DB) Close() error {
	return nil
}

func (DB) RunTransaction(ctx context.Context, handler func(ctx context.Context, tx db.Tx) error) error {
	return handler(ctx, nil)
}

type ErrDB struct{}

func (ErrDB) Open() error {
	return nil
}

func (ErrDB) Close() error {
	return nil
}

func (ErrDB) RunTransaction(ctx context.Context, handler func(ctx context.Context, tx db.Tx) error) error {
	return errors.Wrap(test.MockError, "error with database")
}
