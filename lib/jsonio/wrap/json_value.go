package wrap

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type JsonKey string

func JsonKeyInteger(index int) JsonKey {
	return JsonKey(strconv.FormatInt(int64(index), 10))
}

func (k JsonKey) String() string {
	return string(k)
}

func (k JsonKey) Integer() (int, bool) {
	v, err := strconv.ParseInt(string(k), 10, 64)
	return int(v), err != nil
}

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
	fmt.Printf("%#v:%T\n", v, v)
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

var _ json.Marshaler = (*JsonValue)(nil)
var _ json.Unmarshaler = (*JsonValue)(nil)

func (v *JsonValue) MarshalJSON() ([]byte, error) {
	a := ToAny(v)

	b, err := json.Marshal(a)
	if err != nil {
		return nil, errors.Wrap(errors.Join(err, errors.BadConversion), "fail to marshal JSON")
	}

	return b, nil
}

func (v *JsonValue) UnmarshalJSON(b []byte) error {
	var a any
	decoder := json.NewDecoder(bytes.NewBuffer(b))
	decoder.UseNumber()
	if err := decoder.Decode(&a); err != nil {
		return errors.Wrap(errors.Join(err, errors.BadConversion), "fail to unmarshal JSON")
	}

	u, err := FromAny(a)
	if err != nil {
		return errors.Wrap(errors.Join(err, errors.BadConversion), "fail to unmarshal JSON")
	}

	*v = *u

	return nil
}

func (v *JsonValue) MustString() string {
	if v.Type != JsonTypeString {
		panic(errors.Wrap(errors.BadState, "String() is called with json type %v", v.Type))
	}
	return v.StringValue
}

func (v *JsonValue) MustBool() bool {
	if v.Type != JsonTypeBoolean {
		panic(errors.Wrap(errors.BadState, "Bool() is called with json type %v", v.Type))
	}
	return v.BooleanValue
}

func (v *JsonValue) MustNumber() JsonNumber {
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

func (v *JsonValue) MustObject() JsonObject {
	if v.Type != JsonTypeObject {
		panic(errors.Wrap(errors.BadState, "Object() is called with json type %v", v.Type))
	}
	return v.ObjectValue
}

func (v *JsonValue) MustArray() JsonArray {
	if v.Type != JsonTypeArray {
		panic(errors.Wrap(errors.BadState, "Array() is called with json type %v", v.Type))
	}
	return v.ArrayValue
}

func (v *JsonValue) EnumeratePrimitiveKeys() [][]JsonKey {
	switch v.Type {
	case JsonTypeNull, JsonTypeBoolean, JsonTypeNumber, JsonTypeString:
		return [][]JsonKey{{}}
	case JsonTypeObject:
		keys := [][]JsonKey{}
		for k, vk := range v.MustObject() {
			for _, ck := range vk.EnumeratePrimitiveKeys() {
				keys = append(keys, append([]JsonKey{JsonKey(k)}, ck...))
			}
		}

		return keys
	case JsonTypeArray:
		keys := [][]JsonKey{}
		for i, vi := range v.MustArray() {
			for _, ck := range vi.EnumeratePrimitiveKeys() {
				keys = append(keys, append([]JsonKey{JsonKeyInteger(i)}, ck...))
			}
		}

		return keys
	default:
		return nil
	}
}

func (v *JsonValue) Find(keys ...JsonKey) (*JsonValue, bool) {
	if len(keys) == 0 {
		return v, true
	}

	switch v.Type {
	case JsonTypeArray:
		a := v.MustArray()

		index, ok := keys[0].Integer()
		if !ok {
			return nil, false
		}

		if index >= a.Len() {
			return nil, false
		}

		return a.Get(index).Find(keys[1:]...)
	case JsonTypeObject:
		o := v.MustObject()

		key := keys[0].String()

		if !o.Has(key) {
			return nil, false
		}

		return o.Get(key).Find(keys[1:]...)
	default:
		return nil, false
	}
}
