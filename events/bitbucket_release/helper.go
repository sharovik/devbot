package bitbucket_release

import (
	"fmt"
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/log"
	"regexp"
	"strconv"
	"time"
)

func failedPullRequestsText(failedPullRequests map[string]failedToMerge) string {
	if len(failedPullRequests) == 0 {
		return "There is no pull-requests, which cannot be merged! This is good."
	}

	var text = "Next pull-requests cannot be merged yet:\n"

	for pullRequest, reason := range failedPullRequests {
		text += fmt.Sprintf("%s - %s \n", pullRequest, reason.Reason)
	}

	return text
}

func canBeMergedPullRequestsText(canBeMerged map[string]PullRequest) string {
	if len(canBeMerged) == 0 {
		return "There is no pull-requests, which can be merged."
	}

	var text = "Next pull-requests is about to release:\n"

	for pullRequestURL, pullRequest := range canBeMerged {
		text += fmt.Sprintf("[#%d] %s \n", pullRequest.ID, pullRequestURL)
	}

	return text
}

func checkPullRequests(items []PullRequest) (map[string]PullRequest, map[string][]PullRequest, map[string]failedToMerge) {
	var (
		failedPullRequests         = map[string]failedToMerge{}
		canBeMergedPullRequestList = map[string]PullRequest{}
		canBeMergedByRepository    = map[string][]PullRequest{}
	)
	for _, pullRequest := range items {
		cleanPullRequestURL := fmt.Sprintf("https://bitbucket.org/%s/%s/%d", pullRequest.Workspace, pullRequest.RepositorySlug, pullRequest.ID)
		info, err := container.C.BibBucketClient.PullRequestInfo(pullRequest.Workspace, pullRequest.RepositorySlug, pullRequest.ID)
		pullRequest.Description = info.Description

		if err != nil {
			failedPullRequests[cleanPullRequestURL] = failedToMerge{
				Reason: err.Error(),
				Info:   info,
				Error:  err,
			}

			continue
		}

		if !isPullRequestAlreadyMerged(info) {
			failedPullRequests[cleanPullRequestURL] = failedToMerge{
				Reason: fmt.Sprintf("The state should be %s, instead of it %s received.", pullRequestStateOpen, info.State),
				Info:   info,
				Error:  nil,
			}

			continue
		}

		isRequiredReviewersExistsInPullRequest, reason := checkIfRequiredReviewersExists(info)
		if !isRequiredReviewersExistsInPullRequest {
			failedPullRequests[cleanPullRequestURL] = reason
			continue
		}

		isPullRequestApprovedByReviewers := checkIfOneOfRequiredReviewersApprovedPullRequest(info)
		if !isPullRequestApprovedByReviewers {
			failedPullRequests[cleanPullRequestURL] = failedToMerge{
				Reason: "The pull-request should be approved by one of the required reviewers.",
				Info:   info,
				Error:  nil,
			}

			continue
		}

		canBeMergedPullRequestList[cleanPullRequestURL] = pullRequest
		canBeMergedByRepository[pullRequest.RepositorySlug] = append(canBeMergedByRepository[pullRequest.RepositorySlug], pullRequest)
	}

	return canBeMergedPullRequestList, canBeMergedByRepository, failedPullRequests
}

func startRelease(canBeMergedPullRequestList map[string]PullRequest, canBeMergedByRepository map[string][]PullRequest) (string, error) {
	log.Logger().StartMessage("Release of received pull-requests")
	releaseText := ""
	//In case when we have only one pull-request we will merge it straight to the main branch
	if len(canBeMergedPullRequestList) == 1 {
		log.Logger().Debug().Msg("There is only 1 received pull-request. Trying to merge it.")
		releaseText = fmt.Sprintf("We have only one pull-request, so I will try to merge it stright to the main branch.")
		newText, err := mergePullRequests(canBeMergedPullRequestList)
		releaseText += fmt.Sprintf("\n%s", newText)
		if err != nil {
			log.Logger().AddError(err).Msg("Failed to merge the pull-request")
			log.Logger().FinishMessage("Release of received pull-requests")
			return releaseText, err
		}

		log.Logger().FinishMessage("Release of received pull-requests")
		return releaseText, nil
	}

	//Here we take sorted by repository pull-requests and trying to merge them into main or release branch.
	//We go in for loop into each repository and check how many pull-requests do we have there.
	//If only one, then we merge it into main branch, otherwise we create release branch for selected repository,
	//switch direction of the pull-requests to that release branch and merge all of them.
	for repository, pullRequests := range canBeMergedByRepository {
		//Well, in that case we have only one pull-request so we merge it into main branch
		if len(pullRequests) == 1 {
			log.Logger().Debug().Str("repository", repository).Msg("Only one pull-request received for selected repository")

			releaseText = fmt.Sprintf("There is only one pull-request for selected repository `%s`.", repository)
			newText, err := mergePullRequests(canBeMergedPullRequestList)
			releaseText += fmt.Sprintf("\n%s", newText)
			if err != nil {
				log.Logger().AddError(err).Msg("Received error during pull-request merge")
			}

			continue
		}

		//This is for multiple pull-requests links
		var (
			repositories        = map[string]string{}
			releaseBranchName = fmt.Sprintf("release/%s", time.Now().Format("2006.01.02"))
		)

		//In that case we have multiple pull-requests for that repository, so we have to create a release branch
		for _, pullRequest := range pullRequests {
			//If we don't have any created release branch for this repository the we need to create it
			if repositories[repository] == "" {
				_, err := container.C.BibBucketClient.FindOrCreateBranch(pullRequest.Workspace, pullRequest.RepositorySlug, releaseBranchName)
				if err != nil {
					releaseText += fmt.Sprintf("The release-branch for repository %s cannot be created, because of `%s`", repository, err)
					log.Logger().AddError(err).Msg("Received an error during the release branch creation")
					break
				}

				repositories[repository] = releaseBranchName
			}

			//@todo: need to switch the direction of selected pull-request to release branch and merge it.
		}
	}

	log.Logger().FinishMessage("Release of received pull-requests")
	return releaseText, nil
}

