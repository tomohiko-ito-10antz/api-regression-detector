package io

import (
	"testing"

	"github.com/Jumpaku/api-regression-detector/test/assert"
	"golang.org/x/exp/slices"
)

func TestGetTableNames(t *testing.T) {
	v := Tables{
		"a": Table{},
		"z": Table{},
		"b": Table{},
		"y": Table{},
	}
	a := v.GetTableNames()
	assert.Equal(t, len(a), 4)
	assert.Equal(t, slices.Contains(a, "a"), true)
	assert.Equal(t, slices.Contains(a, "z"), true)
	assert.Equal(t, slices.Contains(a, "b"), true)
	assert.Equal(t, slices.Contains(a, "y"), true)

}
