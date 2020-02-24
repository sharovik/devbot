package bitbucket_release

import (
	"fmt"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/log"
	"regexp"
)

//EventName the name of the event
const (
	EventName = "bitbucket_release"
	pullRequestsRegex = `(?m)https:\/\/bitbucket.org\/(?P<workspace>.+)\/(?P<repository_slug>.+)\/pull-requests\/(?P<pull_request_id>\d+)`

	pullRequestStringAnswer = "I found next pull-requests:\n"
	noPullRequestStringAnswer = `I can't find any pull-request in your message`
)

//BitBucketReleaseEvent event of bitbucket release
type BitBucketReleaseEvent struct {
	EventName string
}

//Event - object which is ready to use
var Event = BitBucketReleaseEvent{
	EventName: EventName,
}

func (BitBucketReleaseEvent) Execute(message dto.SlackRequestChatPostMessage) (dto.SlackRequestChatPostMessage, error) {
	var answer = message

	foundPullRequests := findAllPullRequests(pullRequestsRegex, answer.OriginalMessage.Text)

	answer.Text = prepareReceivedPullRequests(foundPullRequests)
	return answer, nil
}

func prepareReceivedPullRequests(foundPullRequests map[string]map[string]string) string {

	if len(foundPullRequests) == 0 {
		return noPullRequestStringAnswer
	}

	var pullRequestsString = pullRequestStringAnswer
	for _, item := range foundPullRequests {
		pullRequestsString = pullRequestsString + fmt.Sprintf("Pull-request #%s [repository: %s]\n", item["pull_request_id"], item["repository_slug"])
	}

	return pullRequestsString
}

func findAllPullRequests(regex string, subject string) map[string]map[string]string {
	re, err := regexp.Compile(regex)

	if err != nil {
		log.Logger().AddError(err).Msg("Error during the Find Matches operation")
		return map[string]map[string]string{}
	}

	matches := re.FindAllStringSubmatch(subject, -1)
	result := map[string]map[string]string{}

	if len(matches) == 0 {
		return result
	}

	for index, id := range matches {
		if id[1] != "" {
			item := map[string]string{}
			item["workspace"] = id[1]
			item["repository_slug"] = id[2]
			item["pull_request_id"] = id[3]
			result[fmt.Sprintf("%d", index)] = item
		}
	}

	return result
}
