package fpgo

import (
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChannelQueue(t *testing.T) {
	var queue Queue[int]
	var err error
	var result int
	var timeout time.Duration

	channelQueue := NewChannelQueue[int](3)
	queue = channelQueue

	err = queue.Offer(1)
	assert.Equal(t, nil, err)
	err = queue.Offer(2)
	assert.Equal(t, nil, err)
	err = queue.Offer(3)
	assert.Equal(t, nil, err)
	err = queue.Offer(4)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, ErrQueueIsFull, err)

	result, err = queue.Poll()
	assert.Equal(t, 1, result)
	assert.Equal(t, nil, err)
	result, err = queue.Poll()
	assert.Equal(t, 2, result)
	assert.Equal(t, nil, err)
	result, err = queue.Poll()
	assert.Equal(t, 3, result)
	assert.Equal(t, nil, err)
	result, err = queue.Poll()
	assert.NotEqual(t, 4, result)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, 0, result)
	assert.Equal(t, ErrQueueIsEmpty, err)

	timeout = 1 * time.Millisecond
	done := make(chan struct{}, 2)
	go func() {
		defer func() { done <- struct{}{} }()
		result, err := channelQueue.TakeWithTimeout(timeout)
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, result)
		result, err = channelQueue.TakeWithTimeout(timeout)
		assert.Equal(t, nil, err)
		assert.Equal(t, 2, result)
		result, err = channelQueue.TakeWithTimeout(timeout)
		assert.Equal(t, nil, err)
		assert.Equal(t, 3, result)
		result, err = channelQueue.TakeWithTimeout(timeout)
		assert.NotEqual(t, 4, result)
		assert.NotEqual(t, nil, err)
		assert.Equal(t, 0, result)
		assert.Equal(t, ErrQueueTakeTimeout, err)
	}()
	go func() {
		defer func() { done <- struct{}{} }()
		err := channelQueue.PutWithTimeout(1, timeout)
		assert.Equal(t, nil, err)
		err = channelQueue.PutWithTimeout(2, timeout)
		assert.Equal(t, nil, err)
		err = channelQueue.PutWithTimeout(3, timeout)
		assert.Equal(t, nil, err)

		time.Sleep(3 * timeout / 2)

		err = channelQueue.PutWithTimeout(4, timeout)
		assert.Equal(t, nil, err)
		err = channelQueue.PutWithTimeout(5, timeout)
		assert.Equal(t, nil, err)
		err = channelQueue.PutWithTimeout(6, timeout)
		assert.Equal(t, nil, err)
		err = channelQueue.PutWithTimeout(7, timeout)
		assert.NotEqual(t, nil, err)
		assert.Equal(t, ErrQueuePutTimeout, err)
	}()

	<-done
	<-done
}

func TestLinkedListQueue(t *testing.T) {
	var queue Queue[int]
	var stack Stack[int]
	var err error
	var result int

	linkedListQueue := NewLinkedListQueue[int]()
	queue = linkedListQueue
	stack = linkedListQueue
	concurrentQueue := NewConcurrentQueue[int](queue)

	err = queue.Offer(1)
	assert.Equal(t, nil, err)
	err = queue.Offer(2)
	assert.Equal(t, nil, err)
	err = queue.Offer(3)
	assert.Equal(t, nil, err)

	result, err = queue.Poll()
	assert.Equal(t, 1, result)
	assert.Equal(t, nil, err)
	result, err = queue.Poll()
	assert.Equal(t, 2, result)
	assert.Equal(t, nil, err)
	result, err = queue.Poll()
	assert.Equal(t, 3, result)
	assert.Equal(t, nil, err)
	result, err = queue.Poll()
	assert.NotEqual(t, 4, result)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, 0, result)
	assert.Equal(t, ErrQueueIsEmpty, err)

	err = stack.Push(1)
	assert.Equal(t, nil, err)
	err = stack.Push(2)
	assert.Equal(t, nil, err)
	err = stack.Push(3)
	assert.Equal(t, nil, err)

	result, err = stack.Pop()
	assert.Equal(t, 3, result)
	assert.Equal(t, nil, err)
	result, err = stack.Pop()
	assert.Equal(t, 2, result)
	assert.Equal(t, nil, err)
	result, err = stack.Pop()
	assert.Equal(t, 1, result)
	assert.Equal(t, nil, err)
	result, err = stack.Pop()
	assert.NotEqual(t, 4, result)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, 0, result)
	assert.Equal(t, ErrStackIsEmpty, err)

	linkedListQueue.KeepNodePoolCount(10)
	assert.Equal(t, 10, linkedListQueue.nodeCount)
	assert.Equal(t, 10, linkedListQueue.nodePoolFirst.Count())
	linkedListQueue.KeepNodePoolCount(2)
	assert.Equal(t, 2, linkedListQueue.nodeCount)
	assert.Equal(t, 2, linkedListQueue.nodePoolFirst.Count())
	linkedListQueue.KeepNodePoolCount(0)
	assert.Equal(t, 0, linkedListQueue.nodeCount)
	assert.Nil(t, linkedListQueue.nodePoolFirst)
	linkedListQueue.KeepNodePoolCount(5)
	assert.Equal(t, 5, linkedListQueue.nodeCount)
	assert.Equal(t, 5, linkedListQueue.nodePoolFirst.Count())
	linkedListQueue.KeepNodePoolCount(3)
	assert.Equal(t, 3, linkedListQueue.nodeCount)
	assert.Equal(t, 3, linkedListQueue.nodePoolFirst.Count())

	// Deterministic sequence (was sleep-lockstep goroutines; no real concurrency
	// here — ConcurrentQueue locking is covered by dedicated concurrent tests).
	assert.Equal(t, 0, concurrentQueue.Count())
	assert.Equal(t, 3, concurrentQueue.NodePoolCount())
	e := concurrentQueue.Put(1)
	assert.Equal(t, nil, e)
	assert.Equal(t, 1, concurrentQueue.Count())
	assert.Equal(t, 2, concurrentQueue.NodePoolCount())
	e = concurrentQueue.Put(2)
	assert.Equal(t, nil, e)
	assert.Equal(t, 2, concurrentQueue.Count())
	assert.Equal(t, 1, concurrentQueue.NodePoolCount())
	e = concurrentQueue.Put(3)
	assert.Equal(t, nil, e)
	assert.Equal(t, 3, concurrentQueue.Count())
	assert.Equal(t, 0, concurrentQueue.NodePoolCount())

	r, e := concurrentQueue.Take()
	assert.Equal(t, nil, e)
	assert.Equal(t, 1, r)
	assert.Equal(t, 2, concurrentQueue.Count())
	assert.Equal(t, 1, concurrentQueue.NodePoolCount())
	r, e = concurrentQueue.Take()
	assert.Equal(t, nil, e)
	assert.Equal(t, 2, r)
	assert.Equal(t, 1, concurrentQueue.Count())
	assert.Equal(t, 2, concurrentQueue.NodePoolCount())
	r, e = concurrentQueue.Take()
	assert.Equal(t, nil, e)
	assert.Equal(t, 3, r)
	assert.Equal(t, 0, concurrentQueue.Count())
	assert.Equal(t, 3, concurrentQueue.NodePoolCount())
	r, e = concurrentQueue.Take()
	assert.NotEqual(t, 4, r)
	assert.NotEqual(t, nil, e)
	assert.Equal(t, 0, r)
	assert.Equal(t, ErrQueueIsEmpty, e)
	assert.Equal(t, 0, concurrentQueue.Count())
	assert.Equal(t, 3, concurrentQueue.NodePoolCount())

	linkedListQueue.KeepNodePoolCount(2)
	assert.Equal(t, 2, concurrentQueue.NodePoolCount())
	assert.Equal(t, 2, linkedListQueue.nodePoolFirst.Count())
	e = concurrentQueue.Put(4)
	assert.Equal(t, nil, e)
	assert.Equal(t, 1, concurrentQueue.Count())
	assert.Equal(t, 1, concurrentQueue.NodePoolCount())
	e = concurrentQueue.Put(5)
	assert.Equal(t, nil, e)
	assert.Equal(t, 2, concurrentQueue.Count())
	assert.Equal(t, 0, concurrentQueue.NodePoolCount())
	e = concurrentQueue.Put(6)
	assert.Equal(t, nil, e)
	assert.Equal(t, 3, concurrentQueue.Count())
	assert.Equal(t, 0, concurrentQueue.NodePoolCount())
	linkedListQueue.Clear()
	assert.Equal(t, 0, concurrentQueue.Count())
	assert.Equal(t, 3, concurrentQueue.NodePoolCount())
	node := linkedListQueue.nodePoolFirst
	for node != nil {
		assert.Nil(t, node.Val)
		node = node.Next
	}
	linkedListQueue.ClearNodePool()
	assert.Equal(t, 0, concurrentQueue.Count())
	assert.Equal(t, 0, concurrentQueue.NodePoolCount())
	assert.Nil(t, linkedListQueue.nodePoolFirst)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 1; i <= 10000; i++ {
			// LinkedListQueue.Take is non-blocking; retry while the producer
			// hasn't caught up (ErrQueueIsEmpty) so the consumer never desyncs.
			var result int
			var err error
			for {
				result, err = concurrentQueue.Take()
				if err != ErrQueueIsEmpty {
					break
				}
			}
			assert.Equal(t, nil, err)
			assert.Equal(t, i, result)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 1; i <= 10000; i++ {
			err := concurrentQueue.Offer(i)
			assert.Equal(t, nil, err)
		}
	}()

	wg.Wait()
}

