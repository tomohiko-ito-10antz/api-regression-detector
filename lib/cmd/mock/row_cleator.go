package mock

import (
	"context"
	"fmt"

	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
)

type RowCreator struct{}

func (RowCreator) CreateRows(
	ctx context.Context,
	tx db.Tx,
	tableName string,
	schema db.Schema,
	rows []jsonio.Row,
) error {
	if tableName != "mock_table" {
		return fmt.Errorf("table %s not found", tableName)
	}
	return nil
}

type ErrRowCreator struct{}

func (ErrRowCreator) CreateRows(
	ctx context.Context,
	tx db.Tx,
	tableName string,
	schema db.Schema,
	rows []jsonio.Row,
) error {
	return fmt.Errorf("error with table %s", tableName)
}
