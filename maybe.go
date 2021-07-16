package fpgo

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

var (
	// ErrConversionUnsupported Conversion Unsupported
	ErrConversionUnsupported = errors.New("unsupported")
	// ErrConversionNil Conversion Nil
	ErrConversionNil = errors.New("<nil>")
)

// Maybe

// MaybeDef Maybe inspired by Rx/Optional/Guava/Haskell
type MaybeDef interface {
	Just(in interface{}) MaybeDef
	Or(or interface{}) interface{}
	CloneTo(dest interface{}) MaybeDef
	Clone() MaybeDef
	FlatMap(fn func(interface{}) MaybeDef) MaybeDef
	ToString() string
	ToPtr() *interface{}
	ToMaybe() MaybeDef
	ToFloat64() (float64, error)
	ToFloat32() (float32, error)
	ToInt() (int, error)
	ToInt32() (int32, error)
	ToInt64() (int64, error)
	ToBool() (bool, error)
	Let(fn func())
	Unwrap() interface{}
	IsPresent() bool
	IsNil() bool
	IsValid() bool
	IsPtr() bool
	Type() reflect.Type
	Kind() reflect.Kind
	IsType(t reflect.Type) bool
	IsKind(t reflect.Kind) bool
}

// someDef Maybe inspired by Rx/Optional/Guava/Haskell
type someDef struct {
	ref       interface{}
	isNil     bool
	isPresent bool
}

// Just New Maybe by a given value
func (maybeSelf someDef) Just(in interface{}) MaybeDef {
	if IsNil(in) {
		return None
	}

	return someDef{ref: in, isNil: false, isPresent: true}
}

// Or Check the value wrapped by Maybe, if it's nil then return a given fallback value
func (maybeSelf someDef) Or(or interface{}) interface{} {
	// if maybeSelf.IsNil() {
	// 	return or
	// }

	return maybeSelf.ref
}

// CloneTo Clone the Ptr target to an another Ptr target
func (maybeSelf someDef) CloneTo(dest interface{}) MaybeDef {
	// if maybeSelf.IsNil() {
	// 	return maybeSelf.Just(nil)
	// }

	x := reflect.ValueOf(maybeSelf.ref)
	if x.Kind() == reflect.Ptr {
		starX := x.Elem()
		y := reflect.New(starX.Type())
		starY := y.Elem()
		starY.Set(starX)
		reflect.ValueOf(dest).Elem().Set(y.Elem())
		return maybeSelf.Just(dest)
	}
	dest = x.Interface()

	return maybeSelf.Just(dest)
}

// Clone Clone Maybe object & its wrapped value
func (maybeSelf someDef) Clone() MaybeDef {
	return maybeSelf.CloneTo(new(interface{}))
}

// FlatMap FlatMap Maybe by function
func (maybeSelf someDef) FlatMap(fn func(interface{}) MaybeDef) MaybeDef {
	return fn(maybeSelf.ref)
}

