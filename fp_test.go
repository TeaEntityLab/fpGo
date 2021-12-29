package fpgo

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompose(t *testing.T) {
	var expectedinteger int

	var fn01 = func(args ...int) []int {
		val := args[0]
		return SliceOf(val + 1)
	}
	var fn02 = func(args ...int) []int {
		val := args[0]
		return SliceOf(val + 2)
	}
	var fn03 = func(args ...int) []int {
		val := args[0]
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
	assert.Equal(t, expectedinteger, Reduce(func(a, b int) int { return a - b }, 0, Map(func(a int) int { return a + 1 }, Filter(func(a int, i int) bool { return a >= 0 }, -1, 0, 1, 2, 3)...)...))

	assert.Equal(t, []int{1, 2, 3, 4, 5}, SortSlice(func(a, b int) bool { return b-a > 0 }, 1, 4, 5, 2, 3))

	var actualInt int
	var actualInt2 int
	var actualMap map[int]int

	fib := func(n int) int {

		result, _ := Trampoline(func(input ...int) ([]int, bool, error) {
			n := input[0]
			a := input[1]
			b := input[2]

			if n == 0 {
				return []int{0, a, b}, true, nil
			}

			return []int{n - 1, b, a + b}, false, nil
		}, n, 0, 1)

		return result[1]
	}

	actualInt = fib(6)
	assert.Equal(t, 8, actualInt)

	assert.Equal(t, []int{3, 2}, Compose[int](
		CurryParam1[int, []int, int](DropLast[int], 1), Reverse[int], SortOrderedAscending[int], DistinctRandom[int], SortOrderedAscending[int])(
		1, 1, 2, 1, 3, 1, 2, 1,
	),
	)
	assert.Equal(t, []int{3, 2}, Compose[int](
		CurryParam1[int, []int, int](DropLast[int], 1), Reverse[int], SortOrderedAscending[int], Distinct[int], SortOrderedAscending[int])(
		1, 1, 2, 1, 3, 1, 2, 1,
	),
	)
	assert.Equal(t, []int{1, 2, 3}, Pipe[int](
		CurryParam1[int, []int, int](DropLast[int], 1), Reverse[int], SortOrderedAscending[int], Distinct[int], SortOrderedAscending[int])(
		1, 1, 2, 1, 3, 1, 2, 1,
	),
	)
	assert.Equal(t, []int{1}, Compose(
		MakeNumericReturnForParam1ReturnBool1[int, int](IsNeg[int]),
		func(val ...int) []int {
			return SliceOf(Reduce(func(a, b int) int { return a - b }, 0, val...))
		},
	)(1, 2, 3, 4))
	assert.Equal(t, []int{1}, Compose(
		MakeNumericReturnForVariadicParamReturnBool1[int, int](IsDistinct[int]),
	)(1, 2, 3, 4))
	assert.Equal(t, []int{1}, Compose(
		MakeNumericReturnForVariadicParamReturnBool1[int, int](CurryParam1ForSlice1(IsEqual[int], []int{1, 2, 3, 4})),
	)(1, 2, 3, 4))
	assert.Equal(t, []int{1, 2, 3}, Dedupe(1, 1, 2, 2, 3, 3, 3, 3, 3))
	assert.Equal(t, []int{1, 2, 3}, Difference([]int{5, 1, 2, 3}, []int{4, 5, 7, 8}))
	assert.Equal(t, []int{1, 2, 3}, SortOrderedAscending(Distinct(1, 1, 2, 1, 3, 1, 2, 1)...))
	assert.Equal(t, true, IsDistinct(1, 2, 3))
	assert.Equal(t, false, IsDistinct(1, 1, 2, 1, 3, 1, 2, 1))
	assert.Equal(t, []int{2, 3, 2}, DropEq(1, 1, 1, 2, 1, 3, 1, 2, 1))
	assert.Equal(t, []int{1, 2, 1}, Drop(5, 1, 1, 2, 1, 3, 1, 2, 1))
	assert.Equal(t, []int{1, 1, 2}, DropLast(5, 1, 1, 2, 1, 3, 1, 2, 1))
	assert.Equal(t, []int{3, 1, 2, 1}, DropWhile(func(a int) bool { return a < 3 }, 1, 1, 2, 1, 3, 1, 2, 1))
	assert.Equal(t, true, IsEqual([]int{1, 1, 2}, []int{1, 1, 2}))
	assert.Equal(t, false, IsEqual([]int{1, 1, 2}, []int{1, 1, 3}))
	assert.Equal(t, true, IsEqualMap(map[int]int{2: 1, 3: 1, 1: 2}, map[int]int{1: 2, 2: 1, 3: 1}))
	assert.Equal(t, false, IsEqualMap(map[int]int{2: 1, 3: 1, 1: 2}, map[int]int{1: 1, 2: 1, 3: 2}))
	assert.Equal(t, true, Every(func(a int) bool { return a%2 == 0 }, 8, 2, 10, 4))
	assert.Equal(t, false, Every(func(a int) bool { return a%2 == 0 }, 8, 3, 10, 4))
	assert.Equal(t, true, Exists(10, 8, 3, 10, 4))
	assert.Equal(t, false, Exists(9, 8, 3, 10, 4))
	assert.Equal(t, []int{1, 3, 2}, Intersection([]int{5, 1, 3, 2, 8}, []int{7, 6, 4, 3, 1, 2}))
	assert.Equal(t, []int{2, 5}, SortOrderedAscending(Keys(IntersectionMapByKey(map[int]int{2: 11, 5: 11, 1: 12}, map[int]int{41: 1, 2: 77, 42: 1, 5: 66, 43: 2}))...))
	assert.Equal(t, []int{1, 2, 3}, SortOrderedAscending(Keys(map[int]int{2: 8, 1: 5, 3: 4})...))
	assert.Equal(t, []int{4, 5, 8}, SortOrderedAscending(Values(map[int]int{2: 8, 1: 5, 3: 4})...))
	assert.Equal(t, []int{5, 8, 8}, Minus([]int{5, 1, 8, 3, 2, 8}, []int{7, 6, 4, 3, 1, 2}))
	assert.Equal(t, []int{7, 6, 4}, Minus([]int{7, 6, 4, 3, 1, 2}, []int{5, 1, 8, 3, 2, 8}))
	assert.Equal(t, []int{1}, SortOrderedAscending(Keys(MinusMapByKey(map[int]int{2: 11, 5: 11, 1: 12}, map[int]int{41: 1, 2: 77, 42: 1, 5: 66, 43: 2}))...))
	assert.Equal(t, []int{41, 42, 43}, SortOrderedAscending(Keys(MinusMapByKey(map[int]int{41: 1, 2: 77, 42: 1, 5: 66, 43: 2}, map[int]int{2: 11, 5: 11, 1: 12}))...))
	assert.Equal(t, true, IsSubset([]int{1, 2, 3}, []int{4, 5, 1, 2, 3, 6}))
	assert.Equal(t, true, IsSubset([]int{1, 2, 2, 3}, []int{4, 5, 1, 2, 3, 6}))
	assert.Equal(t, false, IsSubset([]int{5, 1, 8, 3, 2, 8}, []int{7, 6, 4, 3, 1, 2}))
	assert.Equal(t, true, IsSuperset([]int{4, 5, 1, 2, 3, 6}, []int{1, 2, 3}))
	assert.Equal(t, true, IsSuperset([]int{4, 5, 1, 2, 3, 6}, []int{1, 2, 2, 3}))
	assert.Equal(t, false, IsSuperset([]int{5, 1, 8, 3, 2, 8}, []int{7, 6, 4, 3, 1, 2}))
	assert.Equal(t, true, IsSubsetMapByKey(map[int]int{1: 3, 2: 4}, map[int]int{5: 6, 1: 4, 3: 5, 2: 7, 8: 9}))
	assert.Equal(t, false, IsSubsetMapByKey(map[int]int{5: 6, 1: 4, 3: 5, 2: 7, 8: 9}, map[int]int{7: 8, 6: 9, 4: 10, 3: 11, 1: 13, 2: 12}))
	assert.Equal(t, true, IsSupersetMapByKey(map[int]int{5: 6, 1: 4, 3: 5, 2: 7, 8: 9}, map[int]int{1: 3, 2: 4}))
	assert.Equal(t, false, IsSupersetMapByKey(map[int]int{5: 6, 1: 4, 3: 5, 2: 7, 8: 9}, map[int]int{7: 8, 6: 9, 4: 10, 3: 11, 1: 13, 2: 12}))
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8}, SortOrderedAscending(Union([]int{5, 1, 3, 2, 8}, []int{7, 6, 4, 3, 1, 2})...))
	assert.Equal(t, 3, Max(2, 1, 2, 1, 3, 1, 2, 1))
	assert.Equal(t, 1, Min(2, 1, 2, 1, 3, 1, 2, 1))
	actualInt, actualInt2 = MinMax(2, 1, 2, 1, 3, 1, 2, 1)
	assert.Equal(t, []int{1, 3}, []int{actualInt, actualInt2})
	actualMap = Merge(map[int]int{2: 11, 5: 11, 1: 12}, map[int]int{41: 1, 42: 1, 43: 2})
	assert.Equal(t, []int{1, 2, 5, 41, 42, 43}, SortOrderedAscending(Keys(actualMap)...))
	assert.Equal(t, []int{1, 1, 2, 11, 11, 12}, SortOrderedAscending(Values(actualMap)...))
	assert.Equal(t, true, IsNeg(-1))
	assert.Equal(t, false, IsPos(-1))
	assert.Equal(t, false, IsNeg(1))
	assert.Equal(t, true, IsPos(1))
	assert.Equal(t, false, IsNeg(0))
	assert.Equal(t, false, IsPos(0))
	assert.Equal(t, true, IsZero(0))
	assert.Equal(t, []int{2, 1, 2}, Take(3, 2, 1, 2, 1, 3, 1, 2, 1))
	assert.Equal(t, []int{1, 2, 1}, TakeLast(3, 2, 1, 2, 1, 3, 1, 2, 1))
	assert.Equal(t, []int{1, 2, 1, 3, 1, 2, 1}, Tail(2, 1, 2, 1, 3, 1, 2, 1))
	assert.Equal(t, 2, Head(2, 1, 2, 1, 3, 1, 2, 1))
	assert.Equal(t, []int{2, 1, 2, 1, 3, 1, 2, 1}, Flatten([]int{2, 1, 2}, []int{1, 3, 1}, []int{2, 1}))
	assert.Equal(t, []int{3, 2, 1, 2, 1, 3, 1, 2, 1}, Prepend(3, []int{2, 1, 2, 1, 3, 1, 2, 1}))
	assert.Equal(t, []int{1, 2, 3, 4}, Range(1, 5))
	assert.Equal(t, []int{1, 3, 5, 7, 9}, Range(1, 10, 2))
	assert.Equal(t, []int{5, 4, 3, 2, 1}, Reverse(1, 2, 3, 4, 5))
	assert.Equal(t, []int{2, 3, 4, 5, 6}, PMap(func(a int) int { return a + 1 }, nil, 1, 2, 3, 4, 5))
	assert.Equal(t, []int{2, 3, 4, 5, 6}, PMap(func(a int) int { return a + 1 }, &PMapOption{FixedPool: 3, RandomOrder: false}, 1, 2, 3, 4, 5))
	assert.Equal(t, []int{2, 3, 4, 5, 6}, SortOrderedAscending(PMap(func(a int) int { return a + 1 }, &PMapOption{FixedPool: 3, RandomOrder: true}, 1, 2, 3, 4, 5)...))
	assert.Equal(t, true, Some(func(a int) bool { return a%2 == 0 }, 1, 2, 3, 4, 5))
	assert.Equal(t, false, Some(func(a int) bool { return a%2 == 0 }, 1, 3, 5, 7, 9))
	assert.Equal(t, map[int]string{1: "a", 2: "b", 3: "c"}, Zip([]int{1, 2, 3}, []string{"a", "b", "c"}))
	assert.Equal(t, [][]int{{1, 3, 5, 7}, {2, 4, 6, 8}}, Partition(func(a int) bool { return a%2 == 1 }, 1, 2, 3, 4, 5, 6, 7, 8))
	assert.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8}}, SplitEvery(3, 1, 2, 3, 4, 5, 6, 7, 8))
	assert.Equal(t, map[int][]int{1: {1, 3, 5, 7}, 0: {2, 4, 6, 8}}, GroupBy(func(a int) int { return a % 2 }, 1, 2, 3, 4, 5, 6, 7, 8))
	assert.Equal(t, []int{1, 2}, UniqBy(func(a int) int { return a % 2 }, 1, 2, 3, 4, 5, 6, 7, 8))
}

