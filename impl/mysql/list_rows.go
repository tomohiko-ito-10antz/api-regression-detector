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

func Select() selectOperation {
	return selectOperation{}
}

var _ cmd.RowLister = selectOperation{}

func (o selectOperation) ListRows(ctx context.Context, tx db.Transaction, table string) (rows db.Table, err error) {
	primaryKeys, err := getPrimaryKeys(ctx, tx, table)
	if err != nil {
		return nil, err
	}
	columnTypes, err := getColumnTypes(ctx, tx, table)
	if err != nil {
		return nil, err
	}
	rows, err = tx.Read(ctx, fmt.Sprintf(`SELECT * FROM %s ORDER BY %s`, table, strings.Join(primaryKeys, ", ")), nil)
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
			typ, exists := columnTypes[columnName]
			if !exists {
				return nil, fmt.Errorf("column %s not found", columnName)
			}
			if colBytes.Valid {
				outRow.SetColumnValue(columnName, string(colBytes.Bytes), typ)
			} else {
				outRow.SetColumnValue(columnName, nil, typ)
			}
		}
		out = append(out, outRow)
	}
	return out, nil
}
