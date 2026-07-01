package fpgo

// HandlerDef Handler inspired by Android/WebWorker
type HandlerDef struct {
	isClosed AtomBool

	ch   chan func()
	done chan struct{}
}

var defaultHandler *HandlerDef

// GetDefault Get Default Handler
func (handlerSelf *HandlerDef) GetDefault() *HandlerDef {
	return defaultHandler
}

// New New Handler instance
func (handlerSelf *HandlerDef) New() *HandlerDef {
	return handlerSelf.NewByCh(make(chan func()))
}

// NewByCh New Handler by its Channel
func (handlerSelf *HandlerDef) NewByCh(ioCh chan func()) *HandlerDef {
	newOne := HandlerDef{ch: ioCh, done: make(chan struct{})}
	go newOne.run()

	return &newOne
}

// Post Post a function to execute on the Handler
func (handlerSelf *HandlerDef) Post(fn func()) {
	if handlerSelf.isClosed.Get() {
		return
	}

	select {
	case <-handlerSelf.done:
	case handlerSelf.ch <- fn:
	}
}

// Close Close the Handler
func (handlerSelf *HandlerDef) Close() {
	if !handlerSelf.isClosed.CompareAndSwap(false, true) {
		return
	}

	close(handlerSelf.done)
}

// IsClosed Check if handler is closed
func (handlerSelf *HandlerDef) IsClosed() bool {
	return handlerSelf.isClosed.Get()
}

func (handlerSelf *HandlerDef) run() {
	for {
		select {
		case <-handlerSelf.done:
			return
		case fn, ok := <-handlerSelf.ch:
			if !ok {
				return
			}
			if fn != nil {
				fn()
			}
		}
	}
}

// Handler Handler utils instance
var Handler HandlerDef

func init() {
	Handler = *Handler.New()
	defaultHandler = &Handler
}