func TestVariadic(t *testing.T) {
	assert.Equal(t, []int{28}, Compose(
		MakeVariadicParam1(func(arg1 int) []int {
			return []int{arg1}
		}),
		MakeVariadicReturn1(func(args ...int) int {
			return args[0] + args[1]
		}),
		MakeVariadicParam2(func(arg1 int, arg2 int) []int {
			return []int{arg1, arg2}
		}),
		MakeVariadicReturn2(func(args ...int) (int, int) {
			return args[0], args[1] + args[2]
		}),
		MakeVariadicParam3(func(arg1 int, arg2 int, arg3 int) []int {
			return []int{arg1, arg2, arg3}
		}),
		MakeVariadicReturn3(func(args ...int) (int, int, int) {
			return args[0], args[1], args[2] + args[3]
		}),
		MakeVariadicParam4(func(arg1 int, arg2 int, arg3 int, arg4 int) []int {
			return []int{arg1, arg2, arg3, arg4}
		}),
		MakeVariadicReturn4(func(args ...int) (int, int, int, int) {
			return args[0], args[1], args[2], args[3] + args[4]
		}),
		MakeVariadicParam5(func(arg1 int, arg2 int, arg3 int, arg4 int, arg5 int) []int {
			return []int{arg1, arg2, arg3, arg4, arg5}
		}),
		MakeVariadicReturn5(func(args ...int) (int, int, int, int, int) {
			return args[0], args[1], args[2], args[3], args[4] + args[5]
		}),
		MakeVariadicParam6(func(arg1 int, arg2 int, arg3 int, arg4 int, arg5 int, arg6 int) []int {
			return []int{arg1, arg2, arg3, arg4, arg5, arg6}
		}),
		MakeVariadicReturn6(func(args ...int) (int, int, int, int, int, int) {
			return args[0], args[1], args[2], args[3], args[4], args[5] + args[6]
		}),
	)(1, 2, 3, 4, 5, 6, 7),
	)
	assert.Equal(t, 28, CurryParam6(func(arg1 int, arg2 int, arg3 int, arg4 int, arg5 int, arg6 int, arg7 ...int) int {
		return arg7[0] + CurryParam5(func(arg1 int, arg2 int, arg3 int, arg4 int, arg5 int, arg6 ...int) int {
			return arg6[0] + CurryParam4(func(arg1 int, arg2 int, arg3 int, arg4 int, arg5 ...int) int {
				return arg5[0] + CurryParam3(func(arg1 int, arg2 int, arg3 int, arg4 ...int) int {
					return arg4[0] + CurryParam2(func(arg1 int, arg2 int, arg3 ...int) int {
						return arg3[0] + CurryParam1(func(arg1 int, arg2 ...int) int {
							return arg2[0] + arg1
						}, arg1)(arg2)
					}, arg1, arg2)(arg3)
				}, arg1, arg2, arg3)(arg4)
			}, arg1, arg2, arg3, arg4)(arg5)
		}, arg1, arg2, arg3, arg4, arg5)(arg6)
	}, 1, 2, 3, 4, 5, 6)(7))
}

func TestCurry(t *testing.T) {
	c := CurryNew(func(c *CurryDef[interface{}, interface{}], args ...interface{}) interface{} {
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
