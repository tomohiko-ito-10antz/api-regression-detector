package cmd

import (
	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/db/impl/mysql"
	"github.com/Jumpaku/api-regression-detector/lib/db/impl/postgres"
	"github.com/Jumpaku/api-regression-detector/lib/db/impl/spanner"
	"github.com/Jumpaku/api-regression-detector/lib/db/impl/sqlite"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

func NewDriver(name string) (*cmd.DatabaseDriver, error) {
	var driver *cmd.DatabaseDriver

	switch name {
	case "mysql":
		driver = &cmd.DatabaseDriver{
			Name:         name,
			RowLister:    mysql.ListRows(),
			RowClearer:   mysql.ClearRows(),
			RowCreator:   mysql.CreateRows(),
			SchemaGetter: mysql.GetSchema(),
		}
	case "postgres":
		driver = &cmd.DatabaseDriver{
			Name:         name,
			RowLister:    postgres.ListRows(),
			RowClearer:   postgres.ClearRows(),
			RowCreator:   postgres.CreateRows(),
			SchemaGetter: postgres.GetSchema(),
		}
	case "sqlite3":
		driver = &cmd.DatabaseDriver{
			Name:         name,
			RowLister:    sqlite.ListRows(),
			RowClearer:   sqlite.ClearRows(),
			RowCreator:   sqlite.Insert(),
			SchemaGetter: sqlite.GetSchema(),
		}
	case "spanner":
		driver = &cmd.DatabaseDriver{
			Name:         name,
			RowLister:    spanner.ListRows(),
			RowClearer:   spanner.ClearRows(),
			RowCreator:   spanner.CreateRows(),
			SchemaGetter: spanner.GetSchema(),
		}
	default:
		return nil, errors.WithTag(errors.New(errors.Info{"name": name}.AppendTo("invalid driver name")),
			errors.BadArgs, errors.Unsupported)
	}

	return driver, nil
}
