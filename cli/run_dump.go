package cli

import (
	"context"
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
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
		err = errors.Wrap(errors.Join(err, driver.Close(), errors.IOFailure), "fail RunDump")
		if err != nil {
			code = 1
		}
	}()

	tableNames, err := jsonio.LoadJson[[]string](os.Stdin)
	if err != nil {
		return 1, errors.Wrap(err, "fail to load JSON from stdin")
	}

	dump, err := cmd.Dump(context.Background(), driver.DB, tableNames, driver.SchemaGetter, driver.ListRows)
	if err != nil {
		return 1, errors.Wrap(err, "fail Dump")
	}

	json, err := jsonio.TableToJson(dump)
	if err != nil {
		return 1, errors.Wrap(err, "fail to convert to JSON")
	}

	err = jsonio.SaveJson(json, os.Stdout)
	if err != nil {
		return 1, errors.Wrap(err, "fail to save JSON to stdout")
	}

	return 0, nil
}
