package client

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/sharovik/devbot/internal/log"
	"io/ioutil"
	"net/http"
	"strings"
)

//BaseHttpClientInterface base interface for all http clients
type BaseHttpClientInterface interface {
	//Configuration methods
	SetOauthToken(token string)
	SetBaseUrl(baseUrl string)
	BasicAuth(username string, password string) string
	GetClientID() string
	GetClientSecret() string
	GetOAuthToken() string

	//Http methods
	Request(string, string, interface{}, map[string]string) ([]byte, int, error)
	Post(string, interface{}, map[string]string) ([]byte, int, error)
	Get(string, map[string]string) ([]byte, int, error)
	Put(string, interface{}, map[string]string) ([]byte, int, error)
}

//HttpClient main http client
type HttpClient struct {
	Client     *http.Client

	//Configuration of client
	OAuthToken string
	BaseURL string
	ClientID string
	ClientSecret string
}

//SetOauthToken method sets the oauth token and retrieves its self
func (client *HttpClient) SetOauthToken(token string) {
	client.OAuthToken = token
}

//GetClientID method retrieves the clientID
func (client HttpClient) GetClientID() string {
	return client.ClientID
}

//GetClientSecret method retrieves the clientSecret
func (client HttpClient) GetClientSecret() string {
	return client.ClientSecret
}

//GetOAuthToken method retrieves the oauth token
func (client HttpClient) GetOAuthToken() string {
	return client.OAuthToken
}

//SetBaseUrl method sets the base url and retrieves its self
func (client *HttpClient) SetBaseUrl(baseUrl string) {
	client.BaseURL = baseUrl
}

func (client HttpClient) BasicAuth(username string, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",username, password)))
}

//Request method for API requests
func (client HttpClient) Request(method string, url string, body interface{}, headers map[string]string) ([]byte, int, error) {

	log.Logger().StartMessage("Http request")

	var (
		resp    *http.Response
		request *http.Request
		err error
	)

	log.Logger().Debug().
		Str("url", url).
		Str("method", method).
		Interface("body", body).
		Msg("Endpoint call")

	switch body.(type) {
	case string:
		request, err = http.NewRequest(method, url, strings.NewReader(fmt.Sprintf("%s", body)))
		if err != nil {
			log.Logger().AddError(err).Msg("Error during the request generation")
			log.Logger().FinishMessage("Http request")
			return nil, 0, err
		}
	default:
		request, err = http.NewRequest(method, url, bytes.NewReader(body.([]byte)))
		if err != nil {
			log.Logger().AddError(err).Msg("Error during the request generation")
			log.Logger().FinishMessage("Http request")
			return nil, 0, err
		}
	}

	request.Header.Set("Content-Type", "application/json")
	for attribute, value := range headers {
		request.Header.Set(attribute, value)
	}

	if client.OAuthToken != "" {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.OAuthToken))
	}

	resp, errorResponse := client.Client.Do(request)

	if resp == nil {
		err = errors.New("Response cannot be null ")
		errMsg := err.Error()
		if errorResponse != nil {
			errMsg = errorResponse.Error()
		}
		log.Logger().AddError(errorResponse).
			Str("response_error", errMsg).
			Msg("Error during response body parse")

		log.Logger().FinishMessage("Http request")
		return nil, 0, err
	}

	defer resp.Body.Close()
	byteResp, errorConversion := ioutil.ReadAll(resp.Body)
	if errorConversion != nil {
		log.Logger().AddError(errorConversion).
			Err(errorConversion).
			Msg("Error during response body parse")
		log.Logger().FinishMessage("Http request")
		return byteResp, 0, errorConversion
	}

	var response []byte
	if string(byteResp) == "" {
		response = []byte(`{}`)
	} else {
		response = byteResp
	}

	log.Logger().FinishMessage("Http request")
	return response, resp.StatusCode, nil
}

//Post method for POST http requests
func (client HttpClient) Post(endpoint string, body interface{}, headers map[string]string) ([]byte, int, error) {
	return client.Request(http.MethodPost, client.generateAPIUrl(endpoint), body, headers)
}

//Put method for PUT http requests
func (client HttpClient) Put(endpoint string, body interface{}, headers map[string]string) ([]byte, int, error) {
	return client.Request(http.MethodPut, client.generateAPIUrl(endpoint), body, headers)
}

//Get method for GET http requests
func (client *HttpClient) Get(endpoint string, query map[string]string) ([]byte, int, error) {
	if client.OAuthToken != "" {
		query["access_token"] = client.OAuthToken
	}

	var queryString = ""
	for fieldName, value := range query {
		if queryString == "" {
			queryString += "?"
		} else {
			queryString += "&"
		}

		queryString += fmt.Sprintf("%s=%s", fieldName, value)
	}

	return client.Request(http.MethodGet, client.generateAPIUrl(endpoint) + fmt.Sprintf("%s", queryString), []byte(``), map[string]string{})
}

func (client HttpClient) generateAPIUrl(endpoint string) string {
	log.Logger().Debug().
		Str("base_url", client.BaseURL).
		Str("endpoint", endpoint).
		Msg("Generate API url")

	return client.BaseURL + endpoint
}
