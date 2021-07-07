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
func (streamSelf *StreamDef) Map(fn func(int) interface{}) *StreamDef {

	var list = make([]interface{}, streamSelf.Len())

	for i := range streamSelf.list {
		list[i] = fn(i)
	}

	return &StreamDef{list: list}
}

// Filter Filter items of Stream by function
func (streamSelf *StreamDef) Filter(fn func(int) bool) *StreamDef {

	var list = make([]interface{}, streamSelf.Len())

	var newLen = 0

	for i := range streamSelf.list {
		if fn(i) {
			newLen++
			list[newLen-1] = streamSelf.list[i]
		}
	}

	return &StreamDef{list: list[:newLen]}
}

// Distinct Filter not nil items and return a new Stream instance
func (streamSelf *StreamDef) Distinct() *StreamDef {
	return streamSelf.Filter(func(i int) bool {
		return Just(streamSelf.list[i]).IsPresent()
	})
}

// Append Append an item into Stream
func (streamSelf *StreamDef) Append(item interface{}) *StreamDef {
	streamSelf.list = append(streamSelf.list, item)
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

// Extend Extend Stream by an another Stream
func (streamSelf *StreamDef) Extend(stream *StreamDef) *StreamDef {
	if stream == nil {
		return streamSelf
	}

	var mine = streamSelf.list
	var mineLen = len(mine)
	var target = stream.ToArray()
	var targetLen = len(target)

	var new = make([]interface{}, mineLen+targetLen)
	for i, item := range mine {
		new[i] = item
	}
	for j, item := range target {
		new[mineLen+j] = item
	}
	streamSelf.list = new

	return streamSelf
}

// Sort Sort Stream items by function
func (streamSelf *StreamDef) Sort(fn func(i, j int) bool) *StreamDef {
	sort.Slice(streamSelf.list, fn)
	return streamSelf
}

// Get Get an item of Stream by its index
func (streamSelf *StreamDef) Get(i int) interface{} {
	return streamSelf.list[i]
}

// Stream Stream utils instance
var Stream StreamDef
