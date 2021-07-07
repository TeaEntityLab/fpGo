package fpgo

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// MaybeDef Maybe inspired by Rx/Optional/Guava/Haskell
type MaybeDef[T any] struct {
	ref T
}

// Just New Maybe by a given value
func Just(in interface{}) MaybeDef[interface{}] {
	return JustGenerics(in)
}

// JustGenerics New Maybe by a given value
func JustGenerics[T any](in T) MaybeDef[T] {
	return MaybeDef[T]{ref: in}
}

// Or Check the value wrapped by Maybe, if it's nil then return a given fallback value
func (maybeSelf MaybeDef[T]) Or(or T) T {
	if maybeSelf.IsNil() {
		return or
	}

	return maybeSelf.ref
}

// CloneTo Clone the Ptr target to an another Ptr target
func CloneTo[T any](maybeSelf MaybeDef[T], dest T) MaybeDef[T] {
	if maybeSelf.IsNil() {
		// return JustGenerics(nil)
		return JustGenerics(maybeSelf.ref)
	}

	x := reflect.ValueOf(maybeSelf.ref)
	if x.Kind() == reflect.Ptr {
		starX := x.Elem()
		y := reflect.New(starX.Type())
		starY := y.Elem()
		starY.Set(starX)
		reflect.ValueOf(dest).Elem().Set(y.Elem())
		return JustGenerics(dest)
	}
	dest = x.Interface().(T)

	return JustGenerics(dest)
}

// Clone Clone Maybe object & its wrapped value
func (maybeSelf MaybeDef[T]) Clone() MaybeDef[T] {
	return CloneTo(maybeSelf, *new(T))
}

// FlatMap FlatMap Maybe by function
func (maybeSelf MaybeDef[T]) FlatMap(fn func(T) *MaybeDef[T]) *MaybeDef[T] {
	return fn(maybeSelf.ref)
}

// ToString Maybe to String
func (maybeSelf MaybeDef[T]) ToString() string {
	if maybeSelf.IsNil() {
		return "<nil>"
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return fmt.Sprintf("%v", ref)
	case int:
		return strconv.Itoa((ref).(int))
	case string:
		return (ref).(string)
	}
}

// ToPtr Maybe to Ptr
func (maybeSelf MaybeDef[T]) ToPtr() *T {
	if maybeSelf.Kind() == reflect.Ptr {
		val := reflect.Indirect(reflect.ValueOf(maybeSelf.ref)).Interface()
		switch val.(type) {
		case *T:
			return val.(*T)
		case T:
			result := val.(T)
			return &result
		}
	}

	return &maybeSelf.ref
}

// ToMaybe Maybe to Maybe
func (maybeSelf MaybeDef[T]) ToMaybe() MaybeDef[T] {
	if maybeSelf.IsNil() {
		return maybeSelf
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return maybeSelf
	case MaybeDef[T]:
		return (ref).(MaybeDef[T])
	}
}

// ToFloat64 Maybe to Float64
func (maybeSelf MaybeDef[T]) ToFloat64() (float64, error) {
	if maybeSelf.IsNil() {
		return float64(0), errors.New("<nil>")
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return float64(0), errors.New("unsupported")
	case string:
		return strconv.ParseFloat(maybeSelf.ToString(), 64)
	case bool:
		val, err := maybeSelf.ToBool()
		if val {
			return float64(1), err
		}
		return float64(0), err
	case int:
		val, err := maybeSelf.ToInt()
		return float64(val), err
	case int32:
		val, err := maybeSelf.ToInt32()
		return float64(val), err
	case int64:
		val, err := maybeSelf.ToInt64()
		return float64(val), err
	case float32:
		val, err := maybeSelf.ToFloat32()
		return float64(val), err
	case float64:
		return (ref).(float64), nil
	}
}

// ToFloat32 Maybe to Float32
func (maybeSelf MaybeDef[T]) ToFloat32() (float32, error) {
	if maybeSelf.IsNil() {
		return float32(0), errors.New("<nil>")
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return float32(0), errors.New("unsupported")
	case string:
		val, err := strconv.ParseFloat(maybeSelf.ToString(), 32)
		return float32(val), err
	case bool:
		val, err := maybeSelf.ToBool()
		if val {
			return float32(1), err
		}
		return float32(0), err
	case int:
		val, err := maybeSelf.ToInt()
		return float32(val), err
	case int32:
		val, err := maybeSelf.ToInt32()
		return float32(val), err
	case int64:
		val, err := maybeSelf.ToInt64()
		return float32(val), err
	case float32:
		return (ref).(float32), nil
	case float64:
		val, err := maybeSelf.ToFloat64()
		return float32(val), err
	}
}

