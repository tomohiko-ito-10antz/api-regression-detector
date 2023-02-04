package io_json_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Jumpaku/api-regression-detector/lib/io_json"
	"github.com/Jumpaku/api-regression-detector/test/assert"
	"golang.org/x/exp/slices"
)

func TestNewJsonString(t *testing.T) {
	s := "abc"
	v := io_json.NewJsonString(s)
	assert.NotEqual(t, v, nil)
	assert.Equal(t, v.Type, io_json.JsonTypeString)
	a, err := v.ToString()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, s)
}

func TestNewJsonBooleanTrue(t *testing.T) {
	b := true
	v := io_json.NewJsonBoolean(b)
	assert.NotEqual(t, v, nil)
	assert.Equal(t, v.Type, io_json.JsonTypeBoolean)
	a, err := v.ToBool()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, b)
}

func TestNewJsonBooleanFalse(t *testing.T) {
	b := false
	v := io_json.NewJsonBoolean(b)
	assert.NotEqual(t, v, nil)
	assert.Equal(t, v.Type, io_json.JsonTypeBoolean)
	a, err := v.ToBool()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, b)
}

func TestNewJsonNumberInt64(t *testing.T) {
	i := int64(123)
	v := io_json.NewJsonNumberInt64(i)
	assert.NotEqual(t, v, nil)
	assert.Equal(t, v.Type, io_json.JsonTypeNumber)
	a, err := v.ToInt64()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, i)
}

func TestNewJsonNumberFloat64(t *testing.T) {
	f := -123.45
	v := io_json.NewJsonNumberFloat64(f)
	assert.NotEqual(t, v, nil)
	assert.Equal(t, v.Type, io_json.JsonTypeNumber)
	a, err := v.ToFloat64()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, f)
}

func TestNewJsonNull(t *testing.T) {
	v := io_json.NewJsonNull()
	assert.NotEqual(t, v, nil)
	assert.Equal(t, v.Type, io_json.JsonTypeNull)
}

func TestNewJsonObjectEmpty(t *testing.T) {
	v := io_json.NewJsonObjectEmpty()
	assert.NotEqual(t, v, nil)
	assert.Equal(t, v.Type, io_json.JsonTypeObject)
}

func TestNewJsonArrayEmpty(t *testing.T) {
	v := io_json.NewJsonArrayEmpty()
	assert.NotEqual(t, v, nil)
	assert.Equal(t, v.Type, io_json.JsonTypeArray)
}

func TestNewJsonNil(t *testing.T) {
	v, err := io_json.NewJson(nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, v.Type, io_json.JsonTypeNull)
}

func TestNewJson_Boolean(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		v, err := io_json.NewJson(true)
		assert.Equal(t, err, nil)
		assert.Equal(t, v.Type, io_json.JsonTypeBoolean)
		b, err := v.ToBool()
		assert.Equal(t, err, nil)
		assert.Equal(t, b, true)
	})
	t.Run("false", func(t *testing.T) {
		v, err := io_json.NewJson(false)
		assert.Equal(t, err, nil)
		assert.Equal(t, v.Type, io_json.JsonTypeBoolean)
		b, err := v.ToBool()
		assert.Equal(t, err, nil)
		assert.Equal(t, b, false)
	})
}
func TestNewJson_Int(t *testing.T) {
	i := int64(123)
	v, err := io_json.NewJson(i)
	assert.Equal(t, err, nil)
	assert.Equal(t, v.Type, io_json.JsonTypeNumber)
	a, err := v.ToInt64()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, i)
}

func TestNewJson_Float(t *testing.T) {
	var f float64 = -123.45
	v, err := io_json.NewJson(f)
	assert.Equal(t, err, nil)
	assert.Equal(t, v.Type, io_json.JsonTypeNumber)
	a, err := v.ToFloat64()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, f)
}

func TestNewJson_JsonNumber(t *testing.T) {
	t.Run("int64", func(t *testing.T) {
		v, err := io_json.NewJson(json.Number("123"))
		assert.Equal(t, err, nil)
		assert.Equal(t, v.Type, io_json.JsonTypeNumber)
		a, err := v.ToInt64()
		assert.Equal(t, err, nil)
		assert.Equal(t, a, int64(123))
	})
	t.Run("float64", func(t *testing.T) {
		v, err := io_json.NewJson(json.Number("-123.45"))
		assert.Equal(t, err, nil)
		assert.Equal(t, v.Type, io_json.JsonTypeNumber)
		a, err := v.ToFloat64()
		assert.Equal(t, err, nil)
		assert.Equal(t, a, float64(-123.45))
	})
}

