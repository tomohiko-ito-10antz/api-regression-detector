package spanner

import (
	"context"
	"fmt"
	"strings"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/impl"
	"github.com/Jumpaku/api-regression-detector/lib/io"
)

type insertOperation struct{}

func CreateRows() insertOperation {
	return insertOperation{}
}

var _ cmd.RowCreator = insertOperation{}

func (o insertOperation) CreateRows(ctx context.Context, tx db.Transaction, tableName string, schema db.Schema, rows []io.Row) (err error) {
	columnTypes := schema.ColumnTypes
	if len(columnTypes) == 0 {
		return nil
	}
	columnNames := columnTypes.GetColumnNames()

	stmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES", tableName, strings.Join(columnNames, ", "))
	params := []any{}
	for i, row := range rows {
		if i > 0 {
			stmt += ","
		}
		stmt += "("
		for j, columnName := range columnNames {
			if j > 0 {
				stmt += ","
			}
			stmt += "?"
			dbType, ok := columnTypes[columnName]
			if !ok {
				return err
			}
			param, err := impl.ExtractColumnValueAsDB(row, columnName, dbType)
			if err != nil {
				return err
			}
			params = append(params, param)
		}
		stmt += ")"
	}
	err = tx.Write(ctx, stmt, params)
	if err != nil {
		return err
	}
	return nil
}
