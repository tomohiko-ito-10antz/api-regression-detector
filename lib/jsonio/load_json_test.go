package jsonio_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/mock"
	"github.com/Jumpaku/api-regression-detector/test/assert"
)

func TestLoadJson_Tables(t *testing.T) {
	v := `{
    "t1": [
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
	],
	"t2": [
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
}`
	reader := mock.NamedBuffer{Buffer: bytes.NewBufferString(v)}
	a, err := jsonio.LoadJson[map[string][]map[string]any](reader)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(a), 2)

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
	assert.Equal(t, aT10["b"], json.Number("123"))
	assert.Equal(t, aT10["c"], json.Number("-123.45"))
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
	assert.Equal(t, aT20["y"], json.Number("123"))
	assert.Equal(t, aT20["z"], json.Number("-123.45"))
}

func TestLoadJson_TableNames(t *testing.T) {
	v := `["t1", "t2"]`
	reader := mock.NamedBuffer{Buffer: bytes.NewBufferString(v)}
	a, err := jsonio.LoadJson[[]any](reader)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(a), 2)

	assert.Equal(t, a[0], "t1")
	assert.Equal(t, a[1], "t2")
}

func TestLoadJson_NG(t *testing.T) {
	reader := mock.ErrNamedBuffer{}
	_, err := jsonio.LoadJson[[]any](reader)
	assert.NotEqual(t, err, nil)
}
