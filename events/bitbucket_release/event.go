package bitbucket_release

import (
	"fmt"
	"github.com/sharovik/devbot/internal/dto"
)

//EventName the name of the event
const (
	EventName         = "bitbucket_release"
	pullRequestsRegex = `(?m)https:\/\/bitbucket.org\/(?P<workspace>.+)\/(?P<repository_slug>.+)\/pull-requests\/(?P<pull_request_id>\d+)`

	pullRequestStringAnswer   = "I found next pull-requests:\n"
	noPullRequestStringAnswer = `I can't find any pull-request in your message`

	pullRequestStateOpen = "OPEN"
)

//ReceivedPullRequests struct for pull-requests list
type ReceivedPullRequests struct {
	Items []PullRequest
}

//PullRequest the pull-request item
type PullRequest struct {
	ID             int64
	RepositorySlug string
	Workspace      string
	Description    string
}

//BitBucketReleaseEvent event of BitBucket release
type BitBucketReleaseEvent struct {
	EventName string
}

//Event - object which is ready to use
var Event = BitBucketReleaseEvent{
	EventName: EventName,
}

type failedToMerge struct {
	Reason string
	Info   dto.BitBucketPullRequestInfoResponse
	Error  error
}

func (BitBucketReleaseEvent) Execute(message dto.SlackRequestChatPostMessage) (dto.SlackRequestChatPostMessage, error) {
	var answer = message

	//First we need to find all the pull-requests in received message
	foundPullRequests := findAllPullRequestsInText(pullRequestsRegex, answer.OriginalMessage.Text)

	//We prepare the text, where we define all the pull-requests which we found in the received message
	answer.Text = receivedPullRequestsText(foundPullRequests)

	//Next step is a pull-request statuses check
	canBeMergedPullRequestsList, canBeMergedByRepository, failedPullRequests := checkPullRequests(foundPullRequests.Items)

	//We generate text for pull-requests which cannot be merged
	answer.Text += fmt.Sprintf("\n%s", failedPullRequestsText(failedPullRequests))
	answer.Text += fmt.Sprintf("\n%s", canBeMergedPullRequestsText(canBeMergedPullRequestsList))

	if len(canBeMergedByRepository) > 0 {
		resultText, err := startRelease(canBeMergedPullRequestsList, canBeMergedByRepository)
		if err != nil {
			answer.Text += fmt.Sprintf("I tried to merge and I failed. Here why: %s", err.Error())
			return answer, err
		}

		answer.Text += fmt.Sprintf("\n%s", resultText)
	}

	return answer, nil
}
