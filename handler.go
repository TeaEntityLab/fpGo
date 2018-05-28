package fpGo

type HandlerDef struct {
	ch *chan func()
}

var defaultHandler *HandlerDef = nil

func (self *HandlerDef) GetDefault() *HandlerDef {
	return defaultHandler
}
func (self *HandlerDef) NewByCh(ioCh *chan func()) *HandlerDef {
	new := HandlerDef{ch: ioCh}
	go new.run()

	return &new
}
func (self *HandlerDef) Post(fn func()) {
	*(self.ch) <- fn
}
func (self *HandlerDef) run() {
	for fn := range *self.ch {
		fn()
	}
}

var Handler HandlerDef

func init() {
	mainCh := make(chan func())
	defaultHandler = Handler.NewByCh(&mainCh)
}
