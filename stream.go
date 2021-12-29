package fpgo

import (
	"sort"
)

// Stream

// StreamDef Stream inspired by Collection utils
type StreamDef[T comparable] []T

// StreamFrom New Stream instance from a T array
func StreamFrom[T comparable](list ...T) *StreamDef[T] {
	return StreamFromArray(list)
}

// StreamFromArray New Stream instance from a T array
func StreamFromArray[T comparable](list []T) *StreamDef[T] {
	result := StreamDef[T](list)
	return &result
}

// StreamFromInterface New Stream instance from an array
func StreamFromInterface(list ...interface{}) *StreamDef[interface{}] {
	return StreamFromArray(list)
}

// StreamFromArrayInterface New Stream instance from an array
func StreamFromArrayInterface(list []interface{}) *StreamDef[interface{}] {
	return StreamFromArray(list)
}

// ToArray Convert Stream to slice
func (streamSelf *StreamDef[T]) ToArray() []T {
	return DuplicateSlice(*streamSelf)
}

// Map Map all items of Stream by function
func (streamSelf *StreamDef[T]) Map(fn func(T, int) T) *StreamDef[T] {

	var result = StreamFromArray(MapIndexed(fn, (*streamSelf)...))

	return result
}

// Filter Filter items of Stream by function
func (streamSelf *StreamDef[T]) Filter(fn func(T, int) bool) *StreamDef[T] {

	var result = StreamFromArray(Filter(fn, (*streamSelf)...))

	return result
}

// Reject Reject items of Stream by function
func (streamSelf *StreamDef[T]) Reject(fn func(T, int) bool) *StreamDef[T] {

	var result = StreamFromArray(Reject(fn, (*streamSelf)...))

	return result
}

// FilterNotNil Filter not nil items and return a new Stream instance
func (streamSelf *StreamDef[T]) FilterNotNil() *StreamDef[T] {
	return streamSelf.Filter(func(val T, i int) bool {
		return Maybe.Just(val).IsPresent()
	})
}

// Distinct Filter duplicated items and return a new Stream instance
func (streamSelf *StreamDef[T]) Distinct() *StreamDef[T] {
	result := StreamDef[T](Distinct[T](*streamSelf...))
	return &result
}

// Contains Check the item exists or not in the Stream
func (streamSelf *StreamDef[T]) Contains(input T) bool {
	return Exists(input, *streamSelf...)
}

// IsSubset returns true or false by checking if stream1 is a subset of stream2
func (streamSelf *StreamDef[T]) IsSubset(input *StreamDef[T]) bool {
	if input == nil || input.Len() == 0 {
		return false
	}

	return IsSubset(*streamSelf, *input)
}

// IsSuperset returns true or false by checking if stream1 is a superset of stream2
func (streamSelf *StreamDef[T]) IsSuperset(input *StreamDef[T]) bool {
	if input == nil || input.Len() == 0 {
		return true
	}

	return IsSuperset(*streamSelf, *input)
}

// Clone Clone this Stream
func (streamSelf *StreamDef[T]) Clone() *StreamDef[T] {
	result := StreamDef[T](DuplicateSlice[T](*streamSelf))

	return &result
}

// Intersection Get the Intersection with this Stream and an another Stream
func (streamSelf *StreamDef[T]) Intersection(input *StreamDef[T]) *StreamDef[T] {
	if input == nil || input.Len() == 0 {
		return new(StreamDef[T])
	}

	result := StreamDef[T](Intersection(*streamSelf, *input))

	return &result
}

// Minus Get all of this Stream but not in the given Stream
func (streamSelf *StreamDef[T]) Minus(input *StreamDef[T]) *StreamDef[T] {
	if input == nil || input.Len() == 0 {
		return streamSelf
	}

	result := StreamDef[T](Minus(*streamSelf, *input))
	return &result
}

// RemoveItem Remove items from the Stream
func (streamSelf *StreamDef[T]) RemoveItem(input ...T) *StreamDef[T] {
	inputLen := len(input)
	if (inputLen > 0) {
		result := StreamDef[T](Minus(*streamSelf, input))

		return &result
	}

	return streamSelf
}

// Append Append an item to Stream
func (streamSelf *StreamDef[T]) Append(item ...T) *StreamDef[T] {
	return streamSelf.Concat(item)
}

// Remove Remove an item by its index
func (streamSelf *StreamDef[T]) Remove(index int) *StreamDef[T] {
	var result StreamDef[T]
	if index >= 0 && index < streamSelf.Len() {
		result = append((*streamSelf)[:index], (*streamSelf)[index+1:]...)
	} else {
		return streamSelf
	}
	return &result
}

