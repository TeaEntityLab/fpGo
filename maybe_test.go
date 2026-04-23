package fpgo

import (
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPresent(t *testing.T) {
	var m MaybeDef[interface{}]

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
	var m MaybeDef[interface{}]

	m = Maybe.Just(1)
	assert.Equal(t, 1, m.Or(3))
	m = Maybe.Just(nil)
	assert.Equal(t, 3, m.Or(3))
}

func TestClone(t *testing.T) {
	var m MaybeDef[interface{}]

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
	assert.Equal(t, nil, m.Clone().UnwrapInterface())
	assert.Equal(t, true, m.IsNil())
	assert.Equal(t, 2, *iptr2)

	iptr = &i
	m = Maybe.Just(iptr)
	assert.Equal(t, 1, *CloneTo(m, interface{}(iptr2)).ToPtr())
	assert.Equal(t, 1, *iptr2)
}

func TestFlatMap(t *testing.T) {
	var m MaybeDef[interface{}]

	m = Maybe.Just(1).FlatMap(func(in interface{}) MaybeDef[interface{}] {
		v, _ := Maybe.Just(in).ToInt()
		result := Maybe.Just(v + 1)
		return result
	})
	assert.Equal(t, 2, m.Unwrap())
}

func TestLet(t *testing.T) {
	var m MaybeDef[interface{}]

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
	var m MaybeDef[interface{}]

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
	var m MaybeDef[interface{}]

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

func TestMaybeToPtr(t *testing.T) {
	m := Maybe.Just(42)

	ptr := m.UnwrapInterface()
	assert.Equal(t, 42, ptr)

	mNil := Maybe.Just(nil)
	assert.Nil(t, mNil.ToPtr())
}

func TestMaybeToMaybe(t *testing.T) {
	m := Maybe.Just(1)
	m2 := m.ToMaybe()

	assert.Equal(t, 1, m2.Unwrap())

	mNil := Maybe.Just(nil)
	m2Nil := mNil.ToMaybe()
	assert.Equal(t, true, m2Nil.IsNil())
}

func TestMaybeUnwrapInterface(t *testing.T) {
	m := Maybe.Just(42)
	assert.Equal(t, 42, m.UnwrapInterface())

	mNil := Maybe.Just(nil)
	assert.Nil(t, mNil.UnwrapInterface())
}

func TestMaybeUnwrap(t *testing.T) {
	m := Maybe.Just(42)
	assert.Equal(t, 42, m.Unwrap())

	mNil := Maybe.Just(nil)
	assert.Nil(t, mNil.UnwrapInterface())
}

func TestMaybeIsValid(t *testing.T) {
	m := Maybe.Just(42)
	assert.Equal(t, true, m.IsValid())

	mNil := Maybe.Just(nil)
	assert.Equal(t, false, mNil.IsValid())
}

func TestMaybeIsPtr(t *testing.T) {
	i := 42
	m := Maybe.Just(&i)
	assert.Equal(t, true, m.IsPtr())

	m = Maybe.Just(42)
	assert.Equal(t, false, m.IsPtr())
}

func TestMaybeClone(t *testing.T) {
	m := Maybe.Just(42)
	cloned := m.Clone()
	assert.Equal(t, 42, cloned.Unwrap())
}

func TestMaybeJustWithNilValue(t *testing.T) {
	m := Maybe.Just(nil)
	assert.True(t, m.IsNil())
}

func TestMaybeOrWithNil(t *testing.T) {
	m := Maybe.Just(nil)
	result := m.Or(42)
	assert.Equal(t, 42, result)
}

func TestMaybeToStringString(t *testing.T) {
	m := Maybe.Just("hello")
	assert.Equal(t, "hello", m.ToString())
}

func TestMaybeToStringInt(t *testing.T) {
	m := Maybe.Just(123)
	assert.Equal(t, "123", m.ToString())
}

func TestMaybeToPtrWithValue(t *testing.T) {
	val := 42
	m := Maybe.Just(&val)
	ptr := m.ToPtr()
	assert.NotNil(t, ptr)
	assert.Equal(t, 42, *ptr)
}

func TestMaybeFlatMapNil(t *testing.T) {
	m := Maybe.Just(nil).FlatMap(func(in interface{}) MaybeDef[interface{}] {
		return Maybe.Just(in)
	})
	assert.True(t, m.IsNil())
}

func TestMaybeToMaybeInterface(t *testing.T) {
	m := Maybe.Just(42)
	m2 := m.ToMaybe()
	assert.Equal(t, 42, m2.Unwrap())

	mNil := Maybe.Just(nil)
	m2Nil := mNil.ToMaybe()
	assert.True(t, m2Nil.IsNil())
}

func TestMaybeLetWithValue(t *testing.T) {
	var called bool
	m := Maybe.Just(42)
	m.Let(func() {
		called = true
	})
	assert.True(t, called)
}

func TestMaybeLetWithNil(t *testing.T) {
	var called bool
	m := Maybe.Just(nil)
	m.Let(func() {
		called = true
	})
	assert.False(t, called)
}

func TestMaybeNoneLet(t *testing.T) {
	var called bool
	var m MaybeDef[interface{}]
	m = None
	m.Let(func() {
		called = true
	})
	assert.False(t, called)
}

func TestJustGenerics(t *testing.T) {
	m := JustGenerics(42)
	assert.Equal(t, 42, m.Unwrap())

	mStr := JustGenerics("hello")
	assert.Equal(t, "hello", mStr.Unwrap())
}

func TestMaybeUnwrapInterfaceAlt(t *testing.T) {
	m := Maybe.Just(42)
	val := m.UnwrapInterface()
	assert.Equal(t, 42, val)

	mNil := Maybe.Just(nil)
	valNil := mNil.UnwrapInterface()
	assert.Nil(t, valNil)
}

func TestMaybeIsKind(t *testing.T) {
	m := Maybe.Just(42)
	assert.True(t, m.IsKind(reflect.Int))
	assert.False(t, m.IsKind(reflect.String))
}

func TestMaybeIsType(t *testing.T) {
	m := Maybe.Just(42)
	assert.True(t, m.IsType(reflect.TypeOf(42)))
	assert.False(t, m.IsType(reflect.TypeOf("")))
}

func TestToByte(t *testing.T) {
	// From int
	m := Maybe.Just(42)
	b, err := m.ToByte()
	assert.Equal(t, byte(42), b)
	assert.Nil(t, err)

	// From int8
	m = Maybe.Just(int8(10))
	b, err = m.ToByte()
	assert.Equal(t, byte(10), b)
	assert.Nil(t, err)

	// From int16
	m = Maybe.Just(int16(100))
	b, err = m.ToByte()
	assert.Equal(t, byte(100), b)
	assert.Nil(t, err)

	// From int32
	m = Maybe.Just(int32(50))
	b, err = m.ToByte()
	assert.Equal(t, byte(50), b)
	assert.Nil(t, err)

	// From int64
	m = Maybe.Just(int64(25))
	b, err = m.ToByte()
	assert.Equal(t, byte(25), b)
	assert.Nil(t, err)

	// From uint
	m = Maybe.Just(uint(200))
	b, err = m.ToByte()
	assert.Equal(t, byte(200), b)
	assert.Nil(t, err)

	// From uint16
	m = Maybe.Just(uint16(150))
	b, err = m.ToByte()
	assert.Equal(t, byte(150), b)
	assert.Nil(t, err)

	// From uint32
	m = Maybe.Just(uint32(75))
	b, err = m.ToByte()
	assert.Equal(t, byte(75), b)
	assert.Nil(t, err)

	// From uint64
	m = Maybe.Just(uint64(125))
	b, err = m.ToByte()
	assert.Equal(t, byte(125), b)
	assert.Nil(t, err)

	// From uintptr
	m = Maybe.Just(uintptr(64))
	b, err = m.ToByte()
	assert.Equal(t, byte(64), b)
	assert.Nil(t, err)

	// From float32
	m = Maybe.Just(float32(3.7))
	b, err = m.ToByte()
	assert.Equal(t, byte(4), b) // rounds to 4
	assert.Nil(t, err)

	// From float64
	m = Maybe.Just(float64(2.3))
	b, err = m.ToByte()
	assert.Equal(t, byte(2), b) // rounds to 2
	assert.Nil(t, err)

	// From bool true
	m = Maybe.Just(true)
	b, err = m.ToByte()
	assert.Equal(t, byte(1), b)
	assert.Nil(t, err)

	// From bool false
	m = Maybe.Just(false)
	b, err = m.ToByte()
	assert.Equal(t, byte(0), b)
	assert.Nil(t, err)

	// From string
	m = Maybe.Just("99")
	b, err = m.ToByte()
	assert.Equal(t, byte(99), b)
	assert.Nil(t, err)

	// From byte
	m = Maybe.Just(byte(255))
	b, err = m.ToByte()
	assert.Equal(t, byte(255), b)
	assert.Nil(t, err)

	// From nil
	m = Maybe.Just(nil)
	b, err = m.ToByte()
	assert.Equal(t, byte(0), b)
	assert.Equal(t, ErrConversionNil, err)

	// Overflow from int
	m = Maybe.Just(256)
	b, err = m.ToByte()
	assert.Equal(t, byte(0), b)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from int (negative)
	m = Maybe.Just(-1)
	b, err = m.ToByte()
	assert.Equal(t, byte(0), b)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint
	m = Maybe.Just(uint(300))
	b, err = m.ToByte()
	assert.Equal(t, byte(0), b)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Unsupported type
	m = Maybe.Just([]int{1, 2, 3})
	b, err = m.ToByte()
	assert.Equal(t, byte(0), b)
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestToInt8(t *testing.T) {
	// From int
	m := Maybe.Just(42)
	i, err := m.ToInt8()
	assert.Equal(t, int8(42), i)
	assert.Nil(t, err)

	// From int8
	m = Maybe.Just(int8(10))
	i, err = m.ToInt8()
	assert.Equal(t, int8(10), i)
	assert.Nil(t, err)

	// From int16
	m = Maybe.Just(int16(100))
	i, err = m.ToInt8()
	assert.Equal(t, int8(100), i)
	assert.Nil(t, err)

	// From int32
	m = Maybe.Just(int32(50))
	i, err = m.ToInt8()
	assert.Equal(t, int8(50), i)
	assert.Nil(t, err)

	// From int64
	m = Maybe.Just(int64(25))
	i, err = m.ToInt8()
	assert.Equal(t, int8(25), i)
	assert.Nil(t, err)

	// From uint
	m = Maybe.Just(uint(100))
	i, err = m.ToInt8()
	assert.Equal(t, int8(100), i)
	assert.Nil(t, err)

	// From uint16
	m = Maybe.Just(uint16(50))
	i, err = m.ToInt8()
	assert.Equal(t, int8(50), i)
	assert.Nil(t, err)

	// From uint32
	m = Maybe.Just(uint32(75))
	i, err = m.ToInt8()
	assert.Equal(t, int8(75), i)
	assert.Nil(t, err)

	// From uint64
	m = Maybe.Just(uint64(125))
	i, err = m.ToInt8()
	assert.Equal(t, int8(125), i)
	assert.Nil(t, err)

	// From uintptr
	m = Maybe.Just(uintptr(64))
	i, err = m.ToInt8()
	assert.Equal(t, int8(64), i)
	assert.Nil(t, err)

	// From float32
	m = Maybe.Just(float32(3.7))
	i, err = m.ToInt8()
	assert.Equal(t, int8(4), i) // rounds to 4
	assert.Nil(t, err)

	// From float64
	m = Maybe.Just(float64(2.3))
	i, err = m.ToInt8()
	assert.Equal(t, int8(2), i) // rounds to 2
	assert.Nil(t, err)

	// From bool true
	m = Maybe.Just(true)
	i, err = m.ToInt8()
	assert.Equal(t, int8(1), i)
	assert.Nil(t, err)

	// From bool false
	m = Maybe.Just(false)
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Nil(t, err)

	// From string
	m = Maybe.Just("99")
	i, err = m.ToInt8()
	assert.Equal(t, int8(99), i)
	assert.Nil(t, err)

	// From byte
	m = Maybe.Just(byte(100))
	i, err = m.ToInt8()
	assert.Equal(t, int8(100), i)
	assert.Nil(t, err)

	// From nil
	m = Maybe.Just(nil)
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Equal(t, ErrConversionNil, err)

	// Overflow from int (positive)
	m = Maybe.Just(128)
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from int (negative)
	m = Maybe.Just(-129)
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint
	m = Maybe.Just(uint(200))
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint16
	m = Maybe.Just(uint16(200))
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint32
	m = Maybe.Just(uint32(200))
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint64
	m = Maybe.Just(uint64(200))
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uintptr
	m = Maybe.Just(uintptr(200))
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Unsupported type
	m = Maybe.Just([]int{1, 2, 3})
	i, err = m.ToInt8()
	assert.Equal(t, int8(0), i)
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestToInt16(t *testing.T) {
	// From int
	m := Maybe.Just(1000)
	i, err := m.ToInt16()
	assert.Equal(t, int16(1000), i)
	assert.Nil(t, err)

	// From int8
	m = Maybe.Just(int8(10))
	i, err = m.ToInt16()
	assert.Equal(t, int16(10), i)
	assert.Nil(t, err)

	// From int16
	m = Maybe.Just(int16(500))
	i, err = m.ToInt16()
	assert.Equal(t, int16(500), i)
	assert.Nil(t, err)

	// From int32
	m = Maybe.Just(int32(250))
	i, err = m.ToInt16()
	assert.Equal(t, int16(250), i)
	assert.Nil(t, err)

	// From int64
	m = Maybe.Just(int64(125))
	i, err = m.ToInt16()
	assert.Equal(t, int16(125), i)
	assert.Nil(t, err)

	// From uint
	m = Maybe.Just(uint(1000))
	i, err = m.ToInt16()
	assert.Equal(t, int16(1000), i)
	assert.Nil(t, err)

	// From uint16
	m = Maybe.Just(uint16(500))
	i, err = m.ToInt16()
	assert.Equal(t, int16(500), i)
	assert.Nil(t, err)

	// From uint32
	m = Maybe.Just(uint32(250))
	i, err = m.ToInt16()
	assert.Equal(t, int16(250), i)
	assert.Nil(t, err)

	// From uint64
	m = Maybe.Just(uint64(125))
	i, err = m.ToInt16()
	assert.Equal(t, int16(125), i)
	assert.Nil(t, err)

	// From uintptr
	m = Maybe.Just(uintptr(64))
	i, err = m.ToInt16()
	assert.Equal(t, int16(64), i)
	assert.Nil(t, err)

	// From float32
	m = Maybe.Just(float32(3.7))
	i, err = m.ToInt16()
	assert.Equal(t, int16(4), i) // rounds to 4
	assert.Nil(t, err)

	// From float64
	m = Maybe.Just(float64(2.3))
	i, err = m.ToInt16()
	assert.Equal(t, int16(2), i) // rounds to 2
	assert.Nil(t, err)

	// From bool true
	m = Maybe.Just(true)
	i, err = m.ToInt16()
	assert.Equal(t, int16(1), i)
	assert.Nil(t, err)

	// From bool false
	m = Maybe.Just(false)
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Nil(t, err)

	// From string
	m = Maybe.Just("999")
	i, err = m.ToInt16()
	assert.Equal(t, int16(999), i)
	assert.Nil(t, err)

	// From byte
	m = Maybe.Just(byte(100))
	i, err = m.ToInt16()
	assert.Equal(t, int16(100), i)
	assert.Nil(t, err)

	// From nil
	m = Maybe.Just(nil)
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Equal(t, ErrConversionNil, err)

	// Overflow from int (positive)
	m = Maybe.Just(32768) // MaxInt16 + 1
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from int (negative)
	m = Maybe.Just(-32769) // MinInt16 - 1
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint
	m = Maybe.Just(uint(40000))
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint16
	m = Maybe.Just(uint16(40000))
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint32
	m = Maybe.Just(uint32(40000))
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint64
	m = Maybe.Just(uint64(40000))
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uintptr
	m = Maybe.Just(uintptr(40000))
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Unsupported type
	m = Maybe.Just([]int{1, 2, 3})
	i, err = m.ToInt16()
	assert.Equal(t, int16(0), i)
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestToUint(t *testing.T) {
	// From int
	m := Maybe.Just(42)
	u, err := m.ToUint()
	assert.Equal(t, uint(42), u)
	assert.Nil(t, err)

	// From int8
	m = Maybe.Just(int8(10))
	u, err = m.ToUint()
	assert.Equal(t, uint(10), u)
	assert.Nil(t, err)

	// From int16
	m = Maybe.Just(int16(100))
	u, err = m.ToUint()
	assert.Equal(t, uint(100), u)
	assert.Nil(t, err)

	// From int32
	m = Maybe.Just(int32(50))
	u, err = m.ToUint()
	assert.Equal(t, uint(50), u)
	assert.Nil(t, err)

	// From int64
	m = Maybe.Just(int64(25))
	u, err = m.ToUint()
	assert.Equal(t, uint(25), u)
	assert.Nil(t, err)

	// From uint
	m = Maybe.Just(uint(200))
	u, err = m.ToUint()
	assert.Equal(t, uint(200), u)
	assert.Nil(t, err)

	// From uint16
	m = Maybe.Just(uint16(150))
	u, err = m.ToUint()
	assert.Equal(t, uint(150), u)
	assert.Nil(t, err)

	// From uint32
	m = Maybe.Just(uint32(75))
	u, err = m.ToUint()
	assert.Equal(t, uint(75), u)
	assert.Nil(t, err)

	// From uint64
	m = Maybe.Just(uint64(125))
	u, err = m.ToUint()
	assert.Equal(t, uint(125), u)
	assert.Nil(t, err)

	// From uintptr
	m = Maybe.Just(uintptr(64))
	u, err = m.ToUint()
	assert.Equal(t, uint(64), u)
	assert.Nil(t, err)

	// From float32
	m = Maybe.Just(float32(3.7))
	u, err = m.ToUint()
	assert.Equal(t, uint(4), u) // rounds to 4
	assert.Nil(t, err)

	// From float64
	m = Maybe.Just(float64(2.3))
	u, err = m.ToUint()
	assert.Equal(t, uint(2), u) // rounds to 2
	assert.Nil(t, err)

	// From bool true
	m = Maybe.Just(true)
	u, err = m.ToUint()
	assert.Equal(t, uint(1), u)
	assert.Nil(t, err)

	// From bool false
	m = Maybe.Just(false)
	u, err = m.ToUint()
	assert.Equal(t, uint(0), u)
	assert.Nil(t, err)

	// From string
	m = Maybe.Just("99")
	u, err = m.ToUint()
	assert.Equal(t, uint(99), u)
	assert.Nil(t, err)

	// From byte
	m = Maybe.Just(byte(255))
	u, err = m.ToUint()
	assert.Equal(t, uint(255), u)
	assert.Nil(t, err)

	// From nil
	m = Maybe.Just(nil)
	u, err = m.ToUint()
	assert.Equal(t, uint(0), u)
	assert.Equal(t, ErrConversionNil, err)

	// Unsupported type
	m = Maybe.Just([]int{1, 2, 3})
	u, err = m.ToUint()
	assert.Equal(t, uint(0), u)
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestToUint8(t *testing.T) {
	// From int
	m := Maybe.Just(42)
	u, err := m.ToUint8()
	assert.Equal(t, uint8(42), u)
	assert.Nil(t, err)

	// From int8
	m = Maybe.Just(int8(10))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(10), u)
	assert.Nil(t, err)

	// From int16
	m = Maybe.Just(int16(100))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(100), u)
	assert.Nil(t, err)

	// From int32
	m = Maybe.Just(int32(50))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(50), u)
	assert.Nil(t, err)

	// From int64
	m = Maybe.Just(int64(25))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(25), u)
	assert.Nil(t, err)

	// From uint
	m = Maybe.Just(uint(200))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(200), u)
	assert.Nil(t, err)

	// From uint16
	m = Maybe.Just(uint16(150))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(150), u)
	assert.Nil(t, err)

	// From uint32
	m = Maybe.Just(uint32(75))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(75), u)
	assert.Nil(t, err)

	// From uint64
	m = Maybe.Just(uint64(125))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(125), u)
	assert.Nil(t, err)

	// From uintptr
	m = Maybe.Just(uintptr(64))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(64), u)
	assert.Nil(t, err)

	// From float32
	m = Maybe.Just(float32(3.7))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(4), u) // rounds to 4
	assert.Nil(t, err)

	// From float64
	m = Maybe.Just(float64(2.3))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(2), u) // rounds to 2
	assert.Nil(t, err)

	// From bool true
	m = Maybe.Just(true)
	u, err = m.ToUint8()
	assert.Equal(t, uint8(1), u)
	assert.Nil(t, err)

	// From bool false
	m = Maybe.Just(false)
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Nil(t, err)

	// From string
	m = Maybe.Just("99")
	u, err = m.ToUint8()
	assert.Equal(t, uint8(99), u)
	assert.Nil(t, err)

	// From byte
	m = Maybe.Just(byte(255))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(255), u)
	assert.Nil(t, err)

	// From nil
	m = Maybe.Just(nil)
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Equal(t, ErrConversionNil, err)

	// Overflow from int
	m = Maybe.Just(256)
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from int (negative)
	m = Maybe.Just(-1)
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint
	m = Maybe.Just(uint(300))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint16
	m = Maybe.Just(uint16(300))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint32
	m = Maybe.Just(uint32(300))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint64
	m = Maybe.Just(uint64(300))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uintptr
	m = Maybe.Just(uintptr(300))
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Unsupported type
	m = Maybe.Just([]int{1, 2, 3})
	u, err = m.ToUint8()
	assert.Equal(t, uint8(0), u)
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestToUint16(t *testing.T) {
	// From int
	m := Maybe.Just(1000)
	u, err := m.ToUint16()
	assert.Equal(t, uint16(1000), u)
	assert.Nil(t, err)

	// From int8
	m = Maybe.Just(int8(10))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(10), u)
	assert.Nil(t, err)

	// From int16
	m = Maybe.Just(int16(500))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(500), u)
	assert.Nil(t, err)

	// From int32
	m = Maybe.Just(int32(250))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(250), u)
	assert.Nil(t, err)

	// From int64
	m = Maybe.Just(int64(125))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(125), u)
	assert.Nil(t, err)

	// From uint
	m = Maybe.Just(uint(1000))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(1000), u)
	assert.Nil(t, err)

	// From uint16
	m = Maybe.Just(uint16(500))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(500), u)
	assert.Nil(t, err)

	// From uint32
	m = Maybe.Just(uint32(250))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(250), u)
	assert.Nil(t, err)

	// From uint64
	m = Maybe.Just(uint64(125))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(125), u)
	assert.Nil(t, err)

	// From uintptr
	m = Maybe.Just(uintptr(64))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(64), u)
	assert.Nil(t, err)

	// From float32
	m = Maybe.Just(float32(3.7))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(4), u) // rounds to 4
	assert.Nil(t, err)

	// From float64
	m = Maybe.Just(float64(2.3))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(2), u) // rounds to 2
	assert.Nil(t, err)

	// From bool true
	m = Maybe.Just(true)
	u, err = m.ToUint16()
	assert.Equal(t, uint16(1), u)
	assert.Nil(t, err)

	// From bool false
	m = Maybe.Just(false)
	u, err = m.ToUint16()
	assert.Equal(t, uint16(0), u)
	assert.Nil(t, err)

	// From string
	m = Maybe.Just("999")
	u, err = m.ToUint16()
	assert.Equal(t, uint16(999), u)
	assert.Nil(t, err)

	// From byte
	m = Maybe.Just(byte(100))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(100), u)
	assert.Nil(t, err)

	// From nil
	m = Maybe.Just(nil)
	u, err = m.ToUint16()
	assert.Equal(t, uint16(0), u)
	assert.Equal(t, ErrConversionNil, err)

	// Overflow from int
	m = Maybe.Just(65536) // MaxUint16 + 1
	u, err = m.ToUint16()
	assert.Equal(t, uint16(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from int (negative)
	m = Maybe.Just(-1)
	u, err = m.ToUint16()
	assert.Equal(t, uint16(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint
	m = Maybe.Just(uint(70000))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint32
	m = Maybe.Just(uint32(70000))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint64
	m = Maybe.Just(uint64(70000))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uintptr
	m = Maybe.Just(uintptr(70000))
	u, err = m.ToUint16()
	assert.Equal(t, uint16(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Unsupported type
	m = Maybe.Just([]int{1, 2, 3})
	u, err = m.ToUint16()
	assert.Equal(t, uint16(0), u)
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestToUint32(t *testing.T) {
	// From int
	m := Maybe.Just(1000)
	u, err := m.ToUint32()
	assert.Equal(t, uint32(1000), u)
	assert.Nil(t, err)

	// From int8
	m = Maybe.Just(int8(10))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(10), u)
	assert.Nil(t, err)

	// From int16
	m = Maybe.Just(int16(500))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(500), u)
	assert.Nil(t, err)

	// From int32
	m = Maybe.Just(int32(250))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(250), u)
	assert.Nil(t, err)

	// From int64
	m = Maybe.Just(int64(125))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(125), u)
	assert.Nil(t, err)

	// From uint
	m = Maybe.Just(uint(1000))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(1000), u)
	assert.Nil(t, err)

	// From uint16
	m = Maybe.Just(uint16(500))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(500), u)
	assert.Nil(t, err)

	// From uint32
	m = Maybe.Just(uint32(250))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(250), u)
	assert.Nil(t, err)

	// From uint64
	m = Maybe.Just(uint64(125))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(125), u)
	assert.Nil(t, err)

	// From uintptr
	m = Maybe.Just(uintptr(64))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(64), u)
	assert.Nil(t, err)

	// From float32
	m = Maybe.Just(float32(3.7))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(4), u) // rounds to 4
	assert.Nil(t, err)

	// From float64
	m = Maybe.Just(float64(2.3))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(2), u) // rounds to 2
	assert.Nil(t, err)

	// From bool true
	m = Maybe.Just(true)
	u, err = m.ToUint32()
	assert.Equal(t, uint32(1), u)
	assert.Nil(t, err)

	// From bool false
	m = Maybe.Just(false)
	u, err = m.ToUint32()
	assert.Equal(t, uint32(0), u)
	assert.Nil(t, err)

	// From string
	m = Maybe.Just("999")
	u, err = m.ToUint32()
	assert.Equal(t, uint32(999), u)
	assert.Nil(t, err)

	// From byte
	m = Maybe.Just(byte(100))
	u, err = m.ToUint32()
	assert.Equal(t, uint32(100), u)
	assert.Nil(t, err)

	// From nil
	m = Maybe.Just(nil)
	u, err = m.ToUint32()
	assert.Equal(t, uint32(0), u)
	assert.Equal(t, ErrConversionNil, err)

	// Overflow from int (negative)
	m = Maybe.Just(-1)
	u, err = m.ToUint32()
	assert.Equal(t, uint32(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from int64
	m = Maybe.Just(int64(4294967296)) // MaxUint32 + 1
	u, err = m.ToUint32()
	assert.Equal(t, uint32(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Overflow from uint64
	m = Maybe.Just(uint64(4294967296)) // MaxUint32 + 1
	u, err = m.ToUint32()
	assert.Equal(t, uint32(0), u)
	assert.Equal(t, ErrConversionSizeOverflow, err)

	// Unsupported type
	m = Maybe.Just([]int{1, 2, 3})
	u, err = m.ToUint32()
	assert.Equal(t, uint32(0), u)
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestToUint64(t *testing.T) {
	// From int
	m := Maybe.Just(1000)
	u, err := m.ToUint64()
	assert.Equal(t, uint64(1000), u)
	assert.Nil(t, err)

	// From int8
	m = Maybe.Just(int8(10))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(10), u)
	assert.Nil(t, err)

	// From int16
	m = Maybe.Just(int16(500))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(500), u)
	assert.Nil(t, err)

	// From int32
	m = Maybe.Just(int32(250))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(250), u)
	assert.Nil(t, err)

	// From int64
	m = Maybe.Just(int64(125))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(125), u)
	assert.Nil(t, err)

	// From uint
	m = Maybe.Just(uint(1000))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(1000), u)
	assert.Nil(t, err)

	// From uint16
	m = Maybe.Just(uint16(500))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(500), u)
	assert.Nil(t, err)

	// From uint32
	m = Maybe.Just(uint32(250))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(250), u)
	assert.Nil(t, err)

	// From uint64
	m = Maybe.Just(uint64(125))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(125), u)
	assert.Nil(t, err)

	// From uintptr
	m = Maybe.Just(uintptr(64))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(64), u)
	assert.Nil(t, err)

	// From float32
	m = Maybe.Just(float32(3.7))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(4), u) // rounds to 4
	assert.Nil(t, err)

	// From float64
	m = Maybe.Just(float64(2.3))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(2), u) // rounds to 2
	assert.Nil(t, err)

	// From bool true
	m = Maybe.Just(true)
	u, err = m.ToUint64()
	assert.Equal(t, uint64(1), u)
	assert.Nil(t, err)

	// From bool false
	m = Maybe.Just(false)
	u, err = m.ToUint64()
	assert.Equal(t, uint64(0), u)
	assert.Nil(t, err)

	// From string
	m = Maybe.Just("999")
	u, err = m.ToUint64()
	assert.Equal(t, uint64(999), u)
	assert.Nil(t, err)

	// From byte
	m = Maybe.Just(byte(100))
	u, err = m.ToUint64()
	assert.Equal(t, uint64(100), u)
	assert.Nil(t, err)

	// From nil
	m = Maybe.Just(nil)
	u, err = m.ToUint64()
	assert.Equal(t, uint64(0), u)
	assert.Equal(t, ErrConversionNil, err)

	// Unsupported type
	m = Maybe.Just([]int{1, 2, 3})
	u, err = m.ToUint64()
	assert.Equal(t, uint64(0), u)
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestToUintptr(t *testing.T) {
	// From int
	m := Maybe.Just(1000)
	u, err := m.ToUintptr()
	assert.Equal(t, uintptr(1000), u)
	assert.Nil(t, err)

	// From int8
	m = Maybe.Just(int8(10))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(10), u)
	assert.Nil(t, err)

	// From int16
	m = Maybe.Just(int16(500))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(500), u)
	assert.Nil(t, err)

	// From int32
	m = Maybe.Just(int32(250))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(250), u)
	assert.Nil(t, err)

	// From int64
	m = Maybe.Just(int64(125))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(125), u)
	assert.Nil(t, err)

	// From uint
	m = Maybe.Just(uint(1000))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(1000), u)
	assert.Nil(t, err)

	// From uint16
	m = Maybe.Just(uint16(500))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(500), u)
	assert.Nil(t, err)

	// From uint32
	m = Maybe.Just(uint32(250))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(250), u)
	assert.Nil(t, err)

	// From uint64
	m = Maybe.Just(uint64(125))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(125), u)
	assert.Nil(t, err)

	// From uintptr
	m = Maybe.Just(uintptr(64))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(64), u)
	assert.Nil(t, err)

	// From float32
	m = Maybe.Just(float32(3.7))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(4), u) // rounds to 4
	assert.Nil(t, err)

	// From float64
	m = Maybe.Just(float64(2.3))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(2), u) // rounds to 2
	assert.Nil(t, err)

	// From bool true
	m = Maybe.Just(true)
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(1), u)
	assert.Nil(t, err)

	// From bool false
	m = Maybe.Just(false)
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(0), u)
	assert.Nil(t, err)

	// From string
	m = Maybe.Just("999")
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(999), u)
	assert.Nil(t, err)

	// From byte
	m = Maybe.Just(byte(100))
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(100), u)
	assert.Nil(t, err)

	// From nil
	m = Maybe.Just(nil)
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(0), u)
	assert.Equal(t, ErrConversionNil, err)

	// Unsupported type
	m = Maybe.Just([]int{1, 2, 3})
	u, err = m.ToUintptr()
	assert.Equal(t, uintptr(0), u)
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestMaybeNoneConversions(t *testing.T) {
	m := None

	assert.False(t, m.IsPresent())
	assert.True(t, m.IsNil())

	_, err := m.ToFloat64()
	assert.Equal(t, ErrConversionNil, err)
	_, err = m.ToFloat32()
	assert.Equal(t, ErrConversionNil, err)
	_, err = m.ToInt()
	assert.Equal(t, ErrConversionNil, err)
	_, err = m.ToInt8()
	assert.Equal(t, ErrConversionNil, err)
	_, err = m.ToInt16()
	assert.Equal(t, ErrConversionNil, err)
	_, err = m.ToInt32()
	assert.Equal(t, ErrConversionNil, err)
	_, err = m.ToInt64()
	assert.Equal(t, ErrConversionNil, err)
	_, err = m.ToUint()
	assert.Equal(t, ErrConversionNil, err)
	_, err = m.ToUint8()
	assert.Equal(t, ErrConversionNil, err)
	_, err = m.ToUint16()
	assert.Equal(t, ErrConversionNil, err)
	_, err = m.ToUint32()
	assert.Equal(t, ErrConversionNil, err)
	_, err = m.ToUint64()
	assert.Equal(t, ErrConversionNil, err)
	_, err = m.ToUintptr()
	assert.Equal(t, ErrConversionNil, err)
	_, err = m.ToByte()
	assert.Equal(t, ErrConversionNil, err)
	_, err = m.ToBool()
	assert.Equal(t, ErrConversionNil, err)

	assert.Equal(t, "<nil>", m.ToString())
	assert.Nil(t, m.ToPtr())
	assert.False(t, m.IsValid())
	assert.False(t, m.IsPtr())

	result := m.Or(42)
	assert.Equal(t, 42, result)

	cloned := m.Clone()
	assert.True(t, cloned.IsNil())

	toMaybe := m.ToMaybe()
	assert.True(t, toMaybe.IsNil())
}

