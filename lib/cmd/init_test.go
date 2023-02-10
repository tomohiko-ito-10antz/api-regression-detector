package cmd_test

import (
	"context"
	"testing"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/cmd/mock"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/tables"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
	"github.com/Jumpaku/api-regression-detector/test/assert"
)

func TestInit_OK(t *testing.T) {
	err := cmd.Init(context.Background(),
		mock.DB{},
		tables.InitTables{
			tables.Table{
				Name: "mock_table",
				Rows: []tables.Row{
					{"a": wrap.Null()},
					{"b": wrap.Boolean(true)},
					{"c": wrap.Number(123)},
					{"x": wrap.Number(-123.45)},
					{"y": wrap.String("abc")},
				}},
		},
		mock.SchemaGetter{},
		mock.RowClearer{},
		mock.RowCreator{})
	assert.Equal(t, err, nil)
}

func TestInit_NG_Table(t *testing.T) {
	err := cmd.Init(context.Background(),
		mock.DB{},
		tables.InitTables{
			tables.Table{
				Name: "invalid_table",
				Rows: []tables.Row{
					{"a": wrap.Null()},
					{"b": wrap.Boolean(true)},
					{"c": wrap.Number(123)},
					{"x": wrap.Number(-123.45)},
					{"y": wrap.String("abc")},
				}},
		},
		mock.SchemaGetter{},
		mock.RowClearer{},
		mock.RowCreator{})
	assert.NotEqual(t, err, nil)
}

func TestInit_NG_DB(t *testing.T) {
	err := cmd.Init(context.Background(),
		mock.ErrDB{},
		tables.InitTables{
			tables.Table{
				Name: "mock_table",
				Rows: []tables.Row{
					{"a": wrap.Null()},
					{"b": wrap.Boolean(true)},
					{"c": wrap.Number(123)},
					{"x": wrap.Number(-123.45)},
					{"y": wrap.String("abc")},
				}},
		},
		mock.SchemaGetter{},
		mock.RowClearer{},
		mock.RowCreator{})
	assert.NotEqual(t, err, nil)
}

func TestInit_NG_SchemaGetter(t *testing.T) {
	err := cmd.Init(context.Background(),
		mock.DB{},
		tables.InitTables{
			tables.Table{
				Name: "mock_table",
				Rows: []tables.Row{
					{"a": wrap.Null()},
					{"b": wrap.Boolean(true)},
					{"c": wrap.Number(123)},
					{"x": wrap.Number(-123.45)},
					{"y": wrap.String("abc")},
				}},
		},
		mock.ErrSchemaGetter{},
		mock.RowClearer{},
		mock.RowCreator{})
	assert.NotEqual(t, err, nil)
}

func TestInit_NG_RowClearer(t *testing.T) {
	err := cmd.Init(context.Background(),
		mock.DB{},
		tables.InitTables{
			tables.Table{
				Name: "mock_table",
				Rows: []tables.Row{
					{"a": wrap.Null()},
					{"b": wrap.Boolean(true)},
					{"c": wrap.Number(123)},
					{"x": wrap.Number(-123.45)},
					{"y": wrap.String("abc")},
				}},
		},
		mock.SchemaGetter{},
		mock.ErrRowClearer{},
		mock.RowCreator{})
	assert.NotEqual(t, err, nil)
}

func TestInit_NG_RowCreator(t *testing.T) {
	err := cmd.Init(context.Background(),
		mock.DB{},
		tables.InitTables{
			tables.Table{
				Name: "mock_table",
				Rows: []tables.Row{
					{"a": wrap.Null()},
					{"b": wrap.Boolean(true)},
					{"c": wrap.Number(123)},
					{"x": wrap.Number(-123.45)},
					{"y": wrap.String("abc")},
				}},
		},
		mock.SchemaGetter{},
		mock.RowClearer{},
		mock.ErrRowCreator{})
	assert.NotEqual(t, err, nil)
}
