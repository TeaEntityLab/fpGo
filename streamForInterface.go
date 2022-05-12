package fpgo

import (
	"sort"
)

// Stream

// StreamForInterfaceDef Stream inspired by Collection utils
type StreamForInterfaceDef []interface{}

// FromArrayMaybe FromArrayMaybe New Stream instance from a Maybe array
func (streamSelf *StreamForInterfaceDef) FromArrayMaybe(old []MaybeDef[interface{}]) *StreamForInterfaceDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayString New Stream instance from a string array
func (streamSelf *StreamForInterfaceDef) FromArrayString(old []string) *StreamForInterfaceDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayBool New Stream instance from a bool array
func (streamSelf *StreamForInterfaceDef) FromArrayBool(old []bool) *StreamForInterfaceDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayInt New Stream instance from an int array
func (streamSelf *StreamForInterfaceDef) FromArrayInt(old []int) *StreamForInterfaceDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayByte New Stream instance from an int8 array
func (streamSelf *StreamForInterfaceDef) FromArrayByte(old []byte) *StreamForInterfaceDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayInt8 New Stream instance from an int8 array
func (streamSelf *StreamForInterfaceDef) FromArrayInt8(old []int8) *StreamForInterfaceDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayInt16 New Stream instance from an int16 array
func (streamSelf *StreamForInterfaceDef) FromArrayInt16(old []int16) *StreamForInterfaceDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayInt32 New Stream instance from an int32 array
func (streamSelf *StreamForInterfaceDef) FromArrayInt32(old []int32) *StreamForInterfaceDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayInt64 New Stream instance from an int64 array
func (streamSelf *StreamForInterfaceDef) FromArrayInt64(old []int64) *StreamForInterfaceDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayFloat32 New Stream instance from a float32 array
func (streamSelf *StreamForInterfaceDef) FromArrayFloat32(old []float32) *StreamForInterfaceDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayFloat64 New Stream instance from a float64 array
func (streamSelf *StreamForInterfaceDef) FromArrayFloat64(old []float64) *StreamForInterfaceDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// From New Stream instance from a interface{} array
func (streamSelf *StreamForInterfaceDef) From(list ...interface{}) *StreamForInterfaceDef {
	return StreamForInterface.FromArray(list)
}

// FromArray New Stream instance from an interface{} array
func (streamSelf *StreamForInterfaceDef) FromArray(list []interface{}) *StreamForInterfaceDef {
	result := StreamForInterfaceDef(list)
	return &result
}

// ToArray Convert Stream to slice
func (streamSelf *StreamForInterfaceDef) ToArray() []interface{} {
	return DuplicateSlice(*streamSelf)
}

// Map Map all items of Stream by function
func (streamSelf *StreamForInterfaceDef) Map(fn func(interface{}, int) interface{}) *StreamForInterfaceDef {
	result := StreamForInterface.FromArray(MapIndexed(fn, (*streamSelf)...))

	return result
}

// Filter Filter items of Stream by function
func (streamSelf *StreamForInterfaceDef) Filter(fn func(interface{}, int) bool) *StreamForInterfaceDef {
	result := StreamForInterface.FromArray(Filter(fn, (*streamSelf)...))

	return result
}

// Reject Reject items of Stream by function
func (streamSelf *StreamForInterfaceDef) Reject(fn func(interface{}, int) bool) *StreamForInterfaceDef {
	result := StreamForInterface.FromArray(Reject(fn, (*streamSelf)...))

	return result
}

// FilterNotNil Filter not nil items and return a new Stream instance
func (streamSelf *StreamForInterfaceDef) FilterNotNil() *StreamForInterfaceDef {
	return streamSelf.Filter(func(val interface{}, i int) bool {
		return Maybe.Just(val).IsPresent()
	})
}

// Distinct Filter duplicated items and return a new Stream instance
func (streamSelf *StreamForInterfaceDef) Distinct() *StreamForInterfaceDef {
	return StreamForInterface.FromArray(DistinctForInterface(*streamSelf...))
}

// Contains Check the item exists or not in the Stream
func (streamSelf *StreamForInterfaceDef) Contains(input interface{}) bool {
	return ExistsForInterface(input, *streamSelf...)
}

// IsSubset returns true or false by checking if stream1 is a subset of stream2
func (streamSelf *StreamForInterfaceDef) IsSubset(input *StreamForInterfaceDef) bool {
	if input == nil || input.Len() == 0 {
		return false
	}

	return IsSubsetForInterface(*streamSelf, *input)
}