// Len Get length of Stream
func (streamSelf *StreamDef[T]) Len() int {
	return len(*streamSelf)
}

// Concat Concat Stream by another slices
func (streamSelf *StreamDef[T]) Concat(slices ...[]T) *StreamDef[T] {
	if len(slices) == 0 {
		return streamSelf
	}

	return StreamFromArray(Concat(streamSelf.ToArray(), slices...))
}

// Extend Extend Stream by another Stream(s)
func (streamSelf *StreamDef[T]) Extend(streams ...*StreamDef[T]) *StreamDef[T] {
	if len(streams) == 0 {
		return streamSelf
	}

	var mine = *streamSelf
	var mineLen = len(mine)
	var totalLen = mineLen

	for _, stream := range(streams) {
		if stream == nil {
			continue
		}

		var targetLen = len(*stream)
		totalLen += targetLen
	}
	var newOne = make(StreamDef[T], totalLen)

	for i, item := range mine {
		newOne[i] = item
	}
	totalIndex := mineLen

	for _, stream := range(streams) {
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

	return &newOne
}

// Reverse Reverse Stream items
func (streamSelf *StreamDef[T]) Reverse() *StreamDef[T] {
	result := StreamDef[T](Reverse(*streamSelf...))
	return &result
}

// SortByIndex Sort Stream items by function(index, index) bool
func (streamSelf *StreamDef[T]) SortByIndex(fn func(a, b int) bool) *StreamDef[T] {
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
func (streamSelf *StreamDef[T]) Sort(fn Comparator[T]) *StreamDef[T] {
	result := streamSelf.Clone()
	Sort(fn, *result)
	return result
}

// Get Get an item of Stream by its index
func (streamSelf *StreamDef[T]) Get(i int) T {
	return (*streamSelf)[i]
}

// // Stream Stream utils instance
// var Stream StreamDef[interface{}]

// Set

// SetDef Set inspired by Collection utils
type SetDef[T comparable, R comparable] interface {
	MapKey(fn TransformerFunctor[T, T]) SetDef[T, R]
	MapValue(fn TransformerFunctor[R, R]) SetDef[T, R]
	ContainsKey(input T) bool
	ContainsValue(input R) bool
	IsSubsetByKey(input SetDef[T, R]) bool
	IsSupersetByKey(input SetDef[T, R]) bool
	Add(input ...T) SetDef[T, R]
	RemoveKeys(input ...T) SetDef[T, R]
	RemoveValues(input ...R) SetDef[T, R]
	Get(key T) R
	Set(key T, value R)
	Clone() SetDef[T, R]
	Union(input SetDef[T, R]) SetDef[T, R]
	Intersection(input SetDef[T, R]) SetDef[T, R]
	Minus(input SetDef[T, R]) SetDef[T, R]
	Size() int
	Keys() []T
	Values() []R
    AsMap() map[T] R
    AsMapSet() *MapSetDef[T, R]
}

// MapSetDef Set inspired by Collection utils
type MapSetDef[T comparable, R comparable] map[T] R

// SetFrom New Set instance from a T array
func SetFrom[T comparable, R comparable](list ...T) *MapSetDef[T, R] {
	return SetFromArray[T, R](list)
}

// SetFromArray New Set instance from a T array
func SetFromArray[T comparable, R comparable](list []T) *MapSetDef[T, R] {
	newOne := MapSetDef[T, R](SliceToMap(*new(R), list...))
	return &newOne
}

// SetFromMap New Set instance from a map[T]R
func SetFromMap[T comparable, R comparable](theMap map[T]R) *MapSetDef[T, R] {
	result := MapSetDef[T, R](theMap)
	return &result
}

// SetFromInterface New Set instance from an array
func SetFromInterface(list ...interface{}) *MapSetDef[interface{}, interface{}] {
	return SetFromArray[interface{}, interface{}](list)
}

// SetFromArrayInterface New Set instance from an array
func SetFromArrayInterface(list []interface{}) *MapSetDef[interface{}, interface{}] {
	return SetFromArray[interface{}, interface{}](list)
}

// MapKey Map all keys of Set by function
func (mapSetSelf *MapSetDef[T, R]) MapKey(fn TransformerFunctor[T, T]) *MapSetDef[T, R] {
	result := make(MapSetDef[T, R], len(*mapSetSelf))
	for k, v := range(*mapSetSelf) {
		result[fn(k)] = v
	}

	return &result
}

// MapValue Map all values of Set by function
func (mapSetSelf *MapSetDef[T, R]) MapValue(fn TransformerFunctor[R, R]) *MapSetDef[T, R] {
	result := make(MapSetDef[T, R], len(*mapSetSelf))
	for k, v := range(*mapSetSelf) {
		result[k] = fn(v)
	}

	return &result
}

// ContainsKey Check the key exists or not in the Set
func (mapSetSelf *MapSetDef[T, R]) ContainsKey(input T) bool {
	_, ok := (*mapSetSelf)[input]
	return ok
}

// ContainsValue Check the value exists or not in the Set
func (mapSetSelf *MapSetDef[T, R]) ContainsValue(input R) bool {
	for _, v := range(*mapSetSelf) {
		if v == input {
			return true
		}
	}
	return false
}

// IsSubsetByKey returns true or false by checking if set1 is a subset of set2
func (mapSetSelf *MapSetDef[T, R]) IsSubsetByKey(input *MapSetDef[T, R]) bool {
	return IsSubsetMapByKey(*mapSetSelf, *input)
}

// IsSupersetByKey returns true or false by checking if set1 is a superset of set2
func (mapSetSelf *MapSetDef[T, R]) IsSupersetByKey(input *MapSetDef[T, R]) bool {
	return IsSupersetMapByKey(*mapSetSelf, *input)
}

// Add Add items into the Set
func (mapSetSelf *MapSetDef[T, R]) Add(input ...T) *MapSetDef[T, R] {
	inputLen := len(input)
	if (inputLen > 0) {
		result := mapSetSelf.Clone()
		for _, v := range(input) {
			if _, ok := (*result)[v]; ok {
				continue
			}
			(*result)[v] = *new(R)
		}

		return result
	}

	return mapSetSelf
}

// RemoveKeys Remove keys from the Set
func (mapSetSelf *MapSetDef[T, R]) RemoveKeys(input ...T) *MapSetDef[T, R] {
	inputLen := len(input)
	if (inputLen > 0) {
		result := mapSetSelf.Clone()
		for _, v := range(input) {
			delete(*result, v)
		}

		return result
	}

	return mapSetSelf
}

// RemoveValues Remove values from the Set
func (mapSetSelf *MapSetDef[T, R]) RemoveValues(input ...R) *MapSetDef[T, R] {
	inputLen := len(input)
	if (inputLen > 0) {
		result := mapSetSelf.Clone()
		valueMap := SliceToMap(0, input...)
		for k, v := range(*mapSetSelf) {
			if _, ok := valueMap[v]; ok {
				delete(*result, k)
			}
		}

		return result
	}

	return mapSetSelf
}

// Get Get items from the Set
func (mapSetSelf *MapSetDef[T, R]) Get(key T) R {
	return (*mapSetSelf)[key]
}

// Set Set items to the Set
func (mapSetSelf *MapSetDef[T, R]) Set(key T, value R) {
	(*mapSetSelf)[key] = value

	// return mapSetSelf
}

// Clone Clone this Set
func (mapSetSelf *MapSetDef[T, R]) Clone() *MapSetDef[T, R] {
	result := MapSetDef[T, R](DuplicateMap[T, R](*mapSetSelf))

	return &result
}

// Union Union an another Set object
func (mapSetSelf *MapSetDef[T, R]) Union(input *MapSetDef[T, R]) *MapSetDef[T, R] {
	if (input == nil || input.Size() == 0) {
		return mapSetSelf
	}

	result := MapSetDef[T, R](Merge(*mapSetSelf, *input))

	return &result
}

// Intersection Get the Intersection with this Set and an another Set
func (mapSetSelf *MapSetDef[T, R]) Intersection(input *MapSetDef[T, R]) *MapSetDef[T, R] {
	if (input == nil || input.Size() == 0) {
		return new(MapSetDef[T, R])
	}

	result := MapSetDef[T, R](IntersectionMapByKey(*mapSetSelf, *input))

	return &result
}

// Minus Get all of this Set but not in the given Set
func (mapSetSelf *MapSetDef[T, R]) Minus(input *MapSetDef[T, R]) *MapSetDef[T, R] {
	if input == nil || input.Size() == 0 {
		return mapSetSelf
	}

	result := mapSetSelf.Clone()
	for k := range(*result) {
		_, exists := (*input)[k]
		if exists {
			delete(*result, k)
		}
	}

	return result
}

// Size Get size
func (mapSetSelf *MapSetDef[T, R]) Size() int {
	return len(*mapSetSelf)
}

// Keys Convert Set to slice
func (mapSetSelf *MapSetDef[T, R]) Keys() []T {
	return Keys(*mapSetSelf)
}

// Values Convert Set to slice
func (mapSetSelf *MapSetDef[T, R]) Values() []R {
	return Values(*mapSetSelf)
}

// AsMap Make Set an object typed as map[T] R
func (mapSetSelf *MapSetDef[T, R]) AsMap() map[T] R {
	return *mapSetSelf
}

// AsMapSet Make Set an object typed as *MapSetDef[T, R]
func (mapSetSelf *MapSetDef[T, R]) AsMapSet() *MapSetDef[T, R] {
	return mapSetSelf
}

// // Set Set utils instance
// var Set MapSetDef[interface{}]


// StreamSet

// StreamSetDef StreamSet inspired by Collection utils
type StreamSetDef[T comparable, R comparable] struct {
	MapSetDef[T, *StreamDef[R]]
}

// NewStreamSet New StreamSet instance
func NewStreamSet[T comparable, R comparable]() *StreamSetDef[T, R] {
	return &StreamSetDef[T, R]{
		MapSetDef: MapSetDef[T, *StreamDef[R]]{},
	}
}

// StreamSetFrom New StreamSet instance from a T array
func StreamSetFrom[T comparable, R comparable](list ...T) *StreamSetDef[T, R] {
	return StreamSetFromArray[T, R](list)
}

// StreamSetFromArray New StreamSet instance from a T array
func StreamSetFromArray[T comparable, R comparable](list []T) *StreamSetDef[T, R] {
	newOne := NewStreamSet[T, R]()
  for _, v := range(list) {
    newOne.MapSetDef[v] = new(StreamDef[R])
  }
	return newOne
}

// StreamSetFromMap New StreamSet instance from a map[T]R
func StreamSetFromMap[T comparable, R comparable](theMap map[T]*StreamDef[R]) *StreamSetDef[T, R] {
	result := StreamSetDef[T, R]{MapSetDef[T, *StreamDef[R]](DuplicateMap(theMap))}
	return &result
}

// StreamSetFromInterface New StreamSet instance from an array
func StreamSetFromInterface(list ...interface{}) *StreamSetDef[interface{}, interface{}] {
	return StreamSetFromArray[interface{}, interface{}](list)
}

// StreamSetFromArrayInterface New StreamSet instance from an array
func StreamSetFromArrayInterface(list []interface{}) *StreamSetDef[interface{}, interface{}] {
	return StreamSetFromArray[interface{}, interface{}](list)
}

// Clone Clone this StreamSet
func (streamSetSelf *StreamSetDef[T, R]) Clone() *StreamSetDef[T, R] {
	result := StreamSetFromMap(DuplicateMap(streamSetSelf.MapSetDef))
  for k, v := range result.MapSetDef {
    if (v != nil) {
      v = v.Clone()
      result.MapSetDef[k] = v
    }
  }

	return result
}

// Union Union an another StreamSet object
func (streamSetSelf *StreamSetDef[T, R]) Union(input *StreamSetDef[T, R]) *StreamSetDef[T, R] {
  if (input == nil || input.Size() == 0) {
    return streamSetSelf
  }

	result := StreamSetFromMap(Merge(streamSetSelf.MapSetDef, input.MapSetDef))

  for k, v := range streamSetSelf.MapSetDef {
    v2, ok := input.MapSetDef[k]
    if (ok && v2 != nil && v2.Len() > 0) {
			if (v == nil) {
				v = new(StreamDef[R])
			}
      v = v.Extend(v2)
      result.MapSetDef[k] = v
    }
  }

	return result
}

// Intersection Get the Intersection with this StreamSet and an another StreamSet
func (streamSetSelf *StreamSetDef[T, R]) Intersection(input *StreamSetDef[T, R]) *StreamSetDef[T, R] {
  if (input == nil || input.Size() == 0) {
    return NewStreamSet[T, R]()
  }

	result := StreamSetFromMap(IntersectionMapByKey(streamSetSelf.MapSetDef, input.MapSetDef))

  for k, v := range result.MapSetDef {
    v2, ok := input.MapSetDef[k]
    if (ok && v2 != nil && v2.Len() > 0) {
			if (v == nil) {
				v = new(StreamDef[R])
			}

      v = v.Intersection(v2)
      result.MapSetDef[k] = v
    }
  }

	return result
}

// MinusStreams Minus the Stream values by their keys(keys will not be changed but Stream values will)
func (streamSetSelf *StreamSetDef[T, R]) MinusStreams(input *StreamSetDef[T, R]) *StreamSetDef[T, R] {
  if (input == nil || input.Size() == 0) {
    return NewStreamSet[T, R]()
  }

	result := streamSetSelf.Clone()

  for k, v := range result.MapSetDef {
    v2, ok := input.MapSetDef[k]
    if (ok && v2 != nil && v2.Len() > 0) {
			if (v == nil) {
				v = new(StreamDef[R])
			}

      v = v.Minus(v2)
      result.MapSetDef[k] = v
    }
  }

	return result
}

// // StreamSet StreamSet utils instance
// var StreamSet StreamSetDef[interface{}]
