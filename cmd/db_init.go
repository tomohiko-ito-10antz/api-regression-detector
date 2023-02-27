package cmd

import (
	"context"
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/cmd/cli"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/tables"
)

func RunInit(stdio *cli.Stdio, databaseDriver string, connectionString string) (code int) {
	errorInfo := errors.Info{"databaseDriver": databaseDriver}

	driver, err := NewDriver(databaseDriver)
	if err != nil {
		PrintError(stdio.Stderr, errors.Wrap(errors.BadArgs.Err(err), errorInfo.AppendTo("fail RunInit")))
		return 1
	}

	errorInfo = errorInfo.With("connectionString", connectionString)

	if err := driver.Open(connectionString); err != nil {
		PrintError(stdio.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail to open database")))
		return 1
	}
	defer func() {
		if err := errors.Join(err, driver.Close()); err != nil {
			PrintError(stdio.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunInit")))
			code = 1
		}
	}()

	initTables, err := tables.LoadInitTables(os.Stdin)
	if err != nil {
		PrintError(stdio.Stderr, errors.Wrap(err, "fail to load init tables from stdin"))
		return 1
	}

	err = cmd.Init(context.Background(), driver.DB, initTables, driver.SchemaGetter, driver.RowClearer, driver.RowCreator)
	if err != nil {
		PrintError(stdio.Stderr, errors.Wrap(err, errorInfo.With("initTables", initTables).AppendTo("fail to init tables in database")))
		return 1
	}

	return 0
}
