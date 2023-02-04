package jsonio

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

type jsonType string

const (
	JsonTypeUnknown jsonType = "UNKNOWN"
	JsonTypeNull    jsonType = "NULL"
	JsonTypeString  jsonType = "STRING"
	JsonTypeNumber  jsonType = "NUMBER"
	JsonTypeBoolean jsonType = "BOOLEAN"
	JsonTypeArray   jsonType = "ARRAY"
	JsonTypeObject  jsonType = "OBJECT"
)

type (
	JsonNull    any
	JsonString  string
	JsonNumber  json.Number
	JsonBoolean bool
	JsonArray   []*JsonValue
	JsonObject  map[string]*JsonValue
	JsonValue   struct {
		Type         jsonType
		stringValue  JsonString
		numberValue  JsonNumber
		booleanValue JsonBoolean
		arrayValue   JsonArray
		objectValue  JsonObject
	}
)

func NewJsonString(v string) *JsonValue {
	return &JsonValue{Type: JsonTypeString, stringValue: JsonString(v)}
}

func NewJsonBoolean(v bool) *JsonValue {
	return &JsonValue{Type: JsonTypeBoolean, booleanValue: JsonBoolean(v)}
}

func NewJsonNumberInt64(v int64) *JsonValue {
	return &JsonValue{Type: JsonTypeNumber, numberValue: JsonNumber(strconv.FormatInt(v, 10))}
}

func NewJsonNumberFloat64(v float64) *JsonValue {
	return &JsonValue{Type: JsonTypeNumber, numberValue: JsonNumber(strconv.FormatFloat(v, 'f', 15, 64))}
}

func NewJsonNull() *JsonValue {
	return &JsonValue{Type: JsonTypeNull}
}

func NewJsonArrayEmpty() *JsonValue {
	return &JsonValue{Type: JsonTypeArray, arrayValue: JsonArray{}}
}

func NewJsonObjectEmpty() *JsonValue {
	return &JsonValue{Type: JsonTypeObject, objectValue: JsonObject{}}
}

func NewJson(valAny any) (jv *JsonValue, err error) {
	switch val := valAny.(type) {
	case nil:
		return NewJsonNull(), nil
	case string:
		return NewJsonString(val), nil
	case int64:
		return NewJsonNumberInt64(val), nil
	case float64:
		return NewJsonNumberFloat64(val), nil
	case json.Number:
		return &JsonValue{Type: JsonTypeNumber, numberValue: JsonNumber(val)}, nil
	case bool:
		return NewJsonBoolean(val), nil
	}

	rv := reflect.ValueOf(valAny)

	switch rv.Kind() {
	case reflect.Slice:
		arrayValue := JsonArray{}

		for i := 0; i < rv.Len(); i++ {
			rvi := rv.Index(i)

			var vi *JsonValue

			if !rvi.IsValid() {
				vi, err = NewJson(nil)
				if err != nil {
					return nil, errors.Wrap(errors.Join(err, errors.BadArgs), "cannot parse object %v:%T as JsonValue", valAny, valAny)
				}
			} else {
				vi, err = NewJson(rvi.Interface())
				if err != nil {
					return nil, errors.Wrap(errors.Join(err, errors.BadArgs), "cannot parse object %v:%T as JsonValue", valAny, valAny)
				}
			}

			arrayValue = append(arrayValue, vi)
		}

		return &JsonValue{Type: JsonTypeArray, arrayValue: arrayValue}, nil
	case reflect.Map:
		objectValue := JsonObject{}

		for _, rvKey := range rv.MapKeys() {
			rvVal := rv.MapIndex(rvKey)

			key, ok := rvKey.Interface().(string)
			if !ok {
				return nil, errors.Wrap(errors.Join(err, errors.BadArgs), "cannot parse object %v:%T as JsonValue", valAny, valAny)
			}

			var val *JsonValue

			if !rvVal.IsValid() {
				val, err = NewJson(nil)
				if err != nil {
					return nil, errors.Wrap(errors.Join(err, errors.BadArgs), "cannot parse object %v:%T as JsonValue", valAny, valAny)
				}
			} else {
				val, err = NewJson(rvVal.Interface())
				if err != nil {
					return nil, errors.Wrap(errors.Join(err, errors.BadArgs), "cannot parse object %v:%T as JsonValue", valAny, valAny)
				}
			}

			objectValue[key] = val
		}

		return &JsonValue{Type: JsonTypeObject, objectValue: objectValue}, nil
	}

	return nil, errors.Wrap(errors.Join(err, errors.BadArgs), "cannot parse object %v:%T as JsonValue", valAny, valAny)
}

func (v *JsonValue) ToString() (string, error) {
	switch v.Type {
	case JsonTypeString:
		return string(v.stringValue), nil
	case JsonTypeBoolean:
		return fmt.Sprintf("%t", v.booleanValue), nil
	case JsonTypeNull:
		return "null", nil
	case JsonTypeNumber:
		return json.Number(v.numberValue).String(), nil
	default:
		return "", errors.Wrap(errors.BadConversion, "cannot convert %v to string", v.Type)
	}
}

