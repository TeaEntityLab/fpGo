package fpGo

// HandlerDef Handler inspired by Android/WebWorker
type HandlerDef struct {
	isClosed bool

	ch *chan func()
}

var defaultHandler *HandlerDef

// GetDefault Get Default Handler
func (handlerSelf *HandlerDef) GetDefault() *HandlerDef {
	return defaultHandler
}

// New New Handler instance
func (handlerSelf *HandlerDef) New() *HandlerDef {
	ch := make(chan func())
	return handlerSelf.NewByCh(&ch)
}

// NewByCh New Handler by its Channel
func (handlerSelf *HandlerDef) NewByCh(ioCh *chan func()) *HandlerDef {
	new := HandlerDef{ch: ioCh}
	go new.run()

	return &new
}

// Post Post a function to execute on the Handler
func (handlerSelf *HandlerDef) Post(fn func()) {
	if handlerSelf.isClosed {
		return
	}

	*(handlerSelf.ch) <- fn
}

// Close Close the Handler
func (handlerSelf *HandlerDef) Close() {
	handlerSelf.isClosed = true

	close(*handlerSelf.ch)
}
func (handlerSelf *HandlerDef) run() {
	for fn := range *handlerSelf.ch {
		fn()
	}
}

// Handler Handler utils instance
var Handler HandlerDef

func init() {
	Handler = *Handler.New()
	defaultHandler = &Handler
}
