package fpgo

import (
	"reflect"
)

// SortDescriptor Define a Transformer Pattern SortDescriptor
type SortDescriptor interface {
	Transformer

	IsAscending() bool
	SetAscending(bool)
}

// SortedListBySortDescriptors Sort items by sortDescriptors and return value
func SortedListBySortDescriptors(sortDescriptors []SortDescriptor, input ...interface{}) []interface{} {
	result := append(input[:0:0], input...)
	SortBySortDescriptors(sortDescriptors, result)

	return result
}

// SortBySortDescriptors Sort items by sortDescriptors
func SortBySortDescriptors(sortDescriptors []SortDescriptor, input []interface{}) {
	Sort(func(item1 interface{}, item2 interface{}) bool {
		return _compareBySortDescriptors(item1, item2, sortDescriptors, 0) >= 0
	}, input)
}

func _compareBySortDescriptors(item1 interface{}, item2 interface{}, sortDescriptors []SortDescriptor, descriptorIndex int) int {
	descriptor := sortDescriptors[descriptorIndex]
	key1 := descriptor.TransformedBy()(item1).(Comparable)
	key2 := descriptor.TransformedBy()(item2).(Comparable)
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
	if result == 0 && _hasNextDescriptor(sortDescriptors, descriptorIndex) {
		return _compareBySortDescriptors(item1, item2, sortDescriptors, descriptorIndex+1)
	}

	return result
}

func _hasNextDescriptor(sortDescriptors []SortDescriptor, index int) bool {
	return (index + 1) < len(sortDescriptors)
}

// SimpleSortDescriptor

// NewSimpleSortDescriptor Generate a new SimpleSortDescriptor by TransformerFunctor & asceding(true)/descending(false)
func NewSimpleSortDescriptor(transformFn TransformerFunctor, ascending bool) SimpleSortDescriptor {
	return SimpleSortDescriptor{
		transformFn: transformFn,
		ascending:   ascending,
	}
}

// SimpleSortDescriptor SimpleSortDescriptor implemented by TransformerFunctor
type SimpleSortDescriptor struct {
	ascending bool

	transformFn TransformerFunctor
}

// IsAscending Check is this SortDescriptor sorting by ascending
func (descriptor SimpleSortDescriptor) IsAscending() bool {
	return descriptor.ascending
}

// SetAscending Set this SortDescriptor sorting by ascending(true) or descending(false)
func (descriptor SimpleSortDescriptor) SetAscending(val bool) {
	descriptor.ascending = val
}

// TransformedBy Get the TransformerFunctor of this SortDescriptor
func (descriptor SimpleSortDescriptor) TransformedBy() TransformerFunctor {
	return descriptor.transformFn
}

// FieldSortDescriptor

// NewFieldSortDescriptor Generate a new FieldSortDescriptor by FieldName & asceding(true)/descending(false)
func NewFieldSortDescriptor(fieldName string, ascending bool) FieldSortDescriptor {
	return FieldSortDescriptor{
		SimpleSortDescriptor: SimpleSortDescriptor{
			ascending: ascending,
		},

		fieldName: fieldName,
	}
}

// FieldSortDescriptor FieldSortDescriptor implemented by Reflection(by FieldName)
type FieldSortDescriptor struct {
	SimpleSortDescriptor

	fieldName string
}

// GetFieldName Get the fieldName to sort
func (descriptor FieldSortDescriptor) GetFieldName() string {
	return descriptor.fieldName
}

// SetFieldName Set the fieldName to sort
func (descriptor FieldSortDescriptor) SetFieldName(val string) {
	descriptor.fieldName = val
}

// TransformedBy Get the TransformerFunctor of this SortDescriptor
func (descriptor FieldSortDescriptor) TransformedBy() TransformerFunctor {
	return func(input interface{}) interface{} {
		r := reflect.ValueOf(input)
		f := reflect.Indirect(r).FieldByName(descriptor.fieldName).Interface()

		return f
	}
}

// SortDescriptorsBuilder

// NewSortDescriptorsBuilder Generate a new SortDescriptorsBuilder
func NewSortDescriptorsBuilder() SortDescriptorsBuilder {
	return SortDescriptorsBuilder{}
}

// SortDescriptorsBuilder SortDescriptorsBuilder for composing SortDescriptor list and sorting data
type SortDescriptorsBuilder []SortDescriptor

// ThenWithTransformerFunctor Use TransformerFunctor as a SortDescriptor
func (builder SortDescriptorsBuilder) ThenWithTransformerFunctor(transformFn TransformerFunctor, ascending bool) SortDescriptorsBuilder {
	result := append(builder, NewSimpleSortDescriptor(transformFn, ascending))
	return result
}

// ThenWithFieldName Use FieldName as a SortDescriptor
func (builder SortDescriptorsBuilder) ThenWithFieldName(fieldName string, ascending bool) SortDescriptorsBuilder {
	result := append(builder, NewFieldSortDescriptor(fieldName, ascending))
	return result
}

// ThenWith Append a SortDescriptor
func (builder SortDescriptorsBuilder) ThenWith(input ...SortDescriptor) SortDescriptorsBuilder {
	result := append(builder, input...)
	return result
}

// GetSortDescriptors Get sortDescriptors
func (builder SortDescriptorsBuilder) GetSortDescriptors() []SortDescriptor {
	return []SortDescriptor(builder)
}

// ToSortedList Get the sorted result
func (builder SortDescriptorsBuilder) ToSortedList(input ...interface{}) []interface{} {
	result := SortedListBySortDescriptors(builder.GetSortDescriptors(), input...)
	return result
}

// Sort Sort by sortDescriptors
func (builder SortDescriptorsBuilder) Sort(input []interface{}) {
	SortBySortDescriptors(builder.GetSortDescriptors(), input)
}
