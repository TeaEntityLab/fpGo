package fpgo

import (
	"errors"
	"sync"
	"time"
)

// Queue Queue inspired by Collection utils
type Queue interface {
	Put(val interface{}) error
	Take() (interface{}, error)
	Offer(val interface{}) error
	Poll() (interface{}, error)
}

// Stack Stack inspired by Collection utils
type Stack interface {
	Push(val interface{}) error
	Pop() (interface{}, error)
}

var (
	// ErrQueueIsEmpty Queue Is Empty
	ErrQueueIsEmpty = errors.New("queue is empty")
	// ErrQueueIsFull Queue Is Full
	ErrQueueIsFull = errors.New("queue is full")
	// ErrQueueIsClosed Queue Is Closed
	ErrQueueIsClosed = errors.New("queue is closed")
	// ErrQueueTakeTimeout Queue Take Timeout
	ErrQueueTakeTimeout = errors.New("queue take timeout")
	// ErrQueuePutTimeout Queue Put Timeout
	ErrQueuePutTimeout = errors.New("queue put timeout")

	// ErrStackIsEmpty Stack Is Empty
	ErrStackIsEmpty = errors.New("stack is empty")
	// ErrStackIsFull Stack Is Full
	ErrStackIsFull = errors.New("stack is full")
)

// ConcurrentQueue & ConcurrentStack

// ConcurrentQueue ConcurrentQueue inspired by Collection utils
type ConcurrentQueue struct {
	lock  sync.RWMutex
	queue Queue
}

// NewConcurrentQueue New ConcurrentQueue instance from a Queue
func NewConcurrentQueue(queue Queue) *ConcurrentQueue {
	return &ConcurrentQueue{
		queue: queue,
	}
}

// Put Put the val(probably blocking)
func (q *ConcurrentQueue) Put(val interface{}) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.queue.Put(val)
}

// Take Take the val(probably blocking)
func (q *ConcurrentQueue) Take() (interface{}, error) {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return q.queue.Take()
}

// Offer Offer the val(non-blocking)
func (q *ConcurrentQueue) Offer(val interface{}) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.queue.Offer(val)
}

// Poll Poll the val(non-blocking)
func (q *ConcurrentQueue) Poll() (interface{}, error) {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return q.queue.Poll()
}

// ConcurrentStack

// ConcurrentStack ConcurrentStack inspired by Collection utils
type ConcurrentStack struct {
	lock  sync.RWMutex
	stack Stack
}

// NewConcurrentStack New ConcurrentStack instance from a Stack
func NewConcurrentStack(stack Stack) *ConcurrentStack {
	return &ConcurrentStack{
		stack: stack,
	}
}

// Put Put the val(probably blocking)
func (q *ConcurrentStack) Push(val interface{}) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.stack.Push(val)
}

// Take Take the val(probably blocking)
func (q *ConcurrentStack) Pop() (interface{}, error) {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return q.stack.Pop()
}

// ChannelQueue

// ChannelQueue ChannelQueue inspired by Collection utils
type ChannelQueue chan interface{}

// NewChannelQueue New ChannelQueue instance with capacity
func NewChannelQueue(capacity int) ChannelQueue {
	return make(ChannelQueue, capacity)
}

// Put Put the val(blocking)
func (q ChannelQueue) Put(val interface{}) error {
	q <- val
	return nil
}

// PutWithTimeout Put the val(blocking), with timeout
func (q ChannelQueue) PutWithTimeout(val interface{}, timeout time.Duration) error {
	select {
	case q <- val:
		return nil
	case <-time.After(timeout):
		return ErrQueuePutTimeout
	}
}

// Take Take the val(blocking)
func (q ChannelQueue) Take() (interface{}, error) {
	val, ok := <-q
	if !ok {
		return *new(interface{}), ErrQueueIsClosed
	}
	return val, nil
}

// TakeWithTimeout Take the val(blocking), with timeout
func (q ChannelQueue) TakeWithTimeout(timeout time.Duration) (interface{}, error) {
	select {
	case val := <-q:
		return val, nil
	case <-time.After(timeout):
		return *new(interface{}), ErrQueueTakeTimeout
	}
}

// Offer Offer the val(non-blocking)
func (q ChannelQueue) Offer(val interface{}) error {
	select {
	case q <- val:
		return nil
	default:
		return ErrQueueIsFull
	}
}

// Poll Poll the val(non-blocking)
func (q ChannelQueue) Poll() (interface{}, error) {
	select {
	case val := <-q:
		return val, nil
	default:
		return *new(interface{}), ErrQueueIsEmpty
	}
}

// LinkedList & DoublyLinkedList

