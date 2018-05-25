// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fpGo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromArrayMapReduce(t *testing.T) {
	var s StreamDef
	var tempString = ""

	s = Stream.FromArrayMonad([]MonadDef{Monad.JustVal("1"), Monad.JustVal("2"), Monad.JustVal("3"), Monad.JustVal("4")})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Monad.JustVal(v).ToMonad().ToString()
	}
	assert.Equal(t, "1234", tempString)
	s.Map(func(index int) interface{} {
		var val = Monad.JustVal(s.Get(index)).ToMonad().ToString()
		return "v" + val
	})
	assert.Equal(t, "{[v1 v2 v3 v4]}", Monad.JustVal(s).ToString())
	assert.Equal(t, "[v1 v2 v3 v4]", Monad.JustVal(s.ToArray()).ToString())
	tempString = ""

	s = Stream.FromArrayString([]string{"1", "2", "3", "4"})
	assert.Equal(t, "{[1 2 3 4]}", Monad.JustVal(s).ToString())
	s.Map(func(index int) interface{} {
		var val = Monad.JustVal(s.Get(index)).ToString()
		return "v" + val
	})
	assert.Equal(t, "{[v1 v2 v3 v4]}", Monad.JustVal(s).ToString())
	assert.Equal(t, "[v1 v2 v3 v4]", Monad.JustVal(s.ToArray()).ToString())

	s = Stream.FromArrayInt([]int{1, 2, 3, 4})
	assert.Equal(t, "{[1 2 3 4]}", Monad.JustVal(s).ToString())
	s.Map(func(index int) interface{} {
		var val, _ = Monad.JustVal(s.Get(index)).ToInt()
		return val * val
	})
	assert.Equal(t, "{[1 4 9 16]}", Monad.JustVal(s).ToString())
	assert.Equal(t, "[1 4 9 16]", Monad.JustVal(s.ToArray()).ToString())

	s = Stream.FromArrayFloat32([]float32{1, 2, 3, 4})
	assert.Equal(t, "{[1 2 3 4]}", Monad.JustVal(s).ToString())
	s.Map(func(index int) interface{} {
		var val, _ = Monad.JustVal(s.Get(index)).ToFloat32()
		return val * val
	})
	assert.Equal(t, "{[1 4 9 16]}", Monad.JustVal(s).ToString())
	assert.Equal(t, "[1 4 9 16]", Monad.JustVal(s.ToArray()).ToString())

	s = Stream.FromArrayFloat64([]float64{1, 2, 3, 4})
	assert.Equal(t, "{[1 2 3 4]}", Monad.JustVal(s).ToString())
	s.Map(func(index int) interface{} {
		var val, _ = Monad.JustVal(s.Get(index)).ToFloat64()
		return val * val
	})
	assert.Equal(t, "{[1 4 9 16]}", Monad.JustVal(s).ToString())
	assert.Equal(t, "[1 4 9 16]", Monad.JustVal(s.ToArray()).ToString())
}

func TestFilter(t *testing.T) {
	var s StreamDef

	s = Stream.FromArrayInt([]int{1}).Extend(Stream.FromArrayInt([]int{2, 3, 4}))
	assert.Equal(t, "{[1 2 3 4]}", Monad.JustVal(s).ToString())
	s = s.Filter(func(index int) bool {
		var val, _ = Monad.JustVal(s.Get(index)).ToInt()
		return val > 1 && val < 4
	})
	assert.Equal(t, "{[2 3]}", Monad.JustVal(s).ToString())
	assert.Equal(t, "[2 3]", Monad.JustVal(s.ToArray()).ToString())
}
