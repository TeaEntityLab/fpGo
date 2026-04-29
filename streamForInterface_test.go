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

func TestStreamForInterfaceFrom(t *testing.T) {
	s := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	assert.Equal(t, 3, s.Len())
}

func TestStreamForInterfaceFromArray(t *testing.T) {
	s := StreamForInterface.FromArray([]interface{}{1, 2})
	assert.Equal(t, 2, s.Len())
}

func TestStreamSetForInterfaceFromMap(t *testing.T) {
	m := map[interface{}]*StreamForInterfaceDef{
		1: StreamForInterface.FromArrayInt([]int{1, 2}),
		2: StreamForInterface.FromArrayInt([]int{3, 4}),
	}
	s := StreamSetForInterfaceFromMap(m)
	assert.Equal(t, 2, s.Size())
}

func TestStreamForInterface_Contains(t *testing.T) {
	s := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	assert.True(t, s.Contains(2))
	assert.False(t, s.Contains(5))
}

func TestStreamForInterface_IsSubset(t *testing.T) {
	s1 := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	s2 := StreamForInterface.FromArrayInt([]int{1, 2, 3, 4, 5})
	assert.True(t, s1.IsSubset(s2))
	assert.False(t, s2.IsSubset(s1))
}

func TestStreamForInterface_IsSuperset(t *testing.T) {
	s1 := StreamForInterface.FromArrayInt([]int{1, 2, 3, 4, 5})
	s2 := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	assert.True(t, s1.IsSuperset(s2))
	assert.False(t, s2.IsSuperset(s1))
}

func TestStreamForInterface_Clone(t *testing.T) {
	s1 := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	s2 := s1.Clone()
	assert.Equal(t, s1.ToArray(), s2.ToArray())
	assert.Equal(t, 3, s2.Len())
}

func TestStreamForInterface_Intersection(t *testing.T) {
	s1 := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	s2 := StreamForInterface.FromArrayInt([]int{2, 3, 4})
	result := s1.Intersection(s2)
	assert.Equal(t, 2, result.Len())
}

func TestStreamForInterface_Minus(t *testing.T) {
	s1 := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	s2 := StreamForInterface.FromArrayInt([]int{2})
	result := s1.Minus(s2)
	assert.Equal(t, 2, result.Len())
}

func TestStreamForInterface_RemoveItem(t *testing.T) {
	s := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	result := s.RemoveItem(2)
	assert.Equal(t, 2, result.Len())
}

func TestStreamForInterface_Append(t *testing.T) {
	s := StreamForInterface.FromArrayInt([]int{1, 2})
	result := s.Append(3, 4)
	assert.Equal(t, 4, result.Len())
}

func TestStreamForInterface_Remove(t *testing.T) {
	s := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	result := s.Remove(1)
	assert.Equal(t, 2, result.Len())
}

func TestStreamForInterface_Len(t *testing.T) {
	s := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	assert.Equal(t, 3, s.Len())
}

func TestStreamForInterface_Concat(t *testing.T) {
	s := StreamForInterface.FromArrayInt([]int{1, 2})
	result := s.Concat([]interface{}{3, 4})
	assert.Equal(t, 4, result.Len())
}

func TestStreamForInterface_Extend(t *testing.T) {
	s1 := StreamForInterface.FromArrayInt([]int{1, 2})
	s2 := StreamForInterface.FromArrayInt([]int{3, 4})
	result := s1.Extend(s2)
	assert.Equal(t, 4, result.Len())
}

func TestStreamForInterface_Reverse(t *testing.T) {
	s := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	result := s.Reverse()
	arr := result.ToArray()
	assert.Equal(t, 3, arr[0].(int))
	assert.Equal(t, 1, arr[2].(int))
}

func TestStreamForInterface_Sort(t *testing.T) {
	s := StreamForInterface.FromArrayInt([]int{3, 1, 2})
	result := s.Sort(func(a, b interface{}) bool {
		return a.(int) < b.(int)
	})
	arr := result.ToArray()
	assert.Equal(t, 1, arr[0].(int))
	assert.Equal(t, 3, arr[2].(int))
}

func TestStreamForInterface_Get(t *testing.T) {
	s := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	assert.Equal(t, 2, s.Get(1))
}

func TestSetForInterfaceFrom(t *testing.T) {
	s := SetForInterfaceFrom(1, 2, 3)
	assert.Equal(t, 3, s.Size())
}

