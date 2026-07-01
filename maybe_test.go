package fpgo

import (
	"errors"
	"math"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPresent(t *testing.T) {
	var m MaybeDef

	m = Maybe.Just(1)
	assert.Equal(t, true, m.IsPresent())
	assert.Equal(t, false, m.IsNil())

	m = Maybe.Just(nil)
	assert.Equal(t, false, m.IsPresent())
	assert.Equal(t, true, m.IsNil())

	i := 1
	var iptr *int

	iptr = nil
	m = Maybe.Just(iptr)
	assert.Equal(t, false, m.IsPresent())
	assert.Equal(t, true, m.IsNil())

	iptr = &i
	m = Maybe.Just(iptr)
	assert.Equal(t, true, m.IsPresent())
	assert.Equal(t, false, m.IsNil())
}

func TestOr(t *testing.T) {
	var m MaybeDef

	m = Maybe.Just(1)
	assert.Equal(t, 1, m.Or(3))
	m = Maybe.Just(nil)
	assert.Equal(t, 3, m.Or(3))
}

func TestClone(t *testing.T) {
	var m MaybeDef

	i := 1
	i2 := 2
	temp := 3
	var iptr *int
	var iptr2 *int
	iptr2 = &i2

	m = Maybe.Just(1)
	assert.Equal(t, 1, m.Clone().Unwrap())
	assert.Equal(t, 3, temp)
	m = Maybe.Just(1)
	assert.Equal(t, 1, m.Clone().Unwrap())
	assert.Equal(t, 2, *iptr2)
	m = Maybe.Just(nil)
	assert.Equal(t, nil, m.Clone().Unwrap())
	assert.Equal(t, 2, *iptr2)

	iptr = nil
	m = Maybe.Just(iptr)
	assert.Equal(t, nil, m.Clone().Unwrap())
	assert.Equal(t, 2, *iptr2)

	iptr = &i
	m = Maybe.Just(iptr)
	assert.Equal(t, 1, *m.CloneTo(iptr2).ToPtr())
	assert.Equal(t, 1, *iptr2)
}

func TestFlatMap(t *testing.T) {
	var m MaybeDef

	m = Maybe.Just(1).FlatMap(func(in interface{}) MaybeDef {
		v, _ := Maybe.Just(in).ToInt()
		result := Maybe.Just(v + 1)
		return result
	})
	assert.Equal(t, 2, m.Unwrap())
}

func TestLet(t *testing.T) {
	var m MaybeDef

	var letVal int

	letVal = 1
	m = Maybe.Just(1)
	m.Let(func() {
		letVal = 2
	})
	assert.Equal(t, 2, letVal)

	letVal = 1
	m = Maybe.Just(nil)
	m.Let(func() {
		letVal = 3
	})
	assert.Equal(t, 1, letVal)
}

func TestType(t *testing.T) {
	var m MaybeDef

	m = Maybe.Just(1)
	assert.Equal(t, reflect.Int, m.Kind())
	assert.Equal(t, true, m.IsKind(reflect.Int))
	assert.Equal(t, false, m.IsKind(reflect.Ptr))

	assert.Equal(t, reflect.TypeOf(1), m.Type())
	assert.Equal(t, true, m.IsType(reflect.TypeOf(1)))
	assert.Equal(t, false, m.IsType(reflect.TypeOf(nil)))

	m = Maybe.Just(nil)
	assert.Equal(t, reflect.Invalid, m.Kind())
	assert.Equal(t, false, m.IsKind(reflect.Int))
	assert.Equal(t, true, m.IsKind(reflect.Invalid))

	assert.Equal(t, reflect.TypeOf(nil), m.Type())
	assert.Equal(t, true, m.IsType(reflect.TypeOf(nil)))
	assert.Equal(t, false, m.IsType(reflect.TypeOf(1)))
}

func TestCast(t *testing.T) {
	var m MaybeDef

	var f32 float32
	var f64 float64
	var b bool
	var i int
	var i32 int32
	var i64 int64
	var err error

	// Int
	m = Maybe.Just(1)
	assert.Equal(t, "1", m.ToString())

	f32, err = m.ToFloat32()
	assert.Equal(t, float32(1), f32)
	assert.Equal(t, nil, err)
	f64, err = m.ToFloat64()
	assert.Equal(t, float64(1), f64)
	assert.Equal(t, nil, err)
	i, err = m.ToInt()
	assert.Equal(t, 1, i)
	assert.Equal(t, nil, err)
	i32, err = m.ToInt32()
	assert.Equal(t, int32(1), i32)
	assert.Equal(t, nil, err)
	i64, err = m.ToInt64()
	assert.Equal(t, int64(1), i64)
	assert.Equal(t, nil, err)
	b, err = m.ToBool()
	assert.Equal(t, true, b)
	assert.Equal(t, nil, err)

	// Int32
	m = Maybe.Just(int32(1))
	assert.Equal(t, "1", m.ToString())

	f32, err = m.ToFloat32()
	assert.Equal(t, float32(1), f32)
	assert.Equal(t, nil, err)
	f64, err = m.ToFloat64()
	assert.Equal(t, float64(1), f64)
	assert.Equal(t, nil, err)
	i, err = m.ToInt()
	assert.Equal(t, 1, i)
	assert.Equal(t, nil, err)
	i32, err = m.ToInt32()
	assert.Equal(t, int32(1), i32)
	assert.Equal(t, nil, err)
	i64, err = m.ToInt64()
	assert.Equal(t, int64(1), i64)
	assert.Equal(t, nil, err)
	b, err = m.ToBool()
	assert.Equal(t, true, b)
	assert.Equal(t, nil, err)

	// Int64
	m = Maybe.Just(int64(1))
	assert.Equal(t, "1", m.ToString())

	f32, err = m.ToFloat32()
	assert.Equal(t, float32(1), f32)
	assert.Equal(t, nil, err)
	f64, err = m.ToFloat64()
	assert.Equal(t, float64(1), f64)
	assert.Equal(t, nil, err)
	i, err = m.ToInt()
	assert.Equal(t, 1, i)
	assert.Equal(t, nil, err)
	i32, err = m.ToInt32()
	assert.Equal(t, int32(1), i32)
	assert.Equal(t, nil, err)
	i64, err = m.ToInt64()
	assert.Equal(t, int64(1), i64)
	assert.Equal(t, nil, err)
	b, err = m.ToBool()
	assert.Equal(t, true, b)
	assert.Equal(t, nil, err)

	// Float32
	m = Maybe.Just(float32(1.1))
	assert.Equal(t, "1.1", m.ToString())

	f32, err = m.ToFloat32()
	assert.Equal(t, float32(1.1), f32)
	assert.Equal(t, nil, err)
	f64, err = m.ToFloat64()
	assert.Equal(t, float64(1.100000023841858), f64)
	assert.Equal(t, nil, err)
	i, err = m.ToInt()
	assert.Equal(t, 1, i)
	assert.Equal(t, nil, err)
	i32, err = m.ToInt32()
	assert.Equal(t, int32(1), i32)
	assert.Equal(t, nil, err)
	i64, err = m.ToInt64()
	assert.Equal(t, int64(1), i64)
	assert.Equal(t, nil, err)
	b, err = m.ToBool()
	assert.Equal(t, true, b)
	assert.Equal(t, nil, err)

	// Float64
	m = Maybe.Just(float64(1.2))
	assert.Equal(t, "1.2", m.ToString())

	f32, err = m.ToFloat32()
	assert.Equal(t, float32(1.2), f32)
	assert.Equal(t, nil, err)
	f64, err = m.ToFloat64()
	assert.Equal(t, float64(1.2), f64)
	assert.Equal(t, nil, err)
	i, err = m.ToInt()
	assert.Equal(t, 1, i)
	assert.Equal(t, nil, err)
	i32, err = m.ToInt32()
	assert.Equal(t, int32(1), i32)
	assert.Equal(t, nil, err)
	i64, err = m.ToInt64()
	assert.Equal(t, int64(1), i64)
	assert.Equal(t, nil, err)
	b, err = m.ToBool()
	assert.Equal(t, true, b)
	assert.Equal(t, nil, err)

	// Bool(true)
	m = Maybe.Just(true)
	assert.Equal(t, "true", m.ToString())

	f32, err = m.ToFloat32()
	assert.Equal(t, float32(1), f32)
	assert.Equal(t, nil, err)
	f64, err = m.ToFloat64()
	assert.Equal(t, float64(1), f64)
	assert.Equal(t, nil, err)
	i, err = m.ToInt()
	assert.Equal(t, 1, i)
	assert.Equal(t, nil, err)
	i32, err = m.ToInt32()
	assert.Equal(t, int32(1), i32)
	assert.Equal(t, nil, err)
	i64, err = m.ToInt64()
	assert.Equal(t, int64(1), i64)
	assert.Equal(t, nil, err)
	b, err = m.ToBool()
	assert.Equal(t, true, b)
	assert.Equal(t, nil, err)

	// Bool(false)
	m = Maybe.Just(false)
	assert.Equal(t, "false", m.ToString())

	f32, err = m.ToFloat32()
	assert.Equal(t, float32(0), f32)
	assert.Equal(t, nil, err)
	f64, err = m.ToFloat64()
	assert.Equal(t, float64(0), f64)
	assert.Equal(t, nil, err)
	i, err = m.ToInt()
	assert.Equal(t, 0, i)
	assert.Equal(t, nil, err)
	i32, err = m.ToInt32()
	assert.Equal(t, int32(0), i32)
	assert.Equal(t, nil, err)
	i64, err = m.ToInt64()
	assert.Equal(t, int64(0), i64)
	assert.Equal(t, nil, err)
	b, err = m.ToBool()
	assert.Equal(t, false, b)
	assert.Equal(t, nil, err)

	// Nil
	m = Maybe.Just(nil)
	assert.Equal(t, "<nil>", m.ToString())

	f32, err = m.ToFloat32()
	assert.Equal(t, float32(0), f32)
	assert.Equal(t, errors.New("<nil>"), err)
	f64, err = m.ToFloat64()
	assert.Equal(t, float64(0), f64)
	assert.Equal(t, errors.New("<nil>"), err)
	i, err = m.ToInt()
	assert.Equal(t, 0, i)
	assert.Equal(t, errors.New("<nil>"), err)
	i32, err = m.ToInt32()
	assert.Equal(t, int32(0), i32)
	assert.Equal(t, errors.New("<nil>"), err)
	i64, err = m.ToInt64()
	assert.Equal(t, int64(0), i64)
	assert.Equal(t, errors.New("<nil>"), err)
	b, err = m.ToBool()
	assert.Equal(t, false, b)
	assert.Equal(t, errors.New("<nil>"), err)
}

func maybeForConversion(v interface{}) someDef {
	m := Maybe.Just(v)
	if s, ok := m.(someDef); ok {
		return s
	}
	return someDef{ref: nil, isNil: true, isPresent: false}
}

func TestMaybeIsValid(t *testing.T) {
	assert.True(t, Maybe.Just(42).IsValid())
	assert.False(t, Maybe.Just(nil).IsValid())
}

func TestNoneDefSurface(t *testing.T) {
	m := None
	assert.False(t, m.IsPresent())
	assert.True(t, m.IsNil())
	assert.False(t, m.IsValid())
	assert.False(t, m.IsPtr())
	assert.Equal(t, reflect.Invalid, m.Kind())
	assert.Nil(t, m.Unwrap())
	assert.Equal(t, "<nil>", m.ToString())
	assert.Nil(t, m.ToPtr())
	assert.Equal(t, 42, m.Or(42))
	var dest interface{}
	assert.True(t, m.CloneTo(dest).IsNil())
	assert.True(t, m.Clone().IsNil())
	assert.True(t, m.ToMaybe().IsNil())
	var called bool
	m.Let(func() { called = true })
	assert.False(t, called)
}

func TestJustNilMaybeDefSurface(t *testing.T) {
	m := Maybe.Just(nil)
	assert.True(t, m.IsNil())
	assert.False(t, m.IsPresent())
	assert.Equal(t, "<nil>", m.ToString())
	assert.Nil(t, m.ToPtr())
	assert.False(t, m.IsPtr())
	assert.False(t, m.IsValid())
	var called bool
	m.Let(func() { called = true })
	assert.False(t, called)
}

func TestMaybeToPtrBranches(t *testing.T) {
	val := 42
	m := Maybe.Just(&val)
	ptr := m.ToPtr()
	assert.NotNil(t, ptr)
	assert.Equal(t, 42, *ptr)
	m2 := Maybe.Just(99)
	ptr2 := m2.ToPtr()
	assert.NotNil(t, ptr2)
	assert.Equal(t, 99, *ptr2)
}

func TestToBytePorted(t *testing.T) {
	// From int
	m := maybeForConversion(42)
	b, err := m.ToByte()
	assert.Equal(t, byte(42), b)
	assert.Nil(t, err)

	// From int8
	m = maybeForConversion(int8(10))
	b, err = m.ToByte()
	assert.Equal(t, byte(10), b)
	assert.Nil(t, err)

	// From int16
	m = maybeForConversion(int16(100))
	b, err = m.ToByte()
	assert.Equal(t, byte(100), b)
	assert.Nil(t, err)

	// From int32
	m = maybeForConversion(int32(50))
	b, err = m.ToByte()
	assert.Equal(t, byte(50), b)
	assert.Nil(t, err)

	// From int64
	m = maybeForConversion(int64(25))
	b, err = m.ToByte()
	assert.Equal(t, byte(25), b)
	assert.Nil(t, err)

	// From uint
	m = maybeForConversion(uint(200))
	b, err = m.ToByte()
	assert.Equal(t, byte(200), b)
	assert.Nil(t, err)

	// From uint16
	m = maybeForConversion(uint16(150))
	b, err = m.ToByte()
	assert.Equal(t, byte(150), b)
	assert.Nil(t, err)

	// From uint32
	m = maybeForConversion(uint32(75))
	b, err = m.ToByte()
	assert.Equal(t, byte(75), b)
	assert.Nil(t, err)

	// From uint64
	m = maybeForConversion(uint64(125))
	b, err = m.ToByte()
	assert.Equal(t, byte(125), b)
	assert.Nil(t, err)

	// From uintptr
	m = maybeForConversion(uintptr(64))
	b, err = m.ToByte()
	assert.Equal(t, byte(64), b)
	assert.Nil(t, err)

	// From float32
	m = maybeForConversion(float32(3.7))
	b, err = m.ToByte()
	assert.Equal(t, byte(4), b) // rounds to 4
	assert.Nil(t, err)

	// From float64
	m = maybeForConversion(float64(2.3))
	b, err = m.ToByte()
	assert.Equal(t, byte(2), b) // rounds to 2
	assert.Nil(t, err)

	// From bool true
	m = maybeForConversion(true)
	b, err = m.ToByte()
	assert.Equal(t, byte(1), b)
	assert.Nil(t, err)

	// From bool false
	m = maybeForConversion(false)
	b, err = m.ToByte()
	assert.Equal(t, byte(0), b)
	assert.Nil(t, err)

	// From string
	m = maybeForConversion("99")
	b, err = m.ToByte()
	assert.Equal(t, byte(99), b)
	assert.Nil(t, err)

	// From byte
	m = maybeForConversion(byte(255))
	b, err = m.ToByte()
	assert.Equal(t, byte(255), b)
	assert.Nil(t, err)

	// From nil
	m = maybeForConversion(nil)
	b, err = m.ToByte()
	assert.Equal(t, byte(0), b)
	assert.Equal(t, ErrConversionUnsupported, err)

	// Overflow from int
	m = maybeForConversion(256)
	b, err = m.ToByte()
	assert.Equal(t, byte(0), b)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from int (negative)
	m = maybeForConversion(-1)
	b, err = m.ToByte()
	assert.Equal(t, byte(0), b)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint
	m = maybeForConversion(uint(300))
	b, err = m.ToByte()
	assert.Equal(t, byte(0), b)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Unsupported type
	m = maybeForConversion([]int{1, 2, 3})
	b, err = m.ToByte()
	assert.Equal(t, byte(0), b)
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestToInt8Ported(t *testing.T) {
	// From int
	m := maybeForConversion(42)
	i, err := m.ToInt8()
	assert.Equal(t, int8(42), i)
	assert.Nil(t, err)

	// From int8
	m = maybeForConversion(int8(10))
	i, err = m.ToInt8()
	assert.Equal(t, int8(10), i)
	assert.Nil(t, err)

	// From int16
	m = maybeForConversion(int16(100))
	i, err = m.ToInt8()
	assert.Equal(t, int8(100), i)
	assert.Nil(t, err)

	// From int32
	m = maybeForConversion(int32(50))
	i, err = m.ToInt8()
	assert.Equal(t, int8(50), i)
	assert.Nil(t, err)

	// From int64
	m = maybeForConversion(int64(25))
	i, err = m.ToInt8()
	assert.Equal(t, int8(25), i)
	assert.Nil(t, err)

	// From uint
	m = maybeForConversion(uint(100))
	i, err = m.ToInt8()
	assert.Equal(t, int8(100), i)
	assert.Nil(t, err)

	// From uint16
	m = maybeForConversion(uint16(50))
	i, err = m.ToInt8()
	assert.Equal(t, int8(50), i)
	assert.Nil(t, err)

	// From uint32
	m = maybeForConversion(uint32(75))
	i, err = m.ToInt8()
	assert.Equal(t, int8(75), i)
	assert.Nil(t, err)

	// From uint64
	m = maybeForConversion(uint64(125))
	i, err = m.ToInt8()
	assert.Equal(t, int8(125), i)
	assert.Nil(t, err)

	// From uintptr
	m = maybeForConversion(uintptr(64))
	i, err = m.ToInt8()
	assert.Equal(t, int8(64), i)
	assert.Nil(t, err)

	// From float32
	m = maybeForConversion(float32(3.7))
	i, err = m.ToInt8()
	assert.Equal(t, int8(4), i) // rounds to 4
	assert.Nil(t, err)

	// From float64
	m = maybeForConversion(float64(2.3))
	i, err = m.ToInt8()
	assert.Equal(t, int8(2), i) // rounds to 2
	assert.Nil(t, err)

	// From bool true
	m = maybeForConversion(true)
	i, err = m.ToInt8()
	assert.Equal(t, int8(1), i)
	assert.Nil(t, err)

	// From bool false
	m = maybeForConversion(false)
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Nil(t, err)

	// From string
	m = maybeForConversion("99")
	i, err = m.ToInt8()
	assert.Equal(t, int8(99), i)
	assert.Nil(t, err)

	// From byte
	m = maybeForConversion(byte(100))
	i, err = m.ToInt8()
	assert.Equal(t, int8(100), i)
	assert.Nil(t, err)

	// From nil
	m = maybeForConversion(nil)
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Equal(t, ErrConversionUnsupported, err)

	// Overflow from int (positive)
	m = maybeForConversion(128)
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from int (negative)
	m = maybeForConversion(-129)
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint
	m = maybeForConversion(uint(200))
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint16
	m = maybeForConversion(uint16(200))
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint32
	m = maybeForConversion(uint32(200))
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint64
	m = maybeForConversion(uint64(200))
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uintptr
	m = maybeForConversion(uintptr(200))
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Unsupported type
	m = maybeForConversion([]int{1, 2, 3})
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestToInt16Ported(t *testing.T) {
	// From int
	m := maybeForConversion(1000)
	i, err := m.ToInt16()
	assert.Equal(t, int16(1000), i)
	assert.Nil(t, err)

	// From int8
	m = maybeForConversion(int8(10))
	i, err = m.ToInt16()
	assert.Equal(t, int16(10), i)
	assert.Nil(t, err)

	// From int16
	m = maybeForConversion(int16(500))
	i, err = m.ToInt16()
	assert.Equal(t, int16(500), i)
	assert.Nil(t, err)

	// From int32
	m = maybeForConversion(int32(250))
	i, err = m.ToInt16()
	assert.Equal(t, int16(250), i)
	assert.Nil(t, err)

	// From int64
	m = maybeForConversion(int64(125))
	i, err = m.ToInt16()
	assert.Equal(t, int16(125), i)
	assert.Nil(t, err)

	// From uint
	m = maybeForConversion(uint(1000))
	i, err = m.ToInt16()
	assert.Equal(t, int16(1000), i)
	assert.Nil(t, err)

	// From uint16
	m = maybeForConversion(uint16(500))
	i, err = m.ToInt16()
	assert.Equal(t, int16(500), i)
	assert.Nil(t, err)

	// From uint32
	m = maybeForConversion(uint32(250))
	i, err = m.ToInt16()
	assert.Equal(t, int16(250), i)
	assert.Nil(t, err)

	// From uint64
	m = maybeForConversion(uint64(125))
	i, err = m.ToInt16()
	assert.Equal(t, int16(125), i)
	assert.Nil(t, err)

	// From uintptr
	m = maybeForConversion(uintptr(64))
	i, err = m.ToInt16()
	assert.Equal(t, int16(64), i)
	assert.Nil(t, err)

	// From float32
	m = maybeForConversion(float32(3.7))
	i, err = m.ToInt16()
	assert.Equal(t, int16(4), i) // rounds to 4
	assert.Nil(t, err)

	// From float64
	m = maybeForConversion(float64(2.3))
	i, err = m.ToInt16()
	assert.Equal(t, int16(2), i) // rounds to 2
	assert.Nil(t, err)

	// From bool true
	m = maybeForConversion(true)
	i, err = m.ToInt16()
	assert.Equal(t, int16(1), i)
	assert.Nil(t, err)

	// From bool false
	m = maybeForConversion(false)
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Nil(t, err)

	// From string
	m = maybeForConversion("999")
	i, err = m.ToInt16()
	assert.Equal(t, int16(999), i)
	assert.Nil(t, err)

	// From byte
	m = maybeForConversion(byte(100))
	i, err = m.ToInt16()
	assert.Equal(t, int16(100), i)
	assert.Nil(t, err)

	// From nil
	m = maybeForConversion(nil)
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Equal(t, ErrConversionUnsupported, err)

	// Overflow from int (positive)
	m = maybeForConversion(32768) // MaxInt16 + 1
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from int (negative)
	m = maybeForConversion(-32769) // MinInt16 - 1
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint
	m = maybeForConversion(uint(40000))
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint16
	m = maybeForConversion(uint16(40000))
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint32
	m = maybeForConversion(uint32(40000))
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint64
	m = maybeForConversion(uint64(40000))
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uintptr
	m = maybeForConversion(uintptr(40000))
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Unsupported type
	m = maybeForConversion([]int{1, 2, 3})
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestToUintPorted(t *testing.T) {
	// From int
	m := maybeForConversion(42)
	u, err := m.ToUint()
	assert.Equal(t, uint(42), u)
	assert.Nil(t, err)

	// From int8
	m = maybeForConversion(int8(10))
	u, err = m.ToUint()
	assert.Equal(t, uint(10), u)
	assert.Nil(t, err)

	// From int16
	m = maybeForConversion(int16(100))
	u, err = m.ToUint()
	assert.Equal(t, uint(100), u)
	assert.Nil(t, err)

	// From int32
	m = maybeForConversion(int32(50))
	u, err = m.ToUint()
	assert.Equal(t, uint(50), u)
	assert.Nil(t, err)

	// From int64
	m = maybeForConversion(int64(25))
	u, err = m.ToUint()
	assert.Equal(t, uint(25), u)
	assert.Nil(t, err)

	// From uint
	m = maybeForConversion(uint(200))
	u, err = m.ToUint()
	assert.Equal(t, uint(200), u)
	assert.Nil(t, err)

	// From uint16
	m = maybeForConversion(uint16(150))
	u, err = m.ToUint()
	assert.Equal(t, uint(150), u)
	assert.Nil(t, err)

	// From uint32
	m = maybeForConversion(uint32(75))
	u, err = m.ToUint()
	assert.Equal(t, uint(75), u)
	assert.Nil(t, err)

	// From uint64
	m = maybeForConversion(uint64(125))
	u, err = m.ToUint()
	assert.Equal(t, uint(125), u)
	assert.Nil(t, err)

	// From uintptr
	m = maybeForConversion(uintptr(64))
	u, err = m.ToUint()
	assert.Equal(t, uint(64), u)
	assert.Nil(t, err)

	// From float32
	m = maybeForConversion(float32(3.7))
	u, err = m.ToUint()
	assert.Equal(t, uint(4), u) // rounds to 4
	assert.Nil(t, err)

	// From float64
	m = maybeForConversion(float64(2.3))
	u, err = m.ToUint()
	assert.Equal(t, uint(2), u) // rounds to 2
	assert.Nil(t, err)

	// From bool true
	m = maybeForConversion(true)
	u, err = m.ToUint()
	assert.Equal(t, uint(1), u)
	assert.Nil(t, err)

	// From bool false
	m = maybeForConversion(false)
	u, err = m.ToUint()
	assert.Equal(t, uint(0), u)
	assert.Nil(t, err)

	// From string
	m = maybeForConversion("99")
	u, err = m.ToUint()
	assert.Equal(t, uint(99), u)
	assert.Nil(t, err)

	// From byte
	m = maybeForConversion(byte(255))
	u, err = m.ToUint()
	assert.Equal(t, uint(255), u)
	assert.Nil(t, err)

	// From nil
	m = maybeForConversion(nil)
	u, err = m.ToUint()
	assert.Equal(t, uint(0), u)
	assert.Equal(t, ErrConversionUnsupported, err)

	// Unsupported type
	m = maybeForConversion([]int{1, 2, 3})
	u, err = m.ToUint()
	assert.Equal(t, uint(0), u)
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestToUint8Ported(t *testing.T) {
	// From int
	m := maybeForConversion(42)
	u, err := m.ToUint8()
	assert.Equal(t, uint8(42), u)
	assert.Nil(t, err)

	// From int8
	m = maybeForConversion(int8(10))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(10), u)
	assert.Nil(t, err)

	// From int16
	m = maybeForConversion(int16(100))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(100), u)
	assert.Nil(t, err)

	// From int32
	m = maybeForConversion(int32(50))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(50), u)
	assert.Nil(t, err)

	// From int64
	m = maybeForConversion(int64(25))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(25), u)
	assert.Nil(t, err)

	// From uint
	m = maybeForConversion(uint(200))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(200), u)
	assert.Nil(t, err)

	// From uint16
	m = maybeForConversion(uint16(150))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(150), u)
	assert.Nil(t, err)

	// From uint32
	m = maybeForConversion(uint32(75))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(75), u)
	assert.Nil(t, err)

	// From uint64
	m = maybeForConversion(uint64(125))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(125), u)
	assert.Nil(t, err)

	// From uintptr
	m = maybeForConversion(uintptr(64))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(64), u)
	assert.Nil(t, err)

	// From float32
	m = maybeForConversion(float32(3.7))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(4), u) // rounds to 4
	assert.Nil(t, err)

	// From float64
	m = maybeForConversion(float64(2.3))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(2), u) // rounds to 2
	assert.Nil(t, err)

	// From bool true
	m = maybeForConversion(true)
	u, err = m.ToUint8()
	assert.Equal(t, uint8(1), u)
	assert.Nil(t, err)

	// From bool false
	m = maybeForConversion(false)
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Nil(t, err)

	// From string
	m = maybeForConversion("99")
	u, err = m.ToUint8()
	assert.Equal(t, uint8(99), u)
	assert.Nil(t, err)

	// From byte
	m = maybeForConversion(byte(255))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(255), u)
	assert.Nil(t, err)

	// From nil
	m = maybeForConversion(nil)
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Equal(t, ErrConversionUnsupported, err)

	// Overflow from int
	m = maybeForConversion(256)
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from int (negative)
	m = maybeForConversion(-1)
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint
	m = maybeForConversion(uint(300))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint16
	m = maybeForConversion(uint16(300))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint32
	m = maybeForConversion(uint32(300))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint64
	m = maybeForConversion(uint64(300))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uintptr
	m = maybeForConversion(uintptr(300))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Unsupported type
	m = maybeForConversion([]int{1, 2, 3})
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestToUint16Ported(t *testing.T) {
	// From int
	m := maybeForConversion(1000)
	u, err := m.ToUint16()
	assert.Equal(t, uint16(1000), u)
	assert.Nil(t, err)

	// From int8
	m = maybeForConversion(int8(10))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(10), u)
	assert.Nil(t, err)

	// From int16
	m = maybeForConversion(int16(500))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(500), u)
	assert.Nil(t, err)

	// From int32
	m = maybeForConversion(int32(250))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(250), u)
	assert.Nil(t, err)

	// From int64
	m = maybeForConversion(int64(125))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(125), u)
	assert.Nil(t, err)

	// From uint
	m = maybeForConversion(uint(1000))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(1000), u)
	assert.Nil(t, err)

	// From uint16
	m = maybeForConversion(uint16(500))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(500), u)
	assert.Nil(t, err)

	// From uint32
	m = maybeForConversion(uint32(250))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(250), u)
	assert.Nil(t, err)

	// From uint64
	m = maybeForConversion(uint64(125))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(125), u)
	assert.Nil(t, err)

	// From uintptr
	m = maybeForConversion(uintptr(64))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(64), u)
	assert.Nil(t, err)

	// From float32
	m = maybeForConversion(float32(3.7))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(4), u) // rounds to 4
	assert.Nil(t, err)

	// From float64
	m = maybeForConversion(float64(2.3))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(2), u) // rounds to 2
	assert.Nil(t, err)

	// From bool true
	m = maybeForConversion(true)
	u, err = m.ToUint16()
	assert.Equal(t, uint16(1), u)
	assert.Nil(t, err)

	// From bool false
	m = maybeForConversion(false)
	u, err = m.ToUint16()
	assert.Equal(t, uint16(0), u)
	assert.Nil(t, err)

	// From string
	m = maybeForConversion("999")
	u, err = m.ToUint16()
	assert.Equal(t, uint16(999), u)
	assert.Nil(t, err)

	// From byte
	m = maybeForConversion(byte(100))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(100), u)
	assert.Nil(t, err)

	// From nil
	m = maybeForConversion(nil)
	u, err = m.ToUint16()
	assert.Equal(t, uint16(0), u)
	assert.Equal(t, ErrConversionUnsupported, err)

	// Overflow from int
	m = maybeForConversion(65536) // MaxUint16 + 1
	u, err = m.ToUint16()
	assert.Equal(t, uint16(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from int (negative)
	m = maybeForConversion(-1)
	u, err = m.ToUint16()
	assert.Equal(t, uint16(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint
	m = maybeForConversion(uint(70000))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint32
	m = maybeForConversion(uint32(70000))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint64
	m = maybeForConversion(uint64(70000))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uintptr
	m = maybeForConversion(uintptr(70000))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Unsupported type
	m = maybeForConversion([]int{1, 2, 3})
	u, err = m.ToUint16()
	assert.Equal(t, uint16(0), u)
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestToUint32Ported(t *testing.T) {
	// From int
	m := maybeForConversion(1000)
	u, err := m.ToUint32()
	assert.Equal(t, uint32(1000), u)
	assert.Nil(t, err)

	// From int8
	m = maybeForConversion(int8(10))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(10), u)
	assert.Nil(t, err)

	// From int16
	m = maybeForConversion(int16(500))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(500), u)
	assert.Nil(t, err)

	// From int32
	m = maybeForConversion(int32(250))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(250), u)
	assert.Nil(t, err)

	// From int64
	m = maybeForConversion(int64(125))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(125), u)
	assert.Nil(t, err)

	// From uint
	m = maybeForConversion(uint(1000))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(1000), u)
	assert.Nil(t, err)

	// From uint16
	m = maybeForConversion(uint16(500))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(500), u)
	assert.Nil(t, err)

	// From uint32
	m = maybeForConversion(uint32(250))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(250), u)
	assert.Nil(t, err)

	// From uint64
	m = maybeForConversion(uint64(125))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(125), u)
	assert.Nil(t, err)

	// From uintptr
	m = maybeForConversion(uintptr(64))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(64), u)
	assert.Nil(t, err)

	// From float32
	m = maybeForConversion(float32(3.7))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(4), u) // rounds to 4
	assert.Nil(t, err)

	// From float64
	m = maybeForConversion(float64(2.3))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(2), u) // rounds to 2
	assert.Nil(t, err)

	// From bool true
	m = maybeForConversion(true)
	u, err = m.ToUint32()
	assert.Equal(t, uint32(1), u)
	assert.Nil(t, err)

	// From bool false
	m = maybeForConversion(false)
	u, err = m.ToUint32()
	assert.Equal(t, uint32(0), u)
	assert.Nil(t, err)

	// From string
	m = maybeForConversion("999")
	u, err = m.ToUint32()
	assert.Equal(t, uint32(999), u)
	assert.Nil(t, err)

	// From byte
	m = maybeForConversion(byte(100))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(100), u)
	assert.Nil(t, err)

	// From nil
	m = maybeForConversion(nil)
	u, err = m.ToUint32()
	assert.Equal(t, uint32(0), u)
	assert.Equal(t, ErrConversionUnsupported, err)

	// Overflow from int (negative)
	m = maybeForConversion(-1)
	u, err = m.ToUint32()
	assert.Equal(t, uint32(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from int64
	m = maybeForConversion(int64(4294967296)) // MaxUint32 + 1
	u, err = m.ToUint32()
	assert.Equal(t, uint32(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint64
	m = maybeForConversion(uint64(4294967296)) // MaxUint32 + 1
	u, err = m.ToUint32()
	assert.Equal(t, uint32(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Unsupported type
	m = maybeForConversion([]int{1, 2, 3})
	u, err = m.ToUint32()
	assert.Equal(t, uint32(0), u)
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestToUint64Ported(t *testing.T) {
	// From int
	m := maybeForConversion(1000)
	u, err := m.ToUint64()
	assert.Equal(t, uint64(1000), u)
	assert.Nil(t, err)

	// From int8
	m = maybeForConversion(int8(10))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(10), u)
	assert.Nil(t, err)

	// From int16
	m = maybeForConversion(int16(500))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(500), u)
	assert.Nil(t, err)

	// From int32
	m = maybeForConversion(int32(250))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(250), u)
	assert.Nil(t, err)

	// From int64
	m = maybeForConversion(int64(125))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(125), u)
	assert.Nil(t, err)

	// From uint
	m = maybeForConversion(uint(1000))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(1000), u)
	assert.Nil(t, err)

	// From uint16
	m = maybeForConversion(uint16(500))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(500), u)
	assert.Nil(t, err)

	// From uint32
	m = maybeForConversion(uint32(250))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(250), u)
	assert.Nil(t, err)

	// From uint64
	m = maybeForConversion(uint64(125))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(125), u)
	assert.Nil(t, err)

	// From uintptr
	m = maybeForConversion(uintptr(64))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(64), u)
	assert.Nil(t, err)

	// From float32
	m = maybeForConversion(float32(3.7))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(4), u) // rounds to 4
	assert.Nil(t, err)

	// From float64
	m = maybeForConversion(float64(2.3))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(2), u) // rounds to 2
	assert.Nil(t, err)

	// From bool true
	m = maybeForConversion(true)
	u, err = m.ToUint64()
	assert.Equal(t, uint64(1), u)
	assert.Nil(t, err)

	// From bool false
	m = maybeForConversion(false)
	u, err = m.ToUint64()
	assert.Equal(t, uint64(0), u)
	assert.Nil(t, err)

	// From string
	m = maybeForConversion("999")
	u, err = m.ToUint64()
	assert.Equal(t, uint64(999), u)
	assert.Nil(t, err)

	// From byte
	m = maybeForConversion(byte(100))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(100), u)
	assert.Nil(t, err)

	// From nil
	m = maybeForConversion(nil)
	u, err = m.ToUint64()
	assert.Equal(t, uint64(0), u)
	assert.Equal(t, ErrConversionUnsupported, err)

	// Unsupported type
	m = maybeForConversion([]int{1, 2, 3})
	u, err = m.ToUint64()
	assert.Equal(t, uint64(0), u)
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestToUintptrPorted(t *testing.T) {
	// From int
	m := maybeForConversion(1000)
	u, err := m.ToUintptr()
	assert.Equal(t, uintptr(1000), u)
	assert.Nil(t, err)

	// From int8
	m = maybeForConversion(int8(10))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(10), u)
	assert.Nil(t, err)

	// From int16
	m = maybeForConversion(int16(500))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(500), u)
	assert.Nil(t, err)

	// From int32
	m = maybeForConversion(int32(250))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(250), u)
	assert.Nil(t, err)

	// From int64
	m = maybeForConversion(int64(125))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(125), u)
	assert.Nil(t, err)

	// From uint
	m = maybeForConversion(uint(1000))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(1000), u)
	assert.Nil(t, err)

	// From uint16
	m = maybeForConversion(uint16(500))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(500), u)
	assert.Nil(t, err)

	// From uint32
	m = maybeForConversion(uint32(250))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(250), u)
	assert.Nil(t, err)

	// From uint64
	m = maybeForConversion(uint64(125))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(125), u)
	assert.Nil(t, err)

	// From uintptr
	m = maybeForConversion(uintptr(64))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(64), u)
	assert.Nil(t, err)

	// From float32
	m = maybeForConversion(float32(3.7))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(4), u) // rounds to 4
	assert.Nil(t, err)

	// From float64
	m = maybeForConversion(float64(2.3))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(2), u) // rounds to 2
	assert.Nil(t, err)

	// From bool true
	m = maybeForConversion(true)
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(1), u)
	assert.Nil(t, err)

	// From bool false
	m = maybeForConversion(false)
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(0), u)
	assert.Nil(t, err)

	// From string
	m = maybeForConversion("999")
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(999), u)
	assert.Nil(t, err)

	// From byte
	m = maybeForConversion(byte(100))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(100), u)
	assert.Nil(t, err)

	// From nil
	m = maybeForConversion(nil)
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(0), u)
	assert.Equal(t, ErrConversionUnsupported, err)

	// Unsupported type
	m = maybeForConversion([]int{1, 2, 3})
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(0), u)
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestMaybeUnsupportedAndOverflowConversionsPorted(t *testing.T) {
	unsupported := maybeForConversion(errors.New("boom"))

	_, err := unsupported.ToFloat32()
	assert.Equal(t, ErrConversionUnsupported, err)
	_, err = unsupported.ToFloat64()
	assert.Equal(t, ErrConversionUnsupported, err)
	_, err = unsupported.ToInt()
	assert.Equal(t, ErrConversionUnsupported, err)

	overflowUint := maybeForConversion(uint(math.MaxUint32))
	_, err = overflowUint.ToInt()
	assert.Equal(t, ErrConversionSizeOverflow, err)

	overflowUint32 := maybeForConversion(uint32(math.MaxUint32))
	_, err = overflowUint32.ToInt()
	assert.Equal(t, ErrConversionSizeOverflow, err)

	overflowUint64 := maybeForConversion(uint64(math.MaxUint64))
	_, err = overflowUint64.ToInt()
	assert.Equal(t, ErrConversionSizeOverflow, err)
	_, err = overflowUint64.ToUint32()
	assert.Equal(t, ErrConversionSizeOverflow, err)
	_, err = overflowUint64.ToUintptr()
	assert.NoError(t, err)

	negativeFloat := maybeForConversion(float64(-1.2))
	_, err = negativeFloat.ToUint()
	assert.Equal(t, ErrConversionSizeOverflow, err)
	valUint64, err := negativeFloat.ToUint64()
	assert.Equal(t, ErrConversionSizeOverflow, err)
	assert.Equal(t, uint64(0), valUint64)

	largeFloat := maybeForConversion(float64(math.MaxFloat64))
	_, err = largeFloat.ToUint64()
	assert.Equal(t, ErrConversionSizeOverflow, err)

	_, err = maybeForConversion(uintptr(7)).ToUint64()
	assert.NoError(t, err)
	_, err = maybeForConversion(uintptr(7)).ToUint32()
	assert.NoError(t, err)
	_, err = maybeForConversion(uintptr(7)).ToUint()
	assert.NoError(t, err)
	_, err = maybeForConversion(uintptr(7)).ToInt32()
	assert.NoError(t, err)
	valUintptr, err := maybeForConversion(int64(-1)).ToUintptr()
	assert.NoError(t, err)
	assert.Equal(t, ^uintptr(0), valUintptr)
	_, err = maybeForConversion(float64(12.6)).ToUint32()
	assert.NoError(t, err)
	_, err = maybeForConversion(float32(12.4)).ToUint64()
	assert.NoError(t, err)
	_, err = maybeForConversion(float64(12.6)).ToUintptr()
	assert.NoError(t, err)
}

