package spanner

import (
	"context"
	"fmt"
	"strings"

	"github.com/Jumpaku/api-regression-detector/cmd"
	"github.com/Jumpaku/api-regression-detector/db"
)

type selectOperation struct {
}

func ListRows() selectOperation {
	return selectOperation{}
}

var _ cmd.RowLister = selectOperation{}

func (o selectOperation) ListRows(ctx context.Context, tx db.Transaction, tableName string, schema db.Schema) (table db.Table, err error) {
	rows, err := tx.Read(ctx, fmt.Sprintf(`SELECT * FROM %s ORDER BY %s`, tableName, strings.Join(schema.PrimaryKeys, ", ")), nil)
	if err != nil {
		return db.Table{}, err
	}
	out := db.Table{}
	for _, row := range rows {
		outRow := db.Row{}
		for _, columnName := range row.GetColumnNames() {
			col, err := row.GetColumnValue(columnName)
			if err != nil {
				return db.Table{}, err
			}
			typ, exists := schema.ColumnTypes[columnName]
			if !exists {
				return db.Table{}, fmt.Errorf("column %s not found", columnName)
			}
			outRow.SetColumnValue(columnName, col, typ)
		}
		out.Rows = append(out.Rows, outRow)
	}
	return out, nil
}
