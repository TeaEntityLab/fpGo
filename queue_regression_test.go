package fpgo

import (
	"testing"
	"time"

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

// Regression: ConcurrentQueue.Count() only delegated to a Count() method and
// returned 0 for anything else. ChannelQueue is a buffered channel with no
// Count() method, so a ConcurrentQueue wrapping one reported 0 items while it
// actually held them. Count must reflect the channel's buffered length.
func TestConcurrentQueueCountOnChannelQueue(t *testing.T) {
	cq := NewConcurrentQueue[int](NewChannelQueue[int](4))

	assert.Equal(t, 0, cq.Count(), "empty channel queue")
	assert.NoError(t, cq.Offer(1))
	assert.NoError(t, cq.Offer(2))
	assert.Equal(t, 2, cq.Count(), "count must track buffered channel items")

	v, err := cq.Poll()
	assert.NoError(t, err)
	assert.Equal(t, 1, v)
	assert.Equal(t, 1, cq.Count(), "count must drop after a poll")
}

// Regression: ChannelQueue.Poll used `case val := <-q` without the ok flag, so
// a closed channel yielded the zero value with a nil error — indistinguishable
// from polling a real zero. Take/TakeWithTimeout already report
// ErrQueueIsClosed; Poll must do the same.
func TestChannelQueuePollAfterClose(t *testing.T) {
	queue := NewChannelQueue[int](3)
	assert.NoError(t, queue.Offer(1))
	close(queue)

	val, err := queue.Poll()
	assert.NoError(t, err)
	assert.Equal(t, 1, val)

	val, err = queue.Poll()
	assert.Equal(t, ErrQueueIsClosed, err)
	assert.Equal(t, 0, val)
}

// Regression: Put, PutWithTimeout, and Offer sent straight onto the channel.
// ChannelQueue is a named chan, so callers can close it, and a send on a
// closed channel panics — including inside select, so Offer's default arm did
// not help. Receives already return ErrQueueIsClosed; sends must too.
func TestChannelQueueSendAfterClose(t *testing.T) {
	queue := NewChannelQueue[int](1)
	close(queue)

	assert.Equal(t, ErrQueueIsClosed, queue.Put(1))
	assert.Equal(t, ErrQueueIsClosed, queue.Offer(1))
	assert.Equal(t, ErrQueueIsClosed, queue.PutWithTimeout(1, time.Millisecond))

	// A closed full channel must also report closed, not full.
	full := NewChannelQueue[int](1)
	assert.NoError(t, full.Offer(1))
	close(full)
	assert.Equal(t, ErrQueueIsClosed, full.Offer(2))
}

// Regression: BufferedChannelQueue.Poll only polled the channel, so it could
// return ErrQueueIsEmpty while items sat in the overflow pool waiting for the
// interval-driven loader to move them. Poll must drain the pool too.
func TestBufferedChannelQueuePollDrainsOverflowPool(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 10, 4)
	defer q.Close()

	// Fill the channel first, then overflow into the pool.
	assert.NoError(t, q.Offer(1))
	assert.NoError(t, q.Offer(2))
	assert.NoError(t, q.Offer(3))

	// Drain synchronously via Poll — no waiting for the loader goroutine.
	got := []int{}
	for i := 0; i < 3; i++ {
		v, err := q.Poll()
		assert.NoError(t, err)
		got = append(got, v)
	}
	assert.Equal(t, []int{1, 2, 3}, got, "Poll must return pool items in FIFO order")

	_, err := q.Poll()
	assert.Equal(t, ErrQueueIsEmpty, err)
}

// Regression: ConcurrentQueue held its write lock across blocking Take. With a
// self-synchronizing inner queue (ChannelQueue), a consumer blocked in Take
// held the lock forever, so a producer Offer could never proceed — deadlock.
// Self-synchronizing queues now bypass the redundant wrapper lock.
func TestConcurrentQueueChannelBackedProducerConsumer(t *testing.T) {
	cq := NewConcurrentQueue[int](NewChannelQueue[int](4))
	done := make(chan int, 1)

	go func() {
		v, err := cq.Take() // blocks until producer offers
		assert.NoError(t, err)
		done <- v
	}()

	// Give the consumer time to block inside Take.
	time.Sleep(20 * time.Millisecond)
	assert.NoError(t, cq.Offer(42))

	select {
	case v := <-done:
		assert.Equal(t, 42, v)
	case <-time.After(2 * time.Second):
		t.Fatal("deadlock: producer Offer blocked behind consumer Take")
	}
}

// Coverage for the remaining self-synchronizing bypass arms and zero-value
// safety.
func TestConcurrentQueueChannelBackedPutAndPoolNodeCount(t *testing.T) {
	cq := NewConcurrentQueue[int](NewChannelQueue[int](4))
	assert.NoError(t, cq.Put(7))
	v, err := cq.Poll()
	assert.NoError(t, err)
	assert.Equal(t, 7, v)

	// Zero-value BufferedChannelQueue: nil pool must not panic.
	var zero BufferedChannelQueue[int]
	assert.Equal(t, 0, zero.PoolNodeCount())
}
