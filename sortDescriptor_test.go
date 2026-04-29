package fpgo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCustomObject struct {
	Name   ComparableString
	Age    int
	Height float64
}

func TestSortDescriptor(t *testing.T) {
	objects := []TestCustomObject{
		{
			Name: NewComparableString("BC"),
			Age:  30,
		},
		{
			Name: NewComparableString("AD"),
			Age:  30,
		},
		{
			Name: NewComparableString("AB"),
			Age:  50,
		},
	}
	objects = NewSortDescriptorsBuilder[TestCustomObject]().
		ThenWithTransformerFunctor(func(obj TestCustomObject) Comparable[interface{}] {
			return NewComparableOrdered(obj.Age)
		}, false).
		ThenWithFieldName("Name", true).
		ToSortedList(objects...)

	assert.Equal(t, 3, len(objects))

	testOrder := ""
	for _, object := range objects {
		testOrder += fmt.Sprintf("%v%v/", object.Name.Val, object.Age)
	}
	assert.Equal(t, "AB50/AD30/BC30/", testOrder)
}

func TestComparableOrdered(t *testing.T) {
	obj1 := NewComparableOrdered(10)
	obj2 := NewComparableOrdered(20)

	assert.Equal(t, 0, obj1.CompareTo(ComparableOrdered[int]{Val: 10}))
	assert.Equal(t, 1, obj1.CompareTo(obj2))
	assert.Equal(t, -1, obj2.CompareTo(obj1))
}

func TestComparableString(t *testing.T) {
	obj1 := NewComparableString("apple")
	obj2 := NewComparableString("banana")

	assert.Equal(t, 0, obj1.CompareTo(ComparableString{Val: "apple"}))
	assert.Equal(t, -1, obj1.CompareTo(obj2))
	assert.Equal(t, 1, obj2.CompareTo(obj1))
}

func TestSortedListBySortDescriptors(t *testing.T) {
	objects := []TestCustomObject{
		{Name: NewComparableString("BC"), Age: 30},
		{Name: NewComparableString("AD"), Age: 20},
		{Name: NewComparableString("AB"), Age: 10},
	}

	descs := NewSortDescriptorsBuilder[TestCustomObject]().
		ThenWithTransformerFunctor(func(obj TestCustomObject) Comparable[interface{}] {
			return NewComparableOrdered(obj.Age)
		}, true).GetSortDescriptors()

	result := SortedListBySortDescriptors(descs, objects...)

	assert.Equal(t, "AB", result[0].Name.Val)
	assert.Equal(t, "AD", result[1].Name.Val)
	assert.Equal(t, "BC", result[2].Name.Val)
}

func TestSimpleSortDescriptor(t *testing.T) {
	descriptor := NewSimpleSortDescriptor(func(obj TestCustomObject) Comparable[interface{}] {
		return NewComparableOrdered(obj.Age)
	}, true)

	assert.True(t, descriptor.IsAscending())

	transformed := descriptor.TransformedBy()(TestCustomObject{Age: 10})
	assert.NotNil(t, transformed)

	// Exercise value-receiver setter path for coverage.
	descriptor.SetAscending(false)
}

func TestFieldSortDescriptor(t *testing.T) {
	descriptor := NewFieldSortDescriptor[TestCustomObject]("Name", true)

	assert.True(t, descriptor.IsAscending())
	assert.Equal(t, "Name", descriptor.GetFieldName())

	transformed := descriptor.TransformedBy()(TestCustomObject{Name: NewComparableString("test"), Age: 25})
	assert.NotNil(t, transformed)

	// Exercise value-receiver setter path for coverage.
	descriptor.SetFieldName("Age")

	desc2 := NewFieldSortDescriptor[TestCustomObject]("Age", false)
	assert.False(t, desc2.IsAscending())
	assert.Equal(t, "Age", desc2.GetFieldName())
}

