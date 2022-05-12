package fpgo

import (
	"errors"
	"sync"
	"time"
)

// Queue Queue inspired by Collection utils
type Queue[T any] interface {
	Put(val T) error
	Take() (T, error)
	Offer(val T) error
	Poll() (T, error)
}

// Stack Stack inspired by Collection utils
type Stack[T any] interface {
	Push(val T) error
	Pop() (T, error)
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

// ConcurrentQueue

// ConcurrentQueue ConcurrentQueue inspired by Collection utils
type ConcurrentQueue[T any] struct {
	lock  sync.RWMutex
	queue Queue[T]
}

// NewConcurrentQueue New ConcurrentQueue instance from a Queue[T]
func NewConcurrentQueue[T any](queue Queue[T]) *ConcurrentQueue[T] {
	return &ConcurrentQueue[T]{
		queue: queue,
	}
}

// Put Put the T val(probably blocking)
func (q *ConcurrentQueue[T]) Put(val T) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.queue.Put(val)
}

// Take Take the T val(probably blocking)
func (q *ConcurrentQueue[T]) Take() (T, error) {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return q.queue.Take()
}

// Offer Offer the T val(non-blocking)
func (q *ConcurrentQueue[T]) Offer(val T) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.queue.Offer(val)
}

// Poll Poll the T val(non-blocking)
func (q *ConcurrentQueue[T]) Poll() (T, error) {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return q.queue.Poll()
}

// ChannelQueue

// ChannelQueue ChannelQueue inspired by Collection utils
type ChannelQueue[T any] chan T

// NewChannelQueue New ChannelQueue instance with capacity
func NewChannelQueue[T any](capacity int) ChannelQueue[T] {
	return make(ChannelQueue[T], capacity)
}

// Put Put the T val(blocking)
func (q ChannelQueue[T]) Put(val T) error {
	q <- val
	return nil
}

// PutWithTimeout Put the T val(blocking), with timeout
func (q ChannelQueue[T]) PutWithTimeout(val T, timeout time.Duration) error {
	select {
	case q <- val:
		return nil
	case <-time.After(timeout):
		return ErrQueuePutTimeout
	}
}

// Take Take the T val(blocking)
func (q ChannelQueue[T]) Take() (T, error) {
	val, ok := <-q
	if !ok {
		return *new(T), ErrQueueIsClosed
	}
	return val, nil
}

// TakeWithTimeout Take the T val(blocking), with timeout
func (q ChannelQueue[T]) TakeWithTimeout(timeout time.Duration) (T, error) {
	select {
	case val := <-q:
		return val, nil
	case <-time.After(timeout):
		return *new(T), ErrQueueTakeTimeout
	}
}

// Offer Offer the T val(non-blocking)
func (q ChannelQueue[T]) Offer(val T) error {
	select {
	case q <- val:
		return nil
	default:
		return ErrQueueIsFull
	}
}

// Poll Poll the T val(non-blocking)
func (q ChannelQueue[T]) Poll() (T, error) {
	select {
	case val := <-q:
		return val, nil
	default:
		return *new(T), ErrQueueIsEmpty
	}
}

// LinkedList & DoublyLinkedList

// LinkedListItem LinkedListItem inspired by Collection utils
type LinkedListItem[T any] struct {
	Next *LinkedListItem[T]

	Val *T
}

// Count Count items
func (listItem *LinkedListItem[T]) Count() int {
	count := 1
	first := listItem
	for first.Next != nil {
		count++
		first = first.Next
	}
	return count
}

// Last Get the Last one
func (listItem *LinkedListItem[T]) Last() *LinkedListItem[T] {
	last := listItem
	for last.Next != nil {
		last = last.Next
	}
	return last
}

// AddLast Add the input item to the last
func (listItem *LinkedListItem[T]) AddLast(input *LinkedListItem[T]) *LinkedListItem[T] {
	last := listItem.Last()
	last.Next = input
	return last
}

// DoublyListItem DoublyListItem inspired by Collection utils
type DoublyListItem[T any] struct {
	Next *DoublyListItem[T]
	Prev *DoublyListItem[T]

	Val *T
}

// Count Count items
func (listItem *DoublyListItem[T]) Count() int {
	count := 1
	first := listItem.First()
	for first.Next != nil {
		count++
		first = first.Next
	}
	return count
}

// Last Get the Last one
func (listItem *DoublyListItem[T]) Last() *DoublyListItem[T] {
	last := listItem
	for last.Next != nil {
		last = last.Next
	}
	return last
}

