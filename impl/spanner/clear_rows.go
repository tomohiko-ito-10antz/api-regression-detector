package spanner

import (
	"context"
	"fmt"

	"github.com/Jumpaku/api-regression-detector/cmd"
	"github.com/Jumpaku/api-regression-detector/db"
)

type truncateOperation struct{}

func ClearRows() truncateOperation {
	return truncateOperation{}
}

var _ cmd.RowClearer = truncateOperation{}

func (o truncateOperation) ClearRows(ctx context.Context, tx db.Transaction, table string) (err error) {
	err = tx.Write(ctx, fmt.Sprintf(`DELETE FROM %s WHERE TRUE`, table), nil)
	if err != nil {
		return err
	}
	return nil
}