func TestSortDescriptorsBuilder(t *testing.T) {
	objects := []TestCustomObject{
		{Name: NewComparableString("BC"), Age: 30},
		{Name: NewComparableString("AD"), Age: 20},
	}

	desc1 := NewSimpleSortDescriptor(func(obj TestCustomObject) Comparable[interface{}] {
		return NewComparableOrdered(obj.Age)
	}, true)
	desc2 := NewFieldSortDescriptor[TestCustomObject]("Name", false)

	builder := NewSortDescriptorsBuilder[TestCustomObject]().
		ThenWithTransformerFunctor(func(obj TestCustomObject) Comparable[interface{}] {
			return NewComparableOrdered(obj.Age)
		}, true).
		ThenWithFieldName("Name", false).
		ThenWith(desc1, desc2)

	descs := builder.GetSortDescriptors()
	assert.Equal(t, 4, len(descs))

	result := builder.ToSortedList(objects...)
	assert.Equal(t, 2, len(result))

	input := make([]TestCustomObject, len(objects))
	copy(input, objects)
	builder.Sort(input)
	assert.Equal(t, 2, len(input))
}

func TestSortBySortDescriptors(t *testing.T) {
	objects := []TestCustomObject{
		{Name: NewComparableString("BC"), Age: 30},
		{Name: NewComparableString("AD"), Age: 20},
		{Name: NewComparableString("AB"), Age: 10},
	}

	descs := []SortDescriptor[TestCustomObject]{
		NewSimpleSortDescriptor(func(obj TestCustomObject) Comparable[interface{}] {
			return NewComparableOrdered(obj.Age)
		}, true),
	}

	SortBySortDescriptors(descs, objects)

	assert.Equal(t, "AB", objects[0].Name.Val)
	assert.Equal(t, "AD", objects[1].Name.Val)
	assert.Equal(t, "BC", objects[2].Name.Val)
}

// Tests for sortDescriptor.go nil key handling coverage

type TestObjectWithOptionalName struct {
	Name *string
	Age  int
}

func TestSortDescriptorWithNilKey1Ascending(t *testing.T) {
	name1 := "A"
	objects := []TestObjectWithOptionalName{
		{Name: nil, Age: 30},
		{Name: &name1, Age: 20},
	}

	descs := []SortDescriptor[TestObjectWithOptionalName]{
		NewSimpleSortDescriptor(func(obj TestObjectWithOptionalName) Comparable[interface{}] {
			if obj.Name == nil {
				return nil
			}
			return NewComparableString(*obj.Name)
		}, true),
	}

	SortBySortDescriptors(descs, objects)

	// With ascending, nil should come after non-nil
	assert.NotNil(t, objects[0].Name)
	assert.Nil(t, objects[1].Name)
}

func TestSortDescriptorWithNilKey1Descending(t *testing.T) {
	name1 := "A"
	objects := []TestObjectWithOptionalName{
		{Name: nil, Age: 30},
		{Name: &name1, Age: 20},
	}

	descs := []SortDescriptor[TestObjectWithOptionalName]{
		NewSimpleSortDescriptor(func(obj TestObjectWithOptionalName) Comparable[interface{}] {
			if obj.Name == nil {
				return nil
			}
			return NewComparableString(*obj.Name)
		}, false),
	}

	SortBySortDescriptors(descs, objects)

	// With descending, nil should come before non-nil
	assert.Nil(t, objects[0].Name)
	assert.NotNil(t, objects[1].Name)
}

func TestSortDescriptorWithNilKey2Ascending(t *testing.T) {
	name1 := "A"
	objects := []TestObjectWithOptionalName{
		{Name: &name1, Age: 30},
		{Name: nil, Age: 20},
	}

	descs := []SortDescriptor[TestObjectWithOptionalName]{
		NewSimpleSortDescriptor(func(obj TestObjectWithOptionalName) Comparable[interface{}] {
			if obj.Name == nil {
				return nil
			}
			return NewComparableString(*obj.Name)
		}, true),
	}

	SortBySortDescriptors(descs, objects)

	// With ascending, nil should come after non-nil
	assert.NotNil(t, objects[0].Name)
	assert.Nil(t, objects[1].Name)
}

func TestSortDescriptorWithNilKey2Descending(t *testing.T) {
	name1 := "A"
	objects := []TestObjectWithOptionalName{
		{Name: &name1, Age: 30},
		{Name: nil, Age: 20},
	}

	descs := []SortDescriptor[TestObjectWithOptionalName]{
		NewSimpleSortDescriptor(func(obj TestObjectWithOptionalName) Comparable[interface{}] {
			if obj.Name == nil {
				return nil
			}
			return NewComparableString(*obj.Name)
		}, false),
	}

	SortBySortDescriptors(descs, objects)

	// With descending, nil should come before non-nil
	assert.Nil(t, objects[0].Name)
	assert.NotNil(t, objects[1].Name)
}
