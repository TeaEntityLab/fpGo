package fpgo

import "sync"

// PublisherDef Publisher inspired by Rx/NotificationCenter/PubSub
type PublisherDef[T any] struct {
	subscribers []*Subscription[T]
	subscribeM  sync.Mutex
	subOn       *HandlerDef

	origin *PublisherDef[T]
}

// New New a Publisher
func (publisherSelf *PublisherDef[T]) New() *PublisherDef[interface{}] {
	return PublisherNewGenerics[interface{}]()
}

// PublisherNewGenerics New a Publisher
func PublisherNewGenerics[T any]() *PublisherDef[T] {
	p := &PublisherDef[T]{}

	return p
}

// Map Map the Publisher in order to make a broadcasting chain
func (publisherSelf *PublisherDef[T]) Map(fn func(T) T) *PublisherDef[T] {
	next := PublisherNewGenerics[T]()
	next.origin = publisherSelf
	publisherSelf.Subscribe(Subscription[T]{
		OnNext: func(in T) {
			next.Publish(fn(in))
		},
	})

	return next
}

// Subscribe Subscribe the Publisher by Subscription[T]
func (publisherSelf *PublisherDef[T]) Subscribe(sub Subscription[T]) *Subscription[T] {
	s := &sub

	publisherSelf.doSubscribeSafe(func() {
		publisherSelf.subscribers = append(publisherSelf.subscribers, s)
	})
	return s
}

// SubscribeOn Subscribe the Publisher on the specific Handler
func (publisherSelf *PublisherDef[T]) SubscribeOn(h *HandlerDef) *PublisherDef[T] {
	publisherSelf.subOn = h
	return publisherSelf
}

// Unsubscribe Unsubscribe the publisher by the Subscription[T]
func (publisherSelf *PublisherDef[T]) Unsubscribe(s *Subscription[T]) {
	isAnyMatching := false

	publisherSelf.doSubscribeSafe(func() {
		subscribers := publisherSelf.subscribers
		for i, v := range subscribers {
			if v == s {
				isAnyMatching = true
				subscribers = append(subscribers[:i], subscribers[i+1:]...)
				publisherSelf.subscribers = subscribers
				break
			}
		}
	})

	// Delete subscriptions recursively
	if isAnyMatching {
		publisherSelf.Unsubscribe(s)
	}
}

// Publish Publish a value to its subscribers or next chains
func (publisherSelf *PublisherDef[T]) Publish(result T) {
	var subscribers []*Subscription[T]
	publisherSelf.doSubscribeSafe(func() {
		subscribers = publisherSelf.subscribers
	})

	for _, s := range subscribers {
		if s.OnNext != nil {

			doSub := func() {
				s.OnNext(result)
			}
			if publisherSelf.subOn != nil {
				publisherSelf.subOn.Post(doSub)
			} else {
				doSub()
			}
		}
	}
}

func (publisherSelf *PublisherDef[T]) doSubscribeSafe(fn func()) {
	publisherSelf.subscribeM.Lock()
	fn()
	publisherSelf.subscribeM.Unlock()
}

// Publisher Publisher utils instance
var Publisher PublisherDef[interface{}]
