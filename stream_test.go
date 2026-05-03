package fpgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromArrayMapReduce(t *testing.T) {
	var s *StreamDef[string]
	var tempString string

	s = StreamFrom("1", "2", "3", "4")
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += v
	}
	assert.Equal(t, "1234", tempString)
	s = s.Map(func(item string, index int) string {
		val := item
		result := "v" + val
		return result
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "v1v2v3v4", tempString)
	tempString = ""

	s2 := StreamFromArray([]string{"1", "2", "3", "4"})
	for _, v := range s2.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)
	s2 = s2.Map(func(item string, index int) string {
		val := Maybe.Just(item).ToString()
		var result string = "v" + val
		return result
	})
	tempString = ""
	for _, v := range s2.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "v1v2v3v4", tempString)

	s3 := StreamFrom(1, 2, 3, 4)
	tempString = ""
	for _, v := range s3.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)
	s3 = s3.Map(func(item int, index int) int {
		return item * item
	})
	tempString = ""
	for _, v := range s3.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "14916", tempString)

	s3 = StreamFrom(1, 2, 3, 4)
	tempString = ""
	for _, v := range s3.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)
	s3 = s3.Map(func(item int, index int) int {
		return item * item
	})
	tempString = ""
	for _, v := range s3.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "14916", tempString)

	s3 = StreamFrom(1, 2, 3, 4)
	tempString = ""
	for _, v := range s3.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)

	s3 = s3.Map(func(item int, index int) int {
		return item * item
	})
	tempString = ""
	for _, v := range s3.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "14916", tempString)

	s4 := StreamFrom(true, false, true, false)
	tempString = ""
	for _, v := range s4.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "truefalsetruefalse", tempString)

	s3 = StreamFrom(1, 2, 3, 4)
	tempString = ""
	for _, v := range s3.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)

	s3 = StreamFrom(1, 2, 3, 4)
	tempString = ""
	for _, v := range s3.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)

	s3 = StreamFrom(1, 2, 3, 4)
	tempString = ""
	for _, v := range s3.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)

	s3 = StreamFrom(1, 2, 3, 4)
	tempString = ""
	for _, v := range s3.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)

	s3 = StreamFrom(1, 2, 3, 4)
	tempString = ""
	for _, v := range s3.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)
}

func TestFilter(t *testing.T) {
	var s *StreamDef[int]
	var tempString string

	s = StreamFromArray([]int{}).Append(1, 1).Extend(StreamFromArray([]int{2, 3, 4}))
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "11234", tempString)
	s = s.Distinct()
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)
	s = s.FilterNotNil()
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)

	s = s.Filter(func(item int, index int) bool {
		return item > 1 && item < 4
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "23", tempString)
	s = s.Reject(func(item int, index int) bool {
		return item > 2
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "2", tempString)
}

func TestSort(t *testing.T) {
	var s *StreamDef[int]
	var tempString string

	s = StreamFrom(11).Extend(StreamFrom(2, 3, 4, 5)).Remove(4)
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "11234", tempString)

	s = StreamFrom(11).Concat([]int{2, 3, 4, 5}).RemoveItem(4)
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "11235", tempString)
	tempString = ""
	for _, v := range s.Reverse().ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "53211", tempString)

	s = StreamFrom(11).Concat([]int{2, 3, 4, 5}).Remove(4)
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "11234", tempString)

	tempString = ""
	for _, v := range s.SortByIndex(func(i, j int) bool {
		vali, _ := Maybe.Just(s.Get(i)).ToInt()
		valj, _ := Maybe.Just(s.Get(j)).ToInt()
		return vali < valj
	}).ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "23411", tempString)
	tempString = ""
	for _, v := range s.Sort(func(a, b int) bool {
		return a < b
	}).ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "23411", tempString)
}

