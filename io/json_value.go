package io

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

type JsonNull any
type JsonString string
type JsonNumber json.Number
type JsonBoolean bool
type JsonArray []*JsonValue
type JsonObject map[string]*JsonValue
type JsonValue struct {
	Type         jsonType
	stringValue  JsonString
	numberValue  JsonNumber
	booleanValue JsonBoolean
	arrayValue   JsonArray
	objectValue  JsonObject
}

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
					return nil, err
				}
			} else {
				vi, err = NewJson(rvi.Interface())
				if err != nil {
					return nil, err
				}
			}
			arrayValue = append(arrayValue, vi)
		}
		return &JsonValue{Type: JsonTypeArray, arrayValue: arrayValue}, nil
	case reflect.Map:
		objectValue := JsonObject{}
		for _, rvKey := range rv.MapKeys() {
			rvVal := rv.MapIndex(rvKey)
			key := rvKey.Interface().(string)
			var val *JsonValue
			if !rvVal.IsValid() {
				val, err = NewJson(nil)
				if err != nil {
					return nil, err
				}
			} else {
				val, err = NewJson(rvVal.Interface())
				if err != nil {
					return nil, err
				}
			}
			objectValue[key] = val
		}
		return &JsonValue{Type: JsonTypeObject, objectValue: objectValue}, nil
	}
	return nil, fmt.Errorf("cannot new JsonValue %v:%T", valAny, valAny)
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
		return "", fmt.Errorf("cannot convert value of %v to string", v.Type)
	}
}

func (v *JsonValue) ToBool() (bool, error) {
	switch v.Type {
	case JsonTypeString:
		switch strings.ToLower(string(v.stringValue)) {
		case "1", "true":
			return true, nil
		case "", "0", "null", "false":
			return false, nil
		}
		return false, fmt.Errorf("cannot convert string value %v to bool", v.stringValue)
	case JsonTypeBoolean:
		return bool(v.booleanValue), nil
	case JsonTypeNull:
		return false, nil
	case JsonTypeNumber:
		switch json.Number(v.numberValue).String() {
		case "1":
			return true, nil
		case "0":
			return false, nil
		}
		return false, fmt.Errorf("cannot convert number value %v to bool", json.Number(v.numberValue).String())
	case JsonTypeArray:
		size := len([]*JsonValue(v.arrayValue))
		if size == 0 {
			return false, nil
		}
		return false, fmt.Errorf("cannot convert array value of length %d to bool", size)
	default:
		return false, fmt.Errorf("cannot convert value of %v to bool", v.Type)
	}
}

func (v *JsonValue) ToInt64() (int64, error) {
	switch v.Type {
	case JsonTypeNumber:
		i, err := json.Number(v.numberValue).Int64()
		if err != nil {
			return 0, fmt.Errorf("cannot convert number value %v to int64", json.Number(v.numberValue).String())
		}
		return i, nil
	case JsonTypeBoolean:
		if v.booleanValue {
			return 1, nil
		} else {
			return 0, nil
		}
	case JsonTypeNull:
		return 0, nil
	case JsonTypeString:
		i, err := strconv.ParseInt(string(v.stringValue), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string value %v to int64", v.stringValue)
		}
		return i, nil
	case JsonTypeArray:
		size := len([]*JsonValue(v.arrayValue))
		if size == 0 {
			return 0, nil
		}
		return 0, fmt.Errorf("cannot convert array value of length %d to int64", size)
	default:
		return 0, fmt.Errorf("cannot convert value of %v to int64", v.Type)
	}
}

func (v *JsonValue) ToFloat64() (float64, error) {
	switch v.Type {
	case JsonTypeNumber:
		f, err := json.Number(v.numberValue).Float64()
		if err != nil {
			return 0, fmt.Errorf("cannot convert number value %v to float64", json.Number(v.numberValue).String())
		}
		return f, nil
	case JsonTypeBoolean:
		if v.booleanValue {
			return 1, nil
		} else {
			return 0, nil
		}
	case JsonTypeNull:
		return 0, nil
	case JsonTypeString:
		f, err := strconv.ParseFloat(string(v.stringValue), 64)
		if err != nil {
			return 0, fmt.Errorf("cannot convert string value %v to float64", v.stringValue)
		}
		return f, nil
	case JsonTypeArray:
		size := len([]*JsonValue(v.arrayValue))
		if size == 0 {
			return 0, nil
		}
		return 0, fmt.Errorf("cannot convert array value of length %d to float64", size)
	default:
		return 0, fmt.Errorf("cannot convert value of %v to float64", v.Type)
	}
}

func (v *JsonValue) ObjectKeys() (keys []string, err error) {
	if v.Type != JsonTypeObject {
		return nil, fmt.Errorf("ObjectKeys must be called with JsonValue of type JsonTypeObject")
	}
	for k := range v.objectValue {
		keys = append(keys, k)
	}
	return keys, nil
}
func (v *JsonValue) ObjectGet(key string) (*JsonValue, error) {
	if v.Type != JsonTypeObject {
		return nil, fmt.Errorf("ObjectGet must be called with JsonValue of type JsonTypeObject")
	}
	val, ok := v.objectValue[key]
	if !ok {
		return nil, fmt.Errorf("value not found for key %s", key)
	}
	return val, nil
}

func (v *JsonValue) ObjectSet(key string, val *JsonValue) (err error) {
	if v.Type != JsonTypeObject {
		return fmt.Errorf("ObjectSet must be called with JsonValue of type JsonTypeObject")
	}
	if val == nil {
		val, err = NewJson(nil)
		if err != nil {
			return err
		}
	}
	m := map[string]*JsonValue(v.objectValue)
	m[key] = val
	v.objectValue = m
	return nil
}

func (v *JsonValue) ArrayLen() (size int, err error) {
	if v.Type != JsonTypeArray {
		return 0, fmt.Errorf("ArrayLen must be called with JsonValue with type JsonTypeArray")
	}
	return len([]*JsonValue(v.arrayValue)), nil
}

func (v *JsonValue) ArrayGet(i int) (val *JsonValue, err error) {
	if v.Type != JsonTypeArray {
		return nil, fmt.Errorf("ArrayGet must be called with JsonValue with type JsonTypeArray")
	}
	return v.arrayValue[i], nil
}

func (v *JsonValue) ArraySet(i int, val *JsonValue) (err error) {
	if v.Type != JsonTypeArray {
		return fmt.Errorf("ArraySet must be called with JsonValue with type JsonTypeArray")
	}
	if val == nil {
		val, err = NewJson(nil)
		if err != nil {
			return err
		}
	}
	v.arrayValue[i] = val
	return nil
}

func (v *JsonValue) ArrayAppend(val *JsonValue) (err error) {
	if v.Type != JsonTypeArray {
		return fmt.Errorf("ArrayAppend must be called with JsonValue with type JsonTypeArray")
	}
	if val == nil {
		val, err = NewJson(nil)
		if err != nil {
			return err
		}
	}
	a := []*JsonValue(v.arrayValue)
	a = append(a, val)
	v.arrayValue = a
	return nil
}
