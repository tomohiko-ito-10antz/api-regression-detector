package cmd_test

import (
	"context"
	"testing"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/cmd/mock"
	"github.com/Jumpaku/api-regression-detector/lib/io_json"
	"github.com/Jumpaku/api-regression-detector/test/assert"
)

func TestInit_OK(t *testing.T) {
	err := cmd.Init(context.Background(),
		mock.MockDB{},
		io_json.Tables{
			"mock_table": io_json.Table{Rows: []io_json.Row{
				{"a": io_json.NewJsonNull()},
				{"b": io_json.NewJsonBoolean(true)},
				{"c": io_json.NewJsonNumberInt64(123)},
				{"x": io_json.NewJsonNumberFloat64(-123.45)},
				{"y": io_json.NewJsonString("abc")},
			}},
		},
		mock.MockSchemaGetter{},
		mock.MockRowClearer{},
		mock.MockRowCreator{})
	assert.Equal(t, err, nil)
}
func TestInit_NG_Table(t *testing.T) {
	err := cmd.Init(context.Background(),
		mock.MockDB{},
		io_json.Tables{
			"invalid_table": io_json.Table{Rows: []io_json.Row{
				{"a": io_json.NewJsonNull()},
				{"b": io_json.NewJsonBoolean(true)},
				{"c": io_json.NewJsonNumberInt64(123)},
				{"x": io_json.NewJsonNumberFloat64(-123.45)},
				{"y": io_json.NewJsonString("abc")},
			}},
		},
		mock.MockSchemaGetter{},
		mock.MockRowClearer{},
		mock.MockRowCreator{})
	assert.NotEqual(t, err, nil)
}

func TestInit_NG_DB(t *testing.T) {
	err := cmd.Init(context.Background(),
		mock.MockDBErr{},
		io_json.Tables{
			"mock_table": io_json.Table{Rows: []io_json.Row{
				{"a": io_json.NewJsonNull()},
				{"b": io_json.NewJsonBoolean(true)},
				{"c": io_json.NewJsonNumberInt64(123)},
				{"x": io_json.NewJsonNumberFloat64(-123.45)},
				{"y": io_json.NewJsonString("abc")},
			}},
		},
		mock.MockSchemaGetter{},
		mock.MockRowClearer{},
		mock.MockRowCreator{})
	assert.NotEqual(t, err, nil)
}
func TestInit_NG_SchemaGetter(t *testing.T) {
	err := cmd.Init(context.Background(),
		mock.MockDB{},
		io_json.Tables{
			"mock_table": io_json.Table{Rows: []io_json.Row{
				{"a": io_json.NewJsonNull()},
				{"b": io_json.NewJsonBoolean(true)},
				{"c": io_json.NewJsonNumberInt64(123)},
				{"x": io_json.NewJsonNumberFloat64(-123.45)},
				{"y": io_json.NewJsonString("abc")},
			}},
		},
		mock.MockSchemaGetterErr{},
		mock.MockRowClearer{},
		mock.MockRowCreator{})
	assert.NotEqual(t, err, nil)
}

func TestInit_NG_RowClearer(t *testing.T) {
	err := cmd.Init(context.Background(),
		mock.MockDB{},
		io_json.Tables{
			"mock_table": io_json.Table{Rows: []io_json.Row{
				{"a": io_json.NewJsonNull()},
				{"b": io_json.NewJsonBoolean(true)},
				{"c": io_json.NewJsonNumberInt64(123)},
				{"x": io_json.NewJsonNumberFloat64(-123.45)},
				{"y": io_json.NewJsonString("abc")},
			}},
		},
		mock.MockSchemaGetter{},
		mock.MockRowClearerErr{},
		mock.MockRowCreator{})
	assert.NotEqual(t, err, nil)
}

func TestInit_NG_RowCreator(t *testing.T) {
	err := cmd.Init(context.Background(),
		mock.MockDB{},
		io_json.Tables{
			"mock_table": io_json.Table{Rows: []io_json.Row{
				{"a": io_json.NewJsonNull()},
				{"b": io_json.NewJsonBoolean(true)},
				{"c": io_json.NewJsonNumberInt64(123)},
				{"x": io_json.NewJsonNumberFloat64(-123.45)},
				{"y": io_json.NewJsonString("abc")},
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
	jsonTables io_json.Tables,
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
func convertTablesJsonToDB(jsonTables io_json.JsonTables) (dbTables db.Tables) {
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
