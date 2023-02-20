package wrap_test

import (
	"encoding/json"
	"testing"

	jw "github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
	"github.com/Jumpaku/api-regression-detector/test/assert"
)

func TestNumberInt64(t *testing.T) {
	v := jw.JsonNumber(json.Number("123"))
	a, err := json.Number(v).Int64()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, int64(123))
}

func TestNumberFloat64(t *testing.T) {
	v := jw.JsonNumber(json.Number("-123.45"))
	a, err := json.Number(v).Float64()
	assert.Equal(t, err, nil)
	assert.Equal(t, a, -123.45)
}
