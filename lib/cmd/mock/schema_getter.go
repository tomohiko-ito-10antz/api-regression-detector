package mock

import (
	"context"
	"fmt"

	"github.com/Jumpaku/api-regression-detector/lib/db"
)

type SchemaGetter struct{}

func (SchemaGetter) GetSchema(ctx context.Context, tx db.Tx, tableName string) (db.Schema, error) {
	if tableName != "mock_table" {
		return db.Schema{}, fmt.Errorf("table %s not found", tableName)
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
	return db.Schema{}, fmt.Errorf("error with table %s", tableName)
}
