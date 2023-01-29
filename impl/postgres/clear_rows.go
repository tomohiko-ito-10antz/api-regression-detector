package postgres

import (
	"context"
	"fmt"

	"github.com/Jumpaku/api-regression-detector/cmd"
	"github.com/Jumpaku/api-regression-detector/db"
)

type truncateOperation struct {
}

func ClearRows() truncateOperation {
	return truncateOperation{}
}

var _ cmd.RowClearer = truncateOperation{}

func (o truncateOperation) ClearRows(ctx context.Context, tx db.Transaction, table string) (err error) {
	err = tx.Write(ctx, fmt.Sprintf(`TRUNCATE TABLE %s RESTART IDENTITY`, table), nil)
	if err != nil {
		return err
	}
	return nil
}
