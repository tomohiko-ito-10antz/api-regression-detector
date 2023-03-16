package db

type Schema struct {
	PrimaryKeys []string
	ColumnTypes ColumnTypes
	References  []string
}

type Tables map[string]Table

type Table struct {
	Name   string
	Schema Schema
	Rows   []Row
}

func (schema Schema) GetColumnNames() []string {
	return schema.ColumnTypes.GetColumnNames()
}