func TestNewJson_Object(t *testing.T) {
	v, err := io_json.NewJson(map[string]any{
		"x": map[string]any{
			"a": int64(123),
			"b": float64(-123.45),
			"c": "abc",
			"d": nil,
			"e": true,
			"f": false,
			"g": map[string]any{},
			"h": []any{},
		},
		"y": []any{
			int64(123),
			float64(-123.45),
			"abc",
			nil,
			true,
			false,
			map[string]any{},
			[]any{},
		},
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, v.Type, io_json.JsonTypeObject)
}

func TestNewJson_Array(t *testing.T) {
	v, err := io_json.NewJson([]any{
		map[string]any{
			"a": int64(123),
			"b": float64(-123.45),
			"c": "abc",
			"d": nil,
			"e": true,
			"f": false,
			"g": map[string]any{},
			"h": []any{},
		},
		[]any{
			int64(123),
			float64(-123.45),
			"abc",
			nil,
			true,
			false,
			map[string]any{},
			[]any{},
		},
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, v.Type, io_json.JsonTypeArray)
}

func TestToString_Null(t *testing.T) {
	v := io_json.NewJsonNull()
	a, err := v.ToString()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, "null")
}

func TestToString_Int64(t *testing.T) {
	v := io_json.NewJsonNumberInt64(123)
	a, err := v.ToString()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, "123")
}

func TestToString_Float64(t *testing.T) {
	v := io_json.NewJsonNumberFloat64(-123.45)
	a, err := v.ToString()
	assert.Equal(t, err, nil)
	if !strings.HasPrefix(a, "-123.45") {
		t.Errorf("expect: %v, actual: %v", "-123.45", a)
	}
}

func TestToString_String(t *testing.T) {
	v := io_json.NewJsonString("abc")
	a, err := v.ToString()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, "abc")
}

func TestToString_True(t *testing.T) {
	v := io_json.NewJsonBoolean(true)
	a, err := v.ToString()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, "true")
}

func TestToString_False(t *testing.T) {
	v := io_json.NewJsonBoolean(false)
	a, err := v.ToString()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, "false")
}

func TestToString_Object(t *testing.T) {
	v := io_json.NewJsonObjectEmpty()
	_, err := v.ToString()
	assert.NotEqual(t, err, nil)
}

func TestToString_Array(t *testing.T) {
	v := io_json.NewJsonArrayEmpty()
	_, err := v.ToString()
	assert.NotEqual(t, err, nil)
}

func TestToInt64_Null(t *testing.T) {
	v := io_json.NewJsonNull()
	a, err := v.ToInt64()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, int64(0))
}

func TestToInt64_Int64(t *testing.T) {
	v := io_json.NewJsonNumberInt64(123)
	a, err := v.ToInt64()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, int64(123))
}

func TestToInt64_Float64(t *testing.T) {
	t.Run("cannot parse as integer", func(t *testing.T) {
		v := io_json.NewJsonNumberFloat64(-123.45)
		_, err := v.ToInt64()
		assert.NotEqual(t, err, nil)
	})
	t.Run("integer", func(t *testing.T) {
		v := io_json.NewJsonNumberFloat64(123.0)
		a, err := v.ToInt64()
		assert.Equal(t, err, nil)
		assert.Equal(t, a, int64(123))
	})
}

func TestToInt64_String(t *testing.T) {
	t.Run("cannot parse as integer", func(t *testing.T) {
		v := io_json.NewJsonString("abc")
		_, err := v.ToInt64()
		assert.NotEqual(t, err, nil)
	})
	t.Run("parse as integer", func(t *testing.T) {
		v := io_json.NewJsonString("123")
		a, err := v.ToInt64()
		assert.Equal(t, err, nil)
		assert.Equal(t, a, int64(123))
	})
}

func TestToInt64_Boolean(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		v := io_json.NewJsonBoolean(true)
		a, err := v.ToInt64()
		assert.Equal(t, err, nil)
		assert.Equal(t, a, int64(1))

	})
	t.Run("false", func(t *testing.T) {
		v := io_json.NewJsonBoolean(false)
		a, err := v.ToInt64()
		assert.Equal(t, err, nil)
		assert.Equal(t, a, int64(0))
	})
}