// IsSuperset returns true or false by checking if stream1 is a superset of stream2
func (streamSelf *StreamForInterfaceDef) IsSuperset(input *StreamForInterfaceDef) bool {
	if input == nil || input.Len() == 0 {
		return true
	}

	return IsSupersetForInterface(*streamSelf, *input)
}

// Clone Clone this Stream
func (streamSelf *StreamForInterfaceDef) Clone() *StreamForInterfaceDef {
	return StreamForInterface.FromArray(DuplicateSlice(*streamSelf))
}

// Intersection Get the Intersection with this Stream and an another Stream
func (streamSelf *StreamForInterfaceDef) Intersection(input *StreamForInterfaceDef) *StreamForInterfaceDef {
	if input == nil || input.Len() == 0 {
		return new(StreamForInterfaceDef)
	}

	result := StreamForInterface.FromArray(IntersectionForInterface(*streamSelf, *input))

	return result
}

// Minus Get all of this Stream but not in the given Stream
func (streamSelf *StreamForInterfaceDef) Minus(input *StreamForInterfaceDef) *StreamForInterfaceDef {
	if input == nil || input.Len() == 0 {
		return streamSelf
	}

	result := StreamForInterface.FromArray(MinusForInterface(*streamSelf, *input))
	return result
}

// RemoveItem Remove items from the Stream
func (streamSelf *StreamForInterfaceDef) RemoveItem(input ...interface{}) *StreamForInterfaceDef {
	inputLen := len(input)
	if inputLen > 0 {
		result := StreamForInterface.FromArray(MinusForInterface(*streamSelf, input))

		return result
	}

	return streamSelf
}

// Append Append an item to Stream
func (streamSelf *StreamForInterfaceDef) Append(item ...interface{}) *StreamForInterfaceDef {
	return streamSelf.Concat(item)
}

// Remove Remove an item by its index
func (streamSelf *StreamForInterfaceDef) Remove(index int) *StreamForInterfaceDef {
	if index >= 0 && index < streamSelf.Len() {
		(*streamSelf) = append((*streamSelf)[:index], (*streamSelf)[index+1:]...)
	}
	return streamSelf
}

// Len Get length of Stream
func (streamSelf *StreamForInterfaceDef) Len() int {
	return len(*streamSelf)
}

// Concat Concat Stream by another slices
func (streamSelf *StreamForInterfaceDef) Concat(slices ...[]interface{}) *StreamForInterfaceDef {
	if len(slices) == 0 {
		return streamSelf
	}

	return StreamForInterface.FromArray(Concat(streamSelf.ToArray(), slices...))
}

// Extend Extend Stream by another Stream(s)
func (streamSelf *StreamForInterfaceDef) Extend(streams ...*StreamForInterfaceDef) *StreamForInterfaceDef {
	if len(streams) == 0 {
		return streamSelf
	}

	mine := *streamSelf
	mineLen := len(mine)
	totalLen := mineLen

	for _, stream := range streams {
		if stream == nil {
			continue
		}

		targetLen := len(*stream)
		totalLen += targetLen
	}
	newOne := make([]interface{}, totalLen)

	for i, item := range mine {
		newOne[i] = item
	}
	totalIndex := mineLen

	for _, stream := range streams {
		if stream == nil {
			continue
		}

		target := *stream
		targetLen := len(target)
		for j, item := range target {
			newOne[totalIndex+j] = item
		}
		totalIndex += targetLen
	}

	return StreamForInterface.FromArray(newOne)
}

// Reverse Reverse Stream items
func (streamSelf *StreamForInterfaceDef) Reverse() *StreamForInterfaceDef {
	return StreamForInterface.FromArray(Reverse(*streamSelf...))
}

// SortByIndex Sort Stream items by function(index, index) bool
func (streamSelf *StreamForInterfaceDef) SortByIndex(fn func(a, b int) bool) *StreamForInterfaceDef {
	// Keep the old value
	oldValue := streamSelf.Clone()
	// Make the target for sorting (original)
	result := *streamSelf
	sort.SliceStable(result, fn)
	// Replace values back
	*streamSelf = *oldValue

	// Return the sorted target
	return &result
}

// Sort Sort Stream items by Comparator
func (streamSelf *StreamForInterfaceDef) Sort(fn Comparator[interface{}]) *StreamForInterfaceDef {
	result := streamSelf.Clone()
	Sort(fn, *result)
	return result
}

