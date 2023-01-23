package io

type JsonColumnValue any
type JsonRow map[string]JsonColumnValue
type JsonRows []JsonRow
type JsonTables map[string]JsonRows

func (tables JsonTables) GetTableNames() (tableNames []string) {
	for tableName := range tables {
		tableNames = append(tableNames, tableName)
	}
	return tableNames
}

func (row JsonRow) GetColumnNames() (columnNames []string) {
	for columnName := range row {
		columnNames = append(columnNames, columnName)
	}
	return columnNames
}
