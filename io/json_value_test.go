package io

import (
	"encoding/json"
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

func assertEqual(t *testing.T, actual any, expect any) {
	if actual != expect {
		t.Errorf(`expect: %v, actual: %v`, expect, actual)
	}
}
func assertNotEqual(t *testing.T, actual any, expect any) {
	if actual == expect {
		t.Errorf(`expect: not %v, actual: %v`, expect, actual)
	}
}
func TestNewJsonString(t *testing.T) {
	s := "abc"
	v := NewJsonString(s)
	assertNotEqual(t, v, nil)
	assertEqual(t, v.Type, JsonTypeString)
	a, err := v.ToString()
	assertEqual(t, err, nil)
	assertEqual(t, a, s)
}

func TestNewJsonBooleanTrue(t *testing.T) {
	b := true
	v := NewJsonBoolean(b)
	assertNotEqual(t, v, nil)
	assertEqual(t, v.Type, JsonTypeBoolean)
	a, err := v.ToBool()
	assertEqual(t, err, nil)
	assertEqual(t, a, b)
}

func TestNewJsonBooleanFalse(t *testing.T) {
	b := false
	v := NewJsonBoolean(b)
	assertNotEqual(t, v, nil)
	assertEqual(t, v.Type, JsonTypeBoolean)
	a, err := v.ToBool()
	assertEqual(t, err, nil)
	assertEqual(t, a, b)
}

func TestNewJsonNumberInt64(t *testing.T) {
	i := int64(123)
	v := NewJsonNumberInt64(i)
	assertNotEqual(t, v, nil)
	assertEqual(t, v.Type, JsonTypeNumber)
	a, err := v.ToInt64()
	assertEqual(t, err, nil)
	assertEqual(t, a, i)
}

func TestNewJsonNumberFloat64(t *testing.T) {
	f := -123.45
	v := NewJsonNumberFloat64(f)
	assertNotEqual(t, v, nil)
	assertEqual(t, v.Type, JsonTypeNumber)
	a, err := v.ToFloat64()
	assertEqual(t, err, nil)
	assertEqual(t, a, f)
}

func TestNewJsonNull(t *testing.T) {
	v := NewJsonNull()
	assertNotEqual(t, v, nil)
	assertEqual(t, v.Type, JsonTypeNull)
}

func TestNewJsonObjectEmpty(t *testing.T) {
	v := NewJsonObjectEmpty()
	assertNotEqual(t, v, nil)
	assertEqual(t, v.Type, JsonTypeObject)
}

func TestNewJsonArrayEmpty(t *testing.T) {
	v := NewJsonArrayEmpty()
	assertNotEqual(t, v, nil)
	assertEqual(t, v.Type, JsonTypeArray)
}

func TestNewJsonNil(t *testing.T) {
	v, err := NewJson(nil)
	assertEqual(t, err, nil)
	assertEqual(t, v.Type, JsonTypeNull)
}

func TestNewJson_Boolean(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		v, err := NewJson(true)
		assertEqual(t, err, nil)
		assertEqual(t, v.Type, JsonTypeBoolean)
		b, err := v.ToBool()
		assertEqual(t, err, nil)
		assertEqual(t, b, true)
	})
	t.Run("false", func(t *testing.T) {
		v, err := NewJson(false)
		assertEqual(t, err, nil)
		assertEqual(t, v.Type, JsonTypeBoolean)
		b, err := v.ToBool()
		assertEqual(t, err, nil)
		assertEqual(t, b, false)
	})
}
func TestNewJson_Int(t *testing.T) {
	i := int64(123)
	v, err := NewJson(i)
	assertEqual(t, err, nil)
	assertEqual(t, v.Type, JsonTypeNumber)
	a, err := v.ToInt64()
	assertEqual(t, err, nil)
	assertEqual(t, a, i)
}

func TestNewJson_Float(t *testing.T) {
	var f float64 = -123.45
	v, err := NewJson(f)
	assertEqual(t, err, nil)
	assertEqual(t, v.Type, JsonTypeNumber)
	a, err := v.ToFloat64()
	assertEqual(t, err, nil)
	assertEqual(t, a, f)
}

