package fpGo

import (
	"fmt"
	"reflect"
	"regexp"
	"sync"
)

type fnObj func(interface{}) interface{}

func Compose(fnList ...fnObj) fnObj {
	return func(s interface{}) interface{} {
		f := fnList[0]
		nextFnList := fnList[1:]

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
	fn     func(c *CurryDef, args ...interface{}) interface{}
	result interface{}
	isDone AtomBool

	callM sync.Mutex
	args  []interface{}
}

func (self *CurryDef) New(fn func(c *CurryDef, args ...interface{}) interface{}) *CurryDef {
	c := &CurryDef{fn: fn}

	return c
}
func (self *CurryDef) Call(args ...interface{}) *CurryDef {
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
func (self *CurryDef) Result() interface{} {
	return self.result
}

var Curry CurryDef

// PatternMatching

type Pattern interface {
	Matches(value interface{}) bool
	Apply(interface{}) interface{}
}
type PatternMatching struct {
	patterns []Pattern
}

// KindPatternDef Pattern which matching when the kind matches
type KindPatternDef struct {
	kind   reflect.Kind
	effect fnObj
}

// CompTypePatternDef Pattern which matching when the SumType matches
type CompTypePatternDef struct {
	compType CompType
	effect   fnObj
}

// EqualPatternDef Pattern which matching when the given object is equal to predefined one
type EqualPatternDef struct {
	value  interface{}
	effect fnObj
}

// RegexPatternDef Pattern which matching when the regex rule matches the given string
type RegexPatternDef struct {
	pattern string
	effect  fnObj
}

// OtherwisePatternDef Pattern which matching when the others didn't match(finally)
type OtherwisePatternDef struct {
	effect fnObj
}

func (self PatternMatching) MatchFor(inValue interface{}) interface{} {
	for _, pattern := range self.patterns {
		value := inValue
		maybe := Maybe.Just(inValue)
		if maybe.IsKind(reflect.Ptr) {
			ptr := maybe.ToPtr()
			if reflect.TypeOf(*ptr).Kind() == (reflect.TypeOf(CompData{}).Kind()) {
				value = *ptr
			}
		}

		if pattern.Matches(value) {
			return pattern.Apply(value)
		}
	}

	panic(fmt.Sprintf("Cannot match %v", inValue))
}

// Matches Match the given value by the pattern
func (patternSelf KindPatternDef) Matches(value interface{}) bool {
	if value == nil {
		return false
	}

	return patternSelf.kind == reflect.TypeOf(value).Kind()
}

// Matches Match the given value by the pattern
func (patternSelf CompTypePatternDef) Matches(value interface{}) bool {
	if value != nil && reflect.TypeOf(value).Kind() == reflect.TypeOf(CompData{}).Kind() {
		return MatchCompType(patternSelf.compType, (value).(CompData))
	}

	return patternSelf.compType.Matches(value)
}

// Matches Match the given value by the pattern
func (patternSelf EqualPatternDef) Matches(value interface{}) bool {
	return patternSelf.value == value
}

// Matches Match the given value by the pattern
func (patternSelf RegexPatternDef) Matches(value interface{}) bool {
	if value == nil || reflect.TypeOf(value).Kind() != reflect.String {
		return false
	}

	matches, err := regexp.MatchString(patternSelf.pattern, (value).(string))
	if err == nil && matches {
		return true
	}

	return false
}

// Matches Match the given value by the pattern
func (patternSelf OtherwisePatternDef) Matches(value interface{}) bool {
	return true
}

// Apply Evaluate the result by its given effect function
func (patternSelf KindPatternDef) Apply(value interface{}) interface{} {
	return patternSelf.effect(value)
}

// Apply Evaluate the result by its given effect function
func (patternSelf CompTypePatternDef) Apply(value interface{}) interface{} {
	return patternSelf.effect(value)
}

// Apply Evaluate the result by its given effect function
func (patternSelf EqualPatternDef) Apply(value interface{}) interface{} {
	return patternSelf.effect(value)
}

// Apply Evaluate the result by its given effect function
func (patternSelf RegexPatternDef) Apply(value interface{}) interface{} {
	return patternSelf.effect(value)
}

// Apply Evaluate the result by its given effect function
func (patternSelf OtherwisePatternDef) Apply(value interface{}) interface{} {
	return patternSelf.effect(value)
}

func DefPattern(patterns ...Pattern) PatternMatching {
	return PatternMatching{patterns: patterns}
}

// InCaseOfKind In case of its Kind matches the given one
func InCaseOfKind(kind reflect.Kind, effect fnObj) Pattern {
	return KindPatternDef{kind: kind, effect: effect}
}

// InCaseOfSumType In case of its SumType matches the given one
func InCaseOfSumType(compType CompType, effect fnObj) Pattern {
	return CompTypePatternDef{compType: compType, effect: effect}
}

// InCaseOfEqual In case of its value is equal to the given one
func InCaseOfEqual(value interface{}, effect fnObj) Pattern {
	return EqualPatternDef{value: value, effect: effect}
}

// InCaseOfRegex In case of the given regex rule matches its value
func InCaseOfRegex(pattern string, effect fnObj) Pattern {
	return RegexPatternDef{pattern: pattern, effect: effect}
}

// Otherwise In case of the other patterns didn't match it
func Otherwise(effect fnObj) Pattern {
	return OtherwisePatternDef{effect: effect}
}

func Either(value interface{}, patterns ...Pattern) interface{} {
	return DefPattern(patterns...).MatchFor(value)
}

// SumType

type CompData struct {
	compType CompType
	objects  []interface{}
}

type CompType interface {
	Matches(value ...interface{}) bool
}

type SumType struct {
	compTypes []CompType
}
type ProductType struct {
	kinds []reflect.Kind
}
type NilTypeDef struct {
}

func (typeSelf SumType) Matches(value ...interface{}) bool {
	for _, compType := range typeSelf.compTypes {
		if compType.Matches(value...) {
			return true
		}
	}

	return false
}
func (typeSelf ProductType) Matches(value ...interface{}) bool {
	if len(value) != len(typeSelf.kinds) {
		return false
	}

	matches := true
	for i, v := range value {
		matches = matches && typeSelf.kinds[i] == Maybe.Just(v).Kind()
	}
	return matches
}
func (typeSelf NilTypeDef) Matches(value ...interface{}) bool {
	if len(value) != 1 {
		return false
	}

	return Maybe.Just(value[0]).IsNil()
}

func DefSum(compTypes ...CompType) CompType {
	return SumType{compTypes: compTypes}
}

func DefProduct(kinds ...reflect.Kind) CompType {
	return ProductType{kinds: kinds}
}

func NewCompData(compType CompType, value ...interface{}) *CompData {
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
