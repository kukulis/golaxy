package ws

type WsEvent struct {
	name    string
	token   string
	payload []byte
}

func NewWsEvent(name string, token string, payload []byte) *WsEvent {
	return &WsEvent{name: name, token: token, payload: payload}
}

func (e *WsEvent) GetName() string    { return e.name }
func (e *WsEvent) GetToken() string   { return e.token }
func (e *WsEvent) GetPayload() []byte { return e.payload }
