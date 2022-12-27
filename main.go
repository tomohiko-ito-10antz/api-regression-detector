package main

import (
	"log"
	"os"

	"github.com/Jumpaku/api-regression-detector/prepare"
	"github.com/Jumpaku/api-regression-detector/prepare/mysql"
	"github.com/Jumpaku/api-regression-detector/prepare/postgres"
	"github.com/Jumpaku/api-regression-detector/prepare/sqlite"
)

func main() {
	tables, err := prepare.ReadTablesFrom(os.Stdin)
	if err != nil {
		log.Fatalln(err.Error())
	}
	sql := ""
	db := "postgres"
	switch db {
	case "postgres":
		sql = postgres.Build(tables)
	case "mysql":
		sql = mysql.Build(tables)
	case "sqlite":
		sql = sqlite.Build(tables)
	default:
		log.Fatalln("no database specified")
	}
	err = prepare.WriteSqlTo(sql, os.Stdout)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