func TestToInt64_Object(t *testing.T) {
	v := io_json.NewJsonObjectEmpty()
	_, err := v.ToInt64()
	assert.NotEqual(t, err, nil)
}

func TestToInt64_Array(t *testing.T) {
	v := io_json.NewJsonArrayEmpty()
	_, err := v.ToInt64()
	assert.NotEqual(t, err, nil)
}

func TestToFloat64_Null(t *testing.T) {
	v := io_json.NewJsonNull()
	a, err := v.ToFloat64()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, float64(0))
}

func TestToFloat64_Int64(t *testing.T) {
	v := io_json.NewJsonNumberInt64(123)
	a, err := v.ToFloat64()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, float64(123.0))
}

func TestToFloat64_Float64(t *testing.T) {
	v := io_json.NewJsonNumberFloat64(-123.45)
	a, err := v.ToFloat64()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, float64(-123.45))
}

func TestToFloat64_String(t *testing.T) {
	t.Run("cannot parse as float", func(t *testing.T) {
		v := io_json.NewJsonString("abc")
		_, err := v.ToFloat64()
		assert.NotEqual(t, err, nil)
	})
	t.Run("parse as float", func(t *testing.T) {
		v := io_json.NewJsonString("-123.45")
		a, err := v.ToFloat64()
		assert.Equal(t, err, nil)
		assert.Equal(t, a, float64(-123.45))
	})
}

func TestToFloat64_Boolean(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		v := io_json.NewJsonBoolean(true)
		a, err := v.ToFloat64()
		assert.Equal(t, err, nil)
		assert.Equal(t, a, float64(1.0))
	})
	t.Run("false", func(t *testing.T) {
		v := io_json.NewJsonBoolean(false)
		a, err := v.ToFloat64()
		assert.Equal(t, err, nil)
		assert.Equal(t, a, float64(0.0))
	})
}

func TestToFloat64_Object(t *testing.T) {
	v := io_json.NewJsonObjectEmpty()
	_, err := v.ToFloat64()
	assert.NotEqual(t, err, nil)
}

func TestToFloat64_Array(t *testing.T) {
	v := io_json.NewJsonArrayEmpty()
	_, err := v.ToFloat64()
	assert.NotEqual(t, err, nil)
}

func TestToBool_Null(t *testing.T) {
	v := io_json.NewJsonNull()
	a, err := v.ToBool()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, false)
}

func TestToBool_Int64(t *testing.T) {
	t.Run("parse 123 as true", func(t *testing.T) {
		v := io_json.NewJsonNumberInt64(123)
		a, err := v.ToBool()
		assert.Equal(t, err, nil)
		assert.Equal(t, a, true)
	})
	t.Run("parse 1 as true", func(t *testing.T) {
		v := io_json.NewJsonNumberInt64(1)
		a, err := v.ToBool()
		assert.Equal(t, err, nil)
		assert.Equal(t, a, true)
	})
	t.Run("parse 0 as false", func(t *testing.T) {
		v := io_json.NewJsonNumberInt64(0)
		a, err := v.ToBool()
		assert.Equal(t, err, nil)
		assert.Equal(t, a, false)
	})
}

func TestToBool_Float64(t *testing.T) {
	t.Run("parse -123.45 as true", func(t *testing.T) {
		v := io_json.NewJsonNumberFloat64(-123.45)
		a, err := v.ToBool()
		assert.Equal(t, err, nil)
		assert.Equal(t, a, true)
	})
	t.Run("parse 1 as true", func(t *testing.T) {
		v := io_json.NewJsonNumberFloat64(1)
		a, err := v.ToBool()
		assert.Equal(t, err, nil)
		assert.Equal(t, a, true)
	})
	t.Run("parse 0 as false", func(t *testing.T) {
		v := io_json.NewJsonNumberFloat64(0)
		a, err := v.ToBool()
		assert.Equal(t, err, nil)
		assert.Equal(t, a, false)
	})
}

