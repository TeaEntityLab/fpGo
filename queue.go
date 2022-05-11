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

var (
	ErrQueueIsEmpty     = errors.New("queue is empty")
	ErrQueueIsFull      = errors.New("queue is full")
	ErrQueueTakeTimeout = errors.New("queue take timeout")
	ErrQueuePutTimeout  = errors.New("queue put timeout")
)

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

func (q *ConcurrentQueue[T]) Put(val T) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.queue.Put(val)
}

func (q *ConcurrentQueue[T]) Take() (T, error) {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return q.queue.Take()
}

func (q *ConcurrentQueue[T]) Offer(val T) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.queue.Offer(val)
}

func (q *ConcurrentQueue[T]) Poll() (T, error) {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return q.queue.Poll()
}

type ChannelQueue[T any] struct {
	channel chan T
}

func NewChannelQueue[T any](capacity int) *ChannelQueue[T] {
	return &ChannelQueue[T]{
		channel: make(chan T, capacity),
	}
}

func (q *ChannelQueue[T]) Put(val T) error {
	q.channel <- val
	return nil
}

func (q *ChannelQueue[T]) PutWithTimeout(val T, timeout *time.Duration) error {
	if timeout == nil {
		return q.Put(val)
	}

	select {
	case q.channel <- val:
		return nil
	case <-time.After(*timeout):
		return ErrQueuePutTimeout
	}
}

func (q *ChannelQueue[T]) Take() (T, error) {
	val := <-q.channel
	return val, nil
}

func (q *ChannelQueue[T]) TakeWithTimeout(timeout *time.Duration) (T, error) {
	if timeout == nil {
		return q.Take()
	}

	select {
	case val := <-q.channel:
		return val, nil
	case <-time.After(*timeout):
		return *new(T), ErrQueueTakeTimeout
	}
}

func (q *ChannelQueue[T]) Offer(val T) error {
	select {
	case q.channel <- val:
		return nil
	default:
		return ErrQueueIsFull
	}
}

func (q *ChannelQueue[T]) Poll() (T, error) {
	select {
	case val := <-q.channel:
		return val, nil
	default:
		return *new(T), ErrQueueIsEmpty
	}
}

type LinkedListItem[T any] struct {
	Next *LinkedListItem[T]

	Val *T
}

func (listItem *LinkedListItem[T]) Count() int {
	count := 1
	first := listItem
	for first.Next != nil {
		count++
		first = first.Next
	}
	return count
}

func (listItem *LinkedListItem[T]) Last() *LinkedListItem[T] {
	last := listItem
	for last.Next != nil {
		last = last.Next
	}
	return last
}

func (listItem *LinkedListItem[T]) AddLast(input *LinkedListItem[T]) *LinkedListItem[T] {
	last := listItem.Last()
	last.Next = input
	return last
}

type DoublyListItem[T any] struct {
	Next *DoublyListItem[T]
	Prev *DoublyListItem[T]

	Val *T
}

func (listItem *DoublyListItem[T]) Count() int {
	count := 1
	first := listItem.First()
	for first.Next != nil {
		count++
		first = first.Next
	}
	return count
}

func (listItem *DoublyListItem[T]) Last() *DoublyListItem[T] {
	last := listItem
	for last.Next != nil {
		last = last.Next
	}
	return last
}

func (listItem *DoublyListItem[T]) First() *DoublyListItem[T] {
	first := listItem
	for first.Prev != nil {
		first = first.Prev
	}
	return first
}

func (listItem *DoublyListItem[T]) AddLast(input *DoublyListItem[T]) *DoublyListItem[T] {
	last := listItem.Last()
	first := input.First()
	last.Next = first
	first.Prev = last
	return last
}

func (listItem *DoublyListItem[T]) AddFirst(input *DoublyListItem[T]) *DoublyListItem[T] {
	last := input.Last()
	first := listItem.First()
	last.Next = first
	first.Prev = last
	return first
}

type LinkedListQueue[T any] struct {
	first *LinkedListItem[T]
	last  *LinkedListItem[T]
	count int

	nodePoolFirst *LinkedListItem[T]
	nodePoolLast  *LinkedListItem[T]
	nodeCount     int
}

// NewLinkedListQueue New LinkedListQueue instance
func NewLinkedListQueue[T any]() *LinkedListQueue[T] {
	return new(LinkedListQueue[T])
}

func (q *LinkedListQueue[T]) Count() int {
	return q.count
}

func (q *LinkedListQueue[T]) ClearNodePool() {
	q.nodeCount = 0
	q.nodePoolFirst = nil
	q.nodePoolLast = nil
}

func (q *LinkedListQueue[T]) KeepNodePoolCount(n int) {
  if n <= 0 {
    q.ClearNodePool()
    return
  }

  q.nodeCount = n

  n--
  last := q.nodePoolFirst
  if last == nil {
    last = new(LinkedListItem[T])
    q.nodePoolFirst = last
  }

  for n > 0 {
    n--
    if last.Next == nil {
      last.Next = new(LinkedListItem[T])
    }
    last = last.Next
  }
  last.Next = nil
}

func (q *LinkedListQueue[T]) Clear() {
	q.nodePoolFirst = q.first
	q.nodePoolLast = q.last
	q.nodeCount = q.count

	q.first = nil
	q.last = nil
	q.count = 0
}

func (q *LinkedListQueue[T]) Put(val T) error {
	return q.Offer(val)
}

func (q *LinkedListQueue[T]) Take() (T, error) {
	return q.Poll()
}

func (q *LinkedListQueue[T]) Offer(val T) error {
	// Try get from pool or new one
	node := q.nodePoolFirst
	if node == nil {
		node = new(LinkedListItem[T])
	} else {
		q.nodeCount--
		q.nodePoolFirst = node.Next
		if q.nodePoolFirst == nil {
			q.nodePoolLast = nil
		}
		node.Next = nil
	}
	node.Val = &val

	q.count++
	if q.first == nil {
		q.first = node
	}
	last := q.last
	if last != nil {
		last.Next = node
	}
	q.last = node

	return nil
}

func (q *LinkedListQueue[T]) Poll() (T, error) {
	node := q.first
	if node == nil {
		return *new(T), ErrQueueIsEmpty
	}
	q.count--
	q.first = node.Next
	if q.first == nil {
		q.last = nil
	}
	val := *node.Val

	// Recycle
	q.nodeCount++
	node.Val = nil
	node.Next = nil
	if q.nodePoolLast == nil {
		q.nodePoolFirst = node
	} else {
		q.nodePoolLast.Next = node
	}
	q.nodePoolLast = node

	return val, nil
}
