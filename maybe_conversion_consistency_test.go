package fpgo

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// These tests pin down Maybe numeric-conversion behavior at the boundaries.
//
// CONTRACT (after the guard fixes): converting an out-of-range value to a
// numeric unsigned type returns ErrConversionSizeOverflow rather than silently
// wrapping. "Out of range" means: a negative signed source, or a value beyond
// the destination's max. This now holds uniformly across ToByte/ToUint/
// ToUint16/ToUint32/ToUint64 regardless of the source width.
//
// EXCEPTION: ToUintptr is an address/bit-pattern type and DELIBERATELY
// reinterprets a negative signed value as its two's-complement bits
// (int64(-1) -> ^uintptr(0)); see TestToUintptrKeepsBitPattern.
//
// Float -> integer out-of-range now also returns an overflow error (the
// float64 branches of ToInt32/ToUint32 gained the bounds checks their float32
// / ToUint64 siblings already had), so those are portable to assert too.

// ---- Negative signed -> numeric unsigned: uniformly rejected --------------

func TestNegativeSignedToUnsignedRejected(t *testing.T) {
	// ToUint across every signed source width.
	for _, tc := range []struct {
		name string
		run  func() (uint, error)
	}{
		{"int(-1)", func() (uint, error) { return Maybe.Just(int(-1)).ToUint() }},
		{"int(-5)", func() (uint, error) { return Maybe.Just(int(-5)).ToUint() }},
		{"int8(-1)", func() (uint, error) { return Maybe.Just(int8(-1)).ToUint() }},
		{"int16(-1)", func() (uint, error) { return Maybe.Just(int16(-1)).ToUint() }},
		{"int32(-1)", func() (uint, error) { return Maybe.Just(int32(-1)).ToUint() }},
		{"int64(-1)", func() (uint, error) { return Maybe.Just(int64(-1)).ToUint() }},
	} {
		got, err := tc.run()
		assert.ErrorIsf(t, err, ErrConversionSizeOverflow, "%s.ToUint", tc.name)
		assert.Equalf(t, uint(0), got, "%s.ToUint value", tc.name)
	}

	// ToUint64 across every signed source width (previously all wrapped).
	for _, tc := range []struct {
		name string
		run  func() (uint64, error)
	}{
		{"int(-1)", func() (uint64, error) { return Maybe.Just(int(-1)).ToUint64() }},
		{"int8(-1)", func() (uint64, error) { return Maybe.Just(int8(-1)).ToUint64() }},
		{"int16(-1)", func() (uint64, error) { return Maybe.Just(int16(-1)).ToUint64() }},
		{"int32(-1)", func() (uint64, error) { return Maybe.Just(int32(-1)).ToUint64() }},
		{"int64(-1)", func() (uint64, error) { return Maybe.Just(int64(-1)).ToUint64() }},
	} {
		got, err := tc.run()
		assert.ErrorIsf(t, err, ErrConversionSizeOverflow, "%s.ToUint64", tc.name)
		assert.Equalf(t, uint64(0), got, "%s.ToUint64 value", tc.name)
	}

	// ToByte / ToUint16 / ToUint32 narrow-signed negatives.
	b, err := Maybe.Just(int8(-1)).ToByte()
	assert.ErrorIs(t, err, ErrConversionSizeOverflow)
	assert.Equal(t, byte(0), b)

	u16, err := Maybe.Just(int16(-1)).ToUint16()
	assert.ErrorIs(t, err, ErrConversionSizeOverflow)
	assert.Equal(t, uint16(0), u16)

	u16b, err := Maybe.Just(int8(-1)).ToUint16()
	assert.ErrorIs(t, err, ErrConversionSizeOverflow)
	assert.Equal(t, uint16(0), u16b)

	u32, err := Maybe.Just(int32(-1)).ToUint32()
	assert.ErrorIs(t, err, ErrConversionSizeOverflow)
	assert.Equal(t, uint32(0), u32)

	u32b, err := Maybe.Just(int8(-1)).ToUint32()
	assert.ErrorIs(t, err, ErrConversionSizeOverflow)
	assert.Equal(t, uint32(0), u32b)

	u32c, err := Maybe.Just(int16(-1)).ToUint32()
	assert.ErrorIs(t, err, ErrConversionSizeOverflow)
	assert.Equal(t, uint32(0), u32c)
}

// Positive in-range signed values still convert cleanly across widths.
func TestPositiveSignedToUnsignedStillWorks(t *testing.T) {
	u, err := Maybe.Just(int8(100)).ToUint()
	assert.NoError(t, err)
	assert.Equal(t, uint(100), u)

	u16, err := Maybe.Just(int16(30000)).ToUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(30000), u16) // guards the former ToInt32 copy-paste

	u64, err := Maybe.Just(int32(1234567)).ToUint64()
	assert.NoError(t, err)
	assert.Equal(t, uint64(1234567), u64)

	b, err := Maybe.Just(int8(127)).ToByte()
	assert.NoError(t, err)
	assert.Equal(t, byte(127), b)
}

// ToUintptr is the deliberate exception: it keeps the two's-complement bits.
func TestToUintptrKeepsBitPattern(t *testing.T) {
	up, err := Maybe.Just(int(-1)).ToUintptr()
	assert.NoError(t, err)
	assert.Equal(t, ^uintptr(0), up)

	up64, err := Maybe.Just(int64(-1)).ToUintptr()
	assert.NoError(t, err)
	assert.Equal(t, ^uintptr(0), up64)
}

// ---- Positive overflow guards (genuine contract, unchanged) ----------------

