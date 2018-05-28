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
	isClosed AtomBool
	closedM  sync.Mutex

	opCh     *chan *CorOp
	resultCh *chan *interface{}

	effect func()
}

func (self *CorDef) New(effect func()) *CorDef {
	opCh := make(chan *CorOp, 3)
	resultCh := make(chan *interface{}, 3)
	return &CorDef{effect: effect, opCh: &opCh, resultCh: &resultCh}
}
func (self *CorDef) Start(in *interface{}) {
	if self.IsDone() {
		return
	}

	go func() {
		self.receive(nil, in)
		self.effect()
		self.close()
	}()
}
func (self *CorDef) Yield() *interface{} {
	return self.YieldRef(nil)
}
func (self *CorDef) YieldRef(out *interface{}) *interface{} {
	if self.IsDone() {
		return nil
	}

	op := <-*self.opCh
	if op.cor != nil {
		cor := op.cor
		*cor.resultCh <- out
	}
	return op.val
}
func (self *CorDef) YieldFrom(target *CorDef, in *interface{}) *interface{} {
	if target.IsDone() {
		return nil
	}

	target.receive(self, in)
	result := <-*self.resultCh
	return result
}
func (self *CorDef) receive(cor *CorDef, in *interface{}) {
	if self.IsDone() {
		return
	}

	if self.opCh != nil {
		*(self.opCh) <- &CorOp{cor: cor, val: in}
	}
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
	self.closedM.Lock()
	self.closedM.Unlock()
	return self.isClosed.Get()
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

var Cor CorDef