func TestNewJson_JsonNumber(t *testing.T) {
	t.Run("int64", func(t *testing.T) {
		v, err := NewJson(json.Number("123"))
		assertEqual(t, err, nil)
		assertEqual(t, v.Type, JsonTypeNumber)
		a, err := v.ToInt64()
		assertEqual(t, err, nil)
		assertEqual(t, a, int64(123))
	})
	t.Run("float64", func(t *testing.T) {
		v, err := NewJson(json.Number("-123.45"))
		assertEqual(t, err, nil)
		assertEqual(t, v.Type, JsonTypeNumber)
		a, err := v.ToFloat64()
		assertEqual(t, err, nil)
		assertEqual(t, a, float64(-123.45))
	})
}

func TestNewJson_Object(t *testing.T) {
	v, err := NewJson(map[string]any{
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
	assertEqual(t, err, nil)
	assertEqual(t, v.Type, JsonTypeObject)
}

func TestNewJson_Array(t *testing.T) {
	v, err := NewJson([]any{
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
	assertEqual(t, err, nil)
	assertEqual(t, v.Type, JsonTypeArray)
}

func TestToString_Null(t *testing.T) {
	v := NewJsonNull()
	a, err := v.ToString()
	assertEqual(t, err, nil)
	assertEqual(t, a, "null")
}

func TestToString_Int64(t *testing.T) {
	v := NewJsonNumberInt64(123)
	a, err := v.ToString()
	assertEqual(t, err, nil)
	assertEqual(t, a, "123")
}

func TestToString_Float64(t *testing.T) {
	v := NewJsonNumberFloat64(-123.45)
	a, err := v.ToString()
	assertEqual(t, err, nil)
	if !strings.HasPrefix(a, "-123.45") {
		t.Errorf("expect: %v, actual: %v", "-123.45", a)
	}
}

func TestToString_String(t *testing.T) {
	v := NewJsonString("abc")
	a, err := v.ToString()
	assertEqual(t, err, nil)
	assertEqual(t, a, "abc")
}

func TestToString_True(t *testing.T) {
	v := NewJsonBoolean(true)
	a, err := v.ToString()
	assertEqual(t, err, nil)
	assertEqual(t, a, "true")
}

func TestToString_False(t *testing.T) {
	v := NewJsonBoolean(false)
	a, err := v.ToString()
	assertEqual(t, err, nil)
	assertEqual(t, a, "false")
}

func TestToString_Object(t *testing.T) {
	v := NewJsonObjectEmpty()
	_, err := v.ToString()
	assertNotEqual(t, err, nil)
}

func TestToString_Array(t *testing.T) {
	v := NewJsonArrayEmpty()
	_, err := v.ToString()
	assertNotEqual(t, err, nil)
}

func TestToInt64_Null(t *testing.T) {
	v := NewJsonNull()
	a, err := v.ToInt64()
	assertEqual(t, err, nil)
	assertEqual(t, a, int64(0))
}

func TestToInt64_Int64(t *testing.T) {
	v := NewJsonNumberInt64(123)
	a, err := v.ToInt64()
	assertEqual(t, err, nil)
	assertEqual(t, a, int64(123))
}

func TestToInt64_Float64(t *testing.T) {
	t.Run("cannot parse as integer", func(t *testing.T) {
		v := NewJsonNumberFloat64(-123.45)
		_, err := v.ToInt64()
		assertNotEqual(t, err, nil)
	})
	t.Run("integer", func(t *testing.T) {
		v := NewJsonNumberFloat64(123.0)
		a, err := v.ToInt64()
		assertEqual(t, err, nil)
		assertEqual(t, a, int64(123))
	})
}

func TestToInt64_String(t *testing.T) {
	t.Run("cannot parse as integer", func(t *testing.T) {
		v := NewJsonString("abc")
		_, err := v.ToInt64()
		assertNotEqual(t, err, nil)
	})
	t.Run("parse as integer", func(t *testing.T) {
		v := NewJsonString("123")
		a, err := v.ToInt64()
		assertEqual(t, err, nil)
		assertEqual(t, a, int64(123))
	})
}

func TestToInt64_Boolean(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		v := NewJsonBoolean(true)
		a, err := v.ToInt64()
		assertEqual(t, err, nil)
		assertEqual(t, a, int64(1))

	})
	t.Run("false", func(t *testing.T) {
		v := NewJsonBoolean(false)
		a, err := v.ToInt64()
		assertEqual(t, err, nil)
		assertEqual(t, a, int64(0))
	})
}

func TestToInt64_Object(t *testing.T) {
	v := NewJsonObjectEmpty()
	_, err := v.ToInt64()
	assertNotEqual(t, err, nil)
}

func TestToInt64_Array(t *testing.T) {
	v := NewJsonArrayEmpty()
	_, err := v.ToInt64()
	assertNotEqual(t, err, nil)
}

func TestToFloat64_Null(t *testing.T) {
	v := NewJsonNull()
	a, err := v.ToFloat64()
	assertEqual(t, err, nil)
	assertEqual(t, a, float64(0))
}

func TestToFloat64_Int64(t *testing.T) {
	v := NewJsonNumberInt64(123)
	a, err := v.ToFloat64()
	assertEqual(t, err, nil)
	assertEqual(t, a, float64(123.0))
}

func TestToFloat64_Float64(t *testing.T) {
	v := NewJsonNumberFloat64(-123.45)
	a, err := v.ToFloat64()
	assertEqual(t, err, nil)
	assertEqual(t, a, float64(-123.45))
}

func TestToFloat64_String(t *testing.T) {
	t.Run("cannot parse as float", func(t *testing.T) {
		v := NewJsonString("abc")
		_, err := v.ToFloat64()
		assertNotEqual(t, err, nil)
	})
	t.Run("parse as float", func(t *testing.T) {
		v := NewJsonString("-123.45")
		a, err := v.ToFloat64()
		assertEqual(t, err, nil)
		assertEqual(t, a, float64(-123.45))
	})
}

func TestToFloat64_Boolean(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		v := NewJsonBoolean(true)
		a, err := v.ToFloat64()
		assertEqual(t, err, nil)
		assertEqual(t, a, float64(1.0))
	})
	t.Run("false", func(t *testing.T) {
		v := NewJsonBoolean(false)
		a, err := v.ToFloat64()
		assertEqual(t, err, nil)
		assertEqual(t, a, float64(0.0))
	})
}

