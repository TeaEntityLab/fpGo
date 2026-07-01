package fpgo

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCustomInt int

func (obj TestCustomInt) CompareTo(input interface{}) int {
	result := input.(TestCustomInt) - obj
	if result > 0 {
		return 1
	} else if result < 0 {
		return -1
	}

	return 0
}

type TestCustomString string

func (obj TestCustomString) CompareTo(input interface{}) int {
	return strings.Compare(string(input.(TestCustomString)), string(obj))
}

type TestCustomObject struct {
	Name TestCustomString
	Age  TestCustomInt
}

func TestSortDescriptor(t *testing.T) {
	objects := []interface{}{
		TestCustomObject{
			Name: TestCustomString("BC"),
			Age:  TestCustomInt(30),
		},
		TestCustomObject{
			Name: TestCustomString("AD"),
			Age:  TestCustomInt(30),
		},
		TestCustomObject{
			Name: TestCustomString("AB"),
			Age:  TestCustomInt(50),
		},
	}
	objects = NewSortDescriptorsBuilder().
		ThenWithFieldName("Age", false).
		ThenWithFieldName("Name", true).
		ToSortedList(objects...)

	assert.Equal(t, 3, len(objects))

	testOrder := ""
	for _, object := range objects {
		val := object.(TestCustomObject)
		testOrder += fmt.Sprintf("%v%v/", val.Name, val.Age)
	}
	assert.Equal(t, "AB50/AD30/BC30/", testOrder)
}

func TestComparableOrdered(t *testing.T) {
	obj1, obj2 := TestCustomInt(10), TestCustomInt(20)
	assert.Equal(t, 0, obj1.CompareTo(TestCustomInt(10)))
	assert.Equal(t, 1, obj1.CompareTo(obj2))
	assert.Equal(t, -1, obj2.CompareTo(obj1))
}
func TestComparableString(t *testing.T) {
	obj1, obj2 := TestCustomString("apple"), TestCustomString("banana")
	assert.Equal(t, 1, obj1.CompareTo(obj2))
}
func TestSimpleSortDescriptor(t *testing.T) {
	d := NewSimpleSortDescriptor(func(obj interface{}) interface{} { return obj.(TestCustomObject).Age }, true)
	assert.True(t, d.IsAscending())
	assert.Equal(t, TestCustomInt(10), d.TransformedBy()(TestCustomObject{Age: TestCustomInt(10)}))
	d.SetAscending(false)
}
func TestFieldSortDescriptor(t *testing.T) {
	d := NewFieldSortDescriptor("Name", true)
	assert.Equal(t, "Name", d.GetFieldName())
	d.SetFieldName("Age")
}
func TestSortedListBySortDescriptors(t *testing.T) {
	objects := []interface{}{
		TestCustomObject{Name: TestCustomString("BC"), Age: TestCustomInt(30)},
		TestCustomObject{Name: TestCustomString("AD"), Age: TestCustomInt(20)},
		TestCustomObject{Name: TestCustomString("AB"), Age: TestCustomInt(10)},
	}
	descs := NewSortDescriptorsBuilder().ThenWithTransformerFunctor(func(obj interface{}) interface{} {
		return obj.(TestCustomObject).Age
	}, true).GetSortDescriptors()
	r := SortedListBySortDescriptors(descs, objects...)
	assert.Equal(t, TestCustomString("AB"), r[0].(TestCustomObject).Name)
}
func TestSortDescriptorsBuilder(t *testing.T) {
	objects := []interface{}{
		TestCustomObject{Name: TestCustomString("BC"), Age: TestCustomInt(30)},
		TestCustomObject{Name: TestCustomString("AD"), Age: TestCustomInt(20)},
	}
	b := NewSortDescriptorsBuilder().ThenWithTransformerFunctor(func(obj interface{}) interface{} {
		return obj.(TestCustomObject).Age
	}, true).ThenWithFieldName("Name", false).ThenWith(
		NewSimpleSortDescriptor(func(obj interface{}) interface{} { return obj.(TestCustomObject).Age }, true),
	)
	assert.Equal(t, 3, len(b.GetSortDescriptors()))
	input := append([]interface{}{}, objects...)
	b.Sort(input)
}
func TestSortBySortDescriptors(t *testing.T) {
	objects := []interface{}{
		TestCustomObject{Name: TestCustomString("BC"), Age: TestCustomInt(30)},
		TestCustomObject{Name: TestCustomString("AD"), Age: TestCustomInt(20)},
		TestCustomObject{Name: TestCustomString("AB"), Age: TestCustomInt(10)},
	}
	SortBySortDescriptors([]SortDescriptor{NewSimpleSortDescriptor(func(obj interface{}) interface{} {
		return obj.(TestCustomObject).Age
	}, true)}, objects)
	assert.Equal(t, TestCustomString("AB"), objects[0].(TestCustomObject).Name)
}

type maybeName struct{ name *TestCustomString }

func (m maybeName) CompareTo(input interface{}) int {
	o := input.(maybeName)
	if m.name == nil && o.name == nil {
		return 0
	}
	if m.name == nil {
		return -1
	}
	if o.name == nil {
		return 1
	}
	return m.name.CompareTo(*o.name)
}

type TestObjectWithOptionalName struct {
	Name *TestCustomString
	Age  int
}

func sortByOptionalName(objects []interface{}, ascending bool) {
	SortBySortDescriptors([]SortDescriptor{NewSimpleSortDescriptor(func(obj interface{}) interface{} {
		o := obj.(TestObjectWithOptionalName)
		return maybeName{o.Name}
	}, ascending)}, objects)
}

func TestSortDescriptorWithNilKey1Ascending(t *testing.T) {
	nameA := TestCustomString("A")
	objects := []interface{}{
		TestObjectWithOptionalName{nil, 30},
		TestObjectWithOptionalName{&nameA, 20},
	}
	sortByOptionalName(objects, true)
	assert.NotNil(t, objects[0].(TestObjectWithOptionalName).Name)
	assert.Nil(t, objects[1].(TestObjectWithOptionalName).Name)
}

func TestSortDescriptorWithNilKey1Descending(t *testing.T) {
	nameA := TestCustomString("A")
	objects := []interface{}{
		TestObjectWithOptionalName{nil, 30},
		TestObjectWithOptionalName{&nameA, 20},
	}
	sortByOptionalName(objects, false)
	assert.Nil(t, objects[0].(TestObjectWithOptionalName).Name)
	assert.NotNil(t, objects[1].(TestObjectWithOptionalName).Name)
}

func TestSortDescriptorWithNilKey2Ascending(t *testing.T) {
	nameA := TestCustomString("A")
	objects := []interface{}{
		TestObjectWithOptionalName{&nameA, 30},
		TestObjectWithOptionalName{nil, 20},
	}
	sortByOptionalName(objects, true)
	assert.NotNil(t, objects[0].(TestObjectWithOptionalName).Name)
	assert.Nil(t, objects[1].(TestObjectWithOptionalName).Name)
}

func TestSortDescriptorWithNilKey2Descending(t *testing.T) {
	nameA := TestCustomString("A")
	objects := []interface{}{
		TestObjectWithOptionalName{&nameA, 30},
		TestObjectWithOptionalName{nil, 20},
	}
	sortByOptionalName(objects, false)
	assert.Nil(t, objects[0].(TestObjectWithOptionalName).Name)
	assert.NotNil(t, objects[1].(TestObjectWithOptionalName).Name)
}
