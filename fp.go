package fpGo

import (
	"fmt"
	"reflect"
	"regexp"
	"sync"
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

// Curry

type CurryDef struct {
	fn     func(c *CurryDef, args ...*interface{}) *interface{}
	result *interface{}
	isDone AtomBool

	callM sync.Mutex
	args  []*interface{}
}

func (self *CurryDef) New(fn func(c *CurryDef, args ...*interface{}) *interface{}) *CurryDef {
	c := &CurryDef{fn: fn}

	return c
}
func (self *CurryDef) Call(args ...*interface{}) *CurryDef {
	self.callM.Lock()
	if !self.isDone.Get() {
		self.args = append(self.args, args...)
		self.result = self.fn(self, self.args...)
	}
	self.callM.Unlock()
	return self
}
func (self *CurryDef) MarkDone() {
	self.isDone.Set(true)
}
func (self *CurryDef) IsDone() bool {
	return self.isDone.Get()
}
func (self *CurryDef) Result() *interface{} {
	return self.result
}

var Curry CurryDef

// PatternMatching

type Pattern interface {
	Matches(value *interface{}) bool
	Apply(*interface{}) *interface{}
}
type PatternMatching struct {
	patterns []Pattern
}

type KindPatternDef struct {
	kind   reflect.Kind
	effect fnObj
}
type CompTypePatternDef struct {
	compType CompType
	effect   fnObj
}
type EqualPatternDef struct {
	value  *interface{}
	effect fnObj
}
type RegexPatternDef struct {
	pattern string
	effect  fnObj
}
type OtherwisePatternDef struct {
	effect fnObj
}

func (self PatternMatching) MatchFor(value *interface{}) *interface{} {
	for _, pattern := range self.patterns {
		if pattern.Matches(value) {
			return pattern.Apply(value)
		}
	}

	if value == nil {
		panic(fmt.Sprintf("Cannot match %v", value))
	} else {
		panic(fmt.Sprintf("Cannot match %v", *value))
	}
}

func (self KindPatternDef) Matches(value *interface{}) bool {
	if value == nil {
		return false
	}

	return self.kind == reflect.TypeOf(*value).Kind()
}
func (self CompTypePatternDef) Matches(value *interface{}) bool {
	if value != nil && reflect.TypeOf(*value).Kind() == reflect.TypeOf(CompData{}).Kind() {
		return MatchCompType(self.compType, (*value).(CompData))
	}

	return self.compType.Matches(value)
}
func (self EqualPatternDef) Matches(value *interface{}) bool {
	if value == nil {
		return self.value == value
	}

	return *self.value == *value
}
func (self RegexPatternDef) Matches(value *interface{}) bool {
	if value == nil || reflect.TypeOf(*value).Kind() != reflect.String {
		return false
	}

	matches, err := regexp.MatchString(self.pattern, (*value).(string))
	if err == nil && matches {
		return true
	}

	return false
}
func (self OtherwisePatternDef) Matches(value *interface{}) bool {
	return true
}

func (self KindPatternDef) Apply(value *interface{}) *interface{} {
	return self.effect(value)
}
func (self CompTypePatternDef) Apply(value *interface{}) *interface{} {
	return self.effect(value)
}
func (self EqualPatternDef) Apply(value *interface{}) *interface{} {
	return self.effect(value)
}
func (self RegexPatternDef) Apply(value *interface{}) *interface{} {
	return self.effect(value)
}
func (self OtherwisePatternDef) Apply(value *interface{}) *interface{} {
	return self.effect(value)
}

func DefPattern(patterns ...Pattern) PatternMatching {
	return PatternMatching{patterns: patterns}
}

func InCaseOfKind(kind reflect.Kind, effect fnObj) Pattern {
	return KindPatternDef{kind: kind, effect: effect}
}
func InCaseOfSumType(compType CompType, effect fnObj) Pattern {
	return CompTypePatternDef{compType: compType, effect: effect}
}
func InCaseOfEqual(value *interface{}, effect fnObj) Pattern {
	return EqualPatternDef{value: value, effect: effect}
}
func InCaseOfRegex(pattern string, effect fnObj) Pattern {
	return RegexPatternDef{pattern: pattern, effect: effect}
}
func Otherwise(effect fnObj) Pattern {
	return OtherwisePatternDef{effect: effect}
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
		matches = matches && self.kinds[i] == Maybe.Just(v).Kind()
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
