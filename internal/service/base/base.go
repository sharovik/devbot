package base

//ServiceInterface base interface for messages APIs services
type ServiceInterface interface {
	InitWebSocketReceiver() error
}
