package mock

import (
	"context"
	"fmt"

	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/io"
)

type MockRowCreator struct{}

func (MockRowCreator) CreateRows(ctx context.Context, tx db.Tx, tableName string, schema db.Schema, rows []io.Row) error {
	if tableName != "mock_table" {
		return fmt.Errorf("table %s not found", tableName)
	}
	return nil
}

type MockRowCreatorErr struct{}

func (MockRowCreatorErr) CreateRows(ctx context.Context, tx db.Tx, tableName string, schema db.Schema, rows []io.Row) error {
	return fmt.Errorf("error with table %s", tableName)
}
