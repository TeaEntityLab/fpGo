package worker

import (
	"errors"
	"log"
	"runtime"
	"sync"
	"time"

	fpgo "github.com/TeaEntityLab/fpGo"
)

var (
	// ErrWorkerPoolJobQueueIsFull WorkerPool JobQueue Is Full
	ErrWorkerPoolJobQueueIsFull = errors.New("workerPool JobQueue is full")
	// ErrWorkerPoolIsClosed WorkerPool Is Closed
	ErrWorkerPoolIsClosed = errors.New("workerPool is closed")
	// ErrWorkerPoolScheduleTimeout WorkerPool Schedule Timeout
	ErrWorkerPoolScheduleTimeout = errors.New("workerPool schedule timeout")
	// ErrWorkerPoolJobQueueIsNil WorkerPool JobQueue Is Nil
	ErrWorkerPoolJobQueueIsNil = errors.New("workerPool JobQueue is nil")
)

// WorkerPool

// WorkerPool WorkerPool inspired by Java ExecutorService
type WorkerPool interface {
	Close()
	IsClosed() bool

	Schedule(func()) error
	ScheduleWithTimeout(func(), time.Duration) error
}

// DefaultWorkerPoolSettings Settings for DefaultWorkerPool
type DefaultWorkerPoolSettings struct {
	// JobQueue

	isJobQueueClosedWhenClose bool
	workerBatchSize           int

	// Worker

	workerSizeStandBy     int
	workerSizeMaximum     int
	spawnWorkerDuration   time.Duration
	workerExpiryDuration  time.Duration
	workerJamDuration     time.Duration
	scheduleRetryInterval time.Duration

	// Panic Handler

	panicHandler func(interface{})
}

var defaultPanicHandler = func(panic interface{}) {
	log.Printf("panic from worker: %v\n", panic)
	buf := make([]byte, 4096)
	log.Printf("panic from worker: %s\n", string(buf[:runtime.Stack(buf, false)]))
}

var defaultDefaultWorkerSettings = &DefaultWorkerPoolSettings{
	isJobQueueClosedWhenClose: true,
	workerBatchSize:           5,
	workerSizeStandBy:         5,
	workerSizeMaximum:         1000,
	spawnWorkerDuration:       100 * time.Millisecond,
	workerExpiryDuration:      5000 * time.Millisecond,
	workerJamDuration:         1000 * time.Millisecond,
	scheduleRetryInterval:     50 * time.Millisecond,
	panicHandler:              defaultPanicHandler,
}

// DefaultWorkerPool DefaultWorkerPool inspired by Java ExecutorService
type DefaultWorkerPool struct {
	isClosed fpgo.AtomBool
	lock     sync.RWMutex

	jobQueue *fpgo.BufferedChannelQueue

	workerCount   int
	workerBusy    int
	spawnWorkerCh fpgo.ChannelQueue
	done          chan struct{}
	lastAliveTime time.Time

	// Settings
	DefaultWorkerPoolSettings
}

// NewDefaultWorkerPool New a DefaultWorkerPool
func NewDefaultWorkerPool(jobQueue *fpgo.BufferedChannelQueue, settings *DefaultWorkerPoolSettings) *DefaultWorkerPool {
	if settings == nil {
		settings = defaultDefaultWorkerSettings
	}
	workerPool := &DefaultWorkerPool{
		jobQueue: jobQueue,

		spawnWorkerCh: fpgo.NewChannelQueue(1),
		done:          make(chan struct{}),
		// Settings
		DefaultWorkerPoolSettings: *settings,
	}
	go workerPool.spawnLoop()

	return workerPool
}

