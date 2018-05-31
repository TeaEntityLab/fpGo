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
		return Monad.JustVal(val + 1).Ref()
	}
	var fn02 = func(obj *interface{}) *interface{} {
		val, _ := Monad.Just(obj).ToInt()
		return Monad.JustVal(val + 2).Ref()
	}
	var fn03 = func(obj *interface{}) *interface{} {
		val, _ := Monad.Just(obj).ToInt()
		return Monad.JustVal(val + 3).Ref()
	}

	expectedinteger = 1
	assert.Equal(t, expectedinteger, *Compose(fn01)(Monad.JustVal(0).Ref()))

	expectedinteger = 2
	assert.Equal(t, expectedinteger, *Compose(fn02)(Monad.JustVal(0).Ref()))

	expectedinteger = 3
	assert.Equal(t, expectedinteger, *Compose(fn03)(Monad.JustVal(0).Ref()))

	expectedinteger = 3
	assert.Equal(t, expectedinteger, *Compose(fn01, fn02)(Monad.JustVal(0).Ref()))

	expectedinteger = 4
	assert.Equal(t, expectedinteger, *Compose(fn01, fn03)(Monad.JustVal(0).Ref()))

	expectedinteger = 5
	assert.Equal(t, expectedinteger, *Compose(fn02, fn03)(Monad.JustVal(0).Ref()))

	expectedinteger = 6
	assert.Equal(t, expectedinteger, *Compose(fn01, fn02, fn03)(Monad.JustVal(0).Ref()))
}

func TestCompType(t *testing.T) {
	var compTypeA CompType = DefProduct(reflect.Int, reflect.String)
	var compTypeB CompType = DefProduct(reflect.String)
	var myType CompType = DefSum(NilType{}, compTypeA, compTypeB)

	assert.Equal(t, true, myType.Matches(Monad.JustVal(1).Ref(), Monad.JustVal("1").Ref()))
}
