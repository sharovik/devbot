package dto

// BitBucketResponseDefaultReviewers the bitbucket reviewers response object
type BitBucketResponseDefaultReviewers struct {
	Pagelen int `json:"pagelen"`
	Values  []struct {
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
		Type      string `json:"type"`
		Nickname  string `json:"nickname"`
		AccountID string `json:"account_id"`
	} `json:"values"`
	Page int `json:"page"`
	Size int `json:"size"`
}
