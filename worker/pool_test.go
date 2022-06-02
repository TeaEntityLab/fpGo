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
