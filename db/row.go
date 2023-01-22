package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
)

type Row map[string]any
type Rows []Row

func (row Row) GetColumns() (columns []string) {
	for column := range row {
		columns = append(columns, column)
	}
	return columns
}
func (row Row) GetString(column string) (val string, err error) {
	nullable, err := row.GetNullString(column)
	if err != nil {
		return "", err
	}
	if !nullable.Valid {
		return "", fmt.Errorf("value is nil")
	}
	return nullable.String, nil
}

func (row Row) IsNull(column string) (isNull bool, err error) {
	valAny, exists := row[column]
	if !exists {
		return false, fmt.Errorf("column not found")
	}
	rv := reflect.ValueOf(valAny)
	if !rv.IsValid() || rv.IsNil() {
		return true, nil
	}
	switch val := valAny.(type) {
	case sql.NullBool:
		return !val.Valid, nil
	case sql.NullByte:
		return !val.Valid, nil
	case sql.NullInt16:
		return !val.Valid, nil
	case *sql.NullInt32:
		return !val.Valid, nil
	case sql.NullInt32:
		return !val.Valid, nil
	case sql.NullInt64:
		return !val.Valid, nil
	case sql.NullFloat64:
		return !val.Valid, nil
	case sql.NullString:
		return !val.Valid, nil
	case sql.NullTime:
		return !val.Valid, nil
	default:
		return false, nil
	}
}
func (row Row) GetNullString(column string) (val sql.NullString, err error) {
	isNull, err := row.IsNull(column)
	if err != nil {
		return sql.NullString{}, err
	}
	if isNull {
		return sql.NullString{}, nil
	}
	valAny, exists := row[column]
	if !exists {
		return sql.NullString{}, fmt.Errorf("column not found")
	}
	switch val := valAny.(type) {
	case *string:
		return sql.NullString{String: *val, Valid: true}, nil
	case string:
		return sql.NullString{String: val, Valid: true}, nil
	case sql.NullString:
		return val, nil
	case sql.NullTime:
		return sql.NullString{String: val.Time.Format(time.RFC3339), Valid: true}, nil
	default:
		return sql.NullString{}, fmt.Errorf("value is not string")
	}
}
func (row Row) GetInteger(column string) (val int64, err error) {
	nullable, err := row.GetNullInteger(column)
	if err != nil {
		return 0, err
	}
	if !nullable.Valid {
		return 0, fmt.Errorf("value is nil")
	}
	return nullable.Int64, nil
}
func (row Row) GetNullInteger(column string) (val sql.NullInt64, err error) {
	isNull, err := row.IsNull(column)
	if err != nil {
		return sql.NullInt64{}, err
	}
	if isNull {
		return sql.NullInt64{}, nil
	}
	valAny, exists := row[column]
	if !exists {
		return sql.NullInt64{}, fmt.Errorf("column not found")
	}
	switch val := valAny.(type) {
	case *int:
		return sql.NullInt64{Int64: int64(*val), Valid: true}, nil
	case *int8:
		return sql.NullInt64{Int64: int64(*val), Valid: true}, nil
	case *int16:
		return sql.NullInt64{Int64: int64(*val), Valid: true}, nil
	case *int32:
		return sql.NullInt64{Int64: int64(*val), Valid: true}, nil
	case *int64:
		return sql.NullInt64{Int64: int64(*val), Valid: true}, nil
	case *uint:
		return sql.NullInt64{Int64: int64(*val), Valid: true}, nil
	case *uint8:
		return sql.NullInt64{Int64: int64(*val), Valid: true}, nil
	case *uint16:
		return sql.NullInt64{Int64: int64(*val), Valid: true}, nil
	case *uint32:
		return sql.NullInt64{Int64: int64(*val), Valid: true}, nil
	case *uint64:
		return sql.NullInt64{Int64: int64(*val), Valid: true}, nil
	case int:
		return sql.NullInt64{Int64: int64(val), Valid: true}, nil
	case int8:
		return sql.NullInt64{Int64: int64(val), Valid: true}, nil
	case int16:
		return sql.NullInt64{Int64: int64(val), Valid: true}, nil
	case int32:
		return sql.NullInt64{Int64: int64(val), Valid: true}, nil
	case int64:
		return sql.NullInt64{Int64: int64(val), Valid: true}, nil
	case uint:
		return sql.NullInt64{Int64: int64(val), Valid: true}, nil
	case uint8:
		return sql.NullInt64{Int64: int64(val), Valid: true}, nil
	case uint16:
		return sql.NullInt64{Int64: int64(val), Valid: true}, nil
	case uint32:
		return sql.NullInt64{Int64: int64(val), Valid: true}, nil
	case uint64:
		return sql.NullInt64{Int64: int64(val), Valid: true}, nil
	case sql.NullByte:
		return sql.NullInt64{Int64: int64(val.Byte), Valid: val.Valid}, nil
	case sql.NullInt16:
		return sql.NullInt64{Int64: int64(val.Int16), Valid: val.Valid}, nil
	case sql.NullInt32:
		return sql.NullInt64{Int64: int64(val.Int32), Valid: val.Valid}, nil
	case sql.NullInt64:
		return val, nil
	default:
		return sql.NullInt64{}, fmt.Errorf("value is not integer")
	}
}
func (row Row) GetFloat(column string) (val float64, err error) {
	nullable, err := row.GetNullFloat(column)
	if err != nil {
		return 0, err
	}
	if !nullable.Valid {
		return 0, fmt.Errorf("value is nil")
	}
	return nullable.Float64, nil
}
func (row Row) GetNullFloat(column string) (val sql.NullFloat64, err error) {
	isNull, err := row.IsNull(column)
	if err != nil {
		return sql.NullFloat64{}, err
	}
	if isNull {
		return sql.NullFloat64{}, nil
	}
	valAny, exists := row[column]
	if !exists {
		return sql.NullFloat64{}, fmt.Errorf("column not found")
	}
	switch val := valAny.(type) {
	case *float32:
		return sql.NullFloat64{Float64: float64(*val), Valid: true}, nil
	case *float64:
		return sql.NullFloat64{Float64: float64(*val), Valid: true}, nil
	case float32:
		return sql.NullFloat64{Float64: float64(val), Valid: true}, nil
	case float64:
		return sql.NullFloat64{Float64: float64(val), Valid: true}, nil
	case sql.NullFloat64:
		return val, nil
	default:
		return sql.NullFloat64{}, fmt.Errorf("value is not float")
	}
}
func (row Row) GetBoolean(column string) (val bool, err error) {
	nullable, err := row.GetNullBoolean(column)
	if err != nil {
		return false, err
	}
	if !nullable.Valid {
		return false, fmt.Errorf("value is nil")
	}
	return nullable.Bool, nil
}
func (row Row) GetNullBoolean(column string) (val sql.NullBool, err error) {
	isNull, err := row.IsNull(column)
	if err != nil {
		return sql.NullBool{}, err
	}
	if isNull {
		return sql.NullBool{}, nil
	}
	valAny, exists := row[column]
	if !exists {
		return sql.NullBool{}, fmt.Errorf("column not found")
	}
	switch val := valAny.(type) {
	case nil:
		return sql.NullBool{}, nil
	case *bool:
		return sql.NullBool{Bool: *val, Valid: true}, nil
	case bool:
		return sql.NullBool{Bool: val, Valid: true}, nil
	case sql.NullBool:
		return val, nil
	default:
		return sql.NullBool{}, fmt.Errorf("value is not boolean")
	}
}
func (row Row) GetByteArray(column string) (val []byte, err error) {
	isNull, err := row.IsNull(column)
	if err != nil {
		return nil, err
	}
	if isNull {
		return nil, nil
	}
	valAny, exists := row[column]
	if !exists {
		return nil, fmt.Errorf("column not found")
	}
	switch val := valAny.(type) {
	case []byte:
		return val, nil
	default:
		return nil, fmt.Errorf("value is not array of bytes")
	}
}
