package fpgo

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Regression: ConcurrentQueue.Take/Poll and ConcurrentStack.Pop are MUTATING
// operations (they remove an element from the wrapped, non-thread-safe
// container). They previously took only a read lock (RLock), so concurrent
// removers ran the underlying mutation simultaneously -> data race and lost/
// duplicated/dropped elements. They now take the write lock. These tests
// exercise concurrent producers/consumers and assert no element is lost or
// duplicated; run under `-race` they also catch the lock regression directly.

func TestConcurrentQueueNoLossUnderConcurrency(t *testing.T) {
	base := NewLinkedListQueue[int]()
	cq := NewConcurrentQueue[int](base)

	const total = 4000
	for i := 0; i < total; i++ {
		assert.NoError(t, cq.Offer(i))
	}

	var mu sync.Mutex
	seen := make(map[int]int, total)
	got := 0

	var wg sync.WaitGroup
	const workers = 8
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				v, err := cq.Poll()
				if err != nil {
					return // empty
				}
				mu.Lock()
				seen[v]++
				got++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	assert.Equal(t, total, got, "every offered element polled exactly once")
	assert.Equal(t, total, len(seen), "no duplicates, no losses")
	for k, c := range seen {
		assert.Equalf(t, 1, c, "value %d polled %d times", k, c)
	}
}

func TestConcurrentStackNoLossUnderConcurrency(t *testing.T) {
	base := NewLinkedListQueue[int]()
	cs := NewConcurrentStack[int](base)

	const total = 4000
	for i := 0; i < total; i++ {
		assert.NoError(t, cs.Push(i))
	}

	var mu sync.Mutex
	seen := make(map[int]int, total)
	got := 0

	var wg sync.WaitGroup
	const workers = 8
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				v, err := cs.Pop()
				if err != nil {
					return
				}
				mu.Lock()
				seen[v]++
				got++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	assert.Equal(t, total, got, "every pushed element popped exactly once")
	assert.Equal(t, total, len(seen), "no duplicates, no losses")
}

// Concurrent producers and consumers interleaved: total conservation.
func TestConcurrentQueueMixedProducersConsumers(t *testing.T) {
	base := NewLinkedListQueue[int]()
	cq := NewConcurrentQueue[int](base)

	const perProducer = 1000
	const producers = 4
	var prodWg sync.WaitGroup
	for p := 0; p < producers; p++ {
		prodWg.Add(1)
		go func(base int) {
			defer prodWg.Done()
			for i := 0; i < perProducer; i++ {
				_ = cq.Offer(base + i)
			}
		}(p * perProducer)
	}

	var mu sync.Mutex
	got := 0
	stop := make(chan struct{})
	var consWg sync.WaitGroup
	for c := 0; c < 4; c++ {
		consWg.Add(1)
		go func() {
			defer consWg.Done()
			for {
				v, err := cq.Poll()
				if err == nil {
					_ = v
					mu.Lock()
					got++
					mu.Unlock()
					continue
				}
				select {
				case <-stop:
					// Drain any stragglers before exiting.
					if _, err := cq.Poll(); err == nil {
						mu.Lock()
						got++
						mu.Unlock()
						continue
					}
					return
				default:
				}
			}
		}()
	}

	prodWg.Wait()
	close(stop)
	consWg.Wait()

	// Drain whatever consumers may have missed after stop.
	for {
		_, err := cq.Poll()
		if err != nil {
			break
		}
		got++
	}

	assert.Equal(t, producers*perProducer, got, "all produced items consumed exactly once")
}
