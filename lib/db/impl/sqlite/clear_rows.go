package sqlite

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
	err := tx.Write(ctx, fmt.Sprintf(`DELETE FROM %s`, tableName), nil)
	if err != nil {
		return errors.Wrap(errors.Join(err, errors.DBFailure), "fail to delete all rows in table %s", tableName)
	}

	err = tx.Write(ctx, `DELETE FROM sqlite_sequence WHERE name = ?`, []any{tableName})
	if err != nil {
		return errors.Wrap(errors.Join(err, errors.DBFailure), "fail to reset auto inclement in table %s", tableName)
	}

	return nil
}