// First Get the First one
func (listItem *DoublyListItem[T]) First() *DoublyListItem[T] {
	first := listItem
	for first.Prev != nil {
		first = first.Prev
	}
	return first
}

// AddLast Add the input item to the last
func (listItem *DoublyListItem[T]) AddLast(input *DoublyListItem[T]) *DoublyListItem[T] {
	last := listItem.Last()
	first := input.First()
	last.Next = first
	first.Prev = last
	return last
}

// AddFirst Add the input item to the first position
func (listItem *DoublyListItem[T]) AddFirst(input *DoublyListItem[T]) *DoublyListItem[T] {
	last := input.Last()
	first := listItem.First()
	last.Next = first
	first.Prev = last
	return first
}

// LinkedListQueue LinkedListQueue inspired by Collection utils
type LinkedListQueue[T any] struct {
	first *DoublyListItem[T]
	last  *DoublyListItem[T]
	count int

	nodePoolFirst *DoublyListItem[T]
	nodeCount     int
}

// NewLinkedListQueue New LinkedListQueue instance
func NewLinkedListQueue[T any]() *LinkedListQueue[T] {
	return new(LinkedListQueue[T])
}

// Count Count Items
func (q *LinkedListQueue[T]) Count() int {
	return q.count
}

// ClearNodePool Clear cached LinkedListItem nodes in nodePool
func (q *LinkedListQueue[T]) ClearNodePool() {
	q.nodeCount = 0
	q.nodePoolFirst = nil
}

// KeepNodePoolCount Decrease/Increase LinkedListItem nodes to n items
func (q *LinkedListQueue[T]) KeepNodePoolCount(n int) {
	if n <= 0 {
		q.ClearNodePool()
		return
	}

	q.nodeCount = n

	n--
	last := q.nodePoolFirst
	if last == nil {
		last = new(DoublyListItem[T])
		q.nodePoolFirst = last
	}

	for n > 0 {
		n--
		if last.Next == nil {
			last.Next = new(DoublyListItem[T])
		}
		last = last.Next
	}
	last.Next = nil
}

