package socket

type handler struct {
	events     map[string]func(Socket, ContentType)
	err        func(Socket, error)
	connect    func(Socket) error
	disconnect func(Socket, string)
}

func newHandler() *handler {
	return &handler{
		events: make(map[string]func(Socket, ContentType)),
	}
}

func (h *handler) addEventHandler(event string, f func(Socket, ContentType)) {
	h.events[event] = f
}

func (h *handler) addErrorHandler(f func(Socket, error)) {
	h.err = f
}

func (h *handler) addConnectHandler(f func(Socket) error) {
	h.connect = f
}

func (h *handler) addDisconnectHandler(f func(Socket, string)) {
	h.disconnect = f
}
