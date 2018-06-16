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
func (self MonadDef) OrVal(or interface{}) *interface{} {
	return self.Or(&or)
}
func (self MonadDef) Or(or *interface{}) *interface{} {
	if self.IsNil() {
		return or
	}

	return self.ref
}

func (self MonadDef) FlatMap(fn func(*interface{}) *MonadDef) *MonadDef {
	return fn(self.ref)
}

func (self MonadDef) ToString() string {
	if self.IsNil() {
		return "<nil>"
	}

	ref := self.ref
	switch (*ref).(type) {
	default:
		return fmt.Sprintf("%v", *ref)
	case int:
		return strconv.Itoa((*ref).(int))
	case string:
		return (*ref).(string)
	}
}
func (self MonadDef) ToMonad() MonadDef {
	if self.IsNil() {
		return self
	}

	var ref = self.ref
	switch (*ref).(type) {
	default:
		return self
	case MonadDef:
		return (*ref).(MonadDef)
	}
}
func (self MonadDef) ToFloat64() (float64, error) {
	if self.IsNil() {
		return float64(0), errors.New("<nil>")
	}

	ref := self.ref
	switch (*ref).(type) {
	default:
		return float64(0), errors.New("unsupported")
	case string:
		return strconv.ParseFloat(self.ToString(), 64)
	case bool:
		val, err := self.ToBool()
		if val {
			return float64(1), err
		} else {
			return float64(0), err
		}
	case int:
		val, err := self.ToInt()
		return float64(val), err
	case int32:
		val, err := self.ToInt32()
		return float64(val), err
	case int64:
		val, err := self.ToInt64()
		return float64(val), err
	case float32:
		val, err := self.ToFloat32()
		return float64(val), err
	case float64:
		return (*ref).(float64), nil
	}
}
func (self MonadDef) ToFloat32() (float32, error) {
	if self.IsNil() {
		return float32(0), errors.New("<nil>")
	}

	ref := self.ref
	switch (*ref).(type) {
	default:
		return float32(0), errors.New("unsupported")
	case string:
		val, err := strconv.ParseFloat(self.ToString(), 32)
		return float32(val), err
	case bool:
		val, err := self.ToBool()
		if val {
			return float32(1), err
		} else {
			return float32(0), err
		}
	case int:
		val, err := self.ToInt()
		return float32(val), err
	case int32:
		val, err := self.ToInt32()
		return float32(val), err
	case int64:
		val, err := self.ToInt64()
		return float32(val), err
	case float32:
		return (*ref).(float32), nil
	case float64:
		val, err := self.ToFloat64()
		return float32(val), err
	}
}
func (self MonadDef) ToInt() (int, error) {
	if self.IsNil() {
		return int(0), errors.New("<nil>")
	}

	ref := self.ref
	switch (*ref).(type) {
	default:
		return int(0), errors.New("unsupported")
	case string:
		return strconv.Atoi(self.ToString())
	case bool:
		val, err := self.ToBool()
		if val {
			return int(1), err
		} else {
			return int(0), err
		}
	case int:
		return (*ref).(int), nil
	case int32:
		val, err := self.ToInt32()
		return int(val), err
	case int64:
		val, err := self.ToInt64()
		return int(val), err
	case float32:
		val, err := self.ToFloat32()
		return int(val), err
	case float64:
		val, err := self.ToFloat64()
		return int(val), err
	}
}
func (self MonadDef) ToInt32() (int32, error) {
	if self.IsNil() {
		return int32(0), errors.New("<nil>")
	}

	ref := self.ref
	switch (*ref).(type) {
	default:
		return int32(0), errors.New("unsupported")
	case string:
		val, err := self.ToInt64()
		return int32(val), err
	case bool:
		val, err := self.ToBool()
		if val {
			return int32(1), err
		} else {
			return int32(0), err
		}
	case int:
		val, err := self.ToInt()
		return int32(val), err
	case int32:
		return (*ref).(int32), nil
	case int64:
		val, err := self.ToInt64()
		return int32(val), err
	case float32:
		val, err := self.ToFloat32()
		return int32(val), err
	case float64:
		val, err := self.ToFloat64()
		return int32(val), err
	}
}
func (self MonadDef) ToInt64() (int64, error) {
	if self.IsNil() {
		return int64(0), errors.New("<nil>")
	}

	ref := self.ref
	switch (*ref).(type) {
	default:
		return int64(0), errors.New("unsupported")
	case string:
		return strconv.ParseInt(self.ToString(), 10, 32)
	case bool:
		val, err := self.ToBool()
		if val {
			return int64(1), err
		} else {
			return int64(0), err
		}
	case int:
		val, err := self.ToInt()
		return int64(val), err
	case int32:
		val, err := self.ToInt32()
		return int64(val), err
	case int64:
		return (*ref).(int64), nil
	case float32:
		val, err := self.ToFloat32()
		return int64(val), err
	case float64:
		val, err := self.ToFloat64()
		return int64(val), err
	}
}
func (self MonadDef) ToBool() (bool, error) {
	if self.IsNil() {
		return bool(false), errors.New("<nil>")
	}

	ref := self.ref
	switch (*ref).(type) {
	default:
		return bool(false), errors.New("unsupported")
	case string:
		return strconv.ParseBool(self.ToString())
	case bool:
		return (*ref).(bool), nil
	case int:
		val, err := self.ToInt()
		return bool(val != 0), err
	case int32:
		val, err := self.ToInt32()
		return bool(val != 0), err
	case int64:
		val, err := self.ToInt64()
		return bool(val != 0), err
	case float32:
		val, err := self.ToFloat32()
		return bool(val != 0), err
	case float64:
		val, err := self.ToFloat64()
		return bool(val != 0), err
	}
}

func (self MonadDef) Let(fn func()) {
	if self.IsPresent() {
		fn()
	}
}

func (self MonadDef) Ref() *interface{} {
	if self.IsNil() {
		return nil
	}

	return self.ref
}
func (self MonadDef) Unwrap() interface{} {
	if self.IsNil() {
		return nil
	}

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
