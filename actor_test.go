package fpgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActorCommon(t *testing.T) {
	var expectedInt int
	var actual int
	var resultChannel chan interface{}

	actual = 0
	expectedInt = 1400
	// Testee
	resultChannel = make(chan interface{}, 5)
	cmdSpawn := "spawn"
	cmdShutdown := "shutdown"
	actorRoot := Actor.New(func(self *ActorDef, input interface{}) {
		if input == cmdSpawn {
			self.Spawn(func(self *ActorDef, input interface{}) {
				if input == cmdShutdown {
					self.Close()
					return
				}

				val, _ := Maybe.Just(input).ToInt()
				resultChannel <- val * 10
			})
			return
		}
		if input == cmdShutdown {
			for _, child := range self.children {
				child.Send(cmdShutdown)
			}
			self.Close()

			close(resultChannel)
			return
		}

		intVal, _ := Maybe.Just(input).ToInt()
		if intVal > 0 {
			for _, child := range self.children {
				child.Send(intVal)
			}
		}
	})
	actorRoot.Send(cmdSpawn)
	actorRoot.Send(10)
	actorRoot.Send(cmdSpawn)
	actorRoot.Send(20)
	actorRoot.Send(cmdSpawn)
	actorRoot.Send(30)

	actorRoot.Send(cmdShutdown)

	for val := range resultChannel {
		intVal, _ := Maybe.Just(val).ToInt()
		actual += intVal
	}

	assert.Equal(t, expectedInt, actual)
}
