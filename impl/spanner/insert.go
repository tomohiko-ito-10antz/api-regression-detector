package spanner

import (
	"context"
	"fmt"
	"strings"

	"github.com/Jumpaku/api-regression-detector/cmd"
	"github.com/Jumpaku/api-regression-detector/db"
	"github.com/Jumpaku/api-regression-detector/json"
)

type insertOperation struct{}

func Insert() insertOperation {
	return insertOperation{}
}

var _ cmd.Insert = insertOperation{}

func (o insertOperation) Insert(ctx context.Context, tx db.Exec, table string, rows db.Rows) (err error) {
	columnTypes, err := getColumnTypes(ctx, tx, table)
	if err != nil {
		return err
	}
	columns := getColumns(rows)
	if len(columns) == 0 {
		return nil
	}
	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES", table, strings.Join(columns, ", "))
	values := []any{}
	for i, row := range rows {
		if i > 0 {
			stmt += ","
		}
		stmt += "("
		for j, column := range columns {
			if j > 0 {
				stmt += ","
			}
			stmt += "?"
			columnType, ok := columnTypes[column]
			if !ok {
				return err
			}
			value, ok := row[column]
			if !ok {
				value = nil
			}
			switch columnType {
			case db.ColumnTypeBoolean:
				v, err := json.ToNullableBoolean(value)
				if err != nil {
					return err
				}
				values = append(values, v)
			case db.ColumnTypeInteger:
				v, err := json.ToNullableInteger(value)
				if err != nil {
					return err
				}
				values = append(values, v)
			case db.ColumnTypeFloat:
				v, err := json.ToNullableFloat(value)
				if err != nil {
					return err
				}
				values = append(values, v)
			default:
				v, err := json.ToNullableString(value)
				if err != nil {
					return err
				}
				values = append(values, v)
			}
		}
		stmt += ")"
	}
	err = tx.Write(ctx, stmt, values)
	if err != nil {
		return err
	}
	return nil
}
