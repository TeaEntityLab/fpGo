package fpgo

import (
	"sync"
	"sync/atomic"
)

// AtomBool Atomic Bool
type AtomBool struct{ flag int32 }

// Set Set the bool atomically
func (atomBoolSelf *AtomBool) Set(value bool) {
	var i int32
	i = 0
	if value {
		i = 1
	}
	atomic.StoreInt32(&(atomBoolSelf.flag), int32(i))
}

// Get Get the bool atomically
func (atomBoolSelf *AtomBool) Get() bool {
	if atomic.LoadInt32(&(atomBoolSelf.flag)) != 0 {
		return true
	}
	return false
}

// CorOp Cor Yield Operation/Delegation/Callback
type CorOp[T any] struct {
	cor *CorDef[T]
	val T
}

// CorDef Cor Coroutine inspired by Python/Ecmascript/Lua
type CorDef[T any] struct {
	isStarted AtomBool
	isClosed  AtomBool
	closedM   sync.Mutex

	opCh     chan *CorOp[T]
	resultCh chan T

	effect func()
}

// New New a Cor instance
func (corSelf *CorDef[T]) New(effect func()) *CorDef[interface{}] {
	return CorNewGenerics[interface{}](effect)
}

// New New a Cor instance
func CorNewGenerics[T any](effect func()) *CorDef[T] {
	opCh := make(chan *CorOp[T], 5)
	resultCh := make(chan T, 5)
	cor := &CorDef[T]{effect: effect, opCh: opCh, resultCh: resultCh, isStarted: AtomBool{flag: 0}}
	return cor
}

// NewAndStart New a Cor instance and start it immediately
func (corSelf *CorDef[T]) NewAndStart(effect func()) *CorDef[T] {
	cor := CorNewGenerics[T](effect)
	cor.Start()
	return cor
}

// DoNotation Do Notation by function (inspired by Haskell one)
func (corSelf *CorDef[T]) DoNotation(effect func(*CorDef[T]) T) T {
	var result T

	var wg sync.WaitGroup
	wg.Add(1)
	var cor *CorDef[T]
	cor = CorNewGenerics[T](func() {
		result = effect(cor)
		wg.Done()
	})
	cor.Start()
	wg.Wait()

	return result
}

// StartWithVal Start the Cor with an initial value
func (corSelf *CorDef[T]) StartWithVal(in T) {
	if corSelf.IsDone() || corSelf.isStarted.Get() {
		return
	}

	corSelf.receive(nil, in)
	corSelf.Start()
}

// Start Start the Cor
func (corSelf *CorDef[T]) Start() {
	if corSelf.IsDone() || corSelf.isStarted.Get() {
		return
	}
	corSelf.isStarted.Set(true)

	go func() {
		corSelf.effect()
		corSelf.close()
	}()
}

// // Yield Yield back(nil)
// func (corSelf *CorDef[T]) Yield() T {
// 	return corSelf.YieldRef(nil)
// }

// YieldRef Yield a value
func (corSelf *CorDef[T]) YieldRef(out T) T {
	var result T
	if corSelf.IsDone() {
		return result
	}

	var op *CorOp[T]
	var more bool
	// fmt.Println(corSelf, "Wait for", "op")
	op, more = <-corSelf.opCh
	// fmt.Println(corSelf, "Wait for", "op", "done")

	if more && op != nil && op.cor != nil {
		cor := op.cor
		cor.doCloseSafe(func() {
			cor.resultCh <- out
		})
	}
	result = op.val

	return result
}

// YieldFrom Yield from a given Cor
func (corSelf *CorDef[T]) YieldFrom(target *CorDef[T], in T) T {
	var result T
	if corSelf.IsDone() {
		return result
	}

	target.receive(corSelf, in)

	// fmt.Println(corSelf, "Wait for", "result")
	result, _ = <-corSelf.resultCh
	// fmt.Println(corSelf, "Wait for", "result", "done")

	return result
}

func (corSelf *CorDef[T]) receive(cor *CorDef[T], in T) {
	corSelf.doCloseSafe(func() {
		if corSelf.opCh != nil {
			// fmt.Println(corSelf, "Wait for", "receive", cor, in)
			corSelf.opCh <- &CorOp[T]{cor: cor, val: in}
			// fmt.Println(corSelf, "Wait for", "receive", "done")
		}
	})
}

// YieldFromIO Yield from a given MonadIO
func (corSelf *CorDef[T]) YieldFromIO(target *MonadIODef[T]) T {
	var result T

	var wg sync.WaitGroup
	wg.Add(1)
	target.SubscribeOn(nil).Subscribe(Subscription[T]{
		OnNext: func(in T) {
			result = in
			wg.Done()
		},
	})
	wg.Wait()

	return result
}

// IsDone Is the Cor done
func (corSelf *CorDef[T]) IsDone() bool {
	return corSelf.isClosed.Get()
}

// IsStarted Is the Cor started
func (corSelf *CorDef[T]) IsStarted() bool {
	return corSelf.isStarted.Get()
}

func (corSelf *CorDef[T]) close() {
	corSelf.isClosed.Set(true)

	corSelf.closedM.Lock()
	if corSelf.resultCh != nil {
		close(corSelf.resultCh)
	}
	if corSelf.opCh != nil {
		close(corSelf.opCh)
	}
	corSelf.closedM.Unlock()
}

func (corSelf *CorDef[T]) doCloseSafe(fn func()) {
	if corSelf.IsDone() {
		return
	}
	corSelf.closedM.Lock()
	fn()
	corSelf.closedM.Unlock()
}

// Cor Cor utils instance
var Cor CorDef[interface{}]
