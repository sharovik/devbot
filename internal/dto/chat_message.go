package dto

import "time"

//ChatMessage request for post.chatMessage
type ChatMessage struct {
	Channel           string `json:"channel"`
	Text              string `json:"text"`
	AsUser            bool   `json:"as_user"`
	Ts                time.Time
	DictionaryMessage DictionaryMessage
	OriginalMessage   EventMessage
}