func TestToBool_String(t *testing.T) {
	t.Run(`cannot parse "abc"`, func(t *testing.T) {
		v := io_json.NewJsonString("abc")
		_, err := v.ToBool()
		assert.NotEqual(t, err, nil)
	})
	t.Run(`cannot parse "null"`, func(t *testing.T) {
		v := io_json.NewJsonString("null")
		_, err := v.ToBool()
		assert.NotEqual(t, err, nil)
	})
	t.Run(`cannot parse ""`, func(t *testing.T) {
		v := io_json.NewJsonString("")
		_, err := v.ToBool()
		assert.NotEqual(t, err, nil)
	})
	t.Run(`parse "true" as true`, func(t *testing.T) {
		v := io_json.NewJsonString("true")
		a, err := v.ToBool()
		assert.Equal(t, err, nil)
		assert.Equal(t, a, true)
	})
	t.Run(`parse "false" as false`, func(t *testing.T) {
		v := io_json.NewJsonString("false")
		a, err := v.ToBool()
		assert.Equal(t, err, nil)
		assert.Equal(t, a, false)
	})
}

func TestToBool_True(t *testing.T) {
	v := io_json.NewJsonBoolean(true)
	a, err := v.ToBool()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, true)
}

func TestToBool_False(t *testing.T) {
	v := io_json.NewJsonBoolean(false)
	a, err := v.ToBool()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, false)
}

func TestToBool_Object(t *testing.T) {
	v := io_json.NewJsonObjectEmpty()
	_, err := v.ToBool()
	assert.NotEqual(t, err, nil)
}

func TestToBool_Array(t *testing.T) {
	v := io_json.NewJsonArrayEmpty()
	_, err := v.ToBool()
	assert.NotEqual(t, err, nil)
}