func TestToFloat64_Object(t *testing.T) {
	v := NewJsonObjectEmpty()
	_, err := v.ToFloat64()
	assertNotEqual(t, err, nil)
}

func TestToFloat64_Array(t *testing.T) {
	v := NewJsonArrayEmpty()
	_, err := v.ToFloat64()
	assertNotEqual(t, err, nil)
}

func TestToBool_Null(t *testing.T) {
	v := NewJsonNull()
	a, err := v.ToBool()
	assertEqual(t, err, nil)
	assertEqual(t, a, false)
}

func TestToBool_Int64(t *testing.T) {
	t.Run("parse 123 as true", func(t *testing.T) {
		v := NewJsonNumberInt64(123)
		a, err := v.ToBool()
		assertEqual(t, err, nil)
		assertEqual(t, a, true)
	})
	t.Run("parse 1 as true", func(t *testing.T) {
		v := NewJsonNumberInt64(1)
		a, err := v.ToBool()
		assertEqual(t, err, nil)
		assertEqual(t, a, true)
	})
	t.Run("parse 0 as false", func(t *testing.T) {
		v := NewJsonNumberInt64(0)
		a, err := v.ToBool()
		assertEqual(t, err, nil)
		assertEqual(t, a, false)
	})
}

func TestToBool_Float64(t *testing.T) {
	t.Run("parse -123.45 as true", func(t *testing.T) {
		v := NewJsonNumberFloat64(-123.45)
		a, err := v.ToBool()
		assertEqual(t, err, nil)
		assertEqual(t, a, true)
	})
	t.Run("parse 1 as true", func(t *testing.T) {
		v := NewJsonNumberFloat64(1)
		a, err := v.ToBool()
		assertEqual(t, err, nil)
		assertEqual(t, a, true)
	})
	t.Run("parse 0 as false", func(t *testing.T) {
		v := NewJsonNumberFloat64(0)
		a, err := v.ToBool()
		assertEqual(t, err, nil)
		assertEqual(t, a, false)
	})
}

