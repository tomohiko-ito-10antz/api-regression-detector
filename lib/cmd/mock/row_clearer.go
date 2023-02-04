package mock

import (
	"context"

	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/test"
)

type RowClearer struct{}

func (RowClearer) ClearRows(ctx context.Context, tx db.Tx, tableName string) error {
	if tableName != "mock_table" {
		return errors.Wrap(test.MockError, "table %s not found", tableName)
	}

	return nil
}

type ErrRowClearer struct{}

func (ErrRowClearer) ClearRows(ctx context.Context, tx db.Tx, tableName string) error {
	return errors.Wrap(test.MockError, "error with database")
}
