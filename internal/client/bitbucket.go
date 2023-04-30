package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/log"
)

// BitBucketClient the bitbucket client struct
type BitBucketClient struct {
	client           BaseHTTPClientInterface
	OauthToken       string
	OauthTokenExpire time.Time
	RefreshToken     string
}

type responseAccessToken struct {
	AccessToken  string `json:"access_token"`
	Scopes       string `json:"scopes"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

const (
	//DefaultBitBucketBaseAPIUrl the base url
	DefaultBitBucketBaseAPIUrl = "https://api.bitbucket.org/2.0"

	//DefaultBitBucketAccessTokenURL the access token endpoint which will be used for token generation
	DefaultBitBucketAccessTokenURL = "https://bitbucket.org/site/oauth2"

	//DefaultBitBucketMainBranch the default main branch
	DefaultBitBucketMainBranch = "master"

	//ErrorBranchExists error message for "branch exists error"
	ErrorBranchExists = "BRANCH_ALREADY_EXISTS"

	//StrategySquash the squash strategy which can be used during the merge
	StrategySquash = "squash"

	//StrategyMerge the default merge strategy which can be used during the merge
	StrategyMerge = "merge_commit"

	//ErrorMsgNoAccess error message response of bot, once he got a bad status code from the API
	ErrorMsgNoAccess = "received unauthorized response. Looks like I'm not permitted to do any actions to that repository"
)

// Init initialise the client
func (b *BitBucketClient) Init(client BaseHTTPClientInterface) {
	b.client = client
}

func (b *BitBucketClient) isTokenInvalid() bool {
	if b.OauthToken == "" {
		log.Logger().Warn().Str("error", "oauth_empty").Msg("Invalid token")
		return true
	}

	if b.OauthTokenExpire.Unix() <= time.Now().Unix() {
		log.Logger().Warn().Time("time", b.OauthTokenExpire).Str("error", "expired").Msg("Invalid token")
		return true
	}

	return false
}

func (b *BitBucketClient) beforeRequest() error {
	log.Logger().StartMessage("Before BitBucket request")
	if b.isTokenInvalid() {
		log.Logger().Warn().Msg("Trying to regenerate the token")
		if err := b.loadAuthToken(); err != nil {
			log.Logger().AddError(err).Msg("Failed to generate new token")
			log.Logger().FinishMessage("Before BitBucket request")
			return err
		}
	}

	log.Logger().FinishMessage("Before BitBucket request")
	return nil
}

// GetAuthToken method retrieves the access token which can be used for custom needs outside of the internal services
func (b *BitBucketClient) GetAuthToken() (string, error) {
	if err := b.beforeRequest(); err != nil {
		log.Logger().FinishMessage("Can't load access token")
		return "", err
	}

	return b.OauthToken, nil
}

func (b *BitBucketClient) loadAuthToken() error {
	log.Logger().StartMessage("Loading OAuth token")

	b.client.SetBaseURL(DefaultBitBucketAccessTokenURL)
	b.client.SetOauthToken("") //this will cleanup the token in the client and will generate new basic auth for specific client_id and client_secret

	formData := url.Values{}
	formData.Add("grant_type", "client_credentials")

	response, _, err := b.client.Post("/access_token", formData.Encode(), map[string]string{
		"Content-Type": "application/x-www-form-urlencoded;",
		"Authorization": fmt.Sprintf("Basic %s", b.client.BasicAuth(
			b.client.GetClientID(),
			b.client.GetClientSecret(),
		)),
	})
	if err != nil {
		return err
	}

	var responseObject responseAccessToken
	err = json.Unmarshal(response, &responseObject)
	if err != nil {
		return err
	}

	b.RefreshToken = responseObject.RefreshToken
	b.OauthTokenExpire = time.Now().Add(time.Second * time.Duration(responseObject.ExpiresIn))
	b.OauthToken = responseObject.AccessToken
	b.client.SetOauthToken(responseObject.AccessToken)

	log.Logger().FinishMessage("Loading OAuth token")
	return nil
}

// CreateBranch creates the branch in API
func (b *BitBucketClient) CreateBranch(workspace string, repositorySlug string, branchName string) (dto.BitBucketResponseBranchCreate, error) {
	log.Logger().StartMessage("Create branch")
	if err := b.beforeRequest(); err != nil {
		log.Logger().FinishMessage("Create branch")
		return dto.BitBucketResponseBranchCreate{}, err
	}

	b.client.SetBaseURL(DefaultBitBucketBaseAPIUrl)

	endpoint := fmt.Sprintf("/repositories/%s/%s/refs/branches/%s", workspace, repositorySlug, branchName)
	response, statusCode, err := b.client.Get(endpoint, map[string]string{})
	if err != nil {
		log.Logger().FinishMessage("Create branch")
		return dto.BitBucketResponseBranchCreate{}, err
	}

	responseObject := dto.BitBucketResponseBranchCreate{}
	if statusCode == http.StatusNotFound {
		log.Logger().Info().Str("branch", branchName).Msg("Release branch wasn't found. Trying to create it.")
		request := dto.BitBucketRequestBranchCreate{
			Name: branchName,
			Target: dto.BitBucketBranchTarget{
				Hash: DefaultBitBucketMainBranch,
			},
		}

		byteRequest, err := json.Marshal(request)
		if err != nil {
			log.Logger().FinishMessage("Create branch")
			return dto.BitBucketResponseBranchCreate{}, err
		}

		endpoint := fmt.Sprintf("/repositories/%s/%s/refs/branches", workspace, repositorySlug)
		response, statusCode, err := b.client.Post(endpoint, byteRequest, map[string]string{})
		if err != nil {
			log.Logger().AddError(err).
				Msg("Failed to trigger request")
			log.Logger().FinishMessage("Create branch")

			return dto.BitBucketResponseBranchCreate{}, err
		}

		if err := json.Unmarshal(response, &responseObject); err != nil {
			log.Logger().Info().
				Str("branch", branchName).
				Int("status_code", statusCode).
				Str("response", string(response)).
				Msg("Failed to unmarshal response")
			log.Logger().FinishMessage("Create branch")
			return dto.BitBucketResponseBranchCreate{}, err
		}

		if statusCode == http.StatusBadRequest {
			if responseObject.Data.Key == ErrorBranchExists {
				log.Logger().Info().
					Str("branch", branchName).
					Int("status_code", statusCode).
					RawJSON("response", response).
					Msg("Branch already exists")
				log.Logger().FinishMessage("Create branch")
				return responseObject, nil
			}
		}

		if statusCode != http.StatusCreated {
			log.Logger().Warn().
				Str("branch", branchName).
				Int("status_code", statusCode).
				Interface("response", responseObject).
				Msg("Bad status code received")

			log.Logger().FinishMessage("Create branch")
			return dto.BitBucketResponseBranchCreate{}, errors.New("wrong status code received during the branch creation. See the logs for more information. ")
		}

		log.Logger().Info().
			Str("branch", branchName).
			Int("status_code", statusCode).
			RawJSON("response", response).Msg("Create branch result")
		log.Logger().FinishMessage("Create branch")
		return responseObject, nil
	}

	err = json.Unmarshal(response, &responseObject)
	if err != nil {
		log.Logger().AddError(err).Msg("Error during response unmarshal")
		log.Logger().FinishMessage("Create branch")
		return dto.BitBucketResponseBranchCreate{}, err
	}

	log.Logger().FinishMessage("Create branch")
	return responseObject, nil
}

// GetBranch creates the branch in API
func (b *BitBucketClient) GetBranch(workspace string, repositorySlug string, branchName string) (dto.BitBucketResponseBranchCreate, error) {
	log.Logger().StartMessage("Get branch")
	if err := b.beforeRequest(); err != nil {
		log.Logger().FinishMessage("Get branch")
		return dto.BitBucketResponseBranchCreate{}, err
	}

	b.client.SetBaseURL(DefaultBitBucketBaseAPIUrl)

	endpoint := fmt.Sprintf("/repositories/%s/%s/refs/branches/%s", workspace, repositorySlug, branchName)
	response, statusCode, err := b.client.Get(endpoint, map[string]string{})
	if err != nil {
		log.Logger().FinishMessage("Get branch")
		return dto.BitBucketResponseBranchCreate{}, err
	}

	if statusCode == http.StatusNotFound {
		log.Logger().FinishMessage("Get branch")
		return dto.BitBucketResponseBranchCreate{}, errors.New("this branch doesn't exist. ")
	}

	if statusCode == http.StatusForbidden {
		log.Logger().FinishMessage("Get branch")
		return dto.BitBucketResponseBranchCreate{}, errors.New("action is not permitted. ")
	}

	responseObject := dto.BitBucketResponseBranchCreate{}
	err = json.Unmarshal(response, &responseObject)
	if err != nil {
		log.Logger().AddError(err).Msg("Error during response unmarshal")
		log.Logger().FinishMessage("Get branch")
		return dto.BitBucketResponseBranchCreate{}, err
	}

	log.Logger().FinishMessage("Get branch")
	return responseObject, nil
}

// PullRequestInfo gets the pull-requests information
func (b *BitBucketClient) PullRequestInfo(workspace string, repositorySlug string, pullRequestID int64) (dto.BitBucketPullRequestInfoResponse, error) {
	log.Logger().StartMessage("Get pull-request status")
	if err := b.beforeRequest(); err != nil {
		log.Logger().FinishMessage("Get pull-request status")
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	b.client.SetBaseURL(DefaultBitBucketBaseAPIUrl)
	endpoint := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d", workspace, repositorySlug, pullRequestID)
	response, _, err := b.client.Get(endpoint, map[string]string{})

	if err != nil {
		log.Logger().FinishMessage("Get pull-request status")
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	var responseObject dto.BitBucketPullRequestInfoResponse
	err = json.Unmarshal(response, &responseObject)
	if err != nil {
		log.Logger().FinishMessage("Get pull-request status")
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	log.Logger().FinishMessage("Get pull-request status")
	return responseObject, nil
}

// MergePullRequest merge the selected pull-request
func (b *BitBucketClient) MergePullRequest(workspace string, repositorySlug string, pullRequestID int64, description string, strategy string) (dto.BitBucketPullRequestInfoResponse, error) {
	log.Logger().StartMessage("Merge pull-request")
	if err := b.beforeRequest(); err != nil {
		log.Logger().FinishMessage("Merge pull-request")
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	formData := map[string]string{
		"merge_strategy":      strategy,
		"message":             description,
		"close_source_branch": "1",
	}

	byteString, err := json.Marshal(formData)
	if err != nil {
		log.Logger().FinishMessage("Merge pull-request")
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	b.client.SetBaseURL(DefaultBitBucketBaseAPIUrl)
	endpoint := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/merge?async=false", workspace, repositorySlug, pullRequestID)

	var dtoResponse = dto.BitBucketPullRequestInfoResponse{}
	response, statusCode, err := b.client.Post(endpoint, byteString, map[string]string{})
	if err != nil {
		log.Logger().FinishMessage("Merge pull-request")
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	//In that case the bitbucket accepts our request and the Pull-request will be merged but in async way(even if we specified async=false)
	if statusCode == http.StatusAccepted {
		log.Logger().
			Debug().
			RawJSON("response", response).
			Int("status_code", statusCode).
			Msg("The response with poll link received.")
		log.Logger().FinishMessage("Merge pull-request")
		return dtoResponse, nil
	}

	if err := json.Unmarshal(response, &dtoResponse); err != nil {
		log.Logger().
			AddError(err).
			RawJSON("response", response).
			Int("status_code", statusCode).
			Msg("Error during the request unmarshal.")
		log.Logger().FinishMessage("Merge pull-request")
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	if statusCode == http.StatusBadRequest {
		log.Logger().FinishMessage("Merge pull-request")
		return dto.BitBucketPullRequestInfoResponse{}, fmt.Errorf("bitbucket response with the error: %s", dtoResponse.Error.Message)
	}

	if statusCode == http.StatusUnauthorized {
		log.Logger().FinishMessage("Merge pull-request")
		return dto.BitBucketPullRequestInfoResponse{}, fmt.Errorf(ErrorMsgNoAccess)
	}

	if statusCode == http.StatusNotFound {
		log.Logger().FinishMessage("Merge pull-request")
		return dto.BitBucketPullRequestInfoResponse{}, errors.New("selected pull-request was not found :( ")
	}

	log.Logger().FinishMessage("Merge pull-request")
	return dtoResponse, nil
}

// ChangePullRequestDestination changes the pull-request destination to selected one
func (b *BitBucketClient) ChangePullRequestDestination(workspace string, repositorySlug string, pullRequestID int64, title string, branchName string) (dto.BitBucketPullRequestInfoResponse, error) {
	log.Logger().StartMessage("Change destination")
	if err := b.beforeRequest(); err != nil {
		log.Logger().FinishMessage("Change destination")
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	byteString, err := json.Marshal(dto.BitBucketPullRequestDestinationUpdateRequest{
		Title: title,
		Destination: dto.BitBucketPullRequestDestination{
			Branch: dto.BitBucketPullRequestDestinationBranch{
				Name: branchName,
			},
		},
	})

	if err != nil {
		log.Logger().FinishMessage("Change destination")
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	b.client.SetBaseURL(DefaultBitBucketBaseAPIUrl)
	endpoint := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d", workspace, repositorySlug, pullRequestID)

	var dtoResponse = dto.BitBucketPullRequestInfoResponse{}
	response, statusCode, err := b.client.Put(endpoint, byteString, map[string]string{})
	if err != nil {
		log.Logger().FinishMessage("Change destination")
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	if err := json.Unmarshal(response, &dtoResponse); err != nil {
		log.Logger().FinishMessage("Change destination")
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	if statusCode == http.StatusBadRequest {
		log.Logger().FinishMessage("Change destination")
		return dto.BitBucketPullRequestInfoResponse{}, errors.New(dtoResponse.Error.Message)
	}

	if statusCode == http.StatusUnauthorized {
		log.Logger().FinishMessage("Change destination")
		return dto.BitBucketPullRequestInfoResponse{}, errors.New(ErrorMsgNoAccess)
	}

	if statusCode == http.StatusNotFound {
		log.Logger().FinishMessage("Change destination")
		return dto.BitBucketPullRequestInfoResponse{}, errors.New("selected pull-request was not found :( ")
	}

	var responseObject dto.BitBucketPullRequestInfoResponse
	err = json.Unmarshal(response, &responseObject)
	if err != nil {
		log.Logger().FinishMessage("Change destination")
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	log.Logger().FinishMessage("Change destination")
	return responseObject, nil
}

// CreatePullRequest creates the pull-request
func (b *BitBucketClient) CreatePullRequest(workspace string, repositorySlug string, request dto.BitBucketRequestPullRequestCreate) (dto.BitBucketPullRequestInfoResponse, error) {
	log.Logger().StartMessage("Create pull-request")
	if err := b.beforeRequest(); err != nil {
		log.Logger().FinishMessage("Create pull-request")
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	byteString, err := json.Marshal(request)
	if err != nil {
		log.Logger().FinishMessage("Create pull-request")
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	b.client.SetBaseURL(DefaultBitBucketBaseAPIUrl)
	endpoint := fmt.Sprintf("/repositories/%s/%s/pullrequests", workspace, repositorySlug)

	response, statusCode, err := b.client.Post(endpoint, byteString, map[string]string{})
	if err != nil {
		log.Logger().AddError(err).Msg("Error during the request")
		log.Logger().FinishMessage("Create pull-request")
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	var dtoResponse = dto.BitBucketPullRequestInfoResponse{}
	if err := json.Unmarshal(response, &dtoResponse); err != nil {
		log.Logger().AddError(err).Msg("Error during the unmarshal")
		log.Logger().FinishMessage("Create pull-request")
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	if statusCode == http.StatusBadRequest {
		log.Logger().Warn().Int("status_code", statusCode).Str("response", string(response)).Msg("Bad request status code")
		log.Logger().FinishMessage("Create pull-request")
		return dto.BitBucketPullRequestInfoResponse{}, errors.New(dtoResponse.Error.Message)
	}

	if statusCode == http.StatusUnauthorized {
		log.Logger().Warn().Int("status_code", statusCode).Str("response", string(response)).Msg("Unauthorized status code")
		log.Logger().FinishMessage("Create pull-request")
		return dto.BitBucketPullRequestInfoResponse{}, errors.New(ErrorMsgNoAccess)
	}

	if statusCode == http.StatusNotFound {
		log.Logger().Warn().Int("status_code", statusCode).Str("response", string(response)).Msg("Not found status code")
		log.Logger().FinishMessage("Create pull-request")
		return dto.BitBucketPullRequestInfoResponse{}, errors.New("endpoint or selected branch was not found :( ")
	}

	log.Logger().FinishMessage("Create pull-request")
	return dtoResponse, nil
}

// RunPipeline runs the selected custom pipeline
func (b *BitBucketClient) RunPipeline(workspace string, repositorySlug string, request dto.BitBucketRequestRunPipeline) (dto.BitBucketResponseRunPipeline, error) {
	log.Logger().StartMessage("Run pipeline")
	if err := b.beforeRequest(); err != nil {
		log.Logger().FinishMessage("Run pipeline")
		return dto.BitBucketResponseRunPipeline{}, err
	}

	if len(request.Variables) == 0 {
		request.Variables = []dto.Variable{}
	}

	byteString, err := json.Marshal(request)
	if err != nil {
		log.Logger().FinishMessage("Run pipeline")
		return dto.BitBucketResponseRunPipeline{}, err
	}

	b.client.SetBaseURL(DefaultBitBucketBaseAPIUrl)
	endpoint := fmt.Sprintf("/repositories/%s/%s/pipelines/", workspace, repositorySlug)
	response, statusCode, err := b.client.Post(endpoint, byteString, map[string]string{})
	if err != nil {
		log.Logger().FinishMessage("Run pipeline")
		return dto.BitBucketResponseRunPipeline{}, err
	}

	fmt.Println(string(response))
	var responseObject dto.BitBucketResponseRunPipeline
	err = json.Unmarshal(response, &responseObject)
	if err != nil {
		log.Logger().FinishMessage("Run pipeline")
		return dto.BitBucketResponseRunPipeline{}, err
	}

	if statusCode == http.StatusBadRequest {
		log.Logger().Warn().Int("status_code", statusCode).Str("response", string(response)).Msg("Bad request status code")
		log.Logger().FinishMessage("Run pipeline")
		return dto.BitBucketResponseRunPipeline{}, fmt.Errorf("%s. (%s)", responseObject.Error.Message, responseObject.Error.Detail)
	}

	if statusCode == http.StatusUnauthorized {
		log.Logger().Warn().Int("status_code", statusCode).Str("response", string(response)).Msg("Unauthorized status code")
		log.Logger().FinishMessage("Run pipeline")
		return dto.BitBucketResponseRunPipeline{}, errors.New(ErrorMsgNoAccess)
	}

	if statusCode == http.StatusNotFound {
		log.Logger().Warn().Int("status_code", statusCode).Str("response", string(response)).Msg("Not found status code")
		log.Logger().FinishMessage("Run pipeline")
		return dto.BitBucketResponseRunPipeline{}, errors.New("endpoint or selected branch was not found :( ")
	}

	log.Logger().FinishMessage("Run pipeline")
	return responseObject, nil
}

// GetDefaultReviewers gets default reviewers
func (b *BitBucketClient) GetDefaultReviewers(workspace string, repositorySlug string) (dto.BitBucketResponseDefaultReviewers, error) {
	if err := b.beforeRequest(); err != nil {
		return dto.BitBucketResponseDefaultReviewers{}, err
	}

	b.client.SetBaseURL(DefaultBitBucketBaseAPIUrl)
	endpoint := fmt.Sprintf("/repositories/%s/%s/default-reviewers", workspace, repositorySlug)

	response, statusCode, err := b.client.Get(endpoint, map[string]string{})
	if err != nil {
		log.Logger().
			AddError(err).
			RawJSON("response", response).
			Int("status_code", statusCode).
			Msg("Error during the request.")
		return dto.BitBucketResponseDefaultReviewers{}, err
	}

	var dtoResponse = dto.BitBucketResponseDefaultReviewers{}
	if err := json.Unmarshal(response, &dtoResponse); err != nil {
		log.Logger().AddError(err).Msg("Error during the unmarshal")
		return dto.BitBucketResponseDefaultReviewers{}, err
	}

	if statusCode == http.StatusUnauthorized {
		log.Logger().Warn().Int("status_code", statusCode).Str("response", string(response)).Msg("Unauthorized status code")
		return dto.BitBucketResponseDefaultReviewers{}, errors.New(ErrorMsgNoAccess)
	}

	if statusCode == http.StatusNotFound {
		log.Logger().Warn().Int("status_code", statusCode).Str("response", string(response)).Msg("Not found status code")
		return dto.BitBucketResponseDefaultReviewers{}, errors.New("endpoint or selected branch was not found :( ")
	}

	return dtoResponse, nil
}
