package cli

import (
	"context"
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/io_json"
	"go.uber.org/multierr"
)

func RunInit(databaseDriver string, connectionString string) (code int, err error) {
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
	json, err := io_json.LoadJson[map[string][]map[string]any](os.Stdin)
	if err != nil {
		return 1, err
	}
	tables, err := io_json.TableFromJson(json)
	if err != nil {
		return 1, err
	}
	err = cmd.Init(context.Background(), driver.DB, tables, driver.SchemaGetter, driver.ClearRows, driver.CreateRows)
	if err != nil {
		return 1, err
	}
	return 0, nil
}
