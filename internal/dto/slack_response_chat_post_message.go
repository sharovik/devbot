package dto

//SlackResponseChatPostMessage response object for slack chat postMessage endpoint response
type SlackResponseChatPostMessage struct {
	Ok      bool   `json:"ok"`
	Channel string `json:"channel"`
	Ts      string `json:"ts"`
	Error   string `json:"error,omitempty"`
	Message struct {
		BotID      string `json:"bot_id"`
		Type       string `json:"type"`
		Text       string `json:"text"`
		User       string `json:"user"`
		Ts         string `json:"ts"`
		Team       string `json:"team"`
		BotProfile struct {
			ID      string `json:"id"`
			Deleted bool   `json:"deleted"`
			Name    string `json:"name"`
			Updated int    `json:"updated"`
			AppID   string `json:"app_id"`
			Icons   struct {
				Image36 string `json:"image_36"`
				Image48 string `json:"image_48"`
				Image72 string `json:"image_72"`
			} `json:"icons"`
			TeamID string `json:"team_id"`
		} `json:"bot_profile"`
	} `json:"message"`
}
