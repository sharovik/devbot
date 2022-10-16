package client

import (
	"github.com/sharovik/devbot/internal/dto"
)

type BaseMessageClient struct {
	HttpClient BaseHTTPClientInterface
}

func (c BaseMessageClient) GetHTTPClient() BaseHTTPClientInterface {
	return c.HttpClient
}

//MessageClientInterface interface for slack client
type MessageClientInterface interface {
	GetHTTPClient() BaseHTTPClientInterface

	//Methods for slackAPI endpoints
	GetConversationsList() (dto.SlackResponseConversationsList, int, error)
	GetUsersList() (dto.SlackResponseUsersList, int, error)

	//PM messages
	SendMessage(dto.SlackRequestChatPostMessage) (dto.SlackResponseChatPostMessage, int, error)
	SendMessageV2(message dto.BaseChatMessage) (response dto.BaseResponseInterface, status int, err error)

	//Send attachment
	AttachFileTo(channel string, pathToFile string, filename string) ([]byte, int, error)
}
