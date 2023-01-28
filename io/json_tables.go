package io

type JsonTable []JsonRow
type JsonTables map[string]JsonTable

func (tables JsonTables) GetTableNames() (tableNames []string) {
	for tableName := range tables {
		tableNames = append(tableNames, tableName)
	}
	return tableNames
}
