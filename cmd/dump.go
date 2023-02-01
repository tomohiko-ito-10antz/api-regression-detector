package cmd

import (
	"context"
	"database/sql"
	"time"

	"github.com/Jumpaku/api-regression-detector/db"
	"github.com/Jumpaku/api-regression-detector/io"
)

func Dump(
	ctx context.Context,
	database *sql.DB,
	tableNames []string,
	schemaGetter SchemaGetter,
	rowLister RowLister,
) (tables io.Tables, err error) {
	tables = io.Tables{}
	err = db.RunTransaction(ctx, database, func(ctx context.Context, tx db.Transaction) error {
		dbTables := db.Tables{}
		for _, tableName := range tableNames {
			schema, err := schemaGetter.GetSchema(ctx, tx, tableName)
			if err != nil {
				return err
			}
			rows, err := rowLister.ListRows(ctx, tx, tableName, schema)
			if err != nil {
				return err
			}
			dbTables[tableName] = db.Table{Name: tableName, Schema: schema, Rows: rows}
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

func convertTablesDBToJson(dbTables db.Tables) (jsonTables io.Tables, err error) {
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

func convertDBColumnValueToJsonValue(dbVal *db.ColumnValue) (*io.JsonValue, error) {
	switch dbVal.Type {
	case db.ColumnTypeBoolean:
		v, err := dbVal.AsBool()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return io.NewJsonNull(), nil
		}
		return io.NewJsonBoolean(v.Bool), nil
	case db.ColumnTypeInteger:
		v, err := dbVal.AsInteger()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return io.NewJsonNull(), nil
		}
		return io.NewJsonNumberInt64(v.Int64), nil
	case db.ColumnTypeFloat:
		v, err := dbVal.AsFloat()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return io.NewJsonNull(), nil
		}
		return io.NewJsonNumberFloat64(v.Float64), nil
	case db.ColumnTypeString:
		v, err := dbVal.AsString()
		if err != nil {
			return nil, err
		}
		if !v.Valid {
			return io.NewJsonNull(), nil
		}
		return io.NewJsonString(v.String), nil
	case db.ColumnTypeTime:
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
