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
	a := v.MustString()
	assert.Equal(t, a, s)
}

func TestBoolean_True(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		b := true
		v := wrap.Boolean(b)
		assert.Equal(t, v.Type, wrap.JsonTypeBoolean)
		assert.Equal(t, v.MustBool(), b)
	})
	t.Run("false", func(t *testing.T) {
		b := false
		v := wrap.Boolean(b)
		assert.Equal(t, v.Type, wrap.JsonTypeBoolean)
		assert.Equal(t, v.MustBool(), b)
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
			a, ok = v.MustNumber().Int64()
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
			a, ok = v.MustNumber().Float64()
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
		a, ok = v.MustNumber().Int64()
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
		a, ok = v.MustNumber().Float64()
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
		assert.Equal(t, len(v.MustObject()), 0)
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
		assert.Equal(t, len(v.MustObject()), 7)
	})
}

func TestArray(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := wrap.Array()
		assert.Equal(t, v.Type, wrap.JsonTypeArray)
		assert.Equal(t, len(v.MustArray()), 0)
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
		assert.Equal(t, len(v.MustArray()), 7)
	})
}
