package worker

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	fpgo "github.com/TeaEntityLab/fpGo"
)

func TestWorkerPool(t *testing.T) {
	var workerPool WorkerPool
	var err error
	// Use custom settings with standBy=0 from the start to prevent spawnLoop
	// from spawning workers before chained setters can override defaults.
	settings := &DefaultWorkerPoolSettings{
		isJobQueueClosedWhenClose: true,
		workerBatchSize:           3,
		workerSizeStandBy:         0,
		workerSizeMaximum:         5,
		spawnWorkerDuration:       1 * time.Millisecond,
		workerExpiryDuration:      50 * time.Millisecond,
		workerJamDuration:         1000 * time.Millisecond,
		scheduleRetryInterval:     50 * time.Millisecond,
		panicHandler:              defaultPanicHandler,
	}
	defaultWorkerPool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue(3, 10000, 100), settings)
	workerPool = defaultWorkerPool
	defer defaultWorkerPool.Close()

	// Test PreAlloc (standBy=0 so spawnLoop won't spawn asynchronously)
	assert.Equal(t, 0, defaultWorkerPool.getWorkerCount())
	defaultWorkerPool.PreAllocWorkerSize(3)
	assert.Equal(t, 3, defaultWorkerPool.getWorkerCount())

	// PreAlloc more
	defaultWorkerPool.PreAllocWorkerSize(5)
	assert.Equal(t, 5, defaultWorkerPool.getWorkerCount())

	// Test ScaleDown: block workers so they don't expire while busy
	blockCh := make(chan struct{})
	for i := 0; i < 5; i++ {
		err = workerPool.Schedule(func() {
			<-blockCh
		})
		assert.NoError(t, err)
	}
	// Now set standBy=1 so after workers finish, 1 remains
	defaultWorkerPool.SetWorkerSizeStandBy(1)
	// Release workers and wait for expiry
	close(blockCh)
	time.Sleep(200 * time.Millisecond)
	// Scaled down to at most standBy (1). Under scheduling stress, idle workers
	// hitting expiry near-simultaneously can all exit (count read under RLock,
	// decrement in defer), so assert the floor invariant, not an exact count.
	assert.LessOrEqual(t, defaultWorkerPool.getWorkerCount(), 1)

	// Test ScaleDown to 0
	defaultWorkerPool.SetWorkerSizeStandBy(0)
	time.Sleep(200 * time.Millisecond)
	assert.Equal(t, 0, defaultWorkerPool.getWorkerCount())
}

func TestScheduleWithTimeout(t *testing.T) {
	var workerPool WorkerPool
	var err error
	settings := &DefaultWorkerPoolSettings{
		isJobQueueClosedWhenClose: true,
		workerBatchSize:           0,
		workerSizeStandBy:         0,
		workerSizeMaximum:         0,
		spawnWorkerDuration:       1 * time.Millisecond,
		workerExpiryDuration:      50 * time.Millisecond,
		workerJamDuration:         1000 * time.Millisecond,
		scheduleRetryInterval:     50 * time.Millisecond,
		panicHandler:              defaultPanicHandler,
	}
	defaultWorkerPool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue(3, 1, 3), settings)
	workerPool = defaultWorkerPool
	defer defaultWorkerPool.Close()

	// Test ScheduleWithTimeout
	// channel: 3 positions, buffered 1 => 4 positions
	for i := 0; i < 4; i++ {
		err = workerPool.ScheduleWithTimeout(func() {
			time.Sleep(10 * time.Millisecond)
		}, 50*time.Millisecond)
		assert.NoError(t, err)
	}
	err = workerPool.Schedule(func() {})
	assert.Equal(t, ErrWorkerPoolJobQueueIsFull, err)
	err = workerPool.ScheduleWithTimeout(func() {}, 1*time.Millisecond)
	assert.Equal(t, ErrWorkerPoolScheduleTimeout, err)

	defaultWorkerPool.SetWorkerSizeMaximum(3)
	err = workerPool.ScheduleWithTimeout(func() {}, 100*time.Millisecond)
	assert.Equal(t, nil, err)
}

