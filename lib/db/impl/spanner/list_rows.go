package spanner

import (
	"context"
	"fmt"
	"strings"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

type selectOperation struct{}

func ListRows() selectOperation {
	return selectOperation{}
}

var _ cmd.RowLister = selectOperation{}

func (o selectOperation) ListRows(
	ctx context.Context,
	tx db.Tx,
	tableName string,
	schema db.Schema,
) ([]db.Row, error) {
	stmt := fmt.Sprintf(`SELECT * FROM %s ORDER BY %s`, tableName, strings.Join(schema.PrimaryKeys, ", "))

	errInfo := errors.Info{"stmt": stmt}

	rows, err := tx.Read(ctx, stmt, nil)
	if err != nil {
		return nil, errors.Wrap(
			errors.DBFailure.Err(err),
			errInfo.AppendTo("fail to select all rows"))
	}

	out := db.Table{}
	for _, row := range rows {
		outRow := db.Row{}
		for _, columnName := range schema.GetColumnNames() {
			errInfo := errInfo.With("columnName", columnName)

			col, ok := row[columnName]
			if !ok {
				return nil, errors.BadKeyAccess.New(
					errInfo.AppendTo("column not found in row"))
			}

			typ, exists := schema.ColumnTypes[columnName]
			if !exists {
				return nil, errors.BadKeyAccess.New(
					errInfo.AppendTo("column not found in table schema"))
			}

			outRow[columnName] = col.WithType(typ)
		}

		out.Rows = append(out.Rows, outRow)
	}
	return out.Rows, nil
}
