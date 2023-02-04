package mock

import (
	"context"

	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/test"
)

type SchemaGetter struct{}

func (SchemaGetter) GetSchema(ctx context.Context, tx db.Tx, tableName string) (db.Schema, error) {
	if tableName != "mock_table" {
		return db.Schema{}, errors.Wrap(test.MockError, "table %s not found", tableName)
	}

	return db.Schema{
		PrimaryKeys: []string{"column_a", "column_b", "column_c"},
		ColumnTypes: db.ColumnTypes{
			"column_a": db.ColumnTypeBoolean,
			"column_b": db.ColumnTypeInteger,
			"column_c": db.ColumnTypeString,
			"column_x": db.ColumnTypeFloat,
			"column_y": db.ColumnTypeTime,
		},
	}, nil
}

type ErrSchemaGetter struct{}

func (ErrSchemaGetter) GetSchema(ctx context.Context, tx db.Tx, tableName string) (db.Schema, error) {
	return db.Schema{}, errors.Wrap(test.MockError, "error with database")
}
