package dto

// Team child struct of SlackResponseRTMConnect
type Team struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
}

// Self child struct of SlackResponseRTMConnect
type Self struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SlackResponseRTMConnect need for decoding an rtm.connect endpoint response
type SlackResponseRTMConnect struct {
	Ok    bool   `json:"ok"`
	URL   string `json:"url"`
	Error string `json:"error"`
	Team  Team   `json:"team"`
	Self  Self   `json:"self"`
}