// ToInt Maybe to Int
func (maybeSelf MaybeDef[T]) ToInt() (int, error) {
	if maybeSelf.IsNil() {
		return int(0), errors.New("<nil>")
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return int(0), errors.New("unsupported")
	case string:
		return strconv.Atoi(maybeSelf.ToString())
	case bool:
		val, err := maybeSelf.ToBool()
		if val {
			return int(1), err
		}
		return int(0), err
	case int:
		return (ref).(int), nil
	case int32:
		val, err := maybeSelf.ToInt32()
		return int(val), err
	case int64:
		val, err := maybeSelf.ToInt64()
		return int(val), err
	case float32:
		val, err := maybeSelf.ToFloat32()
		return int(val), err
	case float64:
		val, err := maybeSelf.ToFloat64()
		return int(val), err
	}
}

// ToInt32 Maybe to Int32
func (maybeSelf MaybeDef[T]) ToInt32() (int32, error) {
	if maybeSelf.IsNil() {
		return int32(0), errors.New("<nil>")
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return int32(0), errors.New("unsupported")
	case string:
		val, err := maybeSelf.ToInt64()
		return int32(val), err
	case bool:
		val, err := maybeSelf.ToBool()
		if val {
			return int32(1), err
		}
		return int32(0), err
	case int:
		val, err := maybeSelf.ToInt()
		return int32(val), err
	case int32:
		return (ref).(int32), nil
	case int64:
		val, err := maybeSelf.ToInt64()
		return int32(val), err
	case float32:
		val, err := maybeSelf.ToFloat32()
		return int32(val), err
	case float64:
		val, err := maybeSelf.ToFloat64()
		return int32(val), err
	}
}

// ToInt64 Maybe to Int64
func (maybeSelf MaybeDef[T]) ToInt64() (int64, error) {
	if maybeSelf.IsNil() {
		return int64(0), errors.New("<nil>")
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return int64(0), errors.New("unsupported")
	case string:
		return strconv.ParseInt(maybeSelf.ToString(), 10, 32)
	case bool:
		val, err := maybeSelf.ToBool()
		if val {
			return int64(1), err
		}
		return int64(0), err
	case int:
		val, err := maybeSelf.ToInt()
		return int64(val), err
	case int32:
		val, err := maybeSelf.ToInt32()
		return int64(val), err
	case int64:
		return (ref).(int64), nil
	case float32:
		val, err := maybeSelf.ToFloat32()
		return int64(val), err
	case float64:
		val, err := maybeSelf.ToFloat64()
		return int64(val), err
	}
}

// ToBool Maybe to Bool
func (maybeSelf MaybeDef[T]) ToBool() (bool, error) {
	if maybeSelf.IsNil() {
		return bool(false), errors.New("<nil>")
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return bool(false), errors.New("unsupported")
	case string:
		return strconv.ParseBool(maybeSelf.ToString())
	case bool:
		return (ref).(bool), nil
	case int:
		val, err := maybeSelf.ToInt()
		return bool(val != 0), err
	case int32:
		val, err := maybeSelf.ToInt32()
		return bool(val != 0), err
	case int64:
		val, err := maybeSelf.ToInt64()
		return bool(val != 0), err
	case float32:
		val, err := maybeSelf.ToFloat32()
		return bool(val != 0), err
	case float64:
		val, err := maybeSelf.ToFloat64()
		return bool(val != 0), err
	}
}

// Let If the wrapped value is not nil, then do the given function
func (maybeSelf MaybeDef[T]) Let(fn func()) {
	if maybeSelf.IsPresent() {
		fn()
	}
}

// Unwrap Unwrap the wrapped value of Maybe
func (maybeSelf MaybeDef[T]) Unwrap() T {
	// if maybeSelf.IsNil() {
	// 	return nil
	// }

	return maybeSelf.ref
}

// IsPresent Check is it present(not nil)
func (maybeSelf MaybeDef[T]) IsPresent() bool {
	return !(maybeSelf.IsNil())
}

// IsNil Check is it nil
func (maybeSelf MaybeDef[T]) IsNil() bool {
	val := reflect.ValueOf(maybeSelf.ref)

	if maybeSelf.Kind() == reflect.Ptr {
		return val.IsNil()
	}
	return !val.IsValid()
}

// IsValid Check is its reflect.ValueOf(ref) valid
func (maybeSelf MaybeDef[T]) IsValid() bool {
	val := reflect.ValueOf(maybeSelf.ref)
	return val.IsValid()
}

// Type Get its Type
func (maybeSelf MaybeDef[T]) Type() reflect.Type {
	if maybeSelf.IsNil() {
		return reflect.TypeOf(nil)
	}
	return reflect.TypeOf(maybeSelf.ref)
}

// Kind Get its Kind
func (maybeSelf MaybeDef[T]) Kind() reflect.Kind {
	return reflect.ValueOf(maybeSelf.ref).Kind()
}

// IsType Check is its Type equal to the given one
func (maybeSelf MaybeDef[T]) IsType(t reflect.Type) bool {
	return maybeSelf.Type() == t
}

// IsKind Check is its Kind equal to the given one
func (maybeSelf MaybeDef[T]) IsKind(t reflect.Kind) bool {
	return maybeSelf.Kind() == t
}

// Maybe Maybe utils instance
var Maybe MaybeDef[interface{}]
