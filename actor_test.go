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
	result, _ = Ask.New(1).AskOnce(actorRoot, nil)
	actual, _ = Maybe.Just(result).ToInt()
	assert.Equal(t, expectedInt, actual)
	// Ask with Timeout
	expectedInt = 20
	result, _ = Ask.New(2).AskOnce(actorRoot, &timeout)
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
	result, err = Ask.New(-1).AskOnce(actorRoot, &timeout)
	actual, _ = Maybe.Just(result).ToInt()
	assert.Equal(t, expectedInt, actual)
	assert.Equal(t, ErrActorAskTimeout, err)
}
