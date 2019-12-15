package dto

//SlackRequestEventMessage object for slack WS event message
type SlackRequestEventMessage struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}
