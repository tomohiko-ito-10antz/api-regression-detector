package mock

import (
	"context"
	"fmt"

	"github.com/Jumpaku/api-regression-detector/lib/db"
)

type MockDB struct {
}

func (MockDB) Open() error {
	return nil
}

func (MockDB) Close() error {
	return nil
}

func (MockDB) RunTransaction(ctx context.Context, handler func(ctx context.Context, tx db.Tx) error) error {
	return handler(ctx, nil)
}

type MockDBErr struct {
}

func (MockDBErr) Open() error {
	return nil
}

func (MockDBErr) Close() error {
	return nil
}

func (MockDBErr) RunTransaction(ctx context.Context, handler func(ctx context.Context, tx db.Tx) error) error {
	return fmt.Errorf("error with database")
}
