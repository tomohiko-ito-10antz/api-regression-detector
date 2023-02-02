package mock

import (
	"context"
	"fmt"

	"github.com/Jumpaku/api-regression-detector/lib/db"
)

type MockRowClearer struct{}

func (MockRowClearer) ClearRows(ctx context.Context, tx db.Tx, tableName string) error {
	if tableName != "mock_table" {
		return fmt.Errorf("table %s not found", tableName)
	}
	return nil
}

type MockRowClearerErr struct{}

func (MockRowClearerErr) ClearRows(ctx context.Context, tx db.Tx, tableName string) error {
	return fmt.Errorf("error with table %s", tableName)
}
