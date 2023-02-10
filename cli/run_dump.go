package cli

import (
	"context"
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/tables"
)

func RunDump(databaseDriver string, connectionString string) (code int, err error) {
	driver, err := NewDriver(databaseDriver)
	if err != nil {
		return 1, errors.Wrap(errors.Join(err, errors.BadArgs), "fail RunDump")
	}

	err = driver.Open(connectionString)
	if err != nil {
		return 1, errors.Wrap(errors.Join(err, errors.IOFailure), "fail RunDump")
	}

	defer func() {
		err = errors.Wrap(errors.Join(err, driver.Close()), "fail RunDump")
		if err != nil {
			code = 1
		}
	}()

	tableNames, err := jsonio.LoadJson[[]string](os.Stdin)
	if err != nil {
		return 1, errors.Wrap(err, "fail to load init tables as JSON from stdin")
	}

	dump, err := cmd.Dump(context.Background(), driver.DB, tableNames, driver.SchemaGetter, driver.RowLister)
	if err != nil {
		return 1, errors.Wrap(err, "fail Dump")
	}

	if err := tables.SaveDumpTables(dump, os.Stdout); err != nil {
		return 1, errors.Wrap(err, "fail to convert dump tables as JSON to stdout")
	}

	return 0, nil
}
