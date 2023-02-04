package jsonio_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
	"github.com/Jumpaku/api-regression-detector/test/assert"
	"golang.org/x/exp/slices"
)

func mustNewJson(val any) *jsonio.JsonValue {
	v, err := jsonio.NewJson(val)
	if err != nil {
		panic(fmt.Sprintf("cannot create jsonio.JsonValue of %v:%T", val, val))
	}

	return v
}

func TestRow_GetColumnNames(t *testing.T) {
	v := jsonio.Row{
		"a": mustNewJson(int64(123)),
		"b": mustNewJson(float64(-123.45)),
		"c": mustNewJson("abc"),
		"d": mustNewJson(nil),
		"e": mustNewJson(true),
		"f": mustNewJson(false),
		"g": mustNewJson(map[string]any{}),
		"h": mustNewJson([]any{}),
	}
	a := v.GetColumnNames()
	assert.Equal(t, len(a), 8)
	assert.Equal(t, slices.Contains(a, "a"), true)
	assert.Equal(t, slices.Contains(a, "b"), true)
	assert.Equal(t, slices.Contains(a, "c"), true)
	assert.Equal(t, slices.Contains(a, "d"), true)
	assert.Equal(t, slices.Contains(a, "e"), true)
	assert.Equal(t, slices.Contains(a, "f"), true)
	assert.Equal(t, slices.Contains(a, "g"), true)
	assert.Equal(t, slices.Contains(a, "h"), true)
}

func TestRow_Has(t *testing.T) {
	v := jsonio.Row{
		"a": mustNewJson(int64(123)),
		"b": mustNewJson(float64(-123.45)),
		"c": mustNewJson("abc"),
		"d": mustNewJson(nil),
		"e": mustNewJson(true),
		"f": mustNewJson(false),
		"g": mustNewJson(map[string]any{}),
		"h": mustNewJson([]any{}),
	}
	assert.Equal(t, v.Has("a"), true)
	assert.Equal(t, v.Has("b"), true)
	assert.Equal(t, v.Has("c"), true)
	assert.Equal(t, v.Has("d"), true)
	assert.Equal(t, v.Has("e"), true)
	assert.Equal(t, v.Has("f"), true)
	assert.Equal(t, v.Has("g"), true)
	assert.Equal(t, v.Has("h"), true)
	assert.Equal(t, v.Has("x"), false)
	assert.Equal(t, v.Has("y"), false)
	assert.Equal(t, v.Has("z"), false)
}

func TestRow_GetColumnTypes(t *testing.T) {
	v := jsonio.Row{
		"a": mustNewJson(int64(123)),
		"b": mustNewJson(float64(-123.45)),
		"c": mustNewJson("abc"),
		"d": mustNewJson(nil),
		"e": mustNewJson(true),
		"f": mustNewJson(false),
		"g": mustNewJson(map[string]any{}),
		"h": mustNewJson([]any{}),
	}
	aA, err := v.GetJsonType("a")
	assert.Equal(t, err, nil)
	assert.Equal(t, aA, jsonio.JsonTypeNumber)

	aB, err := v.GetJsonType("b")
	assert.Equal(t, err, nil)
	assert.Equal(t, aB, jsonio.JsonTypeNumber)

	aC, err := v.GetJsonType("c")
	assert.Equal(t, err, nil)
	assert.Equal(t, aC, jsonio.JsonTypeString)

	aD, err := v.GetJsonType("d")
	assert.Equal(t, err, nil)
	assert.Equal(t, aD, jsonio.JsonTypeNull)

	aE, err := v.GetJsonType("e")
	assert.Equal(t, err, nil)
	assert.Equal(t, aE, jsonio.JsonTypeBoolean)

	aF, err := v.GetJsonType("f")
	assert.Equal(t, err, nil)
	assert.Equal(t, aF, jsonio.JsonTypeBoolean)

	aG, err := v.GetJsonType("g")
	assert.Equal(t, err, nil)
	assert.Equal(t, aG, jsonio.JsonTypeObject)

	aH, err := v.GetJsonType("h")
	assert.Equal(t, err, nil)
	assert.Equal(t, aH, jsonio.JsonTypeArray)

	_, err = v.GetJsonType("z")
	assert.NotEqual(t, err, nil)
}

func TestRow_ToString_Null(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonNull()}
	a, err := v.ToString("a")
	assert.Equal(t, err, nil)
	assert.Equal(t, a, "null")
}

func TestRow_ToString_Int64(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonNumberInt64(123)}
	a, err := v.ToString("a")
	assert.Equal(t, err, nil)
	assert.Equal(t, a, "123")
}

