package fpgo

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandler_New(t *testing.T) {
	handler := new(HandlerDef).New()
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.ch)
	handler.Close()
}

func TestHandler_NewByCh(t *testing.T) {
	ch := make(chan func())
	handler := new(HandlerDef).NewByCh(ch)
	assert.NotNil(t, handler)
	assert.Equal(t, ch, handler.ch)
	assert.False(t, handler.isClosed)
	handler.Close()
}

func TestHandler_Post(t *testing.T) {
	handler := new(HandlerDef).New()
	defer handler.Close()

	var wg sync.WaitGroup
	var result int

	wg.Add(1)
	handler.Post(func() {
		result = 42
		wg.Done()
	})

	wg.Wait()
	assert.Equal(t, 42, result)
}

func TestHandler_PostMultiple(t *testing.T) {
	handler := new(HandlerDef).New()
	defer handler.Close()

	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]int, 0, 3)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		handler.Post(func() {
			mu.Lock()
			results = append(results, i)
			mu.Unlock()
			wg.Done()
		})
	}

	wg.Wait()
	assert.Len(t, results, 3)
}

func TestHandler_PostAfterClose(t *testing.T) {
	handler := new(HandlerDef).New()
	handler.Close()

	handler.Post(func() {})
	handler.Post(func() {})
}

func TestHandler_Close(t *testing.T) {
	handler := new(HandlerDef).New()
	assert.False(t, handler.isClosed)

	handler.Close()
	assert.True(t, handler.isClosed)
}

func TestHandler_CloseMultiple(t *testing.T) {
	handler := new(HandlerDef).New()
	handler.Close()
}

func TestGetDefault(t *testing.T) {
	handler := new(HandlerDef).GetDefault()
	assert.NotNil(t, handler)
}

func TestHandlerIsClosed(t *testing.T) {
	h := new(HandlerDef).New()
	assert.Equal(t, false, h.IsClosed())

	h.Close()
	assert.Equal(t, true, h.IsClosed())
}

func TestHandlerPostAfterClose(t *testing.T) {
	h := new(HandlerDef).New()
	h.Close()

	h.Post(func() {})
	h.Post(func() {})
}

func TestHandlerExecutionOrder(t *testing.T) {
	handler := new(HandlerDef).New()
	defer handler.Close()

	var mu sync.Mutex
	order := make([]int, 0, 5)

	handler.Post(func() {
		mu.Lock()
		order = append(order, 1)
		mu.Unlock()
	})
	handler.Post(func() {
		mu.Lock()
		order = append(order, 2)
		mu.Unlock()
	})
	handler.Post(func() {
		mu.Lock()
		order = append(order, 3)
		mu.Unlock()
	})

	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	assert.Equal(t, []int{1, 2, 3}, order)
	mu.Unlock()
}

func TestHandler_ConcurrentPost(t *testing.T) {
	handler := new(HandlerDef).New()
	defer handler.Close()

	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]int, 0, 100)

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			handler.Post(func() {
				mu.Lock()
				results = append(results, val)
				mu.Unlock()
				wg.Done()
			})
		}(i)
	}

	wg.Wait()
	assert.Len(t, results, 100)
}

func TestHandlerNewByCh(t *testing.T) {
	ch := make(chan func(), 10)
	handler := new(HandlerDef).NewByCh(ch)
	assert.NotNil(t, handler)
	assert.Equal(t, ch, handler.ch)
	handler.Close()
}
