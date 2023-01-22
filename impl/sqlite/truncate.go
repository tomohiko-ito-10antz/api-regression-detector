package sqlite

import (
	"context"
	"fmt"

	"github.com/Jumpaku/api-regression-detector/cmd"
	"github.com/Jumpaku/api-regression-detector/db"
)

type truncateOperation struct {
}

func Truncate() truncateOperation {
	return truncateOperation{}
}

var _ cmd.Truncate = truncateOperation{}

func (o truncateOperation) Truncate(ctx context.Context, tx db.Exec, table string) (err error) {
	err = tx.Write(ctx, fmt.Sprintf(`DELETE FROM %s`, table), nil)
	if err != nil {
		return err
	}
	err = tx.Write(ctx, fmt.Sprintf(`DELETE FROM sqlite_sequence WHERE name = '%s'`, table), nil)
	if err != nil {
		return err
	}
	return nil
}
