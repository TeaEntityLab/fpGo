package fpGo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromArrayMapReduce(t *testing.T) {
	var s *StreamDef
	var tempString = ""

	s = Stream.FromArrayMonad([]MonadDef{Monad.JustVal("1"), Monad.JustVal("2"), Monad.JustVal("3"), Monad.JustVal("4")})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Monad.Just(v).ToMonad().ToString()
	}
	assert.Equal(t, "1234", tempString)
	s = s.Map(func(index int) *interface{} {
		var val = Monad.Just(s.Get(index)).ToMonad().ToString()
		var result interface{} = "v" + val
		return &result
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Monad.Just(v).ToMonad().ToString()
	}
	assert.Equal(t, "v1v2v3v4", tempString)
	tempString = ""

	s = Stream.FromArrayString([]string{"1", "2", "3", "4"})
	for _, v := range s.ToArray() {
		tempString += Monad.Just(v).ToMonad().ToString()
	}
	assert.Equal(t, "1234", tempString)
	s = s.Map(func(index int) *interface{} {
		var val = Monad.Just(s.Get(index)).ToString()
		var result interface{} = "v" + val
		return &result
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Monad.Just(v).ToMonad().ToString()
	}
	assert.Equal(t, "v1v2v3v4", tempString)

	s = Stream.FromArrayInt([]int{1, 2, 3, 4})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Monad.Just(v).ToMonad().ToString()
	}
	assert.Equal(t, "1234", tempString)
	s = s.Map(func(index int) *interface{} {
		var val, _ = Monad.Just(s.Get(index)).ToInt()
		var result interface{} = val * val
		return &result
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Monad.Just(v).ToMonad().ToString()
	}
	assert.Equal(t, "14916", tempString)

	s = Stream.FromArrayFloat32([]float32{1, 2, 3, 4})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Monad.Just(v).ToMonad().ToString()
	}
	assert.Equal(t, "1234", tempString)
	s = s.Map(func(index int) *interface{} {
		var val, _ = Monad.Just(s.Get(index)).ToFloat32()
		var result interface{} = val * val
		return &result
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Monad.Just(v).ToMonad().ToString()
	}
	assert.Equal(t, "14916", tempString)

	s = Stream.FromArrayFloat64([]float64{1, 2, 3, 4})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Monad.Just(v).ToMonad().ToString()
	}
	assert.Equal(t, "1234", tempString)

	s = s.Map(func(index int) *interface{} {
		var val, _ = Monad.Just(s.Get(index)).ToFloat64()
		var result interface{} = val * val
		return &result
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Monad.Just(v).ToMonad().ToString()
	}
	assert.Equal(t, "14916", tempString)
}

func TestFilter(t *testing.T) {
	var s *StreamDef
	var tempString = ""

	s = Stream.FromArrayInt([]int{1}).Extend(Stream.FromArrayInt([]int{2, 3, 4}))
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Monad.Just(v).ToMonad().ToString()
	}
	assert.Equal(t, "1234", tempString)

	s = s.Filter(func(index int) bool {
		var val, _ = Monad.Just(s.Get(index)).ToInt()
		return val > 1 && val < 4
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Monad.Just(v).ToMonad().ToString()
	}
	assert.Equal(t, "23", tempString)
}

func TestSort(t *testing.T) {
	var s *StreamDef
	var tempString = ""

	s = Stream.FromArrayInt([]int{11}).Extend(Stream.FromArrayInt([]int{2, 3, 4, 5})).Remove(4)
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Monad.Just(v).ToMonad().ToString()
	}
	assert.Equal(t, "11234", tempString)

	s = s.Sort(func(i, j int) bool {
		var vali, _ = Monad.Just(s.Get(i)).ToInt()
		var valj, _ = Monad.Just(s.Get(j)).ToInt()
		return vali < valj
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Monad.Just(v).ToMonad().ToString()
	}
	assert.Equal(t, "23411", tempString)
}
