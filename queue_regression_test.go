package fpgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Regression: LinkedListQueue is documented as usable simultaneously as a
// Queue (Offer/Poll/Shift at head) and a Stack (Push/Pop at tail). Removing
// from one end previously left the surviving end node with a dangling link to
// the just-removed-and-recycled node. A subsequent removal from the other end
// then walked into that pooled node (Val == nil) and panicked, even though the
// queue was logically empty. These tests drive mixed end operations to empty
// and assert clean empty-errors instead of a panic or stale value.

func TestLinkedListQueuePopThenShiftDrainsCleanly(t *testing.T) {
	q := NewLinkedListQueue[int]()
	assert.NoError(t, q.Offer(1))
	assert.NoError(t, q.Offer(2))
	assert.NoError(t, q.Offer(3))

	v, err := q.Pop() // tail: 3
	assert.NoError(t, err)
	assert.Equal(t, 3, v)

	v, err = q.Shift() // head: 1
	assert.NoError(t, err)
	assert.Equal(t, 1, v)

	v, err = q.Shift() // head: 2 -> now empty
	assert.NoError(t, err)
	assert.Equal(t, 2, v)
	assert.Equal(t, 0, q.Count())

	// Must report empty, never panic or return a recycled value.
	_, err = q.Peek()
	assert.ErrorIs(t, err, ErrQueueIsEmpty)
	_, err = q.Shift()
	assert.ErrorIs(t, err, ErrQueueIsEmpty)
	_, err = q.Pop()
	assert.ErrorIs(t, err, ErrStackIsEmpty)
}

func TestLinkedListQueueShiftThenPopDrainsCleanly(t *testing.T) {
	q := NewLinkedListQueue[int]()
	assert.NoError(t, q.Offer(10))
	assert.NoError(t, q.Offer(20))
	assert.NoError(t, q.Offer(30))

	v, err := q.Shift() // head: 10
	assert.NoError(t, err)
	assert.Equal(t, 10, v)

	v, err = q.Pop() // tail: 30
	assert.NoError(t, err)
	assert.Equal(t, 30, v)

	v, err = q.Pop() // tail: 20 -> now empty
	assert.NoError(t, err)
	assert.Equal(t, 20, v)
	assert.Equal(t, 0, q.Count())

	_, err = q.Pop()
	assert.ErrorIs(t, err, ErrStackIsEmpty)
	_, err = q.Shift()
	assert.ErrorIs(t, err, ErrQueueIsEmpty)
}

// After draining and refilling, the recycled nodes must be reusable without
// carrying stale links (guards that the sever-on-removal fix cooperates with
// the node pool).
func TestLinkedListQueueRefillAfterMixedDrain(t *testing.T) {
	q := NewLinkedListQueue[int]()
	for _, v := range []int{1, 2, 3, 4} {
		assert.NoError(t, q.Offer(v))
	}
	// Drain via alternating ends.
	_, _ = q.Pop()   // 4
	_, _ = q.Shift() // 1
	_, _ = q.Pop()   // 3
	_, _ = q.Shift() // 2
	assert.Equal(t, 0, q.Count())

	// Refill and verify FIFO integrity end-to-end.
	for _, v := range []int{5, 6, 7} {
		assert.NoError(t, q.Offer(v))
	}
	assert.Equal(t, 3, q.Count())
	got := make([]int, 0, 3)
	for q.Count() > 0 {
		v, err := q.Shift()
		assert.NoError(t, err)
		got = append(got, v)
	}
	assert.Equal(t, []int{5, 6, 7}, got)
}

// Unshift (add head) + Pop (remove tail) mixed drain must also stay clean.
func TestLinkedListQueueUnshiftPopMixed(t *testing.T) {
	q := NewLinkedListQueue[int]()
	assert.NoError(t, q.Unshift(2)) // [2]
	assert.NoError(t, q.Unshift(1)) // [1,2]
	assert.NoError(t, q.Push(3))    // [1,2,3]

	v, err := q.Pop() // 3
	assert.NoError(t, err)
	assert.Equal(t, 3, v)
	v, err = q.Pop() // 2
	assert.NoError(t, err)
	assert.Equal(t, 2, v)
	v, err = q.Shift() // 1 -> empty
	assert.NoError(t, err)
	assert.Equal(t, 1, v)
	assert.Equal(t, 0, q.Count())

	_, err = q.Shift()
	assert.ErrorIs(t, err, ErrQueueIsEmpty)
}
