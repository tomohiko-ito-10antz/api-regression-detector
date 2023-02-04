package sqlite

import (
	"context"
	"fmt"
	"strings"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/impl"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
)

type insertOperation struct{}

func Insert() insertOperation {
	return insertOperation{}
}

var _ cmd.RowCreator = insertOperation{}

func (o insertOperation) CreateRows(
	ctx context.Context,
	tx db.Tx,
	tableName string,
	schema db.Schema,
	rows []jsonio.Row,
) error {
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

			dbType, exists := columnTypes[columnName]
			if !exists {
				return errors.Wrap(
					errors.BadKeyAccess,
					"column %s not found in table", columnName, tableName)
			}

			param, err := impl.ExtractColumnValueAsDB(row, columnName, dbType)
			if err != nil {
				return errors.Wrap(
					errors.BadConversion,
					"fail to convert value of column %s to type %v", columnName, dbType)
			}

			params = append(params, param)
		}

		stmt += ")"
	}

	if err := tx.Write(ctx, stmt, params); err != nil {
		return errors.Wrap(
			errors.BadConversion,
			"fail to insert rows (stmt=%v,params=%v)", stmt, params)
	}

	return nil
}
