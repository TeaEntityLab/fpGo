package fpgo

import (
	"sort"
)

// StreamDef Stream inspired by Collection utils
type StreamDef struct {
	list []interface{}
}

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
	return &StreamDef{list: list}
}

// ToArray Convert Stream to slice
func (streamSelf *StreamDef) ToArray() []interface{} {
	return streamSelf.list
}

// Map Map all items of Stream by function
func (streamSelf *StreamDef) Map(fn func(interface{}, int) interface{}) *StreamDef {

	var result = Stream.FromArray(MapIndexed(fn, (streamSelf.list)...))

	return result
}

// Filter Filter items of Stream by function
func (streamSelf *StreamDef) Filter(fn func(interface{}, int) bool) *StreamDef {

	var result = Stream.FromArray(Filter(fn, (streamSelf.list)...))

	return result
}

// Distinct Filter not nil items and return a new Stream instance
func (streamSelf *StreamDef) Distinct() *StreamDef {
	return streamSelf.Filter(func(val interface{}, i int) bool {
		return Maybe.Just(val).IsPresent()
	})
}

// Append Append an item into Stream
func (streamSelf *StreamDef) Append(item ...interface{}) *StreamDef {
	streamSelf.list = append(streamSelf.list, item...)
	return streamSelf
}

// Remove Remove an item by its index
func (streamSelf *StreamDef) Remove(index int) *StreamDef {
	if index >= 0 && index < streamSelf.Len() {
		streamSelf.list = append(streamSelf.list[:index], streamSelf.list[index+1:]...)
	}
	return streamSelf
}

// Len Get length of Stream
func (streamSelf *StreamDef) Len() int {
	return len(streamSelf.list)
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

	var mine = streamSelf.list
	var mineLen = len(mine)
	var totalLen = mineLen

	for _, stream := range streams {
		if stream == nil {
			continue
		}

		var targetLen = len(stream.list)
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

		var target = stream.list
		var targetLen = len(target)
		for j, item := range target {
			newOne[totalIndex+j] = item
		}
		totalIndex += targetLen
	}

	return Stream.FromArray(newOne)
}

// SortByIndex Sort Stream items by function(index, index) bool
func (streamSelf *StreamDef) SortByIndex(fn func(a, b int) bool) *StreamDef {
	sort.SliceStable(streamSelf.list, fn)
	return streamSelf
}

// Sort Sort Stream items by Comparator
func (streamSelf *StreamDef) Sort(fn Comparator) *StreamDef {
	Sort(fn, streamSelf.list)
	return streamSelf
}

// Get Get an item of Stream by its index
func (streamSelf *StreamDef) Get(i int) interface{} {
	return streamSelf.list[i]
}

// Stream Stream utils instance
var Stream StreamDef
