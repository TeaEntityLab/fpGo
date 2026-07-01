package fpgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReduceIndexed(t *testing.T) {
	result := ReduceIndexed(func(acc, val interface{}, i int) interface{} {
		a, _ := Maybe.Just(acc).ToInt()
		v, _ := Maybe.Just(val).ToInt()
		return a + v + i
	}, 0, 1, 2, 3)
	assert.Equal(t, 9, result)
}

func TestPtrOf(t *testing.T) {
	val := 42
	ptr := PtrOf(val)
	assert.NotNil(t, ptr)
	assert.Equal(t, 42, *ptr)
}

func TestCurryParam1ForSlice1(t *testing.T) {
	fn := func(a interface{}, b []interface{}) interface{} {
		ai, _ := Maybe.Just(a).ToInt()
		bi, _ := Maybe.Just(b[0]).ToInt()
		return ai + bi
	}
	curried := CurryParam1ForSlice1(fn, 10)
	assert.Equal(t, 11, curried(1))
}

func TestMakeNumericReturnForSliceParamReturnBool1(t *testing.T) {
	fn := func(v []interface{}) bool { return len(v) > 0 }
	wrapped := MakeNumericReturnForSliceParamReturnBool1(fn)
	assert.Equal(t, []interface{}{1}, wrapped(1, 2, 3))
	assert.Equal(t, []interface{}{0}, wrapped())
}

func TestComposeIdentity(t *testing.T) {
	assert.Nil(t, Compose()())
	assert.Equal(t, []interface{}{1, 2, 3}, Compose()(1, 2, 3))
}

func TestPipeIdentity(t *testing.T) {
	assert.Nil(t, Pipe()())
	assert.Equal(t, []interface{}{1, 2, 3}, Pipe()(1, 2, 3))
}

func TestDropBranches(t *testing.T) {
	assert.Equal(t, []interface{}{1, 2, 3}, Drop(0, 1, 2, 3))
	assert.Equal(t, []interface{}{3}, Drop(2, 1, 2, 3))
	assert.Equal(t, []interface{}{}, Drop(5, 1, 2))
}

func TestDropLastBranches(t *testing.T) {
	assert.Equal(t, []interface{}{1, 2}, DropLast(0, 1, 2))
	assert.Equal(t, []interface{}{1}, DropLast(2, 1, 2, 3))
	assert.Equal(t, []interface{}{}, DropLast(3, 1, 2))
}

func TestMergeBranches(t *testing.T) {
	assert.Equal(t, map[interface{}]interface{}{}, Merge(nil, nil))
	assert.Equal(t, map[interface{}]interface{}{"b": 2}, Merge(nil, map[interface{}]interface{}{"b": 2}))
	assert.Equal(t, map[interface{}]interface{}{"a": 1}, Merge(map[interface{}]interface{}{"a": 1}, nil))
	result := Merge(map[interface{}]interface{}{"a": 1}, map[interface{}]interface{}{"a": 2, "b": 3})
	assert.Equal(t, 2, result["a"])
	assert.Equal(t, 3, result["b"])
}

func TestTakeBranches(t *testing.T) {
	assert.Equal(t, []interface{}{1, 2, 3}, Take(0, 1, 2, 3))
	assert.Equal(t, []interface{}{1, 2, 3}, Take(10, 1, 2, 3))
	assert.Equal(t, []interface{}{1, 2}, Take(2, 1, 2, 3))
}

func TestIntersectionBranches(t *testing.T) {
	var nilInput [][]interface{}
	assert.Equal(t, []interface{}{}, Intersection(nilInput...))
	assert.Equal(t, []interface{}{1, 2, 3}, Intersection([]interface{}{1, 2, 2, 3}))
	assert.Equal(t, []interface{}{2, 3}, Intersection(
		[]interface{}{1, 2, 3, 4},
		[]interface{}{2, 3, 5},
	))
}
