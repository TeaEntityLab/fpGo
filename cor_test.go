package fpGo

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func logMessage(args ...interface{}) {
	enable := false
	if enable {
		fmt.Println(args...)
	}
}

func TestCorYield(t *testing.T) {
	var expectedInt = 0
	var actualInt = 0
	var testee *CorDef
	var wg sync.WaitGroup

	expectedInt = 5
	actualInt = 0
	wg.Add(1)
	// Cor c1
	var c1 *CorDef
	c1 = Cor.New(func() {
		self := c1

		logMessage(self, "c1 effect")
		initVal := self.Yield()
		logMessage(self, "c1 initVal", initVal)
		logMessage(self, "c1 initVal(unwrap)", *initVal)
		v, _ := Monad.Just(initVal).ToInt()
		// v := 0
		receive := self.YieldRef(PtrOf(v + 1))
		logMessage(self, "c1 yield initVal+1 & receive", receive)
		logMessage(self, "c1", self.Yield())
	})
	// Testee
	testee = Cor.New(func() {
		self := testee

		v := 0
		var m MonadDef

		logMessage(self, "cor", "initialized")

		v, _ = Monad.Just(self.Yield()).ToInt()
		actualInt = v + 1

		v, _ = Monad.Just(self.YieldFromIO(MonadIO.JustVal(1).ObserveOn(&Handler))).ToInt()
		logMessage(self, "s", 5)
		actualInt += v
		logMessage(self, "s", 6)

		logMessage(self, "c1", c1.IsDone())
		logMessage(self, "c1", c1.IsStarted())
		m = Monad.Just(self.YieldFrom(c1, nil)).ToMonad()
		logMessage(self, "c1", c1.IsDone())

		logMessage(m)

		v, _ = m.ToInt()
		actualInt += v

		logMessage(self, "received", v)

		wg.Done()
	})

	c1.StartWithRef(PtrOf(1))
	testee.StartWithRef(PtrOf(1))

	wg.Wait()
	assert.Equal(t, expectedInt, actualInt)
}

func TestCorDoNotation(t *testing.T) {
	var expectedInt = 0
	var actual *interface{} = nil

	expectedInt = 3
	actual = nil
	// Cor c1
	var c1 *CorDef
	c1 = Cor.NewAndStart(func() {
		self := c1

		val := self.YieldRef(PtrOf(1))
		Monad.Just(val).ToInt()
		logMessage(self, "c1 val", val)
	})
	// Testee
	actual = Cor.DoNotation(func(self *CorDef) *interface{} {
		logMessage(self, "Do Notation", "init")

		result := 0
		v := 0
		var m MonadDef

		result = v + 1

		logMessage(self, "Do Notation", "v", v)
		logMessage(self, "Do Notation", "result", result)

		v, _ = Monad.Just(self.YieldFromIO(MonadIO.JustVal(1).ObserveOn(&Handler))).ToInt()
		result += v

		logMessage(self, "Do Notation", "result", result)

		m = Monad.Just(self.YieldFrom(c1, nil)).ToMonad()

		v, _ = m.ToInt()
		result += v

		logMessage(self, "Do Notation", "result", result)

		return PtrOf(result)
	})

	assert.Equal(t, expectedInt, *(actual))
}
