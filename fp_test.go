package fpgo

import (
	"errors"
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

	fn01 := func(args ...int) []int {
		val := args[0]
		return SliceOf(val + 1)
	}
	fn02 := func(args ...int) []int {
		val := args[0]
		return SliceOf(val + 2)
	}
	fn03 := func(args ...int) []int {
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
	result, err := Zip([]int{1, 2, 3}, []string{"a", "b", "c"})
	assert.Nil(t, err)
	assert.Equal(t, map[int]string{1: "a", 2: "b", 3: "c"}, result)
	assert.Equal(t, [][]int{{1, 3, 5, 7}, {2, 4, 6, 8}}, Partition(func(a int) bool { return a%2 == 1 }, 1, 2, 3, 4, 5, 6, 7, 8))
	assert.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8}}, SplitEvery(3, 1, 2, 3, 4, 5, 6, 7, 8))
	assert.Equal(t, map[int][]int{1: {1, 3, 5, 7}, 0: {2, 4, 6, 8}}, GroupBy(func(a int) int { return a % 2 }, 1, 2, 3, 4, 5, 6, 7, 8))
	assert.Equal(t, []int{1, 2}, UniqBy(func(a int) int { return a % 2 }, 1, 2, 3, 4, 5, 6, 7, 8))
}

func TestIsNeg(t *testing.T) {
	assert.Equal(t, true, IsNeg(-1))
	assert.Equal(t, false, IsNeg(1))
	assert.Equal(t, false, IsNeg(0))
	assert.Equal(t, true, IsNeg(int32(-1)))
	assert.Equal(t, false, IsNeg(int64(1)))
	assert.Equal(t, true, IsNeg(float64(-1.5)))
}

func TestIsPos(t *testing.T) {
	assert.Equal(t, true, IsPos(1))
	assert.Equal(t, false, IsPos(-1))
	assert.Equal(t, false, IsPos(0))
	assert.Equal(t, true, IsPos(int32(1)))
	assert.Equal(t, false, IsPos(int64(-1)))
	assert.Equal(t, true, IsPos(float64(1.5)))
}

func TestIsZero(t *testing.T) {
	assert.Equal(t, true, IsZero(0))
	assert.Equal(t, false, IsZero(1))
	assert.Equal(t, true, IsZero(int32(0)))
	assert.Equal(t, false, IsZero(int64(1)))
	assert.Equal(t, true, IsZero(float64(0)))
}

func TestTrampoline(t *testing.T) {
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

	assert.Equal(t, 0, fib(0))
	assert.Equal(t, 1, fib(1))
	assert.Equal(t, 1, fib(2))
	assert.Equal(t, 2, fib(3))
	assert.Equal(t, 55, fib(10))
}

func TestConcat(t *testing.T) {
	result := Concat([]int{1, 2}, []int{3, 4}, []int{5, 6})
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)

	result = Concat([]int{1, 2})
	assert.Equal(t, []int{1, 2}, result)

	result = Concat([]int{1, 2}, []int{})
	assert.Equal(t, []int{1, 2}, result)
}

func TestReject(t *testing.T) {
	result := Reject(func(a int, i int) bool { return a > 2 }, 1, 2, 3, 4, 5)
	assert.Equal(t, []int{1, 2}, result)
}

func TestMapIndexed(t *testing.T) {
	result := MapIndexed(func(a int, i int) int { return a + i }, 1, 2, 3)
	assert.Equal(t, []int{1, 3, 5}, result)
}

func TestReduceIndexed(t *testing.T) {
	result := ReduceIndexed(func(acc int, val int, i int) int { return acc + val + i }, 0, 1, 2, 3)
	assert.Equal(t, 9, result)
}

func TestFlattenAndPrepend(t *testing.T) {
	result := Flatten([]int{1, 2}, []int{3, 4}, []int{5, 6})
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)

	result = Flatten([]int{1, 2})
	assert.Equal(t, []int{1, 2}, result)

	result = Prepend(0, []int{1, 2, 3})
	assert.Equal(t, []int{0, 1, 2, 3}, result)
}

func TestPartition(t *testing.T) {
	result := Partition(func(a int) bool { return a%2 == 0 }, 1, 2, 3, 4)
	assert.Equal(t, [][]int{{2, 4}, {1, 3}}, result)
}