func TestToBool_String(t *testing.T) {
	t.Run(`cannot parse "abc"`, func(t *testing.T) {
		v := NewJsonString("abc")
		_, err := v.ToBool()
		assertNotEqual(t, err, nil)
	})
	t.Run(`cannot parse "null"`, func(t *testing.T) {
		v := NewJsonString("null")
		_, err := v.ToBool()
		assertNotEqual(t, err, nil)
	})
	t.Run(`cannot parse ""`, func(t *testing.T) {
		v := NewJsonString("")
		_, err := v.ToBool()
		assertNotEqual(t, err, nil)
	})
	t.Run(`parse "true" as true`, func(t *testing.T) {
		v := NewJsonString("true")
		a, err := v.ToBool()
		assertEqual(t, err, nil)
		assertEqual(t, a, true)
	})
	t.Run(`parse "false" as false`, func(t *testing.T) {
		v := NewJsonString("false")
		a, err := v.ToBool()
		assertEqual(t, err, nil)
		assertEqual(t, a, false)
	})
}

func TestToBool_True(t *testing.T) {
	v := NewJsonBoolean(true)
	a, err := v.ToBool()
	assertEqual(t, err, nil)
	assertEqual(t, a, true)
}

func TestToBool_False(t *testing.T) {
	v := NewJsonBoolean(false)
	a, err := v.ToBool()
	assertEqual(t, err, nil)
	assertEqual(t, a, false)
}

func TestToBool_Object(t *testing.T) {
	v := NewJsonObjectEmpty()
	_, err := v.ToBool()
	assertNotEqual(t, err, nil)
}

func TestToBool_Array(t *testing.T) {
	v := NewJsonArrayEmpty()
	_, err := v.ToBool()
	assertNotEqual(t, err, nil)
}