// CLear Clear all data
func (q *LinkedListQueue[T]) Clear() {
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

// Put Put the T val to the last position(no-blocking)
func (q *LinkedListQueue[T]) Put(val T) error {
	return q.Offer(val)
}

// Take Take the T val from the first position(no-blocking)
func (q *LinkedListQueue[T]) Take() (T, error) {
	return q.Poll()
}

// Offer Offer the T val to the last position(non-blocking)
func (q *LinkedListQueue[T]) Offer(val T) error {
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

// Poll Poll the T val from the first position(non-blocking)
func (q *LinkedListQueue[T]) Poll() (T, error) {
	return q.Shift()
}

// Peek Peek the T val from the first position without removing it (non-blocking)
func (q *LinkedListQueue[T]) Peek() (T, error) {
	node := q.first
	if node == nil {
		return *new(T), ErrQueueIsEmpty
	}
	return *node.Val, nil
}

// Shift Shift the T val from the first position (non-blocking)
func (q *LinkedListQueue[T]) Shift() (T, error) {
	node := q.first
	if node == nil {
		return *new(T), ErrQueueIsEmpty
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

// Unshift Unshift the T val to the first position(non-blocking)
func (q *LinkedListQueue[T]) Unshift(val T) error {
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
func (q *LinkedListQueue[T]) Pop() (T, error) {
	node := q.last
	if node == nil {
		return *new(T), ErrStackIsEmpty
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
func (q *LinkedListQueue[T]) Push(val T) error {
	// return q.Unshift(val)
	return q.Offer(val)
}

func (q *LinkedListQueue[T]) generateNode() *DoublyListItem[T] {
	node := q.nodePoolFirst
	if node == nil {
		node = new(DoublyListItem[T])
	} else {
		q.nodeCount--
		q.nodePoolFirst = node.Next
		node.Next = nil
		node.Prev = nil
	}

	return node
}

func (q *LinkedListQueue[T]) recycleNode(node *DoublyListItem[T]) {
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
type BufferedChannelQueue[T any] struct {
	lock     sync.RWMutex
	isClosed AtomBool

	loadWorkerCh     ChannelQueue[int]
	freeNodeWorkerCh ChannelQueue[int]

	loadFromPoolDuration             time.Duration
	freeNodeHookPoolIntervalDuration time.Duration
	nodeHookPoolSize                 int
	bufferSizeMaximum                int

	blockingQueue ChannelQueue[T]
	pool          *LinkedListQueue[T]
}

// NewBufferedChannelQueue New BufferedChannelQueue instance from a Queue[T]
func NewBufferedChannelQueue[T any](channelCapacity int, bufferSizeMaximum int, nodeHookPoolSize int) *BufferedChannelQueue[T] {
	pool := NewLinkedListQueue[T]()

	newOne := &BufferedChannelQueue[T]{
		loadWorkerCh:     NewChannelQueue[int](1),
		freeNodeWorkerCh: NewChannelQueue[int](1),

		blockingQueue: NewChannelQueue[T](channelCapacity),
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

func (q *BufferedChannelQueue[T]) freeNodePool() {
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

func (q *BufferedChannelQueue[T]) loadFromPool() {
	for range q.loadWorkerCh {

		if q.isClosed.Get() {
			break
		}

		q.lock.Lock()

		var val T
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

func (q *BufferedChannelQueue[T]) notifyWorkers() {
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
func (q *BufferedChannelQueue[T]) SetBufferSizeMaximum(size int) {
	q.bufferSizeMaximum = size
}

// SetNodeHookPoolSize Set nodeHookPoolSize(the buffering node hooks ideal size)
func (q *BufferedChannelQueue[T]) SetNodeHookPoolSize(size int) {
	q.nodeHookPoolSize = size
}

// SetLoadFromPoolDuration Set loadFromPoolDuration(the interval to take buffered items into the ChannelQueue)
func (q *BufferedChannelQueue[T]) SetLoadFromPoolDuration(duration time.Duration) {
	q.loadFromPoolDuration = duration
}

// SetFreeNodeHookPoolIntervalDuration Set freeNodeHookPoolIntervalDuration(the interval to clear buffering node hooks down to nodeHookPoolSize)
func (q *BufferedChannelQueue[T]) SetFreeNodeHookPoolIntervalDuration(duration time.Duration) {
	q.freeNodeHookPoolIntervalDuration = duration
}

// Close Close the Handler
func (q *BufferedChannelQueue[T]) Close() {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.isClosed.Set(true)
	close(q.loadWorkerCh)
	close(q.blockingQueue)
}

// Put Put the T val(non-blocking)
func (q *BufferedChannelQueue[T]) Put(val T) error {
	// q.lock.Lock()
	// defer q.lock.Unlock()
	//
	// if q.isClosed.Get() {
	// 	return ErrQueueIsClosed
	// }
	//
	// return q.blockingQueue.Put(val)

	return q.Offer(val)
}

// // PutWithTimeout Put the T val(blocking), with timeout
// func (q BufferedChannelQueue[T]) PutWithTimeout(val T, timeout time.Duration) error {
// //  	q.lock.Lock()
// //  	defer q.lock.Unlock()
//
// 	if q.isClosed.Get() {
// 		return ErrQueueIsClosed
// 	}
//
// 	return q.blockingQueue.PutWithTimeout(val, timeout)
// }

// Take Take the T val(blocking)
func (q *BufferedChannelQueue[T]) Take() (T, error) {
	// q.lock.RLock()
	// defer q.lock.RUnlock()

	if q.isClosed.Get() {
		return *new(T), ErrQueueIsClosed
	}

	q.notifyWorkers()

	return q.blockingQueue.Take()
}

// TakeWithTimeout Take the T val(blocking), with timeout
func (q *BufferedChannelQueue[T]) TakeWithTimeout(timeout time.Duration) (T, error) {
	// q.lock.RLock()
	// defer q.lock.RUnlock()

	if q.isClosed.Get() {
		return *new(T), ErrQueueIsClosed
	}

	q.notifyWorkers()

	return q.blockingQueue.TakeWithTimeout(timeout)
}

// Offer Offer the T val(non-blocking)
func (q *BufferedChannelQueue[T]) Offer(val T) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.isClosed.Get() {
		return ErrQueueIsClosed
	}

	// Before +1: >=, After +1: >
	if q.pool.Count() >= q.bufferSizeMaximum {
		return ErrQueueIsFull
	}

	q.pool.Offer(val)
	q.loadWorkerCh.Offer(1)
	return nil

	// err := q.blockingQueue.Offer(val)
	// if err == ErrQueueIsFull {
	// 	q.pool.Offer(val)
	// 	q.loadWorkerCh.Offer(1)
	// 	return nil
	// }
	//
	// return err
}

// Poll Poll the T val(non-blocking)
func (q *BufferedChannelQueue[T]) Poll() (T, error) {
	// q.lock.RLock()
	// defer q.lock.RUnlock()

	if q.isClosed.Get() {
		return *new(T), ErrQueueIsClosed
	}

	q.notifyWorkers()

	return q.blockingQueue.Poll()
}
