package fpGo

type MonadDefAsync struct {
	MonadDef
}

func (self MonadDefAsync) FlatMap(fn func(MonadProto) MonadProto) MonadProto {
	return fn(self)
}