func TestMaybeNoneConversionsPorted(t *testing.T) {
	m := None

	assert.False(t, m.IsPresent())
	assert.True(t, m.IsNil())

	_, err := m.ToFloat64()
	assert.Equal(t, ErrConversionNil, err)
	_, err = m.ToFloat32()
	assert.Equal(t, ErrConversionNil, err)
	_, err = m.ToInt()
	assert.Equal(t, ErrConversionNil, err)
	_, err = m.ToInt32()
	assert.Equal(t, ErrConversionNil, err)
	_, err = m.ToInt64()
	assert.Equal(t, ErrConversionNil, err)
	_, err = m.ToBool()
	assert.Equal(t, ErrConversionNil, err)

	assert.Equal(t, "<nil>", m.ToString())
	assert.Nil(t, m.ToPtr())
	assert.False(t, m.IsValid())
	assert.False(t, m.IsPtr())

	called := false
	m.Let(func() { called = true })
	assert.False(t, called)

	result := m.Or(42)
	assert.Equal(t, 42, result)

	cloned := m.Clone()
	assert.True(t, cloned.IsNil())

	toMaybe := m.ToMaybe()
	assert.True(t, toMaybe.IsNil())
}

func TestMaybeOverflowAndUnsupportedConversionsPorted(t *testing.T) {
	_, err := maybeForConversion(uint64(math.MaxUint64)).ToInt()
	assert.ErrorIs(t, err, ErrConversionSizeOverflow)

	// float64->int32 on non_generic has no upper bound check (unlike float32 path).
	got, err := maybeForConversion(float64(math.MaxFloat64)).ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(math.MaxInt32), got)

	valUint, err := maybeForConversion(int64(-1)).ToUint()
	assert.Equal(t, uint(0), valUint)
	assert.ErrorIs(t, err, ErrConversionSizeOverflow)

	// Negative float -> uint32 now guarded (matches ToUint64), returns overflow.
	valUint32, err := maybeForConversion(float64(-1)).ToUint32()
	assert.Equal(t, uint32(0), valUint32)
	assert.ErrorIs(t, err, ErrConversionSizeOverflow)

	valUint64, err := maybeForConversion(float64(-1)).ToUint64()
	assert.Equal(t, uint64(0), valUint64)
	assert.ErrorIs(t, err, ErrConversionSizeOverflow)

	_, err = maybeForConversion("18446744073709551616").ToUintptr()
	assert.Error(t, err)

	type unsupported struct{}
	_, err = maybeForConversion(unsupported{}).ToInt()
	assert.ErrorIs(t, err, ErrConversionUnsupported)
}

