package fpgo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestActorCommon(t *testing.T) {
	var expectedInt int
	var actual int
	var resultChannel chan interface{}

	// Test difference channel size cases
	for channelSize := 0; channelSize < 5; channelSize++ {

		actual = 0
		expectedInt = 1400
		// Channel for results
		resultChannel = make(chan interface{}, channelSize+1)
		// Message CMDs
		cmdSpawn := "spawn"
		cmdShutdown := "shutdown"
		// Testee
		actorRoot := Actor.New(func(self *ActorDef[interface{}], input interface{}) {
			// SPAWN: for ROOT
			if input == cmdSpawn {
				self.Spawn(func(self *ActorDef[interface{}], input interface{}) {
					// SHUTDOWN: for Children
					if input == cmdShutdown {
						self.Close()
						return
					}

					// INT cases: Children
					val, _ := Maybe.Just(input).ToInt()
					resultChannel <- val * 10
				})
				return
			}
			// SHUTDOWN: for ROOT
			if input == cmdShutdown {
				for _, child := range self.children {
					child.Send(cmdShutdown)
				}
				self.Close()

				close(resultChannel)
				return
			}

			// INT cases: ROOT
			intVal, _ := Maybe.Just(input).ToInt()
			if intVal > 0 {
				for _, child := range self.children {
					child.Send(intVal)
				}
			}
		})

		// Sequential Send messages(async)
		go func() {
			actorRoot.Send(cmdSpawn)
			actorRoot.Send(10)
			actorRoot.Send(cmdSpawn)
			actorRoot.Send(20)
			actorRoot.Send(cmdSpawn)
			actorRoot.Send(30)
		}()

		i := 0
		for val := range resultChannel {
			intVal, _ := Maybe.Just(val).ToInt()
			actual += intVal

			i++
			if i == 5 {
				go actorRoot.Send(cmdShutdown)
			}
		}

		assert.Equal(t, expectedInt, actual)
	}
}

func TestActorAsk(t *testing.T) {
	var expectedInt int
	var actual int
	var err error

	// Testee
	actorRoot := Actor.New(func(self *ActorDef[interface{}], input interface{}) {
		// Ask cases: ROOT
		switch val := input.(type) {
		case *AskDef[interface{}, int]:
			intVal, _ := Maybe.Just(val.Message).ToInt()

			// NOTE If negative, hanging for testing Ask.timeout
			if intVal < 0 {
				break
			}

			val.Reply(intVal * 10)
			break
		}
	})

	// var timer *time.Timer
	var timeout time.Duration
	timeout = 10 * time.Millisecond

	// Normal cases
	expectedInt = 10
	actual = AskNewGenerics[interface{}, int](1).AskOnce(actorRoot)
	assert.Equal(t, expectedInt, actual)
	// Ask with Timeout
	expectedInt = 20
	actual, _ = AskNewGenerics[interface{}, int](2).AskOnceWithTimeout(actorRoot, timeout)
	assert.Equal(t, expectedInt, actual)
	// Ask channel
	expectedInt = 30
	ch := AskNewGenerics[interface{}, int](3).AskChannel(actorRoot)
	actual = <-ch
	close(ch)
	assert.Equal(t, expectedInt, actual)

	// Timeout cases
	expectedInt = 0
	actual, err = AskNewGenerics[interface{}, int](-1).AskOnceWithTimeout(actorRoot, timeout)
	assert.Equal(t, expectedInt, actual)
	assert.Equal(t, ErrActorAskTimeout, err)
}

func TestActorAskOnceWithTimeoutLong(t *testing.T) {
	actorRoot := Actor.New(func(self *ActorDef[interface{}], input interface{}) {
		switch val := input.(type) {
		case *AskDef[interface{}, int]:
			intVal, _ := Maybe.Just(val.Message).ToInt()
			val.Reply(intVal * 10)
		}
	})

	timeout := 10 * time.Second
	actual, err := AskNewGenerics[interface{}, int](5).AskOnceWithTimeout(actorRoot, timeout)
	assert.Equal(t, 50, actual)
	assert.Nil(t, err)
}

func TestActorChildParent(t *testing.T) {
	actorRoot := Actor.New(func(self *ActorDef[interface{}], input interface{}) {
		switch input.(type) {
		case string:
			if input == "ping" {
				self.Send("pong")
			}
		}
	})

	child := actorRoot.Spawn(func(self *ActorDef[interface{}], input interface{}) {
	})

	parent := child.GetParent()
	assert.Equal(t, actorRoot, parent)

	childID := child.GetID()
	foundChild := actorRoot.GetChild(childID)
	assert.Equal(t, child, foundChild)

	assert.Equal(t, false, actorRoot.IsClosed())
	actorRoot.Close()
	assert.Equal(t, true, actorRoot.IsClosed())
}

func TestActorSendAfterClose(t *testing.T) {
	actor := Actor.New(func(self *ActorDef[interface{}], input interface{}) {
	})

	assert.Equal(t, false, actor.IsClosed())
	actor.Close()
	assert.Equal(t, true, actor.IsClosed())

	actor.Send("test")
}

