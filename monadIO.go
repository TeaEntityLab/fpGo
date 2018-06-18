package fpGo

type MonadIODef struct {
	effect func() interface{}

	obOn  *HandlerDef
	subOn *HandlerDef
}
type Subscription struct {
	OnNext func(interface{})
}

func (self MonadIODef) Just(in interface{}) *MonadIODef {
	return &MonadIODef{effect: func() interface{} {
		return in
	}}
}
func (self *MonadIODef) New(effect func() interface{}) *MonadIODef {
	return &MonadIODef{effect: effect}
}

func (self *MonadIODef) FlatMap(fn func(interface{}) *MonadIODef) *MonadIODef {

	return &MonadIODef{effect: func() interface{} {
		next := fn(self.doEffect())
		return next.doEffect()
	}}

}
func (self *MonadIODef) Subscribe(s Subscription) *Subscription {
	obOn := self.obOn
	subOn := self.subOn
	return self.doSubscribe(&s, obOn, subOn)
}
func (self *MonadIODef) SubscribeOn(h *HandlerDef) *MonadIODef {
	self.subOn = h
	return self
}
func (self *MonadIODef) ObserveOn(h *HandlerDef) *MonadIODef {
	self.obOn = h
	return self
}
func (self *MonadIODef) doSubscribe(s *Subscription, obOn *HandlerDef, subOn *HandlerDef) *Subscription {

	if s.OnNext != nil {
		var result interface{}

		doSub := func() {
			s.OnNext(result)
		}
		doOb := func() {
			result = self.doEffect()

			if subOn != nil {
				subOn.Post(doSub)
			} else {
				doSub()
			}
		}
		if obOn != nil {
			obOn.Post(doOb)
		} else {
			doOb()
		}
	}

	return s
}
func (self *MonadIODef) doEffect() interface{} {
	return self.effect()
}

var MonadIO MonadIODef
