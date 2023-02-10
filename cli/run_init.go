package cli

import (
	"context"
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/tables"
)

func RunInit(databaseDriver string, connectionString string) (code int, err error) {
	driver, err := NewDriver(databaseDriver)
	if err != nil {
		return 1, errors.Wrap(err, "fail to new %s", databaseDriver)
	}

	err = driver.Open(connectionString)
	if err != nil {
		return 1, errors.Wrap(err, "fail to connect %s", connectionString)
	}

	defer func() {
		err = errors.Wrap(errors.Join(err, driver.Close()), "fail RunInit")
		if err != nil {
			code = 1
		}
	}()

	initTables, err := tables.LoadInitTables(os.Stdin)
	if err != nil {
		return 1, errors.Wrap(err, "fail to load JSON from stdin")
	}

	err = cmd.Init(context.Background(), driver.DB, initTables, driver.SchemaGetter, driver.RowClearer, driver.RowCreator)
	if err != nil {
		return 1, errors.Wrap(err, "fail Init")
	}

	return 0, nil
}
