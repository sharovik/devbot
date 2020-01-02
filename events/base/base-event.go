package base

import "github.com/sharovik/devbot/internal/dto"

type Event interface {
	Execute(message dto.SlackRequestChatPostMessage) (dto.SlackRequestChatPostMessage, error)
}

type Events struct {
	Events map[string]Event
}
