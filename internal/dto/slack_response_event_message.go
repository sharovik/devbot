package dto

//File child struct of SlackResponseEventMessage
type File struct {
	Created            int    `json:"created"`
	DisplayAsBot       bool   `json:"display_as_bot"`
	Editable           bool   `json:"editable"`
	ExternalType       string `json:"external_type"`
	Filetype           string `json:"filetype"`
	HasRichPreview     bool   `json:"has_rich_preview"`
	ID                 string `json:"id"`
	ImageExifRotation  int    `json:"image_exif_rotation"`
	IsExternal         bool   `json:"is_external"`
	IsPublic           bool   `json:"is_public"`
	IsStarred          bool   `json:"is_starred"`
	Mimetype           string `json:"mimetype"`
	Mode               string `json:"mode"`
	Name               string `json:"name"`
	OriginalH          int    `json:"original_h"`
	OriginalW          int    `json:"original_w"`
	Permalink          string `json:"permalink"`
	PermalinkPublic    string `json:"permalink_public"`
	PrettyType         string `json:"pretty_type"`
	PublicURLShared    bool   `json:"public_url_shared"`
	Size               int    `json:"size"`
	Thumb1024          string `json:"thumb_1024"`
	Thumb1024H         int    `json:"thumb_1024_h"`
	Thumb1024W         int    `json:"thumb_1024_w"`
	Thumb160           string `json:"thumb_160"`
	Thumb360           string `json:"thumb_360"`
	Thumb360H          int    `json:"thumb_360_h"`
	Thumb360W          int    `json:"thumb_360_w"`
	Thumb480           string `json:"thumb_480"`
	Thumb480H          int    `json:"thumb_480_h"`
	Thumb480W          int    `json:"thumb_480_w"`
	Thumb64            string `json:"thumb_64"`
	Thumb720           string `json:"thumb_720"`
	Thumb720H          int    `json:"thumb_720_h"`
	Thumb720W          int    `json:"thumb_720_w"`
	Thumb80            string `json:"thumb_80"`
	Thumb800           string `json:"thumb_800"`
	Thumb800H          int    `json:"thumb_800_h"`
	Thumb800W          int    `json:"thumb_800_w"`
	Thumb960           string `json:"thumb_960"`
	Thumb960H          int    `json:"thumb_960_h"`
	Thumb960W          int    `json:"thumb_960_w"`
	ThumbTiny          string `json:"thumb_tiny"`
	Timestamp          int    `json:"timestamp"`
	Title              string `json:"title"`
	URLPrivate         string `json:"url_private"`
	URLPrivateDownload string `json:"url_private_download"`
	User               string `json:"user"`
	Username           string `json:"username"`
}

//SlackResponseEventMessage main event message object
type SlackResponseEventMessage struct {
	Channel      string `json:"channel"`
	ClientMsgID  string `json:"client_msg_id"`
	DisplayAsBot bool   `json:"display_as_bot"`
	EventTs      string `json:"event_ts"`
	Files        []File `json:"files"`
	SourceTeam   string `json:"source_team"`
	Team         string `json:"team"`
	Text         string `json:"text"`
	Ts           string `json:"ts"`
	Type         string `json:"type"`
	Upload       bool   `json:"upload"`
	User         string `json:"user"`
	UserTeam     string `json:"user_team"`
}
