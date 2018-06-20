package fpGo

// MonadIODef MonadIO inspired by Rx/Observable
type MonadIODef struct {
	effect func() interface{}

	obOn  *HandlerDef
	subOn *HandlerDef
}

// Subscription the delegation/callback of MonadIO/Publisher
type Subscription struct {
	OnNext func(interface{})
}

// Just New MonadIO by a given value
func (monadIOSelf MonadIODef) Just(in interface{}) *MonadIODef {
	return &MonadIODef{effect: func() interface{} {
		return in
	}}
}

// New New MonadIO by effect function
func (monadIOSelf *MonadIODef) New(effect func() interface{}) *MonadIODef {
	return &MonadIODef{effect: effect}
}

// FlatMap FlatMap the MonadIO by function
func (monadIOSelf *MonadIODef) FlatMap(fn func(interface{}) *MonadIODef) *MonadIODef {

	return &MonadIODef{effect: func() interface{} {
		next := fn(monadIOSelf.doEffect())
		return next.doEffect()
	}}

}

// Subscribe Subscribe the MonadIO by Subscription
func (monadIOSelf *MonadIODef) Subscribe(s Subscription) *Subscription {
	obOn := monadIOSelf.obOn
	subOn := monadIOSelf.subOn
	return monadIOSelf.doSubscribe(&s, obOn, subOn)
}

// SubscribeOn Subscribe the MonadIO on the specific Handler
func (monadIOSelf *MonadIODef) SubscribeOn(h *HandlerDef) *MonadIODef {
	monadIOSelf.subOn = h
	return monadIOSelf
}

// ObserveOn Observe the MonadIO on the specific Handler
func (monadIOSelf *MonadIODef) ObserveOn(h *HandlerDef) *MonadIODef {
	monadIOSelf.obOn = h
	return monadIOSelf
}
func (monadIOSelf *MonadIODef) doSubscribe(s *Subscription, obOn *HandlerDef, subOn *HandlerDef) *Subscription {

	if s.OnNext != nil {
		var result interface{}

		doSub := func() {
			s.OnNext(result)
		}
		doOb := func() {
			result = monadIOSelf.doEffect()

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
func (monadIOSelf *MonadIODef) doEffect() interface{} {
	return monadIOSelf.effect()
}

// MonadIO MonadIO utils instance
var MonadIO MonadIODef