func TestStreamSetOperation(t *testing.T) {
	var s *StreamDef[int]
	var s2 *StreamDef[int]
	var s3 *StreamDef[int]
	var tempString string

	s = StreamFrom(11, 2, 3, 4, 5)
	s2 = StreamFrom(9, 2, 5, 6)
	s3 = StreamFrom(2, 5)
	assert.Equal(t, true, s.Contains(4))
	assert.Equal(t, false, s.Contains(6))
	assert.Equal(t, true, s.IsSuperset(s3))
	assert.Equal(t, true, s2.IsSuperset(s3))
	assert.Equal(t, true, s3.IsSubset(s))
	assert.Equal(t, true, s3.IsSubset(s2))
	assert.Equal(t, false, s.IsSuperset(s2))
	tempString = ""
	for _, v := range s.Clone().Intersection(s2).ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "2/5/", tempString)
	assert.Equal(t, 2, s.Intersection(s2).Len())
	tempString = ""
	for _, v := range s.Extend(s2).Distinct().ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "11/2/3/4/5/9/6/", tempString)
	assert.Equal(t, 7, s.Extend(s2).Distinct().Len())
	tempString = ""
	for _, v := range s.Minus(s2).ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "11/3/4/", tempString)
	tempString = ""
	for _, v := range s2.Minus(s).ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "9/6/", tempString)
}

func TestSetSetOperation(t *testing.T) {
	var s *MapSetDef[int, bool]
	var s2 *MapSetDef[int, bool]
	var s3 *MapSetDef[int, bool]
	var tempString string

	s = SetFrom[int, bool](11, 2, 3, 4, 5)
	s2 = SetFrom[int, bool](9, 2, 5, 6)
	s3 = SetFrom[int, bool](2, 5)
	assert.Equal(t, true, s.ContainsKey(4))
	assert.Equal(t, false, s.ContainsKey(6))
	assert.Equal(t, true, s.IsSupersetByKey(s3))
	assert.Equal(t, true, s2.IsSupersetByKey(s3))
	assert.Equal(t, true, s3.IsSubsetByKey(s))
	assert.Equal(t, true, s3.IsSubsetByKey(s2))
	assert.Equal(t, false, s.IsSupersetByKey(s2))
	tempString = ""
	for _, v := range SortOrderedAscending(s.Clone().Intersection(s2).Keys()...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "2/5/", tempString)
	assert.Equal(t, 2, s.Intersection(s2).Size())
	tempString = ""
	for _, v := range SortOrderedAscending(s.Union(s2).Keys()...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "2/3/4/5/6/9/11/", tempString)
	assert.Equal(t, 7, s.Union(s2).Size())
	tempString = ""
	for _, v := range SortOrderedAscending(s.Minus(s2).Keys()...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "3/4/11/", tempString)
	tempString = ""
	for _, v := range SortOrderedAscending(s2.Minus(s).Keys()...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "6/9/", tempString)
}

func TestStreamSetSetOperation(t *testing.T) {
	s1 := StreamFrom(1, 2, 3)

	cloned := s1.Clone()
	assert.Equal(t, s1.Len(), cloned.Len())

	contains := s1.Contains(2)
	assert.Equal(t, true, contains)

	contains = s1.Contains(5)
	assert.Equal(t, false, contains)
}

func TestStreamClone(t *testing.T) {
	s1 := StreamFrom(1, 2, 3)
	s2 := s1.Clone()

	s1 = s1.Map(func(i int, _ int) int { return i * 10 })

	assert.Equal(t, []int{1, 2, 3}, s2.ToArray())
	assert.Equal(t, []int{10, 20, 30}, s1.ToArray())
}

func TestStreamContains(t *testing.T) {
	s := StreamFrom(1, 2, 3)
	assert.Equal(t, true, s.Contains(1))
	assert.Equal(t, true, s.Contains(2))
	assert.Equal(t, true, s.Contains(3))
	assert.Equal(t, false, s.Contains(4))
}

func TestStreamLen(t *testing.T) {
	s := StreamFrom(1, 2, 3, 4, 5)
	assert.Equal(t, 5, s.Len())
}

