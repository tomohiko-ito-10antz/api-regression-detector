package cmd

import (
	"context"
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/cmd/cli"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/tables"
)

func RunDump(stdio *cli.Stdio, databaseDriver string, connectionString string) (code int) {
	errorInfo := errors.Info{"databaseDriver": databaseDriver}

	driver, err := NewDriver(databaseDriver)
	if err != nil {
		PrintError(stdio.Stderr, errors.Wrap(errors.BadArgs.Err(err), errorInfo.AppendTo("fail RunDump")))
		return 1
	}

	errorInfo = errorInfo.With("connectionString", connectionString)

	if err := driver.Open(connectionString); err != nil {
		PrintError(stdio.Stderr, errors.Wrap(errors.IOFailure.Err(err), errorInfo.AppendTo("fail to open database")))
		return 1
	}
	defer func() {
		if err := errors.Join(err, driver.Close()); err != nil {
			PrintError(os.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunDump")))
			code = 1
		}
	}()

	tableNames, err := jsonio.LoadJson[[]string](os.Stdin)
	if err != nil {
		PrintError(os.Stderr, errors.Wrap(err, "fail to load init tables as JSON from stdin"))
		return 1
	}

	dump, err := cmd.Dump(context.Background(), driver.DB, tableNames, driver.SchemaGetter, driver.RowLister)
	if err != nil {
		PrintError(os.Stderr, errors.Wrap(err, errorInfo.With("tableNames", tableNames).AppendTo("fail to get tables from database")))
		return 1
	}

	if err := tables.SaveDumpTables(dump, os.Stdout); err != nil {
		PrintError(os.Stderr, errors.Wrap(err, "fail to dump tables as JSON to stdout"))
		return 1
	}

	return 0
}
