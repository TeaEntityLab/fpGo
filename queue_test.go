package fpgo

import (
	"sync"
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
		r, e := channelQueue.TakeWithTimeout(timeout)
		assert.Equal(t, nil, e)
		assert.Equal(t, 1, r)
		r, e = channelQueue.TakeWithTimeout(timeout)
		assert.Equal(t, nil, e)
		assert.Equal(t, 2, r)
		r, e = channelQueue.TakeWithTimeout(timeout)
		assert.Equal(t, nil, e)
		assert.Equal(t, 3, r)
		r, e = channelQueue.TakeWithTimeout(timeout)
		assert.NotEqual(t, 4, r)
		assert.NotEqual(t, nil, e)
		assert.Equal(t, nil, r)
		assert.Equal(t, ErrQueueTakeTimeout, e)
	}()
	go func() {
		e := channelQueue.PutWithTimeout(1, timeout)
		assert.Equal(t, nil, e)
		e = channelQueue.PutWithTimeout(2, timeout)
		assert.Equal(t, nil, e)
		e = channelQueue.PutWithTimeout(3, timeout)
		assert.Equal(t, nil, e)

		time.Sleep(3 * timeout / 2)

		e = channelQueue.PutWithTimeout(4, timeout)
		assert.Equal(t, nil, e)
		e = channelQueue.PutWithTimeout(5, timeout)
		assert.Equal(t, nil, e)
		e = channelQueue.PutWithTimeout(6, timeout)
		assert.Equal(t, nil, e)
		e = channelQueue.PutWithTimeout(7, timeout)
		assert.NotEqual(t, nil, e)
		assert.Equal(t, ErrQueuePutTimeout, e)
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
	assert.Equal(t, 10, concurrentQueue.NodePoolCount())
	assert.Equal(t, 10, linkedListQueue.nodePoolFirst.Count())
	linkedListQueue.KeepNodePoolCount(2)
	assert.Equal(t, 2, concurrentQueue.NodePoolCount())
	assert.Equal(t, 2, linkedListQueue.nodePoolFirst.Count())
	linkedListQueue.KeepNodePoolCount(0)
	assert.Equal(t, 0, concurrentQueue.NodePoolCount())
	assert.Nil(t, linkedListQueue.nodePoolFirst)
	linkedListQueue.KeepNodePoolCount(5)
	assert.Equal(t, 5, concurrentQueue.NodePoolCount())
	assert.Equal(t, 5, linkedListQueue.nodePoolFirst.Count())
	linkedListQueue.KeepNodePoolCount(3)
	assert.Equal(t, 3, concurrentQueue.NodePoolCount())
	assert.Equal(t, 3, linkedListQueue.nodePoolFirst.Count())

	result = 0
	timeout = 1 * time.Millisecond
	go func() {
		time.Sleep(timeout)
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
		assert.Equal(t, nil, r)
		assert.Equal(t, ErrQueueIsEmpty, e)
		assert.Equal(t, 0, concurrentQueue.Count())
		assert.Equal(t, 3, concurrentQueue.NodePoolCount())
	}()
	go func() {
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

		time.Sleep(3 * timeout / 2)

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
	}()

	time.Sleep(2 * timeout)

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

	bufferedChannelQueue := NewBufferedChannelQueue(3, 10000, 100).
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
		r, e := bufferedChannelQueue.TakeWithTimeout(timeout)
		assert.Equal(t, nil, e)
		assert.Equal(t, 1, r)
		r, e = bufferedChannelQueue.TakeWithTimeout(timeout)
		assert.Equal(t, nil, e)
		assert.Equal(t, 2, r)
		r, e = bufferedChannelQueue.TakeWithTimeout(timeout)
		assert.Equal(t, nil, e)
		assert.Equal(t, 3, r)
		r, e = bufferedChannelQueue.TakeWithTimeout(timeout)
		assert.Equal(t, nil, e)
		assert.Equal(t, 4, r)
		r, e = bufferedChannelQueue.TakeWithTimeout(timeout)
		assert.Equal(t, nil, e)
		assert.Equal(t, 5, r)
		r, e = bufferedChannelQueue.TakeWithTimeout(timeout)
		assert.Equal(t, nil, e)
		assert.Equal(t, 6, r)
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
	}()

	<-asyncTaskDone

	result, err = bufferedChannelQueue.Poll()
	assert.Equal(t, ErrQueueIsEmpty, err)
	assert.Equal(t, nil, result)

	time.Sleep(1 * timeout)

	assert.GreaterOrEqual(t, bufferedChannelQueue.PoolNodeCount(), 100)
	close(asyncTaskDone)
}

func TestBufferedChannelQueueConcurrentCloseNotifyNoPanic(t *testing.T) {
	// Regression: notifyWorkers() sent on loadWorkerCh/freeNodeWorkerCh without
	// checking isClosed under lock. If Close() closed those channels between the
	// isClosed check and the Offer, the Offer would panic with send-on-closed-channel.
	// Fix: notifyWorkers now checks isClosed under RLock.
	for i := 0; i < 100; i++ {
		q := NewBufferedChannelQueue(3, 10000, 100)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			q.Close()
		}()
		go func() {
			defer wg.Done()
			// GetChannel calls notifyWorkers internally
			_ = q.GetChannel()
		}()
		wg.Wait()
	}
}
