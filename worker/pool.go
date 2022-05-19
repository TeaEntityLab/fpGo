package worker

import (
	"errors"
	"log"
	"runtime"
	"sync"
	"time"

	fpgo "github.com/TeaEntityLab/fpGo/v2"
)

var (
	// ErrWorkerPoolJobQueueIsFull WorkerPool JobQueue Is Full
	ErrWorkerPoolJobQueueIsFull = errors.New("workerPool JobQueue is full")
	// ErrWorkerPoolIsClosed WorkerPool Is Closed
	ErrWorkerPoolIsClosed = errors.New("workerPool is closed")
	// ErrWorkerPoolScheduleTimeout WorkerPool Schedule Timeout
	ErrWorkerPoolScheduleTimeout = errors.New("workerPool schedule timeout")
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

	workerSizeStandBy    int
	workerSizeMaximum    int
	spawnWorkerDuration  time.Duration
	workerExpiryDuration time.Duration
	workerJamDuration    time.Duration

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
	panicHandler:              defaultPanicHandler,
}

// DefaultWorkerPool DefaultWorkerPool inspired by Java ExecutorService
type DefaultWorkerPool struct {
	isClosed fpgo.AtomBool
	lock     sync.RWMutex

	jobQueue *fpgo.BufferedChannelQueue[func()]

	workerCount   int
	spawnWorkerCh fpgo.ChannelQueue[int]
	lastAliveTime time.Time

	scheduleWaitCh    fpgo.ChannelQueue[int]
	scheduleWaitCount int

	// Settings
	DefaultWorkerPoolSettings
}

