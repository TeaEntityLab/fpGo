package fpGo

import (
	"sync"
	"sync/atomic"
)

type AtomBool struct{ flag int32 }

func (atomBoolSelf *AtomBool) Set(value bool) {
	var i int32 = 0
	if value {
		i = 1
	}
	atomic.StoreInt32(&(atomBoolSelf.flag), int32(i))
}
func (atomBoolSelf *AtomBool) Get() bool {
	if atomic.LoadInt32(&(atomBoolSelf.flag)) != 0 {
		return true
	}
	return false
}

type CorOp struct {
	cor *CorDef
	val interface{}
}
type CorDef struct {
	isStarted AtomBool
	isClosed  AtomBool
	closedM   sync.Mutex

	opCh     *chan *CorOp
	resultCh *chan interface{}

	effect func()
}

func (corSelf *CorDef) New(effect func()) *CorDef {
	opCh := make(chan *CorOp, 5)
	resultCh := make(chan interface{}, 5)
	cor := &CorDef{effect: effect, opCh: &opCh, resultCh: &resultCh, isStarted: AtomBool{flag: 0}}
	return cor
}
func (corSelf *CorDef) NewAndStart(effect func()) *CorDef {
	cor := corSelf.New(effect)
	cor.Start()
	return cor
}
func (corSelf *CorDef) DoNotation(effect func(*CorDef) interface{}) interface{} {
	var result interface{} = nil

	var wg sync.WaitGroup
	wg.Add(1)
	var cor *CorDef = nil
	cor = corSelf.New(func() {
		result = effect(cor)
		wg.Done()
	})
	cor.Start()
	wg.Wait()

	return result
}
func (corSelf *CorDef) StartWithVal(in interface{}) {
	if corSelf.IsDone() || corSelf.isStarted.Get() {
		return
	}

	corSelf.receive(nil, in)
	corSelf.Start()
}
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
func (corSelf *CorDef) Yield() interface{} {
	return corSelf.YieldRef(nil)
}
func (corSelf *CorDef) YieldRef(out interface{}) interface{} {
	var result interface{} = nil
	if corSelf.IsDone() {
		return result
	}

	var op *CorOp = nil
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
func (corSelf *CorDef) YieldFrom(target *CorDef, in interface{}) interface{} {
	var result interface{} = nil
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
func (corSelf *CorDef) YieldFromIO(target *MonadIODef) interface{} {
	var result interface{} = nil

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
func (corSelf *CorDef) IsDone() bool {
	return corSelf.isClosed.Get()
}
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

var Cor CorDef
