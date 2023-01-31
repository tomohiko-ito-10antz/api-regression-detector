package io

type Table struct {
	Rows []Row
}
type Tables map[string]Table

func (tables Tables) GetTableNames() (tableNames []string) {
	for tableName := range tables {
		tableNames = append(tableNames, tableName)
	}
	return tableNames
}
