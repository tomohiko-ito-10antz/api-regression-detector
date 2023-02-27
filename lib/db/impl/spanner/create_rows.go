package spanner

import (
	"context"
	"fmt"
	"strings"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/db/impl"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/tables"
)

type insertOperation struct{}

func CreateRows() insertOperation {
	return insertOperation{}
}

var _ cmd.RowCreator = insertOperation{}

func (o insertOperation) CreateRows(
	ctx context.Context,
	tx db.Tx,
	tableName string,
	schema db.Schema,
	rows []tables.Row,
) error {
	columnTypes := schema.ColumnTypes
	if len(columnTypes) == 0 {
		return nil
	}

	if len(rows) == 0 {
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

			errInfo := errors.Info{"tableName": tableName, "columnName": columnName}

			dbType, ok := columnTypes[columnName]
			if !ok {
				return errors.BadKeyAccess.New(errInfo.AppendTo("column not found in table"))
			}

			errInfo = errInfo.With("dbType", dbType)

			param, err := impl.ExtractColumnValueAsDB(row, columnName, dbType)
			if err != nil {
				return errors.BadConversion.New(
					errInfo.AppendTo("fail to parse column value to DB type"))
			}

			params = append(params, param)
		}

		stmt += ")"
	}

	if err := tx.Write(ctx, stmt, params); err != nil {
		errInfo := errors.Info{"stmt": stmt, "params": params}

		return errors.Wrap(
			errors.DBFailure.Err(err),
			errInfo.AppendTo("fail to insert rows"))
	}

	return nil
}
