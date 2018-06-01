package fpGo

import (
	"fmt"
	"reflect"
	"regexp"
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

func PtrOf(v interface{}) *interface{} {
	return &v
}

// PatternMatching

type Pattern interface {
	Matches(value *interface{}) bool
	Apply(*interface{}) *interface{}
}
type PatternMatching struct {
	patterns []Pattern
}

type KindPattern struct {
	kind   reflect.Kind
	effect fnObj
}
type CompTypePattern struct {
	compType CompType
	effect   fnObj
}
type EqualPattern struct {
	value  *interface{}
	effect fnObj
}
type RegexPattern struct {
	pattern string
	effect  fnObj
}
type OtherwisePattern struct {
	effect fnObj
}

func (self PatternMatching) MatchFor(value *interface{}) *interface{} {
	for _, pattern := range self.patterns {
		if pattern.Matches(value) {
			return pattern.Apply(value)
		}
	}

	panic(fmt.Sprintf("Cannot match %v", value))
}

func (self KindPattern) Matches(value *interface{}) bool {
	if value == nil {
		return false
	}

	return self.kind == reflect.TypeOf(*value).Kind()
}
func (self CompTypePattern) Matches(value *interface{}) bool {
	return self.compType.Matches(value)
}
func (self EqualPattern) Matches(value *interface{}) bool {
	return self.value == value
}
func (self RegexPattern) Matches(value *interface{}) bool {
	if value == nil || reflect.TypeOf(*value).Kind() != reflect.String {
		return false
	}

	matches, err := regexp.MatchString("p([a-z]+)ch", "peach")
	if err == nil && matches {
		return true
	}

	return false
}
func (self OtherwisePattern) Matches(value *interface{}) bool {
	return true
}

func (self KindPattern) Apply(value *interface{}) *interface{} {
	return self.effect(value)
}
func (self CompTypePattern) Apply(value *interface{}) *interface{} {
	return self.effect(value)
}
func (self EqualPattern) Apply(value *interface{}) *interface{} {
	return self.effect(value)
}
func (self RegexPattern) Apply(value *interface{}) *interface{} {
	return self.effect(value)
}
func (self OtherwisePattern) Apply(value *interface{}) *interface{} {
	return self.effect(value)
}

func DefPattern(patterns ...Pattern) PatternMatching {
	return PatternMatching{patterns: patterns}
}

func NewKindPattern(kind reflect.Kind, effect fnObj) Pattern {
	return KindPattern{kind: kind, effect: effect}
}
func NewSumTypePattern(compType CompType, effect fnObj) Pattern {
	return CompTypePattern{compType: compType, effect: effect}
}
func NewEqualPattern(value *interface{}, effect fnObj) Pattern {
	return EqualPattern{value: value, effect: effect}
}
func NewRegexPattern(pattern string, effect fnObj) Pattern {
	return RegexPattern{pattern: pattern, effect: effect}
}
func NewOtherwisePattern(effect fnObj) Pattern {
	return OtherwisePattern{effect: effect}
}

func Either(value interface{}, patterns ...Pattern) *interface{} {
	return EitherRef(&value, patterns...)
}

func EitherRef(value *interface{}, patterns ...Pattern) *interface{} {
	return DefPattern(patterns...).MatchFor(value)
}

// SumType

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
type NilTypeDef struct {
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
func (self NilTypeDef) Matches(value ...*interface{}) bool {
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

func NewCompDataVal(compType CompType, value ...interface{}) *CompData {
	list := make([]*interface{}, len(value))
	for i, v := range value {
		list[i] = &v
	}

	return NewCompData(compType, list...)
}

func NewCompData(compType CompType, value ...*interface{}) *CompData {
	if compType.Matches(value...) {
		return &CompData{compType: compType, objects: value}
	}

	return nil
}

func MatchCompType(compType CompType, value CompData) bool {
	return MatchCompTypeRef(compType, &value)
}
func MatchCompTypeRef(compType CompType, value *CompData) bool {
	return compType.Matches(value.objects...)
}

var NilType NilTypeDef
