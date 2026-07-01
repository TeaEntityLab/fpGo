package fpgo

import (
	"sync"
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
		actorRoot := Actor.New(func(self *ActorDef, input interface{}) {
			// SPAWN: for ROOT
			if input == cmdSpawn {
				self.Spawn(func(self *ActorDef, input interface{}) {
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
	var result interface{}

	// Testee
	actorRoot := Actor.New(func(self *ActorDef, input interface{}) {
		// Ask cases: ROOT
		switch val := input.(type) {
		case *AskDef:
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
	result = Ask.New(1).AskOnce(actorRoot)
	actual, _ = Maybe.Just(result).ToInt()
	assert.Equal(t, expectedInt, actual)
	// Ask with Timeout
	expectedInt = 20
	result, _ = Ask.New(2).AskOnceWithTimeout(actorRoot, timeout)
	actual, _ = Maybe.Just(result).ToInt()
	assert.Equal(t, expectedInt, actual)
	// Ask channel
	expectedInt = 30
	ch := AskNewGenerics(3).AskChannel(actorRoot)
	actual, _ = Maybe.Just(<-ch).ToInt()
	close(ch)
	assert.Equal(t, expectedInt, actual)

	// Timeout cases
	expectedInt = 0
	result, err = Ask.New(-1).AskOnceWithTimeout(actorRoot, timeout)
	actual, _ = Maybe.Just(result).ToInt()
	assert.Equal(t, expectedInt, actual)
	assert.Equal(t, ErrActorAskTimeout, err)
}

func TestActorChildParent(t *testing.T) {
	actorRoot := Actor.New(func(self *ActorDef, input interface{}) {
		switch input.(type) {
		case string:
			if input == "ping" {
				self.Send("pong")
			}
		}
	})

	child := actorRoot.Spawn(func(self *ActorDef, input interface{}) {
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
	actor := Actor.New(func(self *ActorDef, input interface{}) {
	})

	assert.Equal(t, false, actor.IsClosed())
	actor.Close()
	assert.Equal(t, true, actor.IsClosed())

	actor.Send("test")
}

func TestActorConcurrentSendCloseDoesNotPanic(t *testing.T) {
	actor := Actor.New(func(self *ActorDef, input interface{}) {})
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
	actor := Actor.NewByOptions(func(self *ActorDef, input interface{}) {}, ch, ctx)
	assert.NotNil(t, actor)
	actor.Close()
}

func TestActorSpawn(t *testing.T) {
	actorRoot := Actor.New(func(self *ActorDef, input interface{}) {})

	child := actorRoot.Spawn(func(self *ActorDef, input interface{}) {})
	assert.NotNil(t, child)
	assert.Equal(t, actorRoot, child.GetParent())

	child.Close()
	actorRoot.Close()
}

func TestActorSpawnAfterClose(t *testing.T) {
	actorRoot := Actor.New(func(self *ActorDef, input interface{}) {})

	actorRoot.Close()

	child := actorRoot.Spawn(func(self *ActorDef, input interface{}) {})
	assert.NotNil(t, child)
	assert.Nil(t, child.GetParent())
}

func TestAskDef(t *testing.T) {
	ask := Ask.New(42)
	assert.Equal(t, 42, ask.Message)

	ask2 := AskNewGenerics(100)
	assert.Equal(t, 100, ask2.Message)

	ch := make(chan interface{})
	ask3 := AskNewByOptionsGenerics(200, ch)
	assert.Equal(t, 200, ask3.Message)
}

func TestAskDef_Reply(t *testing.T) {
	ch := make(chan interface{})
	ask := AskNewByOptionsGenerics(42, ch)
	go func() {
		ask.Reply(100)
	}()
	result := <-ch
	assert.Equal(t, 100, result)
}

func TestActorGetID(t *testing.T) {
	actor := Actor.New(func(self *ActorDef, input interface{}) {})
	id := actor.GetID()
	assert.NotNil(t, id)
	actor.Close()
}

func TestActorGetParent(t *testing.T) {
	parent := Actor.New(func(self *ActorDef, input interface{}) {})
	child := parent.Spawn(func(self *ActorDef, input interface{}) {})

	assert.NotNil(t, child.GetParent())
	assert.Equal(t, parent.GetID(), child.GetParent().GetID())

	parent.Close()
}

func TestActorGetChild(t *testing.T) {
	parent := Actor.New(func(self *ActorDef, input interface{}) {})
	child := parent.Spawn(func(self *ActorDef, input interface{}) {})

	childID := child.GetID()
	foundChild := parent.GetChild(childID)
	assert.NotNil(t, foundChild)
	assert.Equal(t, childID, foundChild.GetID())

	parent.Close()
}

func TestActorGetChildNotFound(t *testing.T) {
	parent := Actor.New(func(self *ActorDef, input interface{}) {})

	notFoundChild := parent.GetChild(time.Now().Add(-1 * time.Hour))
	assert.Nil(t, notFoundChild)

	parent.Close()
}

func TestActorCloseMultiple(t *testing.T) {
	actor := Actor.New(func(self *ActorDef, input interface{}) {})
	assert.NotPanics(t, func() {
		actor.Close()
		actor.Close()
	})
	assert.True(t, actor.IsClosed())
}

func TestActorIsClosed(t *testing.T) {
	actor := Actor.New(func(self *ActorDef, input interface{}) {})

	assert.False(t, actor.IsClosed())

	actor.Close()

	assert.True(t, actor.IsClosed())
}

func TestAskNewByOptions(t *testing.T) {
	ch := make(chan interface{}, 1)
	ask := new(AskDef).NewByOptions(42, ch)
	assert.NotNil(t, ask)
	assert.Equal(t, 42, ask.Message)
}

func TestActorAskReplyAfterTimeoutDoesNotPanicOrBlock(t *testing.T) {
	release := make(chan struct{})
	processed := make(chan int, 1)

	actorRoot := Actor.New(func(self *ActorDef, input interface{}) {
		switch val := input.(type) {
		case *AskDef:
			<-release
			val.Reply(999)
			processed <- 1
		}
	})
	defer actorRoot.Close()

	result, err := Ask.New(1).AskOnceWithTimeout(actorRoot, 5*time.Millisecond)
	assert.Equal(t, ErrActorAskTimeout, err)
	actual, _ := Maybe.Just(result).ToInt()
	assert.Equal(t, 0, actual)

	assert.NotPanics(t, func() {
		close(release)
		select {
		case <-processed:
		case <-time.After(time.Second):
			t.Fatal("actor goroutine blocked on Reply after asker timed out")
		}
	})
}

func TestActorAskTimeoutThenReplySucceedsForNextAsk(t *testing.T) {
	release := make(chan struct{})
	actorRoot := Actor.New(func(self *ActorDef, input interface{}) {
		switch val := input.(type) {
		case *AskDef:
			intVal, _ := Maybe.Just(val.Message).ToInt()
			if intVal < 0 {
				<-release
			}
			val.Reply(intVal * 10)
		}
	})
	defer actorRoot.Close()

	_, err := Ask.New(-1).AskOnceWithTimeout(actorRoot, 5*time.Millisecond)
	assert.Equal(t, ErrActorAskTimeout, err)
	close(release)

	result, err := Ask.New(7).AskOnceWithTimeout(actorRoot, time.Second)
	assert.NoError(t, err)
	actual, _ := Maybe.Just(result).ToInt()
	assert.Equal(t, 70, actual)
}

func TestActorConcurrentSpawnGetChild(t *testing.T) {
	parent := Actor.New(func(self *ActorDef, input interface{}) {})

	const workers = 8
	const iterations = 50
	start := make(chan struct{})
	var wg sync.WaitGroup
	var childrenMu sync.Mutex
	children := make([]*ActorDef, 0, workers*iterations)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < iterations; j++ {
				child := parent.Spawn(func(self *ActorDef, input interface{}) {})
				childrenMu.Lock()
				children = append(children, child)
				childrenMu.Unlock()
			}
		}()
	}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < iterations; j++ {
				childrenMu.Lock()
				snapshot := append([]*ActorDef(nil), children...)
				childrenMu.Unlock()
				for _, child := range snapshot {
					if child == nil {
						continue
					}
					found := parent.GetChild(child.GetID())
					if found != nil {
						assert.Equal(t, child.GetID(), found.GetID())
					}
				}
			}
		}()
	}

	close(start)
	wg.Wait()

	childrenMu.Lock()
	for _, child := range children {
		child.Close()
	}
	childrenMu.Unlock()
	parent.Close()
	time.Sleep(100 * time.Millisecond)
}
