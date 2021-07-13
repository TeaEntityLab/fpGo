package fpgo

import (
	"sort"
)

// Stream

// StreamDef Stream inspired by Collection utils
type StreamDef []interface{}

// FromArrayMaybe FromArrayMaybe New Stream instance from a Maybe array
func (streamSelf *StreamDef) FromArrayMaybe(old []MaybeDef) *StreamDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayString New Stream instance from a string array
func (streamSelf *StreamDef) FromArrayString(old []string) *StreamDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayBool New Stream instance from a bool array
func (streamSelf *StreamDef) FromArrayBool(old []bool) *StreamDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayInt New Stream instance from an int array
func (streamSelf *StreamDef) FromArrayInt(old []int) *StreamDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayByte New Stream instance from an int8 array
func (streamSelf *StreamDef) FromArrayByte(old []byte) *StreamDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayInt8 New Stream instance from an int8 array
func (streamSelf *StreamDef) FromArrayInt8(old []int8) *StreamDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayInt16 New Stream instance from an int16 array
func (streamSelf *StreamDef) FromArrayInt16(old []int16) *StreamDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayInt32 New Stream instance from an int32 array
func (streamSelf *StreamDef) FromArrayInt32(old []int32) *StreamDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayInt64 New Stream instance from an int64 array
func (streamSelf *StreamDef) FromArrayInt64(old []int64) *StreamDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayFloat32 New Stream instance from a float32 array
func (streamSelf *StreamDef) FromArrayFloat32(old []float32) *StreamDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArrayFloat64 New Stream instance from a float64 array
func (streamSelf *StreamDef) FromArrayFloat64(old []float64) *StreamDef {
	new := make([]interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = item
	}

	return streamSelf.FromArray(new)
}

// FromArray New Stream instance from an interface{} array
func (streamSelf *StreamDef) FromArray(list []interface{}) *StreamDef {
	result := StreamDef(list)
	return &result
}

// ToArray Convert Stream to slice
func (streamSelf *StreamDef) ToArray() []interface{} {
	return DuplicateSlice(*streamSelf)
}

// Map Map all items of Stream by function
func (streamSelf *StreamDef) Map(fn func(interface{}, int) interface{}) *StreamDef {

	var result = Stream.FromArray(MapIndexed(fn, (*streamSelf)...))

	return result
}

// Filter Filter items of Stream by function
func (streamSelf *StreamDef) Filter(fn func(interface{}, int) bool) *StreamDef {

	var result = Stream.FromArray(Filter(fn, (*streamSelf)...))

	return result
}

// Reject Reject items of Stream by function
func (streamSelf *StreamDef) Reject(fn func(interface{}, int) bool) *StreamDef {

	var result = Stream.FromArray(Reject(fn, (*streamSelf)...))

	return result
}

// FilterNotNil Filter not nil items and return a new Stream instance
func (streamSelf *StreamDef) FilterNotNil() *StreamDef {
	return streamSelf.Filter(func(val interface{}, i int) bool {
		return Maybe.Just(val).IsPresent()
	})
}

// Distinct Filter duplicated items and return a new Stream instance
func (streamSelf *StreamDef) Distinct() *StreamDef {
	return Stream.FromArray(Distinct(*streamSelf...))
}

// Contains Check the item exists or not in the Stream
func (streamSelf *StreamDef) Contains(input interface{}) bool {
	return Exists(input, *streamSelf...)
}

// IsSubset returns true or false by checking if stream1 is a subset of stream2
func (streamSelf *StreamDef) IsSubset(input *StreamDef) bool {
	return IsSubset(*streamSelf, *input)
}

// IsSuperset returns true or false by checking if stream1 is a superset of stream2
func (streamSelf *StreamDef) IsSuperset(input *StreamDef) bool {
	return IsSuperset(*streamSelf, *input)
}

// Clone Clone this Stream
func (streamSelf *StreamDef) Clone() *StreamDef {
	return Stream.FromArray(DuplicateSlice(*streamSelf))
}

// Intersection Get the Intersection with this Stream and an another Stream
func (streamSelf *StreamDef) Intersection(input *StreamDef) *StreamDef {
	result := Stream.FromArray(Intersection(*streamSelf, *input))

	return result
}

// Minus Get all of this Stream but not in the given Stream
func (streamSelf *StreamDef) Minus(input *StreamDef) *StreamDef {
	if input.Len() == 0 {
		return streamSelf
	}

	result := Stream.FromArray(Minus(*streamSelf, *input))
	return result
}

// RemoveItem Remove items from the Stream
func (streamSelf *StreamDef) RemoveItem(input ...interface{}) *StreamDef {
	inputLen := len(input)
	if inputLen > 0 {
		result := Stream.FromArray(Minus(*streamSelf, input))

		return result
	}

	return streamSelf
}

// Append Append an item to Stream
func (streamSelf *StreamDef) Append(item ...interface{}) *StreamDef {
	return streamSelf.Concat(item)
}

