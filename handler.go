package fpGo

type HandlerDef struct {
	isClosed bool

	ch *chan func()
}

var defaultHandler *HandlerDef = nil

func (self *HandlerDef) GetDefault() *HandlerDef {
	return defaultHandler
}
func (self *HandlerDef) New() *HandlerDef {
	ch := make(chan func())
	return self.NewByCh(&ch)
}
func (self *HandlerDef) NewByCh(ioCh *chan func()) *HandlerDef {
	new := HandlerDef{ch: ioCh}
	go new.run()

	return &new
}
func (self *HandlerDef) Post(fn func()) {
	if self.isClosed {
		return
	}

	*(self.ch) <- fn
}
func (self *HandlerDef) Close() {
	self.isClosed = true

	close(*self.ch)
}
func (self *HandlerDef) run() {
	for fn := range *self.ch {
		fn()
	}
}

var Handler HandlerDef

func init() {
	Handler = *Handler.New()
	defaultHandler = &Handler
}