func TestSomeDefToFloat32WithStringPorted(t *testing.T) {
	m := maybeForConversion("42.5")
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42.5), val)
}

func TestSomeDefToFloat64WithStringPorted(t *testing.T) {
	m := maybeForConversion("42.5")
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42.5), val)
}

func TestToBoolWithIntPorted(t *testing.T) {
	m := maybeForConversion(0)
	val, err := m.ToBool()
	assert.NoError(t, err)
	assert.False(t, val)
}

func TestToBoolWithUintPorted(t *testing.T) {
	m := maybeForConversion(uint(1))
	val, err := m.ToBool()
	assert.NoError(t, err)
	assert.True(t, val)
}

func TestSomeDefToBoolWithFloat32Ported(t *testing.T) {
	m := maybeForConversion(float32(1.5))
	val, err := m.ToBool()
	assert.True(t, val)
	assert.Nil(t, err)
}

func TestToInt32WithUint16Ported(t *testing.T) {
	m := maybeForConversion(uint16(42))
	val, err := m.ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(42), val)
}

func TestToInt32WithInt64Ported(t *testing.T) {
	m := maybeForConversion(int64(42))
	val, err := m.ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(42), val)
}

func TestToInt64WithStringPorted(t *testing.T) {
	m := maybeForConversion("999")
	val, err := m.ToInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(999), val)
}