func TestJsonObject(t *testing.T) {
	t.Run("AsObject for empty", func(t *testing.T) {
		v := io_json.NewJsonObjectEmpty()
		_, err := v.AsObject()
		assert.Equal(t, err, nil)
	})
	t.Run("Keys for empty", func(t *testing.T) {
		v := io_json.NewJsonObjectEmpty()
		o, _ := v.AsObject()
		a := o.Keys()
		assert.Equal(t, len(a), 0)
	})
	v, _ := io_json.NewJson(map[string]any{
		"x": map[string]any{
			"a": int64(123),
			"b": float64(-123.45),
			"c": "abc",
			"d": nil,
			"e": true,
			"f": false,
			"g": map[string]any{},
			"h": []any{},
		},
		"y": []any{
			int64(123),
			float64(-123.45),
			"abc",
			nil,
			true,
			false,
			map[string]any{},
			[]any{},
		},
	})
	t.Run("Keys", func(t *testing.T) {
		obj, err := v.AsObject()
		assert.Equal(t, err, nil)
		aKeys := obj.Keys()
		assert.Equal(t, len(aKeys), 2)
		assert.Equal(t, slices.Contains(aKeys, "x"), true)
		assert.Equal(t, slices.Contains(aKeys, "y"), true)
	})
	t.Run("Get", func(t *testing.T) {
		obj, err := v.AsObject()
		assert.Equal(t, err, nil)
		aObjX, err := obj.Get("x")
		assert.Equal(t, err, nil)
		assert.Equal(t, aObjX.Type, io_json.JsonTypeObject)
		aObjY, err := obj.Get("y")
		assert.Equal(t, err, nil)
		assert.Equal(t, aObjY.Type, io_json.JsonTypeArray)

		t.Run("child object Keys", func(t *testing.T) {
			objX, err := aObjX.AsObject()
			assert.Equal(t, err, nil)
			aObjXKeys := objX.Keys()
			assert.Equal(t, len(aObjXKeys), 8)
			assert.Equal(t, slices.Contains(aObjXKeys, "a"), true)
			assert.Equal(t, slices.Contains(aObjXKeys, "b"), true)
			assert.Equal(t, slices.Contains(aObjXKeys, "c"), true)
			assert.Equal(t, slices.Contains(aObjXKeys, "d"), true)
			assert.Equal(t, slices.Contains(aObjXKeys, "e"), true)
			assert.Equal(t, slices.Contains(aObjXKeys, "f"), true)
			assert.Equal(t, slices.Contains(aObjXKeys, "g"), true)
			assert.Equal(t, slices.Contains(aObjXKeys, "h"), true)
		})
		t.Run("child object Get", func(t *testing.T) {
			objX, err := aObjX.AsObject()
			assert.Equal(t, err, nil)
			aObjXA, err := objX.Get("a")
			assert.Equal(t, err, nil)
			assert.Equal(t, aObjXA.Type, io_json.JsonTypeNumber)
			aObjXB, err := objX.Get("b")
			assert.Equal(t, err, nil)
			assert.Equal(t, aObjXB.Type, io_json.JsonTypeNumber)
			aObjXC, err := objX.Get("c")
			assert.Equal(t, err, nil)
			assert.Equal(t, aObjXC.Type, io_json.JsonTypeString)
			aObjXD, err := objX.Get("d")
			assert.Equal(t, err, nil)
			assert.Equal(t, aObjXD.Type, io_json.JsonTypeNull)
			aObjXE, err := objX.Get("e")
			assert.Equal(t, err, nil)
			assert.Equal(t, aObjXE.Type, io_json.JsonTypeBoolean)
			aObjXF, err := objX.Get("f")
			assert.Equal(t, err, nil)
			assert.Equal(t, aObjXF.Type, io_json.JsonTypeBoolean)
			aObjXG, err := objX.Get("g")
			assert.Equal(t, err, nil)
			assert.Equal(t, aObjXG.Type, io_json.JsonTypeObject)
			aObjXH, err := objX.Get("h")
			assert.Equal(t, err, nil)
			assert.Equal(t, aObjXH.Type, io_json.JsonTypeArray)
		})
	})

	t.Run("Set", func(t *testing.T) {
		obj, _ := io_json.NewJsonObjectEmpty().AsObject()
		objX, _ := io_json.NewJsonObjectEmpty().AsObject()
		a := obj.Set("x", objX.
			Set("a", io_json.NewJsonNumberInt64(123)).
			Set("b", io_json.NewJsonNumberFloat64(-123.45)).
			Set("c", io_json.NewJsonString("abc")).
			Set("d", io_json.NewJsonNull()).
			Set("e", io_json.NewJsonBoolean(true)).
			Set("f", io_json.NewJsonBoolean(false)).
			Set("g", io_json.NewJsonObjectEmpty()).
			Set("h", io_json.NewJsonArrayEmpty()).
			AsJsonValue())
		aObjX, err := a.Get("x")
		assert.Equal(t, err, nil)
		assert.Equal(t, aObjX.Type, io_json.JsonTypeObject)

		t.Run("check child object keys", func(t *testing.T) {
			objX, _ := aObjX.AsObject()
			aObjXKeys := objX.Keys()
			assert.Equal(t, len(aObjXKeys), 8)
			assert.Equal(t, slices.Contains(aObjXKeys, "a"), true)
			assert.Equal(t, slices.Contains(aObjXKeys, "b"), true)
			assert.Equal(t, slices.Contains(aObjXKeys, "c"), true)
			assert.Equal(t, slices.Contains(aObjXKeys, "d"), true)
			assert.Equal(t, slices.Contains(aObjXKeys, "e"), true)
			assert.Equal(t, slices.Contains(aObjXKeys, "f"), true)
			assert.Equal(t, slices.Contains(aObjXKeys, "g"), true)
			assert.Equal(t, slices.Contains(aObjXKeys, "h"), true)
		})
		t.Run("check child object values", func(t *testing.T) {
			objX, _ := aObjX.AsObject()
			aObjXA, err := objX.Get("a")
			assert.Equal(t, err, nil)
			assert.Equal(t, aObjXA.Type, io_json.JsonTypeNumber)
			aObjXB, err := objX.Get("b")
			assert.Equal(t, err, nil)
			assert.Equal(t, aObjXB.Type, io_json.JsonTypeNumber)
			aObjXC, err := objX.Get("c")
			assert.Equal(t, err, nil)
			assert.Equal(t, aObjXC.Type, io_json.JsonTypeString)
			aObjXD, err := objX.Get("d")
			assert.Equal(t, err, nil)
			assert.Equal(t, aObjXD.Type, io_json.JsonTypeNull)
			aObjXE, err := objX.Get("e")
			assert.Equal(t, err, nil)
			assert.Equal(t, aObjXE.Type, io_json.JsonTypeBoolean)
			aObjXF, err := objX.Get("f")
			assert.Equal(t, err, nil)
			assert.Equal(t, aObjXF.Type, io_json.JsonTypeBoolean)
			aObjXG, err := objX.Get("g")
			assert.Equal(t, err, nil)
			assert.Equal(t, aObjXG.Type, io_json.JsonTypeObject)
			aObjXH, err := objX.Get("h")
			assert.Equal(t, err, nil)
			assert.Equal(t, aObjXH.Type, io_json.JsonTypeArray)
		})
	})
}

