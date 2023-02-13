package wrap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

type JsonType string

const (
	JsonTypeUnknown JsonType = "JsonTypeUnknown"
	JsonTypeNull    JsonType = "JsonTypeNull"
	JsonTypeString  JsonType = "JsonTypeString"
	JsonTypeNumber  JsonType = "JsonTypeNumber"
	JsonTypeBoolean JsonType = "JsonTypeBoolean"
	JsonTypeArray   JsonType = "JsonTypeArray"
	JsonTypeObject  JsonType = "JsonTypeObject"
)

type JsonValue struct {
	Type         JsonType
	StringValue  string
	NumberValue  JsonNumber
	BooleanValue bool
	ArrayValue   JsonArray
	ObjectValue  JsonObject
}

func String(v string) *JsonValue {
	return &JsonValue{Type: JsonTypeString, StringValue: v}
}

func mustToInt64(v any) int64 {
	switch v := v.(type) {
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return int64(v)
	case uint:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	}
	panic(fmt.Sprintf("must be compatible to int64 but %v", v))
}

func mustToFloat64(v any) float64 {
	switch v := v.(type) {
	case float32:
		return float64(v)
	case float64:
		return float64(v)
	}
	panic(fmt.Sprintf("must be compatible to float64 but %v", v))
}

func Number[T json.Number | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](v T) *JsonValue {
	switch v := any(v).(type) {
	case json.Number:
		return &JsonValue{Type: JsonTypeNumber, NumberValue: JsonNumber(v)}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return Number(json.Number(strconv.FormatInt(mustToInt64(v), 10)))
	case float32, float64:
		return Number(json.Number(strconv.FormatFloat(mustToFloat64(v), 'f', 16, 64)))
	default:
		panic("never match this case")
	}
}

func Boolean(v bool) *JsonValue {
	return &JsonValue{Type: JsonTypeBoolean, BooleanValue: v}
}

var nullValue = &JsonValue{Type: JsonTypeNull}

func Null() *JsonValue {
	return nullValue
}

func Object(mv map[string]*JsonValue) *JsonValue {
	for k, v := range mv {
		if v == nil {
			mv[k] = Null()
		}
	}
	return &JsonValue{Type: JsonTypeObject, ObjectValue: JsonObject(mv)}
}

func Array(av ...*JsonValue) *JsonValue {
	for i, v := range av {
		if v == nil {
			av[i] = Null()
		}
	}
	return &JsonValue{Type: JsonTypeArray, ArrayValue: JsonArray(av)}
}

func FromAny(valAny any) (jv *JsonValue, err error) {
	switch val := valAny.(type) {
	case nil:
		return Null(), nil
	case string:
		return String(val), nil
	case json.Number:
		return Number(val), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return Number(mustToInt64(val)), nil
	case float32, float64:
		return Number(mustToFloat64(val)), nil
	case bool:
		return Boolean(val), nil
	case JsonNumber:
		return &JsonValue{Type: JsonTypeNumber, NumberValue: val}, nil
	case JsonObject:
		return &JsonValue{Type: JsonTypeObject, ObjectValue: val}, nil
	case JsonArray:
		return &JsonValue{Type: JsonTypeArray, ArrayValue: val}, nil
	case JsonValue:
		return &val, nil
	}

	rv := reflect.ValueOf(valAny)
	if !rv.IsValid() {
		return Null(), nil
	}

	switch rv.Kind() {
	case reflect.Pointer, reflect.Interface:
		if rv.IsNil() {
			return Null(), nil
		}

		v, err := FromAny(rv.Elem().Interface())
		if err != nil {
			return nil, errors.Wrap(errors.Join(err, errors.BadArgs), "cannot parse %v:%T as JsonValue", valAny, valAny)
		}

		return v, nil
	case reflect.Slice, reflect.Array:
		arrayValue := JsonArray{}
		for i := 0; i < rv.Len(); i++ {
			vi := Null()

			if rvi := rv.Index(i); rvi.IsValid() {
				vi, err = FromAny(rvi.Interface())
				if err != nil {
					return nil, errors.Wrap(errors.Join(err, errors.BadArgs), "cannot parse %v:%T as JsonValue", valAny, valAny)
				}
			}

			arrayValue = append(arrayValue, vi)
		}

		return &JsonValue{Type: JsonTypeArray, ArrayValue: arrayValue}, nil
	case reflect.Map:
		objectValue := JsonObject{}
		for _, rvKey := range rv.MapKeys() {
			rvVal := rv.MapIndex(rvKey)

			rvKeyStr, ok := rvKey.Interface().(string)
			if !ok {
				return nil, errors.Wrap(errors.Join(err, errors.BadArgs), "cannot parse %v:%T as JsonValue", valAny, valAny)
			}

			rvValJson := Null()

			if rvVal.IsValid() {
				rvValJson, err = FromAny(rvVal.Interface())
				if err != nil {
					return nil, errors.Wrap(errors.Join(err, errors.BadArgs), "cannot parse %v:%T as JsonValue", valAny, valAny)
				}
			}

			objectValue[rvKeyStr] = rvValJson
		}

		return &JsonValue{Type: JsonTypeObject, ObjectValue: objectValue}, nil
	default:
		b := bytes.NewBuffer(nil)
		e := json.NewEncoder(b)
		d := json.NewDecoder(b)

		if err := e.Encode(valAny); err != nil {
			return nil, errors.Wrap(errors.Join(err, errors.BadArgs), "cannot parse %v:%T as JsonValue", valAny, valAny)
		}

		var a any
		if err := d.Decode(&a); err != nil {
			return nil, errors.Wrap(errors.Join(err, errors.BadArgs), "cannot parse %v:%T as JsonValue", valAny, valAny)
		}

		v, err := FromAny(a)
		if err != nil {
			return nil, errors.Wrap(errors.Join(err, errors.BadArgs), "cannot parse %v:%T as JsonValue", valAny, valAny)
		}

		return v, nil
	}
}

