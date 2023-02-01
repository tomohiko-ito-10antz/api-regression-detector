package sqlite

import (
	"context"
	"fmt"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/db"
)

type truncateOperation struct {
}

func ClearRows() truncateOperation {
	return truncateOperation{}
}

var _ cmd.RowClearer = truncateOperation{}

func (o truncateOperation) ClearRows(ctx context.Context, tx db.Transaction, table string) (err error) {
	err = tx.Write(ctx, fmt.Sprintf(`DELETE FROM %s`, table), nil)
	if err != nil {
		return err
	}
	err = tx.Write(ctx, `DELETE FROM sqlite_sequence WHERE name = ?`, []any{table})
	if err != nil {
		return err
	}
	return nil
}