// trySpawn Try Spawn Goroutine as possible
func (workerPoolSelf *DefaultWorkerPool) trySpawn() {
	workerPoolSelf.lock.RLock()
	batchSize := workerPoolSelf.workerBatchSize
	workerCount := workerPoolSelf.workerCount
	workerBusy := workerPoolSelf.workerBusy
	workerSizeStandBy := workerPoolSelf.workerSizeStandBy
	workerSizeMaximum := workerPoolSelf.workerSizeMaximum
	lastAliveTime := workerPoolSelf.lastAliveTime
	workerJamDuration := workerPoolSelf.workerJamDuration
	jobQueue := workerPoolSelf.jobQueue
	workerPoolSelf.lock.RUnlock()

	queueCount := 0
	if jobQueue != nil {
		queueCount = jobQueue.Count()
	}

	var expectedWorkerCount int
	if batchSize > 0 {
		expectedWorkerCount = queueCount / batchSize
		if queueCount%batchSize > 0 {
			expectedWorkerCount++
		}
	}
	if workerSizeStandBy > expectedWorkerCount {
		expectedWorkerCount = workerSizeStandBy
	}
	if workerSizeMaximum > 0 && expectedWorkerCount > workerSizeMaximum {
		expectedWorkerCount = workerSizeMaximum
	}
	// Avoid Jam if (now - lastAliveTime) is over workerJamDuration
	if time.Since(lastAliveTime) > workerJamDuration &&
		workerBusy >= workerCount &&
		workerCount >= expectedWorkerCount {
		expectedWorkerCount = workerCount + 1
	}

	if workerCount < expectedWorkerCount {
		for i := workerCount; i < expectedWorkerCount; i++ {
			workerPoolSelf.generateWorkerWithMaximum(expectedWorkerCount)
		}
	}
}

// PreAllocWorkerSize PreAllocate Workers
func (workerPoolSelf *DefaultWorkerPool) PreAllocWorkerSize(preAllocWorkerSize int) {
	workerPoolSelf.lock.RLock()
	current := workerPoolSelf.workerCount
	workerPoolSelf.lock.RUnlock()
	for i := current; i < preAllocWorkerSize; i++ {
		workerPoolSelf.generateWorkerWithMaximum(preAllocWorkerSize)
	}
}

func (workerPoolSelf *DefaultWorkerPool) spawnLoop() {
	defer func() {
		if panic := recover(); panic != nil {
			defaultPanicHandler(panic)
		}
	}()

	for {
		select {
		case <-workerPoolSelf.done:
			return
		case _, ok := <-workerPoolSelf.spawnWorkerCh:
			if !ok {
				return
			}
			if workerPoolSelf.IsClosed() {
				return
			}

			workerPoolSelf.trySpawn()

			workerPoolSelf.lock.RLock()
			spawnWorkerDuration := workerPoolSelf.spawnWorkerDuration
			workerPoolSelf.lock.RUnlock()
			time.Sleep(spawnWorkerDuration)
		}
	}
}

func (workerPoolSelf *DefaultWorkerPool) notifyWorkers() {
	workerPoolSelf.lock.RLock()
	workerCount := workerPoolSelf.workerCount
	workerSizeStandBy := workerPoolSelf.workerSizeStandBy
	jobQueue := workerPoolSelf.jobQueue
	workerPoolSelf.lock.RUnlock()

	queueCount := 0
	if jobQueue != nil {
		queueCount = jobQueue.Count()
	}
	if workerCount < workerSizeStandBy || queueCount > 0 {
		workerPoolSelf.spawnWorkerCh.Offer(1)
	}
}