func TestToIntWithStringPorted(t *testing.T) {
	m := maybeForConversion("123")
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 123, val)
}

func TestToIntWithFloat32Ported(t *testing.T) {
	m := maybeForConversion(float32(42.5))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 43, val)
}

func TestSomeDefToBoolWithBoolFalsePorted(t *testing.T) {
	m := maybeForConversion(false)
	val, err := m.ToBool()
	assert.NoError(t, err)
	assert.False(t, val)
}

func TestSomeDefToBoolWithBoolTruePorted(t *testing.T) {
	m := maybeForConversion(true)
	val, err := m.ToBool()
	assert.NoError(t, err)
	assert.True(t, val)
}

func TestSomeDefToBoolWithBytePorted(t *testing.T) {
	m := maybeForConversion(byte(1))
	val, err := m.ToBool()
	assert.NoError(t, err)
	assert.True(t, val)
}

func TestSomeDefToBoolWithByteZeroPorted(t *testing.T) {
	m := maybeForConversion(byte(0))
	val, err := m.ToBool()
	assert.NoError(t, err)
	assert.False(t, val)
}

func TestSomeDefToBoolWithFloat64Ported(t *testing.T) {
	m := maybeForConversion(float64(0))
	val, err := m.ToBool()
	assert.False(t, val)
	assert.Nil(t, err)
}

