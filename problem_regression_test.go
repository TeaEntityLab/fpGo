package fpgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// This file collects regression tests that pin down behavior around three
// defects found while auditing fpGo, plus tests documenting genuinely
// surprising-but-intended edge behavior. Each block states the intent it
// guards so a future change that breaks the contract fails loudly.

// ---------------------------------------------------------------------------
// Defect 1 + 2: descriptor-based sorting.
//
// `_compareBySortDescriptors` previously discarded the result of
// `key.CompareTo(...)`, so SimpleSortDescriptor / FieldSortDescriptor /
// SortDescriptorsBuilder sorting on non-nil keys was a no-op (the observed
// order was an artifact of sort.SliceStable seeing an always-"less"
// comparator). Separately, ComparableString used strings.Compare (standard
// sign convention) while the rest of the library (CompareToOrdered,
// ComparableOrdered, SortOrdered) uses the INVERTED convention where CompareTo
// returns +1 when the receiver sorts earlier. Since the descriptor machinery
// has a single comparator for all Comparable types, the lone violator made
// string-keyed sorts run backwards.
//
// These tests use discriminating inputs (already-sorted and shuffled) that a
// degenerate comparator cannot reproduce by coincidence.
// ---------------------------------------------------------------------------

type sortProbe struct {
	Age  int
	Name ComparableString
}

func ageDescriptors(ascending bool) []SortDescriptor[sortProbe] {
	return NewSortDescriptorsBuilder[sortProbe]().
		ThenWithTransformerFunctor(func(o sortProbe) Comparable[interface{}] {
			return NewComparableOrdered(o.Age)
		}, ascending).GetSortDescriptors()
}

func agesOf(items []sortProbe) []int {
	out := make([]int, len(items))
	for i, it := range items {
		out[i] = it.Age
	}
	return out
}

func TestSortDescriptorAscendingActuallySorts(t *testing.T) {
	cases := []struct {
		in   []int
		want []int
	}{
		{[]int{10, 20, 30}, []int{10, 20, 30}}, // already sorted must stay sorted
		{[]int{20, 10, 30}, []int{10, 20, 30}},
		{[]int{30, 20, 10}, []int{10, 20, 30}},
		{[]int{5, 1, 4, 2, 3}, []int{1, 2, 3, 4, 5}},
		{[]int{3, 3, 1, 2, 2}, []int{1, 2, 2, 3, 3}},
	}
	for _, c := range cases {
		items := make([]sortProbe, len(c.in))
		for i, a := range c.in {
			items[i] = sortProbe{Age: a}
		}
		got := SortedListBySortDescriptors(ageDescriptors(true), items...)
		assert.Equal(t, c.want, agesOf(got), "ascending sort of %v", c.in)
	}
}

func TestSortDescriptorDescendingActuallySorts(t *testing.T) {
	cases := []struct {
		in   []int
		want []int
	}{
		{[]int{30, 20, 10}, []int{30, 20, 10}}, // already desc must stay desc
		{[]int{10, 20, 30}, []int{30, 20, 10}},
		{[]int{1, 5, 2, 4, 3}, []int{5, 4, 3, 2, 1}},
	}
	for _, c := range cases {
		items := make([]sortProbe, len(c.in))
		for i, a := range c.in {
			items[i] = sortProbe{Age: a}
		}
		got := SortedListBySortDescriptors(ageDescriptors(false), items...)
		assert.Equal(t, c.want, agesOf(got), "descending sort of %v", c.in)
	}
}

func TestSortDescriptorStringKeyAscending(t *testing.T) {
	// String-keyed ascending sort must order lexicographically (a < b < c).
	in := []sortProbe{
		{Name: NewComparableString("cherry")},
		{Name: NewComparableString("apple")},
		{Name: NewComparableString("banana")},
	}
	descs := NewSortDescriptorsBuilder[sortProbe]().
		ThenWithTransformerFunctor(func(o sortProbe) Comparable[interface{}] {
			return o.Name
		}, true).GetSortDescriptors()
	got := SortedListBySortDescriptors(descs, in...)
	assert.Equal(t, "apple", got[0].Name.Val)
	assert.Equal(t, "banana", got[1].Name.Val)
	assert.Equal(t, "cherry", got[2].Name.Val)
}

