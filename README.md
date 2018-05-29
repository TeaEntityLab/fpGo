# fpGo
Monad, Functional Programming features for Golang

## Optional (IsPresent/IsNil, Or, Let)

```go
var m MonadDef
var orVal int
var boolVal bool

// IsPresent(), IsNil()
m = Monad.JustVal(1)
boolVal = m.IsPresent() // true
boolVal = m.IsNil() // false
m = Monad.Just(nil)
boolVal = m.IsPresent() // false
boolVal = m.IsNil() // true

// Or()
m = Monad.JustVal(1)
orVal = m.OrVal(3).Unwrap() // orVal would be 1
m = Monad.Just(nil)
orVal = m.OrVal(3).Unwrap() // orVal would be 3

// Let()
var letVal int
letVal = 1
m = Monad.JustVal(1)
m.Let(func() {
  letVal = 2
})
fmt.Println(letVal) // letVal would be 2

letVal = 1
m = Monad.Just(nil)
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

m = MonadIO.JustVal(1)
actualInt = 0
m.Subscribe(Subscription{
  OnNext: func(in *interface{}) {
    actualInt, _ = Monad.Just(in).ToInt()
  },
})
fmt.Println(actualInt) // actualInt would be 1

m = MonadIO.JustVal(1).FlatMap(func(in *interface{}) *MonadIODef {
  v, _ := Monad.Just(in).ToInt()
  return MonadIO.JustVal(v + 1)
})
actualInt = 0
m.Subscribe(Subscription{
  OnNext: func(in *interface{}) {
    actualInt, _ = Monad.Just(in).ToInt()
  },
})
fmt.Println(actualInt) // actualInt would be 2
```

## Stream (inspired by Collection libs)

Example:
```go
var s *StreamDef
var tempString = ""

s = Stream.FromArrayInt([]int{}).Append(Monad.JustVal(1).Ref()).Extend(Stream.FromArrayInt([]int{2, 3, 4})).Extend(Stream.FromArray([]*interface{}{Monad.Just(nil).Ref()}))
tempString = ""
for _, v := range s.ToArray() {
  tempString += Monad.Just(v).ToMonad().ToString()
}
fmt.Println(tempString) // tempString would be "1234<nil>"
s = s.Distinct()
tempString = ""
for _, v := range s.ToArray() {
  tempString += Monad.Just(v).ToMonad().ToString()
}
fmt.Println(tempString) // tempString would be "1234"
```

## Compose

Example:

```go
var fn01 = func(obj *interface{}) *interface{} {
  val, _ := Monad.Just(obj).ToInt()
  return Monad.JustVal(val + 1).Ref()
}
var fn02 = func(obj *interface{}) *interface{} {
  val, _ := Monad.Just(obj).ToInt()
  return Monad.JustVal(val + 2).Ref()
}
var fn03 = func(obj *interface{}) *interface{} {
  val, _ := Monad.Just(obj).ToInt()
  return Monad.JustVal(val + 3).Ref()
}

// Result would be 6
result := *Compose(fn01, fn02, fn03)(Monad.JustVal(0).Ref())
```
