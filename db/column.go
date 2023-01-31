package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
)

type ColumnType string
type ColumnTypes map[string]ColumnType

func (cts ColumnTypes) ColumnNames() []string {
	columnNames := make([]string, 0, len(cts))
	for columnName := range cts {
		columnNames = append(columnNames, columnName)
	}
	return columnNames
}

const (
	ColumnTypeUnknown ColumnType = "UNKNOWN"
	ColumnTypeBoolean ColumnType = "BOOL"
	ColumnTypeInteger ColumnType = "INTEGER"
	ColumnTypeFloat   ColumnType = "FLOAT"
	ColumnTypeString  ColumnType = "STRING"
	ColumnTypeTime    ColumnType = "TIME"
)

type NullString sql.NullString
type NullInteger sql.NullInt64
type NullFloat sql.NullFloat64
type NullBool sql.NullBool
type NullBytes struct {
	Bytes []byte
	Valid bool
}
type NullTime sql.NullTime

type ColumnValue struct {
	Type  ColumnType
	value any
}

func NewColumnValue(val any) *ColumnValue {
	return &ColumnValue{value: val, Type: ColumnTypeUnknown}
}

func (v ColumnValue) AsString() (NullString, error) {
	switch val := v.value.(type) {
	case nil:
		return NullString{}, nil
	case *string:
		if val == nil {
			return NullString{}, nil
		}
		return NullString{Valid: true, String: *val}, nil
	case string:
		return NullString{Valid: true, String: val}, nil
	case sql.NullString:
		return NullString(val), nil
	}
	return NullString{}, fmt.Errorf("value not compatible to string")
}

func (v ColumnValue) AsInteger() (NullInteger, error) {
	rv := reflect.ValueOf(v.value)
	if !rv.IsValid() {
		return NullInteger{}, nil
	}
	val := v.value
	if rv.Kind() == reflect.Pointer {
		switch rv.Type().Elem().Kind() {
		default:
			return NullInteger{}, fmt.Errorf("value not compatible to int64")
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		}
		if rv.IsNil() {
			return NullInteger{}, nil
		}
		val = rv.Elem().Interface()
	}
	switch val := val.(type) {
	case int:
		return NullInteger{Valid: true, Int64: int64(val)}, nil
	case int8:
		return NullInteger{Valid: true, Int64: int64(val)}, nil
	case int16:
		return NullInteger{Valid: true, Int64: int64(val)}, nil
	case int32:
		return NullInteger{Valid: true, Int64: int64(val)}, nil
	case int64:
		return NullInteger{Valid: true, Int64: int64(val)}, nil
	case uint:
		return NullInteger{Valid: true, Int64: int64(val)}, nil
	case uint8:
		return NullInteger{Valid: true, Int64: int64(val)}, nil
	case uint16:
		return NullInteger{Valid: true, Int64: int64(val)}, nil
	case uint32:
		return NullInteger{Valid: true, Int64: int64(val)}, nil
	case uint64:
		return NullInteger{Valid: true, Int64: int64(val)}, nil
	case sql.NullByte:
		return NullInteger{Valid: val.Valid, Int64: int64(val.Byte)}, nil
	case sql.NullInt16:
		return NullInteger{Valid: val.Valid, Int64: int64(val.Int16)}, nil
	case sql.NullInt32:
		return NullInteger{Valid: val.Valid, Int64: int64(val.Int32)}, nil
	case sql.NullInt64:
		return NullInteger(val), nil
	}
	return NullInteger{}, fmt.Errorf("value not compatible to int64")
}

func (v ColumnValue) AsFloat() (NullFloat, error) {
	rv := reflect.ValueOf(v.value)
	if !rv.IsValid() {
		return NullFloat{}, nil
	}
	val := v.value
	if rv.Kind() == reflect.Pointer {
		switch rv.Type().Elem().Kind() {
		default:
			return NullFloat{}, fmt.Errorf("value not compatible to float64")
		case reflect.Float32, reflect.Float64:
		}
		if rv.IsNil() {
			return NullFloat{}, nil
		}
		val = rv.Elem().Interface()
	}
	switch val := val.(type) {
	case float32:
		return NullFloat{Valid: true, Float64: float64(val)}, nil
	case float64:
		return NullFloat{Valid: true, Float64: float64(val)}, nil
	case sql.NullFloat64:
		return NullFloat(val), nil
	}
	return NullFloat{}, fmt.Errorf("value not compatible to float64")
}

func (v ColumnValue) AsBytes() (NullBytes, error) {
	switch val := v.value.(type) {
	case nil:
		return NullBytes{}, nil
	case *[]byte:
		if val == nil {
			return NullBytes{}, nil
		}
		return NullBytes{Valid: true, Bytes: *val}, nil
	case []byte:
		if val == nil {
			return NullBytes{}, nil
		}
		return NullBytes{Valid: true, Bytes: val}, nil
	}
	return NullBytes{}, fmt.Errorf("value not compatible to []byte")
}

func (v ColumnValue) AsTime() (NullTime, error) {
	switch val := v.value.(type) {
	case nil:
		return NullTime{}, nil
	case *time.Time:
		if val == nil {
			return NullTime{}, nil
		}
		return NullTime{Valid: true, Time: *val}, nil
	case time.Time:
		return NullTime{Valid: true, Time: val}, nil
	case sql.NullTime:
		return NullTime(val), nil
	default:
		return NullTime{}, fmt.Errorf("value not compatible to time.Time")
	}
}

func (v ColumnValue) AsBool() (NullBool, error) {
	switch val := v.value.(type) {
	case nil:
		return NullBool{}, nil
	case *bool:
		if val == nil {
			return NullBool{}, nil
		}
		return NullBool{Valid: true, Bool: *val}, nil
	case bool:
		return NullBool{Valid: true, Bool: val}, nil
	case sql.NullBool:
		return NullBool(val), nil
	default:
		return NullBool{}, fmt.Errorf("value not compatible to time.Time")
	}
}
