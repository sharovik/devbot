package dto

import "time"

//SlackRequestChatPostMessage request for post.chatMessage
type SlackRequestChatPostMessage struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
	AsUser  bool   `json:"as_user"`
	Ts      time.Time
	DictionaryMessage DictionaryMessage
}
