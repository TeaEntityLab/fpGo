package fpGo

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompose(t *testing.T) {
	var expectedinteger = 0

	var fn01 = func(obj *interface{}) *interface{} {
		val, _ := Monad.Just(obj).ToInt()
		return PtrOf(val + 1)
	}
	var fn02 = func(obj *interface{}) *interface{} {
		val, _ := Monad.Just(obj).ToInt()
		return PtrOf(val + 2)
	}
	var fn03 = func(obj *interface{}) *interface{} {
		val, _ := Monad.Just(obj).ToInt()
		return PtrOf(val + 3)
	}

	expectedinteger = 1
	assert.Equal(t, expectedinteger, *Compose(fn01)(PtrOf(0)))

	expectedinteger = 2
	assert.Equal(t, expectedinteger, *Compose(fn02)(PtrOf(0)))

	expectedinteger = 3
	assert.Equal(t, expectedinteger, *Compose(fn03)(PtrOf(0)))

	expectedinteger = 3
	assert.Equal(t, expectedinteger, *Compose(fn01, fn02)(PtrOf(0)))

	expectedinteger = 4
	assert.Equal(t, expectedinteger, *Compose(fn01, fn03)(PtrOf(0)))

	expectedinteger = 5
	assert.Equal(t, expectedinteger, *Compose(fn02, fn03)(PtrOf(0)))

	expectedinteger = 6
	assert.Equal(t, expectedinteger, *Compose(fn01, fn02, fn03)(PtrOf(0)))
}

func TestCompType(t *testing.T) {
	var compTypeA CompType = DefProduct(reflect.Int, reflect.String)
	var compTypeB CompType = DefProduct(reflect.String)
	var myType CompType = DefSum(NilType, compTypeA, compTypeB)

	assert.Equal(t, true, myType.Matches(PtrOf(1), PtrOf("1")))
	assert.Equal(t, true, myType.Matches(PtrOf("2")))
	assert.Equal(t, true, myType.Matches(nil))
	assert.Equal(t, false, myType.Matches(PtrOf(1), PtrOf(1)))
	assert.Equal(t, true, MatchCompTypeRef(myType, NewCompData(myType, PtrOf(1), PtrOf("1"))))
	assert.Equal(t, true, MatchCompTypeRef(myType, NewCompData(myType, PtrOf("2"))))
	assert.Equal(t, true, MatchCompTypeRef(myType, NewCompData(myType, nil)))
	assert.Equal(t, true, NewCompData(myType, PtrOf(1), PtrOf(1)) == nil)
}
