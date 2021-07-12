package fpgo

import (
	"reflect"
	"strings"
)

// NewComparableOrdered Generate a Ordered Comparable for Comparator
func NewComparableOrdered[T Ordered](val T) ComparableOrdered[T] {
	return ComparableOrdered[T]{
		Val: val,
	}
}
// ComparableOrdered A Ordered Comparable for Comparator
type ComparableOrdered[T Ordered] struct {
	Val T
}
// CompareTo Compare with an another object
func (obj ComparableOrdered[T]) CompareTo(input interface{}) int {
	return CompareToOrdered(obj.Val, input.(ComparableOrdered[T]).Val)
}

// NewComparableString Generate a String Comparable for Comparator
func NewComparableString(val string) ComparableString {
	return ComparableString{
		Val: val,
	}
}
// ComparableString A String Comparable for Comparator
type ComparableString struct {
	Val string
}
// CompareTo Compare with an another object
func (obj ComparableString) CompareTo(input interface{}) int {
	return strings.Compare(string(obj.Val), string(input.(ComparableString).Val))
}

// SortDescriptor Define a Transformer Pattern SortDescriptor
type SortDescriptor[T any] interface {
	// Transformer[T, Comparable] // TODO NOTE go2go still can't use this. ref( https://groups.google.com/g/golang-nuts/c/gX5KC-DF_dQ )
	TransformedBy() TransformerFunctor[T, Comparable[interface{}]] // Use this for go2go temporarily

	IsAscending() bool
	SetAscending(bool)
}

// SortedListBySortDescriptors Sort items by sortDescriptors and return value
func SortedListBySortDescriptors[T any](sortDescriptors []SortDescriptor[T], input ...T) []T {
	result := append(input[:0:0], input...)
	SortBySortDescriptors(sortDescriptors, result)

	return result
}

// SortBySortDescriptors Sort items by sortDescriptors
func SortBySortDescriptors[T any](sortDescriptors []SortDescriptor[T], input []T) {
	Sort(func (item1 T, item2 T) bool {
		return _compareBySortDescriptors(item1, item2, sortDescriptors, 0) >= 0
	}, input)
}

func _compareBySortDescriptors[T any](item1 T, item2 T, sortDescriptors []SortDescriptor[T], descriptorIndex int) int {
	descriptor := sortDescriptors[descriptorIndex]
	key1 := descriptor.TransformedBy()(item1)
	key2 := descriptor.TransformedBy()(item2)
	result := 0
	if key1 != nil && key2 != nil {
		if descriptor.IsAscending() {
			key1.CompareTo(key2)
		} else {
			key2.CompareTo(key1)
		}
	}
	if key1 != nil && key2 == nil {
		if descriptor.IsAscending() {
			return 1
		}

		return -1
	}
	if key1 == nil && key2 != nil {
		if descriptor.IsAscending() {
			return -1
		}

		return 1
	}
	if (result == 0 && _hasNextDescriptor(sortDescriptors, descriptorIndex)) {
		return _compareBySortDescriptors(item1, item2, sortDescriptors, descriptorIndex + 1)
	}

	return result
}

func _hasNextDescriptor[T any](sortDescriptors []SortDescriptor[T], index int) bool {
	return (index + 1) < len(sortDescriptors)
}

// SimpleSortDescriptor

// NewSimpleSortDescriptor Generate a new SimpleSortDescriptor by TransformerFunctor & ascending(true)/descending(false)
func NewSimpleSortDescriptor[T any](transformFn TransformerFunctor[T, Comparable[interface{}]], ascending bool) SimpleSortDescriptor[T] {
	return SimpleSortDescriptor[T]{
		transformFn: transformFn,
		ascending: ascending,
	}
}

// SimpleSortDescriptor SimpleSortDescriptor implemented by TransformerFunctor
type SimpleSortDescriptor[T any] struct {
	ascending bool

	transformFn TransformerFunctor[T, Comparable[interface{}]]
}

// IsAscending Check is this SortDescriptor sorting by ascending
func (descriptor SimpleSortDescriptor[T]) IsAscending() bool {
	return descriptor.ascending
}

// SetAscending Set this SortDescriptor sorting by ascending(true) or descending(false)
func (descriptor SimpleSortDescriptor[T]) SetAscending(val bool) {
	descriptor.ascending = val
}

