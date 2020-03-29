package mock

import (
	"github.com/sharovik/devbot/internal/client"
	"github.com/sharovik/devbot/internal/dto"
)

type MockedBitBucketClient struct {
	IsTokenInvalid bool
	BeforeRequest error
	LoadAuthToken error

	PullRequestInfoResponse dto.BitBucketPullRequestInfoResponse
	PullRequestInfoError error

	MergePullRequestResponse dto.BitBucketPullRequestInfoResponse
	MergePullRequestError error

	CreateBranchResponse dto.BitBucketResponseBranchCreate
	CreateBranchError error

	ChangePullRequestDestinationResponse dto.BitBucketPullRequestInfoResponse
	ChangePullRequestDestinationError error

	CreatePullRequestResponse dto.BitBucketPullRequestInfoResponse
	CreatePullRequestError    error
}

func (b *MockedBitBucketClient) Init(client client.BaseHttpClientInterface)  {

}

func (b *MockedBitBucketClient) isTokenInvalid() bool {
	return b.IsTokenInvalid
}

func (b *MockedBitBucketClient) beforeRequest() error {
	return b.BeforeRequest
}

func (b *MockedBitBucketClient) loadAuthToken() error {
	return b.LoadAuthToken
}

func (b *MockedBitBucketClient) PullRequestInfo(workspace string, repositorySlug string, pullRequestID int64) (dto.BitBucketPullRequestInfoResponse, error) {
	return b.PullRequestInfoResponse, b.PullRequestInfoError
}

func (b *MockedBitBucketClient) MergePullRequest(workspace string, repositorySlug string, pullRequestID int64, description string) (dto.BitBucketPullRequestInfoResponse, error) {
	return b.MergePullRequestResponse, b.MergePullRequestError
}

func (b *MockedBitBucketClient) CreateBranch(workspace string, repositorySlug string, branchName string) (dto.BitBucketResponseBranchCreate, error) {
	return b.CreateBranchResponse, b.CreateBranchError
}

func (b *MockedBitBucketClient) ChangePullRequestDestination(workspace string, repositorySlug string, pullRequestID int64, title string, branchName string) (dto.BitBucketPullRequestInfoResponse, error) {
	return b.ChangePullRequestDestinationResponse, b.ChangePullRequestDestinationError
}

func (b *MockedBitBucketClient) CreatePullRequest(workspace string, repositorySlug string, request dto.BitBucketRequestPullRequestCreate) (dto.BitBucketPullRequestInfoResponse, error) {
	return b.CreatePullRequestResponse, b.CreatePullRequestError
}