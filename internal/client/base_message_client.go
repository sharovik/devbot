package client

import (
	"github.com/sharovik/devbot/internal/dto"
)

//BaseMessageClient the base messages client
type BaseMessageClient struct {
	HTTPClient BaseHTTPClientInterface
}

//GetHTTPClient method for retrieving of the current http client
func (c BaseMessageClient) GetHTTPClient() BaseHTTPClientInterface {
	return c.HTTPClient
}

//MessageClientInterface interface for message client
type MessageClientInterface interface {
	GetHTTPClient() BaseHTTPClientInterface

	//Methods for slackAPI endpoints
	GetConversationsList() (dto.SlackResponseConversationsList, int, error)
	GetUsersList() (dto.SlackResponseUsersList, int, error)

	//SendMessage sends the message to selected channel
	SendMessage(message dto.BaseChatMessage) (response dto.BaseResponseInterface, status int, err error)

	//AttachFileTo send attachment
	AttachFileTo(channel string, pathToFile string, filename string) ([]byte, int, error)
}
