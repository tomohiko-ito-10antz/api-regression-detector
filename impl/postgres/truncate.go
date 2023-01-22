package postgres

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

func (o truncateOperation) Truncate(ctx context.Context, exec db.Exec, table string) (err error) {
	err = exec.Write(ctx, fmt.Sprintf(`TRUNCATE TABLE %s RESTART IDENTITY`, table), nil)
	if err != nil {
		return err
	}
	return nil
}
