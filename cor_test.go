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
		v, _ := Maybe.Just(initVal).ToInt()
		// v := 0
		receive := self.YieldRef((v + 1))
		logMessage(self, "c1 yield initVal+1 & receive", receive)
		logMessage(self, "c1", self.Yield())
	})
	// Testee
	testee = Cor.New(func() {
		self := testee

		var v int
		var m MaybeDef

		logMessage(self, "cor", "initialized")

		v, _ = Maybe.Just(self.Yield()).ToInt()
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

func TestCorDoNotation(t *testing.T) {
	var expectedInt int
	var actual interface{}

	expectedInt = 3
	// Cor c1
	var c1 *CorDef
	c1 = Cor.New(func() {
		self := c1

		val := self.YieldRef((1))
		Maybe.Just(val).ToInt()
		logMessage(self, "c1 val", val)
	})
	c1.Start()
	// Testee
	actual = Cor.DoNotation(func(self *CorDef) interface{} {
		logMessage(self, "Do Notation", "init")

		var result int
		v := 0
		var m MaybeDef

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

	assert.NotPanics(t, func() { cor.Start() })
}

func TestCorStartAlreadyStarted(t *testing.T) {
	c := Cor.New(func() {})
	c.Start()
	assert.NotPanics(t, func() {
		c.Start()
	})
	assert.True(t, c.IsStarted())
}

func TestCorStartWithValAlreadyStarted(t *testing.T) {
	c := Cor.New(func() {})
	c.Start()
	time.Sleep(10 * time.Millisecond)
	assert.NotPanics(t, func() {
		c.StartWithVal(1)
	})
	assert.True(t, c.IsStarted())
}