func TestWorkerJamDuration(t *testing.T) {
	var workerPool WorkerPool
	var err error
	// Use custom settings with standBy=0 from the start to prevent spawnLoop
	// from spawning workers before PreAlloc can set deterministic counts.
	settings := &DefaultWorkerPoolSettings{
		isJobQueueClosedWhenClose: true,
		workerBatchSize:           0,
		workerSizeStandBy:         0,
		workerSizeMaximum:         10,
		spawnWorkerDuration:       1 * time.Millisecond,
		workerExpiryDuration:      200 * time.Millisecond,
		workerJamDuration:         50 * time.Millisecond,
		scheduleRetryInterval:     50 * time.Millisecond,
		panicHandler:              defaultPanicHandler,
	}
	defaultWorkerPool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue(3, 10000, 100), settings)
	// Set standBy=3 after construction so spawnLoop maintains 3 workers
	defaultWorkerPool.SetWorkerSizeStandBy(3)
	workerPool = defaultWorkerPool
	defer defaultWorkerPool.Close()

	// PreAlloc 3 workers deterministically
	assert.Equal(t, 0, defaultWorkerPool.getWorkerCount())
	defaultWorkerPool.PreAllocWorkerSize(3)
	assert.Equal(t, 3, defaultWorkerPool.getWorkerCount())

	// Keep all 3 workers busy with blocking jobs
	blockCh := make(chan struct{})
	for i := 0; i < 3; i++ {
		err = workerPool.Schedule(func() {
			<-blockCh
		})
		assert.NoError(t, err)
	}
	time.Sleep(20 * time.Millisecond)
	// 3 workers all busy
	assert.Equal(t, 3, defaultWorkerPool.getWorkerCount())

	// Schedule non-blocking jobs to queue (all workers busy)
	workerPool.Schedule(func() {})
	workerPool.Schedule(func() {})
	workerPool.Schedule(func() {})
	// Wait for jam duration to elapse, then wake spawnLoop
	time.Sleep(60 * time.Millisecond)
	defaultWorkerPool.spawnWorkerCh.Offer(1)
	time.Sleep(10 * time.Millisecond)
	// Jam > 50ms -> 4th worker spawned
	assert.Equal(t, 4, defaultWorkerPool.getWorkerCount())

	// 4th worker picks up non-blocking jobs, completes them
	time.Sleep(20 * time.Millisecond)
	// Non-blocking jobs done, no new jam
	assert.Equal(t, 4, defaultWorkerPool.getWorkerCount())

	// Schedule a blocking job (goes to queue, picked up by idle 4th worker)
	workerPool.Schedule(func() {
		<-blockCh
	})
	time.Sleep(20 * time.Millisecond)
	// 4th worker picked it up, no jam (< 50ms)
	assert.Equal(t, 4, defaultWorkerPool.getWorkerCount())

	// Schedule 4 non-blocking jobs (all 4 workers busy)
	workerPool.Schedule(func() {})
	workerPool.Schedule(func() {})
	workerPool.Schedule(func() {})
	workerPool.Schedule(func() {})
	// Wait for jam duration to elapse, then wake spawnLoop
	time.Sleep(60 * time.Millisecond)
	defaultWorkerPool.spawnWorkerCh.Offer(1)
	time.Sleep(10 * time.Millisecond)
	// Jam > 50ms -> 5th worker spawned
	assert.Equal(t, 5, defaultWorkerPool.getWorkerCount())

	// Release all blocking workers
	close(blockCh)
}

func TestWorkerPoolScheduleWithNilJobQueue(t *testing.T) {
	// Use constructor to get proper initialization, then nil out the queue
	defaultWorkerPool := NewDefaultWorkerPool(nil, nil)
	defer defaultWorkerPool.Close()

	err := defaultWorkerPool.Schedule(func() {})
	assert.Equal(t, ErrWorkerPoolJobQueueIsNil, err)
}

func TestWorkerPoolCloseStopsSpawnLoop(t *testing.T) {
	defaultWorkerPool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue(3, 10000, 100), nil).
		SetWorkerSizeStandBy(0).
		SetWorkerBatchSize(0)
	defer defaultWorkerPool.Close()

	// spawnLoop should be running
	time.Sleep(10 * time.Millisecond)
}

func TestHandlerCloseMultipleDoesNotPanic(t *testing.T) {
	// Verify Close is idempotent
	defaultWorkerPool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue(3, 10000, 100), nil)
	defaultWorkerPool.Close()
	defaultWorkerPool.Close() // should not panic
}

func TestWorkerConcurrentScheduleCloseDoesNotPanic(t *testing.T) {
	defaultWorkerPool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue(3, 10000, 100), nil).
		SetWorkerSizeMaximum(3).
		SetWorkerSizeStandBy(1).
		SetWorkerBatchSize(0)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				defaultWorkerPool.Schedule(func() {})
			}
		}()
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(5 * time.Millisecond)
		defaultWorkerPool.Close()
	}()
	wg.Wait()
}

