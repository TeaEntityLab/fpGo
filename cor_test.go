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

	expectedInt = 4
	actualInt = 0
	wg.Add(1)
	// Cor c1
	var c1 *CorDef
	c1 = Cor.New(func() {
		self := c1

		logMessage("c1 effect")
		initVal := self.Yield()
		logMessage("c1 initVal", initVal)
		logMessage("c1 initVal(unwrap)", *initVal)
		v, _ := Monad.Just(initVal).ToInt()
		// v := 0
		receive := self.YieldRef(Monad.JustVal(v + 1).Ref())
		logMessage("c1 yield initVal+1 & receive", receive)
		logMessage("c1", self.Yield())
	})
	// Testee
	testee = Cor.New(func() {
		self := testee

		logMessage("cor", "initialized")

		v, _ := Monad.Just(self.Yield()).ToInt()
		actualInt = v + 1

		logMessage("c1", c1.IsDone())
		m := Monad.Just(self.YieldFrom(c1, nil)).ToMonad()
		logMessage("c1", c1.IsDone())

		logMessage(m)

		v, _ = m.ToInt()
		actualInt += v

		logMessage("received", v)

		wg.Done()
	})

	c1.Start(Monad.JustVal(1).Ref())
	testee.Start(Monad.JustVal(1).Ref())

	wg.Wait()
	assert.Equal(t, expectedInt, actualInt)
}
