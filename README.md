# fpGo

[![tag](https://img.shields.io/github/tag/TeaEntityLab/fpGo.svg)](https://github.com/TeaEntityLab/fpGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/TeaEntityLab/fpGo)](https://goreportcard.com/report/github.com/TeaEntityLab/fpGo)
[![codecov](https://codecov.io/gh/TeaEntityLab/fpGo/branch/master/graph/badge.svg)](https://codecov.io/gh/TeaEntityLab/fpGo)
[![Travis CI Build Status](https://travis-ci.com/TeaEntityLab/fpGo.svg?branch=master)](https://travis-ci.com/TeaEntityLab/fpGo)
[![GoDoc](https://godoc.org/github.com/TeaEntityLab/fpGo?status.svg)](https://godoc.org/github.com/TeaEntityLab/fpGo)

[![license](https://img.shields.io/github/license/TeaEntityLab/fpGo.svg?style=social&label=License)](https://github.com/TeaEntityLab/fpGo)
[![stars](https://img.shields.io/github/stars/TeaEntityLab/fpGo.svg?style=social&label=Stars)](https://github.com/TeaEntityLab/fpGo)
[![forks](https://img.shields.io/github/forks/TeaEntityLab/fpGo.svg?style=social&label=Fork)](https://github.com/TeaEntityLab/fpGo)

Monad, Functional Programming features for Golang

# Active Branches:

For *Generics* version(*`>=go1.18`*):[generics](https://github.com/TeaEntityLab/fpGo/tree/generics)

For *NonGenerics* version(*`<=go1.17`*):[non-generics](https://github.com/TeaEntityLab/fpGo/tree/non_generic)

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

* Akka/Erlang-like Actor model(send/receive/spawn/states)

# Special thanks
* fp functions(Dedupe/Difference/Distinct/IsDistinct/DropEq/Drop/DropLast/DropWhile/IsEqual/IsEqualMap/Every/Exists/Intersection/Keys/Values/Max/Min/MinMax/Merge/IsNeg/IsPos/PMap/Range/Reverse/Minus/Some/IsSubset/IsSuperset/Take/TakeLast/Union/IsZero/Zip/GroupBy/UniqBy/Flatten/Prepend/Partition/Tail/Head/SplitEvery)
  *	Credit: https://github.com/logic-building/functional-go
  * Credit: https://github.com/achannarasappa/pneumatic

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

## Actor (inspired by Akka/Erlang)

### Actor common(send/receive/spawn/states)

Example:

```go
actual := 0
// Channel for results
resultChannel := make(chan interface{}, 1)
// Message CMDs
cmdSpawn := "spawn"
cmdShutdown := "shutdown"
// Testee
actorRoot := Actor.New(func(self *ActorDef, input interface{}) {
  // SPAWN: for ROOT
  if input == cmdSpawn {
    self.Spawn(func(self *ActorDef, input interface{}) {
      // SHUTDOWN: for Children
      if input == cmdShutdown {
        self.Close()
        return
      }

      // INT cases: Children
      val, _ := Maybe.Just(input).ToInt()
      resultChannel <- val * 10
    })
    return
  }
  // SHUTDOWN: for ROOT
  if input == cmdShutdown {
    for _, child := range self.children {
      child.Send(cmdShutdown)
    }
    self.Close()

    close(resultChannel)
    return
  }

  // INT cases: ROOT
  intVal, _ := Maybe.Just(input).ToInt()
  if intVal > 0 {
    for _, child := range self.children {
      child.Send(intVal)
    }
  }
})

// Sequential Send messages(async)
go func() {
  actorRoot.Send(cmdSpawn)
  actorRoot.Send(10)
  actorRoot.Send(cmdSpawn)
  actorRoot.Send(20)
  actorRoot.Send(cmdSpawn)
  actorRoot.Send(30)
}()

i := 0
for val := range resultChannel {
  intVal, _ := Maybe.Just(val).ToInt()
  actual += intVal

  i++
  if i == 5 {
    go actorRoot.Send(cmdShutdown)
  }
}

// Result would be 1400 (=10*10+20*10+20*10+30*10+30*10+30*10)
fmt.Println(actual)
```

### Actor Ask (inspired by Akka/Erlang)

```go
actorRoot := Actor.New(func(self *ActorDef, input interface{}) {
    // Ask cases: ROOT
    switch val := input.(type) {
    case *AskDef:
        intVal, _ := Maybe.Just(val.Message).ToInt()

        // NOTE If negative, hanging for testing Ask.timeout
        if intVal < 0 {
            break
        }

        val.Reply(intVal * 10)
        break
    }
})

// var timer *time.Timer
var timeout time.Duration
timeout = 10 * time.Millisecond

// Normal cases
// Result would be 10
actual, _ = Maybe.Just(Ask.New(1, nil).AskOnce(actorRoot)).ToInt()
// Ask with Timeout
// Result would be 20
actual, _ = Maybe.Just(Ask.New(2, &timeout).AskOnce(actorRoot)).ToInt()
// Ask channel
// Result would be 30
ch, timer := Ask.New(3, &timeout).AskChannel(actorRoot)
actual, _ = Maybe.Just(<-*ch).ToInt()
close(*ch)
if timer != nil {
    timer.Stop()
}

// Timeout cases
// Result would be 0 (zero value, timeout)
actual, _ = Maybe.Just(Ask.New(-1, &timeout).AskOnce(actorRoot)).ToInt()
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
