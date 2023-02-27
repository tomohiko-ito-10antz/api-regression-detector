package tables

import (
	"regexp"
	"sort"
	"strconv"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
)

type Row map[string]*wrap.JsonValue

func (row Row) GetColumnNames() []string {
	columnNames := []string{}
	for columnName := range row {
		columnNames = append(columnNames, columnName)
	}

	sort.Slice(columnNames, func(i, j int) bool {
		return columnNames[i] < columnNames[j]
	})

	return columnNames
}

func (row Row) GetJsonType(columnName string) (wrap.JsonType, bool) {
	val, ok := row[columnName]
	if !ok {
		return wrap.JsonTypeNull, false
	}

	return val.Type, true
}

func (row Row) Has(columnName string) bool {
	_, exists := row[columnName]

	return exists
}

func (row Row) ToString(columnName string) (string, error) {
	errInfo := errors.Info{"row": row, "columnName": columnName}

	val, ok := row[columnName]
	if !ok {
		return "", errors.BadKeyAccess.New(
			errInfo.AppendTo("column not found in JsonRow"))
	}

	errInfo = errInfo.With("val", val)

	switch val.Type {
	case wrap.JsonTypeNull:
		return "", nil
	case wrap.JsonTypeString:
		return val.MustString(), nil
	case wrap.JsonTypeNumber:
		return string(val.MustNumber()), nil
	case wrap.JsonTypeBoolean:
		return strconv.FormatBool(val.MustBool()), nil
	default:
		return "", errors.BadConversion.New(
			errInfo.AppendTo("fail to convert value from JSON to string"))
	}
}

func parseAsInteger(text string) (int64, bool) {
	// text can be parsed of the text is a number that may has trailing zeros after decimal point: regexr.com/776pj
	if !regexp.MustCompile(`^-?(0|([1-9][0-9]*))(\.0+)?$`).MatchString(text) {
		return 0, false
	}

	text = regexp.MustCompile(`(\.0+)?$`).ReplaceAllString(text, "")
	v, err := strconv.ParseInt(text, 10, 64)

	return v, err == nil
}

func (row Row) ToBool(columnName string) (bool, error) {
	errInfo := errors.Info{"row": row, "columnName": columnName}

	val, ok := row[columnName]
	if !ok {
		return false, errors.BadKeyAccess.New(
			errInfo.AppendTo("column not found in JsonRow"))
	}

	errInfo = errInfo.With("val", val)

	switch val.Type {
	case wrap.JsonTypeNull:
		return false, nil
	case wrap.JsonTypeString:
		return val.MustString() != "", nil
	case wrap.JsonTypeNumber:
		if v, err := val.MustNumber().Int64(); err == nil {
			return v != 0, nil
		}

		if v, err := val.MustNumber().Float64(); err == nil {
			return v != 0, nil
		}

		return false, errors.BadConversion.New(
			errInfo.AppendTo("fail to convert value from JSON to bool"))
	case wrap.JsonTypeBoolean:
		return val.MustBool(), nil
	default:
		return false, errors.BadConversion.New(
			errInfo.AppendTo("fail to convert value from JSON to bool"))
	}
}

func (row Row) ToInt64(columnName string) (int64, error) {
	errInfo := errors.Info{"row": row, "columnName": columnName}

	val, ok := row[columnName]
	if !ok {
		return 0, errors.BadKeyAccess.New(
			errInfo.AppendTo("column not found in JsonRow"))
	}

	errInfo = errInfo.With("val", val)

	switch val.Type {
	case wrap.JsonTypeNull:
		return 0, nil
	case wrap.JsonTypeString:
		v, ok := parseAsInteger(val.MustString())
		if !ok {
			return 0, errors.BadConversion.New(
				errInfo.AppendTo("fail to convert value from JSON to int64"))
		}

		return v, nil
	case wrap.JsonTypeNumber:
		v, err := val.MustNumber().Int64()
		if err != nil {
			return 0, errors.BadConversion.New(
				errInfo.AppendTo("fail to convert value from JSON to int64"))
		}

		return v, nil
	case wrap.JsonTypeBoolean:
		if val.MustBool() {
			return 1, nil
		} else {
			return 0, nil
		}
	default:
		return 0, errors.BadConversion.New(
			errInfo.AppendTo("fail to convert value from JSON to int64"))
	}
}

func (row Row) ToFloat64(columnName string) (float64, error) {
	errInfo := errors.Info{"row": row, "columnName": columnName}

	val, ok := row[columnName]
	if !ok {
		return 0, errors.BadKeyAccess.New(
			errInfo.AppendTo("column not found in JsonRow"))
	}

	errInfo = errInfo.With("val", val)

	switch val.Type {
	case wrap.JsonTypeNull:
		return 0, nil
	case wrap.JsonTypeString:
		v, err := strconv.ParseFloat(val.MustString(), 64)
		if err != nil {
			return 0, errors.BadConversion.New(
				errInfo.AppendTo("fail to convert value from JSON to float64"))
		}

		return v, nil
	case wrap.JsonTypeNumber:
		v, err := val.MustNumber().Float64()
		if err != nil {
			return 0, errors.BadConversion.New(
				errInfo.AppendTo("fail to convert value from JSON to float64"))
		}

		return v, nil
	case wrap.JsonTypeBoolean:
		if val.MustBool() {
			return 1, nil
		} else {
			return 0, nil
		}
	default:
		return 0, errors.BadConversion.New(
			errInfo.AppendTo("fail to convert value from JSON to float64"))
	}
}

func (row Row) SetString(columnName string, val string) {
	row[columnName] = wrap.String(val)
}

func (row Row) SetBool(columnName string, val bool) {
	row[columnName] = wrap.Boolean(val)
}

func (row Row) SetInt64(columnName string, val int64) {
	row[columnName] = wrap.Number(val)
}

func (row Row) SetFloat64(columnName string, val float64) {
	row[columnName] = wrap.Number(val)
}

func (row Row) SetNil(columnName string) {
	row[columnName] = wrap.Null()
}
