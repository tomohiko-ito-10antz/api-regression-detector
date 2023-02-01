package db

import (
	"testing"

	"github.com/Jumpaku/api-regression-detector/test/assert"
	"golang.org/x/exp/slices"
)

func TestSchema_GetColumnNames(t *testing.T) {
	v := Schema{
		PrimaryKeys: []string{"a", "b", "c"},
		ColumnTypes: ColumnTypes{
			"a": ColumnTypeBoolean,
			"b": ColumnTypeFloat,
			"c": ColumnTypeInteger,
			"x": ColumnTypeString,
			"y": ColumnTypeTime,
			"z": ColumnTypeUnknown,
		},
	}
	a := v.GetColumnNames()
	assert.Equal(t, len(a), 6)
	assert.Equal(t, slices.Contains(a, "a"), true)
	assert.Equal(t, slices.Contains(a, "b"), true)
	assert.Equal(t, slices.Contains(a, "c"), true)
	assert.Equal(t, slices.Contains(a, "x"), true)
	assert.Equal(t, slices.Contains(a, "y"), true)
	assert.Equal(t, slices.Contains(a, "z"), true)
}