// ToString Maybe to String
func (maybeSelf someDef) ToString() string {
	ref := maybeSelf.ref
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
func (maybeSelf someDef) ToPtr() *interface{} {
	if maybeSelf.IsPtr() {
		val := reflect.Indirect(reflect.ValueOf(maybeSelf.ref)).Interface()
		return &val
	}

	return &maybeSelf.ref
}

// ToMaybe Maybe to Maybe
func (maybeSelf someDef) ToMaybe() MaybeDef {
	var ref = maybeSelf.ref
	switch (ref).(type) {
	default:
		return maybeSelf
	case someDef:
		return (ref).(someDef)
	}
}

// ToFloat64 Maybe to Float64
func (maybeSelf someDef) ToFloat64() (float64, error) {
	ref := maybeSelf.ref
	switch (ref).(type) {
	default:
		return float64(0), ErrConversionUnsupported
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
func (maybeSelf someDef) ToFloat32() (float32, error) {
	ref := maybeSelf.ref
	switch (ref).(type) {
	default:
		return float32(0), ErrConversionUnsupported
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
func (maybeSelf someDef) ToInt() (int, error) {
	ref := maybeSelf.ref
	switch (ref).(type) {
	default:
		return int(0), ErrConversionUnsupported
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
func (maybeSelf someDef) ToInt32() (int32, error) {
	ref := maybeSelf.ref
	switch (ref).(type) {
	default:
		return int32(0), ErrConversionUnsupported
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
func (maybeSelf someDef) ToInt64() (int64, error) {
	ref := maybeSelf.ref
	switch (ref).(type) {
	default:
		return int64(0), ErrConversionUnsupported
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
func (maybeSelf someDef) ToBool() (bool, error) {
	ref := maybeSelf.ref
	switch (ref).(type) {
	default:
		return bool(false), ErrConversionUnsupported
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
func (maybeSelf someDef) Let(fn func()) {
	// if maybeSelf.IsPresent() {
	// 	fn()
	// }

	fn()
}

// Unwrap Unwrap the wrapped value of Maybe
func (maybeSelf someDef) Unwrap() interface{} {
	// if maybeSelf.IsNil() {
	// 	return nil
	// }

	return maybeSelf.ref
}

// IsPresent Check is it present(not nil)
func (maybeSelf someDef) IsPresent() bool {
	return maybeSelf.isPresent
	// return !(maybeSelf.IsNil())
}

// IsNil Check is it nil
func (maybeSelf someDef) IsNil() bool {
	return maybeSelf.isNil
	// return IsNil(maybeSelf.ref)
}

// IsValid Check is its reflect.ValueOf(ref) valid
func (maybeSelf someDef) IsValid() bool {
	val := reflect.ValueOf(maybeSelf.ref)
	return val.IsValid()
}

// IsPtr Check is it a Ptr
func (maybeSelf someDef) IsPtr() bool {
	return IsPtr(maybeSelf.ref)
}

// Type Get its Type
func (maybeSelf someDef) Type() reflect.Type {
	// if maybeSelf.IsNil() {
	// 	return reflect.TypeOf(nil)
	// }
	return reflect.TypeOf(maybeSelf.ref)
}

// Kind Get its Kind
func (maybeSelf someDef) Kind() reflect.Kind {
	return Kind(maybeSelf.ref)
}

// IsType Check is its Type equal to the given one
func (maybeSelf someDef) IsType(t reflect.Type) bool {
	return maybeSelf.Type() == t
}

// IsKind Check is its Kind equal to the given one
func (maybeSelf someDef) IsKind(t reflect.Kind) bool {
	return maybeSelf.Kind() == t
}

// Maybe Maybe utils instance
var Maybe someDef

// None

// noneDef None inspired by Rx/Optional/Guava/Haskell
type noneDef struct{ someDef }

// Or Check the value wrapped by Maybe, if it's nil then return a given fallback value
func (noneSelf noneDef) Or(or interface{}) interface{} {
	return or
}

// CloneTo Clone the Ptr target to an another Ptr target
func (noneSelf noneDef) CloneTo(dest interface{}) MaybeDef {
	return None
}

// Clone Clone Maybe object & its wrapped value
func (noneSelf noneDef) Clone() MaybeDef {
	return None
}

// ToString Maybe to String
func (noneSelf noneDef) ToString() string {
	return "<nil>"
}

// ToPtr Maybe to Ptr
func (noneSelf noneDef) ToPtr() *interface{} {
	return nil
}

// ToMaybe Maybe to Maybe
func (noneSelf noneDef) ToMaybe() MaybeDef {
	return noneSelf
}

// ToFloat64 Maybe to Float64
func (noneSelf noneDef) ToFloat64() (float64, error) {
	return float64(0), ErrConversionNil
}

// ToFloat32 Maybe to Float32
func (noneSelf noneDef) ToFloat32() (float32, error) {
	return float32(0), ErrConversionNil
}

// ToInt Maybe to Int
func (noneSelf noneDef) ToInt() (int, error) {
	return int(0), ErrConversionNil
}

// ToInt32 Maybe to Int32
func (noneSelf noneDef) ToInt32() (int32, error) {
	return int32(0), ErrConversionNil
}

// ToInt64 Maybe to Int64
func (noneSelf noneDef) ToInt64() (int64, error) {
	return int64(0), ErrConversionNil
}

// ToBool Maybe to Bool
func (noneSelf noneDef) ToBool() (bool, error) {
	return bool(false), ErrConversionNil
}

// Let If the wrapped value is not nil, then do the given function
func (noneSelf noneDef) Let(fn func()) {}

// Unwrap Unwrap the wrapped value of Maybe
func (noneSelf noneDef) Unwrap() interface{} {
	return nil
}

// IsPresent Check is it present(not nil)
func (noneSelf noneDef) IsPresent() bool {
	return false
}

// IsNil Check is it nil
func (noneSelf noneDef) IsNil() bool {
	return true
}

// IsPtr Check is it nil
func (noneSelf noneDef) IsPtr() bool {
	return false
}

// Type Get its Type
func (noneSelf noneDef) Type() reflect.Type {
	return reflect.TypeOf(nil)
}

// Kind Get its Kind
func (noneSelf noneDef) Kind() reflect.Kind {
	return reflect.Invalid
}

// None None utils instance
var None = noneDef{someDef{isNil: true, isPresent: false}}
