package fpgo

import (
	"fmt"
	"sync"
	"time"
)

var ErrActorAskTimeout = fmt.Errorf("ErrActorAskTimeout")

// ActorHandle A target could send messages
type ActorHandle interface {
	Send(message interface{})
}

// ActorDef Actor model inspired by Erlang/Akka
type ActorDef struct {
	id       time.Time
	isClosed AtomBool
	ch       chan interface{}
	done     chan struct{}
	effect   func(*ActorDef, interface{})

	context map[string]interface{}

	childrenLock sync.RWMutex
	children     map[time.Time]*ActorDef
	parent       *ActorDef
}

var defaultActor *ActorDef

// GetDefault Get Default Actor
func (actorSelf *ActorDef) GetDefault() *ActorDef {
	return defaultActor
}

// New New Actor instance
func (actorSelf *ActorDef) New(effect func(*ActorDef, interface{})) *ActorDef {
	return actorSelf.NewByOptions(effect, make(chan interface{}), map[string]interface{}{})
}

// NewByOptions New Actor by its options
func (actorSelf *ActorDef) NewByOptions(effect func(*ActorDef, interface{}), ioCh chan interface{}, context map[string]interface{}) *ActorDef {
	newOne := ActorDef{
		id:       time.Now(),
		ch:       ioCh,
		done:     make(chan struct{}),
		effect:   effect,
		context:  context,
		children: map[time.Time]*ActorDef{},
	}
	go newOne.run()

	return &newOne
}

// Send Send a message to the Actor
func (actorSelf *ActorDef) Send(message interface{}) {
	if actorSelf.isClosed.Get() {
		return
	}

	select {
	case <-actorSelf.done:
	case actorSelf.ch <- message:
	}
}

// Spawn Spawn a new Actor with parent(this actor)
func (actorSelf *ActorDef) Spawn(effect func(*ActorDef, interface{})) *ActorDef {
	newOne := Actor.New(effect)
	if actorSelf.isClosed.Get() {
		return newOne
	}

	newOne.parent = actorSelf
	actorSelf.childrenLock.Lock()
	actorSelf.children[newOne.id] = newOne
	actorSelf.childrenLock.Unlock()

	return newOne
}

// GetChild Get a child Actor by ID
func (actorSelf *ActorDef) GetChild(id time.Time) *ActorDef {
	actorSelf.childrenLock.RLock()
	defer actorSelf.childrenLock.RUnlock()
	return actorSelf.children[id]
}

// GetParent Get its parent Actor
func (actorSelf *ActorDef) GetParent() *ActorDef {
	return actorSelf.parent
}

// GetID Get ID time.Time
func (actorSelf *ActorDef) GetID() time.Time {
	return actorSelf.id
}

// Close Close the Actor
func (actorSelf *ActorDef) Close() {
	if !actorSelf.isClosed.CompareAndSwap(false, true) {
		return
	}

	close(actorSelf.done)
}

// IsClosed Check is Closed
func (actorSelf *ActorDef) IsClosed() bool {
	return actorSelf.isClosed.Get()
}

func (actorSelf *ActorDef) run() {
	for {
		select {
		case <-actorSelf.done:
			return
		case message, ok := <-actorSelf.ch:
			if !ok {
				return
			}
			actorSelf.effect(actorSelf, message)
		}
	}
}

// Actor Actor utils instance
var Actor ActorDef

// AskDef Ask inspired by Erlang/Akka
type AskDef struct {
	id time.Time
	ch chan interface{}

	Message interface{}
}

// New New Ask instance
func (askSelf *AskDef) New(message interface{}) *AskDef {
	return AskNewGenerics(message)
}

// NewByOptions New Ask by its options
func (askSelf *AskDef) NewByOptions(message interface{}, ioCh chan interface{}) *AskDef {
	return AskNewByOptionsGenerics(message, ioCh)
}

// AskNewGenerics New Ask instance
func AskNewGenerics(message interface{}) *AskDef {
	// Buffered (cap 1) so Reply never blocks and never needs a live receiver:
	// after AskOnceWithTimeout times out there is no reader, and an unbuffered
	// channel would either deadlock the replying actor goroutine or (once the
	// ask side closed it) panic with "send on closed channel".
	return AskNewByOptionsGenerics(message, make(chan interface{}, 1))
}

// AskNewByOptionsGenerics New Ask by its options
func AskNewByOptionsGenerics(message interface{}, ioCh chan interface{}) *AskDef {
	newOne := AskDef{
		id: time.Now(),
		ch: ioCh,

		Message: message,
	}

	return &newOne
}

// AskOnce Sender Ask
func (askSelf *AskDef) AskOnce(target ActorHandle) interface{} {
	ch := askSelf.AskChannel(target)

	return <-ch
}

// AskOnce Sender Ask with timeout
func (askSelf *AskDef) AskOnceWithTimeout(target ActorHandle, timeout time.Duration) (interface{}, error) {
	ch := askSelf.AskChannel(target)
	var result interface{}
	select {
	case result = <-ch:
	case <-time.After(timeout):
		return result, ErrActorAskTimeout
	}

	return result, nil
}

// AskChannel Sender Ask
func (askSelf *AskDef) AskChannel(target ActorHandle) chan interface{} {
	target.Send(askSelf)

	return askSelf.ch
}

// Reply Receiver Reply
func (askSelf *AskDef) Reply(response interface{}) {
	askSelf.ch <- response
}

// Ask Ask utils instance
var Ask AskDef

func init() {
	Actor.isClosed.Set(true)
	defaultActor = &Actor
}