// Get Get an item of Stream by its index
func (streamSelf *StreamForInterfaceDef) Get(i int) interface{} {
	return (*streamSelf)[i]
}

// StreamForInterface Stream utils instance
var StreamForInterface StreamForInterfaceDef

// Set

// SetForInterfaceDef Set inspired by Collection utils
type SetForInterfaceDef map[interface{}]interface{}

// SetForInterfaceFrom New Set instance from a interface{} array
func SetForInterfaceFrom(list ...interface{}) *SetForInterfaceDef {
	return SetForInterfaceFromArray(list)
}

// SetForInterfaceFromArray New Set instance from a interface{} array
func SetForInterfaceFromArray(list []interface{}) *SetForInterfaceDef {
	newOne := SetForInterfaceDef(SliceToMapForInterface(*new(interface{}), list...))
	return &newOne
}

// SetForInterfaceFromMap New Set instance from a map[interface{}]R
func SetForInterfaceFromMap(theMap map[interface{}]interface{}) *SetForInterfaceDef {
	return SetForInterfaceFromArray(KeysForInterface(theMap))
}

// MapKey Map all keys of Set by function
func (setSelf *SetForInterfaceDef) MapKey(fn TransformerFunctor[interface{}, interface{}]) *SetForInterfaceDef {
	result := make(SetForInterfaceDef, len(*setSelf))
	for k, v := range *setSelf {
		result[fn(k)] = v
	}

	return &result
}

// MapValue Map all values of Set by function
func (setSelf *SetForInterfaceDef) MapValue(fn TransformerFunctor[interface{}, interface{}]) *SetForInterfaceDef {
	result := make(SetForInterfaceDef, len(*setSelf))
	for k, v := range *setSelf {
		result[k] = fn(v)
	}

	return &result
}

// ContainsKey Check the key exists or not in the Set
func (setSelf *SetForInterfaceDef) ContainsKey(input interface{}) bool {
	_, ok := (*setSelf)[input]
	return ok
}

// ContainsValue Check the value exists or not in the Set
func (setSelf *SetForInterfaceDef) ContainsValue(input interface{}) bool {
	for _, v := range *setSelf {
		if v == input {
			return true
		}
	}
	return false
}

// IsSubsetByKey returns true or false by checking if set1 is a subset of set2
func (setSelf *SetForInterfaceDef) IsSubsetByKey(input *SetForInterfaceDef) bool {
	return IsSubsetMapByKeyForInterface(*setSelf, *input)
}

// IsSupersetByKey returns true or false by checking if set1 is a superset of set2
func (setSelf *SetForInterfaceDef) IsSupersetByKey(input *SetForInterfaceDef) bool {
	return IsSupersetMapByKeyForInterface(*setSelf, *input)
}

// Add Add items into the Set
func (setSelf *SetForInterfaceDef) Add(input ...interface{}) *SetForInterfaceDef {
	inputLen := len(input)
	if inputLen > 0 {
		result := setSelf.Clone()
		for _, v := range input {
			if _, ok := (*result)[v]; ok {
				continue
			}
			(*result)[v] = true
		}

		return result
	}

	return setSelf
}

// RemoveKeys Remove keys from the Set
func (setSelf *SetForInterfaceDef) RemoveKeys(input ...interface{}) *SetForInterfaceDef {
	inputLen := len(input)
	if inputLen > 0 {
		result := setSelf.Clone()
		for _, v := range input {
			delete(*result, v)
		}

		return result
	}

	return setSelf
}

// RemoveValues Remove values from the Set
func (setSelf *SetForInterfaceDef) RemoveValues(input ...interface{}) *SetForInterfaceDef {
	inputLen := len(input)
	if inputLen > 0 {
		result := setSelf.Clone()
		valueMap := SliceToMapForInterface(0, input...)
		for k, v := range *setSelf {
			if _, ok := valueMap[v]; ok {
				delete(*result, k)
			}
		}

		return result
	}

	return setSelf
}

// Get Get items from the Set
func (setSelf *SetForInterfaceDef) Get(key interface{}) interface{} {
	return (*setSelf)[key]
}

// Set Set items to the Set
func (setSelf *SetForInterfaceDef) Set(key interface{}, value interface{}) {
	(*setSelf)[key] = value

	// return setSelf
}

