package bitbucket_release

import (
	"github.com/sharovik/devbot/internal/config"
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	mock "github.com/sharovik/devbot/test/mock/client"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"runtime"
	"testing"
)

func init() {
	//We switch pointer to the root directory for control the path from which we need to generate test-data file-paths
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	_ = os.Chdir(dir)
}

//@todo: cover all checks
func TestBitBucketReleaseEvent_Execute(t *testing.T) {
	container.C.Config.BitBucketConfig.RequiredReviewers = []config.BitBucketReviewer{
		{
			UUID: "{test-uid}",
			SlackUID: "TESTSLACKID",
		},
		{
			UUID: "{test-second-uid}",
			SlackUID: "TESTSECONDSLACKID",
		},
	}

	//PullRequest status OPEN but no participants
	container.C.BibBucketClient = &mock.MockedBitBucketClient{
		IsTokenInvalid:true,
		PullRequestInfoResponse: dto.BitBucketPullRequestInfoResponse{
			State: pullRequestStateOpen,
			Participants: []dto.Participant{},
		},
	}

	var msg = dto.SlackRequestChatPostMessage{
		OriginalMessage: dto.SlackResponseEventMessage{
			Text: `release these pull-requests to production: https://bitbucket.org/john/test-repo/pull-requests/1/testing-pr-flow`,
		},
	}

	answer, err := Event.Execute(msg)
	assert.NoError(t, err)

	expectedText := "I found next pull-requests:\nPull-request #1 [repository: test-repo]\n\nNext pull-requests cannot be merged yet:\nhttps://bitbucket.org/john/test-repo/1 - One of the required reviewers (@TESTSLACKID) was not found in the reviewers list. \n"
	assert.Equal(t, expectedText, answer.Text)
}