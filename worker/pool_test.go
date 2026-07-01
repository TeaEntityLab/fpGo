package worker

import (
	"runtime"
	"sync/atomic"
	"testing"
	"time"
	// "sync"

	"github.com/stretchr/testify/assert"

	fpgo "github.com/TeaEntityLab/fpGo/v2"
)

func TestWorkerPool(t *testing.T) {
	var workerPool WorkerPool
	var err error
	defaultWorkerPool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue[func()](3, 10000, 100), nil).
		SetSpawnWorkerDuration(1 * time.Millisecond / 10).
		SetWorkerExpiryDuration(2 * time.Millisecond).
		SetWorkerSizeMaximum(5).
		SetWorkerSizeStandBy(1).
		SetWorkerBatchSize(3)
	// defaultWorkerPool.PreAllocWorkerSize(5)
	workerPool = defaultWorkerPool

	// Test Spawn
	assert.Equal(t, 0, defaultWorkerPool.getWorkerCount())
	for i := 0; i < 8; i++ {
		v := i
		err = workerPool.Schedule(func() {
			// Nothing to do
			time.Sleep(3 * time.Millisecond / 2)
			t.Log(v)
		})
		assert.NoError(t, err)
	}
	assert.Eventually(t, func() bool {
		return defaultWorkerPool.getWorkerCount() >= 2
	}, 50*time.Millisecond, time.Millisecond)
	defaultWorkerPool.PreAllocWorkerSize(5)
	assert.Equal(t, 5, defaultWorkerPool.getWorkerCount())

	// Test ScaleDown
	assert.Eventually(t, func() bool {
		return defaultWorkerPool.getWorkerCount() == 1
	}, 100*time.Millisecond, time.Millisecond)
	for i := 0; i < 4; i++ {
		v := i
		err = workerPool.Schedule(func() {
			// Overtime
			time.Sleep(3 * time.Millisecond)
			t.Log(v)
		})
		assert.NoError(t, err)
	}
	assert.Eventually(t, func() bool {
		return defaultWorkerPool.getWorkerCount() >= 1
	}, 50*time.Millisecond, time.Millisecond)
	defaultWorkerPool.SetWorkerSizeStandBy(0)
	assert.Eventually(t, func() bool {
		return defaultWorkerPool.getWorkerCount() == 0
	}, 100*time.Millisecond, time.Millisecond)
}

func TestScheduleWithTimeout(t *testing.T) {
	var workerPool WorkerPool
	var err error
	defaultWorkerPool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue[func()](3, 1, 3), nil).
		SetSpawnWorkerDuration(1 * time.Millisecond / 10).
		SetWorkerExpiryDuration(2 * time.Millisecond).
		SetWorkerSizeMaximum(0).
		SetWorkerSizeStandBy(0).
		SetWorkerBatchSize(0)
	defer defaultWorkerPool.Close()
	// defaultWorkerPool.PreAllocWorkerSize(5)
	workerPool = defaultWorkerPool

	// Test ScheduleWithTimeout
	// channel: 3 positions, buffered 1 => 4 positions
	for i := 0; i < 4; i++ {
		v := i
		err = workerPool.ScheduleWithTimeout(func() {
			// Nothing to do
			time.Sleep(3 * time.Millisecond)
			t.Log(v)
		}, 1*time.Millisecond)
		assert.NoError(t, err)
	}
	err = workerPool.Schedule(func() {})
	assert.Equal(t, ErrWorkerPoolJobQueueIsFull, err)
	err = workerPool.ScheduleWithTimeout(func() {}, 1*time.Millisecond/2)
	assert.Equal(t, ErrWorkerPoolScheduleTimeout, err)

	defaultWorkerPool.SetWorkerSizeMaximum(3)
	// Generous timeout: this asserts that raising the worker ceiling lets the
	// job schedule (a worker spawns and drains a slot), not that the spawn wins
	// a tight millisecond race. A small bound here is scheduler-load sensitive
	// and flakes under `go test ./...`; a real "never schedules" bug still fails.
	err = workerPool.ScheduleWithTimeout(func() {}, 500*time.Millisecond)
	assert.Equal(t, nil, err)
}

