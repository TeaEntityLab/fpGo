package fpgo

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

var (
	// ErrConversionUnsupported Conversion Unsupported
	ErrConversionUnsupported = errors.New("unsupported")
	// ErrConversionNil Conversion Nil
	ErrConversionNil = errors.New("<nil>")
	// ErrConversionNil Conversion Size Overflow
	ErrConversionSizeOverflow = errors.New("size overflow")
)

// Maybe

// MaybeDef Maybe inspired by Rx/Optional/Guava/Haskell
type MaybeDef[T any] interface {
	Just(in interface{}) MaybeDef[interface{}]
	Or(or T) T
	Clone() MaybeDef[T]
	FlatMap(fn func(T) MaybeDef[T]) MaybeDef[T]
	ToString() string
	ToPtr() *T
	ToMaybe() MaybeDef[T]
	ToFloat64() (float64, error)
	ToFloat32() (float32, error)
	ToInt() (int, error)
	ToInt32() (int32, error)
	ToInt64() (int64, error)
	ToBool() (bool, error)
	Let(fn func())
	Unwrap() T
	UnwrapInterface() interface{}
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
type someDef[T any] struct {
	ref       T
	isNil     bool
	isPresent bool
}

// Just New Maybe by a given value
func (maybeSelf someDef[T]) Just(in interface{}) MaybeDef[interface{}] {
	if IsNil(in) {
		return None
	}

	return JustGenerics(in)
}

// JustGenerics New Maybe by a given value
func JustGenerics[T any](in T) MaybeDef[T] {
	isNil := IsNil(in)
	return someDef[T]{ref: in, isNil: isNil, isPresent: !isNil}
}

// Or Check the value wrapped by Maybe, if it's nil then return a given fallback value
func (maybeSelf someDef[T]) Or(or T) T {
	if maybeSelf.IsNil() {
		return or
	}

	return maybeSelf.ref
}

// CloneTo Clone the Ptr target to an another Ptr target
func CloneTo[T any](maybeSelf MaybeDef[T], dest T) MaybeDef[T] {
	if maybeSelf.IsNil() {
		// return JustGenerics(nil)
		return JustGenerics(maybeSelf.Unwrap())
	}

	x := reflect.ValueOf(maybeSelf.Unwrap())
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
func (maybeSelf someDef[T]) Clone() MaybeDef[T] {
	return CloneTo[T](maybeSelf, *new(T))
}

// FlatMap FlatMap Maybe by function
func (maybeSelf someDef[T]) FlatMap(fn func(T) MaybeDef[T]) MaybeDef[T] {
	return fn(maybeSelf.ref)
}

// ToString Maybe to String
func (maybeSelf someDef[T]) ToString() string {
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
func (maybeSelf someDef[T]) ToPtr() *T {
	if maybeSelf.IsPtr() {
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
func (maybeSelf someDef[T]) ToMaybe() MaybeDef[T] {
	if maybeSelf.IsNil() {
		return maybeSelf
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return maybeSelf
	case someDef[T]:
		return (ref).(someDef[T])
	}
}

// ToFloat64 Maybe to Float64
func (maybeSelf someDef[T]) ToFloat64() (float64, error) {
	if maybeSelf.IsNil() {
		return 0, ErrConversionNil
	}

	var ref interface{} = maybeSelf.ref
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
	case uint:
		val, err := maybeSelf.ToUint()
		return float64(val), err
	case uint16:
		val, err := maybeSelf.ToUint16()
		return float64(val), err
	case uint32:
		val, err := maybeSelf.ToUint32()
		return float64(val), err
	case uint64:
		val, err := maybeSelf.ToUint64()
		return float64(val), err
	case uintptr:
		val, err := maybeSelf.ToUintptr()
		return float64(val), err
	case byte:
		val, err := maybeSelf.ToByte()
		return float64(val), err
	case int:
		val, err := maybeSelf.ToInt()
		return float64(val), err
	case int8:
		val, err := maybeSelf.ToInt8()
		return float64(val), err
	case int16:
		val, err := maybeSelf.ToInt16()
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
func (maybeSelf someDef[T]) ToFloat32() (float32, error) {
	if maybeSelf.IsNil() {
		return 0, ErrConversionNil
	}

	var ref interface{} = maybeSelf.ref
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
	case uint:
		val, err := maybeSelf.ToUint()
		return float32(val), err
	case uint16:
		val, err := maybeSelf.ToUint16()
		return float32(val), err
	case uint32:
		val, err := maybeSelf.ToUint32()
		return float32(val), err
	case uint64:
		val, err := maybeSelf.ToUint64()
		return float32(val), err
	case uintptr:
		val, err := maybeSelf.ToUintptr()
		return float32(val), err
	case byte:
		val, err := maybeSelf.ToByte()
		return float32(val), err
	case int:
		val, err := maybeSelf.ToInt()
		return float32(val), err
	case int8:
		val, err := maybeSelf.ToInt8()
		return float32(val), err
	case int16:
		val, err := maybeSelf.ToInt16()
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
func (maybeSelf someDef[T]) ToInt() (int, error) {
	if maybeSelf.IsNil() {
		return 0, ErrConversionNil
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return 0, ErrConversionUnsupported
	case string:
		return strconv.Atoi(maybeSelf.ToString())
	case bool:
		val, err := maybeSelf.ToBool()
		if val {
			return 1, err
		}
		return 0, err
	case uint:
		val, err := maybeSelf.ToUint()
		if val <= math.MaxInt32 {
			return int(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uint16:
		val, err := maybeSelf.ToUint16()
		return int(val), err
	case uint32:
		val, err := maybeSelf.ToUint32()
		if val <= math.MaxInt32 {
			return int(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uint64:
		val, err := maybeSelf.ToUint64()
		if val <= math.MaxInt32 {
			return int(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uintptr:
		val, err := maybeSelf.ToUintptr()
		if val <= math.MaxInt32 {
			return int(val), err
		}
		return 0, ErrConversionSizeOverflow
	case byte:
		val, err := maybeSelf.ToByte()
		return int(val), err
	case int:
		return (ref).(int), nil
	case int8:
		val, err := maybeSelf.ToInt8()
		return int(val), err
	case int16:
		val, err := maybeSelf.ToInt16()
		return int(val), err
	case int32:
		val, err := maybeSelf.ToInt32()
		if val >= math.MinInt32 && val <= math.MaxInt32 {
			return int(val), err
		}
		return 0, ErrConversionSizeOverflow
	case int64:
		val, err := maybeSelf.ToInt64()
		if val >= math.MinInt32 && val <= math.MaxInt32 {
			return int(val), err
		}
		return 0, ErrConversionSizeOverflow
	case float32:
		val, err := maybeSelf.ToFloat32()
		if val >= math.MinInt32 && val <= math.MaxInt32 {
			return int(math.Round(float64(val))), err
		}
		return 0, ErrConversionSizeOverflow
	case float64:
		val, err := maybeSelf.ToFloat64()
		if val >= math.MinInt32 && val <= math.MaxInt32 {
			return int(math.Round(val)), err
		}
		return 0, ErrConversionSizeOverflow
	}
}

// ToInt8 Maybe to Int8
func (maybeSelf someDef[T]) ToInt8() (int8, error) {
	if maybeSelf.IsNil() {
		return 0, ErrConversionNil
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return int8(0), ErrConversionUnsupported
	case string:
		parseInt, err := strconv.ParseInt((ref).(string), 10, 8)
		return int8(parseInt), err
	case bool:
		val, err := maybeSelf.ToBool()
		if val {
			return int8(1), err
		}
		return int8(0), err
	case uint:
		val, err := maybeSelf.ToUint()
		if val <= math.MaxInt8 {
			return int8(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uint16:
		val, err := maybeSelf.ToUint16()
		if val <= math.MaxInt8 {
			return int8(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uint32:
		val, err := maybeSelf.ToUint32()
		if val <= math.MaxInt8 {
			return int8(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uint64:
		val, err := maybeSelf.ToUint64()
		if val <= math.MaxInt8 {
			return int8(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uintptr:
		val, err := maybeSelf.ToUintptr()
		if val <= math.MaxInt8 {
			return int8(val), err
		}
		return 0, ErrConversionSizeOverflow
	case byte:
		val, err := maybeSelf.ToByte()
		return int8(val), err
	case int:
		val, err := maybeSelf.ToInt()
		if val >= math.MinInt8 && val <= math.MaxInt8 {
			return int8(val), err
		}
		return 0, ErrConversionSizeOverflow
	case int8:
		return (ref).(int8), nil
	case int16:
		val, err := maybeSelf.ToInt16()
		if val >= math.MinInt8 && val <= math.MaxInt8 {
			return int8(val), err
		}
		return 0, ErrConversionSizeOverflow
	case int32:
		val, err := maybeSelf.ToInt32()
		if val >= math.MinInt8 && val <= math.MaxInt8 {
			return int8(val), err
		}
		return 0, ErrConversionSizeOverflow
	case int64:
		val, err := maybeSelf.ToInt64()
		if val >= math.MinInt8 && val <= math.MaxInt8 {
			return int8(val), err
		}
		return 0, ErrConversionSizeOverflow
	case float32:
		val, err := maybeSelf.ToFloat32()
		if val >= math.MinInt8 && val <= math.MaxInt8 {
			return int8(math.Round(float64(val))), err
		}
		return 0, ErrConversionSizeOverflow
	case float64:
		val, err := maybeSelf.ToFloat64()
		if val >= math.MinInt8 && val <= math.MaxInt8 {
			return int8(math.Round(val)), err
		}
		return 0, ErrConversionSizeOverflow
	}
}

// ToInt16 Maybe to Int16
func (maybeSelf someDef[T]) ToInt16() (int16, error) {
	if maybeSelf.IsNil() {
		return 0, ErrConversionNil
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return int16(0), ErrConversionUnsupported
	case string:
		parseInt, err := strconv.ParseInt((ref).(string), 10, 16)
		return int16(parseInt), err
	case bool:
		val, err := maybeSelf.ToBool()
		if val {
			return int16(1), err
		}
		return int16(0), err
	case uint:
		val, err := maybeSelf.ToUint()
		if val <= math.MaxInt16 {
			return int16(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uint16:
		val, err := maybeSelf.ToUint16()
		if val <= math.MaxInt16 {
			return int16(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uint32:
		val, err := maybeSelf.ToUint32()
		if val <= math.MaxInt16 {
			return int16(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uint64:
		val, err := maybeSelf.ToUint64()
		if val <= math.MaxInt16 {
			return int16(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uintptr:
		val, err := maybeSelf.ToUintptr()
		if val <= math.MaxInt16 {
			return int16(val), err
		}
		return 0, ErrConversionSizeOverflow
	case byte:
		val, err := maybeSelf.ToByte()
		return int16(val), err
	case int:
		val, err := maybeSelf.ToInt()
		if val >= math.MinInt16 && val <= math.MaxInt16 {
			return int16(val), err
		}
		return 0, ErrConversionSizeOverflow
	case int8:
		val, err := maybeSelf.ToInt8()
		return int16(val), err
	case int16:
		return (ref).(int16), nil
	case int32:
		val, err := maybeSelf.ToInt32()
		if val >= math.MinInt16 && val <= math.MaxInt16 {
			return int16(val), err
		}
		return 0, ErrConversionSizeOverflow
	case int64:
		val, err := maybeSelf.ToInt64()
		if val >= math.MinInt16 && val <= math.MaxInt16 {
			return int16(val), err
		}
		return 0, ErrConversionSizeOverflow
	case float32:
		val, err := maybeSelf.ToFloat32()
		if val >= math.MinInt16 && val <= math.MaxInt16 {
			return int16(math.Round(float64(val))), err
		}
		return 0, ErrConversionSizeOverflow
	case float64:
		val, err := maybeSelf.ToFloat64()
		if val >= math.MinInt16 && val <= math.MaxInt16 {
			return int16(math.Round(val)), err
		}
		return 0, ErrConversionSizeOverflow
	}
}

// ToInt32 Maybe to Int32
func (maybeSelf someDef[T]) ToInt32() (int32, error) {
	if maybeSelf.IsNil() {
		return 0, ErrConversionNil
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return int32(0), ErrConversionUnsupported
	case string:
		parseInt, err := strconv.ParseInt((ref).(string), 10, 32)
		return int32(parseInt), err
	case bool:
		val, err := maybeSelf.ToBool()
		if val {
			return int32(1), err
		}
		return int32(0), err
	case uint:
		val, err := maybeSelf.ToUint()
		if val <= math.MaxInt32 {
			return int32(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uint16:
		val, err := maybeSelf.ToUint16()
		return int32(val), err
	case uint32:
		val, err := maybeSelf.ToUint32()
		if val <= math.MaxInt32 {
			return int32(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uint64:
		val, err := maybeSelf.ToUint64()
		if val <= math.MaxInt32 {
			return int32(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uintptr:
		val, err := maybeSelf.ToUintptr()
		if val <= math.MaxInt32 {
			return int32(val), err
		}
		return 0, ErrConversionSizeOverflow
	case byte:
		val, err := maybeSelf.ToByte()
		return int32(val), err
	case int:
		val, err := maybeSelf.ToInt()
		if val >= math.MinInt32 && val <= math.MaxInt32 {
			return int32(val), err
		}
		return 0, ErrConversionSizeOverflow
	case int8:
		val, err := maybeSelf.ToInt8()
		return int32(val), err
	case int16:
		val, err := maybeSelf.ToInt16()
		return int32(val), err
	case int32:
		return (ref).(int32), nil
	case int64:
		val, err := maybeSelf.ToInt64()
		if val >= math.MinInt32 && val <= math.MaxInt32 {
			return int32(val), err
		}
		return 0, ErrConversionSizeOverflow
	case float32:
		val, err := maybeSelf.ToFloat32()
		if val >= math.MinInt32 && val <= math.MaxInt32 {
			return int32(math.Round(float64(val))), err
		}
		return 0, ErrConversionSizeOverflow
	case float64:
		val, err := maybeSelf.ToFloat64()
		return int32(math.Round(val)), err
	}
}

// ToInt64 Maybe to Int64
func (maybeSelf someDef[T]) ToInt64() (int64, error) {
	if maybeSelf.IsNil() {
		return 0, ErrConversionNil
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return int64(0), ErrConversionUnsupported
	case string:
		return strconv.ParseInt((ref).(string), 10, 64)
	case bool:
		val, err := maybeSelf.ToBool()
		if val {
			return int64(1), err
		}
		return int64(0), err
	case uint:
		val, err := maybeSelf.ToUint()
		if val <= math.MaxInt64 {
			return int64(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uint16:
		val, err := maybeSelf.ToUint16()
		return int64(val), err
	case uint32:
		val, err := maybeSelf.ToUint32()
		return int64(val), err
	case uint64:
		val, err := maybeSelf.ToUint64()
		if val <= math.MaxInt64 {
			return int64(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uintptr:
		val, err := maybeSelf.ToUintptr()
		if val <= math.MaxInt64 {
			return int64(val), err
		}
		return 0, ErrConversionSizeOverflow
	case byte:
		val, err := maybeSelf.ToByte()
		return int64(val), err
	case int:
		val, err := maybeSelf.ToInt()
		return int64(val), err
	case int8:
		val, err := maybeSelf.ToInt8()
		return int64(val), err
	case int16:
		val, err := maybeSelf.ToInt16()
		return int64(val), err
	case int32:
		val, err := maybeSelf.ToInt32()
		return int64(val), err
	case int64:
		return (ref).(int64), nil
	case float32:
		val, err := maybeSelf.ToFloat32()
		if val >= math.MinInt64 && val <= math.MaxInt64 {
			return int64(math.Round(float64(val))), err
		}
		return 0, ErrConversionSizeOverflow
	case float64:
		val, err := maybeSelf.ToFloat64()
		if val >= math.MinInt64 && val <= math.MaxInt64 {
			return int64(math.Round(val)), err
		}
		return 0, ErrConversionSizeOverflow
	}
}

// ToByte Maybe to Byte
func (maybeSelf someDef[T]) ToByte() (byte, error) {
	if maybeSelf.IsNil() {
		return 0, ErrConversionNil
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return uint8(0), ErrConversionUnsupported
	case string:
		parseInt, err := strconv.ParseInt((ref).(string), 10, 8)
		return uint8(parseInt), err
	case bool:
		val, err := maybeSelf.ToBool()
		if val {
			return uint8(1), err
		}
		return uint8(0), err
	case uint:
		val, err := maybeSelf.ToUint()
		if val >= 0 && val <= math.MaxUint8 {
			return uint8(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uint16:
		val, err := maybeSelf.ToUint16()
		if val >= 0 && val <= math.MaxUint8 {
			return uint8(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uint32:
		val, err := maybeSelf.ToUint32()
		if val >= 0 && val <= math.MaxUint8 {
			return uint8(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uint64:
		val, err := maybeSelf.ToUint64()
		if val >= 0 && val <= math.MaxUint8 {
			return uint8(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uintptr:
		val, err := maybeSelf.ToUintptr()
		if val >= 0 && val <= math.MaxUint8 {
			return uint8(val), err
		}
		return 0, ErrConversionSizeOverflow
	case byte:
		return ref.(byte), nil
	case int:
		val, err := maybeSelf.ToInt()
		if val >= 0 && val <= math.MaxUint8 {
			return uint8(val), err
		}
		return 0, ErrConversionSizeOverflow
	case int8:
		val, err := maybeSelf.ToInt8()
		return uint8(val), err
	case int16:
		val, err := maybeSelf.ToInt16()
		if val >= 0 && val <= math.MaxUint8 {
			return uint8(val), err
		}
		return 0, ErrConversionSizeOverflow
	case int32:
		val, err := maybeSelf.ToInt32()
		if val >= 0 && val <= math.MaxUint8 {
			return uint8(val), err
		}
		return 0, ErrConversionSizeOverflow
	case int64:
		val, err := maybeSelf.ToInt64()
		if val >= 0 && val <= math.MaxUint8 {
			return uint8(val), err
		}
		return 0, ErrConversionSizeOverflow
	case float32:
		val, err := maybeSelf.ToFloat32()
		if val >= 0 && val <= math.MaxUint8 {
			return uint8(math.Round(float64(val))), err
		}
		return 0, ErrConversionSizeOverflow
	case float64:
		val, err := maybeSelf.ToFloat64()
		if val >= 0 && val <= math.MaxUint8 {
			return uint8(math.Round(val)), err
		}
		return 0, ErrConversionSizeOverflow
	}
}

// ToUint Maybe to Uint
func (maybeSelf someDef[T]) ToUint() (uint, error) {
	if maybeSelf.IsNil() {
		return 0, ErrConversionNil
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return 0, ErrConversionUnsupported
	case string:
		parseInt, err := strconv.ParseInt((ref).(string), 10, 32)
		return uint(parseInt), err
	case bool:
		val, err := maybeSelf.ToBool()
		if val {
			return 1, err
		}
		return 0, err
	case uint:
		return ref.(uint), nil
	case uint16:
		val, err := maybeSelf.ToUint16()
		return uint(val), err
	case uint32:
		val, err := maybeSelf.ToUint32()
		if val >= 0 && val <= math.MaxUint32 {
			return uint(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uint64:
		val, err := maybeSelf.ToUint64()
		if val >= 0 && val <= math.MaxUint32 {
			return uint(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uintptr:
		val, err := maybeSelf.ToUintptr()
		if val >= 0 && val <= math.MaxUint32 {
			return uint(val), err
		}
		return 0, ErrConversionSizeOverflow
	case byte:
		val, err := maybeSelf.ToByte()
		return uint(val), err
	case int:
		val, err := maybeSelf.ToInt()
		return uint(val), err
	case int8:
		val, err := maybeSelf.ToInt8()
		return uint(val), err
	case int16:
		val, err := maybeSelf.ToInt16()
		return uint(val), err
	case int32:
		val, err := maybeSelf.ToInt32()
		return uint(val), err
	case int64:
		val, err := maybeSelf.ToInt64()
		if val >= 0 && val <= math.MaxUint32 {
			return uint(val), err
		}
		return 0, ErrConversionSizeOverflow
	case float32:
		val, err := maybeSelf.ToFloat32()
		if val >= 0 && val <= math.MaxUint32 {
			return uint(math.Round(float64(val))), err
		}
		return 0, ErrConversionSizeOverflow
	case float64:
		val, err := maybeSelf.ToFloat64()
		if val >= 0 && val <= math.MaxUint32 {
			return uint(math.Round(val)), err
		}
		return 0, ErrConversionSizeOverflow
	}
}

// ToUint8 Maybe to Uint8
func (maybeSelf someDef[T]) ToUint8() (uint8, error) {
	return maybeSelf.ToByte()
}

// ToUint16 Maybe to Uint16
func (maybeSelf someDef[T]) ToUint16() (uint16, error) {
	if maybeSelf.IsNil() {
		return 0, ErrConversionNil
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return uint16(0), ErrConversionUnsupported
	case string:
		parseInt, err := strconv.ParseInt((ref).(string), 10, 16)
		return uint16(parseInt), err
	case bool:
		val, err := maybeSelf.ToBool()
		if val {
			return uint16(1), err
		}
		return uint16(0), err
	case uint:
		val, err := maybeSelf.ToUint()
		if val >= 0 && val <= math.MaxUint16 {
			return uint16(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uint16:
		return (ref).(uint16), nil
	case uint32:
		val, err := maybeSelf.ToUint32()
		if val >= 0 && val <= math.MaxUint16 {
			return uint16(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uint64:
		val, err := maybeSelf.ToUint64()
		if val >= 0 && val <= math.MaxUint16 {
			return uint16(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uintptr:
		val, err := maybeSelf.ToUintptr()
		if val >= 0 && val <= math.MaxUint16 {
			return uint16(val), err
		}
		return 0, ErrConversionSizeOverflow
	case byte:
		val, err := maybeSelf.ToByte()
		return uint16(val), err
	case int:
		val, err := maybeSelf.ToInt()
		if val >= 0 && val <= math.MaxUint16 {
			return uint16(val), err
		}
		return 0, ErrConversionSizeOverflow
	case int8:
		val, err := maybeSelf.ToInt8()
		return uint16(val), err
	case int16:
		val, err := maybeSelf.ToInt32()
		return uint16(val), err
	case int32:
		val, err := maybeSelf.ToInt32()
		if val >= 0 && val <= math.MaxUint16 {
			return uint16(val), err
		}
		return 0, ErrConversionSizeOverflow
	case int64:
		val, err := maybeSelf.ToInt64()
		if val >= 0 && val <= math.MaxUint16 {
			return uint16(val), err
		}
		return 0, ErrConversionSizeOverflow
	case float32:
		val, err := maybeSelf.ToFloat32()
		if val >= 0 && val <= math.MaxUint16 {
			return uint16(math.Round(float64(val))), err
		}
		return 0, ErrConversionSizeOverflow
	case float64:
		val, err := maybeSelf.ToFloat64()
		if val >= 0 && val <= math.MaxUint16 {
			return uint16(math.Round(val)), err
		}
		return 0, ErrConversionSizeOverflow
	}
}

// ToUint32 Maybe to Uint32
func (maybeSelf someDef[T]) ToUint32() (uint32, error) {
	if maybeSelf.IsNil() {
		return 0, ErrConversionNil
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return uint32(0), ErrConversionUnsupported
	case string:
		parseInt, err := strconv.ParseInt((ref).(string), 10, 32)
		return uint32(parseInt), err
	case bool:
		val, err := maybeSelf.ToBool()
		if val {
			return uint32(1), err
		}
		return uint32(0), err
	case uint:
		val, err := maybeSelf.ToUint()
		if val >= 0 && val <= math.MaxUint32 {
			return uint32(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uint16:
		val, err := maybeSelf.ToUint16()
		return uint32(val), err
	case uint32:
		return ref.(uint32), nil
	case uint64:
		val, err := maybeSelf.ToUint64()
		if val >= 0 && val <= math.MaxUint32 {
			return uint32(val), err
		}
		return 0, ErrConversionSizeOverflow
	case uintptr:
		val, err := maybeSelf.ToUintptr()
		if val >= 0 && val <= math.MaxUint32 {
			return uint32(val), err
		}
		return 0, ErrConversionSizeOverflow
	case byte:
		val, err := maybeSelf.ToByte()
		return uint32(val), err
	case int:
		val, err := maybeSelf.ToInt()
		if val >= 0 && val <= math.MaxUint32 {
			return uint32(val), err
		}
		return 0, ErrConversionSizeOverflow
	case int8:
		val, err := maybeSelf.ToInt8()
		return uint32(val), err
	case int16:
		val, err := maybeSelf.ToInt16()
		return uint32(val), err
	case int32:
		val, err := maybeSelf.ToInt32()
		return uint32(val), err
	case int64:
		val, err := maybeSelf.ToInt64()
		if val >= 0 && val <= math.MaxUint32 {
			return uint32(val), err
		}
		return 0, ErrConversionSizeOverflow
	case float32:
		val, err := maybeSelf.ToFloat32()
		if val >= 0 && val <= math.MaxUint32 {
			return uint32(math.Round(float64(val))), err
		}
		return 0, ErrConversionSizeOverflow
	case float64:
		val, err := maybeSelf.ToFloat64()
		return uint32(math.Round(val)), err
	}
}

// ToUint64 Maybe to Uint64
func (maybeSelf someDef[T]) ToUint64() (uint64, error) {
	if maybeSelf.IsNil() {
		return 0, ErrConversionNil
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return uint64(0), ErrConversionUnsupported
	case string:
		parseInt, err := strconv.ParseInt((ref).(string), 10, 64)
		return uint64(parseInt), err
	case bool:
		val, err := maybeSelf.ToBool()
		if val {
			return uint64(1), err
		}
		return uint64(0), err
	case uint:
		val, err := maybeSelf.ToUint()
		return uint64(val), err
	case uint16:
		val, err := maybeSelf.ToUint16()
		return uint64(val), err
	case uint32:
		val, err := maybeSelf.ToUint32()
		return uint64(val), err
	case uint64:
		return ref.(uint64), nil
	case uintptr:
		val, err := maybeSelf.ToUintptr()
		if val >= 0 && val <= math.MaxUint64 {
			return uint64(val), err
		}
		return 0, ErrConversionSizeOverflow
	case byte:
		val, err := maybeSelf.ToByte()
		return uint64(val), err
	case int:
		val, err := maybeSelf.ToInt()
		return uint64(val), err
	case int8:
		val, err := maybeSelf.ToInt8()
		return uint64(val), err
	case int16:
		val, err := maybeSelf.ToInt16()
		return uint64(val), err
	case int32:
		val, err := maybeSelf.ToInt32()
		return uint64(val), err
	case int64:
		val, err := maybeSelf.ToInt64()
		return uint64(val), err
	case float32:
		val, err := maybeSelf.ToFloat32()
		if val >= 0 && val <= math.MaxUint64 {
			return uint64(math.Round(float64(val))), err
		}
		return 0, ErrConversionSizeOverflow
	case float64:
		val, err := maybeSelf.ToFloat64()
		if val >= 0 && val <= math.MaxUint64 {
			return uint64(math.Round(val)), err
		}
		return 0, ErrConversionSizeOverflow
	}
}

// ToUintptr Maybe to Uintptr
func (maybeSelf someDef[T]) ToUintptr() (uintptr, error) {
	if maybeSelf.IsNil() {
		return 0, ErrConversionNil
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return uintptr(0), ErrConversionUnsupported
	case string:
		parseInt, err := strconv.ParseInt((ref).(string), 10, 64)
		return uintptr(parseInt), err
	case bool:
		val, err := maybeSelf.ToBool()
		if val {
			return uintptr(1), err
		}
		return uintptr(0), err
	case uint:
		val, err := maybeSelf.ToUint()
		return uintptr(val), err
	case uint16:
		val, err := maybeSelf.ToUint16()
		return uintptr(val), err
	case uint32:
		val, err := maybeSelf.ToUint32()
		return uintptr(val), err
	case uint64:
		val, err := maybeSelf.ToUint64()
		return uintptr(val), err
	case uintptr:
		return ref.(uintptr), nil
	case byte:
		val, err := maybeSelf.ToByte()
		return uintptr(val), err
	case int:
		val, err := maybeSelf.ToInt()
		return uintptr(val), err
	case int8:
		val, err := maybeSelf.ToInt8()
		return uintptr(val), err
	case int16:
		val, err := maybeSelf.ToInt16()
		return uintptr(val), err
	case int32:
		val, err := maybeSelf.ToInt32()
		return uintptr(val), err
	case int64:
		val, err := maybeSelf.ToInt64()
		return uintptr(val), err
	case float32:
		val, err := maybeSelf.ToFloat32()
		return uintptr(math.Round(float64(val))), err
	case float64:
		val, err := maybeSelf.ToFloat64()
		return uintptr(math.Round(val)), err
	}
}

// ToBool Maybe to Bool
func (maybeSelf someDef[T]) ToBool() (bool, error) {
	if maybeSelf.IsNil() {
		return false, ErrConversionNil
	}

	var ref interface{} = maybeSelf.ref
	switch (ref).(type) {
	default:
		return false, ErrConversionUnsupported
	case string:
		return strconv.ParseBool(maybeSelf.ToString())
	case bool:
		return (ref).(bool), nil
	case uint:
		val, err := maybeSelf.ToUint()
		return val != 0, err
	case uint16:
		val, err := maybeSelf.ToUint16()
		return val != 0, err
	case uint32:
		val, err := maybeSelf.ToUint32()
		return val != 0, err
	case uint64:
		val, err := maybeSelf.ToUint64()
		return val != 0, err
	case uintptr:
		val, err := maybeSelf.ToUintptr()
		return val != 0, err
	case byte:
		val, err := maybeSelf.ToByte()
		return val != 0, err
	case int:
		val, err := maybeSelf.ToInt()
		return val != 0, err
	case int8:
		val, err := maybeSelf.ToInt8()
		return val != 0, err
	case int16:
		val, err := maybeSelf.ToInt16()
		return val != 0, err
	case int32:
		val, err := maybeSelf.ToInt32()
		return val != 0, err
	case int64:
		val, err := maybeSelf.ToInt64()
		return val != 0, err
	case float32:
		val, err := maybeSelf.ToFloat32()
		return val != 0, err
	case float64:
		val, err := maybeSelf.ToFloat64()
		return val != 0, err
	}
}

// Let If the wrapped value is not nil, then do the given function
func (maybeSelf someDef[T]) Let(fn func()) {
	if maybeSelf.IsPresent() {
		fn()
	}
}

// Unwrap Unwrap the wrapped value of Maybe
func (maybeSelf someDef[T]) Unwrap() T {
	// if maybeSelf.IsNil() {
	// 	return nil
	// }

	return maybeSelf.ref
}

// UnwrapInterface Unwrap the wrapped value of Maybe as interface{}
func (maybeSelf someDef[T]) UnwrapInterface() interface{} {
	if maybeSelf.IsNil() {
		return nil
	}

	return maybeSelf.ref
}

// IsPresent Check is it present(not nil)
func (maybeSelf someDef[T]) IsPresent() bool {
	return maybeSelf.isPresent
	// return !(maybeSelf.IsNil())
}

// IsNil Check is it nil
func (maybeSelf someDef[T]) IsNil() bool {
	return maybeSelf.isNil

	// return IsNil(maybeSelf.ref)
	//
	// val := reflect.ValueOf(maybeSelf.ref)
	//
	// if maybeSelf.Kind() == reflect.Ptr {
	// 	return val.IsNil()
	// }
	// return !val.IsValid()
}

// IsValid Check is its reflect.ValueOf(ref) valid
func (maybeSelf someDef[T]) IsValid() bool {
	val := reflect.ValueOf(maybeSelf.ref)
	return val.IsValid()
}

// IsPtr Check is it a Ptr
func (maybeSelf someDef[T]) IsPtr() bool {
	return IsPtr(maybeSelf.ref)
}

// Type Get its Type
func (maybeSelf someDef[T]) Type() reflect.Type {
	if maybeSelf.IsNil() {
		return reflect.TypeOf(nil)
	}
	return reflect.TypeOf(maybeSelf.ref)
}

// Kind Get its Kind
func (maybeSelf someDef[T]) Kind() reflect.Kind {
	return reflect.ValueOf(maybeSelf.ref).Kind()
}

// IsType Check is its Type equal to the given one
func (maybeSelf someDef[T]) IsType(t reflect.Type) bool {
	return maybeSelf.Type() == t
}

// IsKind Check is its Kind equal to the given one
func (maybeSelf someDef[T]) IsKind(t reflect.Kind) bool {
	return maybeSelf.Kind() == t
}

// Maybe Maybe utils instance
var Maybe someDef[interface{}]

// None

// noneDef None inspired by Rx/Optional/Guava/Haskell
type noneDef struct {
	someDef[interface{}]
}

// Or Check the value wrapped by Maybe, if it's nil then return a given fallback value
func (noneSelf noneDef) Or(or interface{}) interface{} {
	return or
}

// CloneTo Clone the Ptr target to an another Ptr target
func (noneSelf noneDef) CloneTo(dest interface{}) MaybeDef[interface{}] {
	return None
}

// Clone Clone Maybe object & its wrapped value
func (noneSelf noneDef) Clone() MaybeDef[interface{}] {
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
func (noneSelf noneDef) ToMaybe() MaybeDef[interface{}] {
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

// UnwrapInterface Unwrap the wrapped value of Maybe as interface{}
func (noneSelf noneDef) UnwrapInterface() interface{} {
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

// IsPtr Check is it Ptr
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
var None = noneDef{someDef[any]{isNil: true, isPresent: false}}

// // var noneAsSome = someDef[interface{}](None)
// var noneAsSome = None.someDef
