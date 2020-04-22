package mock

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

type MockedHttpClient struct {
	BaseURL string

	RequestMethodResponse []byte
	RequestMethodResponseStatusCode int
	RequestMethodError error
}

//SetOauthToken method sets the oauth token and retrieves its self
func (client *MockedHttpClient) SetOauthToken(token string) {
	
}

//GetClientID method retrieves the clientID
func (client MockedHttpClient) GetClientID() string {
	return "client.ClientID"
}

//GetClientSecret method retrieves the clientSecret
func (client MockedHttpClient) GetClientSecret() string {
	return "client.ClientSecret"
}

//GetOAuthToken method retrieves the oauth token
func (client MockedHttpClient) GetOAuthToken() string {
	return "client.OAuthToken"
}

//SetBaseUrl method sets the base url and retrieves its self
func (client *MockedHttpClient) SetBaseUrl(baseUrl string) {

}

func (client MockedHttpClient) BasicAuth(username string, password string) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
}

//Request method for API requests
//
//This method accepts parameters:
//method - the method of request. Ex: POST, GET, PUT, DELETE and etc
//endpoint - endpoint to which we should do a request
//body - it's a request body. Accepted types of body: string, url.Values(for form_data requests), byte
//headers - request headers
func (client MockedHttpClient) Request(method string, endpoint string, body interface{}, headers map[string]string) ([]byte, int, error) {
	return client.RequestMethodResponse, client.RequestMethodResponseStatusCode, client.RequestMethodError
}

//Post method for POST http requests
func (client MockedHttpClient) Post(endpoint string, body interface{}, headers map[string]string) ([]byte, int, error) {
	return client.Request(http.MethodPost, client.generateAPIUrl(endpoint), body, headers)
}

//Put method for PUT http requests
func (client MockedHttpClient) Put(endpoint string, body interface{}, headers map[string]string) ([]byte, int, error) {
	return client.Request(http.MethodPut, client.generateAPIUrl(endpoint), body, headers)
}

//Get method for GET http requests
func (client *MockedHttpClient) Get(endpoint string, query map[string]string) ([]byte, int, error) {
	var queryString = ""
	for fieldName, value := range query {
		if queryString == "" {
			queryString += "?"
		} else {
			queryString += "&"
		}

		queryString += fmt.Sprintf("%s=%s", fieldName, value)
	}

	return client.Request(http.MethodGet, client.generateAPIUrl(endpoint)+fmt.Sprintf("%s", queryString), []byte(``), map[string]string{})
}

func (client MockedHttpClient) generateAPIUrl(endpoint string) string {
	return client.BaseURL + endpoint
}
