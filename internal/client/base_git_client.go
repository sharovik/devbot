package client

import "github.com/sharovik/devbot/internal/dto"

//GitClientInterface the main git client interface which should be used once you connect VCS API
type GitClientInterface interface {
	Init(client BaseHTTPClientInterface)
	PullRequestInfo(workspace string, repositorySlug string, pullRequestID int64) (dto.BitBucketPullRequestInfoResponse, error)
	MergePullRequest(workspace string, repositorySlug string, pullRequestID int64, description string) (dto.BitBucketPullRequestInfoResponse, error)
	CreateBranch(workspace string, repositorySlug string, branchName string) (dto.BitBucketResponseBranchCreate, error)
	ChangePullRequestDestination(workspace string, repositorySlug string, pullRequestID int64, title string, branchName string) (dto.BitBucketPullRequestInfoResponse, error)
	CreatePullRequest(workspace string, repositorySlug string, request dto.BitBucketRequestPullRequestCreate) (dto.BitBucketPullRequestInfoResponse, error)
	RunPipeline(workspace string, repositorySlug string, request dto.BitBucketRequestRunPipeline) (dto.BitBucketResponseRunPipeline, error)
}