func TestJsonArray(t *testing.T) {
	t.Run("AsArray for empty", func(t *testing.T) {
		v := io_json.NewJsonArrayEmpty()
		_, err := v.AsArray()
		assert.Equal(t, err, nil)
	})
	t.Run("Len for empty", func(t *testing.T) {
		v := io_json.NewJsonArrayEmpty()
		arr, _ := v.AsArray()
		assert.Equal(t, arr.Len(), 0)
	})
	v, _ := io_json.NewJson([]any{
		map[string]any{
			"a": int64(123),
			"b": float64(-123.45),
			"c": "abc",
			"d": nil,
			"e": true,
			"f": false,
			"g": map[string]any{},
			"h": []any{},
		},
		[]any{
			int64(123),
			float64(-123.45),
			"abc",
			nil,
			true,
			false,
			map[string]any{},
			[]any{},
		},
	})
	t.Run("Len for array of 2 elements", func(t *testing.T) {
		arr, err := v.AsArray()
		assert.Equal(t, err, nil)
		assert.Equal(t, arr.Len(), 2)
	})
	t.Run("Get", func(t *testing.T) {
		arr, err := v.AsArray()
		assert.Equal(t, err, nil)
		aArr0, err := arr.Get(0)
		assert.Equal(t, err, nil)
		assert.Equal(t, aArr0.Type, io_json.JsonTypeObject)
		aArr1, err := arr.Get(1)
		assert.Equal(t, err, nil)
		assert.Equal(t, aArr1.Type, io_json.JsonTypeArray)

		t.Run("child array Len", func(t *testing.T) {
			arr1, err := aArr1.AsArray()
			assert.Equal(t, err, nil)
			assert.Equal(t, arr1.Len(), 8)
		})
		t.Run("child array Get", func(t *testing.T) {
			arr1, err := aArr1.AsArray()
			assert.Equal(t, err, nil)
			aArr10, err := arr1.Get(0)
			assert.Equal(t, err, nil)
			assert.Equal(t, aArr10.Type, io_json.JsonTypeNumber)
			aArr11, err := arr1.Get(1)
			assert.Equal(t, err, nil)
			assert.Equal(t, aArr11.Type, io_json.JsonTypeNumber)
			aArr12, err := arr1.Get(2)
			assert.Equal(t, err, nil)
			assert.Equal(t, aArr12.Type, io_json.JsonTypeString)
			aArr13, err := arr1.Get(3)
			assert.Equal(t, err, nil)
			assert.Equal(t, aArr13.Type, io_json.JsonTypeNull)
			aArr14, err := arr1.Get(4)
			assert.Equal(t, err, nil)
			assert.Equal(t, aArr14.Type, io_json.JsonTypeBoolean)
			aArr15, err := arr1.Get(5)
			assert.Equal(t, err, nil)
			assert.Equal(t, aArr15.Type, io_json.JsonTypeBoolean)
			aArr16, err := arr1.Get(6)
			assert.Equal(t, err, nil)
			assert.Equal(t, aArr16.Type, io_json.JsonTypeObject)
			aArr17, err := arr1.Get(7)
			assert.Equal(t, err, nil)
			assert.Equal(t, aArr17.Type, io_json.JsonTypeArray)
		})
	})

	t.Run("Append", func(t *testing.T) {
		arr, _ := io_json.NewJsonArrayEmpty().AsArray()
		arr0, _ := io_json.NewJsonArrayEmpty().AsArray()
		a := arr.Append(arr0.
			Append(io_json.NewJsonNumberInt64(123)).
			Append(io_json.NewJsonNumberFloat64(-123.45)).
			Append(io_json.NewJsonString("abc")).
			Append(io_json.NewJsonNull()).
			Append(io_json.NewJsonBoolean(true)).
			Append(io_json.NewJsonBoolean(false)).
			Append(io_json.NewJsonObjectEmpty()).
			Append(io_json.NewJsonArrayEmpty()).
			AsJsonValue())
		aArr0, err := a.Get(0)
		assert.Equal(t, err, nil)
		assert.Equal(t, aArr0.Type, io_json.JsonTypeArray)

		t.Run("check child array length", func(t *testing.T) {
			arr0, _ := aArr0.AsArray()
			assert.Equal(t, arr0.Len(), 8)
		})
		t.Run("check child array values", func(t *testing.T) {
			arr0, _ := aArr0.AsArray()
			aArr00, err := arr0.Get(0)
			assert.Equal(t, err, nil)
			assert.Equal(t, aArr00.Type, io_json.JsonTypeNumber)
			aArr01, err := arr0.Get(1)
			assert.Equal(t, err, nil)
			assert.Equal(t, aArr01.Type, io_json.JsonTypeNumber)
			aArr02, err := arr0.Get(2)
			assert.Equal(t, err, nil)
			assert.Equal(t, aArr02.Type, io_json.JsonTypeString)
			aArr03, err := arr0.Get(3)
			assert.Equal(t, err, nil)
			assert.Equal(t, aArr03.Type, io_json.JsonTypeNull)
			aArr04, err := arr0.Get(4)
			assert.Equal(t, err, nil)
			assert.Equal(t, aArr04.Type, io_json.JsonTypeBoolean)
			aArr05, err := arr0.Get(5)
			assert.Equal(t, err, nil)
			assert.Equal(t, aArr05.Type, io_json.JsonTypeBoolean)
			aArr06, err := arr0.Get(6)
			assert.Equal(t, err, nil)
			assert.Equal(t, aArr06.Type, io_json.JsonTypeObject)
			aArr07, err := arr0.Get(7)
			assert.Equal(t, err, nil)
			assert.Equal(t, aArr07.Type, io_json.JsonTypeArray)
		})
	})
}