func TestCloneToWithNil(t *testing.T) {
	m := JustGenerics[interface{}](nil)
	var dest interface{}
	result := CloneTo(m, dest)
	assert.True(t, result.IsNil())
}

func TestCloneToWithPtr(t *testing.T) {
	val := 42
	m := JustGenerics(&val)
	dest := new(int)
	result := CloneTo(m, dest)
	assert.True(t, result.IsPresent())
	assert.NotNil(t, result.Unwrap())
}

func TestCloneToWithNonPtr(t *testing.T) {
	m := JustGenerics(42)
	var dest int
	result := CloneTo(m, dest)
	assert.Equal(t, 42, result.Unwrap())
}

func TestCloneToNone(t *testing.T) {
	m := JustGenerics[interface{}](nil)
	result := CloneTo(m, nil)
	assert.True(t, result.IsNil())
}

// Tests for coverage improvement - someDef[T] methods

func TestSomeDefOrWithNil(t *testing.T) {
	m := JustGenerics[interface{}](nil)
	assert.Equal(t, 3, m.Or(3))
}

func TestSomeDefToStringWithNil(t *testing.T) {
	m := JustGenerics[interface{}](nil)
	assert.Equal(t, "<nil>", m.ToString())
}

func TestSomeDefToMaybeWithNil(t *testing.T) {
	m := JustGenerics[interface{}](nil)
	m2 := m.ToMaybe()
	assert.True(t, m2.IsNil())
}