// LinkedListItem LinkedListItem inspired by Collection utils
type LinkedListItem struct {
	Next *LinkedListItem

	Val *interface{}
}

// Count Count items
func (listItem *LinkedListItem) Count() int {
	count := 1
	first := listItem
	for first.Next != nil {
		count++
		first = first.Next
	}
	return count
}

// Last Get the Last one
func (listItem *LinkedListItem) Last() *LinkedListItem {
	last := listItem
	for last.Next != nil {
		last = last.Next
	}
	return last
}

// AddLast Add the input item to the last
func (listItem *LinkedListItem) AddLast(input *LinkedListItem) *LinkedListItem {
	last := listItem.Last()
	last.Next = input
	return last
}

// DoublyListItem DoublyListItem inspired by Collection utils
type DoublyListItem struct {
	Next *DoublyListItem
	Prev *DoublyListItem

	Val *interface{}
}

// Count Count items
func (listItem *DoublyListItem) Count() int {
	count := 1
	first := listItem.First()
	for first.Next != nil {
		count++
		first = first.Next
	}
	return count
}

// Last Get the Last one
func (listItem *DoublyListItem) Last() *DoublyListItem {
	last := listItem
	for last.Next != nil {
		last = last.Next
	}
	return last
}

// First Get the First one
func (listItem *DoublyListItem) First() *DoublyListItem {
	first := listItem
	for first.Prev != nil {
		first = first.Prev
	}
	return first
}

// AddLast Add the input item to the last
func (listItem *DoublyListItem) AddLast(input *DoublyListItem) *DoublyListItem {
	last := listItem.Last()
	first := input.First()
	last.Next = first
	first.Prev = last
	return last
}

// AddFirst Add the input item to the first position
func (listItem *DoublyListItem) AddFirst(input *DoublyListItem) *DoublyListItem {
	last := input.Last()
	first := listItem.First()
	last.Next = first
	first.Prev = last
	return first
}

// LinkedListQueue LinkedListQueue inspired by Collection utils
type LinkedListQueue struct {
	first *DoublyListItem
	last  *DoublyListItem
	count int

	nodePoolFirst *DoublyListItem
	nodeCount     int
}

// NewLinkedListQueue New LinkedListQueue instance
func NewLinkedListQueue() *LinkedListQueue {
	return new(LinkedListQueue)
}

// Count Count Items
func (q *LinkedListQueue) Count() int {
	return q.count
}

// ClearNodePool Clear cached LinkedListItem nodes in nodePool
func (q *LinkedListQueue) ClearNodePool() {
	q.nodeCount = 0
	q.nodePoolFirst = nil
}

// KeepNodePoolCount Decrease/Increase LinkedListItem nodes to n items
func (q *LinkedListQueue) KeepNodePoolCount(n int) {
	if n <= 0 {
		q.ClearNodePool()
		return
	}

	q.nodeCount = n

	n--
	last := q.nodePoolFirst
	if last == nil {
		last = new(DoublyListItem)
		q.nodePoolFirst = last
	}

	for n > 0 {
		n--
		if last.Next == nil {
			last.Next = new(DoublyListItem)
		}
		last = last.Next
	}
	last.Next = nil
}

// CLear Clear all data
func (q *LinkedListQueue) Clear() {
	q.nodePoolFirst = q.first
	q.nodeCount = q.count
	node := q.nodePoolFirst
	for node != nil {
		node.Val = nil
		node.Prev = nil
		node = node.Next
	}

	q.first = nil
	q.last = nil
	q.count = 0
}

// Put Put the val to the last position(no-blocking)
func (q *LinkedListQueue) Put(val interface{}) error {
	return q.Offer(val)
}

// Take Take the val from the first position(no-blocking)
func (q *LinkedListQueue) Take() (interface{}, error) {
	return q.Poll()
}

// Offer Offer the val to the last position(non-blocking)
func (q *LinkedListQueue) Offer(val interface{}) error {
	// Try get from pool or new one
	node := q.generateNode()
	node.Val = &val

	q.count++
	if q.first == nil {
		q.first = node
	}
	last := q.last
	if last != nil {
		last.Next = node
		node.Prev = last
	}
	q.last = node

	return nil
}

// Poll Poll the val from the first position(non-blocking)
func (q *LinkedListQueue) Poll() (interface{}, error) {
	return q.Shift()
}

// Peek Peek the val from the first position without removing it (non-blocking)
func (q *LinkedListQueue) Peek() (interface{}, error) {
	node := q.first
	if node == nil {
		return *new(interface{}), ErrQueueIsEmpty
	}
	return *node.Val, nil
}