func mergePullRequests(pullRequests map[string]PullRequest) (string, error) {
	var (
		releaseText string
		repository = ""
	)

	for _, pullRequest := range pullRequests {
		if repository == "" {
			repository = pullRequest.RepositorySlug
		}

		_, err := container.C.BibBucketClient.MergePullRequest(pullRequest.Workspace, pullRequest.RepositorySlug, pullRequest.ID, pullRequest.Description)
		if err != nil {
			releaseText += fmt.Sprintf("I cannot merge the pull-request #%d because of error `%s`", pullRequest.ID, err.Error())
			return releaseText, err
		}
	}

	releaseText += fmt.Sprintf("\nI merged all pull-requests for repository %s :)", repository)

	return releaseText, nil
}

func checkIfRequiredReviewersExists(info dto.BitBucketPullRequestInfoResponse) (bool, failedToMerge) {
	requiredReviewers := container.C.Config.BitBucketConfig.RequiredReviewers

	var (
		existsInReviewers = false
		result            = failedToMerge{}
	)

	for _, reviewerUuID := range requiredReviewers {
		if existsInReviewers == true {
			break
		}

		for _, user := range info.Participants {
			if reviewerUuID.UUID == user.User.UUID {
				existsInReviewers = true
			}
		}

		if existsInReviewers == false {
			result = failedToMerge{
				Reason: fmt.Sprintf("One of the required reviewers (<@%s>) was not found in the reviewers list.", reviewerUuID.SlackUID),
				Info:   info,
				Error:  nil,
			}
		}
	}

	if existsInReviewers != true {
		return false, result
	}

	return true, failedToMerge{}
}

func checkIfOneOfRequiredReviewersApprovedPullRequest(info dto.BitBucketPullRequestInfoResponse) bool {
	requiredReviewers := container.C.Config.BitBucketConfig.RequiredReviewers

	var isPullRequestApprovedByReviewers = false
	for _, reviewerUuID := range requiredReviewers {

		for _, user := range info.Participants {
			if reviewerUuID.UUID == user.User.UUID && user.Approved {
				isPullRequestApprovedByReviewers = true
			}
		}
	}

	return isPullRequestApprovedByReviewers
}

func isPullRequestAlreadyMerged(info dto.BitBucketPullRequestInfoResponse) bool {
	if info.State == pullRequestStateOpen {
		return true
	}

	return false
}

func receivedPullRequestsText(foundPullRequests ReceivedPullRequests) string {

	if len(foundPullRequests.Items) == 0 {
		return noPullRequestStringAnswer
	}

	var pullRequestsString = pullRequestStringAnswer
	for _, item := range foundPullRequests.Items {
		pullRequestsString = pullRequestsString + fmt.Sprintf("Pull-request #%d [repository: %s]\n", item.ID, item.RepositorySlug)
	}

	return pullRequestsString
}

func findAllPullRequestsInText(regex string, subject string) ReceivedPullRequests {
	re, err := regexp.Compile(regex)

	if err != nil {
		log.Logger().AddError(err).Msg("Error during the Find Matches operation")
		return ReceivedPullRequests{}
	}

	matches := re.FindAllStringSubmatch(subject, -1)
	result := ReceivedPullRequests{}

	if len(matches) == 0 {
		return result
	}

	for _, id := range matches {
		if id[1] != "" {
			item := PullRequest{}
			item.Workspace = id[1]
			item.RepositorySlug = id[2]
			item.ID, err = strconv.ParseInt(id[3], 10, 64)
			if err != nil {
				log.Logger().AddError(err).
					Interface("matches", matches).
					Msg("Error during pull-request ID parsing")
				return ReceivedPullRequests{}
			}

			result.Items = append(result.Items, item)
		}
	}

	return result
}
