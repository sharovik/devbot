package dto

import "time"

// BitBucketResponseRunPipeline the response struct of run-pipeline request
type BitBucketResponseRunPipeline struct {
	Type       string `json:"type"`
	UUID       string `json:"uuid"`
	Repository struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		FullName string `json:"full_name"`
		Links    struct {
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
		UUID string `json:"uuid"`
	} `json:"repository"`
	State struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Stage struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"stage"`
	} `json:"state"`
	BuildNumber int `json:"build_number"`
	Creator     struct {
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
	} `json:"creator"`
	CreatedOn time.Time `json:"created_on"`
	Target    struct {
		Type     string `json:"type"`
		RefType  string `json:"ref_type"`
		RefName  string `json:"ref_name"`
		Selector struct {
			Type    string `json:"type"`
			Pattern string `json:"pattern"`
		} `json:"selector"`
		Commit struct {
			Type  string `json:"type"`
			Hash  string `json:"hash"`
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
			} `json:"links"`
		} `json:"commit"`
	} `json:"target"`
	Trigger struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"trigger"`
	RunNumber         int  `json:"run_number"`
	DurationInSeconds int  `json:"duration_in_seconds"`
	BuildSecondsUsed  int  `json:"build_seconds_used"`
	Expired           bool `json:"expired"`
	Links             struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Steps struct {
			Href string `json:"href"`
		} `json:"steps"`
	} `json:"links"`
	HasVariables bool `json:"has_variables"`
	Error        struct {
		Message string `json:"message"`
		Detail  string `json:"detail"`
		Data    struct {
			Key       string `json:"key"`
			Arguments struct {
			} `json:"arguments"`
		} `json:"data"`
	} `json:"error"`
}
