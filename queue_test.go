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
		result, err = channelQueue.TakeWithTimeout(&timeout)
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, result)
		result, err = channelQueue.TakeWithTimeout(&timeout)
		assert.Equal(t, nil, err)
		assert.Equal(t, 2, result)
		result, err = channelQueue.TakeWithTimeout(&timeout)
		assert.Equal(t, nil, err)
		assert.Equal(t, 3, result)
		result, err = channelQueue.TakeWithTimeout(&timeout)
		assert.NotEqual(t, 4, result)
		assert.NotEqual(t, nil, err)
		assert.Equal(t, 0, result)
		assert.Equal(t, ErrQueueTakeTimeout, err)
	}()
	go func() {
		err = channelQueue.PutWithTimeout(1, &timeout)
		assert.Equal(t, nil, err)
		err = channelQueue.PutWithTimeout(2, &timeout)
		assert.Equal(t, nil, err)
		err = channelQueue.PutWithTimeout(3, &timeout)
		assert.Equal(t, nil, err)

		time.Sleep(3 * timeout / 2)

		err = channelQueue.PutWithTimeout(4, &timeout)
		assert.Equal(t, nil, err)
		err = channelQueue.PutWithTimeout(5, &timeout)
		assert.Equal(t, nil, err)
		err = channelQueue.PutWithTimeout(6, &timeout)
		assert.Equal(t, nil, err)
		err = channelQueue.PutWithTimeout(7, &timeout)
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
}
