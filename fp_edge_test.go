package fpgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Edge-case behavior of the fp helpers. Several of these document
// surprising-but-intended contracts (e.g. Take(0) returns the whole list)
// so a refactor that "tidies" them up trips a test and forces a conscious
// decision.

// Take/TakeLast treat count <= 0 as "no limit" and return the input slice
// unchanged. This is unusual (many libs return [] for take(0)) but is the
// documented behavior here (count <= 0 short-circuits to `return list`).
func TestTakeNonPositiveReturnsWholeList(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, Take(0, 1, 2, 3))
	assert.Equal(t, []int{1, 2, 3}, Take(-1, 1, 2, 3))
	assert.Equal(t, []int{1, 2, 3}, TakeLast(0, 1, 2, 3))
	assert.Equal(t, []int{1, 2, 3}, TakeLast(-5, 1, 2, 3))
}

// Drop with count <= 0 returns the list unchanged (mirror of Take).
func TestDropNonPositiveReturnsWholeList(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, Drop(0, 1, 2, 3))
	assert.Equal(t, []int{1, 2, 3}, Drop(-2, 1, 2, 3))
}

// SplitEvery: when size evenly divides len there is no trailing partial group.
func TestSplitEveryExactMultiple(t *testing.T) {
	assert.Equal(t, [][]int{{1, 2}, {3, 4}}, SplitEvery(2, 1, 2, 3, 4))
	assert.Equal(t, [][]int{{1, 2, 3}}, SplitEvery(3, 1, 2, 3))
}

// SplitEvery: size >= len yields a single group; size 1 yields singletons.
func TestSplitEverySizeBoundaries(t *testing.T) {
	assert.Equal(t, [][]int{{1, 2, 3}}, SplitEvery(5, 1, 2, 3))
	assert.Equal(t, [][]int{{1}, {2}, {3}}, SplitEvery(1, 1, 2, 3))
}

// SplitEvery: len <= 1 (including empty) OR size <= 0 short-circuits to a
// single group wrapping the whole input as-is (`[][]T{list}`).
func TestSplitEverySingleOrEmpty(t *testing.T) {
	assert.Equal(t, [][]int{{42}}, SplitEvery(3, 42))
	assert.Equal(t, [][]int{{1, 2, 3}}, SplitEvery(0, 1, 2, 3))  // size<=0 short-circuit
	assert.Equal(t, [][]int{{1, 2, 3}}, SplitEvery(-1, 1, 2, 3)) // negative size too
}

// Range with an integer hop that does not evenly divide the span stops before
// crossing the upper bound (exclusive upper).
func TestRangeUnevenHop(t *testing.T) {
	assert.Equal(t, []int{0, 3, 6, 9}, Range(0, 10, 3))
	assert.Equal(t, []int{3, 5}, Range(3, 7, 2))
}

// Range over floats accumulates the hop; the count is deterministic but the
// values carry float rounding, so assert count + InDelta, never exact equality.
func TestRangeFloatHopAccumulation(t *testing.T) {
	r := Range(0.0, 1.0, 0.3)
	assert.Len(t, r, 4) // 0, 0.3, 0.6, ~0.9
	want := []float64{0, 0.3, 0.6, 0.9}
	for i := range want {
		assert.InDelta(t, want[i], r[i], 1e-9)
	}
}

// Range with lower >= higher is empty; a non-positive hop is empty.
func TestRangeEmptyConditions(t *testing.T) {
	assert.Equal(t, []int{}, Range(5, 5))
	assert.Equal(t, []int{}, Range(5, 2))
	assert.Equal(t, []int{}, Range(0, 10, 0))
	assert.Equal(t, []int{}, Range(0, 10, -2))
}

// MinMax uses an `else if` so it must still be correct for strictly ascending
// and strictly descending inputs (the branch can't update both on one item).
func TestMinMaxMonotonic(t *testing.T) {
	mn, mx := MinMax(1, 2, 3, 4, 5)
	assert.Equal(t, 1, mn)
	assert.Equal(t, 5, mx)

	mn, mx = MinMax(5, 4, 3, 2, 1)
	assert.Equal(t, 1, mn)
	assert.Equal(t, 5, mx)

	mn, mx = MinMax(-3, -1, -2)
	assert.Equal(t, -3, mn)
	assert.Equal(t, -1, mx)
}

// Dedupe removes only CONSECUTIVE duplicates; non-adjacent repeats survive.
func TestDedupeOnlyConsecutive(t *testing.T) {
	assert.Equal(t, []int{1, 2, 1, 2}, Dedupe(1, 2, 1, 2))
	assert.Equal(t, []int{1, 2, 3}, Dedupe(1, 1, 2, 2, 3, 3))
}

// Difference/Intersection collapse duplicates within the first list (set
// semantics) while preserving first-seen order.
func TestDifferenceIntersectionDedupeFirstList(t *testing.T) {
	assert.Equal(t, []int{1, 2}, Difference([]int{1, 1, 2, 3}, []int{3}))
	assert.Equal(t, []int{2, 3}, Intersection([]int{1, 2, 2, 3}, []int{2, 3, 3}))
}

// DropWhile that drops everything returns an empty (nil) slice.
func TestDropWhileDropsAll(t *testing.T) {
	r := DropWhile(func(x int) bool { return true }, 1, 2, 3)
	assert.Empty(t, r)
}

// Concat tolerates a nil "mine" head and nil interior slices.
func TestConcatNilHeadAndInterior(t *testing.T) {
	assert.Equal(t, []int{1, 2}, Concat[int](nil, []int{1, 2}))
	assert.Equal(t, []int{1, 2, 3, 4}, Concat([]int{1, 2}, nil, []int{3, 4}))
}
