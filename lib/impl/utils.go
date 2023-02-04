package impl

import (
	"fmt"
	"time"

	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/io"
)

func ExtractColumnValueAsDB(row io.Row, columnName string, dbType db.ColumnType) (any, error) {
	isNull := false
	if !row.Has(columnName) {
		isNull = true
	} else {
		jsonType, err := row.GetJsonType(columnName)
		if err != nil {
			return nil, err
		}
		if jsonType == io.JsonTypeNull {
			isNull = true
		}
	}
	switch dbType {
	case db.ColumnTypeBoolean:
		if isNull {
			return (*bool)(nil), nil
		}
		val, err := row.ToBool(columnName)
		if err != nil {
			return nil, err
		}
		return val, nil
	case db.ColumnTypeInteger:
		if isNull {
			return (*int64)(nil), nil
		}
		val, err := row.ToInt64(columnName)
		if err != nil {
			return nil, err
		}
		return val, nil
	case db.ColumnTypeFloat:
		if isNull {
			return (*float64)(nil), nil
		}
		val, err := row.ToFloat64(columnName)
		if err != nil {
			return nil, err
		}
		return val, nil
	case db.ColumnTypeString:
		if isNull {
			return (*string)(nil), nil
		}
		val, err := row.ToString(columnName)
		if err != nil {
			return nil, err
		}
		return val, nil
	case db.ColumnTypeTime:
		if isNull {
			return (*time.Time)(nil), nil
		}
		val, err := row.ToString(columnName)
		if err != nil {
			return nil, err
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
