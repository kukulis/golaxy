package ws

const (
	EventClientDisconnected     = "ws.client-disconnected"
	EventMessageArrived         = "ws.message-arrived"
	EventMessageArrivedError    = "ws.message-arrived-error"
	EventMessageSentToHub       = "ws.message-sent-to-hub"
	EventBeforeMessageSend      = "ws.before-message-send"
	EventWriterInitializeError  = "ws.writer-initialize-error"
	EventMessageSendError       = "ws.message-send-error"
	EventMessageSendSuccess     = "ws.message-send-success"
	EventSendMessageChannelFail = "ws.send-message-channel-fail"
	EventCheckReceiver          = "ws.check-receiver"
)

type WsEvent struct {
	Name    string
	Token   string
	Payload []byte
}

func NewWsEvent(name string, token string, payload []byte) *WsEvent {
	return &WsEvent{Name: name, Token: token, Payload: payload}
}

func (e *WsEvent) GetName() string { return e.Name }

type WsCheckSendEvent struct {
	Name          string
	Token         string
	Payload       []byte
	Send          bool
	ReceiverToken string
}

func NewWsCheckSendEvent(name string, token string, payload []byte) *WsCheckSendEvent {
	return &WsCheckSendEvent{Name: name, Token: token, Payload: payload, Send: true, ReceiverToken: ""}
}

func (e *WsCheckSendEvent) GetName() string { return e.Name }
