package fpgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonadIO(t *testing.T) {
	var m *MonadIODef[interface{}]
	var actualInt int

	m = MonadIO.Just(1)
	actualInt = 0
	m.Subscribe(Subscription[interface{}]{
		OnNext: func(in interface{}) {
			actualInt, _ = Maybe.Just(in).ToInt()
		},
	})
	assert.Equal(t, 1, actualInt)

	m = MonadIO.Just(1).FlatMap(func(in interface{}) *MonadIODef[interface{}] {
		v, _ := Maybe.Just(in).ToInt()
		return MonadIO.Just(v + 1)
	})
	actualInt = 0
	m.Subscribe(Subscription[interface{}]{
		OnNext: func(in interface{}) {
			actualInt, _ = Maybe.Just(in).ToInt()
		},
	})
	assert.Equal(t, 2, actualInt)

	actualInt = 0
	m = MonadIO.New(func() interface{} {
		actualInt = 3
		return 0
	})
	assert.Equal(t, 0, actualInt)
	m.Eval()
	assert.Equal(t, 3, actualInt)
}

func TestMonadIOObserveOn(t *testing.T) {
	m := MonadIO.Just(1).ObserveOn(nil)
	assert.NotNil(t, m)
}

func TestMonadIOSubscribeOn(t *testing.T) {
	m := MonadIO.Just(1).SubscribeOn(nil)
	assert.NotNil(t, m)
}

func TestMonadIOFlatMap(t *testing.T) {
	m := MonadIO.Just(1).FlatMap(func(in interface{}) *MonadIODef[interface{}] {
		v, _ := Maybe.Just(in).ToInt()
		return MonadIO.Just(v * 2)
	})
	result := m.Eval()
	assert.Equal(t, 2, result)
}

func TestMonadIOEval(t *testing.T) {
	var called bool
	m := MonadIO.New(func() interface{} {
		called = true
		return 42
	})
	assert.Equal(t, false, called)
	result := m.Eval()
	assert.Equal(t, true, called)
	assert.Equal(t, 42, result)
}

func TestMonadIOGenerics(t *testing.T) {
	m := MonadIOJustGenerics(42)
	assert.NotNil(t, m)

	result := m.Eval()
	assert.Equal(t, 42, result)

	m2 := MonadIONewGenerics(func() int {
		return 100
	})
	assert.NotNil(t, m2)

	result2 := m2.Eval()
	assert.Equal(t, 100, result2)
}

func TestMonadIONew(t *testing.T) {
	var called bool
	m := MonadIO.New(func() interface{} {
		called = true
		return 42
	})
	assert.False(t, called)

	result := m.Eval()
	assert.True(t, called)
	assert.Equal(t, 42, result)
}

func TestMonadIOFlatMapChain(t *testing.T) {
	m := MonadIO.Just(1).
		FlatMap(func(in interface{}) *MonadIODef[interface{}] {
			v, _ := Maybe.Just(in).ToInt()
			return MonadIO.Just(v + 1)
		}).
		FlatMap(func(in interface{}) *MonadIODef[interface{}] {
			v, _ := Maybe.Just(in).ToInt()
			return MonadIO.Just(v * 2)
		}).
		FlatMap(func(in interface{}) *MonadIODef[interface{}] {
			v, _ := Maybe.Just(in).ToInt()
			return MonadIO.Just(v + 10)
		})

	result := m.Eval()
	assert.Equal(t, 14, result)
}

func TestMonadIOSubscribeMultipleTimes(t *testing.T) {
	m := MonadIO.Just(42)

	callCount := 0
	m.Subscribe(Subscription[interface{}]{
		OnNext: func(in interface{}) {
			callCount++
		},
	})
	m.Subscribe(Subscription[interface{}]{
		OnNext: func(in interface{}) {
			callCount++
		},
	})

	assert.Equal(t, 2, callCount)
}

func TestMonadIOWithNilValue(t *testing.T) {
	m := MonadIO.Just(nil)
	result := m.Eval()
	assert.Nil(t, result)
}

func TestMonadIOJustGenericsWithDifferentTypes(t *testing.T) {
	mInt := MonadIOJustGenerics(42)
	assert.Equal(t, 42, mInt.Eval())

	mStr := MonadIOJustGenerics("hello")
	assert.Equal(t, "hello", mStr.Eval())

	mBool := MonadIOJustGenerics(true)
	assert.Equal(t, true, mBool.Eval())
}

func TestMonadIONewGenericsWithDifferentTypes(t *testing.T) {
	mInt := MonadIONewGenerics(func() int { return 42 })
	assert.Equal(t, 42, mInt.Eval())

	mStr := MonadIONewGenerics(func() string { return "hello" })
	assert.Equal(t, "hello", mStr.Eval())

	mBool := MonadIONewGenerics(func() bool { return true })
	assert.Equal(t, true, mBool.Eval())
}

func TestMonadIOSubscribeOnWithHandler(t *testing.T) {
	h := new(HandlerDef).New()
	defer h.Close()

	m := MonadIO.Just(42).SubscribeOn(h)
	assert.NotNil(t, m)
}

func TestMonadIOObserveOnWithHandler(t *testing.T) {
	h := new(HandlerDef).New()
	defer h.Close()

	m := MonadIO.Just(42).ObserveOn(h)
	assert.NotNil(t, m)
}

func TestMonadIOFlatMapWithNil(t *testing.T) {
	m := MonadIO.Just(nil).FlatMap(func(in interface{}) *MonadIODef[interface{}] {
		return MonadIO.Just(in)
	})
	result := m.Eval()
	assert.Nil(t, result)
}

func TestMonadIOFlatMapChainedEval(t *testing.T) {
	callCount := 0
	m := MonadIO.New(func() interface{} {
		callCount++
		return 1
	}).FlatMap(func(in interface{}) *MonadIODef[interface{}] {
		v, _ := Maybe.Just(in).ToInt()
		return MonadIO.Just(v + 1)
	}).FlatMap(func(in interface{}) *MonadIODef[interface{}] {
		v, _ := Maybe.Just(in).ToInt()
		return MonadIO.Just(v * 2)
	})

	assert.Equal(t, 0, callCount)

	result := m.Eval()
	assert.Equal(t, 1, callCount)
	assert.Equal(t, 4, result)
}

func TestMonadIOSubscribeMultipleSubscriptions(t *testing.T) {
	results := []int{}
	m := MonadIO.Just(10)

	m.Subscribe(Subscription[interface{}]{
		OnNext: func(in interface{}) {
			v, _ := Maybe.Just(in).ToInt()
			results = append(results, v)
		},
	})

	m.Subscribe(Subscription[interface{}]{
		OnNext: func(in interface{}) {
			v, _ := Maybe.Just(in).ToInt()
			results = append(results, v*2)
		},
	})

	assert.Equal(t, []int{10, 20}, results)
}
