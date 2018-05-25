// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fpGo

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPresent(t *testing.T) {
	var m MonadDef

	m = Monad.JustVal(1)
	assert.Equal(t, true, m.IsPresent())
	assert.Equal(t, false, m.IsNil())

	m = Monad.Just(nil)
	assert.Equal(t, false, m.IsPresent())
	assert.Equal(t, true, m.IsNil())
}

func TestOr(t *testing.T) {
	var m MonadDef

	m = Monad.JustVal(1)
	assert.Equal(t, 1, m.OrVal(3).Unwrap())
	m = Monad.Just(nil)
	assert.Equal(t, 3, m.OrVal(3).Unwrap())
}

func TestLet(t *testing.T) {
	var m MonadDef

	var letVal int

	letVal = 1
	m = Monad.JustVal(1)
	m.Let(func() {
		letVal = 2
	})
	assert.Equal(t, 2, letVal)

	letVal = 1
	m = Monad.Just(nil)
	m.Let(func() {
		letVal = 3
	})
	assert.Equal(t, 1, letVal)
}

func TestType(t *testing.T) {
	var m MonadDef

	m = Monad.JustVal(1)
	assert.Equal(t, reflect.Int, m.Kind())
	assert.Equal(t, true, m.IsKind(reflect.Int))
	assert.Equal(t, false, m.IsKind(reflect.Ptr))

	assert.Equal(t, reflect.TypeOf(1), m.Type())
	assert.Equal(t, true, m.IsType(reflect.TypeOf(1)))
	assert.Equal(t, false, m.IsType(reflect.TypeOf(nil)))

	m = Monad.Just(nil)
	assert.Equal(t, reflect.Ptr, m.Kind())
	assert.Equal(t, false, m.IsKind(reflect.Int))
	assert.Equal(t, true, m.IsKind(reflect.Ptr))

	assert.Equal(t, reflect.TypeOf(nil), m.Type())
	assert.Equal(t, true, m.IsType(reflect.TypeOf(nil)))
	assert.Equal(t, false, m.IsType(reflect.TypeOf(1)))
}

func TestCast(t *testing.T) {
	var m MonadDef

	var f32 float32
	var f64 float64
	var i int
	var err error

	// Int
	m = Monad.JustVal(1)
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

	// Float32
	m = Monad.JustVal(float32(1.1))
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

	// Float64
	m = Monad.JustVal(float64(1.2))
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

	// Nil
	m = Monad.Just(nil)
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
}
