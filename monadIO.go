package fpGo

type MonadIODef struct {
	effect func() *interface{}
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

func (self *MonadIODef) FlatMap(fn func(*interface{}) *MonadIODef) *MonadIODef {

	return &MonadIODef{effect: func() *interface{} {
		next := fn(self.doEffect())
		return next.doEffect()
	}}

}
func (self *MonadIODef) Subscribe(s Subscription) *Subscription {
	if s.OnNext != nil {
		s.OnNext(self.doEffect())
	}

	return &s
}
func (self *MonadIODef) doEffect() *interface{} {
	return self.effect()
}

var MonadIO MonadIODef