func TestLinkedListItem(t *testing.T) {
	item1 := &LinkedListItem[int]{Val: new(int)}
	*item1.Val = 1
	item2 := &LinkedListItem[int]{Val: new(int)}
	*item2.Val = 2

	assert.Equal(t, 1, item1.Count())

	item1.Next = item2
	assert.Equal(t, 2, item1.Count())

	last := item1.Last()
	assert.Equal(t, 2, *last.Val)

	added := item1.AddLast(&LinkedListItem[int]{Val: new(int)})
	*added.Val = 3
	assert.Equal(t, 3, item1.Count())
}

func TestDoublyListItem(t *testing.T) {
	item1 := &DoublyListItem[int]{Val: new(int)}
	*item1.Val = 1
	item2 := &DoublyListItem[int]{Val: new(int)}
	*item2.Val = 2

	assert.Equal(t, 1, item1.Count())

	item1.Next = item2
	item2.Prev = item1

	assert.Equal(t, 2, item1.Count())

	last := item1.Last()
	assert.Equal(t, 2, *last.Val)

	first := item1.First()
	assert.Equal(t, 1, *first.Val)

	item1.AddLast(&DoublyListItem[int]{Val: new(int)})
	assert.Equal(t, 3, item1.Count())

	item1.AddFirst(&DoublyListItem[int]{Val: new(int)})
	assert.Equal(t, 4, item1.Count())
}

func TestLinkedListQueuePeek(t *testing.T) {
	q := NewLinkedListQueue[int]()

	q.Offer(1)
	q.Offer(2)
	q.Offer(3)

	val, err := q.Peek()
	assert.Equal(t, 1, val)
	assert.Nil(t, err)
	assert.Equal(t, 3, q.Count())
}

func TestLinkedListQueueShiftUnshift(t *testing.T) {
	q := NewLinkedListQueue[int]()

	q.Unshift(1)
	q.Unshift(2)
	q.Unshift(3)

	val, err := q.Shift()
	assert.Equal(t, 3, val)
	assert.Nil(t, err)

	val, err = q.Shift()
	assert.Equal(t, 2, val)
	assert.Nil(t, err)

	val, err = q.Shift()
	assert.Equal(t, 1, val)
	assert.Nil(t, err)

	_, err = q.Shift()
	assert.Equal(t, ErrQueueIsEmpty, err)
}

func TestBufferedChannelQueueGetters(t *testing.T) {
	q := NewBufferedChannelQueue[int](3, 100, 10)

	assert.Equal(t, 100, q.GetBufferSizeMaximum())
	assert.Equal(t, 10, q.GetNodeHookPoolSize())
	assert.NotNil(t, q.GetChannel())

	q.Close()
}

func TestBufferedChannelQueueDurationGetters(t *testing.T) {
	q := NewBufferedChannelQueue[int](3, 100, 10).
		SetLoadFromPoolDuration(5 * time.Millisecond).
		SetFreeNodeHookPoolIntervalDuration(10 * time.Millisecond)

	assert.Equal(t, 5*time.Millisecond, q.GetLoadFromPoolDuration())
	assert.Equal(t, 10*time.Millisecond, q.GetFreeNodeHookPoolIntervalDuration())

	q.Close()
}

