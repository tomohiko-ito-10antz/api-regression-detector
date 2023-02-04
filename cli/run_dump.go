package cli

import (
	"context"
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/io_json"
	"go.uber.org/multierr"
)

func RunDump(databaseDriver string, connectionString string) (code int, err error) {
	driver, err := NewDriver(databaseDriver)
	if err != nil {
		return 1, err
	}
	err = driver.Open(connectionString)
	if err != nil {
		return 1, err
	}
	defer func() {
		err = multierr.Combine(err, driver.Close())
		if err != nil {
			code = 1
		}
	}()
	tableNames, err := io_json.LoadJson[[]string](os.Stdin)
	if err != nil {
		return 1, err
	}
	dump, err := cmd.Dump(context.Background(), driver.DB, tableNames, driver.SchemaGetter, driver.ListRows)
	if err != nil {
		return 1, err
	}
	json, err := io_json.TableToJson(dump)
	if err != nil {
		return 1, err
	}
	err = io_json.SaveJson(json, os.Stdout)
	if err != nil {
		return 1, err
	}
	return 0, nil
}
