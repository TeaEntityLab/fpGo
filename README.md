# fpGo

[![tag](https://img.shields.io/github/tag/TeaEntityLab/fpGo.svg)](https://github.com/TeaEntityLab/fpGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/TeaEntityLab/fpGo)](https://goreportcard.com/report/github.com/TeaEntityLab/fpGo)
[![GoDoc](https://godoc.org/github.com/TeaEntityLab/fpGo?status.svg)](https://godoc.org/github.com/TeaEntityLab/fpGo)

[![license](https://img.shields.io/github/license/TeaEntityLab/fpGo.svg?style=social&label=License)](https://github.com/TeaEntityLab/fpGo)
[![stars](https://img.shields.io/github/stars/TeaEntityLab/fpGo.svg?style=social&label=Stars)](https://github.com/TeaEntityLab/fpGo)
[![forks](https://img.shields.io/github/forks/TeaEntityLab/fpGo.svg?style=social&label=Fork)](https://github.com/TeaEntityLab/fpGo)

Monad, Functional Programming features for Golang

# Why

I love functional programing, Rx-style coding, and Optional usages.

However it's hard to implement them in Golang, and there're few libraries to achieve parts of them.

Thus I implemented fpGo. I hope you would like it :)

# Features

* Optional/Maybe

* Monad, Rx-like

* Publisher



* Pattern matching

* Fp functions



* Java8Stream-like Collection

* PythonicGenerator-like Coroutine(yield/yieldFrom)


# Usage

## Optional (IsPresent/IsNil, Or, Let)

```go
var m MaybeDef
var orVal int
var boolVal bool

// IsPresent(), IsNil()
m = Maybe.Just(1)
boolVal = m.IsPresent() // true
boolVal = m.IsNil() // false
m = Maybe.Just(nil)
boolVal = m.IsPresent() // false
boolVal = m.IsNil() // true

// Or()
m = Maybe.Just(1)
fmt.Println((m.Or(3))) // 1
m = Maybe.Just(nil)
fmt.Println((m.Or(3))) // 3

// Let()
var letVal int
letVal = 1
m = Maybe.Just(1)
m.Let(func() {
  letVal = 2
})
fmt.Println(letVal) // letVal would be 2

letVal = 1
m = Maybe.Just(nil)
m.Let(func() {
  letVal = 3
})
fmt.Println(letVal) // letVal would be still 1
```

## MonadIO (RxObserver-like)

Example:
```go
var m *MonadIODef
var actualInt int

m = MonadIO.Just(1)
actualInt = 0
m.Subscribe(Subscription{
  OnNext: func(in interface{}) {
    actualInt, _ = Maybe.Just(in).ToInt()
  },
})
fmt.Println(actualInt) // actualInt would be 1

m = MonadIO.Just(1).FlatMap(func(in interface{}) *MonadIODef {
  v, _ := Maybe.Just(in).ToInt()
  return MonadIO.Just(v + 1)
})
actualInt = 0
m.Subscribe(Subscription{
  OnNext: func(in interface{}) {
    actualInt, _ = Maybe.Just(in).ToInt()
  },
})
fmt.Println(actualInt) // actualInt would be 2
```

## Stream (inspired by Collection libs)

Example:
```go
var s *StreamDef
var tempString = ""

s = Stream.FromArrayInt([]int{}).Append(1).Extend(Stream.FromArrayInt([]int{2, 3, 4})).Extend(Stream.FromArray([]interface{}{nil}))
tempString = ""
for _, v := range s.ToArray() {
  tempString += Maybe.Just(v).ToMaybe().ToString()
}
fmt.Println(tempString) // tempString would be "1234<nil>"
s = s.Distinct()
tempString = ""
for _, v := range s.ToArray() {
  tempString += Maybe.Just(v).ToMaybe().ToString()
}
fmt.Println(tempString) // tempString would be "1234"
```

## Compose

Example:

```go
var fn01 = func(args ...interface{}) []interface{} {
  val, _ := Maybe.Just(args[0]).ToInt()
  return SliceOf(val + 1)
}
var fn02 = func(args ...interface{}) []interface{} {
  val, _ := Maybe.Just(args[0]).ToInt()
  return SliceOf(val + 2)
}
var fn03 = func(args ...interface{}) []interface{} {
  val, _ := Maybe.Just(args[0]).ToInt()
  return SliceOf(val + 3)
}

// Result would be 6
result := Compose(fn01, fn02, fn03)((0))[0]
```