// TransformedBy Get the TransformerFunctor of this SortDescriptor
func (descriptor SimpleSortDescriptor[T]) TransformedBy() TransformerFunctor[T, Comparable[interface{}]] {
	return descriptor.transformFn
}

// FieldSortDescriptor

// NewFieldSortDescriptor Generate a new FieldSortDescriptor by FieldName & ascending(true)/descending(false)
func NewFieldSortDescriptor[T any](fieldName string, ascending bool) FieldSortDescriptor[T] {
	return FieldSortDescriptor[T]{
		// SimpleSortDescriptor: SimpleSortDescriptor[T]{
		// 	ascending: ascending,
		// },
		ascending: ascending,

		fieldName: fieldName,
	}
}

// FieldSortDescriptor FieldSortDescriptor implemented by Reflection(by FieldName)
type FieldSortDescriptor[T any] struct {
	// SimpleSortDescriptor[T] // TODO NOTE go2go still can't use this. ref( https://groups.google.com/g/golang-nuts/c/gX5KC-DF_dQ )
	ascending bool

	fieldName string
}

// IsAscending Check is this SortDescriptor sorting by ascending
// TODO NOTE Deprecated
func (descriptor FieldSortDescriptor[T]) IsAscending() bool {
	return descriptor.ascending
}

// SetAscending Set this SortDescriptor sorting by ascending(true) or descending(false)
// TODO NOTE Deprecated
func (descriptor FieldSortDescriptor[T]) SetAscending(val bool) {
	descriptor.ascending = val
}

// GetFieldName Get the fieldName to sort
func (descriptor FieldSortDescriptor[T]) GetFieldName() string {
	return descriptor.fieldName
}

// SetFieldName Set the fieldName to sort
func (descriptor FieldSortDescriptor[T]) SetFieldName(val string) {
	descriptor.fieldName = val
}

// TransformedBy Get the TransformerFunctor of this SortDescriptor
func (descriptor FieldSortDescriptor[T]) TransformedBy() TransformerFunctor[T, Comparable[interface{}]] {
	return func (input T) Comparable[interface{}] {
		r := reflect.ValueOf(input)
		f := reflect.Indirect(r).FieldByName(descriptor.fieldName).Interface()

		return f.(Comparable[interface{}])
	}
}

// SortDescriptorsBuilder

// NewSortDescriptorsBuilder Generate a new SortDescriptorsBuilder
func NewSortDescriptorsBuilder[T any]() SortDescriptorsBuilder[T] {
	return SortDescriptorsBuilder[T]{}
}

// SortDescriptorsBuilder SortDescriptorsBuilder for composing SortDescriptor list and sorting data
type SortDescriptorsBuilder[T any] []SortDescriptor[T]

// ThenWithTransformerFunctor Use TransformerFunctor as a SortDescriptor
func (builder SortDescriptorsBuilder[T]) ThenWithTransformerFunctor(transformFn TransformerFunctor[T, Comparable[interface{}]], ascending bool) SortDescriptorsBuilder[T] {
	result := append(builder, NewSimpleSortDescriptor(transformFn, ascending))
	return result
}

// ThenWithFieldName Use FieldName as a SortDescriptor
func (builder SortDescriptorsBuilder[T]) ThenWithFieldName(fieldName string, ascending bool) SortDescriptorsBuilder[T] {
	result := append(builder, NewFieldSortDescriptor[T](fieldName, ascending))
	return result
}

// ThenWith Append a SortDescriptor
func (builder SortDescriptorsBuilder[T]) ThenWith(input ...SortDescriptor[T]) SortDescriptorsBuilder[T] {
	result := append(builder, input...)
	return result
}

// GetSortDescriptors Get sortDescriptors
func (builder SortDescriptorsBuilder[T]) GetSortDescriptors() []SortDescriptor[T] {
	return []SortDescriptor[T](builder)
}

// ToSortedList Get the sorted result
func (builder SortDescriptorsBuilder[T]) ToSortedList(input ...T) []T {
	result := SortedListBySortDescriptors(builder.GetSortDescriptors(), input...)
	return result
}

// Sort Sort by sortDescriptors
func (builder SortDescriptorsBuilder[T]) Sort(input []T) {
	SortBySortDescriptors(builder.GetSortDescriptors(), input)
}
