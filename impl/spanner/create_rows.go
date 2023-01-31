package spanner

import (
	"context"
	"fmt"
	"strings"

	"github.com/Jumpaku/api-regression-detector/cmd"
	"github.com/Jumpaku/api-regression-detector/db"
	"github.com/Jumpaku/api-regression-detector/impl"
	"github.com/Jumpaku/api-regression-detector/io"
)

type insertOperation struct{}

func CreateRows() insertOperation {
	return insertOperation{}
}

var _ cmd.RowCreator = insertOperation{}

func (o insertOperation) CreateRows(ctx context.Context, tx db.Transaction, tableName string, rows []io.Row) (err error) {
	columnTypes, err := getColumnTypes(ctx, tx, tableName)
	if err != nil {
		return err
	}
	if len(columnTypes) == 0 {
		return nil
	}
	columnNames := columnTypes.ColumnNames()

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
