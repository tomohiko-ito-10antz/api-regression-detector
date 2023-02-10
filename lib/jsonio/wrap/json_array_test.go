package wrap_test

import (
	"testing"

	jw "github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
	"github.com/Jumpaku/api-regression-detector/test/assert"
)

func TestArrayLen(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := jw.JsonArray(nil)
		assert.Equal(t, v.Len(), 0)
	})
	t.Run("elements", func(t *testing.T) {
		v := jw.JsonArray([]*jw.JsonValue{
			jw.Number(123),
			jw.Number(-123.45),
			jw.String("abc"),
			nil,
			jw.Boolean(true),
			jw.Boolean(false),
			jw.Object(map[string]*jw.JsonValue{}),
			jw.Array(),
		})
		assert.Equal(t, v.Len(), 8)
	})
}

func TestArrayEmpty(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		v := jw.JsonArray(nil)
		assert.Equal(t, v.Empty(), true)
	})
	t.Run("elements", func(t *testing.T) {
		v := jw.JsonArray([]*jw.JsonValue{
			jw.Number(123),
			jw.Number(-123.45),
			jw.String("abc"),
			nil,
			jw.Boolean(true),
			jw.Boolean(false),
			jw.Object(map[string]*jw.JsonValue{}),
			jw.Array(),
		})
		assert.Equal(t, v.Empty(), false)
	})
}

func TestArrayGet(t *testing.T) {
	v := jw.JsonArray([]*jw.JsonValue{
		jw.Object(map[string]*jw.JsonValue{
			"a": jw.Number(123),
			"b": jw.Number(-123.45),
			"c": jw.String("abc"),
			"d": nil,
			"e": jw.Boolean(true),
			"f": jw.Boolean(false),
			"g": jw.Object(nil),
			"h": jw.Array(),
		}),
		jw.Array(
			jw.Number(123),
			jw.Number(-123.45),
			jw.String("abc"),
			nil,
			jw.Boolean(true),
			jw.Boolean(false),
			jw.Object(nil),
			jw.Array(),
		),
	})

	assert.Equal(t, v.Get(0).Type, jw.JsonTypeObject)
	assert.Equal(t, v.Get(0).Object().Get("a").Type, jw.JsonTypeNumber)
	assert.Equal(t, v.Get(0).Object().Get("b").Type, jw.JsonTypeNumber)
	assert.Equal(t, v.Get(0).Object().Get("c").Type, jw.JsonTypeString)
	assert.Equal(t, v.Get(0).Object().Get("d").Type, jw.JsonTypeNull)
	assert.Equal(t, v.Get(0).Object().Get("e").Type, jw.JsonTypeBoolean)
	assert.Equal(t, v.Get(0).Object().Get("f").Type, jw.JsonTypeBoolean)
	assert.Equal(t, v.Get(0).Object().Get("g").Type, jw.JsonTypeObject)
	assert.Equal(t, v.Get(0).Object().Get("h").Type, jw.JsonTypeArray)

	assert.Equal(t, v.Get(1).Type, jw.JsonTypeArray)
	assert.Equal(t, v.Get(1).Array().Get(0).Type, jw.JsonTypeNumber)
	assert.Equal(t, v.Get(1).Array().Get(1).Type, jw.JsonTypeNumber)
	assert.Equal(t, v.Get(1).Array().Get(2).Type, jw.JsonTypeString)
	assert.Equal(t, v.Get(1).Array().Get(3).Type, jw.JsonTypeNull)
	assert.Equal(t, v.Get(1).Array().Get(4).Type, jw.JsonTypeBoolean)
	assert.Equal(t, v.Get(1).Array().Get(5).Type, jw.JsonTypeBoolean)
	assert.Equal(t, v.Get(1).Array().Get(6).Type, jw.JsonTypeObject)
	assert.Equal(t, v.Get(1).Array().Get(7).Type, jw.JsonTypeArray)
}