// Clone Clone this Set
func (setSelf *SetForInterfaceDef) Clone() *SetForInterfaceDef {
	result := SetForInterfaceDef(DuplicateMapForInterface(*setSelf))

	return &result
}

// Union Union an another Set object
func (setSelf *SetForInterfaceDef) Union(input *SetForInterfaceDef) *SetForInterfaceDef {
	if input == nil || input.Size() == 0 {
		return setSelf
	}

	result := SetForInterfaceDef(MergeForInterface(*setSelf, *input))

	return &result
}

// Intersection Get the Intersection with this Set and an another Set
func (setSelf *SetForInterfaceDef) Intersection(input *SetForInterfaceDef) *SetForInterfaceDef {
	if input == nil || input.Size() == 0 {
		return new(SetForInterfaceDef)
	}

	result := SetForInterfaceDef(IntersectionMapByKeyForInterface(*setSelf, *input))

	return &result
}

// Minus Get all of this Set but not in the given Set
func (setSelf *SetForInterfaceDef) Minus(input *SetForInterfaceDef) *SetForInterfaceDef {
	if input == nil || input.Size() == 0 {
		return setSelf
	}

	result := setSelf.Clone()
	for k := range *result {
		_, exists := (*input)[k]
		if exists {
			delete(*result, k)
		}
	}

	return result
}

// Size Get size
func (setSelf *SetForInterfaceDef) Size() int {
	return len(*setSelf)
}

// Keys Convert Set to slice
func (setSelf *SetForInterfaceDef) Keys() []interface{} {
	return KeysForInterface(*setSelf)
}

// Values Convert Set to slice
func (setSelf *SetForInterfaceDef) Values() []interface{} {
	return ValuesForInterface(*setSelf)
}

// Set Set utils instance
var Set SetForInterfaceDef

// StreamSetForInterface

// StreamSetForInterfaceDef Set inspired by Collection utils
type StreamSetForInterfaceDef struct {
	SetForInterfaceDef
}

// NewStreamSetForInterface New StreamSetForInterface instance
func NewStreamSetForInterface() *StreamSetForInterfaceDef {
	return &StreamSetForInterfaceDef{
		SetForInterfaceDef: SetForInterfaceDef{},
	}
}

// StreamSetForInterfaceFrom New StreamSetForInterface instance from a T array
func StreamSetForInterfaceFrom(list ...interface{}) *StreamSetForInterfaceDef {
	return StreamSetForInterfaceFromArray(list)
}

// StreamSetForInterfaceFromArray New StreamSetForInterface instance from a T array
func StreamSetForInterfaceFromArray(list []interface{}) *StreamSetForInterfaceDef {
	newOne := NewStreamSetForInterface()
	for _, v := range list {
		newOne.SetForInterfaceDef[v] = new(StreamForInterfaceDef)
	}
	return newOne
}

// StreamSetForInterfaceFromMap New StreamSetForInterface instance from a map[T]R
func StreamSetForInterfaceFromMap(theMap map[interface{}]*StreamForInterfaceDef) *StreamSetForInterfaceDef {
	resultMap := make(map[interface{}]interface{}, len(theMap))
	for k, v := range theMap {
		resultMap[k] = v
	}
	result := StreamSetForInterfaceDef{
		SetForInterfaceDef: SetForInterfaceDef(resultMap),
	}
	return &result
}

// StreamSetFromInterface New StreamSetForInterface instance from an array
func StreamSetFromInterface(list ...interface{}) *StreamSetForInterfaceDef {
	return StreamSetForInterfaceFromArray(list)
}

// StreamSetFromArrayInterface New StreamSetForInterface instance from an array
func StreamSetFromArrayInterface(list []interface{}) *StreamSetForInterfaceDef {
	return StreamSetForInterfaceFromArray(list)
}

// Clone Clone this StreamSetForInterface
func (streamSetSelf *StreamSetForInterfaceDef) Clone() *StreamSetForInterfaceDef {
	result := &StreamSetForInterfaceDef{SetForInterfaceDef: SetForInterfaceDef(DuplicateMapForInterface(streamSetSelf.SetForInterfaceDef))}
	for k, v := range result.SetForInterfaceDef {
		if v != nil {
			v = v.(*StreamForInterfaceDef).Clone()
			(result.SetForInterfaceDef)[k] = v
		}
	}

	return result
}

