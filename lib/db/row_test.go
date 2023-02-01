package db

import (
	"testing"
	"time"

	"github.com/Jumpaku/api-regression-detector/test/assert"
)

func TestRow_String(t *testing.T) {
	v := Row{}
	v.SetColumnValue("a", "abc", ColumnTypeString)
	a, ok := v.GetColumnValue("a")
	assert.Equal(t, ok, true)
	assert.Equal(t, a.Type, ColumnTypeString)
	aV, _ := a.AsString()
	assert.Equal(t, aV.String, "abc")
}
func TestRow_True(t *testing.T) {
	v := Row{}
	v.SetColumnValue("a", true, ColumnTypeBoolean)
	a, ok := v.GetColumnValue("a")
	assert.Equal(t, ok, true)
	assert.Equal(t, a.Type, ColumnTypeBoolean)
	aV, _ := a.AsBool()
	assert.Equal(t, aV.Bool, true)
}
func TestRow_False(t *testing.T) {
	v := Row{}
	v.SetColumnValue("a", false, ColumnTypeBoolean)
	a, ok := v.GetColumnValue("a")
	assert.Equal(t, ok, true)
	assert.Equal(t, a.Type, ColumnTypeBoolean)
	aV, _ := a.AsBool()
	assert.Equal(t, aV.Bool, false)
}
func TestRow_Nil(t *testing.T) {
	v := Row{}
	v.SetColumnValue("a", nil, ColumnTypeUnknown)
	a, ok := v.GetColumnValue("a")
	assert.Equal(t, ok, true)
	assert.Equal(t, a.Type, ColumnTypeUnknown)
	aV, _ := a.AsString()
	assert.Equal(t, aV.Valid, false)
}
func TestRow_Integer(t *testing.T) {
	v := Row{}
	v.SetColumnValue("a", 123, ColumnTypeInteger)
	a, ok := v.GetColumnValue("a")
	assert.Equal(t, ok, true)
	assert.Equal(t, a.Type, ColumnTypeInteger)
	aV, _ := a.AsInteger()
	assert.Equal(t, aV.Int64, int64(123))
}
func TestRow_Float(t *testing.T) {
	v := Row{}
	v.SetColumnValue("a", -123.45, ColumnTypeFloat)
	a, ok := v.GetColumnValue("a")
	assert.Equal(t, ok, true)
	assert.Equal(t, a.Type, ColumnTypeFloat)
	aV, _ := a.AsFloat()
	assert.Equal(t, aV.Float64, float64(-123.45))
}
func TestRow_Time(t *testing.T) {
	e := time.Date(2023, time.February, 1, 12, 34, 56, 0, time.UTC)
	v := Row{}
	v.SetColumnValue("a", e, ColumnTypeFloat)
	a, ok := v.GetColumnValue("a")
	assert.Equal(t, ok, true)
	assert.Equal(t, a.Type, ColumnTypeFloat)
	aV, _ := a.AsTime()
	assert.Equal(t, aV.Time, e)
}

/*
func (row Row) GetColumnValue(columnName string) (*ColumnValue, error) {
	val, exists := row[columnName]
	if !exists {
		return nil, fmt.Errorf("column %s not found", columnName)
	}
	return val, nil
}
func (row Row) SetColumnValue(columnName string, val any, typ ColumnType) {
	row[columnName] = &ColumnValue{Type: typ, value: val}
}
*/