func TestWorkerJamDuration(t *testing.T) {
	defaultWorkerPool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue[func()](3, 10000, 100), nil).
		SetSpawnWorkerDuration(1 * time.Millisecond / 10).
		SetWorkerExpiryDuration(100 * time.Millisecond).
		SetWorkerJamDuration(3 * time.Millisecond).
		SetWorkerSizeMaximum(10).
		SetWorkerSizeStandBy(3).
		SetWorkerBatchSize(0)
	workerPool := WorkerPool(defaultWorkerPool)

	started := make(chan struct{}, 3)
	release := make(chan struct{})
	for i := 0; i < 3; i++ {
		err := workerPool.Schedule(func() {
			started <- struct{}{}
			<-release
		})
		assert.NoError(t, err)
	}

	for i := 0; i < 3; i++ {
		select {
		case <-started:
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("worker %d did not start", i)
		}
	}
	assert.Eventually(t, func() bool {
		return defaultWorkerPool.getWorkerCount() >= 3
	}, 50*time.Millisecond, time.Millisecond)
	time.Sleep(5 * time.Millisecond)

	for i := 0; i < 3; i++ {
		assert.NoError(t, workerPool.Schedule(func() {}))
	}

	assert.Eventually(t, func() bool {
		return defaultWorkerPool.getWorkerCount() >= 4
	}, 100*time.Millisecond, time.Millisecond)

	close(release)
	defaultWorkerPool.Close()
}

func TestWorkerPoolIsClosed(t *testing.T) {
	pool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue[func()](3, 10, 5), nil)
	assert.Equal(t, false, pool.IsClosed())

	pool.Close()
	assert.Equal(t, true, pool.IsClosed())
}

func TestNewDefaultWorkerPool(t *testing.T) {
	pool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue[func()](3, 10, 5), nil)
	assert.NotNil(t, pool)

	pool.Close()
}

func TestDefaultWorkerPoolSettings(t *testing.T) {
	pool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue[func()](3, 10, 5), nil)
	assert.NotNil(t, pool)

	pool.SetDefaultWorkerPoolSettings(DefaultWorkerPoolSettings{
		workerSizeMaximum: 20,
	})
	pool.SetJobQueue(fpgo.NewBufferedChannelQueue[func()](3, 10, 5))
	pool.SetIsJobQueueClosedWhenClose(false)
	pool.SetWorkerSizeMaximum(15)
	pool.SetWorkerSizeStandBy(3)
	pool.SetWorkerBatchSize(5)
	pool.SetWorkerExpiryDuration(10 * time.Second)
	pool.SetWorkerJamDuration(20 * time.Second)
	pool.SetSpawnWorkerDuration(200 * time.Millisecond)
	pool.SetScheduleRetryInterval(100 * time.Millisecond)
	pool.SetPanicHandler(func(interface{}) {})

	pool.Close()
}

func TestDefaultInvokable(t *testing.T) {
	pool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue[func()](3, 10, 5), nil)

	done := make(chan int, 1)
	invokable := NewDefaultInvokable(pool, func(val int) {
		done <- val
	})
	assert.NotNil(t, invokable)

	invokable.Invoke(42)
	assert.Equal(t, 42, <-done)

	var timeoutErr error
	timeoutErr = invokable.InvokeWithTimeout(42, 10*time.Millisecond)
	assert.Nil(t, timeoutErr)
	assert.Equal(t, 42, <-done)

	invokable.SetWorkerPool(pool)
	assert.NotNil(t, invokable.workerPool)
	invokable.SetCallee(func(val int) {
		done <- val + 1
	})
	invokable.Invoke(1)
	assert.Equal(t, 2, <-done)

	pool.Close()
}

func TestWorkerPoolCloseAlreadyClosed(t *testing.T) {
	pool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue[func()](3, 10, 5), nil)

	pool.Close()
	assert.True(t, pool.IsClosed())

	assert.NotPanics(t, func() {
		pool.Close()
	})
	assert.True(t, pool.IsClosed())
}

func TestWorkerPoolScheduleWhenClosed(t *testing.T) {
	pool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue[func()](3, 10, 5), nil)
	pool.Close()

	err := pool.Schedule(func() {})
	assert.Equal(t, ErrWorkerPoolIsClosed, err)
}

func TestWorkerPoolScheduleWithNilJobQueue(t *testing.T) {
	pool := NewDefaultWorkerPool(nil, nil)
	defer pool.Close()

	err := pool.Schedule(func() {})
	assert.Equal(t, ErrWorkerPoolJobQueueIsNil, err)
}

