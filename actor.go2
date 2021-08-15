package fpgo

import (
	"time"
)

// ActorDef Actor model inspired by Erlang/Akka
type ActorDef struct {
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
	newOne := ActorDef{ch: ioCh, effect: effect, context: context, children: map[time.Time]*ActorDef{}}
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
	actorSelf.children[time.Now()] = newOne

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

func init() {
	Actor = *Actor.New(func(_ *ActorDef, _ interface{}) {})
	Actor.Close()
	defaultActor = &Actor
}
