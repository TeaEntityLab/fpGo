package fpgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromArrayMapReduceForInterface(t *testing.T) {
	var s *StreamForInterfaceDef
	var tempString string

	s = StreamForInterface.FromArrayMaybe([]MaybeDef[interface{}]{Maybe.Just("1"), Maybe.Just("2"), Maybe.Just("3"), Maybe.Just("4")})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)
	s = s.Map(func(item interface{}, index int) interface{} {
		val := Maybe.Just(s.Get(index)).ToMaybe().ToString()
		var result interface{} = "v" + val
		return result
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "v1v2v3v4", tempString)
	tempString = ""

	s = StreamForInterface.FromArrayString([]string{"1", "2", "3", "4"})
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)
	s = s.Map(func(item interface{}, index int) interface{} {
		val := Maybe.Just(s.Get(index)).ToString()
		var result interface{} = "v" + val
		return result
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "v1v2v3v4", tempString)

	s = StreamForInterface.FromArrayInt([]int{1, 2, 3, 4})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)
	s = s.Map(func(item interface{}, index int) interface{} {
		val, _ := Maybe.Just(s.Get(index)).ToInt()
		var result interface{} = val * val
		return result
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "14916", tempString)

	s = StreamForInterface.FromArrayFloat32([]float32{1, 2, 3, 4})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)
	s = s.Map(func(item interface{}, index int) interface{} {
		val, _ := Maybe.Just(s.Get(index)).ToFloat32()
		var result interface{} = val * val
		return result
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "14916", tempString)

	s = StreamForInterface.FromArrayFloat64([]float64{1, 2, 3, 4})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)

	s = s.Map(func(item interface{}, index int) interface{} {
		val, _ := Maybe.Just(s.Get(index)).ToFloat64()
		var result interface{} = val * val
		return result
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "14916", tempString)

	s = StreamForInterface.FromArrayBool([]bool{true, false, true, false})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "truefalsetruefalse", tempString)

	s = StreamForInterface.FromArrayByte([]byte{1, 2, 3, 4})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)

	s = StreamForInterface.FromArrayInt8([]int8{1, 2, 3, 4})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)

	s = StreamForInterface.FromArrayInt16([]int16{1, 2, 3, 4})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)

	s = StreamForInterface.FromArrayInt32([]int32{1, 2, 3, 4})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)

	s = StreamForInterface.FromArrayInt64([]int64{1, 2, 3, 4})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "1234", tempString)
}

func TestFilterForInterface(t *testing.T) {
	var s *StreamForInterfaceDef
	var tempString string

	s = StreamForInterface.FromArrayInt([]int{}).Append(1, 1).Extend(StreamForInterface.FromArrayInt([]int{2, 3, 4})).Extend(StreamForInterface.FromArray([]interface{}{nil})).Extend(nil)
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
		val, err := Maybe.Just(s.Get(index)).ToInt()

		return err == nil && val > 1 && val < 4
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "23", tempString)
	s = s.Reject(func(item interface{}, index int) bool {
		val, err := Maybe.Just(item).ToInt()

		return err == nil && val > 2
	})
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "2", tempString)
}

func TestSortForInterface(t *testing.T) {
	var s *StreamForInterfaceDef
	var tempString string

	s = StreamForInterface.FromArrayInt([]int{11}).Extend(StreamForInterface.FromArrayInt([]int{2, 3, 4, 5})).Remove(4)
	tempString = ""
	for _, v := range s.ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "11234", tempString)

	s = StreamForInterface.FromArray([]interface{}{11}).Concat([]interface{}{2, 3, 4, 5}).RemoveItem(4)
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

	s = StreamForInterface.FromArrayInt([]int{11}).Concat([]interface{}{2, 3, 4, 5}).Remove(4)
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
	for _, v := range s.Sort(func(a, b interface{}) bool {
		vali, _ := Maybe.Just(a).ToInt()
		valj, _ := Maybe.Just(b).ToInt()
		return vali < valj
	}).ToArray() {
		tempString += Maybe.Just(v).ToMaybe().ToString()
	}
	assert.Equal(t, "23411", tempString)
}

