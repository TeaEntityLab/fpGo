package fpgo

import (
	"fmt"
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
	isClosed bool
	ch       *chan interface{}
	effect   func(*ActorDef, interface{})

	context map[string]interface{}

	children map[time.Time]*ActorDef
	parent   *ActorDef
}

var defaultActor *ActorDef

// GetDefault Get Default Actor
func (actorSelf *ActorDef) GetDefault() *ActorDef {
	return defaultActor
}

// New New Actor instance
func (actorSelf *ActorDef) New(effect func(*ActorDef, interface{})) *ActorDef {
	ch := make(chan interface{})
	return actorSelf.NewByOptions(effect, &ch, map[string]interface{}{})
}

// NewByOptions New Actor by its options
func (actorSelf *ActorDef) NewByOptions(effect func(*ActorDef, interface{}), ioCh *chan interface{}, context map[string]interface{}) *ActorDef {
	newOne := ActorDef{
		id:       time.Now(),
		ch:       ioCh,
		effect:   effect,
		context:  context,
		children: map[time.Time]*ActorDef{},
	}
	go newOne.run()

	return &newOne
}

// Send Send a message to the Actor
func (actorSelf *ActorDef) Send(message interface{}) {
	if actorSelf.isClosed {
		return
	}

	*(actorSelf.ch) <- message
}

// Spawn Spawn a new Actor with parent(this actor)
func (actorSelf *ActorDef) Spawn(effect func(*ActorDef, interface{})) *ActorDef {
	newOne := Actor.New(effect)
	if actorSelf.isClosed {
		return newOne
	}

	newOne.parent = actorSelf
	actorSelf.children[newOne.id] = newOne

	return newOne
}

// GetChild Get a child Actor by ID
func (actorSelf *ActorDef) GetChild(id time.Time) *ActorDef {
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
	actorSelf.isClosed = true

	close(*actorSelf.ch)
}

// IsClosed Check is Closed
func (actorSelf *ActorDef) IsClosed() bool {
	return actorSelf.isClosed
}

func (actorSelf *ActorDef) run() {
	for message := range *actorSelf.ch {
		actorSelf.effect(actorSelf, message)
	}
}

// Actor Actor utils instance
var Actor ActorDef

// AskDef Ask inspired by Erlang/Akka
type AskDef struct {
	id time.Time
	ch *chan interface{}

	Message interface{}
}

// New New Ask instance
func (askSelf *AskDef) New(message interface{}) *AskDef {
	return AskNewGenerics(message)
}

// NewByOptions New Ask by its options
func (askSelf *AskDef) NewByOptions(message interface{}, ioCh *chan interface{}) *AskDef {
	return AskNewByOptionsGenerics(message, ioCh)
}

// AskNewGenerics New Ask instance
func AskNewGenerics(message interface{}) *AskDef {
	ch := make(chan interface{})
	return AskNewByOptionsGenerics(message, &ch)
}

// AskNewByOptionsGenerics New Ask by its options
func AskNewByOptionsGenerics(message interface{}, ioCh *chan interface{}) *AskDef {
	newOne := AskDef{
		id: time.Now(),
		ch: ioCh,

		Message: message,
	}

	return &newOne
}

// AskOnce Sender Ask
func (askSelf *AskDef) AskOnce(target ActorHandle, timeout *time.Duration) (interface{}, error) {
	ch := askSelf.AskChannel(target)
	defer close(*ch)
	var result interface{}
	// var err error
	if timeout == nil {
		result = <-*ch
	} else {
		select {
		case result = <-*ch:
		case <-time.After(*timeout):
			return result, ErrActorAskTimeout
		}
	}

	return result, nil
}

// AskChannel Sender Ask
func (askSelf *AskDef) AskChannel(target ActorHandle) *chan interface{} {
	target.Send(askSelf)

	return askSelf.ch
}

// Reply Receiver Reply
func (askSelf *AskDef) Reply(response interface{}) {
	*askSelf.ch <- response
}

// Ask Ask utils instance
var Ask AskDef

func init() {
	// Ask = *Ask.New(0, nil)
	// Actor = *Actor.New(func(_ *ActorDef, _ interface{}) {})
	// Actor.Close()
	Actor.isClosed = true
	defaultActor = &Actor
}