func TestStreamGet(t *testing.T) {
	s := StreamFrom(10, 20, 30)
	assert.Equal(t, 10, s.Get(0))
	assert.Equal(t, 20, s.Get(1))
	assert.Equal(t, 30, s.Get(2))
}

func TestStreamRemove(t *testing.T) {
	s := StreamFrom(1, 2, 3, 4, 5)
	s = s.Remove(2)
	assert.Equal(t, []int{1, 2, 4, 5}, s.ToArray())
}

func TestStreamAppend(t *testing.T) {
	s := StreamFrom(1, 2, 3)
	s = s.Append(4, 5)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, s.ToArray())
}

func TestStreamConcat(t *testing.T) {
	s1 := StreamFrom(1, 2)
	s1 = s1.Concat([]int{3, 4})
	assert.Equal(t, []int{1, 2, 3, 4}, s1.ToArray())
}

func TestStreamExtend(t *testing.T) {
	s1 := StreamFrom(1, 2)
	_ = StreamFrom(3, 4)
	s1 = s1.Extend(StreamFrom(3, 4))
	assert.Equal(t, []int{1, 2, 3, 4}, s1.ToArray())
}

func TestStreamFilterNotNil(t *testing.T) {
	s := StreamFrom(1, 2, 3)
	s = s.FilterNotNil()
	assert.Equal(t, 3, s.Len())
}

func TestStreamReject(t *testing.T) {
	s := StreamFrom(1, 2, 3, 4, 5)
	s = s.Reject(func(val int, _ int) bool { return val > 3 })
	assert.Equal(t, []int{1, 2, 3}, s.ToArray())
}

func TestStreamDistinct(t *testing.T) {
	s := StreamFrom(1, 2, 2, 3, 3, 3)
	s = s.Distinct()
	assert.Equal(t, 3, s.Len())
}

func TestStreamReverse(t *testing.T) {
	s := StreamFrom(1, 2, 3)
	s = s.Reverse()
	assert.Equal(t, []int{3, 2, 1}, s.ToArray())
}

func TestMapSetOperations(t *testing.T) {
	s := SetFrom[int, bool](1, 2, 3)

	assert.Equal(t, true, s.ContainsKey(1))
	assert.Equal(t, false, s.ContainsKey(99))

	assert.Equal(t, 3, s.Size())

	cloned := s.Clone()
	assert.Equal(t, s.Size(), cloned.Size())

	keys := s.Keys()
	assert.Equal(t, 3, len(keys))
}

func streamIntTransformer(s *StreamDef[int]) string {
	result := ""
	for _, item := range SortOrderedAscending(s.ToArray()...) {
		result += Maybe.Just(item).ToMaybe().ToString() + ","
	}
	return result + "end"
}

func TestStreamIsSubset(t *testing.T) {
	s1 := StreamFrom(1, 2, 3)
	s2 := StreamFrom(1, 2, 3, 4, 5)
	assert.True(t, s1.IsSubset(s2))

	s1 = StreamFrom(1, 2, 6)
	s2 = StreamFrom(1, 2, 3, 4, 5)
	assert.False(t, s1.IsSubset(s2))
}

func TestStreamIsSuperset(t *testing.T) {
	s1 := StreamFrom(1, 2, 3, 4, 5)
	s2 := StreamFrom(1, 2, 3)
	assert.True(t, s1.IsSuperset(s2))

	s1 = StreamFrom(1, 2, 3)
	s2 = StreamFrom(1, 2, 3, 4, 5)
	assert.False(t, s1.IsSuperset(s2))
}

func TestMapSetContainsValue(t *testing.T) {
	s := SetFrom[int, string](1, 2, 3)
	s.Set(1, "one")
	s.Set(2, "two")

	assert.True(t, s.ContainsValue("one"))
	assert.True(t, s.ContainsValue("two"))
	assert.False(t, s.ContainsValue("three"))
}

