package fpgo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMonadIO(t *testing.T) {
	var m *MonadIODef
	var actualInt int

	m = MonadIO.Just(1)
	actualInt = 0
	m.Subscribe(Subscription{
		OnNext: func(in interface{}) {
			actualInt, _ = Maybe.Just(in).ToInt()
		},
	})
	assert.Equal(t, 1, actualInt)

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

func TestMonadIOSubscribeOnWithHandler(t *testing.T) {
	h := new(HandlerDef).New()
	defer h.Close()

	m := MonadIO.Just(42).SubscribeOn(h)
	assert.NotNil(t, m)
}

func TestMonadIOSubscribeWithHandlers(t *testing.T) {
	hSub := new(HandlerDef).New()
	defer hSub.Close()

	m := MonadIO.Just(1).
		SubscribeOn(hSub).
		ObserveOn(Handler.GetDefault())

	done := make(chan struct{})
	var result int
	m.Subscribe(Subscription{
		OnNext: func(in interface{}) {
			result, _ = Maybe.Just(in).ToInt()
			close(done)
		},
	})

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for async handler")
	}
	assert.Equal(t, 1, result)
}

func TestMonadIOSubscribeOnNil(t *testing.T) {
	m := MonadIO.Just(1).SubscribeOn(nil)

	var result int
	m.Subscribe(Subscription{
		OnNext: func(in interface{}) {
			result, _ = Maybe.Just(in).ToInt()
		},
	})

	assert.Equal(t, 1, result)
}
