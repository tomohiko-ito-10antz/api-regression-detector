package sqlite

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

	rows, err := tx.Read(ctx, stmt, nil)
	if err != nil {
		return nil, errors.Wrap(
			errors.DBFailure,
			"fail to select all rows (stmt=%v)", stmt)
	}

	out := db.Table{}
	for _, row := range rows {
		outRow := db.Row{}
		for _, columnName := range schema.GetColumnNames() {
			col, ok := row[columnName]
			if !ok {
				return nil, errors.Wrap(
					errors.BadKeyAccess,
					"column %s not found in schema of table %s", columnName, tableName)
			}

			typ, exists := schema.ColumnTypes[columnName]
			if !exists {
				return nil, errors.Wrap(
					errors.BadKeyAccess,
					"column %s not found in schema of table %s", columnName, tableName)
			}

			outRow[columnName] = col.WithType(typ)
		}

		out.Rows = append(out.Rows, outRow)
	}

	return out.Rows, nil
}