func TestRow_ToString_Float64(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonNumberFloat64(-123.45)}
	a, err := v.ToString("a")
	assert.Equal(t, err, nil)

	if !strings.HasPrefix(a, "-123.45") {
		t.Errorf("expect: %v, actual: %v", "-123.45", a)
	}
}

func TestRow_ToString_String(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonString("abc")}
	a, err := v.ToString("a")
	assert.Equal(t, err, nil)
	assert.Equal(t, a, "abc")
}

func TestRow_ToString_True(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonBoolean(true)}
	a, err := v.ToString("a")
	assert.Equal(t, err, nil)
	assert.Equal(t, a, "true")
}

func TestRow_ToString_False(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonBoolean(false)}
	a, err := v.ToString("a")
	assert.Equal(t, err, nil)
	assert.Equal(t, a, "false")
}

func TestRow_ToString_Object(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonObjectEmpty()}
	_, err := v.ToString("a")
	assert.NotEqual(t, err, nil)
}

func TestRow_ToString_Array(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonArrayEmpty()}
	_, err := v.ToString("a")
	assert.NotEqual(t, err, nil)
}

func TestRow_ToBool_Null(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonNull()}
	a, err := v.ToBool("a")
	assert.Equal(t, err, nil)
	assert.Equal(t, a, false)
}

func TestRow_ToBool_Int64(t *testing.T) {
	t.Run("parse 123 as true", func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonNumberInt64(123)}
		a, err := v.ToBool("a")
		assert.Equal(t, err, nil)
		assert.Equal(t, a, true)
	})
	t.Run("parse 1 as true", func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonNumberInt64(1)}
		a, err := v.ToBool("a")
		assert.Equal(t, err, nil)
		assert.Equal(t, a, true)
	})
	t.Run("parse 0 as false", func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonNumberInt64(0)}
		a, err := v.ToBool("a")
		assert.Equal(t, err, nil)
		assert.Equal(t, a, false)
	})
}

func TestRow_ToBool_Float64(t *testing.T) {
	t.Run("parse -123.45 as true", func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonNumberFloat64(-123.45)}
		a, err := v.ToBool("a")
		assert.Equal(t, err, nil)
		assert.Equal(t, a, true)
	})
	t.Run("parse 1 as true", func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonNumberFloat64(1)}
		a, err := v.ToBool("a")
		assert.Equal(t, err, nil)
		assert.Equal(t, a, true)
	})
	t.Run("parse 0 as false", func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonNumberFloat64(0)}
		a, err := v.ToBool("a")
		assert.Equal(t, err, nil)
		assert.Equal(t, a, false)
	})
}

func TestRow_ToBool_String(t *testing.T) {
	t.Run(`cannot parse "abc"`, func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonString("abc")}
		_, err := v.ToBool("a")
		assert.NotEqual(t, err, nil)
	})
	t.Run(`cannot parse "null"`, func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonString("null")}
		_, err := v.ToBool("a")
		assert.NotEqual(t, err, nil)
	})
	t.Run(`cannot parse ""`, func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonString("")}
		_, err := v.ToBool("a")
		assert.NotEqual(t, err, nil)
	})
	t.Run(`parse "true" as true`, func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonString("true")}
		a, err := v.ToBool("a")
		assert.Equal(t, err, nil)
		assert.Equal(t, a, true)
	})
	t.Run(`parse "false" as false`, func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonString("false")}
		a, err := v.ToBool("a")
		assert.Equal(t, err, nil)
		assert.Equal(t, a, false)
	})
}

func TestRow_ToBool_True(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonBoolean(true)}
	a, err := v.ToBool("a")
	assert.Equal(t, err, nil)
	assert.Equal(t, a, true)
}

func TestRow_ToBool_False(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonBoolean(false)}
	a, err := v.ToBool("a")
	assert.Equal(t, err, nil)
	assert.Equal(t, a, false)
}

func TestRow_ToBool_Object(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonObjectEmpty()}
	_, err := v.ToBool("a")
	assert.NotEqual(t, err, nil)
}

func TestRow_ToBool_Array(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonArrayEmpty()}
	_, err := v.ToBool("a")
	assert.NotEqual(t, err, nil)
}

func TestRow_ToInt64_Null(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonNull()}
	a, err := v.ToInt64("a")
	assert.Equal(t, err, nil)
	assert.Equal(t, a, int64(0))
}

func TestRow_ToInt64_Int64(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonNumberInt64(123)}
	a, err := v.ToInt64("a")
	assert.Equal(t, err, nil)
	assert.Equal(t, a, int64(123))
}

