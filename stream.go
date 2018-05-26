package fpGo

import (
	"sort"
)

type StreamDef struct {
	list []*interface{}
}

func (self *StreamDef) FromArrayMonad(old []MonadDef) *StreamDef {
	new := make([]*interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = &item
	}

	return self.FromArray(new)
}
func (self *StreamDef) FromArrayString(old []string) *StreamDef {
	new := make([]*interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = &item
	}

	return self.FromArray(new)
}
func (self *StreamDef) FromArrayBool(old []bool) *StreamDef {
	new := make([]*interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = &item
	}

	return self.FromArray(new)
}
func (self *StreamDef) FromArrayInt(old []int) *StreamDef {
	new := make([]*interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = &item
	}

	return self.FromArray(new)
}
func (self *StreamDef) FromArrayInt32(old []int32) *StreamDef {
	new := make([]*interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = &item
	}

	return self.FromArray(new)
}
func (self *StreamDef) FromArrayInt64(old []int64) *StreamDef {
	new := make([]*interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = &item
	}

	return self.FromArray(new)
}
func (self *StreamDef) FromArrayFloat32(old []float32) *StreamDef {
	new := make([]*interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = &item
	}

	return self.FromArray(new)
}
func (self *StreamDef) FromArrayFloat64(old []float64) *StreamDef {
	new := make([]*interface{}, len(old))
	for i, v := range old {
		var item interface{} = v
		new[i] = &item
	}

	return self.FromArray(new)
}
func (self *StreamDef) FromArray(list []*interface{}) *StreamDef {
	return &StreamDef{list: list}
}
func (self *StreamDef) ToArray() []*interface{} {
	return self.list
}

func (self *StreamDef) Map(fn func(int) *interface{}) *StreamDef {

	var list = make([]*interface{}, self.Len())

	for i, _ := range self.list {
		list[i] = fn(i)
	}

	return &StreamDef{list: list}
}
func (self *StreamDef) Filter(fn func(int) bool) *StreamDef {

	var new = &StreamDef{}

	for i, _ := range self.list {
		if fn(i) {
			new = new.Append(self.list[i])
		}
	}

	return new
}
func (self *StreamDef) Append(item *interface{}) *StreamDef {
	self.list = append(self.list, item)
	return self
}
func (self *StreamDef) Remove(index int) *StreamDef {
	if index >= 0 && index < self.Len() {
		self.list = append(self.list[:index], self.list[index+1:]...)
	}
	return self
}
func (self *StreamDef) Len() int {
	return len(self.list)
}
func (self *StreamDef) Extend(stream *StreamDef) *StreamDef {
	if stream == nil {
		return self
	}

	var mine = self.list
	var mineLen = len(mine)
	var target = stream.ToArray()
	var targetLen = len(target)

	var new = make([]*interface{}, mineLen+targetLen)
	for i, item := range mine {
		new[i] = item
	}
	for j, item := range target {
		new[mineLen+j] = item
	}
	self.list = new

	return self
}
func (self *StreamDef) Sort(fn func(i, j int) bool) *StreamDef {
	sort.Slice(self.list, fn)
	return self
}

func (self *StreamDef) Get(i int) *interface{} {
	return self.list[i]
}

var Stream StreamDef