func (workerPoolSelf *DefaultWorkerPool) generateWorkerWithMaximum(maximum int) {
	// Initial
	workerPoolSelf.lock.Lock()
	defer workerPoolSelf.lock.Unlock()
	if workerPoolSelf.workerCount >= maximum ||
		workerPoolSelf.workerCount >= workerPoolSelf.workerSizeMaximum {
		return
	}
	// workerID := time.Now()
	workerPoolSelf.lastAliveTime = time.Now()
	workerPoolSelf.workerCount++
	isBusy := false

	go func() {
		// Recover & Recycle
		defer func() {
			if panic := recover(); panic != nil {
				workerPoolSelf.lock.RLock()
				handler := workerPoolSelf.panicHandler
				workerPoolSelf.lock.RUnlock()
				if handler != nil {
					handler(panic)
				}
			}

			workerPoolSelf.lock.Lock()
			workerPoolSelf.workerCount--
			if isBusy {
				workerPoolSelf.workerBusy--
			}
			workerPoolSelf.lock.Unlock()
		}()

		// Do Jobs
	loopLabel:
		for {
			workerPoolSelf.lock.Lock()
			workerPoolSelf.lastAliveTime = time.Now()
			workerPoolSelf.lock.Unlock()

			if workerPoolSelf.IsClosed() {
				return
			}

			workerPoolSelf.lock.RLock()
			expiryDuration := workerPoolSelf.workerExpiryDuration
			jobQueue := workerPoolSelf.jobQueue
			workerPoolSelf.lock.RUnlock()

			var jobCh chan interface{}
			if jobQueue != nil {
				jobCh = jobQueue.GetChannel()
			}

			select {
			case job := <-jobCh:
				if job != nil {
					workerPoolSelf.lock.Lock()
					isBusy = true
					workerPoolSelf.workerBusy++
					workerPoolSelf.lock.Unlock()

					(job.(func()))()

					workerPoolSelf.lock.Lock()
					workerPoolSelf.workerBusy--
					isBusy = false
					workerPoolSelf.lock.Unlock()
				}
			case <-time.After(expiryDuration):
				workerPoolSelf.lock.RLock()
				workerCount := workerPoolSelf.workerCount
				workerSizeStandBy := workerPoolSelf.workerSizeStandBy
				workerSizeMaximum := workerPoolSelf.workerSizeMaximum
				workerPoolSelf.lock.RUnlock()
				if workerCount > workerSizeStandBy ||
					workerCount > workerSizeMaximum {
					break loopLabel
				}
			}
		}
	}()
}

// SetJobQueue Set the JobQueue(WARNING: if the pool has started to use, doing this is not safe)
func (workerPoolSelf *DefaultWorkerPool) SetJobQueue(jobQueue *fpgo.BufferedChannelQueue) *DefaultWorkerPool {
	workerPoolSelf.lock.Lock()
	workerPoolSelf.jobQueue = jobQueue
	workerPoolSelf.lock.Unlock()
	workerPoolSelf.notifyWorkers()
	return workerPoolSelf
}

// SetIsJobQueueClosedWhenClose Set is the JobQueue closed when the WorkerPool.Close()
func (workerPoolSelf *DefaultWorkerPool) SetIsJobQueueClosedWhenClose(isJobQueueClosedWhenClose bool) *DefaultWorkerPool {
	workerPoolSelf.lock.Lock()
	workerPoolSelf.isJobQueueClosedWhenClose = isJobQueueClosedWhenClose
	workerPoolSelf.lock.Unlock()
	return workerPoolSelf
}

// SetPanicHandler Set the panicHandler(handle/log panic inside workers)
func (workerPoolSelf *DefaultWorkerPool) SetPanicHandler(panicHandler func(interface{})) *DefaultWorkerPool {
	workerPoolSelf.lock.Lock()
	workerPoolSelf.panicHandler = panicHandler
	workerPoolSelf.lock.Unlock()
	return workerPoolSelf
}

// SetWorkerBatchSize Set the workerBatchSize(queued jobs number that every worker could have)
func (workerPoolSelf *DefaultWorkerPool) SetWorkerBatchSize(workerBatchSize int) *DefaultWorkerPool {
	workerPoolSelf.lock.Lock()
	workerPoolSelf.workerBatchSize = workerBatchSize
	workerPoolSelf.lock.Unlock()
	workerPoolSelf.notifyWorkers()
	return workerPoolSelf
}