func TestSetForInterfaceFromArray(t *testing.T) {
	s := SetForInterfaceFromArray([]interface{}{1, 2, 3})
	assert.Equal(t, 3, s.Size())
}

func TestSetForInterfaceFromMap(t *testing.T) {
	m := map[interface{}]interface{}{1: "a", 2: "b"}
	s := SetForInterfaceFromMap(m)
	assert.Equal(t, 2, s.Size())
}

func TestSetForInterfaceDef_ContainsKey(t *testing.T) {
	s := SetForInterfaceFrom(1, 2, 3)
	assert.True(t, s.ContainsKey(2))
	assert.False(t, s.ContainsKey(5))
}

func TestSetForInterfaceDef_ContainsValue(t *testing.T) {
	s := SetForInterfaceFrom(1, 2, 3)
	(*s)[10] = "value"
	assert.True(t, s.ContainsValue("value"))
	assert.False(t, s.ContainsValue("nonexistent"))
}

func TestSetForInterfaceDef_IsSubsetByKey(t *testing.T) {
	s1 := SetForInterfaceFrom(1, 2)
	s2 := SetForInterfaceFrom(1, 2, 3, 4)
	assert.True(t, s1.IsSubsetByKey(s2))
	assert.False(t, s2.IsSubsetByKey(s1))
}

func TestSetForInterfaceDef_IsSupersetByKey(t *testing.T) {
	s1 := SetForInterfaceFrom(1, 2, 3, 4)
	s2 := SetForInterfaceFrom(1, 2)
	assert.True(t, s1.IsSupersetByKey(s2))
	assert.False(t, s2.IsSupersetByKey(s1))
}

func TestSetForInterfaceDef_Add(t *testing.T) {
	s := SetForInterfaceFrom(1, 2)
	result := s.Add(3, 4)
	assert.Equal(t, 4, result.Size())
}

func TestSetForInterfaceDef_RemoveKeys(t *testing.T) {
	s := SetForInterfaceFrom(1, 2, 3)
	result := s.RemoveKeys(2)
	assert.Equal(t, 2, result.Size())
}

func TestStreamSetForInterface_Clone(t *testing.T) {
	s := StreamSetForInterfaceFromArray([]interface{}{1, 2})
	cloned := s.Clone()
	assert.Equal(t, s.Size(), cloned.Size())
}

func TestStreamSetForInterface_Union(t *testing.T) {
	s1 := StreamSetForInterfaceFromArray([]interface{}{1, 2})
	s2 := StreamSetForInterfaceFromArray([]interface{}{3, 4})
	result := s1.Union(s2)
	assert.Equal(t, 4, result.Size())
}

func TestStreamSetForInterface_Intersection(t *testing.T) {
	s1 := StreamSetForInterfaceFromArray([]interface{}{1, 2, 3})
	s2 := StreamSetForInterfaceFromArray([]interface{}{2, 3, 4})
	result := s1.Intersection(s2)
	assert.Equal(t, 2, result.Size())
}

func TestStreamSetForInterface_Minus(t *testing.T) {
	s1 := StreamSetForInterfaceFromArray([]interface{}{1, 2, 3})
	s2 := StreamSetForInterfaceFromArray([]interface{}{2})
	result := s1.Minus(s2)
	assert.Equal(t, 2, result.Size())
}

func TestStreamSetForInterface_MinusStreams(t *testing.T) {
	m1 := map[interface{}]*StreamForInterfaceDef{
		1: StreamForInterface.FromArrayInt([]int{1, 2}),
		2: StreamForInterface.FromArrayInt([]int{3, 4}),
	}
	m2 := map[interface{}]*StreamForInterfaceDef{
		1: StreamForInterface.FromArrayInt([]int{1}),
	}
	s1 := StreamSetForInterfaceFromMap(m1)
	s2 := StreamSetForInterfaceFromMap(m2)
	result := s1.MinusStreams(s2)
	assert.Equal(t, 2, result.Size())
}

func TestStreamSetForInterface_IsSubsetByKey(t *testing.T) {
	s1 := StreamSetForInterfaceFromArray([]interface{}{1, 2})
	s2 := StreamSetForInterfaceFromArray([]interface{}{1, 2, 3})
	assert.True(t, s1.IsSubsetByKey(s2))
}

func TestStreamSetForInterface_IsSupersetByKey(t *testing.T) {
	s1 := StreamSetForInterfaceFromArray([]interface{}{1, 2, 3})
	s2 := StreamSetForInterfaceFromArray([]interface{}{1, 2})
	assert.True(t, s1.IsSupersetByKey(s2))
}