func TestRow_ToInt64_Float64(t *testing.T) {
	t.Run("cannot parse as integer", func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonNumberFloat64(-123.45)}
		_, err := v.ToInt64("a")
		assert.NotEqual(t, err, nil)
	})
	t.Run("integer", func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonNumberFloat64(123)}
		a, err := v.ToInt64("a")
		assert.Equal(t, err, nil)
		assert.Equal(t, a, int64(123))
	})
}

func TestRow_ToInt64_String(t *testing.T) {
	t.Run("cannot parse as integer", func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonString("abc")}
		_, err := v.ToInt64("a")
		assert.NotEqual(t, err, nil)
	})
	t.Run("parse as integer", func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonString("123")}
		a, err := v.ToInt64("a")
		assert.Equal(t, err, nil)
		assert.Equal(t, a, int64(123))
	})
}

func TestRow_ToInt64_Boolean(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonBoolean(true)}
		a, err := v.ToInt64("a")
		assert.Equal(t, err, nil)
		assert.Equal(t, a, int64(1))
	})
	t.Run("false", func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonBoolean(false)}
		a, err := v.ToInt64("a")
		assert.Equal(t, err, nil)
		assert.Equal(t, a, int64(0))
	})
}

func TestRow_ToInt64_Object(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonObjectEmpty()}
	_, err := v.ToInt64("a")
	assert.NotEqual(t, err, nil)
}

func TestRow_ToInt64_Array(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonArrayEmpty()}
	_, err := v.ToInt64("a")
	assert.NotEqual(t, err, nil)
}

func TestRow_ToFloat64_Null(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonNull()}
	a, err := v.ToFloat64("a")
	assert.Equal(t, err, nil)
	assert.Equal(t, a, float64(0))
}

func TestRow_ToFloat64_Int64(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonNumberInt64(123)}
	a, err := v.ToFloat64("a")
	assert.Equal(t, err, nil)
	assert.Equal(t, a, float64(123.0))
}

func TestRow_ToFloat64_Float64(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonNumberFloat64(-123.45)}
	a, err := v.ToFloat64("a")
	assert.Equal(t, err, nil)
	assert.Equal(t, a, float64(-123.45))
}

func TestRow_ToFloat64_String(t *testing.T) {
	t.Run("cannot parse as float", func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonString("abc")}
		_, err := v.ToFloat64("a")
		assert.NotEqual(t, err, nil)
	})
	t.Run("parse as float", func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonString("-123.45")}
		a, err := v.ToFloat64("a")
		assert.Equal(t, err, nil)
		assert.Equal(t, a, float64(-123.45))
	})
}

func TestRow_ToFloat64_Boolean(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonBoolean(true)}
		a, err := v.ToFloat64("a")
		assert.Equal(t, err, nil)
		assert.Equal(t, a, float64(1.0))
	})
	t.Run("false", func(t *testing.T) {
		v := jsonio.Row{"a": jsonio.NewJsonBoolean(false)}
		a, err := v.ToFloat64("a")
		assert.Equal(t, err, nil)
		assert.Equal(t, a, float64(0.0))
	})
}

func TestRow_ToFloat64_Object(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonObjectEmpty()}
	_, err := v.ToFloat64("a")
	assert.NotEqual(t, err, nil)
}

func TestRow_ToFloat64_Array(t *testing.T) {
	v := jsonio.Row{"a": jsonio.NewJsonArrayEmpty()}
	_, err := v.ToFloat64("a")
	assert.NotEqual(t, err, nil)
}

func TestRow_SetString(t *testing.T) {
	v := jsonio.Row{}
	v.SetString("a", "abc")
	a, _ := v.ToString("a")
	assert.Equal(t, a, "abc")
}

func TestRow_SetInt64(t *testing.T) {
	v := jsonio.Row{}
	v.SetInt64("a", int64(123))
	a, _ := v.ToInt64("a")
	assert.Equal(t, a, int64(123))
}

func TestRow_SetFloat64(t *testing.T) {
	v := jsonio.Row{}
	v.SetFloat64("a", float64(-123.45))
	a, _ := v.ToFloat64("a")
	assert.Equal(t, a, float64(-123.45))
}

func TestRow_SetBool(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		v := jsonio.Row{}
		v.SetBool("a", true)
		a, _ := v.ToBool("a")
		assert.Equal(t, a, true)
	})
	t.Run("false", func(t *testing.T) {
		v := jsonio.Row{}
		v.SetBool("a", false)
		a, _ := v.ToBool("a")
		assert.Equal(t, a, false)
	})
}

func TestRow_SetNil(t *testing.T) {
	v := jsonio.Row{}
	v.SetNil("a")
	a, _ := v.GetJsonType("a")
	assert.Equal(t, a, jsonio.JsonTypeNull)
}
