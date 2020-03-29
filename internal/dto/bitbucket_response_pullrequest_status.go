package dto

import "time"

//Rendered the rendered information
type Rendered struct {
	Description struct {
		Raw    string `json:"raw"`
		Markup string `json:"markup"`
		HTML   string `json:"html"`
		Type   string `json:"type"`
	} `json:"description"`
	Title struct {
		Raw    string `json:"raw"`
		Markup string `json:"markup"`
		HTML   string `json:"html"`
		Type   string `json:"type"`
	} `json:"title"`
}

//Links the links action object
type Links struct {
	Decline struct {
		Href string `json:"href"`
	} `json:"decline"`
	Diffstat struct {
		Href string `json:"href"`
	} `json:"diffstat"`
	Commits struct {
		Href string `json:"href"`
	} `json:"commits"`
	Self struct {
		Href string `json:"href"`
	} `json:"self"`
	Comments struct {
		Href string `json:"href"`
	} `json:"comments"`
	Merge struct {
		Href string `json:"href"`
	} `json:"merge"`
	HTML struct {
		Href string `json:"href"`
	} `json:"html"`
	Activity struct {
		Href string `json:"href"`
	} `json:"activity"`
	Diff struct {
		Href string `json:"href"`
	} `json:"diff"`
	Approve struct {
		Href string `json:"href"`
	} `json:"approve"`
	Statuses struct {
		Href string `json:"href"`
	} `json:"statuses"`
}

//Reviewer the reviewer object
type Reviewer struct {
	DisplayName string `json:"display_name"`
	UUID        string `json:"uuid"`
	Links       struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
		Avatar struct {
			Href string `json:"href"`
		} `json:"avatar"`
	} `json:"links"`
	Nickname  string `json:"nickname"`
	Type      string `json:"type"`
	AccountID string `json:"account_id"`
}

//Destination the destination object
type Destination struct {
	Commit struct {
		Hash  string `json:"hash"`
		Type  string `json:"type"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
		} `json:"links"`
	} `json:"commit"`
	Repository struct {
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"links"`
		Type     string `json:"type"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		UUID     string `json:"uuid"`
	} `json:"repository"`
	Branch struct {
		Name string `json:"name"`
	} `json:"branch"`
}

//Summary the summary description
type Summary struct {
	Raw    string `json:"raw"`
	Markup string `json:"markup"`
	HTML   string `json:"html"`
	Type   string `json:"type"`
}

//Source the source object from where was made this pull-request
type Source struct {
	Commit struct {
		Hash  string `json:"hash"`
		Type  string `json:"type"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
		} `json:"links"`
	} `json:"commit"`
	Repository struct {
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"links"`
		Type     string `json:"type"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		UUID     string `json:"uuid"`
	} `json:"repository"`
	Branch struct {
		Name string `json:"name"`
	} `json:"branch"`
}

//ParticipantUser the user of participant object
type ParticipantUser struct {
	DisplayName string `json:"display_name"`
	UUID        string `json:"uuid"`
	Links       struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
		Avatar struct {
			Href string `json:"href"`
		} `json:"avatar"`
	} `json:"links"`
	Nickname  string `json:"nickname"`
	Type      string `json:"type"`
	AccountID string `json:"account_id"`
}

//Participant the participant object
type Participant struct {
	Role           string    `json:"role"`
	ParticipatedOn time.Time `json:"participated_on"`
	Type           string    `json:"type"`
	Approved       bool      `json:"approved"`
	User           ParticipantUser `json:"user"`
}

//Author the author object
type Author struct {
	DisplayName string `json:"display_name"`
	UUID        string `json:"uuid"`
	Links       struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
		Avatar struct {
			Href string `json:"href"`
		} `json:"avatar"`
	} `json:"links"`
	Nickname  string `json:"nickname"`
	Type      string `json:"type"`
	AccountID string `json:"account_id"`
}

//BitBucketPullRequestInfoResponse the response object from the BitBucket api
type BitBucketPullRequestInfoResponse struct {
	Rendered          Rendered      `json:"rendered"`
	Type              string        `json:"type"`
	Description       string        `json:"description"`
	Links             Links         `json:"links"`
	Title             string        `json:"title"`
	CloseSourceBranch bool          `json:"close_source_branch"`
	Reviewers         []Reviewer    `json:"reviewers"`
	ID                int           `json:"id"`
	Destination       Destination   `json:"destination"`
	CreatedOn         time.Time     `json:"created_on"`
	Summary           Summary       `json:"summary"`
	Source            Source        `json:"source"`
	CommentCount      int           `json:"comment_count"`
	State             string        `json:"state"`
	TaskCount         int           `json:"task_count"`
	Participants      []Participant `json:"participants"`
	Reason            string        `json:"reason"`
	UpdatedOn         time.Time     `json:"updated_on"`
	Author            Author        `json:"author"`
	MergeCommit       interface{}   `json:"merge_commit"`
	ClosedBy          interface{}   `json:"closed_by"`
	Error struct {
		Fields struct {
			Newstatus []string `json:"newstatus"`
		} `json:"fields"`
		Message string `json:"message"`
	} `json:"error"`
}