func TestActorConcurrentSendCloseDoesNotPanic(t *testing.T) {
	actor := Actor.New(func(self *ActorDef[interface{}], input interface{}) {})
	start := make(chan struct{})
	done := make(chan struct{}, 16)
	panicCh := make(chan interface{}, 16)

	for i := 0; i < 16; i++ {
		go func() {
			defer func() {
				if panicVal := recover(); panicVal != nil {
					panicCh <- panicVal
				}
				done <- struct{}{}
			}()
			<-start
			for j := 0; j < 100; j++ {
				actor.Send(j)
			}
		}()
	}

	close(start)
	actor.Close()
	for i := 0; i < 16; i++ {
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("send goroutine did not finish after actor close")
		}
	}
	assert.Empty(t, panicCh)
	assert.True(t, actor.IsClosed())
}

func TestActorGetDefault(t *testing.T) {
	defaultActor := Actor.GetDefault()
	assert.NotNil(t, defaultActor)
}

func TestActorNewByOptions(t *testing.T) {
	ch := make(chan interface{})
	ctx := map[string]interface{}{"key": "value"}
	actor := Actor.NewByOptions(func(self *ActorDef[interface{}], input interface{}) {}, ch, ctx)
	assert.NotNil(t, actor)
}

func TestActorSpawn(t *testing.T) {
	actorRoot := Actor.New(func(self *ActorDef[interface{}], input interface{}) {})

	child := actorRoot.Spawn(func(self *ActorDef[interface{}], input interface{}) {})
	assert.NotNil(t, child)
	assert.Equal(t, actorRoot, child.GetParent())

	child.Close()
	actorRoot.Close()
}

// TestActorSpawnAfterClose tests spawning a child after the parent actor is closed
// This covers actor.go line 77-79
func TestActorSpawnAfterClose(t *testing.T) {
	actorRoot := Actor.New(func(self *ActorDef[interface{}], input interface{}) {})

	// Close the actor first
	actorRoot.Close()

	// Spawn should still work but child won't have parent reference
	child := actorRoot.Spawn(func(self *ActorDef[interface{}], input interface{}) {})
	assert.NotNil(t, child)

	// Parent should be nil since actor is closed
	assert.Nil(t, child.GetParent())
}

func TestAskDef(t *testing.T) {
	ask := Ask.New(42)
	assert.Equal(t, 42, ask.Message)

	ask2 := AskNewGenerics[int, string](100)
	assert.Equal(t, 100, ask2.Message)

	ch := make(chan string)
	ask3 := AskNewByOptionsGenerics(200, ch)
	assert.Equal(t, 200, ask3.Message)
}

func TestAskDef_Reply(t *testing.T) {
	ch := make(chan int)
	ask := AskNewByOptionsGenerics(42, ch)
	go func() {
		ask.Reply(100)
	}()
	result := <-ch
	assert.Equal(t, 100, result)
}

func TestActorGetID(t *testing.T) {
	actor := Actor.New(func(self *ActorDef[interface{}], input interface{}) {})
	id := actor.GetID()
	assert.NotNil(t, id)
}

func TestActorDef_New(t *testing.T) {
	var a *ActorDef[string]
	a = new(ActorDef[string]).New(func(self *ActorDef[string], input string) {})
	assert.NotNil(t, a)
}

func TestActorGenerics(t *testing.T) {
	actor := ActorNewGenerics(func(self *ActorDef[int], input int) {})
	assert.NotNil(t, actor)

	ch := make(chan int)
	ctx := map[string]interface{}{"key": "value"}
	actor2 := ActorNewByOptionsGenerics(func(self *ActorDef[int], input int) {}, ch, ctx)
	assert.NotNil(t, actor2)

	actor.Close()
	actor2.Close()
}

func TestActorGetParent(t *testing.T) {
	parent := Actor.New(func(self *ActorDef[interface{}], input interface{}) {})
	child := parent.Spawn(func(self *ActorDef[interface{}], input interface{}) {})

	assert.NotNil(t, child.GetParent())
	assert.Equal(t, parent.GetID(), child.GetParent().GetID())

	parent.Close()
}

func TestActorGetChild(t *testing.T) {
	parent := Actor.New(func(self *ActorDef[interface{}], input interface{}) {})
	child := parent.Spawn(func(self *ActorDef[interface{}], input interface{}) {})

	childID := child.GetID()
	foundChild := parent.GetChild(childID)
	assert.NotNil(t, foundChild)
	assert.Equal(t, childID, foundChild.GetID())

	parent.Close()
}

func TestActorGetChildNotFound(t *testing.T) {
	parent := Actor.New(func(self *ActorDef[interface{}], input interface{}) {})

	notFoundChild := parent.GetChild(time.Now().Add(-1 * time.Hour))
	assert.Nil(t, notFoundChild)

	parent.Close()
}

func TestActorIsClosed(t *testing.T) {
	actor := Actor.New(func(self *ActorDef[interface{}], input interface{}) {})

	assert.False(t, actor.IsClosed())

	actor.Close()

	assert.True(t, actor.IsClosed())
}

func TestAskNewByOptions(t *testing.T) {
	ch := make(chan int, 1)
	ask := new(AskDef[interface{}, int]).NewByOptions(42, ch)
	assert.NotNil(t, ask)
	assert.Equal(t, 42, ask.Message)
}
