package spanner

import (
	"context"
	"fmt"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

type truncateOperation struct{}

func ClearRows() truncateOperation {
	return truncateOperation{}
}

var _ cmd.RowClearer = truncateOperation{}

func (o truncateOperation) ClearRows(ctx context.Context, tx db.Tx, tableName string) error {
	err := tx.Write(ctx, fmt.Sprintf(`DELETE FROM %s WHERE TRUE`, tableName), nil)
	if err != nil {
		return errors.Wrap(errors.Join(err, errors.DBFailure), "fail to delete all rows in table %s", tableName)
	}

	return nil
}
