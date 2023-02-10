package mock

import (
	"context"

	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/tables"
	"github.com/Jumpaku/api-regression-detector/test"
)

type RowCreator struct{}

func (RowCreator) CreateRows(
	ctx context.Context,
	tx db.Tx,
	tableName string,
	schema db.Schema,
	rows []tables.Row,
) error {
	if tableName != "mock_table" {
		return errors.Wrap(test.MockError, "table %s not found", tableName)
	}

	return nil
}

type ErrRowCreator struct{}

func (ErrRowCreator) CreateRows(
	ctx context.Context,
	tx db.Tx,
	tableName string,
	schema db.Schema,
	rows []tables.Row,
) error {
	return errors.Wrap(test.MockError, "error with database")
}
