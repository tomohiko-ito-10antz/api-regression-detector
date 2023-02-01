package cli

import (
	"context"
	"os"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/io"
	"go.uber.org/multierr"
)

func RunDump(databaseDriver string, connectionString string) (code int, err error) {
	driver, err := Connect(databaseDriver, connectionString)
	if err != nil {
		return 1, err
	}
	defer func() {
		err = multierr.Combine(err, driver.Close())
		if err != nil {
			code = 1
		}
	}()
	tables, err := io.Load(os.Stdin)
	if err != nil {
		return 1, err
	}
	tableNames := []string{}
	for tableName := range tables {
		tableNames = append(tableNames, tableName)
	}
	dump, err := cmd.Dump(context.Background(), driver.DB, tableNames, driver.SchemaGetter, driver.ListRows)
	if err != nil {
		return 1, err
	}
	err = io.Save(dump, os.Stdout)
	if err != nil {
		return 1, err
	}
	return 0, nil
}