func TestStreamForInterfaceSetOperation(t *testing.T) {
	var s *StreamForInterfaceDef
	var s2 *StreamForInterfaceDef
	var s3 *StreamForInterfaceDef
	var tempString string

	s = StreamForInterface.FromArray([]interface{}{11, 2, 3, 4, 5})
	s2 = StreamForInterface.FromArray([]interface{}{9, 2, 5, 6})
	s3 = StreamForInterface.FromArray([]interface{}{2, 5})
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

func TestSetForInterfaceSetOperation(t *testing.T) {
	var s *SetForInterfaceDef
	var s2 *SetForInterfaceDef
	var s3 *SetForInterfaceDef
	var tempString string

	s = SetForInterfaceFrom(11, 2, 3, 4, 5)
	s2 = SetForInterfaceFrom(9, 2, 5, 6)
	s3 = SetForInterfaceFrom(2, 5)
	assert.Equal(t, true, s.ContainsKey(4))
	assert.Equal(t, false, s.ContainsKey(6))
	assert.Equal(t, true, s.IsSupersetByKey(s3))
	assert.Equal(t, true, s2.IsSupersetByKey(s3))
	assert.Equal(t, true, s3.IsSubsetByKey(s))
	assert.Equal(t, true, s3.IsSubsetByKey(s2))
	assert.Equal(t, false, s.IsSupersetByKey(s2))
	tempString = ""
	for _, v := range SortIntAscending(s.Clone().Intersection(s2).Keys()...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "2/5/", tempString)
	assert.Equal(t, 2, s.Intersection(s2).Size())
	tempString = ""
	for _, v := range SortIntAscending(s.Union(s2).Keys()...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "2/3/4/5/6/9/11/", tempString)
	assert.Equal(t, 7, s.Union(s2).Size())
	tempString = ""
	for _, v := range SortIntAscending(s.Minus(s2).Keys()...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "3/4/11/", tempString)
	tempString = ""
	for _, v := range SortIntAscending(s2.Minus(s).Keys()...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "6/9/", tempString)
}

func TestStreamSetForInterfaceSetOperation(t *testing.T) {
	var s *StreamSetForInterfaceDef
	var s2 *StreamSetForInterfaceDef
	var s3 *StreamSetForInterfaceDef
	var tempString string

	s = StreamSetForInterfaceFrom(11, 2, 3, 4, 5)
	s2 = StreamSetForInterfaceFrom(9, 2, 5, 6)
	s3 = StreamSetForInterfaceFrom(2, 5)
	s.Set(2, StreamForInterface.From(70, 71, 72))
	s2.Set(2, StreamForInterface.From(73, 74, 75))
	s2.Set(6, StreamForInterface.From(6, 6, 6))
	s2.Set(9, StreamForInterface.From(9, 9, 9))
	s3.Set(2, StreamForInterface.From(71, 73, 78))

	assert.Equal(t, true, s.ContainsKey(4))
	assert.Equal(t, false, s.ContainsKey(6))
	assert.Equal(t, true, s.IsSupersetByKey(s3))
	assert.Equal(t, true, s2.IsSupersetByKey(s3))
	assert.Equal(t, true, s3.IsSubsetByKey(s))
	assert.Equal(t, true, s3.IsSubsetByKey(s2))
	assert.Equal(t, false, s.IsSupersetByKey(s2))
	tempString = ""
	for _, v := range SortIntAscending(s.Clone().Intersection(s2).Keys()...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "2/5/", tempString)
	assert.Equal(t, 2, s.Intersection(s2).Size())
	tempString = ""
	for _, v := range SortIntAscending(s.Union(s2).Keys()...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "2/3/4/5/6/9/11/", tempString)
	assert.Equal(t, 7, s.Union(s2).Size())
	tempString = ""
	for _, v := range SortIntAscending(s.Minus(s2).Keys()...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "3/4/11/", tempString)
	tempString = ""
	for _, v := range SortIntAscending(s2.Minus(s).Keys()...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "6/9/", tempString)

	tempString = ""
	for _, v := range SortStringAscending(Map(streamIntTransformerForInterface, s2.Union(s).Values()...)...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "6,6,6,end/70,71,72,73,74,75,end/9,9,9,end/end/end/end/end/", tempString)
	tempString = ""
	for _, v := range SortStringAscending(Map(streamIntTransformerForInterface, s2.Minus(s).Values()...)...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "6,6,6,end/9,9,9,end/", tempString)
	tempString = ""
	for _, v := range SortStringAscending(Map(streamIntTransformerForInterface, s2.Intersection(s).Values()...)...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "end/end/", tempString)
	tempString = ""
	for _, v := range SortStringAscending(Map(streamIntTransformerForInterface, s.Intersection(s2).Values()...)...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "end/end/", tempString)
	tempString = ""
	for _, v := range SortStringAscending(Map(streamIntTransformerForInterface, s.Intersection(s3).Values()...)...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "71,end/end/", tempString)
	tempString = ""
	for _, v := range SortStringAscending(Map(streamIntTransformerForInterface, s2.Intersection(s3).Values()...)...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "73,end/end/", tempString)
	tempString = ""
	for _, v := range SortStringAscending(Map(streamIntTransformerForInterface, s.MinusStreams(s3).Values()...)...) {
		tempString += Maybe.Just(v).ToMaybe().ToString() + "/"
	}
	assert.Equal(t, "70,72,end/end/end/end/end/", tempString)
}

func streamIntTransformerForInterface(s interface{}) interface{} {
	result := ""
	for _, item := range SortIntAscending(s.(*StreamForInterfaceDef).ToArray()...) {
		result += Maybe.Just(item).ToMaybe().ToString() + ","
	}
	return result + "end"
}
