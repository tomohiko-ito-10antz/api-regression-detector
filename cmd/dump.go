package cmd

import (
	"context"
	"database/sql"
	"time"

	"github.com/Jumpaku/api-regression-detector/db"
	"github.com/Jumpaku/api-regression-detector/io"
)

func Dump(ctx context.Context, database *sql.DB, tableNames []string, s RowLister) (tables io.Tables, err error) {
	tables = io.Tables{}
	err = db.ExecuteTransaction(ctx, database, func(ctx context.Context, exec db.Transaction) error {
		dbTables := db.Tables{}
		for _, tableName := range tableNames {
			rows, err := s.ListRows(ctx, exec, tableName)
			if err != nil {
				return err
			}
			dbTables[tableName] = rows
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
		jsonRows := io.Table{}
		for _, dbRow := range dbTable.Rows {
			jsonRow := io.Row{}
			for dbColumnName, dbColumnValue := range dbRow {
				jsonRow[dbColumnName], err = convertDBColumnValueToJsonValue(dbColumnValue)
				if err != nil {
					return nil, err
				}
			}
			jsonRows = append(jsonRows, jsonRow)
		}
		jsonTables[dbTableName] = jsonRows
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
