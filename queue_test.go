package fpgo

import (
	"runtime"
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
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
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
		defer wg.Done()
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

	wg.Wait()
}

func TestLinkedListQueue(t *testing.T) {
	var queue Queue
	var stack Stack
	var err error
	var result interface{}

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
	assert.Equal(t, nil, r)
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
			var result interface{}
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
		for expected := 1; expected <= 6; expected++ {
			// Retry on timeout: items arrive in FIFO order; a transient
			// TakeWithTimeout miss under load must not fail the assertion.
			var r interface{}
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
			var result interface{}
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
	// synchronizing with Close(). If Close() closed those channels between the
	// isClosed check and the Offer, the Offer would panic with send-on-closed-channel.
	// Fix: notifyWorkers() acquires a dedicated notifyLock sync.Mutex before sending.
	// Close() also holds notifyLock while closing loadWorkerCh, freeNodeWorkerCh, and
	// blockingQueue, so notifyWorkers cannot send on a channel being closed.
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

func TestBufferedChannelQueueDoubleClose(t *testing.T) {
	q := NewBufferedChannelQueue(3, 10, 5)
	q.Close()
	q.Close()
	assert.True(t, q.IsClosed())
}

func TestBufferedChannelQueueCloseStopsFreeNodeGoroutine(t *testing.T) {
	runtime.GC()
	time.Sleep(10 * time.Millisecond)
	baseline := runtime.NumGoroutine()

	const n = 50
	queues := make([]*BufferedChannelQueue, n)
	for i := 0; i < n; i++ {
		queues[i] = NewBufferedChannelQueue(3, 10, 5).
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

func TestBufferedChannelQueuePutWithTimeoutClosedMidWait(t *testing.T) {
	q := NewBufferedChannelQueue(1, 0, 1).
		SetLoadFromPoolDuration(time.Millisecond)
	assert.NoError(t, q.Offer(1))

	hold := make(chan struct{})
	go func() {
		q.lock.Lock()
		close(hold)
		time.Sleep(30 * time.Millisecond)
		q.lock.Unlock()
		q.Close()
	}()
	<-hold

	errCh := make(chan error, 1)
	go func() {
		errCh <- q.PutWithTimeout(2, time.Second)
	}()
	for i := 0; i < 100000; i++ {
		runtime.Gosched()
	}

	select {
	case err := <-errCh:
		assert.Equal(t, ErrQueueIsClosed, err)
	case <-time.After(2 * time.Second):
		t.Fatal("PutWithTimeout did not return")
	}
}

func ifacePtr(v interface{}) *interface{} {
	p := new(interface{})
	*p = v
	return p
}

func TestChannelQueuePutTake(t *testing.T) {
	q := NewChannelQueue(5)
	assert.NoError(t, q.Put(1))
	assert.NoError(t, q.Put(2))
	assert.NoError(t, q.Put(3))
	val, err := q.Take()
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
	val, err = q.Take()
	assert.NoError(t, err)
	assert.Equal(t, 2, val)
	val, err = q.Take()
	assert.NoError(t, err)
	assert.Equal(t, 3, val)
	close(q)
	_, err = q.Take()
	assert.Equal(t, ErrQueueIsClosed, err)
}

func TestConcurrentQueuePoll(t *testing.T) {
	q := NewConcurrentQueue(NewLinkedListQueue())
	assert.NoError(t, q.Offer(1))
	assert.NoError(t, q.Offer(2))
	val, err := q.Poll()
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
	val, err = q.Poll()
	assert.NoError(t, err)
	assert.Equal(t, 2, val)
	_, err = q.Poll()
	assert.Equal(t, ErrQueueIsEmpty, err)
}

func TestConcurrentQueueCount(t *testing.T) {
	q := NewConcurrentQueue(NewLinkedListQueue())
	assert.Equal(t, 0, q.Count())
	assert.NoError(t, q.Offer(1))
	assert.NoError(t, q.Offer(2))
	assert.Equal(t, 2, q.Count())
}

func TestConcurrentStack(t *testing.T) {
	stack := NewConcurrentStack(NewLinkedListQueue())
	assert.NoError(t, stack.Push(1))
	assert.NoError(t, stack.Push(2))
	assert.NoError(t, stack.Push(3))
	val, err := stack.Pop()
	assert.NoError(t, err)
	assert.Equal(t, 3, val)
	val, err = stack.Pop()
	assert.NoError(t, err)
	assert.Equal(t, 2, val)
	val, err = stack.Pop()
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
	_, err = stack.Pop()
	assert.Equal(t, ErrStackIsEmpty, err)
}

func TestLinkedListItem(t *testing.T) {
	item1 := &LinkedListItem{Val: ifacePtr(1)}
	item2 := &LinkedListItem{Val: ifacePtr(2)}
	assert.Equal(t, 1, item1.Count())
	item1.Next = item2
	assert.Equal(t, 2, item1.Count())
	item1.AddLast(&LinkedListItem{Val: ifacePtr(3)})
	assert.Equal(t, 3, item1.Count())
	assert.Equal(t, 3, *item1.Last().Val)
}

func TestDoublyListItem(t *testing.T) {
	item1 := &DoublyListItem{Val: ifacePtr(1)}
	item2 := &DoublyListItem{Val: ifacePtr(2)}
	item1.Next = item2
	item2.Prev = item1
	assert.Equal(t, 2, item1.Count())
	item1.AddLast(&DoublyListItem{Val: ifacePtr(3)})
	assert.Equal(t, 3, item1.Count())
	item1.AddFirst(&DoublyListItem{Val: ifacePtr(0)})
	assert.Equal(t, 4, item1.Count())
}

func TestLinkedListQueuePeek(t *testing.T) {
	q := NewLinkedListQueue()
	_, err := q.Peek()
	assert.Equal(t, ErrQueueIsEmpty, err)
	assert.NoError(t, q.Offer(1))
	val, err := q.Peek()
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
}

func TestBufferedChannelQueueGetters(t *testing.T) {
	q := NewBufferedChannelQueue(3, 100, 10).
		SetLoadFromPoolDuration(5 * time.Millisecond).
		SetFreeNodeHookPoolIntervalDuration(10 * time.Millisecond).
		SetNodeHookPoolSize(4)
	defer q.Close()
	assert.Equal(t, 100, q.GetBufferSizeMaximum())
	assert.Equal(t, 4, q.GetNodeHookPoolSize())
	assert.Equal(t, 5*time.Millisecond, q.GetLoadFromPoolDuration())
	assert.Equal(t, 10*time.Millisecond, q.GetFreeNodeHookPoolIntervalDuration())
}

func TestBufferedChannelQueueCountAndIsClosed(t *testing.T) {
	q := NewBufferedChannelQueue(2, 10, 5)
	assert.False(t, q.IsClosed())
	assert.NoError(t, q.Offer(1))
	assert.Equal(t, 1, q.Count())
	q.Close()
	assert.True(t, q.IsClosed())
	assert.Equal(t, 0, q.Count())
}

func TestBufferedChannelQueueTake(t *testing.T) {
	q := NewBufferedChannelQueue(2, 10, 5).
		SetLoadFromPoolDuration(time.Microsecond).
		SetFreeNodeHookPoolIntervalDuration(time.Microsecond)
	defer q.Close()
	assert.NoError(t, q.Offer(1))
	val, err := q.Take()
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
}

func TestBufferedChannelQueueTakeWithTimeout(t *testing.T) {
	q := NewBufferedChannelQueue(2, 10, 5).
		SetLoadFromPoolDuration(time.Microsecond).
		SetFreeNodeHookPoolIntervalDuration(time.Microsecond)
	defer q.Close()
	assert.NoError(t, q.Offer(42))
	val, err := q.TakeWithTimeout(50 * time.Millisecond)
	assert.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestBufferedChannelQueuePoll(t *testing.T) {
	q := NewBufferedChannelQueue(2, 10, 5).
		SetLoadFromPoolDuration(time.Microsecond).
		SetFreeNodeHookPoolIntervalDuration(time.Microsecond)
	defer q.Close()
	assert.NoError(t, q.Offer(1))
	val, err := q.Poll()
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
	_, err = q.Poll()
	assert.Equal(t, ErrQueueIsEmpty, err)
}

func TestBufferedChannelQueueTakeClosed(t *testing.T) {
	q := NewBufferedChannelQueue(2, 10, 5)
	q.Close()
	_, err := q.Take()
	assert.Equal(t, ErrQueueIsClosed, err)
}

func TestBufferedChannelQueuePutWithTimeoutSuccess(t *testing.T) {
	q := NewBufferedChannelQueue(2, 10, 5).
		SetLoadFromPoolDuration(time.Microsecond).
		SetFreeNodeHookPoolIntervalDuration(time.Microsecond)
	defer q.Close()
	assert.NoError(t, q.PutWithTimeout(1, 50*time.Millisecond))
}

func TestBufferedChannelQueuePutWithTimeoutClosed(t *testing.T) {
	q := NewBufferedChannelQueue(2, 10, 5)
	q.Close()
	assert.Equal(t, ErrQueueIsClosed, q.PutWithTimeout(1, 10*time.Millisecond))
}

func TestBufferedChannelQueuePutWithTimeoutIsFull(t *testing.T) {
	q := NewBufferedChannelQueue(1, 0, 1)
	assert.NoError(t, q.Offer(1))
	assert.Equal(t, ErrQueueIsFull, q.PutWithTimeout(2, 5*time.Millisecond))
	q.Close()
}

func TestBufferedChannelQueuePutWithTimeoutTimeout(t *testing.T) {
	q := NewBufferedChannelQueue(1, 0, 1).
		SetLoadFromPoolDuration(time.Hour).
		SetFreeNodeHookPoolIntervalDuration(time.Hour)
	defer q.Close()
	assert.NoError(t, q.Offer(1))
	assert.Equal(t, ErrQueuePutTimeout, q.PutWithTimeout(2, 0))
}

func TestBufferedChannelQueuePutWithTimeoutFullThenSucceed(t *testing.T) {
	q := NewBufferedChannelQueue(1, 100, 1).
		SetLoadFromPoolDuration(time.Microsecond).
		SetFreeNodeHookPoolIntervalDuration(time.Microsecond)
	defer q.Close()
	assert.NoError(t, q.Offer(1))
	assert.NoError(t, q.Offer(2))
	assert.NoError(t, q.PutWithTimeout(3, 100*time.Millisecond))
}

func TestBufferedChannelQueueConcurrentTakePutWithTimeout(t *testing.T) {
	q := NewBufferedChannelQueue(4, 50, 5).
		SetLoadFromPoolDuration(time.Microsecond).
		SetFreeNodeHookPoolIntervalDuration(time.Microsecond)
	defer q.Close()
	const n = 50
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < n; i++ {
			assert.NoError(t, q.PutWithTimeout(i, 200*time.Millisecond))
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < n; i++ {
			_, err := q.Take()
			assert.NoError(t, err)
		}
	}()
	wg.Wait()
}

func TestBufferedChannelQueueConcurrentCount(t *testing.T) {
	q := NewBufferedChannelQueue(8, 100, 10).
		SetLoadFromPoolDuration(time.Microsecond).
		SetFreeNodeHookPoolIntervalDuration(time.Microsecond)
	defer q.Close()
	const n = 100
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < n; i++ {
			assert.NoError(t, q.Put(i))
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < n; i++ {
			_, err := q.Take()
			assert.NoError(t, err)
		}
	}()
	wg.Wait()
	assert.Equal(t, 0, q.Count())
}