func TestWorkerPoolScheduleWithTimeoutQueueFull(t *testing.T) {
	pool := NewDefaultWorkerPool(
		fpgo.NewBufferedChannelQueue[func()](1, 0, 1),
		nil,
	)
	pool.SetWorkerSizeMaximum(0)
	pool.SetWorkerSizeStandBy(0)
	pool.SetWorkerBatchSize(0)

	err := pool.Schedule(func() { time.Sleep(time.Hour) })
	assert.Nil(t, err)

	err = pool.ScheduleWithTimeout(func() {}, 10*time.Millisecond)
	assert.Equal(t, ErrWorkerPoolScheduleTimeout, err)
}

func TestWorkerPoolScheduleWithTimeoutClosedDuringRetry(t *testing.T) {
	pool := NewDefaultWorkerPool(
		fpgo.NewBufferedChannelQueue[func()](1, 0, 1),
		nil,
	)
	pool.SetWorkerSizeMaximum(0)
	pool.SetWorkerSizeStandBy(0)
	pool.SetWorkerBatchSize(0)
	pool.SetScheduleRetryInterval(20 * time.Millisecond)

	assert.NoError(t, pool.Schedule(func() { time.Sleep(50 * time.Millisecond) }))

	done := make(chan error, 1)
	go func() {
		done <- pool.ScheduleWithTimeout(func() {}, 60*time.Millisecond)
	}()

	time.Sleep(5 * time.Millisecond)
	pool.Close()

	assert.Equal(t, ErrWorkerPoolIsClosed, <-done)
}

func TestWorkerPoolScheduleWithTimeoutRetryIntervalClampAndDirectReturn(t *testing.T) {
	pool := NewDefaultWorkerPool(
		fpgo.NewBufferedChannelQueue[func()](1, 0, 1),
		nil,
	)
	pool.SetWorkerSizeMaximum(0)
	pool.SetWorkerSizeStandBy(0)
	pool.SetWorkerBatchSize(0)
	pool.SetScheduleRetryInterval(100 * time.Millisecond)

	assert.NoError(t, pool.Schedule(func() { time.Sleep(20 * time.Millisecond) }))
	assert.Equal(t, ErrWorkerPoolScheduleTimeout, pool.ScheduleWithTimeout(func() {}, 3*time.Millisecond))

	pool.Close()
	assert.Equal(t, ErrWorkerPoolIsClosed, pool.ScheduleWithTimeout(func() {}, 10*time.Millisecond))
}

func TestWorkerPoolScheduleWithTimeoutSuccessAfterRetry(t *testing.T) {
	pool := NewDefaultWorkerPool(
		fpgo.NewBufferedChannelQueue[func()](1, 0, 1),
		nil,
	)
	pool.SetWorkerSizeMaximum(1)
	pool.SetWorkerSizeStandBy(0)
	pool.SetWorkerBatchSize(1)
	pool.SetScheduleRetryInterval(2 * time.Millisecond)
	pool.SetSpawnWorkerDuration(1 * time.Millisecond)
	pool.SetWorkerExpiryDuration(100 * time.Millisecond)

	assert.NoError(t, pool.Schedule(func() { time.Sleep(5 * time.Millisecond) }))
	assert.NoError(t, pool.ScheduleWithTimeout(func() {}, 50*time.Millisecond))
	pool.Close()
}

func TestWorkerPoolGenerateWorkerWithMaximumEarlyReturn(t *testing.T) {
	pool := NewDefaultWorkerPool(
		fpgo.NewBufferedChannelQueue[func()](3, 10, 5),
		nil,
	)
	pool.SetWorkerSizeMaximum(0)
	pool.SetWorkerSizeStandBy(0)
	pool.SetWorkerBatchSize(0)

	pool.PreAllocWorkerSize(5)
	assert.Equal(t, 0, pool.getWorkerCount())
	pool.Close()
}

func TestWorkerPoolTrySpawnWithMaximumCap(t *testing.T) {
	pool := NewDefaultWorkerPool(
		fpgo.NewBufferedChannelQueue[func()](3, 100, 100),
		nil,
	)
	pool.SetWorkerBatchSize(1)
	pool.SetWorkerSizeStandBy(10)
	pool.SetWorkerSizeMaximum(3)
	pool.SetSpawnWorkerDuration(1 * time.Millisecond)
	pool.SetWorkerJamDuration(100 * time.Hour)
	pool.SetWorkerExpiryDuration(100 * time.Hour)

	for i := 0; i < 5; i++ {
		pool.Schedule(func() {})
	}
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, 3, pool.getWorkerCount())
	pool.Close()
}

