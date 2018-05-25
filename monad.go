// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fpGo

import (
	"errors"
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

	switch self.Unwrap().(type) {
	default:
		return fmt.Sprintf("%v", self.Unwrap())
	case int:
		return strconv.Itoa(self.Unwrap().(int))
	case string:
		return self.Unwrap().(string)
	}
}
func (self MonadDef) ToFloat64() (float64, error) {
	if self.IsNil() {
		return float64(0), errors.New("<nil>")
	}

	switch self.Unwrap().(type) {
	default:
		return float64(0), errors.New("unsupported")
	case string:
		return strconv.ParseFloat(self.ToString(), 64)
	case int:
		val, err := self.ToInt()
		return float64(val), err
	case float32:
		val, err := self.ToFloat32()
		return float64(val), err
	case float64:
		return self.Unwrap().(float64), nil
	}
}
func (self MonadDef) ToFloat32() (float32, error) {
	if self.IsNil() {
		return float32(0), errors.New("<nil>")
	}

	switch self.Unwrap().(type) {
	default:
		return float32(0), errors.New("unsupported")
	case string:
		val, err := strconv.ParseFloat(self.ToString(), 32)
		return float32(val), err
	case int:
		val, err := self.ToInt()
		return float32(val), err
	case float32:
		return self.Unwrap().(float32), nil
	case float64:
		val, err := self.ToFloat64()
		return float32(val), err
	}
}
func (self MonadDef) ToInt() (int, error) {
	if self.IsNil() {
		return int(0), errors.New("<nil>")
	}

	switch self.Unwrap().(type) {
	default:
		return int(0), errors.New("unsupported")
	case string:
		return strconv.Atoi(self.ToString())
	case int:
		return self.Unwrap().(int), nil
	case float32:
		val, err := self.ToFloat32()
		return int(val), err
	case float64:
		val, err := self.ToFloat64()
		return int(val), err
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

func (self MonadDef) Type() reflect.Type {
	if self.IsNil() {
		return reflect.TypeOf(nil)
	}

	return reflect.TypeOf(self.Unwrap())
}
func (self MonadDef) Kind() reflect.Kind {
	if self.IsNil() {
		return reflect.TypeOf(self.ref).Kind()
	}

	return self.Type().Kind()
}
func (self MonadDef) IsType(t reflect.Type) bool {
	return self.Type() == t
}
func (self MonadDef) IsKind(t reflect.Kind) bool {
	return self.Kind() == t
}

var Monad MonadDef
