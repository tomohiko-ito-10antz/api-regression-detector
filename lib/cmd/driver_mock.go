package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/io"
)

var MockDriver = Driver{
	Name:         "mock-driver",
	DB:           &sql.DB{},
	ListRows:     rowLister{},
	ClearRows:    rowClearer{},
	CreateRows:   rowCreator{},
	SchemaGetter: schemaGetter{},
}

type rowLister struct{}

func (rowLister) ListRows(ctx context.Context, tx db.Tx, tableName string, schema db.Schema) ([]db.Row, error) {
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

type rowClearer struct{}

func (rowClearer) ClearRows(ctx context.Context, tx db.Tx, tableName string) error {
	if tableName != "mock_table" {
		return fmt.Errorf("table %s not found", tableName)
	}
	return nil
}

type rowCreator struct{}

func (rowCreator) CreateRows(ctx context.Context, tx db.Tx, tableName string, schema db.Schema, rows []io.Row) error {
	if tableName != "mock_table" {
		return fmt.Errorf("table %s not found", tableName)
	}
	return nil
}

type schemaGetter struct{}

func (schemaGetter) GetSchema(ctx context.Context, tx db.Tx, tableName string) (db.Schema, error) {
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