func TestSomeDefToBoolWithInt16Ported(t *testing.T) {
	m := maybeForConversion(int16(0))
	val, err := m.ToBool()
	assert.False(t, val)
	assert.Nil(t, err)
}

func TestSomeDefToBoolWithInt32Ported(t *testing.T) {
	m := maybeForConversion(int32(0))
	val, err := m.ToBool()
	assert.False(t, val)
	assert.Nil(t, err)
}

func TestSomeDefToBoolWithInt64Ported(t *testing.T) {
	m := maybeForConversion(int64(1))
	val, err := m.ToBool()
	assert.True(t, val)
	assert.Nil(t, err)
}

func TestSomeDefToBoolWithInt8Ported(t *testing.T) {
	m := maybeForConversion(int8(1))
	val, err := m.ToBool()
	assert.True(t, val)
	assert.Nil(t, err)
}

func TestSomeDefToBoolWithNilPorted(t *testing.T) {
	m := maybeForConversion(nil)
	_, err := m.ToBool()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToBoolWithStringPorted(t *testing.T) {
	m := maybeForConversion("true")
	val, err := m.ToBool()
	assert.NoError(t, err)
	assert.True(t, val)
}

func TestSomeDefToBoolWithUint16Ported(t *testing.T) {
	m := maybeForConversion(uint16(0))
	val, err := m.ToBool()
	assert.False(t, val)
	assert.Nil(t, err)
}

func TestSomeDefToBoolWithUint32Ported(t *testing.T) {
	m := maybeForConversion(uint32(1))
	val, err := m.ToBool()
	assert.True(t, val)
	assert.Nil(t, err)
}

func TestSomeDefToBoolWithUint64Ported(t *testing.T) {
	m := maybeForConversion(uint64(0))
	val, err := m.ToBool()
	assert.False(t, val)
	assert.Nil(t, err)
}

func TestSomeDefToBoolWithUintptrPorted(t *testing.T) {
	m := maybeForConversion(uintptr(1))
	val, err := m.ToBool()
	assert.True(t, val)
	assert.Nil(t, err)
}

func TestSomeDefToBoolWithUnsupportedTypePorted(t *testing.T) {
	m := maybeForConversion(make(chan int))
	_, err := m.ToBool()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToInt32WithFloat32OverflowPorted(t *testing.T) {
	m := maybeForConversion(float32(math.MaxInt32 * 2))
	_, err := m.ToInt32()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt32WithFloat64Ported(t *testing.T) {
	m := maybeForConversion(float64(42.5))
	val, err := m.ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(43), val)
}

func TestSomeDefToInt32WithInt32Ported(t *testing.T) {
	m := maybeForConversion(int32(42))
	val, err := m.ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(42), val)
}

func TestSomeDefToInt32WithInt64OverflowPorted(t *testing.T) {
	m := maybeForConversion(int64(math.MaxInt32 + 1))
	_, err := m.ToInt32()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt32WithInt64UnderflowPorted(t *testing.T) {
	m := maybeForConversion(int64(math.MinInt32 - 1))
	_, err := m.ToInt32()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt32WithNilPorted(t *testing.T) {
	m := maybeForConversion(nil)
	_, err := m.ToInt32()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToInt32WithStringPorted(t *testing.T) {
	m := maybeForConversion("42")
	val, err := m.ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(42), val)
}

func TestSomeDefToInt32WithUint32OverflowPorted(t *testing.T) {
	m := maybeForConversion(uint32(math.MaxInt32 + 1))
	_, err := m.ToInt32()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt32WithUint64OverflowPorted(t *testing.T) {
	m := maybeForConversion(uint64(math.MaxInt32 + 1))
	_, err := m.ToInt32()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt32WithUint64WithinPorted(t *testing.T) {
	val, err := maybeForConversion(uint64(5)).ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(5), val)
}

func TestSomeDefToInt32WithUintOverflowPorted(t *testing.T) {
	m := maybeForConversion(uint(math.MaxInt32 + 1))
	_, err := m.ToInt32()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt32WithUintptrOverflowPorted(t *testing.T) {
	m := maybeForConversion(uintptr(math.MaxInt32 + 1))
	_, err := m.ToInt32()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt32WithUnsupportedTypePorted(t *testing.T) {
	m := maybeForConversion(make(chan int))
	_, err := m.ToInt32()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToInt64WithFloat32OverflowPorted(t *testing.T) {
	m := maybeForConversion(float32(math.MaxInt32 * 2))
	_, err := m.ToInt64()
	assert.NoError(t, err) // Should succeed since float32*2 < MaxInt64
}

func TestSomeDefToInt64WithFloat64OverflowPorted(t *testing.T) {
	m := maybeForConversion(float64(math.MaxInt64) * 2)
	_, err := m.ToInt64()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt64WithFloat64UnderflowPorted(t *testing.T) {
	m := maybeForConversion(float64(math.MinInt64) * 2)
	_, err := m.ToInt64()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt64WithInt64Ported(t *testing.T) {
	m := maybeForConversion(int64(42))
	val, err := m.ToInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(42), val)
}

func TestSomeDefToInt64WithNilPorted(t *testing.T) {
	m := maybeForConversion(nil)
	_, err := m.ToInt64()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToInt64WithStringPorted(t *testing.T) {
	m := maybeForConversion("42")
	val, err := m.ToInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(42), val)
}

func TestSomeDefToInt64WithUintPorted(t *testing.T) {
	m := maybeForConversion(uint(42))
	val, err := m.ToInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(42), val)
}

func TestSomeDefToInt64WithUint64OverflowPorted(t *testing.T) {
	m := maybeForConversion(uint64(math.MaxInt64 + 1))
	_, err := m.ToInt64()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt64WithUintOverflowPorted(t *testing.T) {
	// On 64-bit systems, uint can be > MaxInt64
	if math.MaxUint > math.MaxInt64 {
		m := maybeForConversion(uint(math.MaxUint))
		_, err := m.ToInt64()
		assert.Equal(t, ErrConversionSizeOverflow, err)
	}
}

func TestSomeDefToInt64WithUintptrOverflowPorted(t *testing.T) {
	// On 64-bit systems, uintptr can be large
	m := maybeForConversion(uintptr(1 << 62))
	val, err := m.ToInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(1<<62), val)
}

func TestSomeDefToInt64WithUnsupportedTypePorted(t *testing.T) {
	m := maybeForConversion(make(chan int))
	_, err := m.ToInt64()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestToBoolWithIntNonZeroPorted(t *testing.T) {
	m := maybeForConversion(5)
	val, err := m.ToBool()
	assert.NoError(t, err)
	assert.True(t, val)
}

func TestToBoolWithStringFalsePorted(t *testing.T) {
	m := maybeForConversion("false")
	val, err := m.ToBool()
	assert.NoError(t, err)
	assert.False(t, val)
}

func TestToInt16WithFloat32OverflowPorted(t *testing.T) {
	val, err := maybeForConversion(float32(70000)).ToInt16()
	assert.Equal(t, int16(0), val)
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToInt16WithFloat64OverflowPorted(t *testing.T) {
	val, err := maybeForConversion(float64(70000)).ToInt16()
	assert.Equal(t, int16(0), val)
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToInt16WithInt32OverflowPorted(t *testing.T) {
	val, err := maybeForConversion(int32(70000)).ToInt16()
	assert.Equal(t, int16(0), val)
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToInt16WithInt64OverflowPorted(t *testing.T) {
	val, err := maybeForConversion(int64(70000)).ToInt16()
	assert.Equal(t, int16(0), val)
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToInt32WithBoolFalsePorted(t *testing.T) {
	m := maybeForConversion(false)
	val, err := m.ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(0), val)
}

func TestToInt32WithBytePorted(t *testing.T) {
	m := maybeForConversion(byte(200))
	val, err := m.ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(200), val)
}

func TestToInt32WithFloat32Ported(t *testing.T) {
	m := maybeForConversion(float32(42.5))
	val, err := m.ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(43), val)
}

func TestToInt32WithIntPorted(t *testing.T) {
	m := maybeForConversion(42)
	val, err := m.ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(42), val)
}

func TestToInt32WithInt16Ported(t *testing.T) {
	m := maybeForConversion(int16(42))
	val, err := m.ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(42), val)
}

func TestToInt32WithInt8Ported(t *testing.T) {
	m := maybeForConversion(int8(42))
	val, err := m.ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(42), val)
}

func TestToInt32WithIntOverflowPorted(t *testing.T) {
	if strconv.IntSize == 64 {
		val, err := maybeForConversion(int(math.MaxInt32) + 1).ToInt32()
		assert.Equal(t, int32(0), val)
		assert.Equal(t, ErrConversionSizeOverflow, err)
	}
}

func TestToInt32WithStringPorted(t *testing.T) {
	m := maybeForConversion("42")
	val, err := m.ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(42), val)
}

func TestToInt32WithUintPorted(t *testing.T) {
	m := maybeForConversion(uint(42))
	val, err := m.ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(42), val)
}

func TestToInt64WithBoolFalsePorted(t *testing.T) {
	m := maybeForConversion(false)
	val, err := m.ToInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), val)
}

func TestToInt64WithBoolTruePorted(t *testing.T) {
	val, err := maybeForConversion(true).ToInt64()
	assert.Equal(t, int64(1), val)
	assert.NoError(t, err)
}

func TestToInt64WithBytePorted(t *testing.T) {
	m := maybeForConversion(byte(42))
	val, err := m.ToInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(42), val)
}

func TestToInt64WithFloat32OverflowLargePorted(t *testing.T) {
	val, err := maybeForConversion(float32(math.MaxUint64)).ToInt64()
	assert.Equal(t, int64(0), val)
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToInt64WithFloat32SuccessPorted(t *testing.T) {
	m := maybeForConversion(float32(42.5))
	val, err := m.ToInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(43), val)
}

func TestToInt64WithFloat64OverflowPorted(t *testing.T) {
	// float64(math.MaxUint64) = 2^64 which exceeds math.MaxInt64 (2^63-1), triggering overflow.
	// Using math.MaxUint64 because float64(math.MaxInt64) rounds to the same float64 as
	// float64(math.MaxInt64) in the comparison (both become 2^63), so it would NOT overflow.
	val, err := maybeForConversion(float64(math.MaxUint64)).ToInt64()
	assert.Equal(t, int64(0), val)
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToInt64WithFloat64SuccessPorted(t *testing.T) {
	m := maybeForConversion(float64(42.5))
	val, err := m.ToInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(43), val)
}

func TestToInt64WithInt16Ported(t *testing.T) {
	m := maybeForConversion(int16(42))
	val, err := m.ToInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(42), val)
}