func TestMapSetMapValue(t *testing.T) {
	s := SetFrom[int, int](1, 2, 3)
	s.Set(1, 10)
	s.Set(2, 20)
	s.Set(3, 30)

	s = s.MapValue(func(v int) int { return v * 2 }).(*MapSetDef[int, int])

	val := s.Get(1)
	assert.Equal(t, 20, val)
	val = s.Get(2)
	assert.Equal(t, 40, val)
	val = s.Get(3)
	assert.Equal(t, 60, val)
}

func TestMapSetMapKey(t *testing.T) {
	s := SetFrom[int, int](1, 2, 3)
	s = s.MapKey(func(x int) int { return x * 10 }).(*MapSetDef[int, int])

	assert.True(t, s.ContainsKey(10))
	assert.True(t, s.ContainsKey(20))
	assert.True(t, s.ContainsKey(30))
	assert.False(t, s.ContainsKey(1))
}

func TestMapSetRemoveKeys(t *testing.T) {
	s := SetFrom[int, string](1, 2, 3, 4, 5)
	s.Set(1, "one")
	s.Set(2, "two")
	s.Set(3, "three")

	s = s.RemoveKeys(2, 4).(*MapSetDef[int, string])

	assert.Equal(t, 3, s.Size())
	assert.False(t, s.ContainsKey(2))
	assert.False(t, s.ContainsKey(4))
	assert.True(t, s.ContainsKey(1))
	assert.True(t, s.ContainsKey(3))
}

func TestMapSetRemoveValues(t *testing.T) {
	s := SetFrom[int, string](1, 2, 3)
	s.Set(1, "a")
	s.Set(2, "b")
	s.Set(3, "c")

	s = s.RemoveValues("a", "c").(*MapSetDef[int, string])

	assert.Equal(t, 1, s.Size())
	assert.True(t, s.ContainsKey(2))
	assert.False(t, s.ContainsKey(1))
	assert.False(t, s.ContainsKey(3))
}

func TestStreamRemoveItem(t *testing.T) {
	s := StreamFrom(1, 2, 3, 2, 4)
	s = s.RemoveItem(2)
	assert.Equal(t, []int{1, 3, 4}, s.ToArray())
}

func TestStreamSort(t *testing.T) {
	s := StreamFrom(3, 1, 4, 1, 5, 9, 2, 6)
	s = s.Sort(func(a, b int) bool { return a < b })
	assert.Equal(t, []int{1, 1, 2, 3, 4, 5, 6, 9}, s.ToArray())
}

func TestStreamMinus(t *testing.T) {
	s1 := StreamFrom(1, 2, 3, 4)
	s2 := StreamFrom(3, 4, 5)
	result := s1.Minus(s2)
	assert.Equal(t, []int{1, 2}, result.ToArray())
}

func TestStreamSortByIndex(t *testing.T) {
	s := StreamFrom(3, 1, 2)
	_ = s.SortByIndex(func(a, b int) bool { return a > b })
}

func TestStreamRemoveByIndex(t *testing.T) {
	s := StreamFrom(1, 2, 3)
	s = s.Remove(1)
	assert.Equal(t, []int{1, 3}, s.ToArray())
}

func TestStreamMap(t *testing.T) {
	s := StreamFrom(1, 2, 3)
	s = s.Map(func(item int, index int) int {
		return item * 2
	})
	assert.Equal(t, []int{2, 4, 6}, s.ToArray())
}

func TestStreamConcatMulti(t *testing.T) {
	s := StreamFrom(1, 2)
	s = s.Concat([]int{3, 4}, []int{5, 6})
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, s.ToArray())
}

func TestStreamExtendMulti(t *testing.T) {
	s1 := StreamFrom(1, 2)
	s2 := StreamFrom(3, 4)
	s := s1.Extend(s2)
	assert.Equal(t, []int{1, 2, 3, 4}, s.ToArray())
}

func TestMapSetAsMap(t *testing.T) {
	m := SetFromArray[int, string]([]int{1, 2})
	m.Set(1, "one")
	m.Set(2, "two")

	result := m.AsMap()
	assert.Equal(t, map[int]string{1: "one", 2: "two"}, result)
}

