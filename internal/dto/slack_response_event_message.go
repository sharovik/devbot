package dto

type SlackResponseEventMessage struct {
	Type       string `json:"type"`
	Ts         string `json:"ts"`
	User       string `json:"user"`
	Team       string `json:"team"`
	SourceTeam string `json:"source_team"`
	Channel    string `json:"channel"`
	UserTeam   string `json:"user_team"`
	Text       string `json:"text"`
}
