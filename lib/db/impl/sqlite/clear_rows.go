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
	errInfo := errors.Info{tableName: tableName}

	err := tx.Write(ctx, fmt.Sprintf(`DELETE FROM %s`, tableName), nil)
	if err != nil {
		return errors.Wrap(
			errors.DBFailure.Err(err),
			errInfo.AppendTo("fail to delete all rows in table"))
	}

	err = tx.Write(ctx, `DELETE FROM sqlite_sequence WHERE name = ?`, []any{tableName})
	if err != nil {
		return errors.Wrap(
			errors.DBFailure.Err(err),
			errInfo.AppendTo("fail to reset auto inclement in table"))
	}

	return nil
}
