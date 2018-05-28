package fpGo

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorYield(t *testing.T) {
	var expectedInt = 0
	var actualInt = 0
	var testee *CorDef
	var wg sync.WaitGroup

	expectedInt = 2
	actualInt = 0
	wg.Add(1)

	var c1 *CorDef
	c1 = Cor.New(func() {
		c1.YieldRef(Monad.JustVal(1).Ref())
	})
	testee = Cor.New(func() {
		fmt.Println("cor 0" + string(len(*testee.opCh)))

		//*
		v, _ := Monad.Just(testee.Yield()).ToInt()
		/*/
		v, _ := Monad.Just(testee.YieldFromIO(MonadIO.JustVal(1))).ToInt()
		//*/
		actualInt = v + 1

		// v, _ = Monad.Just(testee.YieldFrom(c1, nil)).ToInt()
		// actualInt += v

		fmt.Println("cor 1")

		wg.Done()
	})
	testee.Start(Monad.JustVal(1).Ref())
	wg.Wait()
	assert.Equal(t, expectedInt, actualInt)
}
