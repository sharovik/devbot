package dto

import "time"

//File object
type File struct {
	Created            int    `json:"created"`
	Editable           bool   `json:"editable"`
	ExternalType       string `json:"external_type"`
	Filetype           string `json:"filetype"`
	ID                 string `json:"id"`
	IsExternal         bool   `json:"is_external"`
	IsPublic           bool   `json:"is_public"`
	IsStarred          bool   `json:"is_starred"`
	Mimetype           string `json:"mimetype"`
	Mode               string `json:"mode"`
	Name               string `json:"name"`
	Permalink          string `json:"permalink"`
	PermalinkPublic    string `json:"permalink_public"`
	PublicURLShared    bool   `json:"public_url_shared"`
	Size               int    `json:"size"`
	Timestamp          int    `json:"timestamp"`
	Title              string `json:"title"`
	URLPrivate         string `json:"url_private"`
	URLPrivateDownload string `json:"url_private_download"`
	User               string `json:"user"`
	Username           string `json:"username"`
}

//BaseOriginalMessage original message interface which will be used for identification of base original message object
type BaseOriginalMessage struct {
	Text        string
	User        string
	Files       []File
	Channel     string
	ClientMsgID string
	EventTs     string
	ThreadTS    string
	Ts          string
	Type        string
}

//BaseChatMessage the chat message which will be retrieved from websocket api
type BaseChatMessage struct {
	//The channel from where was message received
	Channel string

	//The text of the message, which was received from the channel. Example: Hello bot
	Text string

	//This is an optional value which is used currently for messages.
	AsUser bool

	ThreadTS string

	//The message timestamp. When it was received
	Ts time.Time

	//The dictionary answer which was found in our database
	DictionaryMessage DictionaryMessage

	//The copy of original message received from the system
	OriginalMessage BaseOriginalMessage
}
