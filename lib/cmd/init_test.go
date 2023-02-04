package cmd_test

import (
	"context"
	"testing"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/cmd/mock"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
	"github.com/Jumpaku/api-regression-detector/test/assert"
)

func TestInit_OK(t *testing.T) {
	err := cmd.Init(context.Background(),
		mock.DB{},
		jsonio.Tables{
			"mock_table": jsonio.Table{Rows: []jsonio.Row{
				{"a": jsonio.NewJsonNull()},
				{"b": jsonio.NewJsonBoolean(true)},
				{"c": jsonio.NewJsonNumberInt64(123)},
				{"x": jsonio.NewJsonNumberFloat64(-123.45)},
				{"y": jsonio.NewJsonString("abc")},
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
		jsonio.Tables{
			"invalid_table": jsonio.Table{Rows: []jsonio.Row{
				{"a": jsonio.NewJsonNull()},
				{"b": jsonio.NewJsonBoolean(true)},
				{"c": jsonio.NewJsonNumberInt64(123)},
				{"x": jsonio.NewJsonNumberFloat64(-123.45)},
				{"y": jsonio.NewJsonString("abc")},
			}},
		},
		mock.SchemaGetter{},
		mock.RowClearer{},
		mock.RowCreator{})
	assert.NotEqual(t, err, nil)
}

func TestInit_NG_DB(t *testing.T) {
	err := cmd.Init(context.Background(),
		mock.MockDBErr{},
		jsonio.Tables{
			"mock_table": jsonio.Table{Rows: []jsonio.Row{
				{"a": jsonio.NewJsonNull()},
				{"b": jsonio.NewJsonBoolean(true)},
				{"c": jsonio.NewJsonNumberInt64(123)},
				{"x": jsonio.NewJsonNumberFloat64(-123.45)},
				{"y": jsonio.NewJsonString("abc")},
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
		jsonio.Tables{
			"mock_table": jsonio.Table{Rows: []jsonio.Row{
				{"a": jsonio.NewJsonNull()},
				{"b": jsonio.NewJsonBoolean(true)},
				{"c": jsonio.NewJsonNumberInt64(123)},
				{"x": jsonio.NewJsonNumberFloat64(-123.45)},
				{"y": jsonio.NewJsonString("abc")},
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
		jsonio.Tables{
			"mock_table": jsonio.Table{Rows: []jsonio.Row{
				{"a": jsonio.NewJsonNull()},
				{"b": jsonio.NewJsonBoolean(true)},
				{"c": jsonio.NewJsonNumberInt64(123)},
				{"x": jsonio.NewJsonNumberFloat64(-123.45)},
				{"y": jsonio.NewJsonString("abc")},
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
		jsonio.Tables{
			"mock_table": jsonio.Table{Rows: []jsonio.Row{
				{"a": jsonio.NewJsonNull()},
				{"b": jsonio.NewJsonBoolean(true)},
				{"c": jsonio.NewJsonNumberInt64(123)},
				{"x": jsonio.NewJsonNumberFloat64(-123.45)},
				{"y": jsonio.NewJsonString("abc")},
			}},
		},
		mock.SchemaGetter{},
		mock.RowClearer{},
		mock.ErrRowCreator{})
	assert.NotEqual(t, err, nil)
}

/*
func Init(ctx context.Context,
	db libdb.DB,
	jsonTables jsonio.Tables,
	schemaGetter SchemaGetter,
	clearer RowClearer,
	creator RowCreator,
) (err error) {
	return db.RunTransaction(ctx, func(ctx context.Context, tx libdb.Tx) error {
		for tableName, table := range jsonTables {
			schema, err := schemaGetter.GetSchema(ctx, tx, tableName)
			if err != nil {
				return err
			}
			err = clearer.ClearRows(ctx, tx, tableName)
			if err != nil {
				return err
			}
			err = creator.CreateRows(ctx, tx, tableName, schema, table.Rows)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

/*
func convertTablesJsonToDB(jsonTables jsonio.JsonTables) (dbTables db.Tables) {
	dbTables = db.Tables{}
	for jsonTableName, jsonRows := range jsonTables {
		dbRows := db.Table{}
		for _, jsonRow := range jsonRows {
			dbRow := db.Row{}
			for jsonColumnName, jsonColumnValue := range jsonRow {
				dbRow[jsonColumnName] = jsonColumnValue
			}
			dbRows = append(dbRows, dbRow)
		}
		dbTables[jsonTableName] = dbRows
	}
	return dbTables
}
*/
