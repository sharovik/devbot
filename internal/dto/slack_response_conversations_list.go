package dto

//Channel it's an object which contains information about message channel
type Channel struct {
	ID                      string        `json:"id"`
	Name                    string        `json:"name"`
	IsChannel               bool          `json:"is_channel"`
	IsGroup                 bool          `json:"is_group"`
	IsIm                    bool          `json:"is_im"`
	Created                 int           `json:"created"`
	IsArchived              bool          `json:"is_archived"`
	IsGeneral               bool          `json:"is_general"`
	Unlinked                int           `json:"unlinked"`
	NameNormalized          string        `json:"name_normalized"`
	IsShared                bool          `json:"is_shared"`
	ParentConversation      interface{}   `json:"parent_conversation"`
	Creator                 string        `json:"creator"`
	IsExtShared             bool          `json:"is_ext_shared"`
	IsOrgShared             bool          `json:"is_org_shared"`
	SharedTeamIds           []string      `json:"shared_team_ids"`
	PendingShared           []interface{} `json:"pending_shared"`
	PendingConnectedTeamIds []interface{} `json:"pending_connected_team_ids"`
	IsPendingExtShared      bool          `json:"is_pending_ext_shared"`
	IsMember                bool          `json:"is_member"`
	IsPrivate               bool          `json:"is_private"`
	IsMpim                  bool          `json:"is_mpim"`
	Topic                   struct {
		Value   string `json:"value"`
		Creator string `json:"creator"`
		LastSet int    `json:"last_set"`
	} `json:"topic"`
	Purpose struct {
		Value   string `json:"value"`
		Creator string `json:"creator"`
		LastSet int    `json:"last_set"`
	} `json:"purpose"`
	PreviousNames []interface{} `json:"previous_names"`
	NumMembers    int           `json:"num_members"`
}

//SlackResponseConversationsList the list of Channel objects
type SlackResponseConversationsList struct {
	Ok               bool      `json:"ok"`
	Error            string    `json:"error,omitempty"`
	Channels         []Channel `json:"channels"`
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}
