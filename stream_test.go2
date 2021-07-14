package fpgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromArrayMapReduce(t *testing.T) {
	var s *StreamDef[interface{}]
	var tempString string

	s = StreamFromInterface(Maybe.Just("1"), Maybe.Just("2"), Maybe.Just("3"), Maybe.Just("4"))
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)
	s = s.Map(func(item interface{}, index int) interface{} {
		var val = Maybe.Just(item).ToMaybe().ToString()
		var result interface{} = "v" + val
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
		var val = Maybe.Just(item).ToString()
		var result string = "v" + val
		return result
	})
	tempString = ""
	for _, v := range s2.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "v1v2v3v4", tempString)

	s = StreamFromInterface(1, 2, 3, 4)
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)
	s = s.Map(func(item interface{}, index int) interface{} {
		var val, _ = Maybe.Just(item).ToInt()
		var result interface{} = val * val
		return result
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "14916", tempString)

	s = StreamFromInterface(1, 2, 3, 4)
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)
	s = s.Map(func(item interface{}, index int) interface{} {
		var val, _ = Maybe.Just(item).ToFloat32()
		var result interface{} = val * val
		return result
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "14916", tempString)

	s = StreamFromInterface(1, 2, 3, 4)
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)

	s = s.Map(func(item interface{}, index int) interface{} {
		var val, _ = Maybe.Just(item).ToFloat64()
		var result interface{} = val * val
		return result
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "14916", tempString)

	s = StreamFromInterface(true, false, true, false)
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "truefalsetruefalse", tempString)

	s = StreamFromInterface(1, 2, 3, 4)
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)

	s = StreamFromInterface(1, 2, 3, 4)
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)

	s = StreamFromInterface(1, 2, 3, 4)
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)

	s = StreamFromInterface(1, 2, 3, 4)
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)

	s = StreamFromInterface(1, 2, 3, 4)
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)
}

func TestFilter(t *testing.T) {
	var s *StreamDef[interface{}]
	var tempString string

	s = StreamFromArray([]interface{}{}).Append(1, 1).Extend(StreamFromArray([]interface{}{2, 3, 4})).Extend(StreamFromArray([]interface{}{nil})).Extend(nil)
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "11234<nil>", tempString)
	s = s.Distinct()
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234<nil>", tempString)
	s = s.FilterNotNil()
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)

	s = s.Filter(func(item interface{}, index int) bool {
		var val, err = Maybe.Just(item).ToInt()

		return err == nil && val > 1 && val < 4
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "23", tempString)
	s = s.Reject(func(item interface{}, index int) bool {
		var val, err = Maybe.Just(item).ToInt()

		return err == nil && val > 2
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "2", tempString)
}

func TestSort(t *testing.T) {
	var s *StreamDef[interface{}]
	var tempString string

	s = StreamFromInterface(11).Extend(StreamFromInterface(2, 3, 4, 5)).Remove(4)
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "11234", tempString)

	s = StreamFromInterface(11).Concat([]interface{}{2, 3, 4, 5}).RemoveItem(4)
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

	s = StreamFromInterface(11).Concat([]interface{}{2, 3, 4, 5}).Remove(4)
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "11234", tempString)

	tempString = ""
	for _, v := range s.SortByIndex(func(i, j int) bool {
		var vali, _ = Maybe.Just(s.Get(i)).ToInt()
		var valj, _ = Maybe.Just(s.Get(j)).ToInt()
		return vali < valj
	}).ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "23411", tempString)
	tempString = ""
	for _, v := range s.Sort(func(a, b interface{}) bool {
		var vali, _ = Maybe.Just(a).ToInt()
		var valj, _ = Maybe.Just(b).ToInt()
		return vali < valj
	}).ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "23411", tempString)
}

func TestStreamSetOperation(t *testing.T) {
	var s *StreamDef[interface{}]
	var s2 *StreamDef[interface{}]
	var s3 *StreamDef[interface{}]
	var tempString string

	s = StreamFromInterface(11, 2, 3, 4, 5)
	s2 = StreamFromInterface(9, 2, 5, 6)
	s3 = StreamFromInterface(2, 5)
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
	var s *SetDef[int, bool]
	var s2 *SetDef[int, bool]
	var s3 *SetDef[int, bool]
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
	var s *StreamSetDef[int, int]
	var s2 *StreamSetDef[int, int]
	var s3 *StreamSetDef[int, int]
	var tempString string

	s = StreamSetFrom[int, int](11, 2, 3, 4, 5)
	s2 = StreamSetFrom[int, int](9, 2, 5, 6)
	s3 = StreamSetFrom[int, int](2, 5)
	s.Set(2, StreamFrom(70, 71, 72))
	s2.Set(2, StreamFrom(73, 74, 75))
	s2.Set(6, StreamFrom(6, 6, 6))
	s2.Set(9, StreamFrom(9, 9, 9))
	s3.Set(2, StreamFrom(71, 73, 78))

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

	tempString = ""
	for _, v := range SortOrderedAscending(Map(streamIntTransformer, s2.Union(s).Values()...)...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "6,6,6,end/70,71,72,73,74,75,end/9,9,9,end/end/end/end/end/", tempString)
	tempString = ""
	for _, v := range SortOrderedAscending(Map(streamIntTransformer, s2.Minus(s).Values()...)...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "6,6,6,end/9,9,9,end/", tempString)
	tempString = ""
	for _, v := range SortOrderedAscending(Map(streamIntTransformer, s2.Intersection(s).Values()...)...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "end/end/", tempString)
	tempString = ""
	for _, v := range SortOrderedAscending(Map(streamIntTransformer, s.Intersection(s2).Values()...)...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "end/end/", tempString)
	tempString = ""
	for _, v := range SortOrderedAscending(Map(streamIntTransformer, s.Intersection(s3).Values()...)...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "71,end/end/", tempString)
	tempString = ""
	for _, v := range SortOrderedAscending(Map(streamIntTransformer, s2.Intersection(s3).Values()...)...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "73,end/end/", tempString)
	tempString = ""
	for _, v := range SortOrderedAscending(Map(streamIntTransformer, s.MinusStreams(s3).Values()...)...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "70,72,end/end/end/end/end/", tempString)
}

func streamIntTransformer(s *StreamDef[int]) string {
 result := ""
 for _, item := range(SortOrderedAscending(s.ToArray()...)) {
	 result += Maybe.Just(item).ToMaybe().ToString() + ","
 }
 return result + "end"
}