// Tests for streamForInterface.go nil/empty input coverage

func TestStreamForInterfaceIsSubsetWithNilInput(t *testing.T) {
	s := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	assert.False(t, s.IsSubset(nil))
}

func TestStreamForInterfaceIsSubsetWithEmptyInput(t *testing.T) {
	s := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	empty := StreamForInterface.FromArrayInt([]int{})
	assert.False(t, s.IsSubset(empty))
}

func TestStreamForInterfaceIsSupersetWithNilInput(t *testing.T) {
	s := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	assert.True(t, s.IsSuperset(nil))
}

func TestStreamForInterfaceIsSupersetWithEmptyInput(t *testing.T) {
	s := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	empty := StreamForInterface.FromArrayInt([]int{})
	assert.True(t, s.IsSuperset(empty))
}

func TestStreamForInterfaceIntersectionWithNilInput(t *testing.T) {
	s := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	result := s.Intersection(nil)
	assert.Equal(t, []interface{}{}, result.ToArray())
}

func TestStreamForInterfaceIntersectionWithEmptyInput(t *testing.T) {
	s := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	empty := StreamForInterface.FromArrayInt([]int{})
	result := s.Intersection(empty)
	assert.Equal(t, []interface{}{}, result.ToArray())
}

func TestStreamForInterfaceMinusWithNilInput(t *testing.T) {
	s := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	result := s.Minus(nil)
	assert.Equal(t, []interface{}{1, 2, 3}, result.ToArray())
}

func TestStreamForInterfaceMinusWithEmptyInput(t *testing.T) {
	s := StreamForInterface.FromArrayInt([]int{1, 2, 3})
	empty := StreamForInterface.FromArrayInt([]int{})
	result := s.Minus(empty)
	assert.Equal(t, []interface{}{1, 2, 3}, result.ToArray())
}

func TestSetForInterfaceIntersectionWithNilInput(t *testing.T) {
	s := SetForInterfaceFromArray([]interface{}{1, 2, 3})
	result := s.Intersection(nil)
	assert.Equal(t, 0, result.Size())
}

func TestSetForInterfaceIntersectionWithEmptyInput(t *testing.T) {
	s := SetForInterfaceFromArray([]interface{}{1, 2, 3})
	empty := SetForInterfaceFromArray([]interface{}{})
	result := s.Intersection(empty)
	assert.Equal(t, 0, result.Size())
}

func TestSetForInterfaceMinusWithNilInput(t *testing.T) {
	s := SetForInterfaceFromArray([]interface{}{1, 2, 3})
	result := s.Minus(nil)
	assert.Equal(t, 3, result.Size())
}

func TestSetForInterfaceMinusWithEmptyInput(t *testing.T) {
	s := SetForInterfaceFromArray([]interface{}{1, 2, 3})
	empty := SetForInterfaceFromArray([]interface{}{})
	result := s.Minus(empty)
	assert.Equal(t, 3, result.Size())
}

func TestSetForInterface_HelperBranches(t *testing.T) {
	s := SetForInterfaceFrom(1, 2)
	s.Set(1, "a")
	s.Set(2, "b")

	mappedKey := s.MapKey(func(input interface{}) interface{} { return input.(int) + 10 })
	assert.True(t, mappedKey.ContainsKey(11))
	assert.Equal(t, "a", mappedKey.Get(11))

	mappedValue := s.MapValue(func(input interface{}) interface{} { return input.(string) + "!" })
	assert.Equal(t, "a!", mappedValue.Get(1))
	assert.Equal(t, "b!", mappedValue.Get(2))

	removed := s.RemoveValues("b")
	assert.True(t, removed.ContainsKey(1))
	assert.False(t, removed.ContainsKey(2))

	unchanged := s.RemoveValues()
	assert.True(t, unchanged.ContainsKey(1))
	assert.True(t, unchanged.ContainsKey(2))
	assert.Equal(t, "a", s.Get(1))
}

func TestStreamSetForInterface_AliasConstructors(t *testing.T) {
	s1 := StreamSetFromInterface(1, 2)
	assert.Equal(t, 2, s1.Size())
	assert.True(t, s1.ContainsKey(1))
	assert.True(t, s1.ContainsKey(2))

	s2 := StreamSetFromArrayInterface([]interface{}{"x", "y"})
	assert.Equal(t, 2, s2.Size())
	assert.True(t, s2.ContainsKey("x"))
	assert.True(t, s2.ContainsKey("y"))
}