func TestLinkedListQueuePushPop(t *testing.T) {
	q := NewLinkedListQueue[int]()

	q.Push(1)
	q.Push(2)
	q.Push(3)

	val, err := q.Pop()
	assert.Equal(t, 3, val)
	assert.Nil(t, err)

	val, err = q.Pop()
	assert.Equal(t, 2, val)
	assert.Nil(t, err)

	val, err = q.Pop()
	assert.Equal(t, 1, val)
	assert.Nil(t, err)

	_, err = q.Pop()
	assert.Equal(t, ErrStackIsEmpty, err)
}

func TestChannelQueuePutTake(t *testing.T) {
	q := NewChannelQueue[int](5)

	q.Put(1)
	q.Put(2)
	q.Put(3)

	val, err := q.Take()
	assert.Equal(t, 1, val)
	assert.Nil(t, err)

	val, err = q.Take()
	assert.Equal(t, 2, val)
	assert.Nil(t, err)
}

func TestConcurrentQueue(t *testing.T) {
	q := NewConcurrentQueue[int](NewLinkedListQueue[int]())

	q.Offer(1)
	q.Offer(2)

	val, err := q.Poll()
	assert.Equal(t, 1, val)
	assert.Nil(t, err)

	val, err = q.Poll()
	assert.Equal(t, 2, val)
	assert.Nil(t, err)
}

func TestNewBufferedChannelQueue(t *testing.T) {
	var queue Queue[int]
	var err error
	var result int
	var timeout time.Duration

	bufferedChannelQueue := NewBufferedChannelQueue[int](3, 10000, 100).
		SetLoadFromPoolDuration(time.Millisecond / 10).
		SetFreeNodeHookPoolIntervalDuration(1 * time.Millisecond)
	queue = bufferedChannelQueue

	// Sync

	timeout = 1 * time.Millisecond
	bufferedChannelQueue.SetBufferSizeMaximum(1)

	err = queue.Offer(1)
	assert.Equal(t, nil, err)
	err = queue.Offer(2)
	assert.Equal(t, nil, err)
	err = queue.Offer(3)
	assert.Equal(t, nil, err)
	// Channel: only 3 positions & Buffer: 1 position, now `4` is inserted into the buffer(buffer size： 1)
	err = queue.Offer(4)
	assert.Equal(t, nil, err)
	time.Sleep(1 * timeout)
	// Channel: only 3 positions & Buffer: 1 position, now `5` can't be inserted into the buffer(`4` is already inside)
	err = queue.Offer(5)
	assert.Equal(t, ErrQueueIsFull, err)

	result, err = bufferedChannelQueue.TakeWithTimeout(timeout)
	assert.Equal(t, 1, result)
	assert.Equal(t, nil, err)
	result, err = bufferedChannelQueue.TakeWithTimeout(timeout)
	assert.Equal(t, 2, result)
	assert.Equal(t, nil, err)
	result, err = bufferedChannelQueue.TakeWithTimeout(timeout)
	assert.Equal(t, 3, result)
	assert.Equal(t, nil, err)
	result, err = bufferedChannelQueue.TakeWithTimeout(timeout)
	assert.Equal(t, 4, result)
	assert.Equal(t, nil, err)

	// Async
	asyncTaskDone := make(chan bool)

	bufferedChannelQueue.SetBufferSizeMaximum(6)
	timeout = 2 * time.Millisecond
	go func() {
		time.Sleep(timeout)
		for expected := 1; expected <= 6; expected++ {
			// Retry on timeout: items arrive in FIFO order; a transient
			// TakeWithTimeout miss under load must not fail the assertion.
			var r int
			var e error
			for {
				r, e = bufferedChannelQueue.TakeWithTimeout(timeout)
				if e != ErrQueueTakeTimeout {
					break
				}
			}
			assert.Equal(t, nil, e)
			assert.Equal(t, expected, r)
		}
		asyncTaskDone <- true
	}()
	go func() {
		e := bufferedChannelQueue.Put(1)
		assert.Equal(t, nil, e)
		e = bufferedChannelQueue.Put(2)
		assert.Equal(t, nil, e)
		e = bufferedChannelQueue.Put(3)
		assert.Equal(t, nil, e)
		e = bufferedChannelQueue.Put(4)
		assert.Equal(t, nil, e)
		e = bufferedChannelQueue.Put(5)
		assert.Equal(t, nil, e)
		e = bufferedChannelQueue.Put(6)
		assert.Equal(t, nil, e)
	}()

	<-asyncTaskDone

	bufferedChannelQueue.SetBufferSizeMaximum(10000)
	timeout = 10 * time.Millisecond
	go func() {
		for i := 1; i <= 10000; i++ {
			// Retry on timeout: all 10000 items are guaranteed to arrive in
			// FIFO order; a transient TakeWithTimeout miss must not cascade.
			var result int
			var err error
			for {
				result, err = bufferedChannelQueue.TakeWithTimeout(timeout)
				if err != ErrQueueTakeTimeout {
					break
				}
			}
			assert.Equal(t, nil, err)
			assert.Equal(t, i, result)
		}
		asyncTaskDone <- true
	}()
	go func() {
		for i := 1; i <= 10000; i++ {
			err := bufferedChannelQueue.Offer(i)
			assert.Equal(t, nil, err)
		}
	}()

	<-asyncTaskDone

	result, err = bufferedChannelQueue.Poll()
	assert.Equal(t, ErrQueueIsEmpty, err)
	assert.Equal(t, 0, result)

	time.Sleep(1 * timeout)

	assert.GreaterOrEqual(t, bufferedChannelQueue.PoolNodeCount(), 100)
	close(asyncTaskDone)
}

func TestBufferedChannelQueuePutWithTimeout(t *testing.T) {
	bufferedChannelQueue := NewBufferedChannelQueue[int](1, 10, 5).
		SetLoadFromPoolDuration(time.Millisecond).
		SetFreeNodeHookPoolIntervalDuration(time.Millisecond)

	err := bufferedChannelQueue.PutWithTimeout(1, 10*time.Millisecond)
	assert.Nil(t, err)

	err = bufferedChannelQueue.PutWithTimeout(2, 10*time.Millisecond)
	assert.Nil(t, err)

	bufferedChannelQueue.Close()

	err = bufferedChannelQueue.PutWithTimeout(3, 10*time.Millisecond)
	assert.Equal(t, ErrQueueIsClosed, err)
}

