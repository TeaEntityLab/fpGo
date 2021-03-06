package fpgo

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompose(t *testing.T) {
	var expectedinteger int

	var fn01 = func(args ...interface{}) []interface{} {
		val, _ := Maybe.Just(args[0]).ToInt()
		return SliceOf(val + 1)
	}
	var fn02 = func(args ...interface{}) []interface{} {
		val, _ := Maybe.Just(args[0]).ToInt()
		return SliceOf(val + 2)
	}
	var fn03 = func(args ...interface{}) []interface{} {
		val, _ := Maybe.Just(args[0]).ToInt()
		return SliceOf(val + 3)
	}

	expectedinteger = 1
	assert.Equal(t, expectedinteger, Compose(fn01)((0))[0])

	expectedinteger = 2
	assert.Equal(t, expectedinteger, Compose(fn02)((0))[0])

	expectedinteger = 3
	assert.Equal(t, expectedinteger, Compose(fn03)((0))[0])

	expectedinteger = 3
	assert.Equal(t, expectedinteger, Compose(fn01, fn02)((0))[0])

	expectedinteger = 4
	assert.Equal(t, expectedinteger, Compose(fn01, fn03)((0))[0])

	expectedinteger = 5
	assert.Equal(t, expectedinteger, Compose(fn02, fn03)((0))[0])

	expectedinteger = 6
	assert.Equal(t, expectedinteger, Compose(fn01, fn02, fn03)((0))[0])
}

func TestCurry(t *testing.T) {
	c := Curry.New(func(c *CurryDef, args ...interface{}) interface{} {
		result := 0
		if len(args) == 3 {
			var v int
			v, _ = Maybe.Just(args[0]).ToInt()
			// fmt.Println(v)
			result += v
			v, _ = Maybe.Just(args[1]).ToInt()
			// fmt.Println(v)
			result += v
			v, _ = Maybe.Just(args[2]).ToInt()
			// fmt.Println(v)
			result += v

			c.MarkDone()
		}
		return (result)
	})

	assert.Equal(t, false, c.IsDone())
	c.Call((1))
	assert.Equal(t, false, c.IsDone())
	c.Call((2))
	assert.Equal(t, false, c.IsDone())
	c.Call((3))
	assert.Equal(t, true, c.IsDone())
	assert.Equal(t, 6, c.Result())
}

func TestCompType(t *testing.T) {
	var compTypeA = DefProduct(reflect.Int, reflect.String)
	var compTypeB = DefProduct(reflect.String)
	var myType = DefSum(NilType, compTypeA, compTypeB)

	assert.Equal(t, true, myType.Matches((1), ("1")))
	assert.Equal(t, true, myType.Matches(("2")))
	assert.Equal(t, true, myType.Matches(nil))
	assert.Equal(t, false, myType.Matches((1), (1)))
	assert.Equal(t, true, MatchCompTypeRef(myType, NewCompData(myType, (1), ("1"))))
	assert.Equal(t, true, MatchCompTypeRef(myType, NewCompData(myType, ("2"))))
	assert.Equal(t, true, MatchCompTypeRef(myType, NewCompData(myType, nil)))
	assert.Equal(t, true, NewCompData(myType, (1), (1)) == nil)
}

func TestPatternMatching(t *testing.T) {
	var compTypeA = DefProduct(reflect.Int, reflect.String)
	var compTypeB = DefProduct(reflect.String, reflect.String)
	var myType = DefSum(NilType, compTypeA, compTypeB)

	assert.Equal(t, true, compTypeA.Matches(1, "3"))
	assert.Equal(t, false, compTypeA.Matches(1, 3))
	assert.Equal(t, true, myType.Matches(nil))
	assert.Equal(t, true, myType.Matches(1, "3"))
	assert.Equal(t, true, myType.Matches("1", "3"))
	assert.Equal(t, false, myType.Matches(1, 3))

	var patterns = []Pattern{
		InCaseOfKind(reflect.Int, func(x interface{}) interface{} {
			return (fmt.Sprintf("Integer: %v", x))
		}),
		InCaseOfEqual(("world"), func(x interface{}) interface{} {
			return (fmt.Sprintf("Hello %v", x))
		}),
		InCaseOfSumType(myType, func(x interface{}) interface{} {
			return (fmt.Sprintf("SumType %v %v", (x).(CompData).objects[0], (x).(CompData).objects[0]))
		}),
		InCaseOfRegex("c+", func(x interface{}) interface{} {
			return (fmt.Sprintf("Matched: %v", x))
		}),
		Otherwise(func(x interface{}) interface{} {
			return (fmt.Sprintf("got this object: %v", x))
		}),
	}
	var pm = DefPattern(patterns...)
	assert.Equal(t, "Integer: 42", pm.MatchFor((42)))
	assert.Equal(t, "Hello world", pm.MatchFor(("world")))
	assert.Equal(t, "Matched: ccc", pm.MatchFor(("ccc")))
	assert.Equal(t, "SumType 1 1", pm.MatchFor((NewCompData(myType, ("1"), ("1")))))
	assert.Equal(t, "got this object: TEST", pm.MatchFor(("TEST")))

	assert.Equal(t, "Integer: 42", Either(42, patterns...))
	assert.Equal(t, "Hello world", Either("world", patterns...))
	assert.Equal(t, "Matched: ccc", Either("ccc", patterns...))
	assert.Equal(t, "SumType 1 1", Either(NewCompData(myType, ("1"), ("1")), patterns...))
	assert.Equal(t, "got this object: TEST", Either("TEST", patterns...))
}