// Remove Remove an item by its index
func (streamSelf *StreamDef) Remove(index int) *StreamDef {
	if index >= 0 && index < streamSelf.Len() {
		(*streamSelf) = append((*streamSelf)[:index], (*streamSelf)[index+1:]...)
	}
	return streamSelf
}

// Len Get length of Stream
func (streamSelf *StreamDef) Len() int {
	return len(*streamSelf)
}

// Concat Concat Stream by another slices
func (streamSelf *StreamDef) Concat(slices ...[]interface{}) *StreamDef {
	if len(slices) == 0 {
		return streamSelf
	}

	return Stream.FromArray(Concat(streamSelf.ToArray(), slices...))
}

// Extend Extend Stream by another Stream(s)
func (streamSelf *StreamDef) Extend(streams ...*StreamDef) *StreamDef {
	if len(streams) == 0 {
		return streamSelf
	}

	var mine = *streamSelf
	var mineLen = len(mine)
	var totalLen = mineLen

	for _, stream := range streams {
		if stream == nil {
			continue
		}

		var targetLen = len(*stream)
		totalLen += targetLen
	}
	var newOne = make([]interface{}, totalLen)

	for i, item := range mine {
		newOne[i] = item
	}
	totalIndex := mineLen

	for _, stream := range streams {
		if stream == nil {
			continue
		}

		var target = *stream
		var targetLen = len(target)
		for j, item := range target {
			newOne[totalIndex+j] = item
		}
		totalIndex += targetLen
	}

	return Stream.FromArray(newOne)
}

// Reverse Reverse Stream items
func (streamSelf *StreamDef) Reverse() *StreamDef {
	return Stream.FromArray(Reverse(*streamSelf...))
}

// SortByIndex Sort Stream items by function(index, index) bool
func (streamSelf *StreamDef) SortByIndex(fn func(a, b int) bool) *StreamDef {
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
func (streamSelf *StreamDef) Sort(fn Comparator) *StreamDef {
	result := streamSelf.Clone()
	Sort(fn, *result)
	return result
}

// Get Get an item of Stream by its index
func (streamSelf *StreamDef) Get(i int) interface{} {
	return (*streamSelf)[i]
}

// Stream Stream utils instance
var Stream StreamDef

// Set

type SetDef map[interface{}]interface{}

// SetFrom New Set instance from a interface{} array
func SetFrom(list ...interface{}) *SetDef {
	return SetFromArray(list)
}

// SetFromArray New Set instance from a interface{} array
func SetFromArray(list []interface{}) *SetDef {
	newOne := SetDef(SliceToMap(true, list...))
	return &newOne
}

// SetFromMap New Set instance from a map[interface{}]R
func SetFromMap(theMap map[interface{}]interface{}) *SetDef {
	return SetFromArray(Keys(theMap))
}

// SetFromInterface New Set instance from an array
func SetFromInterface(list ...interface{}) *SetDef {
	return SetFromArray(list)
}

// SetFromArrayInterface New Set instance from an array
func SetFromArrayInterface(list []interface{}) *SetDef {
	return SetFromArray(list)
}

// Map Map all items of Set by function
func (setSelf *SetDef) Map(fn TransformerFunctor) *SetDef {
	result := make(SetDef, len(*setSelf))
	for k := range *setSelf {
		result[fn(k)] = true
	}

	return &result
}

// Contains Check the item exists or not in the Set
func (setSelf *SetDef) Contains(input interface{}) bool {
	_, ok := (*setSelf)[input]
	return ok
}

// IsSubset returns true or false by checking if set1 is a subset of set2
func (setSelf *SetDef) IsSubset(input *SetDef) bool {
	return IsSubsetMapByKey(*setSelf, *input)
}

// IsSuperset returns true or false by checking if set1 is a superset of set2
func (setSelf *SetDef) IsSuperset(input *SetDef) bool {
	return IsSupersetMapByKey(*setSelf, *input)
}

// Add Add items into the Set
func (setSelf *SetDef) Add(input ...interface{}) *SetDef {
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

// Remove Remove items from the Set
func (setSelf *SetDef) Remove(input ...interface{}) *SetDef {
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

// Clone Clone this Set
func (setSelf *SetDef) Clone() *SetDef {
	result := SetDef(DuplicateMap(*setSelf))

	return &result
}

// Union Union an another Set object
func (setSelf *SetDef) Union(input *SetDef) *SetDef {
	result := SetDef(Merge(*setSelf, *input))

	return &result
}

// Intersection Get the Intersection with this Set and an another Set
func (setSelf *SetDef) Intersection(input *SetDef) *SetDef {
	result := SetDef(IntersectionMapByKey(*setSelf, *input))

	return &result
}

// Minus Get all of this Set but not in the given Set
func (setSelf *SetDef) Minus(input *SetDef) *SetDef {
	if input.Size() == 0 {
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
func (setSelf *SetDef) Size() int {
	return len(*setSelf)
}

// ToArray Convert Set to slice
func (setSelf *SetDef) ToArray() []interface{} {
	return Keys(*setSelf)
}

// Set Set utils instance
var Set SetDef