// Shift Shift the val from the first position (non-blocking)
func (q *LinkedListQueue) Shift() (interface{}, error) {
	node := q.first
	if node == nil {
		return *new(interface{}), ErrQueueIsEmpty
	}

	// Remove the first item
	q.count--
	q.first = node.Next
	if q.first == nil {
		q.last = nil
	}
	val := node.Val

	q.recycleNode(node)

	return *val, nil
}

// Unshift Unshift the val to the first position(non-blocking)
func (q *LinkedListQueue) Unshift(val interface{}) error {
	// Try get from pool or new one
	node := q.generateNode()
	node.Val = &val

	q.count++
	if q.last == nil {
		q.last = node
	}
	first := q.first
	q.first = node
	node.Next = first
	if first != nil {
		first.Prev = node
	}

	return nil
}

// Pop Pop the data from the last position(non-blocking)
func (q *LinkedListQueue) Pop() (interface{}, error) {
	node := q.last
	if node == nil {
		return *new(interface{}), ErrStackIsEmpty
	}

	// Remove the last item
	q.count--
	q.last = node.Prev
	if q.last == nil {
		q.first = nil
	}
	val := node.Val
	q.recycleNode(node)

	return *val, nil
}

// Push Push the data to the last position(non-blocking)
func (q *LinkedListQueue) Push(val interface{}) error {
	// return q.Unshift(val)
	return q.Offer(val)
}

func (q *LinkedListQueue) generateNode() *DoublyListItem {
	node := q.nodePoolFirst
	if node == nil {
		node = new(DoublyListItem)
	} else {
		q.nodeCount--
		q.nodePoolFirst = node.Next
		node.Next = nil
		node.Prev = nil
	}

	return node
}

func (q *LinkedListQueue) recycleNode(node *DoublyListItem) {
	if node == nil {
		return
	}

	// Recycle
	q.nodeCount++
	node.Val = nil
	node.Next = q.nodePoolFirst
	node.Prev = nil
	q.nodePoolFirst = node
}

// BufferedChannelQueue BlockingQueue with ChannelQueue & scalable pool, inspired by Collection utils
type BufferedChannelQueue struct {
	lock     sync.RWMutex
	isClosed AtomBool

	loadWorkerCh     ChannelQueue
	freeNodeWorkerCh ChannelQueue

	loadFromPoolDuration             time.Duration
	freeNodeHookPoolIntervalDuration time.Duration
	nodeHookPoolSize                 int
	bufferSizeMaximum                int

	blockingQueue ChannelQueue
	pool          *LinkedListQueue
}

// NewBufferedChannelQueue New BufferedChannelQueue instance from a Queue
func NewBufferedChannelQueue(channelCapacity int, bufferSizeMaximum int, nodeHookPoolSize int) *BufferedChannelQueue {
	pool := NewLinkedListQueue()

	newOne := &BufferedChannelQueue{
		loadWorkerCh:     NewChannelQueue(1),
		freeNodeWorkerCh: NewChannelQueue(1),

		blockingQueue: NewChannelQueue(channelCapacity),
		pool:          pool,

		loadFromPoolDuration:             10 * time.Millisecond,
		freeNodeHookPoolIntervalDuration: 10 * time.Millisecond,

		nodeHookPoolSize:  nodeHookPoolSize,
		bufferSizeMaximum: bufferSizeMaximum,
	}
	go newOne.freeNodePool()
	go newOne.loadFromPool()

	return newOne
}

func (q *BufferedChannelQueue) freeNodePool() {
	for range q.freeNodeWorkerCh {
		time.Sleep(q.freeNodeHookPoolIntervalDuration)

		if q.isClosed.Get() {
			break
		}

		if q.pool.nodeCount > q.nodeHookPoolSize {
			q.lock.Lock()
			q.pool.KeepNodePoolCount(q.nodeHookPoolSize)
			q.lock.Unlock()
		}
	}
}

func (q *BufferedChannelQueue) loadFromPool() {
	for range q.loadWorkerCh {

		if q.isClosed.Get() {
			break
		}

		q.lock.Lock()

		var val interface{}
		var pollErr, offerErr error

		for q.pool.Count() > 0 {
			// Try poll from the pool
			val, pollErr = q.pool.Poll()
			if pollErr != nil {
				break
			}

			offerErr = q.blockingQueue.Offer(val)
			// If failed, unshift it back
			if offerErr != nil {
				q.pool.Unshift(val)
				break
			}
		}
		q.lock.Unlock()

		time.Sleep(q.loadFromPoolDuration)

	}
}

func (q *BufferedChannelQueue) notifyWorkers() {
	q.lock.RLock()
	if q.pool.Count() > 0 {
		q.loadWorkerCh.Offer(1)
	}
	if q.pool.nodeCount > q.nodeHookPoolSize {
		q.freeNodeWorkerCh.Offer(1)
	}
	q.lock.RUnlock()
}

