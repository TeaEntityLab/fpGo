package fpGo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAsync(t *testing.T) {
	var m MonadProto

	m = Monad.AsyncVal(1)
	assert.Equal(t, true, m.IsPresent())
	assert.Equal(t, 1, m.Unwrap())
}