func TestWorkerPoolScheduleWithTimeoutIsClosed(t *testing.T) {
	pool := NewDefaultWorkerPool(
		fpgo.NewBufferedChannelQueue[func()](1, 0, 1),
		nil,
	)
	pool.SetWorkerSizeMaximum(0)
	pool.SetWorkerSizeStandBy(0)
	pool.SetWorkerBatchSize(0)
	pool.SetScheduleRetryInterval(1 * time.Hour)

	err := pool.Schedule(func() { time.Sleep(time.Hour) })
	assert.Nil(t, err)

	go func() {
		time.Sleep(1 * time.Millisecond)
		pool.Close()
	}()

	err = pool.ScheduleWithTimeout(func() {}, 100*time.Millisecond)
	assert.Equal(t, ErrWorkerPoolIsClosed, err)
}

func TestWorkerPoolSpawnLoopClosedGracefully(t *testing.T) {
	pool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue[func()](3, 10, 5), nil)
	assert.Equal(t, false, pool.IsClosed())

	pool.Close()
	time.Sleep(5 * time.Millisecond)

	pool.notifyWorkers()
	assert.True(t, pool.IsClosed())
}

func TestWorkerPoolScheduleWithTimeoutCloseDuringRetry(t *testing.T) {
	pool := NewDefaultWorkerPool(
		fpgo.NewBufferedChannelQueue[func()](1, 0, 1),
		nil,
	)
	pool.SetWorkerSizeMaximum(0)
	pool.SetWorkerSizeStandBy(0)
	pool.SetWorkerBatchSize(0)
	pool.SetScheduleRetryInterval(1 * time.Hour)

	err := pool.Schedule(func() {})
	assert.Nil(t, err)

	go func() {
		time.Sleep(2 * time.Millisecond)
		pool.Close()
	}()

	err = pool.ScheduleWithTimeout(func() {}, 100*time.Millisecond)
	assert.Equal(t, ErrWorkerPoolIsClosed, err)
}

func TestWorkerPoolWorkerPanicRecovered(t *testing.T) {
	pool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue[func()](3, 10, 5), nil)
	pool.SetWorkerExpiryDuration(time.Second)

	err := pool.Schedule(func() { panic("test panic - should be recovered") })
	assert.Nil(t, err)
	time.Sleep(50 * time.Millisecond) // Let worker process and panic/recover

	// Pool should still be usable after worker panic
	assert.False(t, pool.IsClosed())

	// Schedule a normal job to verify pool works after panic
	var called uint32
	err = pool.Schedule(func() { atomic.StoreUint32(&called, 1) })
	assert.Nil(t, err)
	assert.Eventually(t, func() bool {
		return atomic.LoadUint32(&called) == 1
	}, time.Second, time.Millisecond)

	pool.Close()
}

func TestWorkerPoolScheduleWithTimeoutImmediateSuccess(t *testing.T) {
	pool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue[func()](3, 10, 5), nil)
	defer pool.Close()

	ran := make(chan struct{}, 1)
	err := pool.ScheduleWithTimeout(func() {
		ran <- struct{}{}
	}, 20*time.Millisecond)
	assert.NoError(t, err)

	select {
	case <-ran:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected scheduled job to run")
	}
}

func TestWorkerPoolScheduleWithTimeoutTimeout(t *testing.T) {
	pool := NewDefaultWorkerPool(
		fpgo.NewBufferedChannelQueue[func()](1, 0, 1),
		nil,
	)
	defer pool.Close()
	pool.SetWorkerSizeMaximum(0)
	pool.SetWorkerSizeStandBy(0)
	pool.SetWorkerBatchSize(0)
	pool.SetScheduleRetryInterval(time.Millisecond)

	assert.NoError(t, pool.Schedule(func() {}))
	err := pool.ScheduleWithTimeout(func() {}, 5*time.Millisecond)
	assert.Equal(t, ErrWorkerPoolScheduleTimeout, err)
}