// SetWorkerSizeStandBy Set the workerSizeStandBy
func (workerPoolSelf *DefaultWorkerPool) SetWorkerSizeStandBy(workerSizeStandBy int) *DefaultWorkerPool {
	workerPoolSelf.lock.Lock()
	workerPoolSelf.workerSizeStandBy = workerSizeStandBy
	workerPoolSelf.lock.Unlock()
	workerPoolSelf.notifyWorkers()
	return workerPoolSelf
}

// SetWorkerSizeMaximum Set the workerSizeMaximum
func (workerPoolSelf *DefaultWorkerPool) SetWorkerSizeMaximum(workerSizeMaximum int) *DefaultWorkerPool {
	workerPoolSelf.lock.Lock()
	workerPoolSelf.workerSizeMaximum = workerSizeMaximum
	workerPoolSelf.lock.Unlock()
	workerPoolSelf.notifyWorkers()
	return workerPoolSelf
}

// SetSpawnWorkerDuration Set the spawnWorkerDuration(Checking repeating by the interval/duration)
func (workerPoolSelf *DefaultWorkerPool) SetSpawnWorkerDuration(spawnWorkerDuration time.Duration) *DefaultWorkerPool {
	workerPoolSelf.lock.Lock()
	defer workerPoolSelf.lock.Unlock()
	workerPoolSelf.spawnWorkerDuration = spawnWorkerDuration
	return workerPoolSelf
}

// SetWorkerExpiryDuration The worker would be dead if the worker is idle without jobs over the duration
func (workerPoolSelf *DefaultWorkerPool) SetWorkerExpiryDuration(workerExpiryDuration time.Duration) *DefaultWorkerPool {
	workerPoolSelf.lock.Lock()
	workerPoolSelf.workerExpiryDuration = workerExpiryDuration
	workerPoolSelf.lock.Unlock()
	workerPoolSelf.notifyWorkers()
	return workerPoolSelf
}

// SetWorkerJamDuration A new worker would be created if there's no available worker to do jobs over the duration
func (workerPoolSelf *DefaultWorkerPool) SetWorkerJamDuration(workerJamDuration time.Duration) *DefaultWorkerPool {
	workerPoolSelf.lock.Lock()
	workerPoolSelf.workerJamDuration = workerJamDuration
	workerPoolSelf.lock.Unlock()
	workerPoolSelf.notifyWorkers()
	return workerPoolSelf
}

// SetScheduleRetryInterval Retry interval for ScheduleWithTimeout
func (workerPoolSelf *DefaultWorkerPool) SetScheduleRetryInterval(scheduleRetryInterval time.Duration) *DefaultWorkerPool {
	workerPoolSelf.lock.Lock()
	workerPoolSelf.scheduleRetryInterval = scheduleRetryInterval
	workerPoolSelf.lock.Unlock()
	workerPoolSelf.notifyWorkers()
	return workerPoolSelf
}

// SetDefaultWorkerPoolSettings Set the defaultWorkerPoolSettings
func (workerPoolSelf *DefaultWorkerPool) SetDefaultWorkerPoolSettings(defaultWorkerPoolSettings DefaultWorkerPoolSettings) *DefaultWorkerPool {
	workerPoolSelf.lock.Lock()
	workerPoolSelf.DefaultWorkerPoolSettings = defaultWorkerPoolSettings
	workerPoolSelf.lock.Unlock()
	workerPoolSelf.notifyWorkers()
	return workerPoolSelf
}

// IsClosed Is the DefaultWorkerPool closed
func (workerPoolSelf *DefaultWorkerPool) IsClosed() bool {
	return workerPoolSelf.isClosed.Get()
}

func (workerPoolSelf *DefaultWorkerPool) getWorkerCount() int {
	workerPoolSelf.lock.RLock()
	defer workerPoolSelf.lock.RUnlock()
	return workerPoolSelf.workerCount
}

