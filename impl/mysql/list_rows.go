package mysql

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

func (o selectOperation) ListRows(ctx context.Context, tx db.Transaction, tableName string, schema db.Schema) (rows []db.Row, err error) {
	rows, err = tx.Read(ctx, fmt.Sprintf(`SELECT * FROM %s ORDER BY %s`, tableName, strings.Join(schema.PrimaryKeys, ", ")), nil)
	if err != nil {
		return nil, err
	}
	out := db.Table{}
	for _, row := range rows {
		outRow := db.Row{}
		for _, columnName := range row.GetColumnNames() {
			col, err := row.GetColumnValue(columnName)
			if err != nil {
				return nil, err
			}
			colBytes, err := col.AsBytes()
			if err != nil {
				return nil, err
			}
			typ, exists := schema.ColumnTypes[columnName]
			if !exists {
				return nil, fmt.Errorf("column %s not found", columnName)
			}
			if colBytes.Valid {
				outRow.SetColumnValue(columnName, string(colBytes.Bytes), typ)
			} else {
				outRow.SetColumnValue(columnName, nil, typ)
			}
		}
		out.Rows = append(out.Rows, outRow)
	}
	return out.Rows, nil
}