func TestWorkerPoolScheduleWithTimeoutImmediateError(t *testing.T) {
	pool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue[func()](1, 1, 1), nil)
	pool.Close()

	err := pool.ScheduleWithTimeout(func() {}, 20*time.Millisecond)
	assert.Equal(t, ErrWorkerPoolIsClosed, err)
}

func TestWorkerPoolScheduleWithTimeoutZeroTimeout(t *testing.T) {
	pool := NewDefaultWorkerPool(
		fpgo.NewBufferedChannelQueue[func()](1, 0, 1),
		nil,
	)
	defer pool.Close()
	pool.SetWorkerSizeMaximum(0)
	pool.SetWorkerSizeStandBy(0)
	pool.SetWorkerBatchSize(0)

	assert.NoError(t, pool.Schedule(func() {}))
	err := pool.ScheduleWithTimeout(func() {}, 0)
	assert.Equal(t, ErrWorkerPoolScheduleTimeout, err)
}

func TestWorkerPoolScheduleWithTimeoutRetryIntervalClamp(t *testing.T) {
	pool := NewDefaultWorkerPool(
		fpgo.NewBufferedChannelQueue[func()](1, 0, 1),
		nil,
	)
	defer pool.Close()
	pool.SetWorkerSizeMaximum(0)
	pool.SetWorkerSizeStandBy(0)
	pool.SetWorkerBatchSize(0)
	pool.SetScheduleRetryInterval(time.Hour)

	assert.NoError(t, pool.Schedule(func() {}))
	start := time.Now()
	err := pool.ScheduleWithTimeout(func() {}, 9*time.Millisecond)
	elapsed := time.Since(start)
	assert.Equal(t, ErrWorkerPoolScheduleTimeout, err)
	assert.Less(t, elapsed, 200*time.Millisecond)
}

func TestWorkerPoolScheduleWithTimeoutNoRetryNeededOnClosedPool(t *testing.T) {
	pool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue[func()](1, 1, 1), nil)
	pool.Close()
	assert.True(t, pool.IsClosed())

	err := pool.ScheduleWithTimeout(func() {}, time.Millisecond)
	assert.Equal(t, ErrWorkerPoolIsClosed, err)
}

func TestWorkerPoolExpiryStandby(t *testing.T) {
	pool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue[func()](3, 10, 5), nil)
	pool.SetWorkerBatchSize(1)
	pool.SetWorkerSizeStandBy(2)
	pool.SetWorkerExpiryDuration(1 * time.Millisecond)
	pool.SetSpawnWorkerDuration(1 * time.Millisecond)

	for i := 0; i < 5; i++ {
		pool.Schedule(func() { time.Sleep(1 * time.Millisecond) })
	}
	time.Sleep(20 * time.Millisecond) // Let workers spawn, do jobs, then expire down to standby

	pool.Close()
}

func TestWorkerPoolSpawnLoopPanicRecovery(t *testing.T) {
	pool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue[func()](3, 10, 5), nil)
	pool.jobQueue = nil
	pool.isJobQueueClosedWhenClose = false

	pool.spawnWorkerCh.Offer(1)
	time.Sleep(10 * time.Millisecond)

	pool.Close()
	assert.True(t, pool.IsClosed())
}

func TestWorkerPoolCloseStopsSpawnLoop(t *testing.T) {
	baseline := runtime.NumGoroutine()

	jobQueue := fpgo.NewBufferedChannelQueue[func()](10, 10000, 100)
	pool := NewDefaultWorkerPool(jobQueue, &DefaultWorkerPoolSettings{
		isJobQueueClosedWhenClose: true,
		workerBatchSize:           5,
		workerSizeStandBy:         0,
		workerSizeMaximum:         0,
		spawnWorkerDuration:       100 * time.Millisecond,
		workerExpiryDuration:      5000 * time.Millisecond,
		workerJamDuration:         1000 * time.Millisecond,
		scheduleRetryInterval:     50 * time.Millisecond,
		panicHandler:              defaultPanicHandler,
	})

	pool.Close()

	// Before the fix, spawnLoop blocked forever on <-spawnWorkerCh because
	// Close() never signaled it to exit. After the fix, Close() closes the
	// done channel, which spawnLoop selects on.
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if runtime.NumGoroutine() <= baseline {
			break
		}
		time.Sleep(5 * time.Millisecond)
	}

	assert.LessOrEqual(t, runtime.NumGoroutine(), baseline,
		"spawnLoop goroutine should exit after Close(), not leak")
}
