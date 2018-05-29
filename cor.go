package fpGo

import (
	"sync"
	"sync/atomic"
)

type AtomBool struct{ flag int32 }

func (self *AtomBool) Set(value bool) {
	var i int32 = 0
	if value {
		i = 1
	}
	atomic.StoreInt32(&(self.flag), int32(i))
}
func (self *AtomBool) Get() bool {
	if atomic.LoadInt32(&(self.flag)) != 0 {
		return true
	}
	return false
}

type CorOp struct {
	cor *CorDef
	val *interface{}
}
type CorDef struct {
	isStarted AtomBool
	isClosed  AtomBool
	closedM   sync.Mutex
	wgStarted sync.WaitGroup

	opCh     *chan *CorOp
	resultCh *chan *interface{}

	effect func()
}

func (self *CorDef) New(effect func()) *CorDef {
	opCh := make(chan *CorOp, 5)
	resultCh := make(chan *interface{}, 5)
	cor := &CorDef{effect: effect, opCh: &opCh, resultCh: &resultCh, isStarted: AtomBool{flag: 0}}
	cor.wgStarted.Add(1)
	return cor
}
func (self *CorDef) Start(in *interface{}) {
	if self.IsDone() || self.isStarted.Get() {
		return
	}
	self.isStarted.Set(true)

	self.receive(nil, in)
	self.wgStarted.Done()
	go func() {
		self.effect()
		self.close()
	}()
}
func (self *CorDef) Yield() *interface{} {
	return self.YieldRef(nil)
}
func (self *CorDef) YieldRef(out *interface{}) *interface{} {
	var result *interface{} = nil

	self.doCloseSafe(func() {
		// fmt.Println(self, "Wait for", "op")
		op := <-*self.opCh
		// fmt.Println(self, "Wait for", "op", "done")
		if op.cor != nil {
			cor := op.cor
			*cor.resultCh <- out
		}
		result = op.val
	})

	return result
}
func (self *CorDef) YieldFrom(target *CorDef, in *interface{}) *interface{} {
	var result *interface{} = nil

	target.receive(self, in)
	self.doCloseSafe(func() {
		// fmt.Println(self, "Wait for", "result")
		result = <-*self.resultCh
		// fmt.Println(self, "Wait for", "result", "done")
	})
	return result
}
func (self *CorDef) receive(cor *CorDef, in *interface{}) {
	self.doCloseSafe(func() {
		if self.opCh != nil {
			// fmt.Println(self, "Wait for", "receive", cor, in)
			*(self.opCh) <- &CorOp{cor: cor, val: in}
			// fmt.Println(self, "Wait for", "receive", "done")
		}
	})
}
func (self *CorDef) YieldFromIO(target *MonadIODef) *interface{} {
	var result *interface{} = nil

	var wg sync.WaitGroup
	wg.Add(1)
	target.SubscribeOn(nil).Subscribe(Subscription{
		OnNext: func(in *interface{}) {
			result = in
			wg.Done()
		},
	})
	wg.Wait()

	return result
}
func (self *CorDef) IsDone() bool {
	return self.isClosed.Get()
}
func (self *CorDef) WaitStart() {
	self.wgStarted.Wait()
}
func (self *CorDef) close() {
	self.isClosed.Set(true)

	self.closedM.Lock()
	if self.resultCh != nil {
		close(*self.resultCh)
	}
	if self.opCh != nil {
		close(*self.opCh)
	}
	self.closedM.Unlock()
}
func (self *CorDef) doCloseSafe(fn func()) {
	self.closedM.Lock()
	if self.IsDone() {
		return
	}

	fn()
	self.closedM.Unlock()
}

var Cor CorDef