func TestJsonObject(t *testing.T) {
	t.Run("AsObject for empty", func(t *testing.T) {
		v := NewJsonObjectEmpty()
		_, err := v.AsObject()
		assertEqual(t, err, nil)
	})
	t.Run("Keys for empty", func(t *testing.T) {
		v := NewJsonObjectEmpty()
		o, _ := v.AsObject()
		a := o.Keys()
		assertEqual(t, len(a), 0)
	})
	v, _ := NewJson(map[string]any{
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
		assertEqual(t, err, nil)
		aKeys := obj.Keys()
		assertEqual(t, len(aKeys), 2)
		assertEqual(t, slices.Contains(aKeys, "x"), true)
		assertEqual(t, slices.Contains(aKeys, "y"), true)
	})
	t.Run("Get", func(t *testing.T) {
		obj, err := v.AsObject()
		assertEqual(t, err, nil)
		aObjX, err := obj.Get("x")
		assertEqual(t, err, nil)
		assertEqual(t, aObjX.Type, JsonTypeObject)
		aObjY, err := obj.Get("y")
		assertEqual(t, err, nil)
		assertEqual(t, aObjY.Type, JsonTypeArray)

		t.Run("child object Keys", func(t *testing.T) {
			objX, err := aObjX.AsObject()
			assertEqual(t, err, nil)
			aObjXKeys := objX.Keys()
			assertEqual(t, len(aObjXKeys), 8)
			assertEqual(t, slices.Contains(aObjXKeys, "a"), true)
			assertEqual(t, slices.Contains(aObjXKeys, "b"), true)
			assertEqual(t, slices.Contains(aObjXKeys, "c"), true)
			assertEqual(t, slices.Contains(aObjXKeys, "d"), true)
			assertEqual(t, slices.Contains(aObjXKeys, "e"), true)
			assertEqual(t, slices.Contains(aObjXKeys, "f"), true)
			assertEqual(t, slices.Contains(aObjXKeys, "g"), true)
			assertEqual(t, slices.Contains(aObjXKeys, "h"), true)
		})
		t.Run("child object Get", func(t *testing.T) {
			objX, err := aObjX.AsObject()
			assertEqual(t, err, nil)
			aObjXA, err := objX.Get("a")
			assertEqual(t, err, nil)
			assertEqual(t, aObjXA.Type, JsonTypeNumber)
			aObjXB, err := objX.Get("b")
			assertEqual(t, err, nil)
			assertEqual(t, aObjXB.Type, JsonTypeNumber)
			aObjXC, err := objX.Get("c")
			assertEqual(t, err, nil)
			assertEqual(t, aObjXC.Type, JsonTypeString)
			aObjXD, err := objX.Get("d")
			assertEqual(t, err, nil)
			assertEqual(t, aObjXD.Type, JsonTypeNull)
			aObjXE, err := objX.Get("e")
			assertEqual(t, err, nil)
			assertEqual(t, aObjXE.Type, JsonTypeBoolean)
			aObjXF, err := objX.Get("f")
			assertEqual(t, err, nil)
			assertEqual(t, aObjXF.Type, JsonTypeBoolean)
			aObjXG, err := objX.Get("g")
			assertEqual(t, err, nil)
			assertEqual(t, aObjXG.Type, JsonTypeObject)
			aObjXH, err := objX.Get("h")
			assertEqual(t, err, nil)
			assertEqual(t, aObjXH.Type, JsonTypeArray)
		})
	})

	t.Run("Set", func(t *testing.T) {
		obj, _ := NewJsonObjectEmpty().AsObject()
		objX, _ := NewJsonObjectEmpty().AsObject()
		a := obj.Set("x", objX.
			Set("a", NewJsonNumberInt64(123)).
			Set("b", NewJsonNumberFloat64(-123.45)).
			Set("c", NewJsonString("abc")).
			Set("d", NewJsonNull()).
			Set("e", NewJsonBoolean(true)).
			Set("f", NewJsonBoolean(false)).
			Set("g", NewJsonObjectEmpty()).
			Set("h", NewJsonArrayEmpty()).
			AsJsonValue())
		aObjX, err := a.Get("x")
		assertEqual(t, err, nil)
		assertEqual(t, aObjX.Type, JsonTypeObject)

		t.Run("check child object keys", func(t *testing.T) {
			objX, _ := aObjX.AsObject()
			aObjXKeys := objX.Keys()
			assertEqual(t, len(aObjXKeys), 8)
			assertEqual(t, slices.Contains(aObjXKeys, "a"), true)
			assertEqual(t, slices.Contains(aObjXKeys, "b"), true)
			assertEqual(t, slices.Contains(aObjXKeys, "c"), true)
			assertEqual(t, slices.Contains(aObjXKeys, "d"), true)
			assertEqual(t, slices.Contains(aObjXKeys, "e"), true)
			assertEqual(t, slices.Contains(aObjXKeys, "f"), true)
			assertEqual(t, slices.Contains(aObjXKeys, "g"), true)
			assertEqual(t, slices.Contains(aObjXKeys, "h"), true)
		})
		t.Run("check child object values", func(t *testing.T) {
			objX, _ := aObjX.AsObject()
			aObjXA, err := objX.Get("a")
			assertEqual(t, err, nil)
			assertEqual(t, aObjXA.Type, JsonTypeNumber)
			aObjXB, err := objX.Get("b")
			assertEqual(t, err, nil)
			assertEqual(t, aObjXB.Type, JsonTypeNumber)
			aObjXC, err := objX.Get("c")
			assertEqual(t, err, nil)
			assertEqual(t, aObjXC.Type, JsonTypeString)
			aObjXD, err := objX.Get("d")
			assertEqual(t, err, nil)
			assertEqual(t, aObjXD.Type, JsonTypeNull)
			aObjXE, err := objX.Get("e")
			assertEqual(t, err, nil)
			assertEqual(t, aObjXE.Type, JsonTypeBoolean)
			aObjXF, err := objX.Get("f")
			assertEqual(t, err, nil)
			assertEqual(t, aObjXF.Type, JsonTypeBoolean)
			aObjXG, err := objX.Get("g")
			assertEqual(t, err, nil)
			assertEqual(t, aObjXG.Type, JsonTypeObject)
			aObjXH, err := objX.Get("h")
			assertEqual(t, err, nil)
			assertEqual(t, aObjXH.Type, JsonTypeArray)
		})
	})
}

