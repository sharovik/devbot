package dto

//SlackRequestEventMessage object for message WS event message
type SlackRequestEventMessage struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}