func TestSplitEvery(t *testing.T) {
	result := SplitEvery(2, 1, 2, 3, 4, 5)
	assert.Equal(t, [][]int{{1, 2}, {3, 4}, {5}}, result)
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

func TestZipErrors(t *testing.T) {
	result, err := Zip([]int{1, 2, 3}, []string{"a", "b", "c"})
	assert.Nil(t, err)
	assert.Equal(t, map[int]string{1: "a", 2: "b", 3: "c"}, result)

	_, err = Zip([]int{}, []string{"a"})
	assert.Equal(t, ErrZipEmptyList, err)

	_, err = Zip([]int{1, 2}, []string{"a"})
	assert.Equal(t, ErrZipLengthMismatch, err)

	_, err = Zip([]int{1}, []string{})
	assert.Equal(t, ErrZipEmptyList, err)
}

func TestPMapNoOrderConcurrency(t *testing.T) {
	result := PMap(func(a int) int { return a + 1 }, &PMapOption{FixedPool: 3, RandomOrder: true}, 1, 2, 3, 4, 5)
	sorted := SortOrderedAscending(result...)
	assert.Equal(t, []int{2, 3, 4, 5, 6}, sorted)

	largeList := make([]int, 1000)
	for i := range largeList {
		largeList[i] = i
	}
	result = PMap(func(a int) int { return a * 2 }, &PMapOption{FixedPool: 10, RandomOrder: true}, largeList...)
	assert.Len(t, result, 1000)
}

func TestCompDataTypeValidation(t *testing.T) {
	compTypeA := DefProduct(reflect.Int, reflect.String)

	data := NewCompData(compTypeA, 1, "hello")
	assert.NotNil(t, data)

	nilData := NewCompData(compTypeA, "wrong", "types")
	assert.Nil(t, nilData)
}

func TestDefProduct(t *testing.T) {
	productType := DefProduct(reflect.Int, reflect.String)
	assert.True(t, productType.Matches(1, "test"))
	assert.False(t, productType.Matches(1, 2))
	assert.False(t, productType.Matches(1))
	assert.False(t, productType.Matches(1, "test", "extra"))
}

func TestNilTypeMatches(t *testing.T) {
	var nilValue *int = nil
	assert.True(t, NilType.Matches(nilValue))
	assert.False(t, NilType.Matches(1))
	assert.False(t, NilType.Matches("test"))
	assert.False(t, NilType.Matches(nil, nil))
}

func TestMergeForInterface(t *testing.T) {
	result := MergeForInterface[string](map[interface{}]string{1: "a", 2: "b"}, map[interface{}]string{3: "c", 4: "d"})
	assert.Equal(t, "a", result[1])
	assert.Equal(t, "b", result[2])
	assert.Equal(t, "c", result[3])
	assert.Equal(t, "d", result[4])

	result = MergeForInterface[string](nil, map[interface{}]string{3: "c"})
	assert.Equal(t, "c", result[3])

	result = MergeForInterface[string](map[interface{}]string{1: "a"}, nil)
	assert.Equal(t, "a", result[1])

	result = MergeForInterface[string](nil, nil)
	assert.Empty(t, result)
}

func TestMinusForInterface(t *testing.T) {
	result := MinusForInterface([]interface{}{1, 2, 3, 4}, []interface{}{2, 3})
	assert.Len(t, result, 2)

	result = MinusForInterface([]interface{}{}, []interface{}{1, 2})
	assert.Empty(t, result)

	result = MinusForInterface([]interface{}{1, 2}, []interface{}{})
	assert.Len(t, result, 2)
}

func TestKeysForInterface(t *testing.T) {
	result := KeysForInterface[string](map[interface{}]string{1: "a", "key": "b", 3: "c"})
	assert.Len(t, result, 3)
}

func TestValuesForInterface(t *testing.T) {
	result := ValuesForInterface[string](map[interface{}]string{1: "a", "key": "b", 3: "c"})
	assert.Len(t, result, 3)
}

func TestDistinctForInterface(t *testing.T) {
	result := DistinctForInterface(1, 2, 1, 3, 2, 1)
	assert.Equal(t, []interface{}{1, 2, 3}, result)

	result = DistinctForInterface()
	assert.Empty(t, result)
}

func TestExistsForInterface(t *testing.T) {
	assert.True(t, ExistsForInterface(2, 1, 2, 3))
	assert.False(t, ExistsForInterface(4, 1, 2, 3))
}

func TestIntersectionForInterface(t *testing.T) {
	result := IntersectionForInterface([]interface{}{1, 2, 3}, []interface{}{2, 3, 4})
	assert.Equal(t, []interface{}{2, 3}, result)

	result = IntersectionForInterface([]interface{}{})
	assert.Empty(t, result)
}

func TestIntersectionMapByKeyForInterface(t *testing.T) {
	result := IntersectionMapByKeyForInterface[string](map[interface{}]string{1: "a", 2: "b"}, map[interface{}]string{2: "c", 3: "d"})
	assert.Len(t, result, 1)
	assert.Equal(t, "b", result[2])
}

func TestIsSubsetForInterface(t *testing.T) {
	assert.True(t, IsSubsetForInterface([]interface{}{1, 2}, []interface{}{3, 1, 2, 4}))
	assert.False(t, IsSubsetForInterface([]interface{}{1, 2}, []interface{}{3, 4}))
	assert.False(t, IsSubsetForInterface([]interface{}{}, []interface{}{1, 2}))
}

func TestIsSupersetForInterface(t *testing.T) {
	assert.True(t, IsSupersetForInterface([]interface{}{1, 2, 3}, []interface{}{1, 2}))
	assert.False(t, IsSupersetForInterface([]interface{}{1, 2}, []interface{}{1, 2, 3}))
}

func TestIsSubsetMapByKeyForInterface(t *testing.T) {
	assert.True(t, IsSubsetMapByKeyForInterface[string](map[interface{}]string{1: "a", 2: "b"}, map[interface{}]string{1: "a", 2: "b", 3: "c"}))
	assert.False(t, IsSubsetMapByKeyForInterface[string](map[interface{}]string{1: "a", 4: "b"}, map[interface{}]string{1: "a", 2: "b", 3: "c"}))
}

func TestIsSupersetMapByKeyForInterface(t *testing.T) {
	assert.True(t, IsSupersetMapByKeyForInterface[string](map[interface{}]string{1: "a", 2: "b", 3: "c"}, map[interface{}]string{1: "a", 2: "b"}))
	assert.False(t, IsSupersetMapByKeyForInterface[string](map[interface{}]string{1: "a", 2: "b"}, map[interface{}]string{1: "a", 2: "b", 3: "c"}))
}

func TestRange(t *testing.T) {
	assert.Equal(t, []int{}, Range(5, 2))
	assert.Equal(t, []int{}, Range(1, 5, 0))
	assert.Equal(t, []int{}, Range(1, 5, -1))
	assert.Equal(t, []int{1, 3, 5, 7, 9}, Range(1, 10, 2))
}

func TestDrop(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, Drop(0, 1, 2, 3))
	assert.Equal(t, []int{}, Drop(5, 1, 2))
	assert.Equal(t, []int{}, Drop(3, 1, 2))
	assert.Equal(t, []int{3}, Drop(2, 1, 2, 3))
}

func TestDropLast(t *testing.T) {
	assert.Equal(t, []int{}, DropLast(3, 1, 2))
	assert.Equal(t, []int{1, 2}, DropLast(0, 1, 2))
	assert.Equal(t, []int{1}, DropLast(2, 1, 2, 3))
}

func TestReverse(t *testing.T) {
	assert.Equal(t, []int{}, Reverse[int]())
}

func TestSliceToMap(t *testing.T) {
	result := SliceToMap(0, 1, 2, 3)
	assert.Equal(t, map[int]int{1: 0, 2: 0, 3: 0}, result)
}

func TestPipeInterface(t *testing.T) {
	fn1 := func(args ...interface{}) []interface{} {
		val := args[0].(int)
		return []interface{}{val + 1}
	}
	fn2 := func(args ...interface{}) []interface{} {
		val := args[0].(int)
		return []interface{}{val * 2}
	}

	result := PipeInterface(fn1, fn2)(1)
	assert.Equal(t, 4, result[0])
}

func TestPipe(t *testing.T) {
	fn1 := func(args ...int) []int {
		val := args[0]
		return []int{val + 1}
	}
	fn2 := func(args ...int) []int {
		val := args[0]
		return []int{val * 2}
	}

	result := Pipe(fn1, fn2)(1)
	assert.Equal(t, 4, result[0])
}

func TestEvery(t *testing.T) {
	result := Every(func(a int) bool { return a > 0 }, 1, 2, 3)
	assert.Equal(t, true, result)

	result = Every(func(a int) bool { return a > 0 }, 1, -1, 3)
	assert.Equal(t, false, result)
}

func TestSome(t *testing.T) {
	result := Some(func(a int) bool { return a > 0 }, -1, -2, 3)
	assert.Equal(t, true, result)

	result = Some(func(a int) bool { return a > 0 }, -1, -2, -3)
	assert.Equal(t, false, result)
}

func TestCurryDef(t *testing.T) {
	curry := CurryNewGenerics[int, int](func(c *CurryDef[int, int], args ...int) int {
		sum := 0
		for _, v := range args {
			sum += v
		}
		return sum
	})

	curry.Call(1).Call(2).Call(3)
	assert.Equal(t, 6, curry.Result())
	assert.Equal(t, false, curry.IsDone())

	curry.MarkDone()
	assert.Equal(t, true, curry.IsDone())
}

func TestMinMax(t *testing.T) {
	min, max := MinMax(3, 1, 4, 1, 5, 9, 2, 6)
	assert.Equal(t, 1, min)
	assert.Equal(t, 9, max)
}

func TestMax(t *testing.T) {
	result := Max(3, 1, 4, 1, 5, 9, 2, 6)
	assert.Equal(t, 9, result)
}

func TestMin(t *testing.T) {
	result := Min(3, 1, 4, 1, 5, 9, 2, 6)
	assert.Equal(t, 1, result)
}

func TestKeys(t *testing.T) {
	result := Keys(map[int]string{1: "a", 2: "b", 3: "c"})
	assert.ElementsMatch(t, []int{1, 2, 3}, result)
}

func TestValues(t *testing.T) {
	result := Values(map[int]string{1: "a", 2: "b", 3: "c"})
	assert.ElementsMatch(t, []string{"a", "b", "c"}, result)
}