func (v *JsonValue) AsString() string {
	if v.Type != JsonTypeString {
		panic(errors.Wrap(errors.BadState, "String() is called with json type %v", v.Type))
	}
	return v.StringValue
}

func (v *JsonValue) AsBool() bool {
	if v.Type != JsonTypeBoolean {
		panic(errors.Wrap(errors.BadState, "Bool() is called with json type %v", v.Type))
	}
	return v.BooleanValue
}

func (v *JsonValue) AsNumber() JsonNumber {
	if v.Type != JsonTypeNumber {
		panic(errors.Wrap(errors.BadState, "Number() is called with json type %v", v.Type))
	}
	return v.NumberValue
}

func (v *JsonValue) Int64() (int64, bool) {
	if v.Type != JsonTypeNumber {
		panic(errors.Wrap(errors.BadState, "Int64() is called with json type %v", v.Type))
	}
	return v.NumberValue.Int64()
}

func (v *JsonValue) Float64() (float64, bool) {
	if v.Type != JsonTypeNumber {
		panic(errors.Wrap(errors.BadState, "Float64() is called with json type %v", v.Type))
	}
	return v.NumberValue.Float64()
}

func (v *JsonValue) AsObject() JsonObject {
	if v.Type != JsonTypeObject {
		panic(errors.Wrap(errors.BadState, "Object() is called with json type %v", v.Type))
	}
	return v.ObjectValue
}

func (v *JsonValue) AsArray() JsonArray {
	if v.Type != JsonTypeArray {
		panic(errors.Wrap(errors.BadState, "Array() is called with json type %v", v.Type))
	}
	return v.ArrayValue
}

func (v *JsonValue) FindByKeySeq(keys ...any) (*JsonValue, bool) {
	if len(keys) == 0 {
		return v, true
	}

	switch v.Type {
	case JsonTypeArray:
		a := v.AsArray()

		var index int
		switch i := keys[0].(type) {
		case int:
			index = i
		case string:
			i64, err := strconv.ParseInt(i, 10, 64)
			if err != nil {
				return nil, false
			}

			index = int(i64)
		default:
			return nil, false
		}

		if index >= a.Len() {
			return nil, false
		}

		return a.Get(index).FindByKeySeq(keys[1:])
	case JsonTypeObject:
		o := v.AsObject()

		key, ok := keys[0].(string)
		if ok {
			return nil, false
		}

		if !o.Has(key) {
			return nil, false
		}

		return o.Get(key).FindByKeySeq(keys[1:])
	default:
		return nil, false
	}
}

func ToAny(v *JsonValue) any {
	if v == nil {
		return nil
	}

	switch v.Type {
	case JsonTypeNull:
		return nil
	case JsonTypeBoolean:
		return v.BooleanValue
	case JsonTypeNumber:
		return json.Number(v.NumberValue)
	case JsonTypeString:
		return v.StringValue
	case JsonTypeArray:
		arr := []any{}
		for _, vi := range v.ArrayValue {
			arr = append(arr, ToAny(vi))
		}

		return arr
	case JsonTypeObject:
		obj := map[string]any{}
		for key, val := range v.ObjectValue {
			obj[key] = ToAny(val)
		}

		return obj
	default:
		panic(errors.Wrap(errors.BadState, "unexpected case of ToAny(): %v", v))
	}
}
