package dto

//Profile object
type Profile struct {
	Title                 string      `json:"title"`
	Phone                 string      `json:"phone"`
	Skype                 string      `json:"skype"`
	RealName              string      `json:"real_name"`
	RealNameNormalized    string      `json:"real_name_normalized"`
	DisplayName           string      `json:"display_name"`
	DisplayNameNormalized string      `json:"display_name_normalized"`
	Fields                interface{} `json:"fields"`
	StatusText            string      `json:"status_text"`
	StatusEmoji           string      `json:"status_emoji"`
	StatusExpiration      int         `json:"status_expiration"`
	AvatarHash            string      `json:"avatar_hash"`
	AlwaysActive          bool        `json:"always_active"`
	FirstName             string      `json:"first_name"`
	LastName              string      `json:"last_name"`
	Image24               string      `json:"image_24"`
	Image32               string      `json:"image_32"`
	Image48               string      `json:"image_48"`
	Image72               string      `json:"image_72"`
	Image192              string      `json:"image_192"`
	Image512              string      `json:"image_512"`
	StatusTextCanonical   string      `json:"status_text_canonical"`
	Team                  string      `json:"team"`
}

//SlackMember object which is used in SlackResponseUsersList object
type SlackMember struct {
	ID                string      `json:"id"`
	TeamID            string      `json:"team_id"`
	Name              string      `json:"name"`
	Deleted           bool        `json:"deleted"`
	Color             string      `json:"color"`
	RealName          string      `json:"real_name"`
	Tz                interface{} `json:"tz"`
	TzLabel           string      `json:"tz_label"`
	TzOffset          int         `json:"tz_offset"`
	Profile           Profile     `json:"profile"`
	IsAdmin           bool        `json:"is_admin"`
	IsOwner           bool        `json:"is_owner"`
	IsPrimaryOwner    bool        `json:"is_primary_owner"`
	IsRestricted      bool        `json:"is_restricted"`
	IsUltraRestricted bool        `json:"is_ultra_restricted"`
	IsBot             bool        `json:"is_bot"`
	IsAppUser         bool        `json:"is_app_user"`
	Updated           int         `json:"updated"`
}

//ResponseMetadata object of metadata
type ResponseMetadata struct {
	NextCursor string `json:"next_cursor"`
}

//SlackResponseUsersList response object from users.list
type SlackResponseUsersList struct {
	Ok               bool             `json:"ok"`
	Error            string           `json:"error,omitempty"`
	Members          []SlackMember    `json:"members"`
	CacheTs          int              `json:"cache_ts"`
	ResponseMetadata ResponseMetadata `json:"response_metadata"`
}
