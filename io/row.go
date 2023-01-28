package io

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

func (row Row) GetColumnType(columnName string) (jsonType jsonType, err error) {
	val, exists := row[columnName]
	if !exists {
		return JsonTypeUnknown, fmt.Errorf("column not found in JsonRow")
	}
	return val.Type, nil
}

func (row Row) ToString(columnName string) (string, error) {
	val, ok := row[columnName]
	if !ok {
		return "", fmt.Errorf("column not found in JsonRow")
	}
	return val.ToString()
}

func (row Row) ToBool(columnName string) (bool, error) {
	val, ok := row[columnName]
	if !ok {
		return false, fmt.Errorf("column not found in JsonRow")
	}
	return val.ToBool()
}

func (row Row) ToInt64(columnName string) (int64, error) {
	val, ok := row[columnName]
	if !ok {
		return 0, fmt.Errorf("column not found in JsonRow")
	}
	return val.ToInt64()
}

func (row Row) ToFloat64(columnName string) (float64, error) {
	val, ok := row[columnName]
	if !ok {
		return 0, fmt.Errorf("column not found in JsonRow")
	}
	return val.ToFloat64()
}

func (row Row) SetString(columnName string, val string) (err error) {
	row[columnName], err = NewJson(val)
	return err
}

func (row Row) SetBool(columnName string, val bool) (err error) {
	row[columnName], err = NewJson(val)
	return err
}

func (row Row) SetInt64(columnName string, val int64) (err error) {
	row[columnName], err = NewJson(val)
	return err
}

func (row Row) SetFloat64(columnName string, val float64) (err error) {
	row[columnName], err = NewJson(val)
	return err
}

func (row Row) SetNil(columnName string) (err error) {
	row[columnName], err = NewJson(nil)
	return err
}
