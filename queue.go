package fpgo

import (
	"errors"
	"time"
)

type Queue[T any] interface {
	Put(val T, timeout *time.Duration) error
	Take(timeout *time.Duration) (T, error)
	Offer(val T) error
	Poll() (T, error)
}

var (
	ErrQueueIsEmpty     = errors.New("queue is empty")
	ErrQueueIsFull      = errors.New("queue is full")
	ErrQueueTakeTimeout = errors.New("queue take timeout")
	ErrQueuePutTimeout  = errors.New("queue put timeout")
)

type ChannelQueue[T any] struct {
	channel chan T
}

func NewChannelQueue[T any](capacity int) *ChannelQueue[T] {
	return &ChannelQueue[T]{
		channel: make(chan T, capacity),
	}
}

func (q *ChannelQueue[T]) Put(val T, timeout *time.Duration) error {
	if timeout == nil {
		q.channel <- val
		return nil
	}

	select {
	case q.channel <- val:
		return nil
	case <-time.After(*timeout):
		return ErrQueuePutTimeout
	}
}

func (q *ChannelQueue[T]) Take(timeout *time.Duration) (T, error) {
	if timeout == nil {
		val := <-q.channel
		return val, nil
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
