package dto

//MessageBlock child struct of SlackResponseEventMessage
type MessageBlock struct {
	BlockID  string `json:"block_id"`
	Elements []struct {
		Elements []struct {
			Type   string `json:"type"`
			UserID string `json:"user_id,omitempty"`
			Text   string `json:"text,omitempty"`
		} `json:"elements"`
		Type string `json:"type"`
	} `json:"elements"`
	Type string `json:"type"`
}

//SlackResponseEventMessage main event message object
type SlackResponseEventMessage struct {
	Channel      string         `json:"channel"`
	ClientMsgID  string         `json:"client_msg_id"`
	DisplayAsBot bool           `json:"display_as_bot"`
	EventTs      string         `json:"event_ts"`
	ThreadTS     string         `json:"thread_ts"`
	Files        []File         `json:"files"`
	SourceTeam   string         `json:"source_team"`
	Team         string         `json:"team"`
	Text         string         `json:"text"`
	Ts           string         `json:"ts"`
	Type         string         `json:"type"`
	Upload       bool           `json:"upload"`
	User         string         `json:"user"`
	UserTeam     string         `json:"user_team"`
	Blocks       []MessageBlock `json:"blocks"`
}
