package fpgo

import (
	"fmt"
	"time"
)

var ErrActorAskTimeout = fmt.Errorf("ErrActorAskTimeout")

// ActorHandle A target could send messages
type ActorHandle[T any] interface {
	Send(message T)
}

// ActorDef[T] Actor model inspired by Erlang/Akka
type ActorDef[T any] struct {
	id       time.Time
	isClosed bool
	ch       chan T
	effect   func(*ActorDef[T], T)

	context map[string]interface{}

	children map[time.Time]*ActorDef[T]
	parent   *ActorDef[T]
}

var defaultActor *ActorDef[interface{}]

// GetDefault Get Default Actor
func (actorSelf *ActorDef[T]) GetDefault() *ActorDef[interface{}] {
	return defaultActor
}

// New New Actor instance
func (actorSelf *ActorDef[T]) New(effect func(*ActorDef[T], T)) *ActorDef[T] {
	return ActorNewGenerics(effect)
}

// NewByOptions New Actor by its options
func (actorSelf *ActorDef[T]) NewByOptions(effect func(*ActorDef[T], T), ioCh chan T, context map[string]interface{}) *ActorDef[T] {
	return ActorNewByOptionsGenerics(effect, ioCh, context)
}

// ActorNewGenerics New Actor instance
func ActorNewGenerics[T any](effect func(*ActorDef[T], T)) *ActorDef[T] {
	return ActorNewByOptionsGenerics(effect, make(chan T), map[string]interface{}{})
}

// ActorNewByOptionsGenerics New Actor by its options
func ActorNewByOptionsGenerics[T any](effect func(*ActorDef[T], T), ioCh chan T, context map[string]interface{}) *ActorDef[T] {
	newOne := ActorDef[T]{
		id:       time.Now(),
		ch:       ioCh,
		effect:   effect,
		context:  context,
		children: map[time.Time]*ActorDef[T]{},
	}

	go newOne.run()

	return &newOne
}

// Send Send a message to the Actor
func (actorSelf *ActorDef[T]) Send(message T) {
	if actorSelf.isClosed {
		return
	}

	actorSelf.ch <- message
}

// Spawn Spawn a new Actor with parent(this actor)
func (actorSelf *ActorDef[T]) Spawn(effect func(*ActorDef[T], T)) *ActorDef[T] {
	newOne := actorSelf.New(effect)
	if actorSelf.isClosed {
		return newOne
	}

	newOne.parent = actorSelf
	actorSelf.children[newOne.id] = newOne

	return newOne
}

// GetChild Get a child Actor by ID
func (actorSelf *ActorDef[T]) GetChild(id time.Time) *ActorDef[T] {
	return actorSelf.children[id]
}

// GetParent Get its parent Actor
func (actorSelf *ActorDef[T]) GetParent() *ActorDef[T] {
	return actorSelf.parent
}

// GetID Get ID time.Time
func (actorSelf *ActorDef[T]) GetID() time.Time {
	return actorSelf.id
}

// Close Close the Actor
func (actorSelf *ActorDef[T]) Close() {
	actorSelf.isClosed = true

	close(actorSelf.ch)
}

// IsClosed Check is Closed
func (actorSelf *ActorDef[T]) IsClosed() bool {
	return actorSelf.isClosed
}

func (actorSelf *ActorDef[T]) run() {
	for message := range actorSelf.ch {
		actorSelf.effect(actorSelf, message)
	}
}

// Actor Actor utils instance
var Actor ActorDef[interface{}]

// AskDef[T, R] Ask inspired by Erlang/Akka
type AskDef[T any, R any] struct {
	id time.Time
	ch chan R

	Message T
}

// New New Ask instance
func (askSelf *AskDef[T, R]) New(message T) *AskDef[T, R] {
	return AskNewGenerics[T, R](message)
}

// NewByOptions New Ask by its options
func (askSelf *AskDef[T, R]) NewByOptions(message T, ioCh chan R) *AskDef[T, R] {
	return AskNewByOptionsGenerics[T, R](message, ioCh)
}

// AskNewGenerics New Ask instance
func AskNewGenerics[T any, R any](message T) *AskDef[T, R] {
	return AskNewByOptionsGenerics[T, R](message, make(chan R))
}

// AskNewByOptionsGenerics New Ask by its options
func AskNewByOptionsGenerics[T any, R any](message T, ioCh chan R) *AskDef[T, R] {
	newOne := AskDef[T, R]{
		id: time.Now(),
		ch: ioCh,

		Message: message,
	}

	return &newOne
}

// AskOnce Sender Ask
func (askSelf *AskDef[T, R]) AskOnce(target ActorHandle[interface{}]) (R, error) {
	ch := askSelf.AskChannel(target)
	defer close(ch)
	// var err error

	return <-ch, nil
}

// AskOnceWithTimeout Sender Ask with timeout
func (askSelf *AskDef[T, R]) AskOnceWithTimeout(target ActorHandle[interface{}], timeout time.Duration) (R, error) {
	ch := askSelf.AskChannel(target)
	defer close(ch)
	var result R
	select {
	case result = <-ch:
	case <-time.After(timeout):
		return result, ErrActorAskTimeout
	}

	return result, nil
}

// AskChannel Sender Ask
func (askSelf *AskDef[T, R]) AskChannel(target ActorHandle[interface{}]) chan R {
	target.Send(askSelf)

	return askSelf.ch
}

// Reply Receiver Reply
func (askSelf *AskDef[T, R]) Reply(response R) {
	askSelf.ch <- response
}

// Ask Ask utils instance
var Ask AskDef[interface{}, interface{}]

func init() {
	// Ask = *Ask.New(0, nil)
	// Actor = *Actor.New(func(_ *ActorDef[interface{}], _ interface{}) {})
	// Actor.Close()
	Actor.isClosed = true
	defaultActor = &Actor
}
