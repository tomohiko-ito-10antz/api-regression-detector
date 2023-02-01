package mysql

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
	err = tx.Write(ctx, fmt.Sprintf(`TRUNCATE TABLE %s`, table), nil)
	if err != nil {
		return err
	}
	err = tx.Write(ctx, fmt.Sprintf(`ALTER TABLE %s AUTO_INCREMENT = 1`, table), nil)
	if err != nil {
		return err
	}
	return nil
}