func TestMerge(t *testing.T) {
	result := Merge(map[int]string{1: "a"}, map[int]string{2: "b"})
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "a", result[1])
	assert.Equal(t, "b", result[2])
}

func TestTake(t *testing.T) {
	result := Take(3, 1, 2, 3, 4, 5)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestTakeLast(t *testing.T) {
	result := TakeLast(3, 1, 2, 3, 4, 5)
	assert.Equal(t, []int{3, 4, 5}, result)
}

func TestTakeAndTakeLast(t *testing.T) {
	result := Take(3, 1, 2, 3, 4, 5)
	assert.Equal(t, []int{1, 2, 3}, result)

	result = TakeLast(3, 1, 2, 3, 4, 5)
	assert.Equal(t, []int{3, 4, 5}, result)
}

func TestSliceToMapForInterface(t *testing.T) {
	result := SliceToMapForInterface(true, 1, "key", 3)
	assert.Equal(t, map[interface{}]bool{1: true, "key": true, 3: true}, result)
}

func TestDuplicateSlice(t *testing.T) {
	result := DuplicateSlice[int]([]int{1, 2, 3})
	assert.Equal(t, []int{1, 2, 3}, result)
	assert.NotSame(t, []int{1, 2, 3}, result)

	result = DuplicateSlice[int]([]int{})
	assert.Empty(t, result)

	result = DuplicateSlice[int](nil)
	assert.Empty(t, result)
}

func TestDuplicateMap(t *testing.T) {
	result := DuplicateMap[int, string](map[int]string{1: "a", 2: "b"})
	assert.Equal(t, map[int]string{1: "a", 2: "b"}, result)
	assert.NotSame(t, map[int]string{1: "a", 2: "b"}, result)

	result = DuplicateMap[int, string](map[int]string{})
	assert.Empty(t, result)

	result = DuplicateMap[int, string](nil)
	assert.Empty(t, result)
}

func TestDuplicateMapForInterface(t *testing.T) {
	result := DuplicateMapForInterface[string](map[interface{}]string{1: "a", "key": "b"})
	assert.Equal(t, map[interface{}]string{1: "a", "key": "b"}, result)

	result = DuplicateMapForInterface[string](nil)
	assert.Empty(t, result)
}

func TestIsNil(t *testing.T) {
	var nilPtr *int = nil
	assert.True(t, IsNil(nilPtr))
	assert.False(t, IsNil(1))
	assert.False(t, IsNil("test"))
}

func TestIsPtr(t *testing.T) {
	var ptr *int = new(int)
	assert.True(t, IsPtr(ptr))
	assert.False(t, IsPtr(1))
	assert.False(t, IsPtr("test"))
}

func TestPtrOf(t *testing.T) {
	val := 42
	ptr := PtrOf(val)
	assert.Equal(t, &val, ptr)
	assert.Equal(t, 42, *ptr)
}

func TestHead(t *testing.T) {
	var zeroInt int
	assert.Equal(t, zeroInt, Head[int]())
	assert.Equal(t, 1, Head(1, 2, 3))
}

func TestTail(t *testing.T) {
	assert.Equal(t, []int{}, Tail[int]())
	assert.Equal(t, []int{2, 3}, Tail(1, 2, 3))
}

func TestSortDescending(t *testing.T) {
	result := SortOrderedDescending(3, 1, 2)
	assert.Equal(t, []int{3, 2, 1}, result)
}

func TestCompareToOrdered(t *testing.T) {
	assert.Equal(t, 0, CompareToOrdered(1, 1))
	assert.Equal(t, 1, CompareToOrdered(1, 2))
	assert.Equal(t, -1, CompareToOrdered(2, 1))
}

func TestDropEq(t *testing.T) {
	result := DropEq(2, 1, 2, 3, 2, 4)
	assert.Equal(t, []int{1, 3, 4}, result)

	result = DropEq(5, 1, 2, 3)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestDropWhile(t *testing.T) {
	result := DropWhile(func(x int) bool { return x < 3 }, 1, 2, 3, 4, 5)
	assert.Equal(t, []int{3, 4, 5}, result)

	result = DropWhile(func(x int) bool { return x < 10 }, 1, 2, 3)
	assert.Nil(t, result)

	result = DropWhile(func(x int) bool { return x < 0 }, 1, 2, 3)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestGroupBy(t *testing.T) {
	result := GroupBy(func(x int) string { return "mod2_" + string(rune('0'+x%2)) }, 1, 2, 3, 4, 5)
	assert.Equal(t, []int{1, 3, 5}, result["mod2_1"])
	assert.Equal(t, []int{2, 4}, result["mod2_0"])

	result2 := GroupBy(func(x int) int { return x % 3 }, 1, 2, 3, 4, 5, 6)
	assert.Equal(t, []int{1, 4}, result2[1])
	assert.Equal(t, []int{2, 5}, result2[2])
	assert.Equal(t, []int{3, 6}, result2[0])
}

func TestPrepend(t *testing.T) {
	result := Prepend(0, []int{1, 2, 3})
	assert.Equal(t, []int{0, 1, 2, 3}, result)

	result = Prepend(1, []int{})
	assert.Equal(t, []int{1}, result)
}

func TestUniqBy(t *testing.T) {
	result := UniqBy(func(x int) int { return x % 2 }, 1, 2, 3, 4, 5, 6)
	assert.Equal(t, []int{1, 2}, result)
}

func TestDedupe(t *testing.T) {
	result := Dedupe(1, 2, 2, 3, 3, 3, 4)
	assert.Equal(t, []int{1, 2, 3, 4}, result)

	result = Dedupe(1, 1, 1)
	assert.Equal(t, []int{1}, result)
}

func TestDifference(t *testing.T) {
	result := Difference([]int{1, 2, 3, 4}, []int{2, 4})
	assert.Equal(t, []int{1, 3}, result)

	result = Difference([]int{1, 2, 3}, []int{4, 5, 6})
	assert.Equal(t, []int{1, 2, 3}, result)

	result = Difference([]int{1, 2}, []int{1, 2, 3, 4})
	assert.Nil(t, result)
}

func TestDistinct(t *testing.T) {
	result := Distinct(1, 2, 2, 3, 3, 3, 4, 4)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
}

func TestDistinctString(t *testing.T) {
	result := Distinct("a", "b", "a", "c")
	assert.Equal(t, []string{"a", "b", "c"}, result)
}

func TestDistinctRandom(t *testing.T) {
	result := DistinctRandom(1, 2, 2, 3, 3, 3, 4, 4)
	assert.Equal(t, 4, len(result))
	seen := make(map[int]bool)
	for _, v := range result {
		seen[v] = true
	}
	assert.True(t, seen[1] && seen[2] && seen[3] && seen[4])
}

func TestIntersection(t *testing.T) {
	result := Intersection([]int{1, 2, 3, 4}, []int{2, 4, 5})
	assert.Equal(t, []int{2, 4}, result)

	result = Intersection([]int{1, 2}, []int{3, 4})
	assert.Nil(t, result)
}

func TestUnion(t *testing.T) {
	result := Union([]int{1, 2, 3}, []int{2, 4, 5})
	assert.Equal(t, 5, len(result))
	seen := make(map[int]bool)
	for _, v := range result {
		seen[v] = true
	}
	assert.True(t, seen[1] && seen[2] && seen[3] && seen[4] && seen[5])
}

func TestIsSubset(t *testing.T) {
	assert.True(t, IsSubset([]int{1, 2}, []int{1, 2, 3, 4}))
	assert.False(t, IsSubset([]int{1, 5}, []int{1, 2, 3, 4}))
}

func TestIsSuperset(t *testing.T) {
	assert.True(t, IsSuperset([]int{1, 2, 3, 4}, []int{1, 2}))
	assert.False(t, IsSuperset([]int{1, 2, 3}, []int{1, 5}))
}

func TestIsEqual(t *testing.T) {
	assert.True(t, IsEqual([]int{1, 2, 3}, []int{1, 2, 3}))
	assert.False(t, IsEqual([]int{1, 2, 3}, []int{1, 2, 4}))
	assert.False(t, IsEqual([]int{1, 2}, []int{1, 2, 3}))
}

func TestDuplicateMapStr(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	duplicate := DuplicateMap(m)
	assert.Equal(t, m, duplicate)
	duplicate["a"] = 100
	assert.Equal(t, 1, m["a"])
	assert.Equal(t, 100, duplicate["a"])
}

func TestSliceOf(t *testing.T) {
	result := SliceOf[int](1, 2, 3)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestKeysStr(t *testing.T) {
	result := Keys(map[string]int{"a": 1, "b": 2})
	assert.Equal(t, 2, len(result))
}

func TestMinus(t *testing.T) {
	result := Minus([]int{1, 2, 3, 4}, []int{2, 4})
	assert.Equal(t, []int{1, 3}, result)

	result = Minus([]int{1, 2}, []int{3, 4})
	assert.Equal(t, []int{1, 2}, result)
}

func TestIsEqualMap(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"a": 1, "b": 2}
	m3 := map[string]int{"a": 1, "b": 3}
	m4 := map[string]int{"a": 1}

	assert.True(t, IsEqualMap(m1, m2))
	assert.False(t, IsEqualMap(m1, m3))
	assert.False(t, IsEqualMap(m1, m4))
}

func TestZipSuccess(t *testing.T) {
	result, err := Zip([]int{1, 2, 3}, []string{"a", "b", "c"})
	assert.Nil(t, err)
	assert.Equal(t, map[int]string{1: "a", 2: "b", 3: "c"}, result)
}

func TestMakeVariadicReturn1(t *testing.T) {
	fn := func(args ...int) int {
		sum := 0
		for _, v := range args {
			sum += v
		}
		return sum
	}
	variadicFn := MakeVariadicReturn1(fn)
	result := variadicFn(1, 2, 3)
	assert.Equal(t, []int{6}, result)
}

func TestMakeVariadicParam1(t *testing.T) {
	fn := func(a int) []int {
		return []int{a, a * 2, a * 3}
	}
	variadicFn := MakeVariadicParam1(fn)
	result := variadicFn(5)
	assert.Equal(t, []int{5, 10, 15}, result)
}

func TestMakeNumericReturnForParam1ReturnBool1(t *testing.T) {
	fn := func(x int) bool {
		return x > 5
	}
	variadicFn := MakeNumericReturnForParam1ReturnBool1[int, int](fn)
	result := variadicFn(6)
	assert.Equal(t, []int{1}, result)

	result = variadicFn(3)
	assert.Equal(t, []int{0}, result)
}

func TestComposeInterface(t *testing.T) {
	fn1 := func(args ...interface{}) []interface{} {
		val := args[0].(int)
		return []interface{}{val + 1}
	}
	fn2 := func(args ...interface{}) []interface{} {
		val := args[0].(int)
		return []interface{}{val * 2}
	}

	result := ComposeInterface(fn1, fn2)(1)
	assert.Equal(t, 3, result[0])
}

func TestConcatWithNil(t *testing.T) {
	result := Concat([]int{1, 2}, []int{3, 4}, nil)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
}

func TestIsDistinct(t *testing.T) {
	assert.False(t, IsDistinct[int]())
	assert.True(t, IsDistinct(1, 2, 3))
	assert.False(t, IsDistinct(1, 2, 1))
	assert.True(t, IsDistinct("a", "b", "c"))
	assert.False(t, IsDistinct("a", "b", "a"))
}

func TestEveryEmpty(t *testing.T) {
	result := Every(func(a int) bool { return a > 0 })
	assert.False(t, result)
}

func TestSomeEmpty(t *testing.T) {
	result := Some(func(a int) bool { return a > 0 })
	assert.False(t, result)
}

func TestFilterEmptyInput(t *testing.T) {
	result := Filter(func(a int, i int) bool { return a > 2 })
	assert.Equal(t, []int{}, result)
}

func TestRejectEmptyInput(t *testing.T) {
	result := Reject(func(a int, i int) bool { return a > 2 })
	assert.Equal(t, []int{}, result)
}

func TestReduceEmptyInput(t *testing.T) {
	result := Reduce(func(a, b int) int { return a + b }, 0)
	assert.Equal(t, 0, result)
}

func TestMapEmptyInput(t *testing.T) {
	result := Map(func(a int) int { return a * 2 })
	assert.Equal(t, []int{}, result)
}

func TestHeadEmpty(t *testing.T) {
	var zeroInt int
	assert.Equal(t, zeroInt, Head[int]())
}

func TestTailSingleElement(t *testing.T) {
	result := Tail(1)
	assert.Equal(t, []int{}, result)
}

func TestTakeMoreThanLength(t *testing.T) {
	result := Take(10, 1, 2, 3)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestTakeLastMoreThanLength(t *testing.T) {
	result := TakeLast(10, 1, 2, 3)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestDropMoreThanLength(t *testing.T) {
	result := Drop(10, 1, 2, 3)
	assert.Equal(t, []int{}, result)
}

func TestDropLastMoreThanLength(t *testing.T) {
	result := DropLast(10, 1, 2, 3)
	assert.Equal(t, []int{}, result)
}

func TestReverseSingle(t *testing.T) {
	result := Reverse(1)
	assert.Equal(t, []int{1}, result)
}

func TestRangeWithHops(t *testing.T) {
	result := Range(0, 10, 3)
	assert.Equal(t, []int{0, 3, 6, 9}, result)
}

func TestRangeSingleElement(t *testing.T) {
	result := Range(5, 6, 2)
	assert.Equal(t, []int{5}, result)
}

func TestMinMaxSingle(t *testing.T) {
	min, max := MinMax(5)
	assert.Equal(t, 5, min)
	assert.Equal(t, 5, max)
}

func TestMaxSingle(t *testing.T) {
	result := Max(5)
	assert.Equal(t, 5, result)
}

func TestMinSingle(t *testing.T) {
	result := Min(5)
	assert.Equal(t, 5, result)
}

func TestSortOrderedAscendingEmpty(t *testing.T) {
	result := SortOrderedAscending[int]()
	assert.Nil(t, result)
}

func TestSortOrderedDescendingEmpty(t *testing.T) {
	result := SortOrderedDescending[int]()
	assert.Nil(t, result)
}

func TestDifferenceNilList(t *testing.T) {
	result := Difference[int](nil)
	assert.Equal(t, []int{}, result)
}

func TestDifferenceSingleList(t *testing.T) {
	result := Difference([]int{1, 2, 3})
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestIntersectionNilList(t *testing.T) {
	result := Intersection[int](nil)
	assert.Nil(t, result)
}

func TestIntersectionSingleList(t *testing.T) {
	result := Intersection([]int{1, 2, 2, 3})
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestIntersectionForInterfaceNilList(t *testing.T) {
	result := IntersectionForInterface(nil)
	assert.Nil(t, result)
}

func TestIntersectionForInterfaceSingleList(t *testing.T) {
	result := IntersectionForInterface([]interface{}{1, 2, 2, 3})
	assert.Equal(t, []interface{}{1, 2, 3}, result)
}

func TestUnionNilLists(t *testing.T) {
	result := Union[int](nil)
	assert.Equal(t, []int{}, result)
}

func TestMinusNilLists(t *testing.T) {
	result := Minus[int](nil, nil)
	assert.Equal(t, []int{}, result)
}

func TestIsSubsetNil(t *testing.T) {
	assert.False(t, IsSubset[int](nil, nil))
	assert.False(t, IsSubset[int](nil, []int{1, 2}))
}

func TestIsSupersetNil(t *testing.T) {
	assert.False(t, IsSuperset[int](nil, nil))
	assert.False(t, IsSuperset[int]([]int{1, 2}, nil))
}

func TestIsEqualNil(t *testing.T) {
	assert.False(t, IsEqual[int](nil, nil))
	assert.False(t, IsEqual[int](nil, []int{1}))
	assert.False(t, IsEqual[int]([]int{1}, nil))
}

func TestSortSlice(t *testing.T) {
	input := []int{3, 1, 4, 1, 5}
	Sort(func(a, b int) bool { return a < b }, input)
	assert.Equal(t, []int{1, 1, 3, 4, 5}, input)
}

func TestSortSliceEmpty(t *testing.T) {
	input := []int{}
	Sort(func(a, b int) bool { return a < b }, input)
	assert.Equal(t, []int{}, input)
}

func TestMaxEmptyList(t *testing.T) {
	result := Max[int]()
	assert.Equal(t, 0, result)
}

func TestMinEmptyList(t *testing.T) {
	result := Min[int]()
	assert.Equal(t, 0, result)
}

func TestKeysEmptyMap(t *testing.T) {
	result := Keys(map[string]int{})
	assert.Equal(t, []string{}, result)
}

func TestValuesEmptyMap(t *testing.T) {
	result := Values(map[string]int{})
	assert.Equal(t, []int{}, result)
}

func TestMergeEmptyMaps(t *testing.T) {
	result := Merge(map[string]int{}, map[string]int{})
	assert.Equal(t, map[string]int{}, result)
}

func TestMergeWithOverwrite(t *testing.T) {
	result := Merge(map[string]int{"a": 1}, map[string]int{"a": 2, "b": 3})
	assert.Equal(t, 2, result["a"])
	assert.Equal(t, 3, result["b"])
}

func TestTakeEmpty(t *testing.T) {
	result := Take[int](2)
	assert.Nil(t, result)
}

func TestTakeLastEmpty(t *testing.T) {
	result := TakeLast[int](2)
	assert.Nil(t, result)
}

func TestDropEmpty(t *testing.T) {
	result := Drop[int](2)
	assert.Equal(t, []int{}, result)
}

func TestDropLastEmpty(t *testing.T) {
	result := DropLast[int](2)
	assert.Equal(t, []int{}, result)
}

func TestReverseEmpty(t *testing.T) {
	result := Reverse[int]()
	assert.Equal(t, []int{}, result)
}

func TestKind(t *testing.T) {
	assert.Equal(t, reflect.Int, Kind(42))
	assert.Equal(t, reflect.String, Kind("hello"))
	assert.Equal(t, reflect.Slice, Kind([]int{}))
	assert.Equal(t, reflect.Map, Kind(map[string]int{}))
	assert.Equal(t, reflect.Ptr, Kind(new(int)))
	assert.Equal(t, reflect.Struct, Kind(struct{}{}))
}

func TestIntersectionMapByKey(t *testing.T) {
	result := IntersectionMapByKey[string, int]()
	assert.Equal(t, map[string]int{}, result)

	result = IntersectionMapByKey(map[string]int{"a": 1, "b": 2})
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, result)

	result = IntersectionMapByKey(map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1, "c": 3})
	assert.Equal(t, map[string]int{"a": 1}, result)

	result = IntersectionMapByKey(map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1, "b": 2})
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, result)
}

func TestMinusMapByKey(t *testing.T) {
	result := MinusMapByKey(map[string]int{"a": 1, "b": 2}, map[string]int{"b": 2})
	assert.Equal(t, map[string]int{"a": 1}, result)

	result = MinusMapByKey(map[string]int{"a": 1, "b": 2}, map[string]int{"c": 3})
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, result)

	result = MinusMapByKey(map[string]int{"a": 1}, map[string]int{"a": 1, "b": 2})
	assert.Equal(t, map[string]int{}, result)
}

func TestIsSubsetMapByKey(t *testing.T) {
	assert.True(t, IsSubsetMapByKey(map[string]int{"a": 1}, map[string]int{"a": 1, "b": 2}))
	assert.False(t, IsSubsetMapByKey(map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1}))
	assert.False(t, IsSubsetMapByKey(map[string]int{"c": 3}, map[string]int{"a": 1}))
	assert.False(t, IsSubsetMapByKey(map[string]int{}, map[string]int{"a": 1}))
	assert.False(t, IsSubsetMapByKey(map[string]int{"a": 1}, map[string]int{}))
}

func TestIsSupersetMapByKey(t *testing.T) {
	assert.True(t, IsSupersetMapByKey(map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1}))
	assert.False(t, IsSupersetMapByKey(map[string]int{"a": 1}, map[string]int{"a": 1, "b": 2}))
	assert.False(t, IsSupersetMapByKey(map[string]int{"a": 1}, map[string]int{"c": 3}))
	assert.False(t, IsSupersetMapByKey(map[string]int{}, map[string]int{"a": 1}))
	assert.False(t, IsSupersetMapByKey(map[string]int{"a": 1}, map[string]int{}))
}

