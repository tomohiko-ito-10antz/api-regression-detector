package tables_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/mock"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/tables"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
	"github.com/Jumpaku/api-regression-detector/test/assert"
	"golang.org/x/exp/slices"
)

func TestInitTables_GetTableNames(t *testing.T) {
	v := tables.InitTables{
		tables.Table{Name: "a"},
		tables.Table{Name: "z"},
		tables.Table{Name: "a"},
		tables.Table{Name: "y"},
	}
	a := v.GetTableNames()
	assert.Equal(t, len(a), 3)
	assert.Equal(t, slices.Contains(a, "a"), true)
	assert.Equal(t, slices.Contains(a, "z"), true)
	assert.Equal(t, slices.Contains(a, "y"), true)
}

func TestDumpTables_GetTableNames(t *testing.T) {
	v := tables.DumpTables{
		"a": []tables.Row{},
		"z": []tables.Row{},
		"b": []tables.Row{},
		"y": []tables.Row{},
	}
	a := v.GetTableNames()
	assert.Equal(t, len(a), 4)
	assert.Equal(t, slices.Contains(a, "a"), true)
	assert.Equal(t, slices.Contains(a, "z"), true)
	assert.Equal(t, slices.Contains(a, "b"), true)
	assert.Equal(t, slices.Contains(a, "y"), true)
}

func TestLoadInitTables(t *testing.T) {
	v := `
[
	{
		"name": "t1",
		"rows": [
			{
				"a": null,
				"b": 123,
				"c": -123.45,
				"d": "abc",
				"e": false,
				"f": true
			},
			{
				"a": null,
				"b": 123,
				"c": -123.45,
				"d": "abc",
				"e": false,
				"f": true
			}
		]
	},
	{
		"name": "t2",
		"rows": [
			{
				"x": null,
				"y": 123,
				"z": -123.45
			},
			{
				"x": null,
				"y": 123,
				"z": -123.45
			}
		]
	},
	{
		"name": "t1",
		"rows": [
			{
				"b": 123
			},
			{
				"d": "abc"
			}
		]
	}
]`
	reader := mock.NamedBuffer{Buffer: bytes.NewBufferString(v)}
	a, err := tables.LoadInitTables(reader)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(a), 3)

	a0 := a[0]
	assert.Equal(t, a0.Name, "t1")
	assert.Equal(t, len(a0.Rows), 2)

	a00 := a0.Rows[0]
	assert.Equal(t, len(a00), 6)
	assert.Equal(t, a00["a"].Type, wrap.JsonTypeNull)
	assert.Equal(t, a00["b"].Type, wrap.JsonTypeNumber)
	assert.Equal(t, a00["c"].Type, wrap.JsonTypeNumber)
	assert.Equal(t, a00["d"].Type, wrap.JsonTypeString)
	assert.Equal(t, a00["e"].Type, wrap.JsonTypeBoolean)
	assert.Equal(t, a00["f"].Type, wrap.JsonTypeBoolean)

	a01 := a0.Rows[1]
	assert.Equal(t, len(a01), 6)
	assert.Equal(t, a01["a"].Type, wrap.JsonTypeNull)
	assert.Equal(t, a01["b"].Type, wrap.JsonTypeNumber)
	assert.Equal(t, a01["c"].Type, wrap.JsonTypeNumber)
	assert.Equal(t, a01["d"].Type, wrap.JsonTypeString)
	assert.Equal(t, a01["e"].Type, wrap.JsonTypeBoolean)
	assert.Equal(t, a01["f"].Type, wrap.JsonTypeBoolean)

	a1 := a[1]
	assert.Equal(t, a1.Name, "t2")
	assert.Equal(t, len(a1.Rows), 2)

	a10 := a1.Rows[0]
	assert.Equal(t, len(a10), 3)
	assert.Equal(t, a10["x"].Type, wrap.JsonTypeNull)
	assert.Equal(t, a10["y"].Type, wrap.JsonTypeNumber)
	assert.Equal(t, a10["z"].Type, wrap.JsonTypeNumber)

	a11 := a1.Rows[1]
	assert.Equal(t, len(a11), 3)
	assert.Equal(t, a11["x"].Type, wrap.JsonTypeNull)
	assert.Equal(t, a11["y"].Type, wrap.JsonTypeNumber)
	assert.Equal(t, a11["z"].Type, wrap.JsonTypeNumber)

	a2 := a[2]
	assert.Equal(t, a0.Name, "t1")
	assert.Equal(t, len(a1.Rows), 2)

	a20 := a2.Rows[0]
	assert.Equal(t, len(a20), 1)
	assert.Equal(t, a20["b"].Type, wrap.JsonTypeNumber)

	a21 := a2.Rows[1]
	assert.Equal(t, len(a21), 1)
	assert.Equal(t, a21["d"].Type, wrap.JsonTypeString)
}

