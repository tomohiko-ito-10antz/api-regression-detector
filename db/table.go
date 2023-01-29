package db

type Schema struct {
	PrimaryKeys []string
	ColumnTypes ColumnTypes
}

type Tables map[string]Table

type Table struct {
	Name   string
	Schema Schema
	Rows   []Row
}