// Close Close the DefaultWorkerPool
func (workerPoolSelf *DefaultWorkerPool) Close() {
	if !workerPoolSelf.isClosed.CompareAndSwap(false, true) {
		return
	}
	close(workerPoolSelf.done)

	workerPoolSelf.lock.RLock()
	isJobQueueClosedWhenClose := workerPoolSelf.isJobQueueClosedWhenClose
	jobQueue := workerPoolSelf.jobQueue
	workerPoolSelf.lock.RUnlock()

	if isJobQueueClosedWhenClose && jobQueue != nil {
		jobQueue.Close()
	}
}

// Schedule Schedule the Job
func (workerPoolSelf *DefaultWorkerPool) Schedule(fn func()) error {
	if workerPoolSelf.IsClosed() {
		return ErrWorkerPoolIsClosed
	}

	workerPoolSelf.lock.RLock()
	jobQueue := workerPoolSelf.jobQueue
	workerPoolSelf.lock.RUnlock()
	if jobQueue == nil {
		return ErrWorkerPoolJobQueueIsNil
	}

	defer workerPoolSelf.spawnWorkerCh.Offer(1)

	err := jobQueue.Offer(fn)
	if err == fpgo.ErrQueueIsFull {
		return ErrWorkerPoolJobQueueIsFull
	}

	return err
}

// ScheduleWithTimeout Schedule the Job with timeout
func (workerPoolSelf *DefaultWorkerPool) ScheduleWithTimeout(fn func(), timeout time.Duration) error {
	err := workerPoolSelf.Schedule(fn)
	if err != ErrWorkerPoolJobQueueIsFull {
		return err
	}

	workerPoolSelf.lock.RLock()
	retryInterval := workerPoolSelf.scheduleRetryInterval
	workerPoolSelf.lock.RUnlock()
	if retryInterval > timeout/3 {
		// retryInterval = timeout * 95 / 100 / 3
		retryInterval = timeout / 3
	}
	deadline := time.Now().Add(timeout)

	for {
		if workerPoolSelf.IsClosed() {
			return ErrWorkerPoolIsClosed
		}

		err = workerPoolSelf.Schedule(fn)
		if err != ErrWorkerPoolJobQueueIsFull {
			return err
		}

		if time.Now().After(deadline) {
			return ErrWorkerPoolScheduleTimeout
		}
		time.Sleep(retryInterval)
	}
}

// Invokable

// Invokable Invokable inspired by Java ExecutorService
type Invokable interface {
	Invoke(val interface{})
	InvokeWithTimeout(val interface{}, timeout time.Duration) error
}

// DefaultInvokable DefaultInvokable inspired by Java ExecutorService
type DefaultInvokable struct {
	workerPool WorkerPool
	callee     func(interface{})
}

// NewDefaultInvokable New a DefaultInvokable on the workerPool
func NewDefaultInvokable(workerPool WorkerPool, callee func(interface{})) *DefaultInvokable {
	return &DefaultInvokable{
		workerPool: workerPool,
		callee:     callee,
	}
}

// SetWorkerPool Set the WorkerPool
func (invokableSelf *DefaultInvokable) SetWorkerPool(workerPool WorkerPool) *DefaultInvokable {
	invokableSelf.workerPool = workerPool
	return invokableSelf
}

// SetCallee Set the Callee
func (invokableSelf *DefaultInvokable) SetCallee(callee func(interface{})) *DefaultInvokable {
	invokableSelf.callee = callee
	return invokableSelf
}

// Invoke Invoke the job (non-blocking)
func (invokableSelf *DefaultInvokable) Invoke(val interface{}) {
	callee := invokableSelf.callee
	invokableSelf.workerPool.Schedule(func() {
		callee(val)
	})
}

// InvokeWithTimeout Invoke the job with timeout (blocking, by workerPool.ScheduleWithTimeout())
func (invokableSelf *DefaultInvokable) InvokeWithTimeout(val interface{}, timeout time.Duration) error {
	callee := invokableSelf.callee
	return invokableSelf.workerPool.ScheduleWithTimeout(func() {
		callee(val)
	}, timeout)
}
