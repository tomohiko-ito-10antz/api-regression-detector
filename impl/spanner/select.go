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

func Select() selectOperation {
	return selectOperation{}
}

var _ cmd.Select = selectOperation{}

func (o selectOperation) Select(ctx context.Context, tx db.Exec, table string) (rows db.Rows, err error) {
	primaryKeys, err := getPrimaryKeys(ctx, tx, table)
	if err != nil {
		return nil, err
	}
	rows, err = tx.Read(ctx, fmt.Sprintf(`SELECT * FROM %s ORDER BY %v`, table, strings.Join(primaryKeys, ", ")), nil)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
