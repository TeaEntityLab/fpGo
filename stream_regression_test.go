package fpgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Regression: StreamDef.Remove returns a NEW *StreamDef (like Filter/Map/
// Reverse), so it must be non-destructive. It previously did
// `append((*streamSelf)[:index], (*streamSelf)[index+1:]...)`, which reuses
// the receiver's backing array: the returned slice aliased the receiver AND
// the receiver's tail was corrupted (last element duplicated). These tests
// assert the receiver is untouched and the result does not alias it.

func TestStreamRemoveDoesNotMutateReceiver(t *testing.T) {
	s := StreamFrom(1, 2, 3, 4, 5)
	result := s.Remove(1) // drop value 2

	assert.Equal(t, []int{1, 3, 4, 5}, result.ToArray(), "result content")
	// Receiver must be unchanged (previously became [1 3 4 5 5]).
	assert.Equal(t, []int{1, 2, 3, 4, 5}, s.ToArray(), "receiver must be intact")
}

func TestStreamRemoveResultDoesNotAliasReceiver(t *testing.T) {
	s := StreamFrom(1, 2, 3, 4)
	result := s.Remove(1) // [1,3,4]
	// Mutating the result must not reach back into the receiver.
	(*result)[0] = 999
	assert.Equal(t, []int{1, 2, 3, 4}, s.ToArray(), "receiver unaffected by result mutation")
	assert.Equal(t, []int{999, 3, 4}, result.ToArray())
}

func TestStreamRemoveHeadTailMiddle(t *testing.T) {
	base := StreamFrom(10, 20, 30, 40)

	head := base.Remove(0)
	assert.Equal(t, []int{20, 30, 40}, head.ToArray())

	tail := base.Remove(3)
	assert.Equal(t, []int{10, 20, 30}, tail.ToArray())

	mid := base.Remove(2)
	assert.Equal(t, []int{10, 20, 40}, mid.ToArray())

	// Receiver still intact after several Removes.
	assert.Equal(t, []int{10, 20, 30, 40}, base.ToArray())
}

func TestStreamRemoveOutOfRangeReturnsReceiver(t *testing.T) {
	s := StreamFrom(1, 2, 3)
	assert.Same(t, s, s.Remove(-1))
	assert.Same(t, s, s.Remove(3))
	assert.Same(t, s, s.Remove(99))
	assert.Equal(t, []int{1, 2, 3}, s.ToArray())
}

// Chained Removes should compose correctly without corruption.
func TestStreamRemoveChained(t *testing.T) {
	s := StreamFrom(1, 2, 3, 4, 5)
	got := s.Remove(0).Remove(0).ToArray() // drop 1, then drop 2
	assert.Equal(t, []int{3, 4, 5}, got)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, s.ToArray(), "original intact")
}
