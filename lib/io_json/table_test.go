package io_json_test

import (
	"encoding/json"
	"testing"

	"github.com/Jumpaku/api-regression-detector/lib/io_json"
	"github.com/Jumpaku/api-regression-detector/test/assert"
	"golang.org/x/exp/slices"
)

func TestGetTableNames(t *testing.T) {
	v := io_json.Tables{
		"a": io_json.Table{},
		"z": io_json.Table{},
		"b": io_json.Table{},
		"y": io_json.Table{},
	}
	a := v.GetTableNames()
	assert.Equal(t, len(a), 4)
	assert.Equal(t, slices.Contains(a, "a"), true)
	assert.Equal(t, slices.Contains(a, "z"), true)
	assert.Equal(t, slices.Contains(a, "b"), true)
	assert.Equal(t, slices.Contains(a, "y"), true)
}

func TestTableFromJson(t *testing.T) {
	v := map[string][]map[string]any{
		"t1": {
			{
				"a": nil,
				"b": int64(123),
				"c": float64(-123.45),
				"d": "abc",
				"e": false,
				"f": true,
			},
			{
				"a": nil,
				"b": int64(123),
				"c": float64(-123.45),
				"d": "abc",
				"e": false,
				"f": true,
			},
		},
		"t2": {
			{
				"x": nil,
				"y": json.Number("123"),
				"z": json.Number("-123.45"),
			},
			{
				"x": nil,
				"y": json.Number("123"),
				"z": json.Number("-123.45"),
			},
		},
	}
	a, err := io_json.TableFromJson(v)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(a), 2)

	aT1, ok := a["t1"]
	assert.Equal(t, ok, true)
	assert.Equal(t, len(aT1.Rows), 2)

	aT10 := aT1.Rows[0]
	assert.Equal(t, len(aT10), 6)
	assert.Equal(t, aT10["a"].Type, io_json.JsonTypeNull)
	assert.Equal(t, aT10["b"].Type, io_json.JsonTypeNumber)
	assert.Equal(t, aT10["c"].Type, io_json.JsonTypeNumber)
	assert.Equal(t, aT10["d"].Type, io_json.JsonTypeString)
	assert.Equal(t, aT10["e"].Type, io_json.JsonTypeBoolean)
	assert.Equal(t, aT10["f"].Type, io_json.JsonTypeBoolean)

	aT11 := aT1.Rows[1]
	assert.Equal(t, len(aT11), 6)
	assert.Equal(t, aT11["a"].Type, io_json.JsonTypeNull)
	assert.Equal(t, aT11["b"].Type, io_json.JsonTypeNumber)
	assert.Equal(t, aT11["c"].Type, io_json.JsonTypeNumber)
	assert.Equal(t, aT11["d"].Type, io_json.JsonTypeString)
	assert.Equal(t, aT11["e"].Type, io_json.JsonTypeBoolean)
	assert.Equal(t, aT11["f"].Type, io_json.JsonTypeBoolean)

	aT2, ok := a["t2"]
	assert.Equal(t, ok, true)
	assert.Equal(t, len(aT2.Rows), 2)

	aT20 := aT2.Rows[0]
	assert.Equal(t, len(aT20), 3)
	assert.Equal(t, aT20["x"].Type, io_json.JsonTypeNull)
	assert.Equal(t, aT20["y"].Type, io_json.JsonTypeNumber)
	assert.Equal(t, aT20["z"].Type, io_json.JsonTypeNumber)

	aT21 := aT2.Rows[1]
	assert.Equal(t, len(aT21), 3)
	assert.Equal(t, aT21["x"].Type, io_json.JsonTypeNull)
	assert.Equal(t, aT21["y"].Type, io_json.JsonTypeNumber)
	assert.Equal(t, aT21["z"].Type, io_json.JsonTypeNumber)
}

func TestTableToJson(t *testing.T) {
	v := map[string][]map[string]any{
		"t1": {
			{
				"a": nil,
				"b": int64(123),
				"c": float64(-123.45),
				"d": "abc",
				"e": false,
				"f": true,
			},
			{
				"a": nil,
				"b": int64(123),
				"c": float64(-123.45),
				"d": "abc",
				"e": false,
				"f": true,
			},
		},
		"t2": {
			{
				"x": nil,
				"y": json.Number("123"),
				"z": json.Number("-123.45"),
			},
			{
				"x": nil,
				"y": json.Number("123"),
				"z": json.Number("-123.45"),
			},
		},
	}
	e, _ := io_json.TableFromJson(v)
	a, err := io_json.TableToJson(e)
	assert.Equal(t, err, nil)

	aT1, ok := a["t1"]
	assert.Equal(t, ok, true)
	assert.Equal(t, len(aT1), 2)

	aT10 := aT1[0]
	assert.Equal(t, len(aT10), 6)
	assert.Equal(t, aT10["a"], nil)
	assert.Equal(t, aT10["b"], int64(123))
	assert.Equal(t, aT10["c"], float64(-123.45))
	assert.Equal(t, aT10["d"], "abc")
	assert.Equal(t, aT10["e"], false)
	assert.Equal(t, aT10["f"], true)

	aT11 := aT1[1]
	assert.Equal(t, len(aT11), 6)
	assert.Equal(t, aT11["a"], nil)
	assert.Equal(t, aT11["b"], int64(123))
	assert.Equal(t, aT11["c"], float64(-123.45))
	assert.Equal(t, aT11["d"], "abc")
	assert.Equal(t, aT11["e"], false)
	assert.Equal(t, aT11["f"], true)

	aT2, ok := a["t2"]
	assert.Equal(t, ok, true)
	assert.Equal(t, len(aT2), 2)

	aT20 := aT2[0]
	assert.Equal(t, len(aT20), 3)
	assert.Equal(t, aT20["x"], nil)
	assert.Equal(t, aT20["y"], int64(123))
	assert.Equal(t, aT20["z"], float64(-123.45))

	aT21 := aT2[1]
	assert.Equal(t, len(aT21), 3)
	assert.Equal(t, aT21["x"], nil)
	assert.Equal(t, aT21["y"], int64(123))
	assert.Equal(t, aT20["z"], float64(-123.45))
}
