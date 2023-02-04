package db_test

import (
	"testing"

	"github.com/Jumpaku/api-regression-detector/lib/db"
	"github.com/Jumpaku/api-regression-detector/test/assert"
	"golang.org/x/exp/slices"
)

func TestSchema_GetColumnNames(t *testing.T) {
	v := db.Schema{
		PrimaryKeys: []string{"a", "b", "c"},
		ColumnTypes: db.ColumnTypes{
			"a": db.ColumnTypeBoolean,
			"b": db.ColumnTypeFloat,
			"c": db.ColumnTypeInteger,
			"x": db.ColumnTypeString,
			"y": db.ColumnTypeTime,
			"z": db.ColumnTypeUnknown,
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
