package db

import (
	"database/sql"
	"testing"
	"time"

	"github.com/Jumpaku/api-regression-detector/test/assert"
	"golang.org/x/exp/slices"
)

func TestColumnNames(t *testing.T) {
	a := ColumnTypes{
		"a": ColumnTypeBoolean,
		"b": ColumnTypeFloat,
		"c": ColumnTypeInteger,
		"d": ColumnTypeString,
		"e": ColumnTypeTime,
	}.ColumnNames()

	assert.Equal(t, len(a), 5)
	assert.Equal(t, slices.Contains(a, "a"), true)
	assert.Equal(t, slices.Contains(a, "b"), true)
	assert.Equal(t, slices.Contains(a, "c"), true)
	assert.Equal(t, slices.Contains(a, "d"), true)
	assert.Equal(t, slices.Contains(a, "e"), true)
}

func TestNewColumnValue(t *testing.T) {
	a := NewColumnValue(3)
	assert.Equal(t, a.Type, ColumnTypeUnknown)
}
func TestColumnValue_AsString(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		a, err := NewColumnValue(nil).AsString()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run("(*string)(nil)", func(t *testing.T) {
		a, err := NewColumnValue((*string)(nil)).AsString()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run(`pointer to "abc"`, func(t *testing.T) {
		v := "abc"
		a, err := NewColumnValue(&v).AsString()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.String, v)
	})
	t.Run(`invalid sql.NullString`, func(t *testing.T) {
		v := sql.NullString{}
		a, err := NewColumnValue(v).AsString()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run(`valid sql.NullString`, func(t *testing.T) {
		v := sql.NullString{String: "abc", Valid: true}
		a, err := NewColumnValue(v).AsString()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.String, "abc")
	})
	t.Run(`int`, func(t *testing.T) {
		_, err := NewColumnValue(1).AsString()
		assert.NotEqual(t, err, nil)
	})
}

func TestColumnValue_AsInteger(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		a, err := NewColumnValue(nil).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run("(*int)(nil)", func(t *testing.T) {
		a, err := NewColumnValue((*int)(nil)).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run("(*int8)(nil)", func(t *testing.T) {
		a, err := NewColumnValue((*int8)(nil)).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run("(*int16)(nil)", func(t *testing.T) {
		a, err := NewColumnValue((*int16)(nil)).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run("(*int32)(nil)", func(t *testing.T) {
		a, err := NewColumnValue((*int32)(nil)).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run("(*int64)(nil)", func(t *testing.T) {
		a, err := NewColumnValue((*int64)(nil)).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run("(*uint)(nil)", func(t *testing.T) {
		a, err := NewColumnValue((*uint)(nil)).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run("(*uint8)(nil)", func(t *testing.T) {
		a, err := NewColumnValue((*uint8)(nil)).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run("(*uint16)(nil)", func(t *testing.T) {
		a, err := NewColumnValue((*uint16)(nil)).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run("(*uint32)(nil)", func(t *testing.T) {
		a, err := NewColumnValue((*uint32)(nil)).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run("(*uint64)(nil)", func(t *testing.T) {
		a, err := NewColumnValue((*uint64)(nil)).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run(`pointer to int(1)`, func(t *testing.T) {
		v := int(1)
		a, err := NewColumnValue(&v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`pointer to int8(1)`, func(t *testing.T) {
		v := int8(1)
		a, err := NewColumnValue(&v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`pointer to int16(1)`, func(t *testing.T) {
		v := int16(1)
		a, err := NewColumnValue(&v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`pointer to int32(1)`, func(t *testing.T) {
		v := int32(1)
		a, err := NewColumnValue(&v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`pointer to int64(1)`, func(t *testing.T) {
		v := int64(1)
		a, err := NewColumnValue(&v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`pointer to uint(1)`, func(t *testing.T) {
		v := uint(1)
		a, err := NewColumnValue(&v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`pointer to uint8(1)`, func(t *testing.T) {
		v := uint8(1)
		a, err := NewColumnValue(&v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`pointer to uint16(1)`, func(t *testing.T) {
		v := uint16(1)
		a, err := NewColumnValue(&v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`pointer to uint32(1)`, func(t *testing.T) {
		v := uint32(1)
		a, err := NewColumnValue(&v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`pointer to uint64(1)`, func(t *testing.T) {
		v := uint64(1)
		a, err := NewColumnValue(&v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`int(1)`, func(t *testing.T) {
		v := int(1)
		a, err := NewColumnValue(v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`int8(1)`, func(t *testing.T) {
		v := int8(1)
		a, err := NewColumnValue(v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`int16(1)`, func(t *testing.T) {
		v := int16(1)
		a, err := NewColumnValue(v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`int32(1)`, func(t *testing.T) {
		v := int32(1)
		a, err := NewColumnValue(v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`int64(1)`, func(t *testing.T) {
		v := int64(1)
		a, err := NewColumnValue(v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`uint(1)`, func(t *testing.T) {
		v := uint(1)
		a, err := NewColumnValue(v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`uint8(1)`, func(t *testing.T) {
		v := uint8(1)
		a, err := NewColumnValue(v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`uint16(1)`, func(t *testing.T) {
		v := uint16(1)
		a, err := NewColumnValue(v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`uint32(1)`, func(t *testing.T) {
		v := uint32(1)
		a, err := NewColumnValue(v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`uint64(1)`, func(t *testing.T) {
		v := uint64(1)
		a, err := NewColumnValue(v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v))
	})
	t.Run(`invalid sql.NullByte`, func(t *testing.T) {
		v := sql.NullByte{}
		a, err := NewColumnValue(v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run(`invalid sql.NullInt16`, func(t *testing.T) {
		v := sql.NullInt16{}
		a, err := NewColumnValue(v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run(`invalid sql.NullInt32`, func(t *testing.T) {
		v := sql.NullInt32{}
		a, err := NewColumnValue(v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run(`invalid sql.NullInt64`, func(t *testing.T) {
		v := sql.NullInt64{}
		a, err := NewColumnValue(v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run(`valid sql.NullByte`, func(t *testing.T) {
		v := sql.NullByte{Valid: true, Byte: 1}
		a, err := NewColumnValue(v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v.Byte))
	})
	t.Run(`valid sql.NullInt16`, func(t *testing.T) {
		v := sql.NullInt16{Valid: true, Int16: 1}
		a, err := NewColumnValue(v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v.Int16))
	})
	t.Run(`valid sql.NullInt32`, func(t *testing.T) {
		v := sql.NullInt32{Valid: true, Int32: 1}
		a, err := NewColumnValue(v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v.Int32))
	})
	t.Run(`valid sql.NullInt64`, func(t *testing.T) {
		v := sql.NullInt64{Valid: true, Int64: 1}
		a, err := NewColumnValue(v).AsInteger()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Int64, int64(v.Int64))
	})
	t.Run(`string`, func(t *testing.T) {
		_, err := NewColumnValue("abc").AsInteger()
		assert.NotEqual(t, err, nil)
	})
}

func TestColumnValue_AsFloat(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		a, err := NewColumnValue(nil).AsFloat()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run("(*float32)(nil)", func(t *testing.T) {
		a, err := NewColumnValue((*float32)(nil)).AsFloat()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run("(*float64)(nil)", func(t *testing.T) {
		a, err := NewColumnValue((*float64)(nil)).AsFloat()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run(`pointer to float32(1)`, func(t *testing.T) {
		v := float32(1)
		a, err := NewColumnValue(&v).AsFloat()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Float64, float64(v))
	})
	t.Run(`pointer to float64(1)`, func(t *testing.T) {
		v := float64(1)
		a, err := NewColumnValue(&v).AsFloat()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Float64, float64(v))
	})
	t.Run(`float32(1)`, func(t *testing.T) {
		a, err := NewColumnValue(float32(1)).AsFloat()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Float64, float64(1))
	})
	t.Run(`float64(1)`, func(t *testing.T) {
		a, err := NewColumnValue(float64(1)).AsFloat()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Float64, float64(1))
	})
	t.Run(`invalid sql.NullFloat64`, func(t *testing.T) {
		v := sql.NullFloat64{}
		a, err := NewColumnValue(v).AsFloat()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run(`valid sql.NullFloat64`, func(t *testing.T) {
		v := sql.NullFloat64{Valid: true, Float64: 1}
		a, err := NewColumnValue(v).AsFloat()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Float64, float64(1))
	})
	t.Run(`string`, func(t *testing.T) {
		_, err := NewColumnValue("abc").AsFloat()
		assert.NotEqual(t, err, nil)
	})
}

func TestColumnValue_AsBool(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		a, err := NewColumnValue(nil).AsBool()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run("(*bool)(nil)", func(t *testing.T) {
		a, err := NewColumnValue((*bool)(nil)).AsBool()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run(`pointer to true`, func(t *testing.T) {
		v := true
		a, err := NewColumnValue(&v).AsBool()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Bool, v)
	})
	t.Run(`pointer to false`, func(t *testing.T) {
		v := false
		a, err := NewColumnValue(&v).AsBool()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Bool, v)
	})
	t.Run(`invalid sql.NullBool`, func(t *testing.T) {
		v := sql.NullBool{}
		a, err := NewColumnValue(v).AsBool()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run(`valid true sql.NullBool`, func(t *testing.T) {
		v := sql.NullBool{Valid: true, Bool: true}
		a, err := NewColumnValue(v).AsBool()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Bool, true)
	})
	t.Run(`valid false sql.NullBool`, func(t *testing.T) {
		v := sql.NullBool{Valid: true, Bool: false}
		a, err := NewColumnValue(v).AsBool()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Bool, false)
	})
	t.Run(`string`, func(t *testing.T) {
		_, err := NewColumnValue("abc").AsBool()
		assert.NotEqual(t, err, nil)
	})
}

func TestColumnValue_AsBytes(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		a, err := NewColumnValue(nil).AsBytes()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run("(*[]byte)(nil)", func(t *testing.T) {
		a, err := NewColumnValue((*[]byte)(nil)).AsBytes()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run(`pointer to []byte{}`, func(t *testing.T) {
		v := []byte{}
		a, err := NewColumnValue(&v).AsBytes()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.NotEqual(t, a.Bytes, nil)
		assert.Equal(t, string(a.Bytes), "")
	})
	t.Run(`pointer to []byte("abc")`, func(t *testing.T) {
		v := []byte("abc")
		a, err := NewColumnValue(&v).AsBytes()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, string(a.Bytes), "abc")
	})
	t.Run("([]byte)(nil)", func(t *testing.T) {
		a, err := NewColumnValue(([]byte)(nil)).AsBytes()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run("[]byte{}", func(t *testing.T) {
		a, err := NewColumnValue([]byte{}).AsBytes()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, string(a.Bytes), "")
	})
	t.Run(`[]byte("abc")`, func(t *testing.T) {
		a, err := NewColumnValue([]byte("abc")).AsBytes()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, string(a.Bytes), "abc")
	})
	t.Run(`string`, func(t *testing.T) {
		_, err := NewColumnValue("abc").AsBytes()
		assert.NotEqual(t, err, nil)
	})
	t.Run(`int`, func(t *testing.T) {
		_, err := NewColumnValue(1).AsBytes()
		assert.NotEqual(t, err, nil)
	})
}

func TestColumnValue_AsTime(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		a, err := NewColumnValue(nil).AsTime()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run("(*time.Time)(nil)", func(t *testing.T) {
		a, err := NewColumnValue((*time.Time)(nil)).AsTime()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run(`pointer to time.Time`, func(t *testing.T) {
		v, _ := time.Parse(time.RFC3339, "1995-01-10T12:34:56+09:00")
		a, err := NewColumnValue(&v).AsTime()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Time, v)
	})
	t.Run(`time.Time`, func(t *testing.T) {
		v, _ := time.Parse(time.RFC3339, "1995-01-10T12:34:56+09:00")
		a, err := NewColumnValue(v).AsTime()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Time, v)
	})
	t.Run(`invalid sql.NullTime`, func(t *testing.T) {
		v := sql.NullTime{}
		a, err := NewColumnValue(v).AsTime()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, false)
	})
	t.Run(`valid true sql.NullTime`, func(t *testing.T) {
		e, _ := time.Parse(time.RFC3339, "1995-01-10T12:34:56+09:00")
		v := sql.NullTime{Valid: true, Time: e}
		a, err := NewColumnValue(v).AsTime()
		assert.Equal(t, err, nil)
		assert.Equal(t, a.Valid, true)
		assert.Equal(t, a.Time, e)
	})
	t.Run(`string`, func(t *testing.T) {
		_, err := NewColumnValue("abc").AsTime()
		assert.NotEqual(t, err, nil)
	})
}
