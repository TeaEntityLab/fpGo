package fpgo

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

var (
	ErrConversionSizeOverflow = errors.New("size overflow")
)

// MaybeDef Maybe inspired by Rx/Optional/Guava/Haskell
type MaybeDef struct {
	ref interface{}
}

// Just New Maybe by a given value
func (maybeSelf MaybeDef) Just(in interface{}) MaybeDef {
	return MaybeDef{ref: in}
}

// Or Check the value wrapped by Maybe, if it's nil then return a given fallback value
func (maybeSelf MaybeDef) Or(or interface{}) interface{} {
	if maybeSelf.IsNil() {
		return or
	}

	return maybeSelf.ref
}

// CloneTo Clone the Ptr target to an another Ptr target
func (maybeSelf MaybeDef) CloneTo(dest interface{}) MaybeDef {
	if maybeSelf.IsNil() {
		return maybeSelf.Just(nil)
	}

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
func (maybeSelf MaybeDef) Clone() MaybeDef {
	return maybeSelf.CloneTo(new(interface{}))
}

// FlatMap FlatMap Maybe by function
func (maybeSelf MaybeDef) FlatMap(fn func(interface{}) *MaybeDef) *MaybeDef {
	return fn(maybeSelf.ref)
}

// ToString Maybe to String
func (maybeSelf MaybeDef) ToString() string {
	if maybeSelf.IsNil() {
		return "<nil>"
	}

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
func (maybeSelf MaybeDef) ToPtr() *interface{} {
	if maybeSelf.Kind() == reflect.Ptr {
		val := reflect.Indirect(reflect.ValueOf(maybeSelf.ref)).Interface()
		return &val
	}

	return &maybeSelf.ref
}

// ToMaybe Maybe to Maybe
func (maybeSelf MaybeDef) ToMaybe() MaybeDef {
	if maybeSelf.IsNil() {
		return maybeSelf
	}

	var ref = maybeSelf.ref
	switch (ref).(type) {
	default:
		return maybeSelf
	case MaybeDef:
		return (ref).(MaybeDef)
	}
}

// ToFloat64 Maybe to Float64
func (maybeSelf MaybeDef) ToFloat64() (float64, error) {
	if maybeSelf.IsNil() {
		return float64(0), errors.New("<nil>")
	}

	ref := maybeSelf.ref
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
func (maybeSelf MaybeDef) ToFloat32() (float32, error) {
	if maybeSelf.IsNil() {
		return float32(0), errors.New("<nil>")
	}

	ref := maybeSelf.ref
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
func (maybeSelf MaybeDef) ToInt() (int, error) {
	if maybeSelf.IsNil() {
		return int(0), errors.New("<nil>")
	}

	ref := maybeSelf.ref
	switch (ref).(type) {
	default:
		return int(0), errors.New("unsupported")
	case string:
		return strconv.Atoi(maybeSelf.ToString())
	case bool:
		val, err := maybeSelf.ToBool()
		if val {
			return 1, err
		}
		return 0, err
	case int:
		return (ref).(int), nil
	case int32:
		val, err := maybeSelf.ToInt32()
		if val > math.MaxInt32 {
			return 0, ErrConversionSizeOverflow
		}
		return int(val), err
	case int64:
		val, err := maybeSelf.ToInt64()
		if val > math.MaxInt32 {
			return 0, ErrConversionSizeOverflow
		}
		return int(val), err
	case float32:
		val, err := maybeSelf.ToFloat32()
		if val > math.MaxInt32 {
			return 0, ErrConversionSizeOverflow
		}
		return int(math.Round(float64(val))), err
	case float64:
		val, err := maybeSelf.ToFloat64()
		if val > math.MaxInt32 {
			return 0, ErrConversionSizeOverflow
		}
		return int(math.Round(val)), err
	}
}

// ToInt32 Maybe to Int32
func (maybeSelf MaybeDef) ToInt32() (int32, error) {
	if maybeSelf.IsNil() {
		return int32(0), errors.New("<nil>")
	}

	ref := maybeSelf.ref
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
		if val > math.MaxInt32 {
			return 0, ErrConversionSizeOverflow
		}
		return int32(val), err
	case float32:
		val, err := maybeSelf.ToFloat32()
		if val > math.MaxInt32 {
			return 0, ErrConversionSizeOverflow
		}
		return int32(math.Round(float64(val))), err
	case float64:
		val, err := maybeSelf.ToFloat64()
		return int32(math.Round(val)), err
	}
}

// ToInt64 Maybe to Int64
func (maybeSelf MaybeDef) ToInt64() (int64, error) {
	if maybeSelf.IsNil() {
		return int64(0), errors.New("<nil>")
	}

	ref := maybeSelf.ref
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
		if val > math.MaxInt64 {
			return 0, ErrConversionSizeOverflow
		}
		return int64(math.Round(float64(val))), err
	case float64:
		val, err := maybeSelf.ToFloat64()
		if val > math.MaxInt64 {
			return 0, ErrConversionSizeOverflow
		}
		return int64(math.Round(val)), err
	}
}

// ToBool Maybe to Bool
func (maybeSelf MaybeDef) ToBool() (bool, error) {
	if maybeSelf.IsNil() {
		return bool(false), errors.New("<nil>")
	}

	ref := maybeSelf.ref
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
func (maybeSelf MaybeDef) Let(fn func()) {
	if maybeSelf.IsPresent() {
		fn()
	}
}

// Unwrap Unwrap the wrapped value of Maybe
func (maybeSelf MaybeDef) Unwrap() interface{} {
	if maybeSelf.IsNil() {
		return nil
	}

	return maybeSelf.ref
}

// IsPresent Check is it present(not nil)
func (maybeSelf MaybeDef) IsPresent() bool {
	return !(maybeSelf.IsNil())
}

// IsNil Check is it nil
func (maybeSelf MaybeDef) IsNil() bool {
	val := reflect.ValueOf(maybeSelf.ref)

	if maybeSelf.Kind() == reflect.Ptr {
		return val.IsNil()
	}
	return !val.IsValid()
}

// IsValid Check is its reflect.ValueOf(ref) valid
func (maybeSelf MaybeDef) IsValid() bool {
	val := reflect.ValueOf(maybeSelf.ref)
	return val.IsValid()
}

// Type Get its Type
func (maybeSelf MaybeDef) Type() reflect.Type {
	if maybeSelf.IsNil() {
		return reflect.TypeOf(nil)
	}
	return reflect.TypeOf(maybeSelf.ref)
}

// Kind Get its Kind
func (maybeSelf MaybeDef) Kind() reflect.Kind {
	return reflect.ValueOf(maybeSelf.ref).Kind()
}

// IsType Check is its Type equal to the given one
func (maybeSelf MaybeDef) IsType(t reflect.Type) bool {
	return maybeSelf.Type() == t
}

// IsKind Check is its Kind equal to the given one
func (maybeSelf MaybeDef) IsKind(t reflect.Kind) bool {
	return maybeSelf.Kind() == t
}

// Maybe Maybe utils instance
var Maybe MaybeDef
