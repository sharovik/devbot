package dto

// BitBucketBranchTarget the branch target struct
type BitBucketBranchTarget struct {
	Hash string `json:"hash"`
}

// BitBucketRequestBranchCreate need to be used for branch create requests
type BitBucketRequestBranchCreate struct {
	Name   string                `json:"name"`
	Target BitBucketBranchTarget `json:"target"`
}
