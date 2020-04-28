package dto

import "time"

//BaseOriginalMessage original message interface which will be used for identification of base original message object
type BaseOriginalMessage struct {
	Text string
	User string
}

//BaseChatMessage the chat message which will be retrieved from websocket api
type BaseChatMessage struct {
	Channel         string
	Text            string
	AsUser          bool
	Ts              time.Time
	DictionaryMessage DictionaryMessage
	OriginalMessage BaseOriginalMessage
}