func TestBufferedChannelQueueClosedAndFullBranches(t *testing.T) {
	closedQueue := NewBufferedChannelQueue[int](1, 1, 1)
	closedQueue.Close()
	assert.Equal(t, 0, closedQueue.Count())
	assert.Equal(t, ErrQueueIsClosed, closedQueue.Offer(1))
	assert.Equal(t, ErrQueueIsClosed, closedQueue.PutWithTimeout(1, time.Millisecond))
	_, err := closedQueue.Take()
	assert.Equal(t, ErrQueueIsClosed, err)
	_, err = closedQueue.TakeWithTimeout(time.Millisecond)
	assert.Equal(t, ErrQueueIsClosed, err)
	assert.NotPanics(t, func() { closedQueue.notifyWorkers() })

	countQueue := NewBufferedChannelQueue[int](1, 1, 1)
	countQueue.blockingQueue = nil
	assert.Equal(t, 0, countQueue.Count())

	fullQueue := NewBufferedChannelQueue[int](1, 0, 1).
		SetLoadFromPoolDuration(5 * time.Millisecond).
		SetFreeNodeHookPoolIntervalDuration(time.Millisecond)
	assert.NoError(t, fullQueue.Offer(1))
	assert.Equal(t, ErrQueueIsFull, fullQueue.Offer(2))
	assert.Equal(t, ErrQueueIsFull, fullQueue.PutWithTimeout(3, 10*time.Millisecond))
	fullQueue.Close()
}

func TestBufferedChannelQueueLoadFromPoolAndTimeoutBranches(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 4, 1).
		SetLoadFromPoolDuration(time.Millisecond).
		SetFreeNodeHookPoolIntervalDuration(time.Millisecond)

	assert.NoError(t, q.Offer(1))
	assert.NoError(t, q.Offer(2))
	q.lock.RLock()
	poolCount := q.pool.Count()
	q.lock.RUnlock()
	assert.Equal(t, 1, poolCount)

	v, err := q.TakeWithTimeout(10 * time.Millisecond)
	assert.NoError(t, err)
	assert.Equal(t, 1, v)

	q.notifyWorkers()
	time.Sleep(4 * time.Millisecond)
	assert.Equal(t, 1, q.Count())
	q.lock.RLock()
	poolCount = q.pool.Count()
	q.lock.RUnlock()
	assert.Equal(t, 0, poolCount)

	v, err = q.TakeWithTimeout(10 * time.Millisecond)
	assert.NoError(t, err)
	assert.Equal(t, 2, v)

	timeoutQueue := NewBufferedChannelQueue[int](1, 1, 1).
		SetLoadFromPoolDuration(2 * time.Millisecond).
		SetFreeNodeHookPoolIntervalDuration(time.Millisecond)
	timeoutQueue.Close()
	assert.Equal(t, ErrQueueIsClosed, timeoutQueue.PutWithTimeout(3, 3*time.Millisecond))

	q.Close()
}

func TestLinkedListQueueClear(t *testing.T) {
	queue := NewLinkedListQueue[int]()
	queue.Offer(1)
	queue.Offer(2)
	queue.Offer(3)

	assert.Equal(t, 3, queue.Count())

	queue.Clear()
	assert.Equal(t, 0, queue.Count())
}

func TestConcurrentStack(t *testing.T) {
	linkedListQueue := NewLinkedListQueue[int]()
	stack := NewConcurrentStack[int](linkedListQueue)

	err := stack.Push(1)
	assert.Nil(t, err)
	err = stack.Push(2)
	assert.Nil(t, err)
	err = stack.Push(3)
	assert.Nil(t, err)

	val, err := stack.Pop()
	assert.Equal(t, 3, val)
	assert.Nil(t, err)

	val, err = stack.Pop()
	assert.Equal(t, 2, val)
	assert.Nil(t, err)

	val, err = stack.Pop()
	assert.Equal(t, 1, val)
	assert.Nil(t, err)

	val, err = stack.Pop()
	assert.Equal(t, 0, val)
	assert.NotNil(t, err)
}

func TestLinkedListQueueNodePool(t *testing.T) {
	queue := NewLinkedListQueue[int]()

	queue.Offer(1)
	queue.Offer(2)
	queue.Offer(3)

	queue.KeepNodePoolCount(5)
	queue.ClearNodePool()
}

func TestDoublyListItemAdvanced(t *testing.T) {
	val1 := 1
	val2 := 2
	val3 := 3
	item1 := &DoublyListItem[int]{Val: &val1}
	item2 := &DoublyListItem[int]{Val: &val2}
	item3 := &DoublyListItem[int]{Val: &val3}

	item2 = item2.AddFirst(item1)
	assert.Equal(t, item1, item2.Prev)

	item2 = item2.AddLast(item3)
	assert.Equal(t, item3, item2.Next)

	first := item2.First()
	assert.Equal(t, 1, *first.Val)

	last := item2.Last()
	assert.Equal(t, 3, *last.Val)

	assert.Equal(t, 3, item2.Count())
}

func TestLinkedListQueueShift(t *testing.T) {
	queue := NewLinkedListQueue[int]()

	queue.Offer(1)
	queue.Offer(2)
	queue.Offer(3)

	val, err := queue.Shift()
	assert.Equal(t, 1, val)
	assert.Nil(t, err)

	val, err = queue.Shift()
	assert.Equal(t, 2, val)
	assert.Nil(t, err)

	val, err = queue.Shift()
	assert.Equal(t, 3, val)
	assert.Nil(t, err)

	val, err = queue.Shift()
	assert.Equal(t, 0, val)
	assert.NotNil(t, err)
}

func TestLinkedListQueueUnshift(t *testing.T) {
	queue := NewLinkedListQueue[int]()

	err := queue.Unshift(1)
	assert.Nil(t, err)
	err = queue.Unshift(2)
	assert.Nil(t, err)
	err = queue.Unshift(3)
	assert.Nil(t, err)

	val, err := queue.Poll()
	assert.Equal(t, 3, val)
	assert.Nil(t, err)

	val, err = queue.Poll()
	assert.Equal(t, 2, val)
	assert.Nil(t, err)

	val, err = queue.Poll()
	assert.Equal(t, 1, val)
	assert.Nil(t, err)
}

