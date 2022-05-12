package fpgo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChannelQueue(t *testing.T) {
	var queue Queue
	var err error
	var result interface{}
	var timeout time.Duration

	channelQueue := NewChannelQueue(3)
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
	assert.Equal(t, nil, result)
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
		assert.Equal(t, nil, result)
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
	var queue Queue
	var stack Stack
	var err error
	var result interface{}
	var timeout time.Duration

	linkedListQueue := NewLinkedListQueue()
	queue = linkedListQueue
	stack = linkedListQueue
	concurrentQueue := NewConcurrentQueue(queue)

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
	assert.Equal(t, nil, result)
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
	assert.Equal(t, nil, result)
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
		assert.Equal(t, nil, result)
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

func TestNewBufferedChannelQueue(t *testing.T) {
	var queue Queue
	var err error
	var result interface{}
	var timeout time.Duration

	bufferedChannelQueue := NewBufferedChannelQueue(3, 10000, 100)
	bufferedChannelQueue.SetLoadFromPoolDuration(time.Millisecond / 10)
	bufferedChannelQueue.SetFreeNodeHookPoolIntervalDuration(1 * time.Millisecond)
	queue = bufferedChannelQueue

	// Sync

	timeout = 1 * time.Millisecond
	bufferedChannelQueue.SetBufferSizeMaximum(1)

	err = queue.Offer(1)
	assert.Equal(t, nil, err)
	time.Sleep(1 * timeout)
	err = queue.Offer(2)
	assert.Equal(t, nil, err)
	time.Sleep(1 * timeout)
	err = queue.Offer(3)
	assert.Equal(t, nil, err)
	time.Sleep(1 * timeout)
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

	bufferedChannelQueue.SetBufferSizeMaximum(10000)
	timeout = 1 * time.Millisecond
	asyncTaskDone := make(chan bool)
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
			err := bufferedChannelQueue.Offer(i)
			assert.Equal(t, nil, err)
		}
		assert.Equal(t, 0, bufferedChannelQueue.pool.nodeCount)
	}()

	<-asyncTaskDone

	result, err = bufferedChannelQueue.Poll()
	assert.Equal(t, ErrQueueIsEmpty, err)
	assert.Equal(t, nil, result)

	time.Sleep(1 * timeout)

	assert.GreaterOrEqual(t, 100, bufferedChannelQueue.pool.nodeCount)
	close(asyncTaskDone)
}