func TestNewCompData(t *testing.T) {
	compType := DefProduct(reflect.String, reflect.Int)
	compData := NewCompData(compType, "hello", 42)
	assert.NotNil(t, compData)
	assert.True(t, MatchCompTypeRef(compType, compData))

	compDataNil := NewCompData(compType)
	assert.Nil(t, compDataNil)

	compDataInvalid := NewCompData(compType, 42, "hello")
	assert.Nil(t, compDataInvalid)
}

func TestEither(t *testing.T) {
	result := Either(42, InCaseOfKind(reflect.Int, func(e interface{}) interface{} { return e }))
	assert.Equal(t, 42, result)

	result = Either("hello", InCaseOfKind(reflect.String, func(e interface{}) interface{} { return e }))
	assert.Equal(t, "hello", result)

	result = Either(42, InCaseOfKind(reflect.String, func(e interface{}) interface{} { return e }), Otherwise(func(e interface{}) interface{} { return e }))
	assert.Equal(t, 42, result)
}

func TestInCaseOfEqual(t *testing.T) {
	patterns := []Pattern{
		InCaseOfEqual("hello", func(e interface{}) interface{} { return "world" }),
		Otherwise(func(e interface{}) interface{} { return "fallback" }),
	}
	pm := DefPattern(patterns...)
	result := pm.MatchFor("hello")
	assert.Equal(t, "world", result)

	result = pm.MatchFor("other")
	assert.Equal(t, "fallback", result)
}

