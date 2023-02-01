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

func (schema Schema) GetColumnNames() (columnNames []string) {
	for columnName := range schema.ColumnTypes {
		columnNames = append(columnNames, columnName)
	}
	return columnNames
}
