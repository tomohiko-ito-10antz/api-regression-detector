package cmd

import (
	"context"
	"testing"

	"github.com/Jumpaku/api-regression-detector/lib/cmd/mock"
	"github.com/Jumpaku/api-regression-detector/lib/io"
	"github.com/Jumpaku/api-regression-detector/test/assert"
)

func TestInit_OK(t *testing.T) {
	err := Init(context.Background(),
		mock.MockDB{},
		io.Tables{
			"mock_table": io.Table{Rows: []io.Row{
				{"a": io.NewJsonNull()},
				{"b": io.NewJsonBoolean(true)},
				{"c": io.NewJsonNumberInt64(123)},
				{"x": io.NewJsonNumberFloat64(-123.45)},
				{"y": io.NewJsonString("abc")},
			}},
		},
		mock.MockSchemaGetter{},
		mock.MockRowClearer{},
		mock.MockRowCreator{})
	assert.Equal(t, err, nil)
}
func TestInit_NG_Table(t *testing.T) {
	err := Init(context.Background(),
		mock.MockDB{},
		io.Tables{
			"invalid_table": io.Table{Rows: []io.Row{
				{"a": io.NewJsonNull()},
				{"b": io.NewJsonBoolean(true)},
				{"c": io.NewJsonNumberInt64(123)},
				{"x": io.NewJsonNumberFloat64(-123.45)},
				{"y": io.NewJsonString("abc")},
			}},
		},
		mock.MockSchemaGetter{},
		mock.MockRowClearer{},
		mock.MockRowCreator{})
	assert.NotEqual(t, err, nil)
}

func TestInit_NG_DB(t *testing.T) {
	err := Init(context.Background(),
		mock.MockDBErr{},
		io.Tables{
			"mock_table": io.Table{Rows: []io.Row{
				{"a": io.NewJsonNull()},
				{"b": io.NewJsonBoolean(true)},
				{"c": io.NewJsonNumberInt64(123)},
				{"x": io.NewJsonNumberFloat64(-123.45)},
				{"y": io.NewJsonString("abc")},
			}},
		},
		mock.MockSchemaGetter{},
		mock.MockRowClearer{},
		mock.MockRowCreator{})
	assert.NotEqual(t, err, nil)
}
func TestInit_NG_SchemaGetter(t *testing.T) {
	err := Init(context.Background(),
		mock.MockDB{},
		io.Tables{
			"mock_table": io.Table{Rows: []io.Row{
				{"a": io.NewJsonNull()},
				{"b": io.NewJsonBoolean(true)},
				{"c": io.NewJsonNumberInt64(123)},
				{"x": io.NewJsonNumberFloat64(-123.45)},
				{"y": io.NewJsonString("abc")},
			}},
		},
		mock.MockSchemaGetterErr{},
		mock.MockRowClearer{},
		mock.MockRowCreator{})
	assert.NotEqual(t, err, nil)
}

func TestInit_NG_RowClearer(t *testing.T) {
	err := Init(context.Background(),
		mock.MockDB{},
		io.Tables{
			"mock_table": io.Table{Rows: []io.Row{
				{"a": io.NewJsonNull()},
				{"b": io.NewJsonBoolean(true)},
				{"c": io.NewJsonNumberInt64(123)},
				{"x": io.NewJsonNumberFloat64(-123.45)},
				{"y": io.NewJsonString("abc")},
			}},
		},
		mock.MockSchemaGetter{},
		mock.MockRowClearerErr{},
		mock.MockRowCreator{})
	assert.NotEqual(t, err, nil)
}

func TestInit_NG_RowCreator(t *testing.T) {
	err := Init(context.Background(),
		mock.MockDB{},
		io.Tables{
			"mock_table": io.Table{Rows: []io.Row{
				{"a": io.NewJsonNull()},
				{"b": io.NewJsonBoolean(true)},
				{"c": io.NewJsonNumberInt64(123)},
				{"x": io.NewJsonNumberFloat64(-123.45)},
				{"y": io.NewJsonString("abc")},
			}},
		},
		mock.MockSchemaGetter{},
		mock.MockRowClearer{},
		mock.MockRowCreatorErr{})
	assert.NotEqual(t, err, nil)
}

/*
func Init(ctx context.Context,
	db lib_db.DB,
	jsonTables io.Tables,
	schemaGetter SchemaGetter,
	clearer RowClearer,
	creator RowCreator,
) (err error) {
	return db.RunTransaction(ctx, func(ctx context.Context, tx lib_db.Tx) error {
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
func convertTablesJsonToDB(jsonTables io.JsonTables) (dbTables db.Tables) {
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
