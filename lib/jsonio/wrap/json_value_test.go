package wrap_test

import (
	"encoding/json"
	"testing"

	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
	"github.com/Jumpaku/api-regression-detector/test/assert"
)

func TestString(t *testing.T) {
	s := "abc"
	v := wrap.String(s)
	assert.Equal(t, v.Type, wrap.JsonTypeString)
	a := v.String()
	assert.Equal(t, a, s)
}

func TestBoolean_True(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		b := true
		v := wrap.Boolean(b)
		assert.Equal(t, v.Type, wrap.JsonTypeBoolean)
		assert.Equal(t, v.Bool(), b)
	})
	t.Run("false", func(t *testing.T) {
		b := false
		v := wrap.Boolean(b)
		assert.Equal(t, v.Type, wrap.JsonTypeBoolean)
		assert.Equal(t, v.Bool(), b)
	})
}

func TestNumber_JsonNumber(t *testing.T) {
	t.Run("json.Number", func(t *testing.T) {
		t.Run("int64", func(t *testing.T) {
			n := json.Number("123")
			v := wrap.Number(n)
			assert.Equal(t, v.Type, wrap.JsonTypeNumber)
			a, ok := v.Int64()
			assert.Equal(t, ok, true)
			assert.Equal(t, a, int64(123))
			a, ok = v.Number().Int64()
			assert.Equal(t, ok, true)
			assert.Equal(t, a, int64(123))
		})
		t.Run("float64", func(t *testing.T) {
			n := json.Number("-123.45")
			v := wrap.Number(n)
			assert.Equal(t, v.Type, wrap.JsonTypeNumber)
			a, ok := v.Float64()
			assert.Equal(t, ok, true)
			assert.Equal(t, a, -123.45)
			a, ok = v.Number().Float64()
			assert.Equal(t, ok, true)
			assert.Equal(t, a, -123.45)
		})
	})
	t.Run("int64", func(t *testing.T) {
		i := int64(123)
		v := wrap.Number(i)
		assert.Equal(t, v.Type, wrap.JsonTypeNumber)
		a, ok := v.Int64()
		assert.Equal(t, ok, true)
		assert.Equal(t, a, i)
		a, ok = v.Number().Int64()
		assert.Equal(t, ok, true)
		assert.Equal(t, a, i)
	})
	t.Run("float64", func(t *testing.T) {
		f := -123.45
		v := wrap.Number(f)
		assert.Equal(t, v.Type, wrap.JsonTypeNumber)
		a, ok := v.Float64()
		assert.Equal(t, ok, true)
		assert.Equal(t, a, f)
		a, ok = v.Number().Float64()
		assert.Equal(t, ok, true)
		assert.Equal(t, a, f)
	})
}

func TestNull(t *testing.T) {
	v := wrap.Null()
	assert.NotEqual(t, v, nil)
	assert.Equal(t, v.Type, wrap.JsonTypeNull)
}

func TestObject(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := wrap.Object(nil)
		assert.Equal(t, v.Type, wrap.JsonTypeObject)
		assert.Equal(t, len(v.Object()), 0)
	})
	t.Run("elements", func(t *testing.T) {
		v := wrap.Object(map[string]*wrap.JsonValue{
			"a": nil,
			"g": wrap.Null(),
			"b": wrap.String("xyz"),
			"c": wrap.Number(123),
			"d": wrap.Number(-123.45),
			"e": wrap.Boolean(true),
			"f": wrap.Boolean(false),
		})
		assert.Equal(t, v.Type, wrap.JsonTypeObject)
		assert.Equal(t, len(v.Object()), 7)
	})
}

func TestArray(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := wrap.Array()
		assert.Equal(t, v.Type, wrap.JsonTypeArray)
		assert.Equal(t, len(v.Array()), 0)
	})
	t.Run("elements", func(t *testing.T) {
		v := wrap.Array(
			nil,
			wrap.Null(),
			wrap.String("xyz"),
			wrap.Number(123),
			wrap.Number(-123.45),
			wrap.Boolean(true),
			wrap.Boolean(false),
		)
		assert.Equal(t, v.Type, wrap.JsonTypeArray)
		assert.Equal(t, len(v.Array()), 7)
	})
}

func TestFromAny_Nil(t *testing.T) {
	v, err := wrap.FromAny(nil)
	assert.Equal(t, err, nil)
	assert.Equal(t, v.Type, wrap.JsonTypeNull)
}

func TestFromAny_Boolean(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		v, err := wrap.FromAny(true)
		assert.Equal(t, err, nil)
		assert.Equal(t, v.Type, wrap.JsonTypeBoolean)
		assert.Equal(t, v.BooleanValue, true)
	})
	t.Run("false", func(t *testing.T) {
		v, err := wrap.FromAny(false)
		assert.Equal(t, err, nil)
		assert.Equal(t, v.Type, wrap.JsonTypeBoolean)
		assert.Equal(t, v.BooleanValue, false)
	})
}

func TestFromAny_Int64(t *testing.T) {
	i := int64(123)
	v, err := wrap.FromAny(i)
	assert.Equal(t, err, nil)
	assert.Equal(t, v.Type, wrap.JsonTypeNumber)
	a, ok := v.NumberValue.Int64()
	assert.Equal(t, ok, true)
	assert.Equal(t, a, i)
}

func TestFromAny_Float(t *testing.T) {
	f := -123.45
	v, err := wrap.FromAny(f)
	assert.Equal(t, err, nil)
	assert.Equal(t, v.Type, wrap.JsonTypeNumber)
	a, ok := v.NumberValue.Float64()
	assert.Equal(t, ok, true)
	assert.Equal(t, a, f)
}

func TestFromAny_JsonNumber(t *testing.T) {
	t.Run("int64", func(t *testing.T) {
		v, err := wrap.FromAny(json.Number("123"))
		assert.Equal(t, err, nil)
		assert.Equal(t, v.Type, wrap.JsonTypeNumber)
		a, ok := v.NumberValue.Int64()
		assert.Equal(t, ok, true)
		assert.Equal(t, a, int64(123))
	})
	t.Run("float64", func(t *testing.T) {
		v, err := wrap.FromAny(json.Number("-123.45"))
		assert.Equal(t, err, nil)
		assert.Equal(t, v.Type, wrap.JsonTypeNumber)
		a, ok := v.NumberValue.Float64()
		assert.Equal(t, ok, true)
		assert.Equal(t, a, float64(-123.45))
	})
}

func TestFromAny_Object(t *testing.T) {
	v, err := wrap.FromAny(map[string]any{
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
	assert.Equal(t, v.Type, wrap.JsonTypeObject)
	assert.Equal(t, len(v.ObjectValue), 2)
}

func TestFromAny_Array(t *testing.T) {
	v, err := wrap.FromAny([]any{
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
	assert.Equal(t, v.Type, wrap.JsonTypeArray)
	assert.Equal(t, len(v.ArrayValue), 2)
}