func TestJsonArray(t *testing.T) {
	t.Run("AsArray for empty", func(t *testing.T) {
		v := NewJsonArrayEmpty()
		_, err := v.AsArray()
		assertEqual(t, err, nil)
	})
	t.Run("Len for empty", func(t *testing.T) {
		v := NewJsonArrayEmpty()
		arr, _ := v.AsArray()
		assertEqual(t, arr.Len(), 0)
	})
	v, _ := NewJson([]any{
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
		assertEqual(t, err, nil)
		assertEqual(t, arr.Len(), 2)
	})
	t.Run("Get", func(t *testing.T) {
		arr, err := v.AsArray()
		assertEqual(t, err, nil)
		aArr0, err := arr.Get(0)
		assertEqual(t, err, nil)
		assertEqual(t, aArr0.Type, JsonTypeObject)
		aArr1, err := arr.Get(1)
		assertEqual(t, err, nil)
		assertEqual(t, aArr1.Type, JsonTypeArray)

		t.Run("child array Len", func(t *testing.T) {
			arr1, err := aArr1.AsArray()
			assertEqual(t, err, nil)
			assertEqual(t, arr1.Len(), 8)
		})
		t.Run("child array Get", func(t *testing.T) {
			arr1, err := aArr1.AsArray()
			assertEqual(t, err, nil)
			aArr10, err := arr1.Get(0)
			assertEqual(t, err, nil)
			assertEqual(t, aArr10.Type, JsonTypeNumber)
			aArr11, err := arr1.Get(1)
			assertEqual(t, err, nil)
			assertEqual(t, aArr11.Type, JsonTypeNumber)
			aArr12, err := arr1.Get(2)
			assertEqual(t, err, nil)
			assertEqual(t, aArr12.Type, JsonTypeString)
			aArr13, err := arr1.Get(3)
			assertEqual(t, err, nil)
			assertEqual(t, aArr13.Type, JsonTypeNull)
			aArr14, err := arr1.Get(4)
			assertEqual(t, err, nil)
			assertEqual(t, aArr14.Type, JsonTypeBoolean)
			aArr15, err := arr1.Get(5)
			assertEqual(t, err, nil)
			assertEqual(t, aArr15.Type, JsonTypeBoolean)
			aArr16, err := arr1.Get(6)
			assertEqual(t, err, nil)
			assertEqual(t, aArr16.Type, JsonTypeObject)
			aArr17, err := arr1.Get(7)
			assertEqual(t, err, nil)
			assertEqual(t, aArr17.Type, JsonTypeArray)
		})
	})

	t.Run("Append", func(t *testing.T) {
		arr, _ := NewJsonArrayEmpty().AsArray()
		arr0, _ := NewJsonArrayEmpty().AsArray()
		a := arr.Append(arr0.
			Append(NewJsonNumberInt64(123)).
			Append(NewJsonNumberFloat64(-123.45)).
			Append(NewJsonString("abc")).
			Append(NewJsonNull()).
			Append(NewJsonBoolean(true)).
			Append(NewJsonBoolean(false)).
			Append(NewJsonObjectEmpty()).
			Append(NewJsonArrayEmpty()).
			AsJsonValue())
		aArr0, err := a.Get(0)
		assertEqual(t, err, nil)
		assertEqual(t, aArr0.Type, JsonTypeArray)

		t.Run("check child array length", func(t *testing.T) {
			arr0, _ := aArr0.AsArray()
			assertEqual(t, arr0.Len(), 8)
		})
		t.Run("check child array values", func(t *testing.T) {
			arr0, _ := aArr0.AsArray()
			aArr00, err := arr0.Get(0)
			assertEqual(t, err, nil)
			assertEqual(t, aArr00.Type, JsonTypeNumber)
			aArr01, err := arr0.Get(1)
			assertEqual(t, err, nil)
			assertEqual(t, aArr01.Type, JsonTypeNumber)
			aArr02, err := arr0.Get(2)
			assertEqual(t, err, nil)
			assertEqual(t, aArr02.Type, JsonTypeString)
			aArr03, err := arr0.Get(3)
			assertEqual(t, err, nil)
			assertEqual(t, aArr03.Type, JsonTypeNull)
			aArr04, err := arr0.Get(4)
			assertEqual(t, err, nil)
			assertEqual(t, aArr04.Type, JsonTypeBoolean)
			aArr05, err := arr0.Get(5)
			assertEqual(t, err, nil)
			assertEqual(t, aArr05.Type, JsonTypeBoolean)
			aArr06, err := arr0.Get(6)
			assertEqual(t, err, nil)
			assertEqual(t, aArr06.Type, JsonTypeObject)
			aArr07, err := arr0.Get(7)
			assertEqual(t, err, nil)
			assertEqual(t, aArr07.Type, JsonTypeArray)
		})
	})
}

/*
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
*/
