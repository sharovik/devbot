package dto

//BitBucketPullRequestDestinationUpdateRequest request struct for pull-request destination update
type BitBucketPullRequestDestinationUpdateRequest struct {
	Title string `json:"title"`
	Destination BitBucketPullRequestDestination `json:"destination"`
}

//BitBucketPullRequestDestinationBranch the destination branch struct
type BitBucketPullRequestDestinationBranch struct {
	Name string `json:"name"`
}

//BitBucketPullRequestDestination the destination struct
type BitBucketPullRequestDestination struct {
	Branch BitBucketPullRequestDestinationBranch `json:"branch"`
}