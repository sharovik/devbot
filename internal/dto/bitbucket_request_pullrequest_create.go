package dto

// BitBucketRequestPullRequestCreate used for pull-request create requests
type BitBucketRequestPullRequestCreate struct {
	Title       string                          `json:"title"`
	Description string                          `json:"description"`
	Source      BitBucketPullRequestDestination `json:"source"`
	Reviewers   []struct {
		UUID string `json:"uuid"`
	} `json:"reviewers"`
}

// BitBucketReviewer the reviewer object
type BitBucketReviewer struct {
	UUID string `json:"uuid"`
}