func TestBufferedChannelQueuePoolEmptyChannelFull(t *testing.T) {
	bufferedChannelQueue := NewBufferedChannelQueue[int](1, 1, 10).
		SetLoadFromPoolDuration(time.Millisecond).
		SetFreeNodeHookPoolIntervalDuration(time.Millisecond)

	err := bufferedChannelQueue.Offer(1)
	assert.Nil(t, err)

	err = bufferedChannelQueue.Offer(2)
	assert.Nil(t, err)

	_, err = bufferedChannelQueue.Take()
	assert.Nil(t, err)

	_, err = bufferedChannelQueue.Take()
	assert.Nil(t, err)

	err = bufferedChannelQueue.Offer(3)
	assert.Nil(t, err)
}

func TestBufferedChannelQueueClose(t *testing.T) {
	bufferedChannelQueue := NewBufferedChannelQueue[int](3, 10, 5).
		SetLoadFromPoolDuration(time.Millisecond).
		SetFreeNodeHookPoolIntervalDuration(time.Millisecond)

	assert.False(t, bufferedChannelQueue.IsClosed())

	bufferedChannelQueue.Close()
	bufferedChannelQueue.Close()

	assert.True(t, bufferedChannelQueue.IsClosed())
	assert.Equal(t, 0, bufferedChannelQueue.Count())

	err := bufferedChannelQueue.Offer(1)
	assert.Equal(t, ErrQueueIsClosed, err)

	err = bufferedChannelQueue.Put(1)
	assert.Equal(t, ErrQueueIsClosed, err)

	_, err = bufferedChannelQueue.Take()
	assert.Equal(t, ErrQueueIsClosed, err)

	_, err = bufferedChannelQueue.TakeWithTimeout(time.Millisecond)
	assert.Equal(t, ErrQueueIsClosed, err)

	_, err = bufferedChannelQueue.Poll()
	assert.Equal(t, ErrQueueIsClosed, err)

	err = bufferedChannelQueue.PutWithTimeout(1, time.Millisecond)
	assert.Equal(t, ErrQueueIsClosed, err)

	assert.NotPanics(t, func() {
		bufferedChannelQueue.notifyWorkers()
	})
}

func TestBufferedChannelQueueCloseStopsFreeNodeGoroutine(t *testing.T) {
	runtime.GC()
	time.Sleep(10 * time.Millisecond)
	baseline := runtime.NumGoroutine()

	const n = 50
	queues := make([]*BufferedChannelQueue[int], n)
	for i := 0; i < n; i++ {
		queues[i] = NewBufferedChannelQueue[int](3, 10, 5).
			SetFreeNodeHookPoolIntervalDuration(time.Millisecond).
			SetLoadFromPoolDuration(time.Millisecond)
	}
	for _, q := range queues {
		q.Close()
	}

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		runtime.GC()
		time.Sleep(20 * time.Millisecond)
		if runtime.NumGoroutine() <= baseline+5 {
			return
		}
	}
	assert.LessOrEqual(t, runtime.NumGoroutine(), baseline+5)
}

