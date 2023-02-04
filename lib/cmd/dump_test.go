package cmd_test

import (
	"context"
	"testing"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/cmd/mock"
	"github.com/Jumpaku/api-regression-detector/test/assert"
)

func TestDump_OK(t *testing.T) {
	v, err := cmd.Dump(context.Background(),
		mock.DB{},
		[]string{"mock_table"},
		mock.SchemaGetter{},
		mock.RowLister{})
	assert.Equal(t, err, nil)
	assert.Equal(t, len(v), 1)
	assert.Equal(t, len(v["mock_table"].Rows), 3)
}

func TestDump_NG_Table(t *testing.T) {
	_, err := cmd.Dump(context.Background(),
		mock.DB{},
		[]string{"invalid_table"},
		mock.SchemaGetter{},
		mock.RowLister{})
	assert.NotEqual(t, err, nil)
}

func TestDump_NG_DB(t *testing.T) {
	_, err := cmd.Dump(context.Background(),
		mock.MockDBErr{},
		[]string{"mock_table"},
		mock.SchemaGetter{},
		mock.RowLister{})
	assert.NotEqual(t, err, nil)
}

func TestDump_NG_SchemaGetter(t *testing.T) {
	_, err := cmd.Dump(context.Background(),
		mock.DB{},
		[]string{"mock_table"},
		mock.ErrSchemaGetter{},
		mock.RowLister{})
	assert.NotEqual(t, err, nil)
}

func TestDump_NG_RowLister(t *testing.T) {
	_, err := cmd.Dump(context.Background(),
		mock.DB{},
		[]string{"mock_table"},
		mock.SchemaGetter{},
		mock.ErrRowLister{})
	assert.NotEqual(t, err, nil)
}

/*
) (tables jsonio.Tables, err error) {
	tables = jsonio.Tables{}
	err = libdb.RunTransaction(ctx, db, func(ctx context.Context, tx libdb.Tx) error {
		dbTables := libdb.Tables{}
		for _, tableName := range tableNames {
			schema, err := schemaGetter.GetSchema(ctx, tx, tableName)
			if err != nil {
				return err
			}
			rows, err := rowLister.ListRows(ctx, tx, tableName, schema)
			if err != nil {
				return err
			}
			dbTables[tableName] = libdb.Table{Name: tableName, Schema: schema, Rows: rows}
		}
		tables, err = convertTablesDBToJson(dbTables)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tables, nil
}

func convertTablesDBToJson(dbTables libdb.Tables) (jsonTables jsonio.Tables, err error) {
	jsonTables = jsonio.Tables{}
	for dbTableName, dbTable := range dbTables {
		jsonTable := jsonio.Table{}
		for _, dbRow := range dbTable.Rows {
			jsonRow := jsonio.Row{}
			for dbColumnName, dbColumnValue := range dbRow {
				jsonRow[dbColumnName], err = convertDBColumnValueToJsonValue(dbColumnValue)
				if err != nil {
					return nil, err
				}
			}
			jsonTable.Rows = append(jsonTable.Rows, jsonRow)
		}
		jsonTables[dbTableName] = jsonTable
	}
	return jsonTables, nil
}

func convertDBColumnValueToJsonValue(dbVal *libdb.ColumnValue) (*jsonio.JsonValue, error) {
	switch dbVal.Type {
	case libdb.ColumnTypeBoolean:
		v, err := dbVal.AsBool()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return jsonio.NewJsonNull(), nil
		}
		return jsonio.NewJsonBoolean(v.Bool), nil
	case libdb.ColumnTypeInteger:
		v, err := dbVal.AsInteger()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return jsonio.NewJsonNull(), nil
		}
		return jsonio.NewJsonNumberInt64(v.Int64), nil
	case libdb.ColumnTypeFloat:
		v, err := dbVal.AsFloat()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return jsonio.NewJsonNull(), nil
		}
		return jsonio.NewJsonNumberFloat64(v.Float64), nil
	case libdb.ColumnTypeString:
		v, err := dbVal.AsString()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return jsonio.NewJsonNull(), nil
		}
		return jsonio.NewJsonString(v.String), nil
	case libdb.ColumnTypeTime:
		v, err := dbVal.AsTime()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return jsonio.NewJsonNull(), nil
		}
		return jsonio.NewJsonString(v.Time.Format(time.RFC3339)), nil
	default:
		v, err := dbVal.AsBytes()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return jsonio.NewJsonNull(), nil
		}
		return jsonio.NewJsonString(string(v.Bytes)), nil
	}
}
*/
