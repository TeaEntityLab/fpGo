package fpgo

import (
	"sort"
)

// StreamDef Stream inspired by Collection utils
type StreamDef[T any] []T

// StreamFrom New Stream instance from a T array
func StreamFrom[T any](list ...T) *StreamDef[T] {
	return StreamFromArray(list)
}

// StreamFromArray New Stream instance from a T array
func StreamFromArray[T any](list []T) *StreamDef[T] {
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
	return []T(*streamSelf)
}

// Map Map all items of Stream by function
func (streamSelf *StreamDef[T]) Map(fn func(int) T) *StreamDef[T] {

	var list = make(StreamDef[T], streamSelf.Len())

	for i := range *streamSelf {
		list[i] = fn(i)
	}

	return &list
}

// Filter Filter items of Stream by function
func (streamSelf *StreamDef[T]) Filter(fn func(int) bool) *StreamDef[T] {

	var list = make(StreamDef[T], streamSelf.Len())

	var newLen = 0

	for i := range *streamSelf {
		if fn(i) {
			newLen++
			list[newLen-1] = (*streamSelf)[i]
		}
	}

	result := list[:newLen]
	return &result
}

// Distinct Filter not nil items and return a new Stream instance
func (streamSelf *StreamDef[T]) Distinct() *StreamDef[T] {
	return streamSelf.Filter(func(i int) bool {
		return Just((*streamSelf)[i]).IsPresent()
	})
}

// Append Append an item into Stream
func (streamSelf *StreamDef[T]) Append(item T) *StreamDef[T] {
	result := StreamDef[T](append(*streamSelf, item))
	return &result
}

// Remove Remove an item by its index
func (streamSelf *StreamDef[T]) Remove(index int) *StreamDef[T] {
	var result StreamDef[T]
	if index >= 0 && index < streamSelf.Len() {
		result = append((*streamSelf)[:index], (*streamSelf)[index+1:]...)
	} else {
		result = *streamSelf
	}
	return &result
}

// Len Get length of Stream
func (streamSelf *StreamDef[T]) Len() int {
	return len(*streamSelf)
}

// Extend Extend Stream by an another Stream
func (streamSelf *StreamDef[T]) Extend(stream *StreamDef[T]) *StreamDef[T] {
	if stream == nil {
		return streamSelf
	}

	var mine = *streamSelf
	var mineLen = len(mine)
	var target = stream.ToArray()
	var targetLen = len(target)

	var new = make(StreamDef[T], mineLen+targetLen)
	for i, item := range mine {
		new[i] = item
	}
	for j, item := range target {
		new[mineLen+j] = item
	}

	return &new
}

// Sort Sort Stream items by function
func (streamSelf *StreamDef[T]) Sort(fn func(i, j int) bool) *StreamDef[T] {
	sort.Slice(streamSelf, fn)
	return streamSelf
}

// Get Get an item of Stream by its index
func (streamSelf *StreamDef[T]) Get(i int) T {
	return (*streamSelf)[i]
}

// Stream Stream utils instance
// var Stream StreamDef[interface{}]