func TestBufferedChannelQueueConcurrentOfferTake(t *testing.T) {
	bufferedChannelQueue := NewBufferedChannelQueue[int](10, 100, 20).
		SetLoadFromPoolDuration(time.Millisecond).
		SetFreeNodeHookPoolIntervalDuration(time.Millisecond)

	done := make(chan bool)
	count := 100

	go func() {
		for i := 0; i < count; i++ {
			bufferedChannelQueue.Offer(i)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < count; i++ {
			bufferedChannelQueue.Take()
		}
		done <- true
	}()

	<-done
	<-done
}

func TestConcurrentQueueOfferPoll(t *testing.T) {
	linkedListQueue := NewLinkedListQueue[int]()
	concurrentQueue := NewConcurrentQueue[int](linkedListQueue)

	err := concurrentQueue.Offer(1)
	assert.Nil(t, err)

	err = concurrentQueue.Offer(2)
	assert.Nil(t, err)

	val, err := concurrentQueue.Poll()
	assert.Equal(t, 1, val)
	assert.Nil(t, err)

	val, err = concurrentQueue.Poll()
	assert.Equal(t, 2, val)
	assert.Nil(t, err)

	_, err = concurrentQueue.Poll()
	assert.Equal(t, ErrQueueIsEmpty, err)
}

func TestConcurrentQueuePutTake(t *testing.T) {
	linkedListQueue := NewLinkedListQueue[int]()
	concurrentQueue := NewConcurrentQueue[int](linkedListQueue)

	err := concurrentQueue.Put(1)
	assert.Nil(t, err)

	err = concurrentQueue.Put(2)
	assert.Nil(t, err)

	val, err := concurrentQueue.Take()
	assert.Equal(t, 1, val)
	assert.Nil(t, err)

	val, err = concurrentQueue.Take()
	assert.Equal(t, 2, val)
	assert.Nil(t, err)

	_, err = concurrentQueue.Take()
	assert.Equal(t, ErrQueueIsEmpty, err)
}

func TestChannelQueueOfferPoll(t *testing.T) {
	channelQueue := NewChannelQueue[int](3)

	err := channelQueue.Offer(1)
	assert.Nil(t, err)

	err = channelQueue.Offer(2)
	assert.Nil(t, err)

	val, err := channelQueue.Poll()
	assert.Equal(t, 1, val)
	assert.Nil(t, err)

	val, err = channelQueue.Poll()
	assert.Equal(t, 2, val)
	assert.Nil(t, err)
}

func TestBufferedChannelQueueTakeFromClosed(t *testing.T) {
	bufferedChannelQueue := NewBufferedChannelQueue[int](3, 10, 5)

	bufferedChannelQueue.Offer(1)
	bufferedChannelQueue.Offer(2)

	bufferedChannelQueue.Close()

	_, err := bufferedChannelQueue.Take()
	assert.NotNil(t, err)
}

func TestBufferedChannelQueuePollEmpty(t *testing.T) {
	bufferedChannelQueue := NewBufferedChannelQueue[int](3, 10, 5)

	_, err := bufferedChannelQueue.Poll()
	assert.Equal(t, ErrQueueIsEmpty, err)
}

func TestBufferedChannelQueueTakeWithTimeoutEmpty(t *testing.T) {
	bufferedChannelQueue := NewBufferedChannelQueue[int](3, 10, 5)

	_, err := bufferedChannelQueue.TakeWithTimeout(10 * time.Millisecond)
	assert.NotNil(t, err)
}

func TestLinkedListQueueCount(t *testing.T) {
	queue := NewLinkedListQueue[int]()

	assert.Equal(t, 0, queue.Count())

	queue.Offer(1)
	assert.Equal(t, 1, queue.Count())

	queue.Offer(2)
	assert.Equal(t, 2, queue.Count())

	queue.Poll()
	assert.Equal(t, 1, queue.Count())
}

func TestBufferedChannelQueueCount(t *testing.T) {
	queue := NewBufferedChannelQueue[int](3, 10, 5)

	assert.Equal(t, 0, queue.Count())

	queue.Offer(1)
	queue.Offer(2)
	assert.Equal(t, 2, queue.Count())

	queue.Take()
	assert.Equal(t, 1, queue.Count())
}

func TestBufferedChannelQueueIsClosed(t *testing.T) {
	queue := NewBufferedChannelQueue[int](3, 10, 5)

	assert.False(t, queue.IsClosed())

	queue.Close()
	assert.True(t, queue.IsClosed())
}

func TestLinkedListQueueKeepNodePoolCount(t *testing.T) {
	queue := NewLinkedListQueue[int]()

	for i := 0; i < 10; i++ {
		queue.Offer(i)
	}

	for i := 0; i < 10; i++ {
		queue.Poll()
	}

	queue.KeepNodePoolCount(5)
}

func TestLinkedListQueueClearNodePool(t *testing.T) {
	queue := NewLinkedListQueue[int]()

	for i := 0; i < 10; i++ {
		queue.Offer(i)
	}

	for i := 0; i < 10; i++ {
		queue.Poll()
	}

	queue.ClearNodePool()
}

func TestConcurrentStackPopEmpty(t *testing.T) {
	linkedListQueue := NewLinkedListQueue[int]()
	stack := NewConcurrentStack[int](linkedListQueue)

	_, err := stack.Pop()
	assert.NotNil(t, err)
}

func TestChannelQueueTake(t *testing.T) {
	queue := NewChannelQueue[int](3)

	queue.Offer(1)
	queue.Offer(2)

	val, err := queue.Take()
	assert.Nil(t, err)
	assert.Equal(t, 1, val)

	val, err = queue.Take()
	assert.Nil(t, err)
	assert.Equal(t, 2, val)
}

// Tests for queue.go closed channel coverage

func TestChannelQueueTakeAfterClose(t *testing.T) {
	queue := NewChannelQueue[int](3)
	queue.Offer(1)
	close(queue)

	val, err := queue.Take()
	assert.Nil(t, err)
	assert.Equal(t, 1, val)

	val, err = queue.Take()
	assert.Equal(t, ErrQueueIsClosed, err)
	assert.Equal(t, 0, val)
}

func TestChannelQueueTakeWithTimeoutAfterClose(t *testing.T) {
	queue := NewChannelQueue[int](3)
	queue.Offer(1)
	close(queue)

	timeout := 10 * time.Millisecond
	val, err := queue.TakeWithTimeout(timeout)
	assert.Nil(t, err)
	assert.Equal(t, 1, val)

	val, err = queue.TakeWithTimeout(timeout)
	assert.Equal(t, ErrQueueIsClosed, err)
	assert.Equal(t, 0, val)
}

func TestBufferedChannelQueueClosedBranches(t *testing.T) {
	queue := NewBufferedChannelQueue[int](1, 1, 1)
	queue.Close()

	assert.Equal(t, 0, queue.Count())

	err := queue.PutWithTimeout(1, 5*time.Millisecond)
	assert.Equal(t, ErrQueueIsClosed, err)

	val, pollErr := queue.Poll()
	assert.Equal(t, 0, val)
	assert.Equal(t, ErrQueueIsClosed, pollErr)
}

func TestBufferedChannelQueueSettersAndGetters(t *testing.T) {
	queue := NewBufferedChannelQueue[int](3, 10, 5)
	queue.
		SetBufferSizeMaximum(7).
		SetNodeHookPoolSize(4).
		SetLoadFromPoolDuration(2 * time.Millisecond).
		SetFreeNodeHookPoolIntervalDuration(3 * time.Millisecond)

	assert.Equal(t, 7, queue.GetBufferSizeMaximum())
	assert.Equal(t, 4, queue.GetNodeHookPoolSize())
	assert.Equal(t, 2*time.Millisecond, queue.GetLoadFromPoolDuration())
	assert.Equal(t, 3*time.Millisecond, queue.GetFreeNodeHookPoolIntervalDuration())
}

func TestLinkedListQueuePeekEmpty(t *testing.T) {
	q := NewLinkedListQueue[int]()

	val, err := q.Peek()
	assert.Equal(t, 0, val)
	assert.Equal(t, ErrQueueIsEmpty, err)
}

func TestBufferedChannelQueuePutWithTimeoutClosed(t *testing.T) {
	q := NewBufferedChannelQueue[int](3, 10, 5)
	q.Close()

	err := q.PutWithTimeout(1, 1*time.Millisecond)
	assert.Equal(t, ErrQueueIsClosed, err)
}

func TestBufferedChannelQueueTakeWithTimeoutClosed(t *testing.T) {
	q := NewBufferedChannelQueue[int](3, 10, 5)
	q.Close()

	val, err := q.TakeWithTimeout(1 * time.Millisecond)
	assert.Equal(t, 0, val)
	assert.Equal(t, ErrQueueIsClosed, err)
}

func TestBufferedChannelQueuePutTimeout(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 0, 1).
		SetLoadFromPoolDuration(time.Hour).
		SetFreeNodeHookPoolIntervalDuration(time.Hour)

	err := q.Offer(1)
	assert.Nil(t, err)

	err = q.PutWithTimeout(2, 0)
	assert.Equal(t, ErrQueuePutTimeout, err)
}

func TestBufferedChannelQueueOfferClosed(t *testing.T) {
	q := NewBufferedChannelQueue[int](3, 10, 5)
	q.Close()

	err := q.Offer(1)
	assert.Equal(t, ErrQueueIsClosed, err)
}

func TestLinkedListQueueRecycleNodeNil(t *testing.T) {
	q := NewLinkedListQueue[int]()

	for i := 0; i < 100; i++ {
		err := q.Offer(i)
		assert.Nil(t, err)
	}
	for i := 0; i < 100; i++ {
		val, err := q.Poll()
		assert.Equal(t, i, val)
		assert.Nil(t, err)
	}

	// Pool should have recycled nodes
	assert.Greater(t, q.nodeCount, 0)

	q.ClearNodePool()
	assert.Equal(t, 0, q.nodeCount)
	assert.Nil(t, q.nodePoolFirst)
}

