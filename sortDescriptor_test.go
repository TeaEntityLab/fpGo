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
	return strings.Compare(string(obj), string(input.(TestCustomString)))
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
