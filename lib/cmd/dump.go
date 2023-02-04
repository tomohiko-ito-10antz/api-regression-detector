package cmd

import (
	"context"
	"time"

	lib_db "github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/io_json"
)

func Dump(
	ctx context.Context,
	db lib_db.DB,
	tableNames []string,
	schemaGetter SchemaGetter,
	rowLister RowLister,
) (tables io_json.Tables, err error) {
	tables = io_json.Tables{}
	err = db.RunTransaction(ctx, func(ctx context.Context, tx lib_db.Tx) error {
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

func convertTablesDBToJson(dbTables lib_db.Tables) (jsonTables io_json.Tables, err error) {
	jsonTables = io_json.Tables{}
	for dbTableName, dbTable := range dbTables {
		jsonTable := io_json.Table{}
		for _, dbRow := range dbTable.Rows {
			jsonRow := io_json.Row{}
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

func convertDBColumnValueToJsonValue(dbVal *lib_db.ColumnValue) (*io_json.JsonValue, error) {
	switch dbVal.Type {
	case lib_db.ColumnTypeBoolean:
		v, err := dbVal.AsBool()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return io_json.NewJsonNull(), nil
		}
		return io_json.NewJsonBoolean(v.Bool), nil
	case lib_db.ColumnTypeInteger:
		v, err := dbVal.AsInteger()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return io_json.NewJsonNull(), nil
		}
		return io_json.NewJsonNumberInt64(v.Int64), nil
	case lib_db.ColumnTypeFloat:
		v, err := dbVal.AsFloat()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return io_json.NewJsonNull(), nil
		}
		return io_json.NewJsonNumberFloat64(v.Float64), nil
	case lib_db.ColumnTypeString:
		v, err := dbVal.AsString()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return io_json.NewJsonNull(), nil
		}
		return io_json.NewJsonString(v.String), nil
	case lib_db.ColumnTypeTime:
		v, err := dbVal.AsTime()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return io_json.NewJsonNull(), nil
		}
		return io_json.NewJsonString(v.Time.Format(time.RFC3339)), nil
	default:
		v, err := dbVal.AsBytes()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return io_json.NewJsonNull(), nil
		}
		return io_json.NewJsonString(string(v.Bytes)), nil
	}
}
