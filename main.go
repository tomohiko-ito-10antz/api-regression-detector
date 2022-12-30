package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Jumpaku/api-regression-detector/mysql"
	"github.com/Jumpaku/api-regression-detector/postgres"
	"github.com/Jumpaku/api-regression-detector/prepare"
	"github.com/Jumpaku/api-regression-detector/sqlite"
)

func main() {
	tables, err := prepare.Load(os.Stdin)
	if err != nil {
		log.Fatalln(err.Error())
	}
	driver := "mysql"
	ctx := context.Background()
	db, err := sql.Open(driver, "main:password@(mysql)/main")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.Close()
	switch driver {
	case "mysql":
		err = prepare.Prepare(ctx, db, tables, mysql.Truncate(), mysql.Insert())
	case "postgres":
		err = prepare.Prepare(ctx, db, tables, postgres.Truncate(), postgres.Insert())
	case "sqlite":
		err = prepare.Prepare(ctx, db, tables, sqlite.Truncate(), sqlite.Insert())
	default:
		log.Fatalln("driver specified")
	}
	if err != nil {
		log.Fatalln(err)
	}
}