func TestMapSetAsMapSet(t *testing.T) {
	m := SetFromArray[int, string]([]int{1, 2})
	m.Set(1, "one")
	m.Set(2, "two")

	result := m.AsMapSet()
	assert.NotNil(t, result)
}

func TestStreamSetClone(t *testing.T) {
	s := StreamFrom(1, 2, 3)

	cloned := s.Clone()
	assert.Equal(t, s.ToArray(), cloned.ToArray())
}

func TestStreamSetContainsKey(t *testing.T) {
	m := SetFromArray[int, string]([]int{1, 2})
	m.Set(1, "one")
	m.Set(2, "two")

	assert.True(t, m.ContainsKey(1))
	assert.False(t, m.ContainsKey(3))
}

func TestStreamSetKeys(t *testing.T) {
	m := SetFromArray[int, string]([]int{1, 2})
	m.Set(1, "one")
	m.Set(2, "two")

	keys := m.Keys()
	assert.ElementsMatch(t, []int{1, 2}, keys)
}

func TestStreamSetValues(t *testing.T) {
	m := SetFromArray[int, string]([]int{1, 2})
	m.Set(1, "one")
	m.Set(2, "two")

	values := m.Values()
	assert.ElementsMatch(t, []string{"one", "two"}, values)
}

func TestStreamSetDefMinusStreams(t *testing.T) {
	ss1 := StreamSetFrom[int, int](1, 2, 3)
	ss1.Get(1).Append(10, 20)
	ss1.Get(2).Append(30, 40)
	ss1.Get(3).Append(50, 60)

	ss2 := StreamSetFrom[int, int](2, 3)
	ss2.Get(2).Append(30)
	ss2.Get(3).Append(50)

	result := ss1.MinusStreams(ss2)
	assert.NotNil(t, result)
}

func TestStreamSetFromMap(t *testing.T) {
	m := map[int]*StreamDef[string]{
		1: StreamFrom("a", "b"),
		2: StreamFrom("c", "d"),
	}

	ss := StreamSetFromMap(m)
	assert.NotNil(t, ss)
	assert.Equal(t, 2, ss.Get(1).Len())
	assert.Equal(t, 2, ss.Get(2).Len())
}

func TestNewStreamSet(t *testing.T) {
	ss := NewStreamSet[int, string]()
	assert.NotNil(t, ss)
	assert.Equal(t, 0, ss.Size())
}

func TestSetFromMap(t *testing.T) {
	m := map[int]string{1: "one", 2: "two"}
	s := SetFromMap(m)
	assert.NotNil(t, s)
	assert.True(t, s.ContainsKey(1))
	assert.True(t, s.ContainsKey(2))
}

func TestStreamGetOutOfRange(t *testing.T) {
	s := StreamFrom(1, 2, 3)
	val := s.Get(0)
	assert.Equal(t, 1, val)

	val = s.Get(2)
	assert.Equal(t, 3, val)
}

func TestStreamRemoveItemMultiple(t *testing.T) {
	s := StreamFrom(1, 2, 3, 2, 4, 2, 5)
	s = s.RemoveItem(2)
	assert.Equal(t, []int{1, 3, 4, 5}, s.ToArray())
}

func TestStreamFilter(t *testing.T) {
	s := StreamFrom(1, 2, 3, 4, 5)
	s = s.Filter(func(val int, idx int) bool {
		return val%2 == 0
	})
	assert.Equal(t, []int{2, 4}, s.ToArray())
}

func TestMapSetGet(t *testing.T) {
	s := SetFrom[int, string](1, 2, 3)
	s.Set(1, "one")
	s.Set(2, "two")

	assert.Equal(t, "one", s.Get(1))
	assert.Equal(t, "two", s.Get(2))
	assert.Equal(t, "", s.Get(99))
}

// Tests for stream.go nil/empty input coverage

func TestStreamIsSubsetWithNilInput(t *testing.T) {
	s := StreamFrom(1, 2, 3)
	assert.False(t, s.IsSubset(nil))
}