func TestInCaseOfRegex(t *testing.T) {
	patterns := []Pattern{
		InCaseOfRegex("^hello.*", func(e interface{}) interface{} { return "matched" }),
		Otherwise(func(e interface{}) interface{} { return "fallback" }),
	}
	pm := DefPattern(patterns...)
	result := pm.MatchFor("hello world")
	assert.Equal(t, "matched", result)

	result = pm.MatchFor("bye world")
	assert.Equal(t, "fallback", result)
}

func TestMatchCompType(t *testing.T) {
	compType := DefProduct(reflect.String, reflect.Int)
	compData := NewCompData(compType, "hello", 42)
	assert.NotNil(t, compData)
	assert.True(t, MatchCompType(compType, *compData))

	compDataInvalid := NewCompData(compType, 42, "hello")
	assert.Nil(t, compDataInvalid)
}

func TestSortOrdered(t *testing.T) {
	result := SortOrdered(true, 3, 1, 2)
	assert.Equal(t, []int{1, 2, 3}, result)

	result = SortOrdered(false, 3, 1, 2)
	assert.Equal(t, []int{3, 2, 1}, result)
}

func TestPMap(t *testing.T) {
	option := &PMapOption{
		FixedPool:   2,
		RandomOrder: false,
	}
	result := PMap(func(i int) int { return i * 2 }, option, 1, 2, 3, 4, 5)
	assert.Equal(t, []int{2, 4, 6, 8, 10}, result)
}

