package fpGo

type MonadIODef struct {
	effect  func() *interface{}
	origin  *MonadIODef
	subList []*Subscription
}
type Subscription struct {
	OnNext func(*interface{})
}

func (self MonadIODef) JustVal(in interface{}) *MonadIODef {
	return self.Just(&in)
}
func (self MonadIODef) Just(in *interface{}) *MonadIODef {
	return &MonadIODef{effect: func() *interface{} {
		return in
	}}
}

func (self *MonadIODef) FlatMap(fn func(*interface{}) MonadIODef) *MonadIODef {

	return &MonadIODef{effect: func() *interface{} {
		next := fn(self.doEffect())
		next.setOrigin(self.GetOrigin())
		return next.doEffect()
	}, origin: self.GetOrigin()}

}
func (self *MonadIODef) Subscribe(s Subscription) *Subscription {
	if s.OnNext != nil {
		origin := self.GetOrigin()
		origin.AddSubscription(&s)
		origin.doEffect()
	}

	return &s
}
func (self *MonadIODef) AddSubscription(s *Subscription) {
	self.subList = append(self.subList, s)
}
func (self MonadIODef) GetSubscriptionList() []*Subscription {
	return self.subList
}
func (self MonadIODef) setOrigin(m *MonadIODef) {
	self.origin = m
}
func (self *MonadIODef) GetOrigin() *MonadIODef {
	if self.origin == nil {
		return self
	}

	return self.origin
}
func (self *MonadIODef) doEffect() *interface{} {
	v := self.effect()
	for _, s := range self.GetSubscriptionList() {
		s.OnNext(v)
	}
	return v
}

var MonadIO MonadIODef