func TestStreamIsSubsetWithEmptyInput(t *testing.T) {
	s := StreamFrom(1, 2, 3)
	empty := StreamFrom[int]()
	assert.False(t, s.IsSubset(empty))
}

func TestStreamIsSupersetWithNilInput(t *testing.T) {
	s := StreamFrom(1, 2, 3)
	assert.True(t, s.IsSuperset(nil))
}

func TestStreamIsSupersetWithEmptyInput(t *testing.T) {
	s := StreamFrom(1, 2, 3)
	empty := StreamFrom[int]()
	assert.True(t, s.IsSuperset(empty))
}

func TestStreamIntersectionWithNilInput(t *testing.T) {
	s := StreamFrom(1, 2, 3)
	result := s.Intersection(nil)
	assert.Equal(t, []int{}, result.ToArray())
}

func TestStreamIntersectionWithEmptyInput(t *testing.T) {
	s := StreamFrom(1, 2, 3)
	empty := StreamFrom[int]()
	result := s.Intersection(empty)
	assert.Equal(t, []int{}, result.ToArray())
}

func TestStreamMinusWithNilInput(t *testing.T) {
	s := StreamFrom(1, 2, 3)
	result := s.Minus(nil)
	assert.Equal(t, []int{1, 2, 3}, result.ToArray())
}

func TestStreamMinusWithEmptyInput(t *testing.T) {
	s := StreamFrom(1, 2, 3)
	empty := StreamFrom[int]()
	result := s.Minus(empty)
	assert.Equal(t, []int{1, 2, 3}, result.ToArray())
}

func TestMapSetDefBranchCoverage(t *testing.T) {
	s := SetFrom[int, bool](1, 2)

	added := s.Add(3).(*MapSetDef[int, bool])
	assert.True(t, added.ContainsKey(3))
	same := s.Add().(*MapSetDef[int, bool])
	assert.True(t, same.ContainsKey(1))

	unionNil := s.Union(nil).(*MapSetDef[int, bool])
	assert.Equal(t, s.Size(), unionNil.Size())
	unionEmpty := s.Union(SetFromArray[int, bool]([]int{})).(*MapSetDef[int, bool])
	assert.Equal(t, s.Size(), unionEmpty.Size())

	interNil := s.Intersection(nil).(*MapSetDef[int, bool])
	assert.Equal(t, 0, interNil.Size())
	interEmpty := s.Intersection(SetFromArray[int, bool]([]int{})).(*MapSetDef[int, bool])
	assert.Equal(t, 0, interEmpty.Size())
}

func TestStreamConcatEmptySlices(t *testing.T) {
	s := StreamFromArray([]int{1, 2})
	result := s.Concat()
	assert.Equal(t, []int{1, 2}, result.ToArray())
}

func TestStreamRemoveNegativeIndex(t *testing.T) {
	s := StreamFromArray([]int{1, 2, 3})
	result := s.Remove(-1)
	assert.Equal(t, []int{1, 2, 3}, result.ToArray())
}

func TestStreamRemoveOutOfBounds(t *testing.T) {
	s := StreamFromArray([]int{1, 2, 3})
	result := s.Remove(10)
	assert.Equal(t, []int{1, 2, 3}, result.ToArray())
}

func TestStreamRemoveItemEmpty(t *testing.T) {
	s := StreamFromArray([]int{1, 2, 3})
	result := s.RemoveItem()
	assert.Equal(t, []int{1, 2, 3}, result.ToArray())
}

func TestStreamExtendEmpty(t *testing.T) {
	s := StreamFromArray([]int{1, 2, 3})
	result := s.Extend()
	assert.Equal(t, []int{1, 2, 3}, result.ToArray())
}

func TestStreamSetAddEmpty(t *testing.T) {
	s := SetFromArray[int, int]([]int{1, 2})
	result := s.Add()
	assert.Equal(t, 2, result.Size())
}

