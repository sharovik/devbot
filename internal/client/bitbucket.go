package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/log"
	"net/http"
	"net/url"
	"time"
)

type BitBucketClient struct {
	client           BaseHttpClientInterface
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
	DefaultBitBucketBaseAPIUrl     = "https://api.bitbucket.org/2.0"
	DefaultBitBucketAccessTokenUrl = "https://bitbucket.org/site/oauth2"

	DefaultBitBucketMainBranch = "master"

	ErrorBranchExists = "BRANCH_ALREADY_EXISTS"
)

type GitClientInterface interface {
	Init(client BaseHttpClientInterface)
	PullRequestInfo(workspace string, repositorySlug string, pullRequestID int64) (dto.BitBucketPullRequestInfoResponse, error)
	MergePullRequest(workspace string, repositorySlug string, pullRequestID int64, description string) (dto.BitBucketPullRequestInfoResponse, error)
	FindOrCreateBranch(workspace string, repositorySlug string, branchName string) (string, error)
}

func (b *BitBucketClient) Init(client BaseHttpClientInterface) {
	b.client = client
}

func (b *BitBucketClient) isTokenInvalid() bool {
	if b.OauthToken == "" {
		log.Logger().Warn().Str("error", "oauth_empty").Msg("Invalid token")
		return false
	}

	if b.OauthTokenExpire.Unix() <= time.Now().Unix() {
		log.Logger().Warn().Time("time", b.OauthTokenExpire).Str("error", "expired").Msg("Invalid token")
		return true
	}

	return false
}

func (b *BitBucketClient) beforeRequest() error {
	log.Logger().StartMessage("Before BitBucket request")
	if !b.isTokenInvalid() {
		if err := b.loadAuthToken(); err != nil {
			log.Logger().FinishMessage("Before BitBucket request")
			return err
		}
	}

	log.Logger().FinishMessage("Before BitBucket request")
	return nil
}

func (b *BitBucketClient) loadAuthToken() error {
	log.Logger().StartMessage("Loading OAuth token")

	b.client.SetBaseUrl(DefaultBitBucketAccessTokenUrl)
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

	log.Logger().Debug().Interface("response", responseObject).Msg("Successfully updated token data")
	log.Logger().FinishMessage("Loading OAuth token")
	return nil
}

func (b *BitBucketClient) FindOrCreateBranch(workspace string, repositorySlug string, branchName string) (string, error) {
	log.Logger().StartMessage("Get pull-request status")
	if err := b.beforeRequest(); err != nil {
		log.Logger().FinishMessage("Get pull-request status")
		return "", err
	}

	b.client.SetBaseUrl(DefaultBitBucketBaseAPIUrl)

	endpoint := fmt.Sprintf("/repositories/%s/%s/refs/branches/%s", workspace, repositorySlug, branchName)
	_, statusCode, err := b.client.Get(endpoint, map[string]string{})
	if err != nil {
		log.Logger().FinishMessage("Get pull-request status")
		return "", err
	}

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
			return "", err
		}

		response, statusCode, err := b.client.Post(endpoint, byteRequest, map[string]string{})
		if statusCode == http.StatusBadRequest {
			responseObject := dto.BitBucketErrorResponseBranchCreate{}
			if err := json.Unmarshal(response, &responseObject); err != nil {
				return "", err
			}

			if responseObject.Data.Key == ErrorBranchExists {
				log.Logger().Info().
					Str("branch", branchName).
					Int("status_code", statusCode).
					RawJSON("response", response).
					Msg("Branch already exists")
				return branchName, nil
			}
		}

		if statusCode != http.StatusCreated {
			log.Logger().Warn().
				Str("branch", branchName).
				Int("status_code", statusCode).
				RawJSON("response", response).
				Msg("Bad status code received")

			return "", errors.New("Wrong status code received during the branch creation. See the logs for more information. ")
		}
	}

	log.Logger().FinishMessage("Get pull-request status")
	return branchName, nil
}

func (b *BitBucketClient) PullRequestInfo(workspace string, repositorySlug string, pullRequestID int64) (dto.BitBucketPullRequestInfoResponse, error) {
	log.Logger().StartMessage("Get pull-request status")
	if err := b.beforeRequest(); err != nil {
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	b.client.SetBaseUrl(DefaultBitBucketBaseAPIUrl)
	endpoint := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d", workspace, repositorySlug, pullRequestID)
	response, _, err := b.client.Get(endpoint, map[string]string{})

	if err != nil {
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	var responseObject dto.BitBucketPullRequestInfoResponse
	err = json.Unmarshal(response, &responseObject)
	if err != nil {
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	log.Logger().FinishMessage("Get pull-request status")
	return responseObject, nil
}

func (b *BitBucketClient) MergePullRequest(workspace string, repositorySlug string, pullRequestID int64, description string) (dto.BitBucketPullRequestInfoResponse, error) {
	log.Logger().StartMessage("Merge pull-request")
	if err := b.beforeRequest(); err != nil {
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	formData := map[string]string{
		"merge_strategy":      "squash",
		"message":             description,
		"close_source_branch": "1",
	}

	byteString, err := json.Marshal(formData)
	if err != nil {
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	b.client.SetBaseUrl(DefaultBitBucketBaseAPIUrl)
	endpoint := fmt.Sprintf("/repositories/%s/%s/pullrequests/%d/merge", workspace, repositorySlug, pullRequestID)

	var dtoResponse = dto.BitBucketPullRequestInfoResponse{}
	response, statusCode, err := b.client.Post(endpoint, byteString, map[string]string{})
	if err != nil {
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	if err := json.Unmarshal(response, &dtoResponse); err != nil {
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	if statusCode == http.StatusBadRequest {
		return dto.BitBucketPullRequestInfoResponse{}, errors.New(dtoResponse.Error.Message)
	}

	if statusCode == http.StatusUnauthorized {
		return dto.BitBucketPullRequestInfoResponse{}, errors.New("Hm, I received unauthorized response. It means that someone changes the access for me. ")
	}

	if statusCode == http.StatusNotFound {
		return dto.BitBucketPullRequestInfoResponse{}, errors.New("Selected pull-request was not found :( ")
	}

	var responseObject dto.BitBucketPullRequestInfoResponse
	err = json.Unmarshal(response, &responseObject)
	if err != nil {
		return dto.BitBucketPullRequestInfoResponse{}, err
	}

	log.Logger().FinishMessage("Merge pull-request")
	return responseObject, nil
}
