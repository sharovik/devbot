package dto

//SlackResponseEventApiMessage the main object of the response
type SlackResponseEventApiMessage struct {
	AcceptsResponsePayload bool   `json:"accepts_response_payload"`
	EnvelopeID             string `json:"envelope_id"`
	Payload                struct {
		APIAppID       string `json:"api_app_id"`
		Authorizations []struct {
			EnterpriseID        interface{} `json:"enterprise_id"`
			IsBot               bool        `json:"is_bot"`
			IsEnterpriseInstall bool        `json:"is_enterprise_install"`
			TeamID              string      `json:"team_id"`
			UserID              string      `json:"user_id"`
		} `json:"authorizations"`
		Event struct {
			Blocks []struct {
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
			} `json:"blocks"`
			Channel     string `json:"channel"`
			BotID     string `json:"bot_id"`
			ClientMsgID string `json:"client_msg_id"`
			EventTs     string `json:"event_ts"`
			Team        string `json:"team"`
			Text        string `json:"text"`
			Ts          string `json:"ts"`
			Type        string `json:"type"`
			User        string `json:"user"`
		} `json:"event"`
		EventContext       string `json:"event_context"`
		EventID            string `json:"event_id"`
		EventTime          int    `json:"event_time"`
		IsExtSharedChannel bool   `json:"is_ext_shared_channel"`
		TeamID             string `json:"team_id"`
		Token              string `json:"token"`
		Type               string `json:"type"`
	} `json:"payload"`
	RetryAttempt int    `json:"retry_attempt"`
	RetryReason  string `json:"retry_reason"`
	Type         string `json:"type"`
}