package fpGo

import (
	"reflect"
)

type fnObj func(*interface{}) *interface{}

func Compose(fnList ...fnObj) fnObj {
	return func(s *interface{}) *interface{} {
		f := fnList[0]
		nextFnList := fnList[1:len(fnList)]

		if len(fnList) == 1 {
			return f(s)
		}

		return f(Compose(nextFnList...)(s))
	}
}

type CompData struct {
	compType CompType
	objects  []*interface{}
}

type CompType interface {
	Matches(value ...*interface{}) bool
}

type SumType struct {
	compTypes []CompType
}
type ProductType struct {
	kinds []reflect.Kind
}
type NilType struct {
}

func (self SumType) Matches(value ...*interface{}) bool {
	for _, compType := range self.compTypes {
		if compType.Matches(value...) {
			return true
		}
	}

	return false
}
func (self ProductType) Matches(value ...*interface{}) bool {
	if len(value) != len(self.kinds) {
		return false
	}

	matches := true
	for i, v := range value {
		matches = matches && self.kinds[i] == Monad.Just(v).Kind()
	}
	return matches
}
func (self NilType) Matches(value ...*interface{}) bool {
	if len(value) != 1 {
		return false
	}

	return value[0] == nil
}

func DefSum(compTypes ...CompType) CompType {
	return SumType{compTypes: compTypes}
}

func DefProduct(kinds ...reflect.Kind) CompType {
	return ProductType{kinds: kinds}
}

func NewCompData(compType CompType, value ...*interface{}) *CompData {
	if compType.Matches(value...) {
		return &CompData{compType: compType, objects: value}
	}

	return nil
}