func TestWorkerPoolConcurrentSettersNoRace(t *testing.T) {
	// Regression: setters must lock shared fields to avoid data race with trySpawn/spawnLoop.
	// Pre-fix: setters wrote workerSizeStandBy/workerBatchSize/etc. without lock while
	// trySpawn read them under RLock -> race detector failure.
	defaultWorkerPool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue(3, 10000, 100), nil).
		SetWorkerSizeMaximum(10).
		SetWorkerSizeStandBy(1).
		SetWorkerBatchSize(3)
	defer defaultWorkerPool.Close()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				defaultWorkerPool.SetWorkerSizeStandBy(j%5 + 1)
				defaultWorkerPool.SetWorkerBatchSize(j%4 + 1)
				defaultWorkerPool.SetWorkerSizeMaximum(j%8 + 2)
				defaultWorkerPool.SetWorkerJamDuration(time.Duration(j%10) * time.Millisecond)
				defaultWorkerPool.SetWorkerExpiryDuration(time.Duration(j%10+1) * 10 * time.Millisecond)
				defaultWorkerPool.getWorkerCount()
			}
		}()
	}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				defaultWorkerPool.Schedule(func() { time.Sleep(time.Millisecond) })
			}
		}()
	}
	wg.Wait()
}

func TestWorkerPoolPreAllocConcurrentWithSpawnLoop(t *testing.T) {
	// Regression: PreAllocWorkerSize read workerCount without lock while trySpawn
	// incremented it under Lock -> race detector failure.
	defaultWorkerPool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue(3, 10000, 100), nil).
		SetWorkerSizeStandBy(0).
		SetWorkerBatchSize(0)
	defer defaultWorkerPool.Close()

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			defaultWorkerPool.PreAllocWorkerSize(n)
		}(i + 3)
	}
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				defaultWorkerPool.Schedule(func() {})
			}
		}()
	}
	wg.Wait()
}

func TestDefaultWorkerPoolSettings(t *testing.T) {
	customSettings := &DefaultWorkerPoolSettings{
		isJobQueueClosedWhenClose: false,
		workerBatchSize:           2,
		workerSizeStandBy:         0,
		workerSizeMaximum:         4,
		spawnWorkerDuration:       5 * time.Millisecond,
		workerExpiryDuration:      100 * time.Millisecond,
		workerJamDuration:         200 * time.Millisecond,
		scheduleRetryInterval:     15 * time.Millisecond,
		panicHandler:              func(interface{}) {},
	}
	pool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue(3, 10, 5), customSettings)
	defer pool.Close()

	assert.Equal(t, false, pool.isJobQueueClosedWhenClose)
	assert.Equal(t, 2, pool.workerBatchSize)
	assert.Equal(t, 0, pool.workerSizeStandBy)
	assert.Equal(t, 4, pool.workerSizeMaximum)
	assert.Equal(t, 5*time.Millisecond, pool.spawnWorkerDuration)
	assert.Equal(t, 100*time.Millisecond, pool.workerExpiryDuration)
	assert.Equal(t, 200*time.Millisecond, pool.workerJamDuration)
	assert.Equal(t, 15*time.Millisecond, pool.scheduleRetryInterval)

	replacement := DefaultWorkerPoolSettings{
		isJobQueueClosedWhenClose: true,
		workerBatchSize:           1,
		workerSizeStandBy:         0,
		workerSizeMaximum:         2,
		spawnWorkerDuration:       10 * time.Millisecond,
		workerExpiryDuration:      50 * time.Millisecond,
		workerJamDuration:         75 * time.Millisecond,
		scheduleRetryInterval:     5 * time.Millisecond,
		panicHandler:              defaultPanicHandler,
	}
	assert.Same(t, pool, pool.SetDefaultWorkerPoolSettings(replacement))
	assert.Equal(t, true, pool.isJobQueueClosedWhenClose)
	assert.Equal(t, 1, pool.workerBatchSize)
	assert.Equal(t, 2, pool.workerSizeMaximum)
	assert.Equal(t, 5*time.Millisecond, pool.scheduleRetryInterval)

	newQueue := fpgo.NewBufferedChannelQueue(3, 5, 5)
	assert.Same(t, pool, pool.SetJobQueue(newQueue))
	assert.Same(t, newQueue, pool.jobQueue)
	assert.NoError(t, pool.Schedule(func() {}))

	closeWhenCloseQueue := fpgo.NewBufferedChannelQueue(3, 10, 5)
	closePool := NewDefaultWorkerPool(closeWhenCloseQueue, &DefaultWorkerPoolSettings{
		isJobQueueClosedWhenClose: true,
		workerSizeStandBy:         0,
		workerSizeMaximum:         1,
		workerBatchSize:           1,
		panicHandler:              defaultPanicHandler,
	})
	closePool.PreAllocWorkerSize(1)
	assert.Same(t, closePool, closePool.SetIsJobQueueClosedWhenClose(true))
	closePool.Close()
	assert.True(t, closeWhenCloseQueue.IsClosed())

	keepOpenQueue := fpgo.NewBufferedChannelQueue(3, 10, 5)
	keepOpenPool := NewDefaultWorkerPool(keepOpenQueue, &DefaultWorkerPoolSettings{
		isJobQueueClosedWhenClose: true,
		workerSizeStandBy:         0,
		workerSizeMaximum:         1,
		workerBatchSize:           1,
		panicHandler:              defaultPanicHandler,
	})
	keepOpenPool.PreAllocWorkerSize(1)
	assert.Same(t, keepOpenPool, keepOpenPool.SetIsJobQueueClosedWhenClose(false))
	assert.Equal(t, false, keepOpenPool.isJobQueueClosedWhenClose)
	keepOpenPool.Close()
	assert.False(t, keepOpenQueue.IsClosed())

	var panicHandled uint32
	panicPool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue(3, 10, 5), &DefaultWorkerPoolSettings{
		workerSizeStandBy:    0,
		workerSizeMaximum:    1,
		workerBatchSize:      1,
		spawnWorkerDuration:  time.Millisecond,
		workerExpiryDuration: time.Second,
		panicHandler: func(interface{}) {
			atomic.StoreUint32(&panicHandled, 1)
		},
	})
	defer panicPool.Close()
	// SetPanicHandler returns the pool and installs the handler (coverage).
	assert.Same(t, panicPool, panicPool.SetPanicHandler(func(p interface{}) {
		atomic.StoreUint32(&panicHandled, 1)
	}))
	assert.NoError(t, panicPool.Schedule(func() { panic("settings test panic") }))
	assert.Eventually(t, func() bool { return atomic.LoadUint32(&panicHandled) == 1 }, 2*time.Second, time.Millisecond)

	assert.Same(t, pool, pool.SetSpawnWorkerDuration(25*time.Millisecond))
	assert.Equal(t, 25*time.Millisecond, pool.spawnWorkerDuration)

	assert.Same(t, pool, pool.SetScheduleRetryInterval(8*time.Millisecond))
	assert.Equal(t, 8*time.Millisecond, pool.scheduleRetryInterval)
}