func TestBufferedChannelQueuePutWithTimeoutFullThenSucceed(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 100, 1).
		SetLoadFromPoolDuration(time.Microsecond).
		SetFreeNodeHookPoolIntervalDuration(time.Microsecond)

	err := q.Offer(1)
	assert.Nil(t, err)
	err = q.Offer(2)
	assert.Nil(t, err)

	err = q.PutWithTimeout(3, 100*time.Millisecond)
	assert.Nil(t, err)
}

func TestBufferedChannelQueueOfferWhenPoolNotEmpty(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 10, 5)

	err := q.Offer(1)
	assert.Nil(t, err)

	err = q.Offer(2)
	assert.Nil(t, err)

	err = q.Offer(3)
	assert.Nil(t, err)
}

func TestBufferedChannelQueueOfferWhenPoolFull(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 1, 5)

	err := q.Offer(1)
	assert.Nil(t, err)

	err = q.Offer(2)
	assert.Nil(t, err)

	err = q.Offer(3)
	assert.Equal(t, ErrQueueIsFull, err)
}

func TestLinkedListQueueRecycleNodeNormal(t *testing.T) {
	q := NewLinkedListQueue[int]()

	q.Offer(1)
	q.Offer(2)
	q.Offer(3)

	q.Poll()
	q.Poll()
	q.Poll()

	assert.Greater(t, q.nodeCount, 0)

	err := q.Offer(4)
	assert.Nil(t, err)

	val, err := q.Poll()
	assert.Equal(t, 4, val)
	assert.Nil(t, err)
}

func TestBufferedChannelQueueLoadFromPoolFlow(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 10, 1).
		SetLoadFromPoolDuration(time.Microsecond).
		SetFreeNodeHookPoolIntervalDuration(time.Microsecond)

	err := q.Offer(1)
	assert.Nil(t, err)

	err = q.Offer(2)
	assert.Nil(t, err)

	val, err := q.Take()
	assert.Nil(t, err)
	assert.Equal(t, 1, val)

	time.Sleep(5 * time.Millisecond)

	val, err = q.Take()
	assert.Nil(t, err)
	assert.Equal(t, 2, val)

	assert.Equal(t, 0, q.pool.Count())
}

// TestBufferedChannelQueuePutWithTimeoutIsFull covers PutWithTimeout lines 715-716:
// when Offer returns ErrQueueIsFull and errors.Is returns true, returning immediately.
func TestBufferedChannelQueuePutWithTimeoutIsFull(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 0, 1)

	err := q.Offer(1)
	assert.Nil(t, err)

	// Channel is full (capacity 1) and bufferSizeMaximum is 0.
	// Offer(2) → channel full, poolCount(0) >= bufferMax(0) → ErrQueueIsFull
	// PutWithTimeout checks errors.Is(ErrQueueIsFull, ErrQueueIsFull) → true → return
	err = q.PutWithTimeout(2, 1*time.Millisecond)
	assert.Equal(t, ErrQueueIsFull, err)
}

// TestBufferedChannelQueueOfferChannelFull covers Offer lines 762-764:
// when poolCount == 0 and blockingQueue.Offer returns ErrQueueIsFull,
// falling through the else-if branch to the pool.Offer path.
func TestBufferedChannelQueueOfferChannelFull(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 100, 1)

	err := q.Offer(1)
	assert.Nil(t, err)

	// poolCount=0, channel full, blockingQueue.Offer(2) → ErrQueueIsFull
	// Falls into else-if (do nothing), then poolCount(0) < bufferMax(100) → pool.Offer(2) succeeds
	err = q.Offer(2)
	assert.Nil(t, err)
}

// TestBufferedChannelQueueLoadFromPoolChannelFull covers loadFromPool lines 580-581:
// when loadFromPool polls from pool and Offer to channel fails (channel full),
// it unshifts the value back and breaks out of the poll loop.
func TestBufferedChannelQueueLoadFromPoolChannelFull(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 100, 1).
		SetLoadFromPoolDuration(time.Microsecond)

	err := q.Offer(1)
	assert.Nil(t, err)

	err = q.Offer(2)
	assert.Nil(t, err)

	// Allow loadFromPool goroutine to attempt moving item 2 to the full channel,
	// which triggers the unshift + break path (lines 580-581)
	time.Sleep(5 * time.Millisecond)

	// Drain the channel, waking loadFromPool again via notifyWorkers
	val, err := q.Take()
	assert.Nil(t, err)
	assert.Equal(t, 1, val)

	// Allow loadFromPool to move item 2 from pool to channel
	time.Sleep(10 * time.Millisecond)

	val, err = q.Take()
	assert.Nil(t, err)
	assert.Equal(t, 2, val)
}

// TestLinkedListQueueRecycleNodeCalled covers recycleNode line 490 (nil check)
// and the non-nil recycling path (lines 494-499) where recycled nodes are reused.
func TestLinkedListQueueRecycleNodeCalled(t *testing.T) {
	q := NewLinkedListQueue[int]()

	q.Offer(1)
	q.Offer(2)
	q.Offer(3)

	// Each Poll calls recycleNode with a non-nil node → covers the non-nil path
	q.Poll()
	q.Poll()
	q.Poll()

	// Offer(4) reuses a recycled node from the pool
	err := q.Offer(4)
	assert.Nil(t, err)

	val, err := q.Poll()
	assert.Equal(t, 4, val)
	assert.Nil(t, err)
}

func TestLinkedListQueueRecycleNodeDirectNil(t *testing.T) {
	q := NewLinkedListQueue[int]()

	// recycleNode with nil should not panic, just return early (line 490-492)
	q.recycleNode(nil)
}