func TestPMapNoOrder(t *testing.T) {
	option := &PMapOption{
		FixedPool:   2,
		RandomOrder: true,
	}
	result := PMap(func(i int) int { return i * 2 }, option, 1, 2, 3, 4, 5)
	sorted := SortOrderedAscending(result...)
	assert.Equal(t, []int{2, 4, 6, 8, 10}, sorted)
}

func TestCurryParam1ForSlice1(t *testing.T) {
	fn := func(a int, b []int) int {
		return a + b[0]
	}
	curried := CurryParam1ForSlice1(fn, 10)
	result := curried(1, 2, 3)
	assert.Equal(t, 11, result)
}

func TestMakeNumericReturnForSliceParamReturnBool1(t *testing.T) {
	fn := func(v []int) bool {
		return len(v) > 0
	}
	curried := MakeNumericReturnForSliceParamReturnBool1[int, int](fn)
	result := curried(1, 2, 3)
	assert.Equal(t, []int{1}, result)
}

func TestMakeNumericReturnForVariadicParamReturnBool1(t *testing.T) {
	fn := func(v ...int) bool {
		return len(v) > 1
	}
	curried := MakeNumericReturnForVariadicParamReturnBool1[int, int](fn)
	result := curried(1, 2, 3)
	assert.Equal(t, []int{1}, result)
}

func TestInCaseOfSumType(t *testing.T) {
	myType := DefSum(
		DefProduct(reflect.String),
		DefProduct(reflect.Int),
	)
	patterns := []Pattern{
		InCaseOfSumType(myType, func(x interface{}) interface{} {
			cd := x.(CompData)
			return fmt.Sprintf("SumType: %v", cd.objects)
		}),
		Otherwise(func(x interface{}) interface{} {
			return "fallback"
		}),
	}
	pm := DefPattern(patterns...)
	result := pm.MatchFor(NewCompData(myType, "hello"))
	assert.Equal(t, "SumType: [hello]", result)
}

func TestDefSum(t *testing.T) {
	sumType := DefSum(
		DefProduct(reflect.String),
		DefProduct(reflect.Int),
	)
	assert.NotNil(t, sumType)

	data := NewCompData(sumType, "hello")
	assert.NotNil(t, data)
	assert.True(t, MatchCompTypeRef(sumType, data))

	dataInt := NewCompData(sumType, 42)
	assert.NotNil(t, dataInt)
	assert.True(t, MatchCompTypeRef(sumType, dataInt))
}

func TestMatchCompTypeRef(t *testing.T) {
	compType := DefProduct(reflect.String, reflect.Int)
	compData := NewCompData(compType, "hello", 42)
	assert.NotNil(t, compData)
	assert.True(t, MatchCompTypeRef(compType, compData))
}

func TestCurryNewGenerics(t *testing.T) {
	curried := CurryNewGenerics(func(c *CurryDef[int, int], args ...int) int {
		return 10 + args[0]
	})
	curried.Call(5)
	result := curried.Result()
	assert.Equal(t, 15, result)
}

func TestMakeVariadicParam2(t *testing.T) {
	fn := func(a, b int) []int {
		return []int{a + b}
	}
	curried := MakeVariadicParam2(fn)
	result := curried(1, 2, 3)
	assert.Equal(t, []int{3}, result)
}

func TestMakeVariadicParam3(t *testing.T) {
	fn := func(a, b, c int) []int {
		return []int{a + b + c}
	}
	curried := MakeVariadicParam3(fn)
	result := curried(1, 2, 3)
	assert.Equal(t, []int{6}, result)
}

func TestMakeVariadicParam4(t *testing.T) {
	fn := func(a, b, c, d int) []int {
		return []int{a + b + c + d}
	}
	curried := MakeVariadicParam4(fn)
	result := curried(1, 2, 3, 4)
	assert.Equal(t, []int{10}, result)
}

func TestMakeVariadicParam5(t *testing.T) {
	fn := func(a, b, c, d, e int) []int {
		return []int{a + b + c + d + e}
	}
	curried := MakeVariadicParam5(fn)
	result := curried(1, 2, 3, 4, 5)
	assert.Equal(t, []int{15}, result)
}

func TestMakeVariadicParam6(t *testing.T) {
	fn := func(a, b, c, d, e, f int) []int {
		return []int{a + b + c + d + e + f}
	}
	curried := MakeVariadicParam6(fn)
	result := curried(1, 2, 3, 4, 5, 6)
	assert.Equal(t, []int{21}, result)
}

func TestMakeVariadicReturn2(t *testing.T) {
	fn := func(...int) (int, int) {
		return 1, 2
	}
	curried := MakeVariadicReturn2(fn)
	result := curried(1, 2, 3)
	assert.Equal(t, []int{1, 2}, result)
}

