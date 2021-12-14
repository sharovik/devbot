package client

import (
	"github.com/sharovik/devbot/internal/dto"
)

//MessageClientInterface interface for slack client
type MessageClientInterface interface {
	//Http methods
	Request(string, string, []byte) ([]byte, int, error)
	Post(string, []byte) ([]byte, int, error)
	Get(string) ([]byte, int, error)
	Put(string, []byte) ([]byte, int, error)

	//Methods for slackAPI endpoints
	GetConversationsList() (dto.SlackResponseConversationsList, int, error)
	GetUsersList() (dto.SlackResponseUsersList, int, error)

	//PM messages
	SendMessage(dto.SlackRequestChatPostMessage) (dto.SlackResponseChatPostMessage, int, error)

	//Send attachment
	AttachFileTo(channel string, pathToFile string, filename string) ([]byte, int, error)
}