func (v *JsonValue) ToBool() (bool, error) {
	switch v.Type {
	case JsonTypeString:
		switch strings.ToLower(string(v.stringValue)) {
		case "true":
			return true, nil
		case "false":
			return false, nil
		}

		return false, errors.Wrap(errors.BadConversion, "cannot convert string value %v to bool", v.stringValue)
	case JsonTypeBoolean:
		return bool(v.booleanValue), nil
	case JsonTypeNull:
		return false, nil
	case JsonTypeNumber:
		if v, e := json.Number(v.numberValue).Int64(); e == nil {
			return v != 0, nil
		}

		if v, e := json.Number(v.numberValue).Float64(); e == nil {
			return v != 0.0, nil
		}

		return false, errors.Wrap(errors.BadConversion, "cannot convert number value %v to bool", v.numberValue)
	default:
		return false, errors.Wrap(errors.BadConversion, "cannot convert %v to bool", v.Type)
	}
}

func (v *JsonValue) ToInt64() (int64, error) {
	switch v.Type {
	case JsonTypeNumber:
		text := string(v.numberValue)
		// integer that may has trailing zeros after decimal point: regexr.com/776pj
		if regexp.MustCompile(`^-?(0|([1-9][0-9]*))(\.0+)?$`).MatchString(text) {
			text = regexp.MustCompile(`(\.0+)?$`).ReplaceAllString(text, "")

			i, err := json.Number(text).Int64()
			if err != nil {
				return 0, errors.Wrap(errors.BadConversion, "cannot convert number value %v to int64", v.numberValue)
			}

			return i, nil
		}

		return 0, errors.Wrap(errors.BadConversion, "cannot convert number value %v to int64", v.numberValue)
	case JsonTypeBoolean:
		if v.booleanValue {
			return 1, nil
		}

		return 0, nil
	case JsonTypeNull:
		return 0, nil
	case JsonTypeString:
		i, err := strconv.ParseInt(string(v.stringValue), 10, 64)
		if err != nil {
			return 0, errors.Wrap(errors.BadConversion, "cannot convert string value %v to int64", v.stringValue)
		}

		return i, nil
	default:
		return 0, errors.Wrap(errors.BadConversion, "cannot convert %v to int64", v.Type)
	}
}

func (v *JsonValue) ToFloat64() (float64, error) {
	switch v.Type {
	case JsonTypeNumber:
		f, err := json.Number(v.numberValue).Float64()
		if err != nil {
			return 0, errors.Wrap(errors.BadConversion, "cannot convert number value %v to float64", v.numberValue)
		}

		return f, nil
	case JsonTypeBoolean:
		if v.booleanValue {
			return 1, nil
		}

		return 0, nil
	case JsonTypeNull:
		return 0, nil
	case JsonTypeString:
		f, err := strconv.ParseFloat(string(v.stringValue), 64)
		if err != nil {
			return 0, errors.Wrap(errors.BadConversion, "cannot convert string value %v to float64", v.stringValue)
		}

		return f, nil
	default:
		return 0, errors.Wrap(errors.BadConversion, "cannot convert %v to float64", v.Type)
	}
}

func (v *JsonValue) AsObject() (JsonObject, error) {
	if v.Type != JsonTypeObject {
		return nil, errors.Wrap(errors.BadConversion, "cannot convert %v to JsonTypeObject", v.Type)
	}

	return v.objectValue, nil
}

func (v *JsonValue) AsArray() (JsonArray, error) {
	if v.Type != JsonTypeArray {
		return nil, errors.Wrap(errors.BadConversion, "cannot convert %v to JsonTypeArray", v.Type)
	}

	return v.arrayValue, nil
}

func (o JsonObject) Keys() []string {
	keys := []string{}
	for k := range o {
		keys = append(keys, k)
	}

	return keys
}

func (o JsonObject) Get(key string) (*JsonValue, bool) {
	val, ok := o[key]
	if !ok {
		return nil, false
	}

	return val, true
}

func (o JsonObject) Set(key string, val *JsonValue) JsonObject {
	if val == nil {
		val = NewJsonNull()
	}

	o[key] = val

	return o
}

func (o JsonObject) AsJsonValue() *JsonValue {
	return &JsonValue{Type: JsonTypeObject, objectValue: o}
}

func (a JsonArray) Len() int {
	return len(a)
}

func (a JsonArray) Get(i int) (*JsonValue, bool) {
	if i >= len(a) {
		return nil, false
	}

	return a[i], true
}

func (a JsonArray) Set(i int, val *JsonValue) bool {
	if i >= len(a) {
		return false
	}

	if val == nil {
		val = NewJsonNull()
	}

	a[i] = val

	return true
}

func (a JsonArray) Append(val *JsonValue) JsonArray {
	if val == nil {
		val = NewJsonNull()
	}

	a = append(a, val)

	return a
}

func (a JsonArray) AsJsonValue() *JsonValue {
	return &JsonValue{Type: JsonTypeArray, arrayValue: a}
}
