package mock

import (
	"context"
	"fmt"

	"github.com/Jumpaku/api-regression-detector/lib/db"
)

type RowClearer struct{}

func (RowClearer) ClearRows(ctx context.Context, tx db.Tx, tableName string) error {
	if tableName != "mock_table" {
		return fmt.Errorf("table %s not found", tableName)
	}

	return nil
}

type ErrRowClearer struct{}

func (ErrRowClearer) ClearRows(ctx context.Context, tx db.Tx, tableName string) error {
	return fmt.Errorf("error with table %s", tableName)
}
