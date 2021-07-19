package fpgo

// MonadIODef MonadIO inspired by Rx/Observable
type MonadIODef[T any] struct {
	effect func() T

	obOn  *HandlerDef
	subOn *HandlerDef
}

// Subscription the delegation/callback of MonadIO/Publisher
type Subscription[T any] struct {
	OnNext func(T)
}

// Just New MonadIO by a given value
func (monadIOSelf MonadIODef[T]) Just(in interface{}) *MonadIODef[interface{}] {
	return MonadIOJustGenerics(in)
}

// MonadIOJustGenerics New MonadIO by a given value
func MonadIOJustGenerics[T any](in T) *MonadIODef[T] {
	return &MonadIODef[T]{effect: func() T {
		return in
	}}
}

// New New MonadIO by effect function
func (monadIOSelf *MonadIODef[T]) New(effect func() T) *MonadIODef[T] {
	return MonadIONewGenerics(effect)
}

// MonadIONewGenerics New MonadIO by effect function
func MonadIONewGenerics[T any](effect func() T) *MonadIODef[T] {
	return &MonadIODef[T]{effect: effect}
}

// FlatMap FlatMap the MonadIO by function
func (monadIOSelf *MonadIODef[T]) FlatMap(fn func(T) *MonadIODef[T]) *MonadIODef[T] {

	return &MonadIODef[T]{effect: func() T {
		next := fn(monadIOSelf.doEffect())
		return next.doEffect()
	}}

}

// Subscribe Subscribe the MonadIO by Subscription
func (monadIOSelf *MonadIODef[T]) Subscribe(s Subscription[T]) *Subscription[T] {
	obOn := monadIOSelf.obOn
	subOn := monadIOSelf.subOn
	return monadIOSelf.doSubscribe(&s, obOn, subOn)
}

// SubscribeOn Subscribe the MonadIO on the specific Handler
func (monadIOSelf *MonadIODef[T]) SubscribeOn(h *HandlerDef) *MonadIODef[T] {
	monadIOSelf.subOn = h
	return monadIOSelf
}

// ObserveOn Observe the MonadIO on the specific Handler
func (monadIOSelf *MonadIODef[T]) ObserveOn(h *HandlerDef) *MonadIODef[T] {
	monadIOSelf.obOn = h
	return monadIOSelf
}
func (monadIOSelf *MonadIODef[T]) doSubscribe(s *Subscription[T], obOn *HandlerDef, subOn *HandlerDef) *Subscription[T] {

	if s.OnNext != nil {
		var result T

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
func (monadIOSelf *MonadIODef[T]) doEffect() T {
	return monadIOSelf.effect()
}

// Eval Eval the value right now(sync)
func (monadIOSelf *MonadIODef[T]) Eval() T {
	return monadIOSelf.doEffect()
}

// MonadIO MonadIO utils instance
var MonadIO MonadIODef[interface{}]
