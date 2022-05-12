package fpgo

import (
	"errors"
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