func TestSaveDumpTables(t *testing.T) {
	v := tables.DumpTables{
		"t1": {
			tables.Row{
				"a": nil,
				"b": wrap.Number(123),
				"c": wrap.Number(-123.45),
				"d": wrap.String("abc"),
				"e": wrap.Boolean(false),
				"f": wrap.Boolean(true),
			},
			tables.Row{
				"a": nil,
				"b": wrap.Number(123),
				"c": wrap.Number(-123.45),
				"d": wrap.String("abc"),
				"e": wrap.Boolean(false),
				"f": wrap.Boolean(true),
			},
		},
		"t2": {
			tables.Row{
				"x": nil,
				"y": wrap.Number(json.Number("123")),
				"z": wrap.Number(json.Number("-123.45")),
			},
			tables.Row{
				"x": nil,
				"y": wrap.Number(json.Number("123")),
				"z": wrap.Number(json.Number("-123.45")),
			},
		},
	}
	buffer := mock.NamedBuffer{Buffer: bytes.NewBuffer(nil)}
	_ = tables.SaveDumpTables(v, buffer)

	a, err := jsonio.LoadJson[map[string][]map[string]any](mock.NamedBuffer{Buffer: buffer.Buffer})
	assert.Equal(t, err, nil)

	aT1, ok := a["t1"]
	assert.Equal(t, ok, true)
	assert.Equal(t, len(aT1), 2)

	aT10 := aT1[0]
	assert.Equal(t, len(aT10), 6)
	assert.Equal(t, aT10["a"], nil)
	assert.Equal(t, aT10["b"], json.Number("123"))
	assert.Equal(t, aT10["c"], json.Number("-123.45"))
	assert.Equal(t, aT10["d"], "abc")
	assert.Equal(t, aT10["e"], false)
	assert.Equal(t, aT10["f"], true)

	aT11 := aT1[1]
	assert.Equal(t, len(aT11), 6)
	assert.Equal(t, aT11["a"], nil)
	assert.Equal(t, aT11["b"], json.Number("123"))
	assert.Equal(t, aT11["c"], json.Number("-123.45"))
	assert.Equal(t, aT11["d"], "abc")
	assert.Equal(t, aT11["e"], false)
	assert.Equal(t, aT11["f"], true)

	aT2, ok := a["t2"]
	assert.Equal(t, ok, true)
	assert.Equal(t, len(aT2), 2)

	aT20 := aT2[0]
	assert.Equal(t, len(aT20), 3)
	assert.Equal(t, aT20["x"], nil)
	assert.Equal(t, aT20["y"], json.Number("123"))
	assert.Equal(t, aT20["z"], json.Number("-123.45"))

	aT21 := aT2[1]
	assert.Equal(t, len(aT21), 3)
	assert.Equal(t, aT21["x"], nil)
	assert.Equal(t, aT21["y"], json.Number("123"))
	assert.Equal(t, aT20["z"], json.Number("-123.45"))
}
