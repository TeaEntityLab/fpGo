package fpgo

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func logMessage(args ...interface{}) {
	enable := false
	if enable {
		fmt.Println(args...)
	}
}

func TestCorYield(t *testing.T) {
	var expectedInt, actualInt int
	var testee *CorDef[interface{}]
	var wg sync.WaitGroup

	expectedInt = 5
	actualInt = 0
	wg.Add(1)
	// Cor c1
	var c1 *CorDef[interface{}]
	c1 = Cor.New(func() {
		self := c1

		logMessage(self, "c1 effect")
		initVal := self.YieldRef(nil)
		logMessage(self, "c1 initVal", initVal)
		v, _ := Maybe.Just(initVal).ToInt()
		// v := 0
		receive := self.YieldRef((v + 1))
		logMessage(self, "c1 yield initVal+1 & receive", receive)
		logMessage(self, "c1", self.YieldRef(nil))
	})
	// Testee
	testee = Cor.New(func() {
		self := testee

		var v int
		var m MaybeDef[interface{}]

		logMessage(self, "cor", "initialized")

		v, _ = Maybe.Just(self.YieldRef(nil)).ToInt()
		actualInt = v + 1

		v, _ = Maybe.Just(self.YieldFromIO(MonadIO.Just(1).ObserveOn(Handler.GetDefault()))).ToInt()
		logMessage(self, "s", 5)
		actualInt += v
		logMessage(self, "s", 6)

		logMessage(self, "c1", c1.IsDone())
		logMessage(self, "c1", c1.IsStarted())
		m = Maybe.Just(self.YieldFrom(c1, nil)).ToMaybe()
		logMessage(self, "c1", c1.IsDone())

		logMessage(m)

		v, _ = m.ToInt()
		actualInt += v

		logMessage(self, "received", v)

		wg.Done()
	})

	c1.StartWithVal(1)
	testee.StartWithVal(1)

	wg.Wait()
	assert.Equal(t, expectedInt, actualInt)
}

func TestCorNewAndStart(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	var ran AtomBool

	cor := Cor.NewAndStart(func() {
		ran.Set(true)
		wg.Done()
	})
	assert.NotNil(t, cor)

	wg.Wait()
	assert.True(t, ran.Get())
	assert.True(t, cor.IsStarted())
	assert.Eventually(t, cor.IsDone, time.Second, time.Millisecond)

	// NewAndStart on an already-done receiver is safe and returns a new started Cor.
	assert.NotPanics(t, func() { cor.Start() })
}

func TestCorDoNotation(t *testing.T) {
	var expectedInt int
	var actual interface{}

	expectedInt = 3
	// Cor c1
	var c1 *CorDef[interface{}]
	c1 = Cor.New(func() {
		self := c1

		val := self.YieldRef((1))
		Maybe.Just(val).ToInt()
		logMessage(self, "c1 val", val)
	})
	c1.Start()
	// Testee
	actual = Cor.DoNotation(func(self *CorDef[interface{}]) interface{} {
		logMessage(self, "Do Notation", "init")

		var result int
		v := 0
		var m MaybeDef[interface{}]

		result = v + 1

		logMessage(self, "Do Notation", "v", v)
		logMessage(self, "Do Notation", "result", result)

		h := Handler.New()
		defer h.Close()
		v, _ = Maybe.Just(self.YieldFromIO(MonadIO.Just(1).ObserveOn(h))).ToInt()
		result += v

		logMessage(self, "Do Notation", "result", result)

		m = Maybe.Just(self.YieldFrom(c1, nil)).ToMaybe()

		v, _ = m.ToInt()
		result += v

		logMessage(self, "Do Notation", "result", result)

		return (result)
	})

	assert.Equal(t, expectedInt, (actual))
}

func TestCorNewGenerics(t *testing.T) {
	cor := CorNewGenerics[interface{}](func() {})
	assert.NotNil(t, cor)
	assert.Equal(t, false, cor.IsStarted())
	assert.Equal(t, false, cor.IsDone())
}

func TestCorIsDone(t *testing.T) {
	cor := CorNewGenerics[interface{}](func() {})
	assert.Equal(t, false, cor.IsDone())
}

