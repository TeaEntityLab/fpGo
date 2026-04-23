package fpgo

import (
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

	result = 0
	timeout = 1 * time.Millisecond
	go func() {
		result, err = channelQueue.TakeWithTimeout(timeout)
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
		err = channelQueue.PutWithTimeout(1, timeout)
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

	time.Sleep(2 * timeout)
}

func TestLinkedListQueue(t *testing.T) {
	var queue Queue[int]
	var stack Stack[int]
	var err error
	var result int
	var timeout time.Duration

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

	result = 0
	timeout = 1 * time.Millisecond
	go func() {
		time.Sleep(timeout)
		assert.Equal(t, 3, linkedListQueue.Count())
		assert.Equal(t, 0, linkedListQueue.nodeCount)
		result, err = concurrentQueue.Take()
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, result)
		assert.Equal(t, 2, linkedListQueue.Count())
		assert.Equal(t, 1, linkedListQueue.nodeCount)
		result, err = concurrentQueue.Take()
		assert.Equal(t, nil, err)
		assert.Equal(t, 2, result)
		assert.Equal(t, 1, linkedListQueue.Count())
		assert.Equal(t, 2, linkedListQueue.nodeCount)
		result, err = concurrentQueue.Take()
		assert.Equal(t, nil, err)
		assert.Equal(t, 3, result)
		assert.Equal(t, 0, linkedListQueue.Count())
		assert.Equal(t, 3, linkedListQueue.nodeCount)
		result, err = concurrentQueue.Take()
		assert.NotEqual(t, 4, result)
		assert.NotEqual(t, nil, err)
		assert.Equal(t, 0, result)
		assert.Equal(t, ErrQueueIsEmpty, err)
		assert.Equal(t, 0, linkedListQueue.Count())
		assert.Equal(t, 3, linkedListQueue.nodeCount)
	}()
	go func() {
		assert.Equal(t, 0, linkedListQueue.Count())
		assert.Equal(t, 3, linkedListQueue.nodeCount)
		err = concurrentQueue.Put(1)
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, linkedListQueue.Count())
		assert.Equal(t, 2, linkedListQueue.nodeCount)
		err = concurrentQueue.Put(2)
		assert.Equal(t, nil, err)
		assert.Equal(t, 2, linkedListQueue.Count())
		assert.Equal(t, 1, linkedListQueue.nodeCount)
		err = concurrentQueue.Put(3)
		assert.Equal(t, nil, err)
		assert.Equal(t, 3, linkedListQueue.Count())
		assert.Equal(t, 0, linkedListQueue.nodeCount)

		time.Sleep(3 * timeout / 2)

		assert.Equal(t, 0, linkedListQueue.Count())
		assert.Equal(t, 3, linkedListQueue.nodeCount)
		linkedListQueue.KeepNodePoolCount(2)
		assert.Equal(t, 2, linkedListQueue.nodeCount)
		assert.Equal(t, 2, linkedListQueue.nodePoolFirst.Count())
		err = concurrentQueue.Put(4)
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, linkedListQueue.Count())
		assert.Equal(t, 1, linkedListQueue.nodeCount)
		err = concurrentQueue.Put(5)
		assert.Equal(t, nil, err)
		assert.Equal(t, 2, linkedListQueue.Count())
		assert.Equal(t, 0, linkedListQueue.nodeCount)
		err = concurrentQueue.Put(6)
		assert.Equal(t, nil, err)
		assert.Equal(t, 3, linkedListQueue.Count())
		assert.Equal(t, 0, linkedListQueue.nodeCount)
	}()

	time.Sleep(2 * timeout)

	assert.Equal(t, 3, linkedListQueue.Count())
	assert.Equal(t, 0, linkedListQueue.nodeCount)
	linkedListQueue.Clear()
	assert.Equal(t, 0, linkedListQueue.Count())
	assert.Equal(t, 3, linkedListQueue.nodeCount)
	node := linkedListQueue.nodePoolFirst
	for node != nil {
		assert.Nil(t, node.Val)
		node = node.Next
	}
	linkedListQueue.ClearNodePool()
	assert.Equal(t, 0, linkedListQueue.Count())
	assert.Equal(t, 0, linkedListQueue.nodeCount)
	assert.Nil(t, linkedListQueue.nodePoolFirst)

	go func() {
		time.Sleep(1 * time.Millisecond)

		for i := 1; i <= 10000; i++ {
			result, err := concurrentQueue.Take()
			assert.Equal(t, nil, err)
			assert.Equal(t, i, result)
		}
	}()
	go func() {
		for i := 1; i <= 10000; i++ {
			err := concurrentQueue.Offer(i)
			assert.Equal(t, nil, err)
		}
	}()

	time.Sleep(2 * timeout)
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
		result, err = bufferedChannelQueue.TakeWithTimeout(timeout)
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, result)
		result, err = bufferedChannelQueue.TakeWithTimeout(timeout)
		assert.Equal(t, nil, err)
		assert.Equal(t, 2, result)
		result, err = bufferedChannelQueue.TakeWithTimeout(timeout)
		assert.Equal(t, nil, err)
		assert.Equal(t, 3, result)
		result, err = bufferedChannelQueue.TakeWithTimeout(timeout)
		assert.Equal(t, nil, err)
		assert.Equal(t, 4, result)
		result, err = bufferedChannelQueue.TakeWithTimeout(timeout)
		assert.Equal(t, nil, err)
		assert.Equal(t, 5, result)
		result, err = bufferedChannelQueue.TakeWithTimeout(timeout)
		assert.Equal(t, nil, err)
		assert.Equal(t, 6, result)
		asyncTaskDone <- true
	}()
	go func() {
		err = bufferedChannelQueue.Put(1)
		assert.Equal(t, nil, err)
		err = bufferedChannelQueue.Put(2)
		assert.Equal(t, nil, err)
		err = bufferedChannelQueue.Put(3)
		assert.Equal(t, nil, err)
		err = bufferedChannelQueue.Put(4)
		assert.Equal(t, nil, err)
		err = bufferedChannelQueue.Put(5)
		assert.Equal(t, nil, err)
		err = bufferedChannelQueue.Put(6)
		assert.Equal(t, nil, err)
	}()

	<-asyncTaskDone

	bufferedChannelQueue.SetBufferSizeMaximum(10000)
	timeout = 10 * time.Millisecond
	go func() {
		for i := 1; i <= 10000; i++ {
			result, err := bufferedChannelQueue.TakeWithTimeout(timeout)
			assert.Equal(t, nil, err)
			assert.Equal(t, i, result)
		}
		asyncTaskDone <- true
	}()
	go func() {
		for i := 1; i <= 10000; i++ {
			// err := bufferedChannelQueue.PutWithTimeout(i, timeout)
			// err := bufferedChannelQueue.Put(i)
			err := bufferedChannelQueue.Offer(i)
			assert.Equal(t, nil, err)
		}
		assert.Equal(t, 0, bufferedChannelQueue.pool.nodeCount)
	}()

	<-asyncTaskDone

	result, err = bufferedChannelQueue.Poll()
	assert.Equal(t, ErrQueueIsEmpty, err)
	assert.Equal(t, 0, result)

	time.Sleep(1 * timeout)

	assert.GreaterOrEqual(t, bufferedChannelQueue.pool.nodeCount, 100)
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

	assert.True(t, bufferedChannelQueue.IsClosed())

	err := bufferedChannelQueue.Offer(1)
	assert.Equal(t, ErrQueueIsClosed, err)
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
