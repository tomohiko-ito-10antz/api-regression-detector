package impl

import (
	"time"

	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
)

func ExtractColumnValueAsDB(row jsonio.Row, columnName string, dbType db.ColumnType) (any, error) {
	jsonType, ok := row.GetJsonType(columnName)
	isNull := !ok || jsonType == jsonio.JsonTypeNull

	switch dbType {
	case db.ColumnTypeBoolean:
		if isNull {
			return (*bool)(nil), nil
		}

		val, err := row.ToBool(columnName)
		if err != nil {
			return nil, errors.Wrap(
				errors.Join(err, errors.BadConversion),
				"fail to convert value of column %s to bool (json type=%v,db type=%v)", columnName, jsonType, dbType)
		}

		return val, nil
	case db.ColumnTypeInteger:
		if isNull {
			return (*int64)(nil), nil
		}

		val, err := row.ToInt64(columnName)
		if err != nil {
			return nil, errors.Wrap(
				errors.Join(err, errors.BadConversion),
				"fail to convert value of column %s to int64 (json type=%v,db type=%v)", columnName, jsonType, dbType)
		}

		return val, nil
	case db.ColumnTypeFloat:
		if isNull {
			return (*float64)(nil), nil
		}

		val, err := row.ToFloat64(columnName)
		if err != nil {
			return nil, errors.Wrap(
				errors.Join(err, errors.BadConversion),
				"fail to convert value of column %s to float64 (json type=%v,db type=%v)", columnName, jsonType, dbType)
		}

		return val, nil
	case db.ColumnTypeString:
		if isNull {
			return (*string)(nil), nil
		}

		val, err := row.ToString(columnName)
		if err != nil {
			return nil, errors.Wrap(
				errors.Join(err, errors.BadConversion),
				"fail to convert value of column %s to string (json type=%v,db type=%v)", columnName, jsonType, dbType)
		}

		return val, nil
	case db.ColumnTypeTime:
		if isNull {
			return (*time.Time)(nil), nil
		}

		val, err := row.ToString(columnName)
		if err != nil {
			return nil, errors.Wrap(
				errors.Join(err, errors.BadConversion),
				"fail to convert value of column %s to time.Time (json type=%v,db type=%v)", columnName, jsonType, dbType)
		}

		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, errors.Wrap(
				errors.Join(err, errors.BadConversion),
				"fail to convert value of column %s to time.Time (json type=%v,db type=%v)", columnName, jsonType, dbType)
		}

		return t, nil
	default:
		return nil, errors.Wrap(errors.Join(errors.BadArgs), "unexpected DB column type %v of column %s", dbType, columnName)
	}
}
