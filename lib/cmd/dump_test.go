package cmd

import (
	"context"
	"testing"

	"github.com/Jumpaku/api-regression-detector/test/assert"
)

func TestDump_OK(t *testing.T) {
	v, err := Dump(context.Background(), nil,
		[]string{"mock_table"},
		MockDriver.SchemaGetter,
		MockDriver.ListRows)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(v), 1)
	assert.Equal(t, len(v["mock_table"].Rows), 3)
}
func TestDump_NG(t *testing.T) {}

/*
) (tables io.Tables, err error) {
	tables = io.Tables{}
	err = lib_db.RunTransaction(ctx, db, func(ctx context.Context, tx lib_db.Tx) error {
		dbTables := lib_db.Tables{}
		for _, tableName := range tableNames {
			schema, err := schemaGetter.GetSchema(ctx, tx, tableName)
			if err != nil {
				return err
			}
			rows, err := rowLister.ListRows(ctx, tx, tableName, schema)
			if err != nil {
				return err
			}
			dbTables[tableName] = lib_db.Table{Name: tableName, Schema: schema, Rows: rows}
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

func convertTablesDBToJson(dbTables lib_db.Tables) (jsonTables io.Tables, err error) {
	jsonTables = io.Tables{}
	for dbTableName, dbTable := range dbTables {
		jsonTable := io.Table{}
		for _, dbRow := range dbTable.Rows {
			jsonRow := io.Row{}
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

func convertDBColumnValueToJsonValue(dbVal *lib_db.ColumnValue) (*io.JsonValue, error) {
	switch dbVal.Type {
	case lib_db.ColumnTypeBoolean:
		v, err := dbVal.AsBool()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return io.NewJsonNull(), nil
		}
		return io.NewJsonBoolean(v.Bool), nil
	case lib_db.ColumnTypeInteger:
		v, err := dbVal.AsInteger()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return io.NewJsonNull(), nil
		}
		return io.NewJsonNumberInt64(v.Int64), nil
	case lib_db.ColumnTypeFloat:
		v, err := dbVal.AsFloat()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return io.NewJsonNull(), nil
		}
		return io.NewJsonNumberFloat64(v.Float64), nil
	case lib_db.ColumnTypeString:
		v, err := dbVal.AsString()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return io.NewJsonNull(), nil
		}
		return io.NewJsonString(v.String), nil
	case lib_db.ColumnTypeTime:
		v, err := dbVal.AsTime()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return io.NewJsonNull(), nil
		}
		return io.NewJsonString(v.Time.Format(time.RFC3339)), nil
	default:
		v, err := dbVal.AsBytes()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return io.NewJsonNull(), nil
		}
		return io.NewJsonString(string(v.Bytes)), nil
	}
}
*/
