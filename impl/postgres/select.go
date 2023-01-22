package postgres

import (
	"context"
	"fmt"

	"github.com/Jumpaku/api-regression-detector/cmd"
	"github.com/Jumpaku/api-regression-detector/db"
)

type selectOperation struct {
}

func Select() selectOperation {
	return selectOperation{}
}

var _ cmd.Select = selectOperation{}

func (o selectOperation) Select(ctx context.Context, exec db.Exec, table string) (rows db.Rows, err error) {
	rows, err = exec.Read(ctx, fmt.Sprintf(`SELECT * FROM %s`, table), nil)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
