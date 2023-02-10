package jsonio

import (
	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

type Row map[string]*JsonValue

func (row Row) GetColumnNames() []string {
	columnNames := []string{}
	for columnName := range row {
		columnNames = append(columnNames, columnName)
	}

	return columnNames
}

func (row Row) Has(columnName string) bool {
	_, exists := row[columnName]

	return exists
}

func (row Row) GetJsonType(columnName string) (JsonType, bool) {
	val, exists := row[columnName]
	if !exists {
		return JsonTypeUnknown, false
	}

	return val.Type, true
}

func (row Row) ToString(columnName string) (string, error) {
	val, ok := row[columnName]
	if !ok {
		return "", errors.Wrap(
			errors.BadKeyAccess,
			"column %s not found in JsonRow", columnName)
	}

	v, err := val.ToString()
	if err != nil {
		return "", errors.Wrap(
			errors.BadConversion,
			"fail to convert value %v:%T of column %s to string", val, val, columnName)
	}

	return v, nil
}

func (row Row) ToBool(columnName string) (bool, error) {
	val, ok := row[columnName]
	if !ok {
		return false, errors.Wrap(
			errors.BadKeyAccess,
			"column %s not found in JsonRow", columnName)
	}

	v, err := val.ToBool()
	if err != nil {
		return false, errors.Wrap(
			errors.BadConversion,
			"fail to convert value %v:%T of column %s to bool", val, val, columnName)
	}

	return v, nil
}

func (row Row) ToInt64(columnName string) (int64, error) {
	val, ok := row[columnName]
	if !ok {
		return 0, errors.Wrap(
			errors.BadKeyAccess,
			"column %s not found in JsonRow", columnName)
	}

	v, err := val.ToInt64()
	if err != nil {
		return 0, errors.Wrap(
			errors.BadConversion,
			"fail to convert value %v:%T of column %s to int64", val, val, columnName)
	}

	return v, nil
}

func (row Row) ToFloat64(columnName string) (float64, error) {
	val, ok := row[columnName]
	if !ok {
		return 0, errors.Wrap(
			errors.BadKeyAccess,
			"column %s not found in JsonRow", columnName)
	}

	v, err := val.ToFloat64()
	if err != nil {
		return 0, errors.Wrap(
			errors.BadConversion,
			"fail to convert value %v:%T of column %s to float64", val, val, columnName)
	}

	return v, nil
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