// SetBufferSizeMaximum Set MaximumBufferSize(maximum number of buffered items outside the ChannelQueue)
func (q *BufferedChannelQueue) SetBufferSizeMaximum(size int) *BufferedChannelQueue {
	q.bufferSizeMaximum = size
	return q
}

// SetNodeHookPoolSize Set nodeHookPoolSize(the buffering node hooks ideal size)
func (q *BufferedChannelQueue) SetNodeHookPoolSize(size int) *BufferedChannelQueue {
	q.nodeHookPoolSize = size
	return q
}

// SetLoadFromPoolDuration Set loadFromPoolDuration(the interval to take buffered items into the ChannelQueue)
func (q *BufferedChannelQueue) SetLoadFromPoolDuration(duration time.Duration) *BufferedChannelQueue {
	q.loadFromPoolDuration = duration
	return q
}

// SetFreeNodeHookPoolIntervalDuration Set freeNodeHookPoolIntervalDuration(the interval to clear buffering node hooks down to nodeHookPoolSize)
func (q *BufferedChannelQueue) SetFreeNodeHookPoolIntervalDuration(duration time.Duration) *BufferedChannelQueue {
	q.freeNodeHookPoolIntervalDuration = duration
	return q
}

// IsClosed Is the BufferedChannelQueue closed
func (q *BufferedChannelQueue) IsClosed() bool {
	return q.isClosed.Get()
}

// Close Close the BufferedChannelQueue
func (q *BufferedChannelQueue) Close() {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.isClosed.Set(true)
	close(q.loadWorkerCh)
	close(q.blockingQueue)
}

// Put Put the val(non-blocking)
func (q *BufferedChannelQueue) Put(val interface{}) error {
	// q.lock.Lock()
	// defer q.lock.Unlock()
	//
	// if q.isClosed.Get() {
	// 	return ErrQueueIsClosed
	// }
	//
	// q.lock.Lock()
	// poolCount := q.pool.Count()
	//
	// // If appearing nothing in the pool
	// if poolCount == 0 {
	// 	defer q.lock.Unlock()
	// 	// Try channel
	// 	return q.blockingQueue.Put(val)
	// }
	// q.lock.Unlock()

	return q.Offer(val)
}

// // PutWithTimeout Put the val(blocking), with timeout
// func (q BufferedChannelQueue) PutWithTimeout(val interface{}, timeout time.Duration) error {
// //  	q.lock.Lock()
// //  	defer q.lock.Unlock()
//
// 	if q.isClosed.Get() {
// 		return ErrQueueIsClosed
// 	}
//
// 	return q.blockingQueue.PutWithTimeout(val, timeout)
// }

// Take Take the val(blocking)
func (q *BufferedChannelQueue) Take() (interface{}, error) {
	// q.lock.RLock()
	// defer q.lock.RUnlock()

	if q.isClosed.Get() {
		return *new(interface{}), ErrQueueIsClosed
	}

	q.notifyWorkers()

	return q.blockingQueue.Take()
}

// TakeWithTimeout Take the val(blocking), with timeout
func (q *BufferedChannelQueue) TakeWithTimeout(timeout time.Duration) (interface{}, error) {
	// q.lock.RLock()
	// defer q.lock.RUnlock()

	if q.isClosed.Get() {
		return *new(interface{}), ErrQueueIsClosed
	}

	q.notifyWorkers()

	return q.blockingQueue.TakeWithTimeout(timeout)
}

// Offer Offer the val(non-blocking)
func (q *BufferedChannelQueue) Offer(val interface{}) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.isClosed.Get() {
		return ErrQueueIsClosed
	}

	poolCount := q.pool.Count()

	// If appearing nothing in the pool
	if poolCount == 0 {
		// Try channel
		err := q.blockingQueue.Offer(val)
		if err == nil {
			// Success
			return nil
		} else if err == ErrQueueIsFull {
			// Do nothing and let pool.Offer(val)
		} else {
			// Other
			return err
		}
	}

	// Before +1: >=, After +1: >
	if poolCount >= q.bufferSizeMaximum {
		return ErrQueueIsFull
	}

	q.pool.Offer(val)
	q.loadWorkerCh.Offer(1)
	return nil
}

// Poll Poll the val(non-blocking)
func (q *BufferedChannelQueue) Poll() (interface{}, error) {
	// q.lock.RLock()
	// defer q.lock.RUnlock()

	if q.isClosed.Get() {
		return *new(interface{}), ErrQueueIsClosed
	}

	q.notifyWorkers()

	return q.blockingQueue.Poll()
}