func TestPositiveOverflowGuardsReturnError(t *testing.T) {
	cases := []struct {
		name string
		run  func() (interface{}, error)
	}{
		{"256.ToByte", func() (interface{}, error) { return Maybe.Just(256).ToByte() }},
		{"uint(300).ToByte", func() (interface{}, error) { return Maybe.Just(uint(300)).ToByte() }},
		{"65536.ToUint16", func() (interface{}, error) { return Maybe.Just(65536).ToUint16() }},
		{"int64(1<<32).ToUint32", func() (interface{}, error) { return Maybe.Just(int64(1 << 32)).ToUint32() }},
		{"int64(MaxInt64).ToInt", func() (interface{}, error) { return Maybe.Just(int64(math.MaxInt64)).ToInt() }},
		{"int64(MaxInt64).ToInt32", func() (interface{}, error) { return Maybe.Just(int64(math.MaxInt64)).ToInt32() }},
		{"int64(MaxInt64).ToInt16", func() (interface{}, error) { return Maybe.Just(int64(math.MaxInt64)).ToInt16() }},
		{"int64(MaxInt64).ToInt8", func() (interface{}, error) { return Maybe.Just(int64(math.MaxInt64)).ToInt8() }},
		{"-129.ToInt8", func() (interface{}, error) { return Maybe.Just(-129).ToInt8() }},
		{"uint64(MaxUint64).ToInt64", func() (interface{}, error) { return Maybe.Just(uint64(math.MaxUint64)).ToInt64() }},
	}
	for _, c := range cases {
		_, err := c.run()
		assert.ErrorIsf(t, err, ErrConversionSizeOverflow, "%s should overflow-error", c.name)
	}
}

// ---- Float -> integer out-of-range: now a portable overflow error ----------

func TestFloatToIntegerOutOfRangeErrors(t *testing.T) {
	// float32 source already had the range check.
	_, err := Maybe.Just(float32(1e18)).ToInt32()
	assert.ErrorIs(t, err, ErrConversionSizeOverflow, "float32(1e18).ToInt32")

	// float64 source now range-checks too (previously implementation-defined).
	got, err := Maybe.Just(float64(1e18)).ToInt32()
	assert.ErrorIs(t, err, ErrConversionSizeOverflow, "float64(1e18).ToInt32")
	assert.Equal(t, int32(0), got)

	// Negative float -> unsigned now guarded on both ToUint32 and ToUint64.
	u32, err := Maybe.Just(float64(-1)).ToUint32()
	assert.ErrorIs(t, err, ErrConversionSizeOverflow, "float64(-1).ToUint32")
	assert.Equal(t, uint32(0), u32)

	u64, err := Maybe.Just(float64(-1)).ToUint64()
	assert.ErrorIs(t, err, ErrConversionSizeOverflow, "float64(-1).ToUint64")
	assert.Equal(t, uint64(0), u64)

	_, err = Maybe.Just(float64(math.MaxFloat64)).ToUint64()
	assert.ErrorIs(t, err, ErrConversionSizeOverflow)
}

// ---- Rounding contract for in-range floats (unchanged) ---------------------

func TestFloatToIntegerRounding(t *testing.T) {
	i, err := Maybe.Just(float64(2.6)).ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 3, i)

	i, err = Maybe.Just(float64(2.4)).ToInt()
	assert.NoError(t, err)
	assert.Equal(t, 2, i)

	i, err = Maybe.Just(float64(-2.6)).ToInt()
	assert.NoError(t, err)
	assert.Equal(t, -3, i)

	b, err := Maybe.Just(float32(3.7)).ToByte()
	assert.NoError(t, err)
	assert.Equal(t, byte(4), b)

	// In-range positive float -> uint32/uint64 still rounds and succeeds.
	u32, err := Maybe.Just(float64(12.6)).ToUint32()
	assert.NoError(t, err)
	assert.Equal(t, uint32(13), u32)
}

// ---- ToUintptr float range-check (bit-pattern rule does NOT apply to floats) -
//
// Integer sources keep their two's-complement bits (see TestToUintptrKeepsBitPattern),
// but float sources are range-checked: float->integer overflow is implementation-
// defined in Go, so out-of-range / negative floats must return an error rather
// than a platform-dependent value.
func TestToUintptrFloatOutOfRangeErrors(t *testing.T) {
	// Out-of-range positive floats overflow-error (previously silently saturated).
	for _, tc := range []struct {
		name string
		run  func() (uintptr, error)
	}{
		{"float64(1e40)", func() (uintptr, error) { return Maybe.Just(float64(1e40)).ToUintptr() }},
		{"float32(1e30)", func() (uintptr, error) { return Maybe.Just(float32(1e30)).ToUintptr() }},
		{"float64(MaxFloat64)", func() (uintptr, error) { return Maybe.Just(float64(math.MaxFloat64)).ToUintptr() }},
		{"float64(-1)", func() (uintptr, error) { return Maybe.Just(float64(-1)).ToUintptr() }},
		{"float32(-2.5)", func() (uintptr, error) { return Maybe.Just(float32(-2.5)).ToUintptr() }},
	} {
		got, err := tc.run()
		assert.ErrorIsf(t, err, ErrConversionSizeOverflow, "%s.ToUintptr", tc.name)
		assert.Equalf(t, uintptr(0), got, "%s.ToUintptr value", tc.name)
	}

	// In-range positive floats still round and succeed (contract preserved).
	up, err := Maybe.Just(float64(2.3)).ToUintptr()
	assert.NoError(t, err)
	assert.Equal(t, uintptr(2), up)

	up, err = Maybe.Just(float32(3.7)).ToUintptr()
	assert.NoError(t, err)
	assert.Equal(t, uintptr(4), up)
}
