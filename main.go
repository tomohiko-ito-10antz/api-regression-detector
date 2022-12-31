package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/Jumpaku/api-regression-detector/cmd"
	"github.com/Jumpaku/api-regression-detector/io"
	"github.com/Jumpaku/api-regression-detector/log"
	"github.com/Jumpaku/api-regression-detector/mysql"
	"github.com/Jumpaku/api-regression-detector/postgres"
	"github.com/Jumpaku/api-regression-detector/sqlite"
)

func fail(err error) {
	log.Stderr("Error\n%v", err)
	panic(err)
}

type Driver struct {
	Name     string
	DB       *sql.DB
	Select   cmd.Select
	Truncate cmd.Truncate
	Insert   cmd.Insert
}

func (d *Driver) Close() error {
	return d.DB.Close()
}
func Connect(name string, connectionString string) (*Driver, error) {
	switch name {
	case "mysql":
		db, err := sql.Open(name, connectionString)
		if err != nil {
			return nil, err
		}
		return &Driver{
			Name:     name,
			DB:       db,
			Select:   mysql.Select(),
			Truncate: mysql.Truncate(),
			Insert:   mysql.Insert(),
		}, nil
	case "postgres":
		db, err := sql.Open(name, connectionString)
		if err != nil {
			return nil, err
		}
		return &Driver{
			Name:     name,
			DB:       db,
			Select:   postgres.Select(),
			Truncate: postgres.Truncate(),
			Insert:   postgres.Insert(),
		}, nil
	case "sqlite3":
		db, err := sql.Open(name, connectionString)
		if err != nil {
			return nil, err
		}
		return &Driver{
			Name:     name,
			DB:       db,
			Select:   sqlite.Select(),
			Truncate: sqlite.Truncate(),
			Insert:   sqlite.Insert(),
		}, nil
	default:
		return nil, fmt.Errorf("invalid driver name")
	}
}
func main() {
	ctx := context.Background()
	//driverName := "mysql"
	driverName := "postgres"
	//driverName := "sqlite3"
	//connectionString := "root:password@(mysql)/main"
	connectionString := "user=postgres password=password host=postgres dbname=main sslmode=disable"
	//connectionString := "file:examples/sqlite/sqlite.db";
	connection, err := Connect(driverName, connectionString)
	if err != nil {
		fail(err)
	}
	defer connection.Close()
	tables, err := io.Load(os.Stdin)
	if err != nil {
		fail(err)
	}
	err = cmd.Prepare(ctx, connection.DB, tables, connection.Truncate, connection.Insert)
	if err != nil {
		fail(err)
	}
	tableNames := []string{}
	for table := range tables {
		tableNames = append(tableNames, table)
	}
	tables, err = cmd.Dump(ctx, connection.DB, tableNames, connection.Select)
	if err != nil {
		fail(err)
	}
	/*err = io.Save(tables, os.Stdout)
	if err != nil {
		fail(err)
	}*/
}
