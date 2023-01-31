package impl

import (
	"fmt"
	"time"

	"github.com/Jumpaku/api-regression-detector/db"
	"github.com/Jumpaku/api-regression-detector/io"
)

func ExtractColumnValueAsDB(row io.Row, columnName string, dbType db.ColumnType) (any, error) {
	isNull := false
	jsonType, err := row.GetJsonType(columnName)
	if err != nil {
		return nil, err
	}
	if jsonType == io.JsonTypeNull {
		isNull = true
	}
	switch dbType {
	case db.ColumnTypeBoolean:
		val, err := row.ToBool(columnName)
		if err != nil {
			return nil, err
		}
		if isNull {
			return (*bool)(nil), nil
		}
		return val, nil
	case db.ColumnTypeInteger:
		val, err := row.ToInt64(columnName)
		if err != nil {
			return nil, err
		}
		if isNull {
			return (*int64)(nil), nil
		}
		return val, nil
	case db.ColumnTypeFloat:
		val, err := row.ToFloat64(columnName)
		if err != nil {
			return nil, err
		}
		if isNull {
			return (*float64)(nil), nil
		}
		return val, nil
	case db.ColumnTypeString:
		val, err := row.ToString(columnName)
		if err != nil {
			return nil, err
		}
		if isNull {
			return (*string)(nil), nil
		}
		return val, nil
	case db.ColumnTypeTime:
		val, err := row.ToString(columnName)
		if err != nil {
			return nil, err
		}
		if isNull {
			return (*time.Time)(nil), nil
		}
		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, err
		}
		return t, nil
	default:
		return nil, fmt.Errorf("unexpected database column type %s", dbType)
	}
}
