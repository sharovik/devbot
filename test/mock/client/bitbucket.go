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