func TestMakeVariadicReturn3(t *testing.T) {
	fn := func(...int) (int, int, int) {
		return 1, 2, 3
	}
	curried := MakeVariadicReturn3(fn)
	result := curried(1, 2, 3)
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestMakeVariadicReturn4(t *testing.T) {
	fn := func(...int) (int, int, int, int) {
		return 1, 2, 3, 4
	}
	curried := MakeVariadicReturn4(fn)
	result := curried(1, 2, 3)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
}

func TestMakeVariadicReturn5(t *testing.T) {
	fn := func(...int) (int, int, int, int, int) {
		return 1, 2, 3, 4, 5
	}
	curried := MakeVariadicReturn5(fn)
	result := curried(1, 2, 3)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
}

func TestMakeVariadicReturn6(t *testing.T) {
	fn := func(...int) (int, int, int, int, int, int) {
		return 1, 2, 3, 4, 5, 6
	}
	curried := MakeVariadicReturn6(fn)
	result := curried(1, 2, 3)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
}

func TestMakeNumericReturnForParam1ReturnBool1Alt(t *testing.T) {
	fn := func(v int) bool {
		return v > 0
	}
	curried := MakeNumericReturnForParam1ReturnBool1[int, int](fn)
	result := curried(1)
	assert.Equal(t, []int{1}, result)

	result = curried(-1)
	assert.Equal(t, []int{0}, result)
}

func TestCurryParam2(t *testing.T) {
	fn := func(a, b int, extra ...int) int {
		return a + b + extra[0]
	}
	curried := CurryParam2(fn, 1, 2)
	result := curried(3)
	assert.Equal(t, 6, result)
}

func TestCurryParam3(t *testing.T) {
	fn := func(a, b, c int, extra ...int) int {
		return a + b + c + extra[0]
	}
	curried := CurryParam3(fn, 1, 2, 3)
	result := curried(4)
	assert.Equal(t, 10, result)
}

func TestCurryParam4(t *testing.T) {
	fn := func(a, b, c, d int, extra ...int) int {
		return a + b + c + d + extra[0]
	}
	curried := CurryParam4(fn, 1, 2, 3, 4)
	result := curried(5)
	assert.Equal(t, 15, result)
}

func TestCurryParam5(t *testing.T) {
	fn := func(a, b, c, d, e int, extra ...int) int {
		return a + b + c + d + e + extra[0]
	}
	curried := CurryParam5(fn, 1, 2, 3, 4, 5)
	result := curried(6)
	assert.Equal(t, 21, result)
}

func TestCurryParam6(t *testing.T) {
	fn := func(a, b, c, d, e, f int, extra ...int) int {
		return a + b + c + d + e + f + extra[0]
	}
	curried := CurryParam6(fn, 1, 2, 3, 4, 5, 6)
	result := curried(7)
	assert.Equal(t, 28, result)
}

func TestCurryNew(t *testing.T) {
	curry := CurryNew(func(c *CurryDef[interface{}, interface{}], args ...interface{}) interface{} {
		return nil
	})
	assert.NotNil(t, curry)
}

func TestIsEqualMapAlt(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"a": 1, "b": 2}
	m3 := map[string]int{"a": 1, "b": 3}

	assert.True(t, IsEqualMap(m1, m2))
	assert.False(t, IsEqualMap(m1, m3))
	assert.False(t, IsEqualMap(m1, nil))
	assert.False(t, IsEqualMap(nil, m2))
}

func TestPMapEmptyInput(t *testing.T) {
	result := PMap(func(a int) int { return a + 1 }, nil)
	assert.Equal(t, []int{}, result)
}

func TestPMapWithZeroWorkers(t *testing.T) {
	result := PMap(func(a int) int { return a + 1 }, &PMapOption{FixedPool: 0}, 1, 2, 3)
	assert.Equal(t, []int{2, 3, 4}, result)
}

func TestFilterNilInput(t *testing.T) {
	result := Filter(func(a int, i int) bool { return a > 2 })
	assert.Equal(t, []int{}, result)
}

func TestMapWithNilInput(t *testing.T) {
	result := Map(func(a int) int { return a * 2 })
	assert.Equal(t, []int{}, result)
}

func TestReduceWithNilInput(t *testing.T) {
	result := Reduce(func(a, b int) int { return a + b }, 0)
	assert.Equal(t, 0, result)
}

func TestFlattenEmpty(t *testing.T) {
	result := Flatten([]int{}, []int{}, []int{})
	assert.Equal(t, []int{}, result)
}

func TestFlattenNil(t *testing.T) {
	result := Flatten[int]()
	assert.Equal(t, []int{}, result)
}

func TestZipWithEmptyLists(t *testing.T) {
	_, err := Zip([]int{}, []string{})
	assert.Equal(t, ErrZipEmptyList, err)
}

func TestZipWithMismatchedLengths(t *testing.T) {
	_, err := Zip([]int{1, 2}, []string{"a"})
	assert.Equal(t, ErrZipLengthMismatch, err)
}

func TestPartitionEmptyInput(t *testing.T) {
	result := Partition(func(a int) bool { return a%2 == 0 })
	assert.Equal(t, [][]int{{}, {}}, result)
}

func TestSplitEveryEmptyInput(t *testing.T) {
	result := SplitEvery(3, []int{}...)
	assert.Equal(t, [][]int{{}}, result)
}

func TestGroupByEmptyInput(t *testing.T) {
	result := GroupBy(func(a int) int { return a % 2 }, []int{}...)
	assert.Equal(t, map[int][]int{}, result)
}

func TestUniqByEmptyInput(t *testing.T) {
	result := UniqBy(func(a int) int { return a % 2 }, []int{}...)
	assert.Equal(t, []int{}, result)
}

func TestDedupeEmptyInput(t *testing.T) {
	result := Dedupe([]int{}...)
	assert.Nil(t, result)
}

func TestDifferenceEmptyInput(t *testing.T) {
	result := Difference([]int{})
	assert.Equal(t, []int{}, result)
}

func TestDistinctEmptyInput(t *testing.T) {
	result := Distinct([]int{}...)
	assert.Equal(t, []int{}, result)
}

func TestIntersectionEmptyInput(t *testing.T) {
	result := Intersection([]int{})
	assert.Nil(t, result)
}

func TestUnionEmptyInput(t *testing.T) {
	result := Union([]int{})
	assert.Equal(t, []int{}, result)
}

func TestMinusEmptyInput(t *testing.T) {
	result := Minus([]int{}, []int{})
	assert.Equal(t, []int{}, result)
}

func TestIsSubsetEmptyInput(t *testing.T) {
	assert.False(t, IsSubset([]int{}, []int{}))
	assert.False(t, IsSubset([]int{1}, []int{}))
}

func TestIsSupersetEmptyInput(t *testing.T) {
	assert.False(t, IsSuperset([]int{}, []int{}))
	assert.False(t, IsSuperset([]int{}, []int{1}))
}

func TestIsEqualEmptyInput(t *testing.T) {
	assert.False(t, IsEqual([]int{}, []int{}))
	assert.False(t, IsEqual([]int{1}, []int{}))
	assert.False(t, IsEqual([]int{}, []int{1}))
}

// Tests for nil input checks - coverage improvement

func TestDifferenceNilInput(t *testing.T) {
	// Test nil input slice (line 271-273)
	var nilSlice []int
	result := Difference(nilSlice)
	assert.Equal(t, []int{}, result)
}

func TestDropWhileNilPredicate(t *testing.T) {
	// Test nil predicate (line 447-449)
	var fn func(int) bool
	result := DropWhile(fn, 1, 2, 3)
	assert.Equal(t, []int{}, result)
}

func TestIntersectionNilInput(t *testing.T) {
	// Test nil input (line 577-579)
	var nilInput [][]int
	result := Intersection(nilInput...)
	assert.Equal(t, []int{}, result)
}

func TestIntersectionForInterfaceNilInput(t *testing.T) {
	// Test nil input (line 626-628)
	var nilInput [][]interface{}
	result := IntersectionForInterface(nilInput...)
	assert.Equal(t, []interface{}{}, result)
}

func TestIntersectionMapByKeyForInterfaceEmptyInput(t *testing.T) {
	// Test empty input (line 710-712)
	result := IntersectionMapByKeyForInterface[int]()
	assert.Equal(t, map[interface{}]int{}, result)
}

func TestIntersectionMapByKeyForInterfaceSingleInput(t *testing.T) {
	// Test single input (line 714-719)
	input := map[interface{}]int{"a": 1, "b": 2}
	result := IntersectionMapByKeyForInterface(input)
	assert.Equal(t, input, result)
}

