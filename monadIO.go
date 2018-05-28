package fpGo

type MonadIODef struct {
	effect func() *interface{}

	obOn  *HandlerDef
	subOn *HandlerDef
}
type Subscription struct {
	OnNext func(*interface{})

	originMonadIO *MonadIODef
}

func (self MonadIODef) JustVal(in interface{}) *MonadIODef {
	return self.Just(&in)
}
func (self MonadIODef) Just(in *interface{}) *MonadIODef {
	return &MonadIODef{effect: func() *interface{} {
		return in
	}}
}

func (self *MonadIODef) FlatMap(fn func(*interface{}) *MonadIODef) *MonadIODef {

	return &MonadIODef{effect: func() *interface{} {
		next := fn(self.doEffect())
		return next.doEffect()
	}}

}
func (self *MonadIODef) Subscribe(s Subscription) *Subscription {
	s.originMonadIO = self
	if s.OnNext != nil {
		var result *interface{}

		doSub := func() {
			s.OnNext(result)
		}
		doOb := func() {
			result = self.doEffect()

			if self.subOn != nil {
				self.subOn.Post(doSub)
			} else {
				doSub()
			}
		}
		if self.obOn != nil {
			self.obOn.Post(doOb)
		} else {
			doOb()
		}
	}

	return &s
}
func (self *MonadIODef) doEffect() *interface{} {
	return self.effect()
}

var MonadIO MonadIODef
