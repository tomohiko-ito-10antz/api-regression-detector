package topological_sort_test

import (
	"testing"

	"github.com/Jumpaku/api-regression-detector/lib/topological_sort"
	"github.com/Jumpaku/api-regression-detector/test/assert"
)

func TestTopologicalSort_OK(t *testing.T) {
	g := topological_sort.NewGraph[int]()
	g.Arrow(1, 2)
	g.Arrow(3, 2)
	g.Arrow(2, 4)
	g.Arrow(2, 5)
	g.Arrow(2, 6)
	g.Arrow(4, 7)
	g.Arrow(5, 7)
	g.Arrow(8, 9)

	a, ok := topological_sort.Perform(g)
	assert.Equal(t, ok, true)
	assert.Equal(t, a[1] < a[2], true)
	assert.Equal(t, a[3] < a[2], true)
	assert.Equal(t, a[2] < a[4], true)
	assert.Equal(t, a[2] < a[5], true)
	assert.Equal(t, a[2] < a[6], true)
	assert.Equal(t, a[4] < a[7], true)
	assert.Equal(t, a[5] < a[7], true)
	assert.Equal(t, a[8] < a[9], true)
}

func TestTopologicalSort_NG(t *testing.T) {
	g := topological_sort.NewGraph[int]()
	g.Arrow(1, 2)
	g.Arrow(2, 1)

	_, ok := topological_sort.Perform(g)
	assert.Equal(t, ok, false)
}