// Union Union an another StreamSetForInterface object
func (streamSetSelf *StreamSetForInterfaceDef) Union(input *StreamSetForInterfaceDef) *StreamSetForInterfaceDef {
	if input == nil || input.Size() == 0 {
		return streamSetSelf
	}

	result := &StreamSetForInterfaceDef{SetForInterfaceDef: SetForInterfaceDef(MergeForInterface(streamSetSelf.SetForInterfaceDef, input.SetForInterfaceDef))}

	for k, v := range streamSetSelf.SetForInterfaceDef {
		v2, ok := (input.SetForInterfaceDef)[k]
		if ok && v2 != nil && v2.(*StreamForInterfaceDef).Len() > 0 {
			if v == nil {
				v = new(StreamForInterfaceDef)
			}
			v = v.(*StreamForInterfaceDef).Extend(v2.(*StreamForInterfaceDef))
			(result.SetForInterfaceDef)[k] = v
		}
	}

	return result
}

// Intersection Get the Intersection with this StreamSetForInterface and an another StreamSetForInterface
func (streamSetSelf *StreamSetForInterfaceDef) Intersection(input *StreamSetForInterfaceDef) *StreamSetForInterfaceDef {
	if input == nil || input.Size() == 0 {
		return NewStreamSetForInterface()
	}

	result := &StreamSetForInterfaceDef{SetForInterfaceDef: SetForInterfaceDef(IntersectionMapByKeyForInterface(streamSetSelf.SetForInterfaceDef, input.SetForInterfaceDef))}

	for k, v := range result.SetForInterfaceDef {
		v2, ok := (input.SetForInterfaceDef)[k]
		if ok && v2 != nil && v2.(*StreamForInterfaceDef).Len() > 0 {
			if v == nil {
				v = new(StreamForInterfaceDef)
			}

			v = v.(*StreamForInterfaceDef).Intersection(v2.(*StreamForInterfaceDef))
			(result.SetForInterfaceDef)[k] = v
		}
	}

	return result
}

// MinusStreams Minus the Stream values by their keys(keys will not be changed but Stream values will)
func (streamSetSelf *StreamSetForInterfaceDef) MinusStreams(input *StreamSetForInterfaceDef) *StreamSetForInterfaceDef {
	if input == nil || input.Size() == 0 {
		return NewStreamSetForInterface()
	}

	result := streamSetSelf.Clone()

	for k, v := range result.SetForInterfaceDef {
		v2, ok := (input.SetForInterfaceDef)[k]
		if ok && v2 != nil && v2.(*StreamForInterfaceDef).Len() > 0 {
			if v == nil {
				v = new(StreamForInterfaceDef)
			}

			v = v.(*StreamForInterfaceDef).Minus(v2.(*StreamForInterfaceDef))
			(result.SetForInterfaceDef)[k] = v
		}
	}

	return result
}

/**
TODO DUPLICATED ZONE (temporarily)
BEGIN
**/

// IsSubsetByKey TODO NOTE !!Duplicated!! returns true or false by checking if set1 is a subset of set2
func (streamSetSelf *StreamSetForInterfaceDef) IsSubsetByKey(input *StreamSetForInterfaceDef) bool {
	if input == nil || input.Size() == 0 {
		return false
	}

	return streamSetSelf.SetForInterfaceDef.IsSubsetByKey(&input.SetForInterfaceDef)
}

// IsSupersetByKey TODO NOTE !!Duplicated!! returns true or false by checking if set1 is a superset of set2
func (streamSetSelf *StreamSetForInterfaceDef) IsSupersetByKey(input *StreamSetForInterfaceDef) bool {
	if input == nil || input.Size() == 0 {
		return true
	}

	return streamSetSelf.SetForInterfaceDef.IsSupersetByKey(&input.SetForInterfaceDef)
}

// Minus TODO NOTE !!Duplicated!! Get all of this StreamSetForInterface but not in the given StreamSetForInterface
func (streamSetSelf *StreamSetForInterfaceDef) Minus(input *StreamSetForInterfaceDef) *StreamSetForInterfaceDef {
	if input == nil || input.Size() == 0 {
		return NewStreamSetForInterface()
	}

	return &StreamSetForInterfaceDef{SetForInterfaceDef: *streamSetSelf.SetForInterfaceDef.Minus(&input.SetForInterfaceDef)}
}

/**
TODO DUPLICATED ZONE (temporarily)
END
**/

// StreamSetForInterface StreamSetForInterface utils instance
var StreamSetForInterface StreamSetForInterfaceDef
