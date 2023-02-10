package db

import (
	"database/sql"
	"reflect"
	"sort"
	"time"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

type (
	ColumnType  string
	ColumnTypes map[string]ColumnType
)

func (columnTypes ColumnTypes) GetColumnNames() []string {
	columnNames := []string{}
	for columnName := range columnTypes {
		columnNames = append(columnNames, columnName)
	}

	sort.Slice(columnNames, func(i, j int) bool {
		return columnNames[i] < columnNames[j]
	})

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

type (
	NullString  sql.NullString
	NullInteger sql.NullInt64
	NullFloat   sql.NullFloat64
	NullBool    sql.NullBool
	NullTime    sql.NullTime
	NullBytes   struct {
		Bytes []byte
		Valid bool
	}
)

type ColumnValue struct {
	Type  ColumnType
	value any
}

func UnknownTypeColumnValue(val any) *ColumnValue {
	return NewColumnValue(val, ColumnTypeUnknown)
}

func NewColumnValue(val any, typ ColumnType) *ColumnValue {
	return &ColumnValue{value: val, Type: typ}
}

func (v ColumnValue) WithType(typ ColumnType) *ColumnValue {
	v.Type = typ

	return &v
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

	return NullString{}, errors.Wrap(errors.BadConversion, "value %v:%T not compatible to string", v.value, v.value)
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
			return NullInteger{}, errors.Wrap(errors.BadConversion, "value %v:%T not compatible to int64", v.value, v.value)
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
		return NullInteger{Valid: true, Int64: val}, nil
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

	return NullInteger{}, errors.Wrap(errors.BadConversion, "value %v:%T not compatible to int64", v.value, v.value)
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
			return NullFloat{}, errors.Wrap(errors.BadConversion, "value %v:%T not compatible to float64", v.value, v.value)
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
		return NullFloat{Valid: true, Float64: val}, nil
	case sql.NullFloat64:
		return NullFloat(val), nil
	}

	return NullFloat{}, errors.Wrap(errors.BadConversion, "value %v:%T not compatible to float64", v.value, v.value)
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

	return NullBytes{}, errors.Wrap(errors.BadConversion, "value %v:%T not compatible to []byte", v.value, v.value)
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
		return NullTime{}, errors.Wrap(errors.BadConversion, "value %v:%T not compatible to time.Time", v.value, v.value)
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
		return NullBool{}, errors.Wrap(errors.BadConversion, "value %v:%T not compatible to bool", v.value, v.value)
	}
}
