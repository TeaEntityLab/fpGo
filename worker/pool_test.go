package worker

import (
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
	assert.Equal(t, 0, defaultWorkerPool.workerCount)
	for i := 0; i < 8; i++ {
		v := i
		err = workerPool.Schedule(func() {
			// Nothing to do
			time.Sleep(3 * time.Millisecond / 2)
			t.Log(v)
		})
		assert.NoError(t, err)
	}
	time.Sleep(5 * time.Millisecond / 4)
	// BatchSize: 3, Jobs: 8 -> ceil(8/3) = 3 workers
	assert.Equal(t, 3, defaultWorkerPool.workerCount)
	defaultWorkerPool.PreAllocWorkerSize(5)
	assert.Equal(t, 5, defaultWorkerPool.workerCount)

	// Test ScaleDown
	time.Sleep(10 * time.Millisecond)
	// workerSizeStandBy: 1
	assert.Equal(t, 1, defaultWorkerPool.workerCount)
	for i := 0; i < 4; i++ {
		v := i
		err = workerPool.Schedule(func() {
			// Overtime
			time.Sleep(3 * time.Millisecond)
			t.Log(v)
		})
		assert.NoError(t, err)
	}
	time.Sleep(1 * time.Millisecond)
	// BatchSize: 3, Jobs: 4 -> ceil(4/3) = 2 workers
	assert.GreaterOrEqual(t, defaultWorkerPool.workerCount, 2)
	defaultWorkerPool.SetWorkerSizeStandBy(0)
	time.Sleep(10 * time.Millisecond)
	// workerSizeStandBy: 1
	assert.Equal(t, 0, defaultWorkerPool.workerCount)
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
	err = workerPool.ScheduleWithTimeout(func() {}, 10*time.Millisecond)
	assert.Equal(t, nil, err)
}

func TestWorkerJamDuration(t *testing.T) {
	var workerPool WorkerPool
	var err error
	defaultWorkerPool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue[func()](3, 10000, 100), nil).
		SetSpawnWorkerDuration(1 * time.Millisecond / 10).
		SetWorkerExpiryDuration(5 * time.Millisecond).
		SetWorkerJamDuration(3 * time.Millisecond).
		SetWorkerSizeMaximum(10).
		SetWorkerSizeStandBy(3).
		SetWorkerBatchSize(0)
	// defaultWorkerPool.PreAllocWorkerSize(5)
	workerPool = defaultWorkerPool

	// Test Spawn
	assert.Equal(t, 0, defaultWorkerPool.workerCount)
	anyOneDone := false
	for i := 0; i < 3; i++ {
		v := i
		err = workerPool.Schedule(func() {
			// Nothing to do
			time.Sleep(20 * time.Millisecond)
			t.Log(v)
			anyOneDone = true
		})
		assert.NoError(t, err)
	}
	time.Sleep(3 * time.Millisecond)
	// BatchSize: 0, SetWorkerSizeStandBy: 3 -> 3 workers
	assert.Equal(t, 3, defaultWorkerPool.workerCount)
	time.Sleep(3 * time.Millisecond)
	// Though there're blocking jobs, but no newest job goes into the queue
	assert.Equal(t, 3, defaultWorkerPool.workerCount)
	// There're new jobs going to the queue, and all goroutines are busy
	workerPool.Schedule(func() {})
	workerPool.Schedule(func() {})
	workerPool.Schedule(func() {})
	time.Sleep(3 * time.Millisecond)
	// A new expected goroutine is generated
	assert.Equal(t, 4, defaultWorkerPool.workerCount)
	workerPool.Schedule(func() {})
	workerPool.Schedule(func() {})
	workerPool.Schedule(func() {})
	time.Sleep(3 * time.Millisecond)
	// Only non blocking jobs, thus keep the same amount
	assert.Equal(t, 4, defaultWorkerPool.workerCount)
	// There's a blocking jobs going to the queue
	workerPool.Schedule(func() {
		time.Sleep(20 * time.Millisecond)
		t.Log(3)
		anyOneDone = true
	})
	time.Sleep(3 * time.Millisecond)
	// Though there're blocking jobs, but no newest job goes into the queue
	assert.Equal(t, 4, defaultWorkerPool.workerCount)
	// There're new jobs going to the queue, and all goroutines are busy
	workerPool.Schedule(func() {})
	workerPool.Schedule(func() {})
	workerPool.Schedule(func() {})
	workerPool.Schedule(func() {})
	assert.Equal(t, false, anyOneDone)
	time.Sleep(1 * time.Millisecond)
	// A new expected goroutine is generated
	assert.Equal(t, 5, defaultWorkerPool.workerCount)
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
	assert.Equal(t, 0, pool.workerCount)
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
	assert.Equal(t, 3, pool.workerCount)
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
	var called bool
	err = pool.Schedule(func() { called = true })
	assert.Nil(t, err)
	time.Sleep(50 * time.Millisecond)
	assert.True(t, called)

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
