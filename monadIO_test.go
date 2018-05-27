package fpGo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonadIO(t *testing.T) {
	var m *MonadIODef
	var actualInt int

	m = MonadIO.JustVal(1)
	actualInt = 0
	m.Subscribe(Subscription{
		OnNext: func(in *interface{}) {
			actualInt, _ = Monad.Just(in).ToInt()
		},
	})
	assert.Equal(t, 1, actualInt)
}
