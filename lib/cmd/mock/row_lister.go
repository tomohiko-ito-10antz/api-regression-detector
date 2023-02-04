package mock

import (
	"context"
	"fmt"
	"time"

	"github.com/Jumpaku/api-regression-detector/lib/db"
)

type RowLister struct{}

func (RowLister) ListRows(ctx context.Context, tx db.Tx, tableName string, schema db.Schema) ([]db.Row, error) {
	if tableName != "mock_table" {
		return nil, fmt.Errorf("table %s not found", tableName)
	}
	return []db.Row{
		{
			"column_a": db.NewColumnValue(true, db.ColumnTypeBoolean),
			"column_b": db.NewColumnValue(123, db.ColumnTypeInteger),
			"column_c": db.NewColumnValue("abc", db.ColumnTypeString),
			"column_x": db.NewColumnValue(-123.45, db.ColumnTypeFloat),
			"column_y": db.NewColumnValue(time.Now(), db.ColumnTypeTime),
		},
		{
			"column_a": db.NewColumnValue(false, db.ColumnTypeBoolean),
			"column_b": db.NewColumnValue(0, db.ColumnTypeInteger),
			"column_c": db.NewColumnValue("", db.ColumnTypeString),
			"column_x": db.NewColumnValue(0.0, db.ColumnTypeFloat),
			"column_y": db.NewColumnValue(time.Time{}, db.ColumnTypeTime),
		},
		{
			"column_a": db.NewColumnValue(nil, db.ColumnTypeBoolean),
			"column_b": db.NewColumnValue(nil, db.ColumnTypeInteger),
			"column_c": db.NewColumnValue(nil, db.ColumnTypeString),
			"column_x": db.NewColumnValue(nil, db.ColumnTypeFloat),
			"column_y": db.NewColumnValue(nil, db.ColumnTypeTime),
		},
	}, nil
}

type ErrRowLister struct{}

func (ErrRowLister) ListRows(ctx context.Context, tx db.Tx, tableName string, schema db.Schema) ([]db.Row, error) {
	return nil, fmt.Errorf("error table %s", tableName)
}
