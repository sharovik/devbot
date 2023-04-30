package message

import (
	"github.com/sharovik/devbot/internal/config"
	"github.com/sharovik/devbot/internal/container"
)

// BaseServiceInterface base interface for messages APIs services
type BaseServiceInterface interface {
	InitWebSocketReceiver() error
	BeforeWSConnectionStart() error
	ProcessMessage(message interface{}) error
}

// S message service object
var S BaseServiceInterface

// InitService initialize the events-api service
func InitService() {
	switch container.C.Config.MessagesAPIConfig.Type {
	case config.MessagesAPITypeSlack:
		S = SlackService{}
	default:
		panic("The messages api type is not supported")
	}
}