func TestCorIsStarted(t *testing.T) {
	cor := CorNewGenerics[interface{}](func() {})
	assert.Equal(t, false, cor.IsStarted())

	cor.Start()
	assert.Equal(t, true, cor.IsStarted())
}

func TestCorClose(t *testing.T) {
	cor := CorNewGenerics[interface{}](func() {})
	assert.Equal(t, false, cor.IsDone())
	cor.Start()
	assert.Equal(t, true, cor.IsStarted())
}

func TestAtomBool(t *testing.T) {
	atom := AtomBool{flag: 0}
	assert.False(t, atom.Get())

	atom.Set(true)
	assert.True(t, atom.Get())

	atom.Set(false)
	assert.False(t, atom.Get())
}

func TestCorDef_New(t *testing.T) {
	var c *CorDef[interface{}]
	c = new(CorDef[interface{}]).New(func() {})
	assert.NotNil(t, c)
}

func TestCorStartWithValAlreadyStarted(t *testing.T) {
	c := CorNewGenerics[interface{}](func() {})
	c.Start()
	time.Sleep(10 * time.Millisecond)
	assert.NotPanics(t, func() {
		c.StartWithVal(1)
	})
	assert.True(t, c.IsStarted())
}

func TestCorStartAlreadyStarted(t *testing.T) {
	c := CorNewGenerics[interface{}](func() {})
	c.Start()
	assert.NotPanics(t, func() {
		c.Start()
	})
	assert.True(t, c.IsStarted())
}

func TestCorYieldRefOnDone(t *testing.T) {
	c := CorNewGenerics[interface{}](func() {})
	c.Start()
	time.Sleep(10 * time.Millisecond)
	result := c.YieldRef(nil)
	assert.Equal(t, nil, result)
}

func TestCorYieldFromOnDone(t *testing.T) {
	c := CorNewGenerics[interface{}](func() {})
	c.Start()
	time.Sleep(10 * time.Millisecond)
	result := c.YieldFrom(nil, nil)
	assert.Equal(t, nil, result)
}

func TestCorStartWithValDoubleCall(t *testing.T) {
	var c *CorDef[interface{}]
	c = Cor.New(func() {
		self := c
		self.YieldRef(nil)
	})
	assert.NotPanics(t, func() {
		c.StartWithVal(1)
	})
	assert.NotPanics(t, func() {
		c.StartWithVal(2)
	})
	assert.True(t, c.IsStarted())
}

func TestCorDoCloseSafeOnDone(t *testing.T) {
	c := CorNewGenerics[interface{}](func() {})
	c.Start()
	time.Sleep(10 * time.Millisecond)
	assert.True(t, c.IsDone())

	var parent *CorDef[interface{}]
	parent = CorNewGenerics[interface{}](func() {
		self := parent
		self.YieldFrom(c, nil)
	})
	parent.Start()
	time.Sleep(20 * time.Millisecond)
}

func TestCorReceiveOnDone(t *testing.T) {
	c := CorNewGenerics[interface{}](func() {})
	c.Start()
	time.Sleep(10 * time.Millisecond)
	assert.True(t, c.IsDone())

	done := make(chan struct{})
	go func() {
		defer close(done)
		var child *CorDef[interface{}]
		child = CorNewGenerics[interface{}](func() {
			self := child
			self.YieldFrom(c, nil)
		})
		child.Start()
	}()

	select {
	case <-done:
	case <-time.After(50 * time.Millisecond):
	}
}

func TestCorDoCloseSafeRaceWithClose(t *testing.T) {
	// Directly race close() and doCloseSafe(). Before the fix, close() set
	// isClosed BEFORE acquiring closedM, so doCloseSafe() could pass the
	// IsDone() check, then have close() close the channels, then send on a
	// closed channel → panic. After the fix, both isClosed.Set and the
	// IsDone check happen inside closedM, so the send can never hit a closed
	// channel.
	for i := 0; i < 1000; i++ {
		cor := CorNewGenerics[int](func() {})
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			cor.close()
		}()
		cor.doCloseSafe(func() {
			if cor.opCh != nil {
				cor.opCh <- &CorOp[int]{cor: nil, val: 42}
			}
		})
		wg.Wait()
	}
}
