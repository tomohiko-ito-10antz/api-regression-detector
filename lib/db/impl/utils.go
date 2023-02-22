package impl

import (
	"time"

	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/tables"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
)

func ExtractColumnValueAsDB(row tables.Row, columnName string, dbType db.ColumnType) (any, error) {
	jsonType, ok := row.GetJsonType(columnName)
	isNull := !ok || jsonType == wrap.JsonTypeNull
	errInfo := errors.Info{"columnName": columnName, "jsonType": jsonType, "dbType": dbType}

	switch dbType {
	case db.ColumnTypeBoolean:
		if isNull {
			return (*bool)(nil), nil
		}

		val, err := row.ToBool(columnName)
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.AppendTo("fail to convert column value from JSON to boolean"))
		}

		return val, nil
	case db.ColumnTypeInteger:
		if isNull {
			return (*int64)(nil), nil
		}

		val, err := row.ToInt64(columnName)
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.AppendTo("fail to convert column value from JSON to integer"))
		}

		return val, nil
	case db.ColumnTypeFloat:
		if isNull {
			return (*float64)(nil), nil
		}

		val, err := row.ToFloat64(columnName)
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.AppendTo("fail to convert column value from JSON to float"))
		}

		return val, nil
	case db.ColumnTypeString:
		if isNull {
			return (*string)(nil), nil
		}

		val, err := row.ToString(columnName)
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.AppendTo("fail to convert column value from JSON to string"))
		}

		return val, nil
	case db.ColumnTypeTime:
		if isNull {
			return (*time.Time)(nil), nil
		}

		val, err := row.ToString(columnName)
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.AppendTo("fail to convert column value from JSON to string"))
		}

		t, err := time.Parse(time.RFC3339, val)
		if err != nil {
			return nil, errors.Wrap(
				errors.BadConversion.Err(err),
				errInfo.AppendTo("fail to convert column value from JSON to time.Time"))
		}

		return t, nil
	default:
		return nil, errors.Unsupported.New(errInfo.AppendTo("unsupported DB column type"))
	}
}
