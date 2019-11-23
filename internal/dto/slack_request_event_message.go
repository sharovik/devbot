package dto

type SlackRequestEventMessage struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}
