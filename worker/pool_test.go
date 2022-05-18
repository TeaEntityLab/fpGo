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
	defaultWorkerPool := NewDefaultWorkerPool(fpgo.NewBufferedChannelQueue[func()](3, 10000, 100)).
		SetSpawnWorkerDuration(1 * time.Millisecond / 10).
		SetFreeWorkerDuration(1 * time.Millisecond / 10).
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
	// assert.Equal(t, 2, defaultWorkerPool.workerCount)
	defaultWorkerPool.SetWorkerSizeStandBy(0)
	time.Sleep(10 * time.Millisecond)
	// workerSizeStandBy: 1
	assert.Equal(t, 0, defaultWorkerPool.workerCount)
}