func TestToInt64WithInt32Ported(t *testing.T) {
	m := maybeForConversion(int32(42))
	val, err := m.ToInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(42), val)
}

func TestToInt64WithInt8Ported(t *testing.T) {
	m := maybeForConversion(int8(42))
	val, err := m.ToInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(42), val)
}

func TestToInt64WithUint16Ported(t *testing.T) {
	m := maybeForConversion(uint16(42))
	val, err := m.ToInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(42), val)
}

func TestToInt64WithUint32Ported(t *testing.T) {
	m := maybeForConversion(uint32(42))
	val, err := m.ToInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(42), val)
}

func TestToInt64WithUint64SuccessPorted(t *testing.T) {
	val, err := maybeForConversion(uint64(42)).ToInt64()
	assert.Equal(t, int64(42), val)
	assert.NoError(t, err)
}

func TestToInt64WithUintOverflowPorted(t *testing.T) {
	if math.MaxUint > math.MaxInt64 {
		m := maybeForConversion(uint(math.MaxUint))
		_, err := m.ToInt64()
		assert.Equal(t, ErrConversionSizeOverflow, err)
	}
}

func TestToInt64WithUintSuccessPorted(t *testing.T) {
	val, err := maybeForConversion(uint(42)).ToInt64()
	assert.Equal(t, int64(42), val)
	assert.NoError(t, err)
}