func TestSortDescriptorStringKeyDescending(t *testing.T) {
	in := []sortProbe{
		{Name: NewComparableString("apple")},
		{Name: NewComparableString("cherry")},
		{Name: NewComparableString("banana")},
	}
	descs := NewSortDescriptorsBuilder[sortProbe]().
		ThenWithTransformerFunctor(func(o sortProbe) Comparable[interface{}] {
			return o.Name
		}, false).GetSortDescriptors()
	got := SortedListBySortDescriptors(descs, in...)
	assert.Equal(t, "cherry", got[0].Name.Val)
	assert.Equal(t, "banana", got[1].Name.Val)
	assert.Equal(t, "apple", got[2].Name.Val)
}

func TestSortDescriptorMixedKeysTieBreak(t *testing.T) {
	// Primary: Age descending. Secondary tie-break: Name ascending.
	// Both numeric and string descriptors must agree on direction for this to
	// be correct, which only holds once both defects are fixed.
	in := []sortProbe{
		{Age: 30, Name: NewComparableString("BC")},
		{Age: 30, Name: NewComparableString("AD")},
		{Age: 50, Name: NewComparableString("AB")},
		{Age: 30, Name: NewComparableString("AA")},
	}
	descs := NewSortDescriptorsBuilder[sortProbe]().
		ThenWithTransformerFunctor(func(o sortProbe) Comparable[interface{}] {
			return NewComparableOrdered(o.Age)
		}, false).
		ThenWithTransformerFunctor(func(o sortProbe) Comparable[interface{}] {
			return o.Name
		}, true).GetSortDescriptors()
	got := SortedListBySortDescriptors(descs, in...)
	// 50 first, then the three 30s ordered by Name ascending: AA, AD, BC.
	assert.Equal(t, []int{50, 30, 30, 30}, agesOf(got))
	assert.Equal(t, "AB", got[0].Name.Val)
	assert.Equal(t, "AA", got[1].Name.Val)
	assert.Equal(t, "AD", got[2].Name.Val)
	assert.Equal(t, "BC", got[3].Name.Val)
}

// ComparableString must follow the same inverted house convention as
// ComparableOrdered so a single comparator works across Comparable types.
func TestComparableStringHouseConvention(t *testing.T) {
	a := NewComparableString("apple")
	b := NewComparableString("banana")
	// receiver sorts earlier => +1 (NOT strings.Compare's -1)
	assert.Equal(t, 1, a.CompareTo(b))
	assert.Equal(t, -1, b.CompareTo(a))
	assert.Equal(t, 0, a.CompareTo(NewComparableString("apple")))
	// Must agree in sign with ComparableOrdered for the equivalent ordering.
	oa := NewComparableOrdered("apple")
	ob := NewComparableOrdered("banana")
	assert.Equal(t, oa.CompareTo(ob), a.CompareTo(b))
}

// ---------------------------------------------------------------------------
// Defect 3: Compose()/Pipe() with zero functions.
//
// Composing zero functions is the identity function by definition. Previously
// the returned closure indexed fnList[0] and panicked when called with no
// functions. These guard the identity contract.
// ---------------------------------------------------------------------------

func TestComposeNoFunctionsIsIdentity(t *testing.T) {
	id := Compose[int]()
	assert.Equal(t, []int{5, 6, 7}, id(5, 6, 7))
	assert.Empty(t, id()) // no args in => identity returns the (nil) varargs slice
}

func TestPipeNoFunctionsIsIdentity(t *testing.T) {
	id := Pipe[int]()
	assert.Equal(t, []int{5, 6, 7}, id(5, 6, 7))
	assert.Empty(t, id()) // no args in => identity returns the (nil) varargs slice
}

func TestComposeInterfaceNoFunctionsIsIdentity(t *testing.T) {
	id := ComposeInterface()
	assert.Equal(t, []interface{}{1, "x"}, id(1, "x"))
	idp := PipeInterface()
	assert.Equal(t, []interface{}{1, "x"}, idp(1, "x"))
}

// Single-function Compose/Pipe still behaves, and composition order is
// preserved (Compose right-to-left, Pipe left-to-right) — guards that the
// zero-arg guard did not disturb the general path.
func TestComposePipeOrderingUnaffected(t *testing.T) {
	addOne := func(s ...int) []int { return []int{s[0] + 1} }
	timesTwo := func(s ...int) []int { return []int{s[0] * 2} }

	// Compose(addOne, timesTwo)(3) = addOne(timesTwo(3)) = (3*2)+1 = 7
	assert.Equal(t, []int{7}, Compose(addOne, timesTwo)(3))
	// Pipe(addOne, timesTwo)(3) = timesTwo(addOne(3)) = (3+1)*2 = 8
	assert.Equal(t, []int{8}, Pipe(addOne, timesTwo)(3))
}
