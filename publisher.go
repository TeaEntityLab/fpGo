package fpGo

import "sync"

type PublisherDef struct {
	subscribers []*Subscription
	subscribeM  sync.Mutex
	subOn       *HandlerDef

	origin *PublisherDef
}

func (self *PublisherDef) New() *PublisherDef {
	p := &PublisherDef{}

	return p
}

func (self *PublisherDef) Map(fn func(*interface{}) *interface{}) *PublisherDef {
	next := self.New()
	next.origin = self
	self.Subscribe(Subscription{
		OnNext: func(in *interface{}) {
			next.Publish(fn(in))
		},
	})

	return next
}
func (self *PublisherDef) Subscribe(sub Subscription) *Subscription {
	s := &sub

	self.doSubscribeSafe(func() {
		self.subscribers = append(self.subscribers, s)
	})
	return s
}
func (self *PublisherDef) SubscribeOn(h *HandlerDef) *PublisherDef {
	self.subOn = h
	return self
}
func (self *PublisherDef) Unsubscribe(s *Subscription) {
	isAnyMatching := false

	self.doSubscribeSafe(func() {
		subscribers := self.subscribers
		for i, v := range subscribers {
			if v == s {
				isAnyMatching = true
				subscribers = append(subscribers[:i], subscribers[i+1:]...)
				self.subscribers = subscribers
				break
			}
		}
	})

	// Delete subscriptions recursively
	if isAnyMatching {
		self.Unsubscribe(s)
	}
}
func (self *PublisherDef) Publish(result *interface{}) {
	var subscribers []*Subscription
	self.doSubscribeSafe(func() {
		subscribers = self.subscribers
	})

	for _, s := range subscribers {
		if s.OnNext != nil {

			doSub := func() {
				s.OnNext(result)
			}
			if self.subOn != nil {
				self.subOn.Post(doSub)
			} else {
				doSub()
			}
		}
	}
}
func (self *PublisherDef) doSubscribeSafe(fn func()) {
	self.subscribeM.Lock()
	fn()
	self.subscribeM.Unlock()
}

var Publisher PublisherDef
