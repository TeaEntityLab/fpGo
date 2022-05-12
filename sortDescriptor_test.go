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