func TestSomeDefToFloat64WithNil(t *testing.T) {
	m := JustGenerics[interface{}](nil)
	_, err := m.ToFloat64()
	assert.Equal(t, ErrConversionNil, err)
}

func TestSomeDefToFloat64WithUnsupportedType(t *testing.T) {
	m := JustGenerics(make(chan int))
	_, err := m.ToFloat64()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToFloat64WithUint(t *testing.T) {
	m := JustGenerics(uint(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

func TestSomeDefToFloat64WithUint16(t *testing.T) {
	m := JustGenerics(uint16(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

func TestSomeDefToFloat64WithUint32(t *testing.T) {
	m := JustGenerics(uint32(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

func TestSomeDefToFloat64WithUint64(t *testing.T) {
	m := JustGenerics(uint64(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

func TestSomeDefToFloat64WithUintptr(t *testing.T) {
	m := JustGenerics(uintptr(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

func TestSomeDefToFloat64WithByte(t *testing.T) {
	m := JustGenerics(byte(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

func TestSomeDefToFloat64WithInt8(t *testing.T) {
	m := JustGenerics(int8(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

func TestSomeDefToFloat64WithInt16(t *testing.T) {
	m := JustGenerics(int16(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

func TestSomeDefToFloat64WithInt32(t *testing.T) {
	m := JustGenerics(int32(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

func TestSomeDefToFloat64WithInt64(t *testing.T) {
	m := JustGenerics(int64(42))
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42), val)
}

// Tests for ToInt coverage improvement

func TestSomeDefToIntWithUint(t *testing.T) {
	m := JustGenerics(uint(42))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestSomeDefToIntWithUint16(t *testing.T) {
	m := JustGenerics(uint16(42))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestSomeDefToIntWithUint32(t *testing.T) {
	m := JustGenerics(uint32(42))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestSomeDefToIntWithUint64(t *testing.T) {
	m := JustGenerics(uint64(42))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestSomeDefToIntWithUintptr(t *testing.T) {
	m := JustGenerics(uintptr(42))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestSomeDefToIntWithByte(t *testing.T) {
	m := JustGenerics(byte(42))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestSomeDefToIntWithInt8(t *testing.T) {
	m := JustGenerics(int8(42))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestSomeDefToIntWithInt16(t *testing.T) {
	m := JustGenerics(int16(42))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestSomeDefToIntWithInt32(t *testing.T) {
	m := JustGenerics(int32(42))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestSomeDefToIntWithInt64(t *testing.T) {
	m := JustGenerics(int64(42))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestSomeDefToIntWithFloat32(t *testing.T) {
	m := JustGenerics(float32(42.5))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 43, val)
}

func TestSomeDefToIntWithFloat64(t *testing.T) {
	m := JustGenerics(float64(42.5))
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 43, val)
}

// Tests for ToPtr coverage improvement

func TestSomeDefToPtrWithPtrValue(t *testing.T) {
	val := 42
	m := JustGenerics(&val)
	ptr := m.ToPtr()
	assert.NotNil(t, ptr)
	assert.Equal(t, 42, **ptr)
}

func TestSomeDefToPtrWithNonPtrValue(t *testing.T) {
	m := JustGenerics(42)
	ptr := m.ToPtr()
	assert.NotNil(t, ptr)
	assert.Equal(t, 42, *ptr)
}

// Tests for ToInt overflow cases

func TestSomeDefToIntWithUintOverflow(t *testing.T) {
	m := JustGenerics(uint(math.MaxUint32 + 1))
	_, err := m.ToInt()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToIntWithUint32Overflow(t *testing.T) {
	m := JustGenerics(uint32(math.MaxUint32))
	_, err := m.ToInt()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToIntWithUint64Overflow(t *testing.T) {
	m := JustGenerics(uint64(math.MaxUint64))
	_, err := m.ToInt()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToIntWithUintptrOverflow(t *testing.T) {
	m := JustGenerics(uintptr(math.MaxUint32 + 1))
	_, err := m.ToInt()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

// Tests for ToFloat32 coverage

func TestSomeDefToFloat32WithNil(t *testing.T) {
	m := JustGenerics[interface{}](nil)
	_, err := m.ToFloat32()
	assert.Equal(t, ErrConversionNil, err)
}

func TestSomeDefToFloat32WithUnsupportedType(t *testing.T) {
	m := JustGenerics(make(chan int))
	_, err := m.ToFloat32()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToFloat32WithString(t *testing.T) {
	m := JustGenerics("42.5")
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42.5), val)
}

func TestSomeDefToFloat32WithBoolTrue(t *testing.T) {
	m := JustGenerics(true)
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(1), val)
}

func TestSomeDefToFloat32WithBoolFalse(t *testing.T) {
	m := JustGenerics(false)
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(0), val)
}

func TestSomeDefToFloat32WithUint(t *testing.T) {
	m := JustGenerics(uint(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithUint16(t *testing.T) {
	m := JustGenerics(uint16(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithUint32(t *testing.T) {
	m := JustGenerics(uint32(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithUint64(t *testing.T) {
	m := JustGenerics(uint64(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithUintptr(t *testing.T) {
	m := JustGenerics(uintptr(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithByte(t *testing.T) {
	m := JustGenerics(byte(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithInt(t *testing.T) {
	m := JustGenerics(42)
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithInt8(t *testing.T) {
	m := JustGenerics(int8(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithInt16(t *testing.T) {
	m := JustGenerics(int16(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithInt32(t *testing.T) {
	m := JustGenerics(int32(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithInt64(t *testing.T) {
	m := JustGenerics(int64(42))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42), val)
}

func TestSomeDefToFloat32WithFloat32(t *testing.T) {
	m := JustGenerics(float32(42.5))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42.5), val)
}

func TestSomeDefToFloat32WithFloat64(t *testing.T) {
	m := JustGenerics(float64(42.5))
	val, err := m.ToFloat32()
	assert.NoError(t, err)
	assert.Equal(t, float32(42.5), val)
}

// Tests for ToUint coverage

func TestSomeDefToUintWithNil(t *testing.T) {
	m := JustGenerics[interface{}](nil)
	_, err := m.ToUint()
	assert.Equal(t, ErrConversionNil, err)
}

func TestSomeDefToUintWithUnsupportedType(t *testing.T) {
	m := JustGenerics(make(chan int))
	_, err := m.ToUint()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToUintWithString(t *testing.T) {
	m := JustGenerics("42")
	val, err := m.ToUint()
	assert.NoError(t, err)
	assert.Equal(t, uint(42), val)
}

func TestSomeDefToUintWithBoolTrue(t *testing.T) {
	m := JustGenerics(true)
	val, err := m.ToUint()
	assert.NoError(t, err)
	assert.Equal(t, uint(1), val)
}

func TestSomeDefToUintWithBoolFalse(t *testing.T) {
	m := JustGenerics(false)
	val, err := m.ToUint()
	assert.NoError(t, err)
	assert.Equal(t, uint(0), val)
}

func TestSomeDefToUintWithUint(t *testing.T) {
	m := JustGenerics(uint(42))
	val, err := m.ToUint()
	assert.NoError(t, err)
	assert.Equal(t, uint(42), val)
}

func TestSomeDefToUintWithUint16(t *testing.T) {
	m := JustGenerics(uint16(42))
	val, err := m.ToUint()
	assert.NoError(t, err)
	assert.Equal(t, uint(42), val)
}

func TestSomeDefToUintWithUint32(t *testing.T) {
	m := JustGenerics(uint32(42))
	val, err := m.ToUint()
	assert.NoError(t, err)
	assert.Equal(t, uint(42), val)
}

func TestSomeDefToUintWithUint64(t *testing.T) {
	m := JustGenerics(uint64(42))
	val, err := m.ToUint()
	assert.NoError(t, err)
	assert.Equal(t, uint(42), val)
}

func TestSomeDefToUintWithUintptr(t *testing.T) {
	m := JustGenerics(uintptr(42))
	val, err := m.ToUint()
	assert.NoError(t, err)
	assert.Equal(t, uint(42), val)
}

func TestSomeDefToUintWithByte(t *testing.T) {
	m := JustGenerics(byte(42))
	val, err := m.ToUint()
	assert.NoError(t, err)
	assert.Equal(t, uint(42), val)
}

func TestSomeDefToUintWithInt(t *testing.T) {
	m := JustGenerics(42)
	val, err := m.ToUint()
	assert.NoError(t, err)
	assert.Equal(t, uint(42), val)
}

func TestSomeDefToUintWithInt8(t *testing.T) {
	m := JustGenerics(int8(42))
	val, err := m.ToUint()
	assert.NoError(t, err)
	assert.Equal(t, uint(42), val)
}

func TestSomeDefToUintWithInt16(t *testing.T) {
	m := JustGenerics(int16(42))
	val, err := m.ToUint()
	assert.NoError(t, err)
	assert.Equal(t, uint(42), val)
}

func TestSomeDefToUintWithInt32(t *testing.T) {
	m := JustGenerics(int32(42))
	val, err := m.ToUint()
	assert.NoError(t, err)
	assert.Equal(t, uint(42), val)
}

func TestSomeDefToUintWithInt64(t *testing.T) {
	m := JustGenerics(int64(42))
	val, err := m.ToUint()
	assert.NoError(t, err)
	assert.Equal(t, uint(42), val)
}

func TestSomeDefToUintWithFloat32(t *testing.T) {
	m := JustGenerics(float32(42.5))
	val, err := m.ToUint()
	assert.NoError(t, err)
	assert.Equal(t, uint(43), val)
}

func TestSomeDefToUintWithFloat64(t *testing.T) {
	m := JustGenerics(float64(42.5))
	val, err := m.ToUint()
	assert.NoError(t, err)
	assert.Equal(t, uint(43), val)
}

func TestSomeDefToUintWithInt64Overflow(t *testing.T) {
	m := JustGenerics(int64(math.MaxUint32 + 1))
	_, err := m.ToUint()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToUintWithFloat32Overflow(t *testing.T) {
	m := JustGenerics(float32(math.MaxUint32 * 2))
	_, err := m.ToUint()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToUintWithFloat64Overflow(t *testing.T) {
	m := JustGenerics(float64(math.MaxUint32 + 1))
	_, err := m.ToUint()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

// Tests for ToPtr case *T branch

func TestSomeDefToPtrWithDoublePtr(t *testing.T) {
	val := 42
	ptr := &val
	m := JustGenerics(&ptr)
	result := m.ToPtr()
	assert.NotNil(t, result)
	// result is **int, so we need to dereference twice
	assert.Equal(t, 42, ***result)
}

// Tests for ToFloat64 string case

func TestSomeDefToFloat64WithString(t *testing.T) {
	m := JustGenerics("42.5")
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(42.5), val)
}

func TestSomeDefToFloat64WithBoolTrue(t *testing.T) {
	m := JustGenerics(true)
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(1), val)
}

func TestSomeDefToFloat64WithBoolFalse(t *testing.T) {
	m := JustGenerics(false)
	val, err := m.ToFloat64()
	assert.NoError(t, err)
	assert.Equal(t, float64(0), val)
}

// Tests for ToInt nil/string/unsupported cases

func TestSomeDefToIntWithNil(t *testing.T) {
	m := JustGenerics[interface{}](nil)
	_, err := m.ToInt()
	assert.Equal(t, ErrConversionNil, err)
}

func TestSomeDefToIntWithUnsupportedType(t *testing.T) {
	m := JustGenerics(make(chan int))
	_, err := m.ToInt()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToIntWithString(t *testing.T) {
	m := JustGenerics("42")
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestSomeDefToIntWithBoolTrue(t *testing.T) {
	m := JustGenerics(true)
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
}

func TestSomeDefToIntWithBoolFalse(t *testing.T) {
	m := JustGenerics(false)
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 0, val)
}

func TestSomeDefToIntWithInt(t *testing.T) {
	m := JustGenerics(42)
	val, err := m.ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

// Tests for ToInt8 coverage

func TestSomeDefToInt8WithNil(t *testing.T) {
	m := JustGenerics[interface{}](nil)
	_, err := m.ToInt8()
	assert.Equal(t, ErrConversionNil, err)
}

func TestSomeDefToInt8WithUnsupportedType(t *testing.T) {
	m := JustGenerics(make(chan int))
	_, err := m.ToInt8()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToInt8WithString(t *testing.T) {
	m := JustGenerics("42")
	val, err := m.ToInt8()
	assert.NoError(t, err)
	assert.Equal(t, int8(42), val)
}

func TestSomeDefToInt8WithBoolTrue(t *testing.T) {
	m := JustGenerics(true)
	val, err := m.ToInt8()
	assert.NoError(t, err)
	assert.Equal(t, int8(1), val)
}

func TestSomeDefToInt8WithBoolFalse(t *testing.T) {
	m := JustGenerics(false)
	val, err := m.ToInt8()
	assert.NoError(t, err)
	assert.Equal(t, int8(0), val)
}

func TestSomeDefToInt8WithInt8(t *testing.T) {
	m := JustGenerics(int8(42))
	val, err := m.ToInt8()
	assert.NoError(t, err)
	assert.Equal(t, int8(42), val)
}

func TestSomeDefToInt8WithIntOverflow(t *testing.T) {
	m := JustGenerics(128)
	_, err := m.ToInt8()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt8WithIntUnderflow(t *testing.T) {
	m := JustGenerics(-129)
	_, err := m.ToInt8()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

// Tests for ToInt16 coverage

func TestSomeDefToInt16WithNil(t *testing.T) {
	m := JustGenerics[interface{}](nil)
	_, err := m.ToInt16()
	assert.Equal(t, ErrConversionNil, err)
}

func TestSomeDefToInt16WithUnsupportedType(t *testing.T) {
	m := JustGenerics(make(chan int))
	_, err := m.ToInt16()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToInt16WithString(t *testing.T) {
	m := JustGenerics("42")
	val, err := m.ToInt16()
	assert.NoError(t, err)
	assert.Equal(t, int16(42), val)
}

func TestSomeDefToInt16WithInt16(t *testing.T) {
	m := JustGenerics(int16(42))
	val, err := m.ToInt16()
	assert.NoError(t, err)
	assert.Equal(t, int16(42), val)
}

// Tests for ToInt32 coverage

func TestSomeDefToInt32WithNil(t *testing.T) {
	m := JustGenerics[interface{}](nil)
	_, err := m.ToInt32()
	assert.Equal(t, ErrConversionNil, err)
}

func TestSomeDefToInt32WithUnsupportedType(t *testing.T) {
	m := JustGenerics(make(chan int))
	_, err := m.ToInt32()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToInt32WithString(t *testing.T) {
	m := JustGenerics("42")
	val, err := m.ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(42), val)
}

func TestSomeDefToInt32WithInt32(t *testing.T) {
	m := JustGenerics(int32(42))
	val, err := m.ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(42), val)
}

// Tests for ToInt64 coverage

func TestSomeDefToInt64WithNil(t *testing.T) {
	m := JustGenerics[interface{}](nil)
	_, err := m.ToInt64()
	assert.Equal(t, ErrConversionNil, err)
}

func TestSomeDefToInt64WithUnsupportedType(t *testing.T) {
	m := JustGenerics(make(chan int))
	_, err := m.ToInt64()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToInt64WithString(t *testing.T) {
	m := JustGenerics("42")
	val, err := m.ToInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(42), val)
}

func TestSomeDefToInt64WithInt64(t *testing.T) {
	m := JustGenerics(int64(42))
	val, err := m.ToInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(42), val)
}

// Tests for ToUint16 coverage

func TestSomeDefToUint16WithNil(t *testing.T) {
	m := JustGenerics[interface{}](nil)
	_, err := m.ToUint16()
	assert.Equal(t, ErrConversionNil, err)
}

func TestSomeDefToUint16WithUnsupportedType(t *testing.T) {
	m := JustGenerics(make(chan int))
	_, err := m.ToUint16()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToUint16WithString(t *testing.T) {
	m := JustGenerics("42")
	val, err := m.ToUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(42), val)
}

func TestSomeDefToUint16WithBoolTrue(t *testing.T) {
	m := JustGenerics(true)
	val, err := m.ToUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(1), val)
}

func TestSomeDefToUint16WithBoolFalse(t *testing.T) {
	m := JustGenerics(false)
	val, err := m.ToUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0), val)
}

func TestSomeDefToUint16WithUint16(t *testing.T) {
	m := JustGenerics(uint16(42))
	val, err := m.ToUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(42), val)
}

func TestSomeDefToUint16WithByte(t *testing.T) {
	m := JustGenerics(byte(42))
	val, err := m.ToUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(42), val)
}

func TestSomeDefToUint16WithInt8(t *testing.T) {
	m := JustGenerics(int8(42))
	val, err := m.ToUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(42), val)
}

func TestSomeDefToUint16WithIntOverflow(t *testing.T) {
	m := JustGenerics(65536)
	_, err := m.ToUint16()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToUint16WithIntNegative(t *testing.T) {
	m := JustGenerics(-1)
	_, err := m.ToUint16()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

// Tests for ToUint32 coverage

func TestSomeDefToUint32WithNil(t *testing.T) {
	m := JustGenerics[interface{}](nil)
	_, err := m.ToUint32()
	assert.Equal(t, ErrConversionNil, err)
}

func TestSomeDefToUint32WithUnsupportedType(t *testing.T) {
	m := JustGenerics(make(chan int))
	_, err := m.ToUint32()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToUint32WithString(t *testing.T) {
	m := JustGenerics("42")
	val, err := m.ToUint32()
	assert.NoError(t, err)
	assert.Equal(t, uint32(42), val)
}

func TestSomeDefToUint32WithBoolTrue(t *testing.T) {
	m := JustGenerics(true)
	val, err := m.ToUint32()
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), val)
}

func TestSomeDefToUint32WithBoolFalse(t *testing.T) {
	m := JustGenerics(false)
	val, err := m.ToUint32()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), val)
}

func TestSomeDefToUint32WithUint32(t *testing.T) {
	m := JustGenerics(uint32(42))
	val, err := m.ToUint32()
	assert.NoError(t, err)
	assert.Equal(t, uint32(42), val)
}

func TestSomeDefToUint32WithByte(t *testing.T) {
	m := JustGenerics(byte(42))
	val, err := m.ToUint32()
	assert.NoError(t, err)
	assert.Equal(t, uint32(42), val)
}

func TestSomeDefToUint32WithIntOverflow(t *testing.T) {
	m := JustGenerics(int64(math.MaxUint32 + 1))
	_, err := m.ToUint32()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

// Tests for ToUint64 coverage

func TestSomeDefToUint64WithNil(t *testing.T) {
	m := JustGenerics[interface{}](nil)
	_, err := m.ToUint64()
	assert.Equal(t, ErrConversionNil, err)
}

func TestSomeDefToUint64WithUnsupportedType(t *testing.T) {
	m := JustGenerics(make(chan int))
	_, err := m.ToUint64()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToUint64WithString(t *testing.T) {
	m := JustGenerics("42")
	val, err := m.ToUint64()
	assert.NoError(t, err)
	assert.Equal(t, uint64(42), val)
}

func TestSomeDefToUint64WithBoolTrue(t *testing.T) {
	m := JustGenerics(true)
	val, err := m.ToUint64()
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), val)
}

func TestSomeDefToUint64WithBoolFalse(t *testing.T) {
	m := JustGenerics(false)
	val, err := m.ToUint64()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), val)
}

func TestSomeDefToUint64WithUint64(t *testing.T) {
	m := JustGenerics(uint64(42))
	val, err := m.ToUint64()
	assert.NoError(t, err)
	assert.Equal(t, uint64(42), val)
}

func TestSomeDefToUint64WithByte(t *testing.T) {
	m := JustGenerics(byte(42))
	val, err := m.ToUint64()
	assert.NoError(t, err)
	assert.Equal(t, uint64(42), val)
}

// Tests for ToUintptr coverage

func TestSomeDefToUintptrWithNil(t *testing.T) {
	m := JustGenerics[interface{}](nil)
	_, err := m.ToUintptr()
	assert.Equal(t, ErrConversionNil, err)
}

func TestSomeDefToUintptrWithUnsupportedType(t *testing.T) {
	m := JustGenerics(make(chan int))
	_, err := m.ToUintptr()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToUintptrWithString(t *testing.T) {
	m := JustGenerics("42")
	val, err := m.ToUintptr()
	assert.NoError(t, err)
	assert.Equal(t, uintptr(42), val)
}

func TestSomeDefToUintptrWithBoolTrue(t *testing.T) {
	m := JustGenerics(true)
	val, err := m.ToUintptr()
	assert.NoError(t, err)
	assert.Equal(t, uintptr(1), val)
}

func TestSomeDefToUintptrWithBoolFalse(t *testing.T) {
	m := JustGenerics(false)
	val, err := m.ToUintptr()
	assert.NoError(t, err)
	assert.Equal(t, uintptr(0), val)
}

func TestSomeDefToUintptrWithUintptr(t *testing.T) {
	m := JustGenerics(uintptr(42))
	val, err := m.ToUintptr()
	assert.NoError(t, err)
	assert.Equal(t, uintptr(42), val)
}

func TestSomeDefToUintptrWithByte(t *testing.T) {
	m := JustGenerics(byte(42))
	val, err := m.ToUintptr()
	assert.NoError(t, err)
	assert.Equal(t, uintptr(42), val)
}

// Tests for ToByte coverage

func TestSomeDefToByteWithNil(t *testing.T) {
	m := JustGenerics[interface{}](nil)
	_, err := m.ToByte()
	assert.Equal(t, ErrConversionNil, err)
}

func TestSomeDefToByteWithUnsupportedType(t *testing.T) {
	m := JustGenerics(make(chan int))
	_, err := m.ToByte()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToByteWithString(t *testing.T) {
	m := JustGenerics("42")
	val, err := m.ToByte()
	assert.NoError(t, err)
	assert.Equal(t, byte(42), val)
}

func TestSomeDefToByteWithBoolTrue(t *testing.T) {
	m := JustGenerics(true)
	val, err := m.ToByte()
	assert.NoError(t, err)
	assert.Equal(t, byte(1), val)
}

func TestSomeDefToByteWithBoolFalse(t *testing.T) {
	m := JustGenerics(false)
	val, err := m.ToByte()
	assert.NoError(t, err)
	assert.Equal(t, byte(0), val)
}

func TestSomeDefToByteWithByte(t *testing.T) {
	m := JustGenerics(byte(42))
	val, err := m.ToByte()
	assert.NoError(t, err)
	assert.Equal(t, byte(42), val)
}

func TestSomeDefToByteWithIntOverflow(t *testing.T) {
	m := JustGenerics(256)
	_, err := m.ToByte()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToByteWithIntNegative(t *testing.T) {
	m := JustGenerics(-1)
	_, err := m.ToByte()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

// Tests for ToBool coverage

func TestSomeDefToBoolWithNil(t *testing.T) {
	m := JustGenerics[interface{}](nil)
	_, err := m.ToBool()
	assert.Equal(t, ErrConversionNil, err)
}

func TestSomeDefToBoolWithUnsupportedType(t *testing.T) {
	m := JustGenerics(make(chan int))
	_, err := m.ToBool()
	assert.Equal(t, ErrConversionUnsupported, err)
}

func TestSomeDefToBoolWithString(t *testing.T) {
	m := JustGenerics("true")
	val, err := m.ToBool()
	assert.NoError(t, err)
	assert.True(t, val)
}

func TestSomeDefToBoolWithBoolTrue(t *testing.T) {
	m := JustGenerics(true)
	val, err := m.ToBool()
	assert.NoError(t, err)
	assert.True(t, val)
}

func TestSomeDefToBoolWithBoolFalse(t *testing.T) {
	m := JustGenerics(false)
	val, err := m.ToBool()
	assert.NoError(t, err)
	assert.False(t, val)
}

func TestSomeDefToBoolWithByte(t *testing.T) {
	m := JustGenerics(byte(1))
	val, err := m.ToBool()
	assert.NoError(t, err)
	assert.True(t, val)
}

func TestSomeDefToBoolWithByteZero(t *testing.T) {
	m := JustGenerics(byte(0))
	val, err := m.ToBool()
	assert.NoError(t, err)
	assert.False(t, val)
}

// Tests for ToInt32 overflow cases

func TestSomeDefToInt32WithUintOverflow(t *testing.T) {
	m := JustGenerics(uint(math.MaxInt32 + 1))
	_, err := m.ToInt32()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt32WithUint32Overflow(t *testing.T) {
	m := JustGenerics(uint32(math.MaxInt32 + 1))
	_, err := m.ToInt32()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt32WithUint64Overflow(t *testing.T) {
	m := JustGenerics(uint64(math.MaxInt32 + 1))
	_, err := m.ToInt32()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt32WithUintptrOverflow(t *testing.T) {
	m := JustGenerics(uintptr(math.MaxInt32 + 1))
	_, err := m.ToInt32()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt32WithInt64Overflow(t *testing.T) {
	m := JustGenerics(int64(math.MaxInt32 + 1))
	_, err := m.ToInt32()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt32WithInt64Underflow(t *testing.T) {
	m := JustGenerics(int64(math.MinInt32 - 1))
	_, err := m.ToInt32()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt32WithFloat32Overflow(t *testing.T) {
	m := JustGenerics(float32(math.MaxInt32 * 2))
	_, err := m.ToInt32()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt32WithFloat64(t *testing.T) {
	m := JustGenerics(float64(42.5))
	val, err := m.ToInt32()
	assert.NoError(t, err)
	assert.Equal(t, int32(43), val)
}

// Tests for ToInt64 overflow cases

func TestSomeDefToInt64WithUint64Overflow(t *testing.T) {
	m := JustGenerics(uint64(math.MaxInt64 + 1))
	_, err := m.ToInt64()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt64WithUintptrOverflow(t *testing.T) {
	// On 64-bit systems, uintptr can be large
	m := JustGenerics(uintptr(1 << 62))
	val, err := m.ToInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(1<<62), val)
}

func TestSomeDefToInt64WithFloat32Overflow(t *testing.T) {
	m := JustGenerics(float32(math.MaxInt32 * 2))
	_, err := m.ToInt64()
	assert.NoError(t, err) // Should succeed since float32*2 < MaxInt64
}

func TestSomeDefToInt64WithFloat64Overflow(t *testing.T) {
	m := JustGenerics(float64(math.MaxInt64) * 2)
	_, err := m.ToInt64()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt64WithFloat64Underflow(t *testing.T) {
	m := JustGenerics(float64(math.MinInt64) * 2)
	_, err := m.ToInt64()
	assert.Equal(t, ErrConversionSizeOverflow, err)
}

func TestSomeDefToInt64WithUint(t *testing.T) {
	m := JustGenerics(uint(42))
	val, err := m.ToInt64()
	assert.NoError(t, err)
	assert.Equal(t, int64(42), val)
}

func TestSomeDefToInt64WithUintOverflow(t *testing.T) {
	// On 64-bit systems, uint can be > MaxInt64
	if math.MaxUint > math.MaxInt64 {
		m := JustGenerics(uint(math.MaxUint))
		_, err := m.ToInt64()
		assert.Equal(t, ErrConversionSizeOverflow, err)
	}
}
