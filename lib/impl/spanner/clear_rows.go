package spanner

import (
	"context"
	"fmt"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/db"
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