/*
func (v *JsonValue) ObjectKeys() (keys []string, err error) {
	if v.Type != io_json.JsonTypeObject {
		return nil, fmt.Errorf("ObjectKeys must be called with JsonValue of type io_json.JsonTypeObject")
	}
	for k := range v.objectValue {
		keys = append(keys, k)
	}
	return keys, nil
}
func (v *JsonValue) ObjectGet(key string) (*JsonValue, error) {
	if v.Type != io_json.JsonTypeObject {
		return nil, fmt.Errorf("ObjectGet must be called with JsonValue of type io_json.JsonTypeObject")
	}
	val, ok := v.objectValue[key]
	if !ok {
		return nil, fmt.Errorf("value not found for key %s", key)
	}
	return val, nil
}

func (v *JsonValue) ObjectSet(key string, val *JsonValue) (err error) {
	if v.Type != io_json.JsonTypeObject {
		return fmt.Errorf("ObjectSet must be called with JsonValue of type io_json.JsonTypeObject")
	}
	if val == nil {
		val, err = io_json.NewJson(nil)
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
	if v.Type != io_json.JsonTypeArray {
		return 0, fmt.Errorf("ArrayLen must be called with JsonValue with type io_json.JsonTypeArray")
	}
	return len([]*JsonValue(v.arrayValue)), nil
}

func (v *JsonValue) ArrayGet(i int) (val *JsonValue, err error) {
	if v.Type != io_json.JsonTypeArray {
		return nil, fmt.Errorf("ArrayGet must be called with JsonValue with type io_json.JsonTypeArray")
	}
	return v.arrayValue[i], nil
}

func (v *JsonValue) ArraySet(i int, val *JsonValue) (err error) {
	if v.Type != io_json.JsonTypeArray {
		return fmt.Errorf("ArraySet must be called with JsonValue with type io_json.JsonTypeArray")
	}
	if val == nil {
		val, err = io_json.NewJson(nil)
		if err != nil {
			return err
		}
	}
	v.arrayValue[i] = val
	return nil
}

func (v *JsonValue) ArrayAppend(val *JsonValue) (err error) {
	if v.Type != io_json.JsonTypeArray {
		return fmt.Errorf("ArrayAppend must be called with JsonValue with type io_json.JsonTypeArray")
	}
	if val == nil {
		val, err = io_json.NewJson(nil)
		if err != nil {
			return err
		}
	}
	a := []*JsonValue(v.arrayValue)
	a = append(a, val)
	v.arrayValue = a
	return nil
}
*/
