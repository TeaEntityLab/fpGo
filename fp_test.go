package fpGo

import (
	"fmt"
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

func TestPatternMatching(t *testing.T) {
	var patterns = []Pattern{
		InCaseOfKind(reflect.Int, func(x *interface{}) *interface{} {
			return PtrOf(fmt.Sprintf("Integer: %v", *x))
		}),
		InCaseOfEqual(PtrOf("world"), func(x *interface{}) *interface{} {
			return PtrOf(fmt.Sprintf("Hello %v", *x))
		}),
		InCaseOfRegex("c+", func(x *interface{}) *interface{} {
			return PtrOf(fmt.Sprintf("Matched: %v", *x))
		}),
		Otherwise(func(x *interface{}) *interface{} {
			return PtrOf(fmt.Sprintf("got this object: %v", *x))
		}),
	}
	var pm PatternMatching = DefPattern(patterns...)
	assert.Equal(t, "Integer: 42", *pm.MatchFor(PtrOf(42)))
	assert.Equal(t, "Hello world", *pm.MatchFor(PtrOf("world")))
	assert.Equal(t, "Matched: ccc", *pm.MatchFor(PtrOf("ccc")))
	assert.Equal(t, "got this object: TEST", *pm.MatchFor(PtrOf("TEST")))

	assert.Equal(t, "Integer: 42", *Either(42, patterns...))
	assert.Equal(t, "Hello world", *Either("world", patterns...))
	assert.Equal(t, "Matched: ccc", *Either("ccc", patterns...))
	assert.Equal(t, "got this object: TEST", *Either("TEST", patterns...))
}
