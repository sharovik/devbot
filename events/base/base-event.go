package base

import "github.com/sharovik/devbot/internal/dto"

//Event main interface for events
type Event interface {
	Execute(message dto.SlackRequestChatPostMessage) (dto.SlackRequestChatPostMessage, error)
}

//Events the struct object which is used for events storing
type Events struct {
	Events map[string]Event
}
