// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fpGo

import (
	"fmt"
	"reflect"
	"strconv"
)

type MonadDef struct {
	ref *interface{}
}

func (self MonadDef) JustVal(in interface{}) MonadDef {
	return MonadDef{ref: &in}
}
func (self MonadDef) Just(in *interface{}) MonadDef {
	return MonadDef{ref: in}
}
func (self MonadDef) OrVal(or interface{}) MonadDef {
	return self.Or(&or)
}
func (self MonadDef) Or(or *interface{}) MonadDef {
	if self.IsNil() {
		return MonadDef{ref: or}
	}

	return self
}
func (self MonadDef) ToString() string {
	if self.IsNil() {
		return "<nil>"
	}

	switch (self.Unwrap()).(type) {
	default:
		return fmt.Sprintf("%v", self.Unwrap())
	case int:
		return strconv.Itoa(self.Unwrap().(int))
	case string:
		return self.Unwrap().(string)
	}
}

func (self MonadDef) Let(fn func()) {
	if self.IsPresent() {
		fn()
	}
}

func (self MonadDef) Ref() *interface{} {
	return self.ref
}
func (self MonadDef) Unwrap() interface{} {
	return *self.ref
}
func (self MonadDef) IsPresent() bool {
	return !(self.IsNil())
}
func (self MonadDef) IsNil() bool {
	return self.ref == nil
}

func (self MonadDef) Type() string {
	if self.IsNil() {
		return reflect.TypeOf(self.ref).String()
	}

	return reflect.TypeOf(*self.ref).String()
}

var Monad MonadDef
