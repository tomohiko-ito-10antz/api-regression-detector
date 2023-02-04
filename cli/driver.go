package cli

import (
	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/impl/mysql"
	"github.com/Jumpaku/api-regression-detector/lib/impl/postgres"
	"github.com/Jumpaku/api-regression-detector/lib/impl/spanner"
	"github.com/Jumpaku/api-regression-detector/lib/impl/sqlite"
)

func NewDriver(name string) (*cmd.Driver, error) {
	var driver *cmd.Driver

	switch name {
	case "mysql":
		driver = &cmd.Driver{
			Name:         name,
			ListRows:     mysql.ListRows(),
			ClearRows:    mysql.ClearRows(),
			CreateRows:   mysql.CreateRows(),
			SchemaGetter: mysql.GetSchema(),
		}
	case "postgres":
		driver = &cmd.Driver{
			Name:         name,
			ListRows:     postgres.ListRows(),
			ClearRows:    postgres.ClearRows(),
			CreateRows:   postgres.CreateRows(),
			SchemaGetter: postgres.GetSchema(),
		}
	case "sqlite3":
		driver = &cmd.Driver{
			Name:         name,
			ListRows:     sqlite.ListRows(),
			ClearRows:    sqlite.ClearRows(),
			CreateRows:   sqlite.Insert(),
			SchemaGetter: sqlite.GetSchema(),
		}
	case "spanner":
		driver = &cmd.Driver{
			Name:         name,
			ListRows:     spanner.ListRows(),
			ClearRows:    spanner.ClearRows(),
			CreateRows:   spanner.CreateRows(),
			SchemaGetter: spanner.GetSchema(),
		}
	default:
		return nil, errors.Wrap(errors.BadArgs, "invalid driver name %s", name)
	}
	return driver, nil
}