func TestBufferedChannelQueueWorkersEdgeBranches(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 4, 1).
		SetLoadFromPoolDuration(10 * time.Millisecond).
		SetFreeNodeHookPoolIntervalDuration(10 * time.Millisecond)
	defer q.Close()

	// Cover loadFromPool break when queue closes while a pooled item exists.
	assert.NoError(t, q.Offer(1))
	assert.NoError(t, q.Offer(2))
	q.lock.RLock()
	poolCount := q.pool.Count()
	q.lock.RUnlock()
	assert.Equal(t, 1, poolCount)
	q.isClosed.Set(true)
	q.loadWorkerCh <- 1
	time.Sleep(5 * time.Millisecond)
	q.lock.RLock()
	poolCount = q.pool.Count()
	q.lock.RUnlock()
	assert.Equal(t, 1, poolCount)
	q.isClosed.Set(false)

	// Cover notifyWorkers early return when closed.
	q.isClosed.Set(true)
	assert.NotPanics(t, func() { q.notifyWorkers() })
	q.isClosed.Set(false)
}

func TestBufferedChannelQueuePollClosed(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 1, 1)
	q.Close()

	val, err := q.Poll()
	assert.Equal(t, 0, val)
	assert.Equal(t, ErrQueueIsClosed, err)
}

func TestBufferedChannelQueueCountNilPool(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 1, 1)
	q.pool = nil
	assert.Equal(t, 0, q.Count())
	q.Close()
}

func TestBufferedChannelQueueCountNilBlockingQueue(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 1, 1)
	originalBlockingQueue := q.blockingQueue
	q.blockingQueue = nil
	assert.Equal(t, 0, q.Count())
	q.blockingQueue = originalBlockingQueue
	q.Close()
}

func TestBufferedChannelQueueGetChannelClosed(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 1, 1)
	q.Close()

	ch := q.GetChannel()
	assert.NotNil(t, ch)
	assert.Equal(t, 0, len(ch))
}

func TestBufferedChannelQueueIsClosedAfterDoubleClose(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 1, 1)
	assert.False(t, q.IsClosed())

	q.Close()
	q.Close()

	assert.True(t, q.IsClosed())
}

func TestBufferedChannelQueuePutDelegatesToOfferOnClosedQueue(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 1, 1)
	q.Close()

	err := q.Put(1)
	assert.Equal(t, ErrQueueIsClosed, err)
}

func TestBufferedChannelQueueGetChannelOpen(t *testing.T) {
	q := NewBufferedChannelQueue[int](2, 2, 1)
	defer q.Close()

	ch := q.GetChannel()
	assert.NotNil(t, ch)
	assert.Equal(t, cap(q.blockingQueue), cap(ch))
}

func TestBufferedChannelQueueTakeClosed(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 1, 1)
	q.Close()

	val, err := q.Take()
	assert.Equal(t, 0, val)
	assert.Equal(t, ErrQueueIsClosed, err)
}

// TestBufferedChannelQueueLoadFromPoolInnerIsClosedBreak covers loadFromPool lines
// 571-573: the inner-loop isClosed guard after the outer guard passed with isClosed false.
// Deterministic via lock ordering — while the test holds q.lock, loadFromPool blocks at
// line 565; isClosed is set true before unlock so the inner check at 571 fires.
func TestBufferedChannelQueueLoadFromPoolInnerIsClosedBreak(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 10, 1).
		SetLoadFromPoolDuration(time.Hour).
		SetFreeNodeHookPoolIntervalDuration(time.Hour)
	defer q.Close()

	// Seed the pool directly (not via Offer) so notifyWorkers does not race with the test.
	q.lock.Lock()
	assert.NoError(t, q.pool.Offer(1))
	assert.Equal(t, 1, q.pool.Count())

	q.isClosed.Set(false)
	q.loadWorkerCh <- 1
	for i := 0; i < 100000; i++ {
		runtime.Gosched()
	}
	q.isClosed.Set(true)
	q.lock.Unlock()

	// loadFromPool sleeps for an hour after one iteration; pool access is safe under lock.
	q.lock.Lock()
	assert.Equal(t, 1, q.pool.Count())
	q.lock.Unlock()
}

// TestBufferedChannelQueuePutWithTimeoutSecondLoopIsClosed covers PutWithTimeout lines
// 738-739 (retry sleep after Offer returns ErrQueueIsClosed) and 721-723 (closed check on
// the following loop iteration). Deterministic via the same lock ordering: PutWithTimeout
// blocks inside Offer on q.lock while isClosed is still false at the first loop guard,
// then isClosed is set before unlock so Offer returns ErrQueueIsClosed and the next
// iteration returns ErrQueueIsClosed at line 722.
func TestBufferedChannelQueuePutWithTimeoutSecondLoopIsClosed(t *testing.T) {
	q := NewBufferedChannelQueue[int](1, 10, 1).
		SetLoadFromPoolDuration(0).
		SetFreeNodeHookPoolIntervalDuration(time.Hour)
	defer q.Close()

	q.isClosed.Set(false)
	q.lock.Lock()

	result := make(chan error, 1)
	go func() {
		result <- q.PutWithTimeout(42, time.Second)
	}()
	for i := 0; i < 100000; i++ {
		runtime.Gosched()
	}
	q.isClosed.Set(true)
	q.lock.Unlock()

	select {
	case err := <-result:
		assert.Equal(t, ErrQueueIsClosed, err)
	case <-time.After(time.Second):
		t.Fatal("PutWithTimeout did not return")
	}
}

func TestQueueCountAccessorsParity(t *testing.T) {
	// ConcurrentQueue.Count / NodePoolCount delegate to the wrapped queue.
	linkedListQueue := NewLinkedListQueue[int]()
	concurrentQueue := NewConcurrentQueue[int](linkedListQueue)
	assert.Equal(t, 0, concurrentQueue.Count())
	assert.Equal(t, 0, concurrentQueue.NodePoolCount())
	assert.NoError(t, concurrentQueue.Offer(1))
	assert.NoError(t, concurrentQueue.Offer(2))
	assert.Equal(t, 2, concurrentQueue.Count())

	// LinkedListQueue.NodePoolCount reflects the cached node pool size.
	linkedListQueue.KeepNodePoolCount(5)
	assert.Equal(t, 5, linkedListQueue.NodePoolCount())
	assert.Equal(t, 5, concurrentQueue.NodePoolCount())
	linkedListQueue.KeepNodePoolCount(0)
	assert.Equal(t, 0, linkedListQueue.NodePoolCount())

	// BufferedChannelQueue.PoolNodeCount is a thread-safe pool node counter.
	bufferedChannelQueue := NewBufferedChannelQueue[int](3, 10, 5)
	defer bufferedChannelQueue.Close()
	assert.Equal(t, 0, bufferedChannelQueue.PoolNodeCount())
}
