package fpgo

import (
	"testing"

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
}
