package db

import (
	"fmt"
)

type Row map[string]*ColumnValue

func (row Row) GetColumnNames() (columns []string) {
	for column := range row {
		columns = append(columns, column)
	}
	return columns
}
func (row Row) GetColumnValue(columnName string) (*ColumnValue, error) {
	val, exists := row[columnName]
	if !exists {
		return nil, fmt.Errorf("column %s not found", columnName)
	}
	return val, nil
}
func (row Row) SetColumnValue(columnName string, val any, typ ColumnType) {
	row[columnName] = NewColumnValue(val)
	row[columnName].Type = typ
}