func TestToInt64WithUintptrOverflowLargePorted(t *testing.T) {
	val, err := maybeForConversion(uintptr(^uintptr(0))).ToInt64()
	if uint64(^uintptr(0)) > math.MaxInt64 {
		assert.Equal(t, int64(0), val)
		assert.Equal(t, ErrConversionSizeOverflow, err)
	} else {
		assert.NoError(t, err)
	}
}

func TestToInt8WithByteOverflowPorted(t *testing.T) {
	m := maybeForConversion(byte(200))
	_, err := m.ToInt8()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToInt8WithFloat32OverflowPorted(t *testing.T) {
	val, err := maybeForConversion(float32(200)).ToInt8()
	assert.Equal(t, int8(0), val)
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToInt8WithFloat64Ported(t *testing.T) {
	m := maybeForConversion(float64(42.5))
	val, err := m.ToInt8()
	assert.NoError(t, err)
	assert.Equal(t, int8(43), val)
}

func TestToInt8WithFloat64OverflowPorted(t *testing.T) {
	val, err := maybeForConversion(float64(200)).ToInt8()
	assert.Equal(t, int8(0), val)
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToInt8WithInt16OverflowPorted(t *testing.T) {
	val, err := maybeForConversion(int16(200)).ToInt8()
	assert.Equal(t, int8(0), val)
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToInt8WithInt32OverflowPorted(t *testing.T) {
	val, err := maybeForConversion(int32(200)).ToInt8()
	assert.Equal(t, int8(0), val)
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToInt8WithInt32SuccessPorted(t *testing.T) {
	m := maybeForConversion(int32(42))
	val, err := m.ToInt8()
	assert.NoError(t, err)
	assert.Equal(t, int8(42), val)
}

func TestToInt8WithInt64OverflowPorted(t *testing.T) {
	val, err := maybeForConversion(int64(200)).ToInt8()
	assert.Equal(t, int8(0), val)
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToInt8WithUint16Ported(t *testing.T) {
	m := maybeForConversion(uint16(42))
	val, err := m.ToInt8()
	assert.NoError(t, err)
	assert.Equal(t, int8(42), val)
}

func TestToInt8WithUint32Ported(t *testing.T) {
	m := maybeForConversion(uint32(42))
	val, err := m.ToInt8()
	assert.NoError(t, err)
	assert.Equal(t, int8(42), val)
}

func TestToIntWithBoolFalsePorted(t *testing.T) {
	m := maybeForConversion(false)
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 0, val)
}

func TestToIntWithBoolTruePorted(t *testing.T) {
	m := maybeForConversion(true)
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
}

func TestToIntWithBytePorted(t *testing.T) {
	m := maybeForConversion(byte(42))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestToIntWithFloat32OverflowPorted(t *testing.T) {
	val, err := maybeForConversion(float32(math.MaxInt32 * 2)).ToInt()
	assert.Equal(t, 0, val)
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToIntWithFloat32UnderflowPorted(t *testing.T) {
	val, err := maybeForConversion(float32(math.MinInt32) * 2).ToInt()
	assert.Equal(t, 0, val)
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToIntWithFloat64OverflowPorted(t *testing.T) {
	val, err := maybeForConversion(float64(math.MaxInt32 * 2)).ToInt()
	assert.Equal(t, 0, val)
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToIntWithFloat64UnderflowPorted(t *testing.T) {
	val, err := maybeForConversion(float64(math.MinInt32) * 2).ToInt()
	assert.Equal(t, 0, val)
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToIntWithInt16Ported(t *testing.T) {
	m := maybeForConversion(int16(42))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestToIntWithInt32Ported(t *testing.T) {
	m := maybeForConversion(int32(42))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestToIntWithInt32OverflowPorted(t *testing.T) {
	// int32 always fits in int on all Go platforms, so this should succeed
	m := maybeForConversion(int32(math.MinInt32))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, int(math.MinInt32), val)
}

func TestToIntWithInt64OverflowPorted(t *testing.T) {
	val, err := maybeForConversion(int64(math.MaxInt32) + 1).ToInt()
	assert.Equal(t, 0, val)
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToIntWithInt64UnderflowPorted(t *testing.T) {
	val, err := maybeForConversion(int64(math.MinInt32) - 1).ToInt()
	assert.Equal(t, 0, val)
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToIntWithInt8Ported(t *testing.T) {
	m := maybeForConversion(int8(42))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestToIntWithStringValuePorted(t *testing.T) {
	m := maybeForConversion("42")
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestToIntWithUint16Ported(t *testing.T) {
	m := maybeForConversion(uint16(42))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestToIntWithUint32Ported(t *testing.T) {
	m := maybeForConversion(uint32(42))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestToIntWithUint32OverflowPorted(t *testing.T) {
	val, err := maybeForConversion(uint32(math.MaxInt32 + 1)).ToInt()
	assert.Equal(t, 0, val)
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToIntWithUintOverflowPorted(t *testing.T) {
	val, err := maybeForConversion(uint(math.MaxInt32 + 1)).ToInt()
	assert.Equal(t, 0, val)
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestToIntWithUintptrPorted(t *testing.T) {
	m := maybeForConversion(uintptr(42))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestSomeDefToFloat64WithNilPorted(t *testing.T) {
	m := maybeForConversion(nil)
	_, err := m.ToFloat64()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToFloat64WithUnsupportedTypePorted(t *testing.T) {
	m := maybeForConversion(make(chan int))
	_, err := m.ToFloat64()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToFloat64WithUintPorted(t *testing.T) {
	m := maybeForConversion(uint(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

func TestSomeDefToFloat64WithUint16Ported(t *testing.T) {
	m := maybeForConversion(uint16(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

func TestSomeDefToFloat64WithUint32Ported(t *testing.T) {
	m := maybeForConversion(uint32(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

func TestSomeDefToFloat64WithUint64Ported(t *testing.T) {
	m := maybeForConversion(uint64(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

func TestSomeDefToFloat64WithUintptrPorted(t *testing.T) {
	m := maybeForConversion(uintptr(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

func TestSomeDefToFloat64WithBytePorted(t *testing.T) {
	m := maybeForConversion(byte(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

func TestSomeDefToFloat64WithInt8Ported(t *testing.T) {
	m := maybeForConversion(int8(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

func TestSomeDefToFloat64WithInt16Ported(t *testing.T) {
	m := maybeForConversion(int16(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

func TestSomeDefToFloat64WithInt32Ported(t *testing.T) {
	m := maybeForConversion(int32(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

func TestSomeDefToFloat64WithInt64Ported(t *testing.T) {
	m := maybeForConversion(int64(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

func TestSomeDefToFloat32WithNilPorted(t *testing.T) {
	m := maybeForConversion(nil)
	_, err := m.ToFloat32()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToFloat32WithUnsupportedTypePorted(t *testing.T) {
	m := maybeForConversion(make(chan int))
	_, err := m.ToFloat32()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToFloat32WithBoolTruePorted(t *testing.T) {
	m := maybeForConversion(true)
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(1), val)
}

func TestSomeDefToFloat32WithBoolFalsePorted(t *testing.T) {
	m := maybeForConversion(false)
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(0), val)
}

func TestSomeDefToFloat32WithUintPorted(t *testing.T) {
	m := maybeForConversion(uint(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithUint16Ported(t *testing.T) {
	m := maybeForConversion(uint16(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithUint32Ported(t *testing.T) {
	m := maybeForConversion(uint32(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithUint64Ported(t *testing.T) {
	m := maybeForConversion(uint64(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithUintptrPorted(t *testing.T) {
	m := maybeForConversion(uintptr(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithBytePorted(t *testing.T) {
	m := maybeForConversion(byte(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithIntPorted(t *testing.T) {
	m := maybeForConversion(42)
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithInt8Ported(t *testing.T) {
	m := maybeForConversion(int8(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithInt16Ported(t *testing.T) {
	m := maybeForConversion(int16(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithInt32Ported(t *testing.T) {
	m := maybeForConversion(int32(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithInt64Ported(t *testing.T) {
	m := maybeForConversion(int64(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithFloat32Ported(t *testing.T) {
	m := maybeForConversion(float32(42.5))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42.5), val)
}

func TestSomeDefToFloat32WithFloat64Ported(t *testing.T) {
	m := maybeForConversion(float64(42.5))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42.5), val)
}

func TestSomeDefToFloat64WithBoolTruePorted(t *testing.T) {
	m := maybeForConversion(true)
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(1), val)
}

func TestSomeDefToFloat64WithBoolFalsePorted(t *testing.T) {
	m := maybeForConversion(false)
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(0), val)
}
