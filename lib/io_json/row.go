package io_json

import (
	"fmt"
)

type Row map[string]*JsonValue

func (row Row) GetColumnNames() (columnNames []string) {
	for columnName := range row {
		columnNames = append(columnNames, columnName)
	}
	return columnNames
}

func (row Row) Has(columnName string) bool {
	_, exists := row[columnName]
	return exists
}

func (row Row) GetJsonType(columnName string) (jsonType jsonType, err error) {
	val, exists := row[columnName]
	if !exists {
		return "", fmt.Errorf("column %s not found in JsonRow", columnName)
	}
	return val.Type, nil
}

func (row Row) ToString(columnName string) (string, error) {
	val, ok := row[columnName]
	if !ok {
		return "", fmt.Errorf("column %s not found in JsonRow", columnName)
	}
	return val.ToString()
}

func (row Row) ToBool(columnName string) (bool, error) {
	val, ok := row[columnName]
	if !ok {
		return false, fmt.Errorf("column %s not found in JsonRow", columnName)
	}
	return val.ToBool()
}

func (row Row) ToInt64(columnName string) (int64, error) {
	val, ok := row[columnName]
	if !ok {
		return 0, fmt.Errorf("column %s not found in JsonRow", columnName)
	}
	return val.ToInt64()
}

func (row Row) ToFloat64(columnName string) (float64, error) {
	val, ok := row[columnName]
	if !ok {
		return 0, fmt.Errorf("column %s not found in JsonRow", columnName)
	}
	return val.ToFloat64()
}

func (row Row) SetString(columnName string, val string) {
	row[columnName] = NewJsonString(val)
}

func (row Row) SetBool(columnName string, val bool) {
	row[columnName] = NewJsonBoolean(val)
}

func (row Row) SetInt64(columnName string, val int64) {
	row[columnName] = NewJsonNumberInt64(val)
}

func (row Row) SetFloat64(columnName string, val float64) {
	row[columnName] = NewJsonNumberFloat64(val)
}

func (row Row) SetNil(columnName string) {
	row[columnName] = NewJsonNull()
}