func TestStreamSetRemoveKeysEmpty(t *testing.T) {
	s := SetFromArray[int, int]([]int{1, 2})
	result := s.RemoveKeys()
	assert.Equal(t, 2, result.Size())
}

func TestStreamSetRemoveValuesEmpty(t *testing.T) {
	s := SetFromArray[int, int]([]int{1, 2})
	result := s.RemoveValues()
	assert.Equal(t, 2, result.Size())
}

func TestStreamSetMinusNilInput(t *testing.T) {
	s := SetFromArray[int, int]([]int{1, 2})
	result := s.Minus(nil)
	assert.Equal(t, 2, result.Size())
}

func TestStreamSetUnionNilInput(t *testing.T) {
	s := SetFromArray[int, int]([]int{1, 2})
	result := s.Union(nil)
	assert.Equal(t, 2, result.Size())
}

func TestStreamSetIntersectionNilInput(t *testing.T) {
	s := SetFromArray[int, int]([]int{1, 2})
	result := s.Intersection(nil)
	assert.Equal(t, 0, result.Size())
}

func TestStreamSetUnionWithData(t *testing.T) {
	ss1 := StreamSetFromArray[int, int]([]int{1, 2})
	ss2 := StreamSetFromArray[int, int]([]int{2, 3})
	result := ss1.Union(ss2)
	assert.Equal(t, 3, result.Size())
	assert.True(t, result.ContainsKey(1))
	assert.True(t, result.ContainsKey(2))
	assert.True(t, result.ContainsKey(3))
}

func TestStreamSetIntersectionWithData(t *testing.T) {
	ss1 := StreamSetFromArray[int, int]([]int{1, 2, 3})
	ss2 := StreamSetFromArray[int, int]([]int{2, 3, 4})
	result := ss1.Intersection(ss2)
	assert.Equal(t, 2, result.Size())
	assert.True(t, result.ContainsKey(2))
	assert.True(t, result.ContainsKey(3))
}

func TestStreamSetMinusStreamsWithData(t *testing.T) {
	ss1 := StreamSetFromArray[int, int]([]int{1, 2})
	ss1.Get(1).Append(10, 20, 30)
	ss1.Get(2).Append(40, 50)

	ss2 := StreamSetFromArray[int, int]([]int{1})
	ss2.Get(1).Append(20, 40)

	result := ss1.MinusStreams(ss2)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Size())
	assert.True(t, result.ContainsKey(1))
	assert.True(t, result.ContainsKey(2))
}

func TestStreamSetDefMinusStreamsNilInput(t *testing.T) {
	ss := StreamSetFromArray[int, int]([]int{1, 2})
	result := ss.MinusStreams(nil)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.Size())
}

func TestStreamSetDefIntersectionNilInput(t *testing.T) {
	ss := StreamSetFromArray[int, int]([]int{1, 2})
	result := ss.Intersection(nil)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.Size())
}

func TestStreamSetDefUnionNilInput(t *testing.T) {
	ss := StreamSetFromArray[int, int]([]int{1, 2})
	result := ss.Union(nil)
	assert.Equal(t, 2, result.Size())
}

func TestMapSetDefAddWithExistingKey(t *testing.T) {
	s := SetFromArray[int, int]([]int{1, 2})
	result := s.Add(1)
	assert.Equal(t, 2, result.Size())
	assert.True(t, result.ContainsKey(1))
	assert.True(t, result.ContainsKey(2))
}

func TestStreamSetDefUnionWithStreamData(t *testing.T) {
	ss1 := StreamSetFromArray[int, int]([]int{1, 2})
	ss1.Set(1, StreamFrom(10, 20))

	ss2 := StreamSetFromArray[int, int]([]int{1, 3})
	ss2.Set(1, StreamFrom(30))

	result := ss1.Union(ss2)
	assert.Equal(t, 3, result.Size())
	assert.True(t, result.ContainsKey(1))
	assert.True(t, result.ContainsKey(2))
	assert.True(t, result.ContainsKey(3))

	v := result.Get(1)
	assert.NotNil(t, v)
	assert.Equal(t, 3, v.Len())
	assert.Equal(t, 10, v.Get(0))
	assert.Equal(t, 20, v.Get(1))
	assert.Equal(t, 30, v.Get(2))
}

