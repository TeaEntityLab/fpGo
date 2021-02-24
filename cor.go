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
type CorOp struct {
	cor *CorDef
	val interface{}
}

// CorDef Cor Coroutine inspired by Python/Ecmascript/Lua
type CorDef struct {
	isStarted AtomBool
	isClosed  AtomBool
	closedM   sync.Mutex

	opCh     *chan *CorOp
	resultCh *chan interface{}

	effect func()
}

// New New a Cor instance
func (corSelf *CorDef) New(effect func()) *CorDef {
	opCh := make(chan *CorOp, 5)
	resultCh := make(chan interface{}, 5)
	cor := &CorDef{effect: effect, opCh: &opCh, resultCh: &resultCh, isStarted: AtomBool{flag: 0}}
	return cor
}

// NewAndStart New a Cor instance and start it immediately
func (corSelf *CorDef) NewAndStart(effect func()) *CorDef {
	cor := corSelf.New(effect)
	cor.Start()
	return cor
}

// DoNotation Do Notation by function (inspired by Haskell one)
func (corSelf *CorDef) DoNotation(effect func(*CorDef) interface{}) interface{} {
	var result interface{}

	var wg sync.WaitGroup
	wg.Add(1)
	var cor *CorDef
	cor = corSelf.New(func() {
		result = effect(cor)
		wg.Done()
	})
	cor.Start()
	wg.Wait()

	return result
}

// StartWithVal Start the Cor with an initial value
func (corSelf *CorDef) StartWithVal(in interface{}) {
	if corSelf.IsDone() || corSelf.isStarted.Get() {
		return
	}

	corSelf.receive(nil, in)
	corSelf.Start()
}

// Start Start the Cor
func (corSelf *CorDef) Start() {
	if corSelf.IsDone() || corSelf.isStarted.Get() {
		return
	}
	corSelf.isStarted.Set(true)

	go func() {
		corSelf.effect()
		corSelf.close()
	}()
}

// Yield Yield back(nil)
func (corSelf *CorDef) Yield() interface{} {
	return corSelf.YieldRef(nil)
}

// YieldRef Yield a value
func (corSelf *CorDef) YieldRef(out interface{}) interface{} {
	var result interface{}
	if corSelf.IsDone() {
		return result
	}

	var op *CorOp
	var more bool
	// fmt.Println(corSelf, "Wait for", "op")
	op, more = <-*corSelf.opCh
	// fmt.Println(corSelf, "Wait for", "op", "done")

	if more && op != nil && op.cor != nil {
		cor := op.cor
		cor.doCloseSafe(func() {
			*cor.resultCh <- out
		})
	}
	result = op.val

	return result
}

// YieldFrom Yield from a given Cor
func (corSelf *CorDef) YieldFrom(target *CorDef, in interface{}) interface{} {
	var result interface{}
	if corSelf.IsDone() {
		return result
	}

	target.receive(corSelf, in)

	// fmt.Println(corSelf, "Wait for", "result")
	result, _ = <-*corSelf.resultCh
	// fmt.Println(corSelf, "Wait for", "result", "done")

	return result
}
func (corSelf *CorDef) receive(cor *CorDef, in interface{}) {
	corSelf.doCloseSafe(func() {
		if corSelf.opCh != nil {
			// fmt.Println(corSelf, "Wait for", "receive", cor, in)
			*(corSelf.opCh) <- &CorOp{cor: cor, val: in}
			// fmt.Println(corSelf, "Wait for", "receive", "done")
		}
	})
}

// YieldFromIO Yield from a given MonadIO
func (corSelf *CorDef) YieldFromIO(target *MonadIODef) interface{} {
	var result interface{}

	var wg sync.WaitGroup
	wg.Add(1)
	target.SubscribeOn(nil).Subscribe(Subscription{
		OnNext: func(in interface{}) {
			result = in
			wg.Done()
		},
	})
	wg.Wait()

	return result
}

// IsDone Is the Cor done
func (corSelf *CorDef) IsDone() bool {
	return corSelf.isClosed.Get()
}

// IsStarted Is the Cor started
func (corSelf *CorDef) IsStarted() bool {
	return corSelf.isStarted.Get()
}
func (corSelf *CorDef) close() {
	corSelf.isClosed.Set(true)

	corSelf.closedM.Lock()
	if corSelf.resultCh != nil {
		close(*corSelf.resultCh)
	}
	if corSelf.opCh != nil {
		close(*corSelf.opCh)
	}
	corSelf.closedM.Unlock()
}
func (corSelf *CorDef) doCloseSafe(fn func()) {
	if corSelf.IsDone() {
		return
	}
	corSelf.closedM.Lock()
	fn()
	corSelf.closedM.Unlock()
}

// Cor Cor utils instance
var Cor CorDef
