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

	queue = NewChannelQueue[int](3)

	err = queue.Offer(1)
	assert.Equal(t, err, nil)
	err = queue.Offer(2)
	assert.Equal(t, err, nil)
	err = queue.Offer(3)
	assert.Equal(t, err, nil)
	err = queue.Offer(4)
	assert.NotEqual(t, err, nil)
	assert.Equal(t, err, ErrQueueIsFull)

	result, err = queue.Poll()
	assert.Equal(t, result, 1)
	assert.Equal(t, err, nil)
	result, err = queue.Poll()
	assert.Equal(t, result, 2)
	assert.Equal(t, err, nil)
	result, err = queue.Poll()
	assert.Equal(t, result, 3)
	assert.Equal(t, err, nil)
	result, err = queue.Poll()
	assert.NotEqual(t, result, 4)
	assert.NotEqual(t, err, nil)
	assert.Equal(t, result, 0)
	assert.Equal(t, err, ErrQueueIsEmpty)

	result = 0
	timeout = 1 * time.Millisecond
	go func() {
		result, err = queue.Take(&timeout)
		assert.Equal(t, err, nil)
		assert.Equal(t, result, 1)
		result, err = queue.Take(&timeout)
		assert.Equal(t, err, nil)
		assert.Equal(t, result, 2)
		result, err = queue.Take(&timeout)
		assert.Equal(t, err, nil)
		assert.Equal(t, result, 3)
		result, err = queue.Take(&timeout)
		assert.NotEqual(t, result, 4)
		assert.NotEqual(t, err, nil)
		assert.Equal(t, result, 0)
		assert.Equal(t, err, ErrQueueTakeTimeout)
	}()
	go func() {
		err = queue.Put(1, &timeout)
		assert.Equal(t, err, nil)
		err = queue.Put(2, &timeout)
		assert.Equal(t, err, nil)
		err = queue.Put(3, &timeout)
		assert.Equal(t, err, nil)

		time.Sleep(3 * timeout / 2)

		err = queue.Put(4, &timeout)
		assert.Equal(t, err, nil)
		err = queue.Put(5, &timeout)
		assert.Equal(t, err, nil)
		err = queue.Put(6, &timeout)
		assert.Equal(t, err, nil)
		err = queue.Put(7, &timeout)
		assert.NotEqual(t, err, nil)
		assert.Equal(t, err, ErrQueuePutTimeout)
	}()

	time.Sleep(2 * timeout)
}
