package cli

import (
	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/impl/mysql"
	"github.com/Jumpaku/api-regression-detector/lib/impl/postgres"
	"github.com/Jumpaku/api-regression-detector/lib/impl/spanner"
	"github.com/Jumpaku/api-regression-detector/lib/impl/sqlite"
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
		return nil, errors.Wrap(errors.BadArgs, "invalid driver name %s", name)
	}

	return driver, nil
}
