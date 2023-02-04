package mock

import (
	"context"
	"fmt"

	"github.com/Jumpaku/api-regression-detector/lib/db"
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

type MockDBErr struct{}

func (MockDBErr) Open() error {
	return nil
}

func (MockDBErr) Close() error {
	return nil
}

func (MockDBErr) RunTransaction(ctx context.Context, handler func(ctx context.Context, tx db.Tx) error) error {
	return fmt.Errorf("error with database")
}