// NewDefaultWorkerPool New a DefaultWorkerPool
func NewDefaultWorkerPool(jobQueue *fpgo.BufferedChannelQueue[func()], settings *DefaultWorkerPoolSettings) *DefaultWorkerPool {
	if settings == nil {
		settings = defaultDefaultWorkerSettings
	}
	workerPool := &DefaultWorkerPool{
		jobQueue: jobQueue,

		spawnWorkerCh:  fpgo.NewChannelQueue[int](1),
		scheduleWaitCh: fpgo.NewChannelQueue[int](1),

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
	if batchSize < 1 {
		batchSize = 1
	}
	expectedWorkerCount := workerPoolSelf.jobQueue.Count() / batchSize
	if workerPoolSelf.jobQueue.Count()%batchSize > 0 {
		expectedWorkerCount++
	}
	if workerPoolSelf.workerSizeStandBy > expectedWorkerCount {
		expectedWorkerCount = workerPoolSelf.workerSizeStandBy
	}
	if workerPoolSelf.workerSizeMaximum > 0 &&
		expectedWorkerCount > workerPoolSelf.workerSizeMaximum {
		expectedWorkerCount = workerPoolSelf.workerSizeMaximum
	}
	// Avoid Jam if (now - lastAliveTime) is over workerJamDuration
	if time.Now().Sub(workerPoolSelf.lastAliveTime) > workerPoolSelf.workerJamDuration &&
		workerPoolSelf.workerCount >= expectedWorkerCount {
		expectedWorkerCount++
	}
	workerPoolSelf.lock.RUnlock()

	if workerPoolSelf.workerCount < expectedWorkerCount {
		for i := workerPoolSelf.workerCount; i < expectedWorkerCount; i++ {
			workerPoolSelf.generateWorkerWithMaximum(expectedWorkerCount)
		}
	}
}

// PreAllocWorkerSize PreAllocate Workers
func (workerPoolSelf *DefaultWorkerPool) PreAllocWorkerSize(preAllocWorkerSize int) {
	for i := workerPoolSelf.workerCount; i < preAllocWorkerSize; i++ {
		workerPoolSelf.generateWorkerWithMaximum(preAllocWorkerSize)
	}
}

func (workerPoolSelf *DefaultWorkerPool) spawnLoop() {
	defer func() {
		if panic := recover(); panic != nil {
			defaultPanicHandler(panic)
		}
	}()

	for range workerPoolSelf.spawnWorkerCh {
		if workerPoolSelf.IsClosed() {
			break
		}

		workerPoolSelf.trySpawn()

		time.Sleep(workerPoolSelf.spawnWorkerDuration)
	}
}

func (workerPoolSelf *DefaultWorkerPool) notifyWorkers() {
	if workerPoolSelf.workerCount < workerPoolSelf.workerSizeStandBy || workerPoolSelf.jobQueue.Count() > 0 {
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

	go func() {
		// Recover & Recycle
		defer func() {
			if panic := recover(); panic != nil {
				if handler := workerPoolSelf.panicHandler; handler != nil {
					handler(panic)
				}
			}

			workerPoolSelf.lock.Lock()
			workerPoolSelf.workerCount--
			workerPoolSelf.lock.Unlock()
			// fmt.Println("Terminated")
		}()

		// Do Jobs
	loopLabel:
		for {
			workerPoolSelf.lastAliveTime = time.Now()
			if workerPoolSelf.scheduleWaitCount > 0 {
				workerPoolSelf.scheduleWaitCh.Offer(1)
			}

			select {
			case job := <-workerPoolSelf.jobQueue.GetChannel():
				// fmt.Println("GetJob")
				if job != nil {
					job()
					// fmt.Println("DoJob")
				}
			case <-time.After(workerPoolSelf.workerExpiryDuration):
				workerPoolSelf.lock.RLock()
				workerCount := workerPoolSelf.workerCount
				if workerCount > workerPoolSelf.workerSizeStandBy ||
					workerCount > workerPoolSelf.workerSizeMaximum {
					workerPoolSelf.lock.RUnlock()
					break loopLabel
				}
				workerPoolSelf.lock.RUnlock()
			}
		}
	}()
}

// SetJobQueue Set the JobQueue(WARNING: if the pool has started to use, doing this is not safe)
func (workerPoolSelf *DefaultWorkerPool) SetJobQueue(jobQueue *fpgo.BufferedChannelQueue[func()]) *DefaultWorkerPool {
	workerPoolSelf.jobQueue = jobQueue
	return workerPoolSelf
}

// SetIsJobQueueClosedWhenClose Set is the JobQueue closed when the WorkerPool.Close()
func (workerPoolSelf *DefaultWorkerPool) SetIsJobQueueClosedWhenClose(isJobQueueClosedWhenClose bool) *DefaultWorkerPool {
	workerPoolSelf.isJobQueueClosedWhenClose = isJobQueueClosedWhenClose
	return workerPoolSelf
}

// SetPanicHandler Set the panicHandler(handle/log panic inside workers)
func (workerPoolSelf *DefaultWorkerPool) SetPanicHandler(panicHandler func(interface{})) *DefaultWorkerPool {
	workerPoolSelf.panicHandler = panicHandler
	return workerPoolSelf
}

// SetWorkerBatchSize Set the workerBatchSize(queued jobs number that every worker could have)
func (workerPoolSelf *DefaultWorkerPool) SetWorkerBatchSize(workerBatchSize int) *DefaultWorkerPool {
	workerPoolSelf.workerBatchSize = workerBatchSize
	workerPoolSelf.notifyWorkers()
	return workerPoolSelf
}

// SetWorkerSizeStandBy Set the workerSizeStandBy
func (workerPoolSelf *DefaultWorkerPool) SetWorkerSizeStandBy(workerSizeStandBy int) *DefaultWorkerPool {
	workerPoolSelf.workerSizeStandBy = workerSizeStandBy
	workerPoolSelf.notifyWorkers()
	return workerPoolSelf
}

// SetWorkerSizeMaximum Set the workerSizeMaximum
func (workerPoolSelf *DefaultWorkerPool) SetWorkerSizeMaximum(workerSizeMaximum int) *DefaultWorkerPool {
	workerPoolSelf.workerSizeMaximum = workerSizeMaximum
	workerPoolSelf.notifyWorkers()
	return workerPoolSelf
}

// SetSpawnWorkerDuration Set the spawnWorkerDuration(Checking repeating by the interval/duration)
func (workerPoolSelf *DefaultWorkerPool) SetSpawnWorkerDuration(spawnWorkerDuration time.Duration) *DefaultWorkerPool {
	workerPoolSelf.spawnWorkerDuration = spawnWorkerDuration
	return workerPoolSelf
}

// SetWorkerExpiryDuration The worker would be dead if the worker is idle without jobs over the duration
func (workerPoolSelf *DefaultWorkerPool) SetWorkerExpiryDuration(workerExpiryDuration time.Duration) *DefaultWorkerPool {
	workerPoolSelf.workerExpiryDuration = workerExpiryDuration
	workerPoolSelf.notifyWorkers()
	return workerPoolSelf
}

// SetWorkerJamDuration A new worker would be created if there's no available worker to do jobs over the duration
func (workerPoolSelf *DefaultWorkerPool) SetWorkerJamDuration(workerJamDuration time.Duration) *DefaultWorkerPool {
	workerPoolSelf.workerJamDuration = workerJamDuration
	workerPoolSelf.notifyWorkers()
	return workerPoolSelf
}

// SetDefaultWorkerPoolSettings Set the defaultWorkerPoolSettings
func (workerPoolSelf *DefaultWorkerPool) SetDefaultWorkerPoolSettings(defaultWorkerPoolSettings DefaultWorkerPoolSettings) *DefaultWorkerPool {
	workerPoolSelf.DefaultWorkerPoolSettings = defaultWorkerPoolSettings
	workerPoolSelf.notifyWorkers()
	return workerPoolSelf
}

// IsClosed Is the DefaultWorkerPool closed
func (workerPoolSelf *DefaultWorkerPool) IsClosed() bool {
	return workerPoolSelf.isClosed.Get()
}

// Close Close the DefaultWorkerPool
func (workerPoolSelf *DefaultWorkerPool) Close() {
	if workerPoolSelf.IsClosed() {
		return
	}
	workerPoolSelf.isClosed.Set(true)

	if workerPoolSelf.isJobQueueClosedWhenClose {
		workerPoolSelf.jobQueue.Close()
	}
}

// Schedule Schedule the Job
func (workerPoolSelf *DefaultWorkerPool) Schedule(fn func()) error {
	if workerPoolSelf.IsClosed() {
		return ErrWorkerPoolIsClosed
	}
	defer workerPoolSelf.spawnWorkerCh.Offer(1)

	err := workerPoolSelf.jobQueue.Offer(fn)
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

	deadline := time.Now().Add(timeout)
	workerPoolSelf.addScheduleWaitCount(1)
	go workerPoolSelf.addScheduleWaitCount(-1)

	for {
		if workerPoolSelf.IsClosed() {
			return ErrWorkerPoolIsClosed
		}

		err = workerPoolSelf.Schedule(fn)
		if err != ErrWorkerPoolJobQueueIsFull {
			return err
		}

		select {
		case <-workerPoolSelf.scheduleWaitCh:
			continue
		case <-time.After(deadline.Sub(time.Now())):
			return ErrWorkerPoolScheduleTimeout
		}
	}
	return err
}

func (workerPoolSelf *DefaultWorkerPool) addScheduleWaitCount(amount int) {
	workerPoolSelf.lock.Lock()
	workerPoolSelf.scheduleWaitCount += amount
	workerPoolSelf.lock.Unlock()
}

// Invokable

// Invokable Invokable inspired by Java ExecutorService
type Invokable[T any] interface {
	Invoke(val T)
	InvokeWithTimeout(val T, timeout time.Duration) error
}

// DefaultInvokable DefaultInvokable inspired by Java ExecutorService
type DefaultInvokable[T any] struct {
	workerPool WorkerPool
	callee     func(T)
}

// NewDefaultInvokable New a DefaultInvokable on the workerPool
func NewDefaultInvokable[T any](workerPool WorkerPool, callee func(T)) *DefaultInvokable[T] {
	return &DefaultInvokable[T]{
		workerPool: workerPool,
		callee:     callee,
	}
}

// SetWorkerPool Set the WorkerPool
func (invokableSelf *DefaultInvokable[T]) SetWorkerPool(workerPool WorkerPool) *DefaultInvokable[T] {
	invokableSelf.workerPool = workerPool
	return invokableSelf
}

// SetCallee Set the Callee
func (invokableSelf *DefaultInvokable[T]) SetCallee(callee func(T)) *DefaultInvokable[T] {
	invokableSelf.callee = callee
	return invokableSelf
}

// Invoke Invoke the job (non-blocking)
func (invokableSelf *DefaultInvokable[T]) Invoke(val T) {
	callee := invokableSelf.callee
	invokableSelf.workerPool.Schedule(func() {
		callee(val)
	})
}

// InvokeWithTimeout Invoke the job with timeout (blocking, by workerPool.ScheduleWithTimeout())
func (invokableSelf *DefaultInvokable[T]) InvokeWithTimeout(val T, timeout time.Duration) error {
	callee := invokableSelf.callee
	return invokableSelf.workerPool.ScheduleWithTimeout(func() {
		callee(val)
	}, timeout)
}