func TestMinMaxNilInput(t *testing.T) {
	// Test nil input (line 868-870)
	var nilList []int
	min, max := MinMax(nilList...)
	assert.Equal(t, 0, min)
	assert.Equal(t, 0, max)
}

func TestMinMaxEmptyInput(t *testing.T) {
	// Test empty input (line 868-870)
	min, max := MinMax([]int{}...)
	assert.Equal(t, 0, min)
	assert.Equal(t, 0, max)
}

func TestMergeBothNilMaps(t *testing.T) {
	// Test both maps nil (line 886-888)
	var m1 map[string]int
	var m2 map[string]int
	result := Merge(m1, m2)
	assert.Equal(t, map[string]int{}, result)
}

func TestMergeFirstNilSecondValid(t *testing.T) {
	// Test first nil, second valid (line 892-897)
	var m1 map[string]int
	m2 := map[string]int{"b": 2}
	result := Merge(m1, m2)
	assert.Equal(t, map[string]int{"b": 2}, result)
}

func TestMergeFirstValidSecondNil(t *testing.T) {
	// Test first valid, second nil
	m1 := map[string]int{"a": 1}
	var m2 map[string]int
	result := Merge(m1, m2)
	assert.Equal(t, map[string]int{"a": 1}, result)
}

func TestPMapNilFunction(t *testing.T) {
	var nilFn func(int) int
	result := PMap(nilFn, nil, 1, 2, 3)
	assert.Equal(t, []int{}, result)
}

func TestSomeNilFunction(t *testing.T) {
	var nilFn func(int) bool
	result := Some(nilFn, 1, 2, 3)
	assert.Equal(t, false, result)
}

func TestTrampolineError(t *testing.T) {
	result, err := Trampoline(func(input ...int) ([]int, bool, error) {
		return nil, false, errors.New("test error")
	}, 1)
	assert.Nil(t, result)
	assert.EqualError(t, err, "test error")
}

func TestDifferenceFirstElementDupes(t *testing.T) {
	result := Difference([]int{1, 1, 2, 3}, []int{4, 5})
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestDifferenceAllMatch(t *testing.T) {
	result := Difference([]int{1, 2}, []int{1, 2})
	assert.Nil(t, result)
}

func TestMakeNumericReturnForVariadicParamReturnBool1False(t *testing.T) {
	fn := func(v ...int) bool {
		return false
	}
	result := MakeNumericReturnForVariadicParamReturnBool1[int, int](fn)(1, 2, 3)
	assert.Equal(t, []int{0}, result)
}

func TestMakeNumericReturnForSliceParamReturnBool1False(t *testing.T) {
	fn := func(v []int) bool {
		return false
	}
	result := MakeNumericReturnForSliceParamReturnBool1[int, int](fn)(1, 2, 3)
	assert.Equal(t, []int{0}, result)
}

func TestIsSubsetMapByKeyForInterfaceBothNil(t *testing.T) {
	result := IsSubsetMapByKeyForInterface[string](nil, nil)
	assert.Equal(t, false, result)
}

func TestIsSubsetMapByKeyForInterfaceFirstEmpty(t *testing.T) {
	result := IsSubsetMapByKeyForInterface[string](map[interface{}]string{}, map[interface{}]string{1: "a"})
	assert.Equal(t, false, result)
}

func TestIsSubsetMapByKeyForInterfaceSecondEmpty(t *testing.T) {
	result := IsSubsetMapByKeyForInterface[string](map[interface{}]string{1: "a"}, map[interface{}]string{})
	assert.Equal(t, false, result)
}

func TestPMapWithNilOption(t *testing.T) {
	result := PMap(func(a int) int { return a + 1 }, nil, 1, 2)
	assert.Equal(t, []int{2, 3}, result)
}

func TestMatchForWithPanic(t *testing.T) {
	pm := DefPattern()
	assert.Panics(t, func() {
		pm.MatchFor(42)
	})
}

func TestDropWhileNilFunc(t *testing.T) {
	result := DropWhile[int](nil, 1, 2, 3)
	assert.Equal(t, []int{}, result)
}

func TestKindPatternMatchesNil(t *testing.T) {
	pattern := InCaseOfKind(reflect.Int, func(x interface{}) interface{} { return x })
	result := pattern.Matches(nil)
	assert.False(t, result)
}

func TestRegexPatternMatchesNonString(t *testing.T) {
	pattern := InCaseOfRegex(".*", func(x interface{}) interface{} { return x })

	result := pattern.Matches(42)
	assert.False(t, result)

	result = pattern.Matches(nil)
	assert.False(t, result)
}

func TestRegexPatternMatchesInvalidRegex(t *testing.T) {
	pattern := InCaseOfRegex("[invalid", func(x interface{}) interface{} { return x })
	result := pattern.Matches("test")
	assert.False(t, result)
}

func TestMatchForWithCompTypePtr(t *testing.T) {
	compType := DefProduct(reflect.Int, reflect.String)
	data := NewCompData(compType, 1, "hello")
	patterns := []Pattern{
		InCaseOfSumType(compType, func(x interface{}) interface{} { return "matched" }),
	}
	result := DefPattern(patterns...).MatchFor(data)
	assert.Equal(t, "matched", result)
}

func TestMatchForWithNilPtrInput(t *testing.T) {
	var nilPtr *int = nil
	patterns := []Pattern{
		InCaseOfKind(reflect.Ptr, func(x interface{}) interface{} { return "ptr" }),
		Otherwise(func(x interface{}) interface{} { return "other" }),
	}
	result := DefPattern(patterns...).MatchFor(nilPtr)
	assert.Equal(t, "other", result)
}

func TestPMapNoOrderEmptyList(t *testing.T) {
	result := PMap(func(a int) int { return a + 1 }, &PMapOption{RandomOrder: true})
	assert.Equal(t, []int{}, result)
}

func TestPMapNoOrderWithEmptyListFixedPool(t *testing.T) {
	result := PMap(func(a int) int { return a + 1 }, &PMapOption{FixedPool: 2, RandomOrder: true})
	assert.Equal(t, []int{}, result)
}

func TestDifferenceMultipleArrays(t *testing.T) {
	result := Difference([]int{1, 2, 3}, []int{4, 5}, []int{6, 7})
	assert.Equal(t, []int{1, 2, 3}, result)
}

func TestDifferenceNoArgs(t *testing.T) {
	result := Difference[int]()
	assert.Equal(t, []int{}, result)
}

func TestPMapNoOrderWorkerGtListLen(t *testing.T) {
	result := PMap(func(a int) int { return a + 1 }, &PMapOption{FixedPool: 100, RandomOrder: true}, 1, 2, 3)
	sorted := SortOrderedAscending(result...)
	assert.Equal(t, []int{2, 3, 4}, sorted)
}

func TestPMapNoOrderEmptyListRandomOrder(t *testing.T) {
	result := PMap(func(a int) int { return a + 1 }, &PMapOption{RandomOrder: true})
	assert.Len(t, result, 0)
}

func TestPMapNoOrderCalledDirectly(t *testing.T) {
	result := pMapNoOrder(func(a int) int { return a + 1 }, []int{1, 2, 3}, 10)
	sorted := SortOrderedAscending(result...)
	assert.Equal(t, []int{2, 3, 4}, sorted)
}