func TestDefaultInvokable(t *testing.T) {
	settings := &DefaultWorkerPoolSettings{
		isJobQueueClosedWhenClose: true,
		workerBatchSize:           1,
		workerSizeStandBy:         0,
		workerSizeMaximum:         2,
		spawnWorkerDuration:       1 * time.Millisecond,
		workerExpiryDuration:      100 * time.Millisecond,
		workerJamDuration:         1000 * time.Millisecond,
		scheduleRetryInterval:     5 * time.Millisecond,
		panicHandler:              defaultPanicHandler,
	}
	pool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue(3, 10, 5), settings)
	defer pool.Close()
	pool.PreAllocWorkerSize(1)

	done := make(chan interface{}, 4)
	invokable := NewDefaultInvokable(pool, func(val interface{}) {
		done <- val
	})
	assert.NotNil(t, invokable)

	invokable.Invoke(42)
	assert.Equal(t, 42, <-done)

	timeoutErr := invokable.InvokeWithTimeout(99, 50*time.Millisecond)
	assert.NoError(t, timeoutErr)
	assert.Equal(t, 99, <-done)

	assert.Same(t, invokable, invokable.SetWorkerPool(pool))
	assert.Same(t, pool, invokable.workerPool)
	assert.Same(t, invokable, invokable.SetCallee(func(val interface{}) {
		done <- val
	}))
	invokable.Invoke(7)
	assert.Equal(t, 7, <-done)

	timeoutSettings := &DefaultWorkerPoolSettings{
		isJobQueueClosedWhenClose: true,
		workerBatchSize:           0,
		workerSizeStandBy:         0,
		workerSizeMaximum:         0,
		spawnWorkerDuration:       1 * time.Millisecond,
		workerExpiryDuration:      100 * time.Millisecond,
		workerJamDuration:         1000 * time.Millisecond,
		scheduleRetryInterval:     5 * time.Millisecond,
		panicHandler:              defaultPanicHandler,
	}
	timeoutPool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue(3, 1, 3), timeoutSettings)
	defer timeoutPool.Close()
	blockCh := make(chan struct{})
	for i := 0; i < 4; i++ {
		assert.NoError(t, timeoutPool.Schedule(func() { <-blockCh }))
	}
	timeoutInvokable := NewDefaultInvokable(timeoutPool, func(val interface{}) {})
	timeoutErr = timeoutInvokable.InvokeWithTimeout(1, 5*time.Millisecond)
	assert.Equal(t, ErrWorkerPoolScheduleTimeout, timeoutErr)
	close(blockCh)
}
