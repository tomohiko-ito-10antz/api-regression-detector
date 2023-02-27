package cmd_test

import (
	"context"
	"testing"

	"github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/cmd/mock"
	"github.com/Jumpaku/api-regression-detector/test/assert"
)

func TestDump_OK(t *testing.T) {
	v, err := cmd.Dump(context.Background(),
		mock.DB{},
		[]string{"mock_table"},
		mock.SchemaGetter{},
		mock.RowLister{})
	assert.Equal(t, err, nil)
	assert.Equal(t, len(v), 1)
	assert.Equal(t, len(v["mock_table"]), 3)
}

func TestDump_NG_Table(t *testing.T) {
	_, err := cmd.Dump(context.Background(),
		mock.DB{},
		[]string{"invalid_table"},
		mock.SchemaGetter{},
		mock.RowLister{})
	assert.NotEqual(t, err, nil)
}

func TestDump_NG_DB(t *testing.T) {
	_, err := cmd.Dump(context.Background(),
		mock.ErrDB{},
		[]string{"mock_table"},
		mock.SchemaGetter{},
		mock.RowLister{})
	assert.NotEqual(t, err, nil)
}

func TestDump_NG_SchemaGetter(t *testing.T) {
	_, err := cmd.Dump(context.Background(),
		mock.DB{},
		[]string{"mock_table"},
		mock.ErrSchemaGetter{},
		mock.RowLister{})
	assert.NotEqual(t, err, nil)
}

func TestDump_NG_RowLister(t *testing.T) {
	_, err := cmd.Dump(context.Background(),
		mock.DB{},
		[]string{"mock_table"},
		mock.SchemaGetter{},
		mock.ErrRowLister{})
	assert.NotEqual(t, err, nil)
}
