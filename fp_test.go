package fpgo

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SortIntAscending(input ...interface{}) []interface{} {
	Sort(func(a, b interface{}) bool {
		aVal, _ := Maybe.Just(a).ToInt()
		bVal, _ := Maybe.Just(b).ToInt()
		return bVal > aVal
	}, input)

	return input
}

func SortStringAscending(input ...interface{}) []interface{} {
	Sort(func(a, b interface{}) bool {
		aVal := Maybe.Just(a).ToString()
		bVal := Maybe.Just(b).ToString()
		return strings.Compare(aVal, bVal) < 0
	}, input)

	return input
}

func TestCompose(t *testing.T) {
	var expectedinteger int

	fn01 := func(args ...interface{}) []interface{} {
		val, _ := Maybe.Just(args[0]).ToInt()
		return SliceOf(val + 1)
	}
	fn02 := func(args ...interface{}) []interface{} {
		val, _ := Maybe.Just(args[0]).ToInt()
		return SliceOf(val + 2)
	}
	fn03 := func(args ...interface{}) []interface{} {
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

	expectedinteger = 6
	assert.Equal(t, expectedinteger, Pipe(fn01, fn02, fn03)((0))[0])
}

func TestFPFunctions(t *testing.T) {
	expectedinteger := 0

	expectedinteger = -10
	assert.Equal(t, expectedinteger, Reduce(func(a, b interface{}) interface{} {
		aVal, _ := Maybe.Just(a).ToInt()
		bVal, _ := Maybe.Just(b).ToInt()
		return aVal - bVal
	}, 0, Map(func(a interface{}) interface{} {
		v, _ := Maybe.Just(a).ToInt()
		return (v + 1)
	}, Filter(func(a interface{}, i int) bool {
		v, _ := Maybe.Just(a).ToInt()
		return v >= 0
	}, -1, 0, 1, 2, 3)...)...))

	assert.Equal(t, []interface{}{1, 2, 3, 4, 5}, SortSlice(func(a, b interface{}) bool {
		aVal, _ := Maybe.Just(a).ToInt()
		bVal, _ := Maybe.Just(b).ToInt()
		return bVal-aVal > 0
	}, 1, 4, 5, 2, 3))

	var actualInt int
	// var actualInt2 int
	var actualMap map[interface{}]interface{}

	fib := func(n int) int {
		result, _ := Trampoline(func(input ...interface{}) ([]interface{}, bool, error) {
			n, _ := Maybe.Just(input[0]).ToInt()
			a, _ := Maybe.Just(input[1]).ToInt()
			b, _ := Maybe.Just(input[2]).ToInt()

			if n == 0 {
				return []interface{}{0, a, b}, true, nil
			}

			return []interface{}{n - 1, b, a + b}, false, nil
		}, n, 0, 1)

		val, _ := Maybe.Just(result[1]).ToInt()
		return val
	}

	actualInt = fib(6)
	assert.Equal(t, 8, actualInt)

	assert.Equal(t, []interface{}{3, 2, 1}, Compose(
		Reverse, SortIntAscending, DistinctRandom, SortIntAscending)(
		1, 1, 2, 1, 3, 1, 2, 1,
	),
	)
	assert.Equal(t, []interface{}{3, 2, 1}, Compose(
		Reverse, SortIntAscending, Distinct, SortIntAscending)(
		1, 1, 2, 1, 3, 1, 2, 1,
	),
	)
	assert.Equal(t, []interface{}{1, 2, 3}, Pipe(
		Reverse, SortIntAscending, Distinct, SortIntAscending)(
		1, 1, 2, 1, 3, 1, 2, 1,
	),
	)
	assert.Equal(t, []interface{}{1}, Compose(
		MakeNumericReturnForParam1ReturnBool1(IsZero),
	)(0))
	assert.Equal(t, []interface{}{1}, Compose(
		MakeNumericReturnForVariadicParamReturnBool1(IsDistinct),
	)(1, 2, 3, 4))
	/*
		assert.Equal(t, []interface{}{1}, Compose(
			MakeNumericReturnForVariadicParamReturnBool1(CurryParam1ForSlice1(IsEqual, []interface{}{1, 2, 3, 4})),
		)(1, 2, 3, 4))
	*/
	assert.Equal(t, []interface{}{1, 2, 3}, Dedupe(1, 1, 2, 2, 3, 3, 3, 3, 3))
	assert.Equal(t, []interface{}{1, 2, 3}, Difference([]interface{}{5, 1, 2, 3}, []interface{}{4, 5, 7, 8}))
	assert.Equal(t, []interface{}{1, 2, 3}, SortIntAscending(Distinct(1, 1, 2, 1, 3, 1, 2, 1)...))
	assert.Equal(t, true, IsDistinct(1, 2, 3))
	assert.Equal(t, false, IsDistinct(1, 1, 2, 1, 3, 1, 2, 1))
	assert.Equal(t, []interface{}{2, 3, 2}, DropEq(1, 1, 1, 2, 1, 3, 1, 2, 1))
	assert.Equal(t, []interface{}{1, 2, 1}, Drop(5, 1, 1, 2, 1, 3, 1, 2, 1))
	assert.Equal(t, []interface{}{1, 1, 2}, DropLast(5, 1, 1, 2, 1, 3, 1, 2, 1))
	assert.Equal(t, []interface{}{3, 1, 2, 1}, DropWhile(func(a interface{}) bool {
		aVal, _ := Maybe.Just(a).ToInt()
		return aVal < 3
	}, 1, 1, 2, 1, 3, 1, 2, 1))
	assert.Equal(t, true, IsEqual([]interface{}{1, 1, 2}, []interface{}{1, 1, 2}))
	assert.Equal(t, false, IsEqual([]interface{}{1, 1, 2}, []interface{}{1, 1, 3}))
	assert.Equal(t, true, IsEqualMap(map[interface{}]interface{}{2: 1, 3: 1, 1: 2}, map[interface{}]interface{}{1: 2, 2: 1, 3: 1}))
	assert.Equal(t, false, IsEqualMap(map[interface{}]interface{}{2: 1, 3: 1, 1: 2}, map[interface{}]interface{}{1: 1, 2: 1, 3: 2}))
	assert.Equal(t, true, Every(func(a interface{}) bool {
		aVal, _ := Maybe.Just(a).ToInt()
		return aVal%2 == 0
	}, 8, 2, 10, 4))
	assert.Equal(t, false, Every(func(a interface{}) bool {
		aVal, _ := Maybe.Just(a).ToInt()
		return aVal%2 == 0
	}, 8, 3, 10, 4))
	assert.Equal(t, true, Exists(10, 8, 3, 10, 4))
	assert.Equal(t, false, Exists(9, 8, 3, 10, 4))
	assert.Equal(t, true, IsZero(0))
	assert.Equal(t, false, IsZero(nil))
	assert.Equal(t, false, IsZero(1))
	assert.Equal(t, false, IsZero(""))
	assert.Equal(t, []interface{}{1, 3, 2}, Intersection([]interface{}{5, 1, 3, 2, 8}, []interface{}{7, 6, 4, 3, 1, 2}))
	assert.Equal(t, []interface{}{2, 5}, SortIntAscending(Keys(IntersectionMapByKey(map[interface{}]interface{}{2: 11, 5: 11, 1: 12}, map[interface{}]interface{}{41: 1, 2: 77, 42: 1, 5: 66, 43: 2}))...))
	assert.Equal(t, []interface{}{1, 2, 3}, SortIntAscending(Keys(map[interface{}]interface{}{2: 8, 1: 5, 3: 4})...))
	assert.Equal(t, []interface{}{4, 5, 8}, SortIntAscending(Values(map[interface{}]interface{}{2: 8, 1: 5, 3: 4})...))
	assert.Equal(t, []interface{}{5, 8, 8}, Minus([]interface{}{5, 1, 8, 3, 2, 8}, []interface{}{7, 6, 4, 3, 1, 2}))
	assert.Equal(t, []interface{}{7, 6, 4}, Minus([]interface{}{7, 6, 4, 3, 1, 2}, []interface{}{5, 1, 8, 3, 2, 8}))
	assert.Equal(t, []interface{}{1}, SortIntAscending(Keys(MinusMapByKey(map[interface{}]interface{}{2: 11, 5: 11, 1: 12}, map[interface{}]interface{}{41: 1, 2: 77, 42: 1, 5: 66, 43: 2}))...))
	assert.Equal(t, []interface{}{41, 42, 43}, SortIntAscending(Keys(MinusMapByKey(map[interface{}]interface{}{41: 1, 2: 77, 42: 1, 5: 66, 43: 2}, map[interface{}]interface{}{2: 11, 5: 11, 1: 12}))...))
	assert.Equal(t, true, IsSubset([]interface{}{1, 2, 3}, []interface{}{4, 5, 1, 2, 3, 6}))
	assert.Equal(t, true, IsSubset([]interface{}{1, 2, 2, 3}, []interface{}{4, 5, 1, 2, 3, 6}))
	assert.Equal(t, false, IsSubset([]interface{}{5, 1, 8, 3, 2, 8}, []interface{}{7, 6, 4, 3, 1, 2}))
	assert.Equal(t, true, IsSuperset([]interface{}{4, 5, 1, 2, 3, 6}, []interface{}{1, 2, 3}))
	assert.Equal(t, true, IsSuperset([]interface{}{4, 5, 1, 2, 3, 6}, []interface{}{1, 2, 2, 3}))
	assert.Equal(t, false, IsSuperset([]interface{}{5, 1, 8, 3, 2, 8}, []interface{}{7, 6, 4, 3, 1, 2}))
	assert.Equal(t, true, IsSubsetMapByKey(map[interface{}]interface{}{1: 3, 2: 4}, map[interface{}]interface{}{5: 6, 1: 4, 3: 5, 2: 7, 8: 9}))
	assert.Equal(t, false, IsSubsetMapByKey(map[interface{}]interface{}{5: 6, 1: 4, 3: 5, 2: 7, 8: 9}, map[interface{}]interface{}{7: 8, 6: 9, 4: 10, 3: 11, 1: 13, 2: 12}))
	assert.Equal(t, true, IsSupersetMapByKey(map[interface{}]interface{}{5: 6, 1: 4, 3: 5, 2: 7, 8: 9}, map[interface{}]interface{}{1: 3, 2: 4}))
	assert.Equal(t, false, IsSupersetMapByKey(map[interface{}]interface{}{5: 6, 1: 4, 3: 5, 2: 7, 8: 9}, map[interface{}]interface{}{7: 8, 6: 9, 4: 10, 3: 11, 1: 13, 2: 12}))
	assert.Equal(t, []interface{}{1, 2, 3, 4, 5, 6, 7, 8}, SortIntAscending(Union([]interface{}{5, 1, 3, 2, 8}, []interface{}{7, 6, 4, 3, 1, 2})...))
	actualMap = Merge(map[interface{}]interface{}{2: 11, 5: 11, 1: 12}, map[interface{}]interface{}{41: 1, 42: 1, 43: 2})
	assert.Equal(t, []interface{}{1, 2, 5, 41, 42, 43}, SortIntAscending(Keys(actualMap)...))
	assert.Equal(t, []interface{}{1, 1, 2, 11, 11, 12}, SortIntAscending(Values(actualMap)...))
	assert.Equal(t, []interface{}{2, 1, 2}, Take(3, 2, 1, 2, 1, 3, 1, 2, 1))
	assert.Equal(t, []interface{}{1, 2, 1}, TakeLast(3, 2, 1, 2, 1, 3, 1, 2, 1))
	assert.Equal(t, []interface{}{1, 2, 1, 3, 1, 2, 1}, Tail(2, 1, 2, 1, 3, 1, 2, 1))
	assert.Equal(t, 2, Head(2, 1, 2, 1, 3, 1, 2, 1))
	assert.Equal(t, []interface{}{2, 1, 2, 1, 3, 1, 2, 1}, Flatten([]interface{}{2, 1, 2}, []interface{}{1, 3, 1}, []interface{}{2, 1}))
	assert.Equal(t, []interface{}{3, 2, 1, 2, 1, 3, 1, 2, 1}, Prepend(3, []interface{}{2, 1, 2, 1, 3, 1, 2, 1}))
	assert.Equal(t, []interface{}{5, 4, 3, 2, 1}, Reverse(1, 2, 3, 4, 5))
	assert.Equal(t, []interface{}{2, 3, 4, 5, 6}, PMap(func(a interface{}) interface{} {
		aVal, _ := Maybe.Just(a).ToInt()
		return aVal + 1
	}, nil, 1, 2, 3, 4, 5))
	assert.Equal(t, []interface{}{2, 3, 4, 5, 6}, PMap(func(a interface{}) interface{} {
		aVal, _ := Maybe.Just(a).ToInt()
		return aVal + 1
	}, &PMapOption{FixedPool: 3, RandomOrder: false}, 1, 2, 3, 4, 5))
	assert.Equal(t, []interface{}{2, 3, 4, 5, 6}, SortIntAscending(PMap(func(a interface{}) interface{} {
		aVal, _ := Maybe.Just(a).ToInt()
		return aVal + 1
	}, &PMapOption{FixedPool: 3, RandomOrder: true}, 1, 2, 3, 4, 5)...))
	assert.Equal(t, true, Some(func(a interface{}) bool {
		aVal, _ := Maybe.Just(a).ToInt()
		return aVal%2 == 0
	}, 1, 2, 3, 4, 5))
	assert.Equal(t, false, Some(func(a interface{}) bool {
		aVal, _ := Maybe.Just(a).ToInt()
		return aVal%2 == 0
	}, 1, 3, 5, 7, 9))
	assert.Equal(t, map[interface{}]interface{}{1: "a", 2: "b", 3: "c"}, Zip([]interface{}{1, 2, 3}, []interface{}{"a", "b", "c"}))
	assert.Equal(t, [][]interface{}{{1, 3, 5, 7}, {2, 4, 6, 8}}, Partition(func(a interface{}) bool {
		aVal, _ := Maybe.Just(a).ToInt()
		return aVal%2 == 1
	}, 1, 2, 3, 4, 5, 6, 7, 8))
	assert.Equal(t, [][]interface{}{{1, 2, 3}, {4, 5, 6}, {7, 8}}, SplitEvery(3, 1, 2, 3, 4, 5, 6, 7, 8))
	assert.Equal(t, map[interface{}][]interface{}{1: {1, 3, 5, 7}, 0: {2, 4, 6, 8}}, GroupBy(func(a interface{}) interface{} {
		aVal, _ := Maybe.Just(a).ToInt()
		return aVal % 2
	}, 1, 2, 3, 4, 5, 6, 7, 8))
	assert.Equal(t, []interface{}{1, 2}, UniqBy(func(a interface{}) interface{} {
		aVal, _ := Maybe.Just(a).ToInt()
		return aVal % 2
	}, 1, 2, 3, 4, 5, 6, 7, 8))
}

func TestVariadic(t *testing.T) {
	assert.Equal(t, []interface{}{28}, Compose(
		MakeVariadicParam1(func(arg1 interface{}) []interface{} {
			return []interface{}{arg1.(int)}
		}),
		MakeVariadicReturn1(func(args ...interface{}) interface{} {
			return args[0].(int) + args[1].(int)
		}),
		MakeVariadicParam2(func(arg1 interface{}, arg2 interface{}) []interface{} {
			return []interface{}{arg1.(int), arg2.(int)}
		}),
		MakeVariadicReturn2(func(args ...interface{}) (interface{}, interface{}) {
			return args[0].(int), args[1].(int) + args[2].(int)
		}),
		MakeVariadicParam3(func(arg1 interface{}, arg2 interface{}, arg3 interface{}) []interface{} {
			return []interface{}{arg1.(int), arg2.(int), arg3.(int)}
		}),
		MakeVariadicReturn3(func(args ...interface{}) (interface{}, interface{}, interface{}) {
			return args[0].(int), args[1].(int), args[2].(int) + args[3].(int)
		}),
		MakeVariadicParam4(func(arg1 interface{}, arg2 interface{}, arg3 interface{}, arg4 interface{}) []interface{} {
			return []interface{}{arg1.(int), arg2.(int), arg3.(int), arg4.(int)}
		}),
		MakeVariadicReturn4(func(args ...interface{}) (interface{}, interface{}, interface{}, interface{}) {
			return args[0].(int), args[1].(int), args[2].(int), args[3].(int) + args[4].(int)
		}),
		MakeVariadicParam5(func(arg1 interface{}, arg2 interface{}, arg3 interface{}, arg4 interface{}, arg5 interface{}) []interface{} {
			return []interface{}{arg1.(int), arg2.(int), arg3.(int), arg4.(int), arg5.(int)}
		}),
		MakeVariadicReturn5(func(args ...interface{}) (interface{}, interface{}, interface{}, interface{}, interface{}) {
			return args[0].(int), args[1].(int), args[2].(int), args[3].(int), args[4].(int) + args[5].(int)
		}),
		MakeVariadicParam6(func(arg1 interface{}, arg2 interface{}, arg3 interface{}, arg4 interface{}, arg5 interface{}, arg6 interface{}) []interface{} {
			return []interface{}{arg1.(int), arg2.(int), arg3.(int), arg4.(int), arg5.(int), arg6.(int)}
		}),
		MakeVariadicReturn6(func(args ...interface{}) (interface{}, interface{}, interface{}, interface{}, interface{}, interface{}) {
			return args[0].(int), args[1].(int), args[2].(int), args[3].(int), args[4].(int), args[5].(int) + args[6].(int)
		}),
	)(1, 2, 3, 4, 5, 6, 7),
	)
	assert.Equal(t, 28, CurryParam6(func(arg1 interface{}, arg2 interface{}, arg3 interface{}, arg4 interface{}, arg5 interface{}, arg6 interface{}, arg7 ...interface{}) interface{} {
		return arg7[0].(int) + CurryParam5(func(arg1 interface{}, arg2 interface{}, arg3 interface{}, arg4 interface{}, arg5 interface{}, arg6 ...interface{}) interface{} {
			return arg6[0].(int) + CurryParam4(func(arg1 interface{}, arg2 interface{}, arg3 interface{}, arg4 interface{}, arg5 ...interface{}) interface{} {
				return arg5[0].(int) + CurryParam3(func(arg1 interface{}, arg2 interface{}, arg3 interface{}, arg4 ...interface{}) interface{} {
					return arg4[0].(int) + CurryParam2(func(arg1 interface{}, arg2 interface{}, arg3 ...interface{}) interface{} {
						return arg3[0].(int) + CurryParam1(func(arg1 interface{}, arg2 ...interface{}) interface{} {
							return arg2[0].(int) + arg1.(int)
						}, arg1)(arg2).(int)
					}, arg1, arg2)(arg3).(int)
				}, arg1, arg2, arg3)(arg4).(int)
			}, arg1, arg2, arg3, arg4)(arg5).(int)
		}, arg1, arg2, arg3, arg4, arg5)(arg6).(int)
	}, 1, 2, 3, 4, 5, 6)(7))
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
	compTypeA := DefProduct(reflect.Int, reflect.String)
	compTypeB := DefProduct(reflect.String)
	myType := DefSum(NilType, compTypeA, compTypeB)

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
	compTypeA := DefProduct(reflect.Int, reflect.String)
	compTypeB := DefProduct(reflect.String, reflect.String)
	myType := DefSum(NilType, compTypeA, compTypeB)

	assert.Equal(t, true, compTypeA.Matches(1, "3"))
	assert.Equal(t, false, compTypeA.Matches(1, 3))
	assert.Equal(t, true, myType.Matches(nil))
	assert.Equal(t, true, myType.Matches(1, "3"))
	assert.Equal(t, true, myType.Matches("1", "3"))
	assert.Equal(t, false, myType.Matches(1, 3))

	patterns := []Pattern{
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
	pm := DefPattern(patterns...)
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
