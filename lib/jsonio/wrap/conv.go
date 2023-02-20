package wrap

import (
	"bytes"
	"encoding/json"
	"reflect"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

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
