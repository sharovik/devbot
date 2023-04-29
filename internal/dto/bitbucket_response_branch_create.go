package dto

import "time"

// DataKey the key struct
type DataKey struct {
	Key string `json:"key"`
}

// Error the error struct
type Error struct {
	Message string  `json:"message"`
	Data    DataKey `json:"data"`
}

// BitBucketErrorResponseBranchCreate the bad response from the bitbucket api
type BitBucketErrorResponseBranchCreate struct {
	Data  DataKey `json:"data"`
	Type  string  `json:"type"`
	Error Error   `json:"error"`
}

// BitBucketResponseBranchCreate will be used for parsing of the response for branch-create request
type BitBucketResponseBranchCreate struct {
	Name  string `json:"name"`
	Links struct {
		Commits struct {
			Href string `json:"href"`
		} `json:"commits"`
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
	DefaultMergeStrategy string   `json:"default_merge_strategy"`
	MergeStrategies      []string `json:"merge_strategies"`
	Type                 string   `json:"type"`
	Target               struct {
		Hash       string `json:"hash"`
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
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			Comments struct {
				Href string `json:"href"`
			} `json:"comments"`
			Patch struct {
				Href string `json:"href"`
			} `json:"patch"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Diff struct {
				Href string `json:"href"`
			} `json:"diff"`
			Approve struct {
				Href string `json:"href"`
			} `json:"approve"`
			Statuses struct {
				Href string `json:"href"`
			} `json:"statuses"`
		} `json:"links"`
		Author struct {
			Raw  string `json:"raw"`
			Type string `json:"type"`
			User Author `json:"user"`
		} `json:"author"`
		Parents []struct {
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
		} `json:"parents"`
		Date    time.Time `json:"date"`
		Message string    `json:"message"`
		Type    string    `json:"type"`
	} `json:"target"`
	Data  DataKey `json:"data"`
	Error Error   `json:"error"`
}
