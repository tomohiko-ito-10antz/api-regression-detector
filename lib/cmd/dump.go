package cmd

import (
	"context"
	"database/sql"
	"time"

	lib_db "github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/io"
)

func Dump(
	ctx context.Context,
	db *sql.DB,
	tableNames []string,
	schemaGetter SchemaGetter,
	rowLister RowLister,
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
