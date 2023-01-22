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

var _ cmd.Select = selectOperation{}

func (o selectOperation) Select(ctx context.Context, tx db.Exec, table string) (rows db.Rows, err error) {
	primaryKeys, err := getPrimaryKeys(ctx, tx, table)
	if err != nil {
		return nil, err
	}
	rows, err = tx.Read(ctx, fmt.Sprintf(`SELECT * FROM %s ORDER BY %s`, table, strings.Join(primaryKeys, ", ")), nil)
	if err != nil {
		return nil, err
	}
	out := db.Rows{}
	for i, row := range rows {
		out = append(out, db.Row{})
		for _, col := range row.GetColumns() {
			isNull, err := row.IsNull(col)
			if err != nil {
				return nil, err
			}
			if isNull {
				out[i][col] = nil
			} else {
				byteArray, err := rows[i].GetByteArray(col)
				if err != nil {
					return nil, err
				}
				out[i][col] = string(byteArray)
			}
		}
	}
	return out, nil
}