func TestStreamSetDefUnionWithNilStream(t *testing.T) {
	ss1 := StreamSetFromArray[int, int]([]int{1, 2})
	ss1.Set(1, nil)

	ss2 := StreamSetFromArray[int, int]([]int{1, 3})
	ss2.Set(1, StreamFrom(10))

	result := ss1.Union(ss2)
	assert.Equal(t, 3, result.Size())
	assert.True(t, result.ContainsKey(1))
	assert.True(t, result.ContainsKey(2))
	assert.True(t, result.ContainsKey(3))

	v := result.Get(1)
	assert.NotNil(t, v)
	assert.Equal(t, 1, v.Len())
	assert.Equal(t, 10, v.Get(0))
}

func TestStreamSetDefIntersectionWithStreamData(t *testing.T) {
	ss1 := StreamSetFromArray[int, int]([]int{1, 2, 3})
	ss1.Set(1, StreamFrom(10, 20, 30))
	ss1.Set(2, StreamFrom(40, 50))

	ss2 := StreamSetFromArray[int, int]([]int{1, 2, 4})
	ss2.Set(1, StreamFrom(20, 30, 40))
	ss2.Set(2, StreamFrom(60))

	result := ss1.Intersection(ss2)
	assert.Equal(t, 2, result.Size())
	assert.True(t, result.ContainsKey(1))
	assert.True(t, result.ContainsKey(2))

	v := result.Get(1)
	assert.NotNil(t, v)
	assert.Equal(t, 2, v.Len())
	assert.Equal(t, 20, v.Get(0))
	assert.Equal(t, 30, v.Get(1))
}

func TestStreamSetDefMinusStreamsWithStreamData(t *testing.T) {
	ss1 := StreamSetFromArray[int, int]([]int{1, 2})
	ss1.Set(1, StreamFrom(10, 20, 30))
	ss1.Set(2, StreamFrom(40, 50))

	ss2 := StreamSetFromArray[int, int]([]int{1, 2, 3})
	ss2.Set(1, StreamFrom(20, 40))

	result := ss1.MinusStreams(ss2)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Size())
	assert.True(t, result.ContainsKey(1))
	assert.True(t, result.ContainsKey(2))

	v := result.Get(1)
	assert.NotNil(t, v)
	assert.Equal(t, 2, v.Len())
	assert.Equal(t, 10, v.Get(0))
	assert.Equal(t, 30, v.Get(1))
}

func TestStreamSetDefMinusStreamsWithNilStream(t *testing.T) {
	ss1 := StreamSetFromArray[int, int]([]int{1, 2})
	ss1.Set(1, nil)

	ss2 := StreamSetFromArray[int, int]([]int{1})
	ss2.Set(1, StreamFrom(10, 20))

	result := ss1.MinusStreams(ss2)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Size())
	assert.True(t, result.ContainsKey(1))
	assert.True(t, result.ContainsKey(2))
}

func TestStreamSetDefIntersectionWithNilStream(t *testing.T) {
	ss1 := StreamSetFromArray[int, int]([]int{1, 2})
	ss1.Set(1, nil)

	ss2 := StreamSetFromArray[int, int]([]int{1, 2})
	ss2.Set(1, StreamFrom(10, 20, 30))

	result := ss1.Intersection(ss2)
	assert.Equal(t, 2, result.Size())
	assert.True(t, result.ContainsKey(1))
	assert.True(t, result.ContainsKey(2))

	v := result.Get(1)
	assert.NotNil(t, v)
	assert.Equal(t, 0, v.Len())
}

func TestStreamExtendWithNilStream(t *testing.T) {
	s := StreamFromArray([]int{1, 2})
	result := s.Extend(nil, nil)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Len())
	assert.Equal(t, 1, result.Get(0))
	assert.Equal(t, 2, result.Get(1))
}
