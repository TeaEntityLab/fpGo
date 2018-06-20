package fpGo

import "sync"

// PublisherDef Publisher inspired by Rx/NotificationCenter/PubSub
type PublisherDef struct {
	subscribers []*Subscription
	subscribeM  sync.Mutex
	subOn       *HandlerDef

	origin *PublisherDef
}

// New New a Publisher
func (publisherSelf *PublisherDef) New() *PublisherDef {
	p := &PublisherDef{}

	return p
}

// Map Map the Publisher in order to make a broadcasting chain
func (publisherSelf *PublisherDef) Map(fn func(interface{}) interface{}) *PublisherDef {
	next := publisherSelf.New()
	next.origin = publisherSelf
	publisherSelf.Subscribe(Subscription{
		OnNext: func(in interface{}) {
			next.Publish(fn(in))
		},
	})

	return next
}

// Subscribe Subscribe the Publisher by Subscription
func (publisherSelf *PublisherDef) Subscribe(sub Subscription) *Subscription {
	s := &sub

	publisherSelf.doSubscribeSafe(func() {
		publisherSelf.subscribers = append(publisherSelf.subscribers, s)
	})
	return s
}

// SubscribeOn Subscribe the Publisher on the specific Handler
func (publisherSelf *PublisherDef) SubscribeOn(h *HandlerDef) *PublisherDef {
	publisherSelf.subOn = h
	return publisherSelf
}

// Unsubscribe Unsubscribe the publisher by the Subscription
func (publisherSelf *PublisherDef) Unsubscribe(s *Subscription) {
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
func (publisherSelf *PublisherDef) Publish(result interface{}) {
	var subscribers []*Subscription
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
func (publisherSelf *PublisherDef) doSubscribeSafe(fn func()) {
	publisherSelf.subscribeM.Lock()
	fn()
	publisherSelf.subscribeM.Unlock()
}

// Publisher Publisher utils instance
var Publisher PublisherDef