func TestArrayAppend(t *testing.T) {
	arr := jw.JsonArray([]*jw.JsonValue{})
	arr = arr.Append(jw.Object(map[string]*jw.JsonValue{
		"a": jw.Number(123),
		"b": jw.Number(-123.45),
		"c": jw.String("abc"),
		"d": nil,
		"e": jw.Boolean(true),
		"f": jw.Boolean(false),
		"g": jw.Object(nil),
		"h": jw.Array(),
	}))
	arr1 := jw.JsonArray(nil).
		Append(jw.Number(123)).
		Append(jw.Number(-123.45)).
		Append(jw.String("abc")).
		Append(nil).
		Append(jw.Boolean(true)).
		Append(jw.Boolean(false)).
		Append(jw.Object(nil)).
		Append(jw.Array())
	arr = arr.Append(jw.Array(arr1...))

	assert.Equal(t, arr.Get(0).Type, jw.JsonTypeObject)
	assert.Equal(t, arr.Get(0).Object().Get("a").Type, jw.JsonTypeNumber)
	assert.Equal(t, arr.Get(0).Object().Get("b").Type, jw.JsonTypeNumber)
	assert.Equal(t, arr.Get(0).Object().Get("c").Type, jw.JsonTypeString)
	assert.Equal(t, arr.Get(0).Object().Get("d").Type, jw.JsonTypeNull)
	assert.Equal(t, arr.Get(0).Object().Get("e").Type, jw.JsonTypeBoolean)
	assert.Equal(t, arr.Get(0).Object().Get("f").Type, jw.JsonTypeBoolean)
	assert.Equal(t, arr.Get(0).Object().Get("g").Type, jw.JsonTypeObject)
	assert.Equal(t, arr.Get(0).Object().Get("h").Type, jw.JsonTypeArray)

	assert.Equal(t, arr.Get(1).Type, jw.JsonTypeArray)
	assert.Equal(t, arr.Get(1).Array().Get(0).Type, jw.JsonTypeNumber)
	assert.Equal(t, arr.Get(1).Array().Get(1).Type, jw.JsonTypeNumber)
	assert.Equal(t, arr.Get(1).Array().Get(2).Type, jw.JsonTypeString)
	assert.Equal(t, arr.Get(1).Array().Get(3).Type, jw.JsonTypeNull)
	assert.Equal(t, arr.Get(1).Array().Get(4).Type, jw.JsonTypeBoolean)
	assert.Equal(t, arr.Get(1).Array().Get(5).Type, jw.JsonTypeBoolean)
	assert.Equal(t, arr.Get(1).Array().Get(6).Type, jw.JsonTypeObject)
	assert.Equal(t, arr.Get(1).Array().Get(7).Type, jw.JsonTypeArray)
}

func TestArraSet(t *testing.T) {
	arr := jw.JsonArray([]*jw.JsonValue{nil, nil, nil, nil, nil, nil, nil, nil, nil})

	arr.Set(0, jw.Number(123))
	arr.Set(1, jw.Number(-123.45))
	arr.Set(2, jw.String("abc"))
	arr.Set(3, nil)
	arr.Set(4, jw.Boolean(true))
	arr.Set(5, jw.Boolean(false))
	arr.Set(6, jw.Object(nil))
	arr.Set(7, jw.Array())

	assert.Equal(t, arr.Get(0).Type, jw.JsonTypeNumber)
	assert.Equal(t, arr.Get(1).Type, jw.JsonTypeNumber)
	assert.Equal(t, arr.Get(2).Type, jw.JsonTypeString)
	assert.Equal(t, arr.Get(3).Type, jw.JsonTypeNull)
	assert.Equal(t, arr.Get(4).Type, jw.JsonTypeBoolean)
	assert.Equal(t, arr.Get(5).Type, jw.JsonTypeBoolean)
	assert.Equal(t, arr.Get(6).Type, jw.JsonTypeObject)
	assert.Equal(t, arr.Get(7).Type, jw.JsonTypeArray)
}
