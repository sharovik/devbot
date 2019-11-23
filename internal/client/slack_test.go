package client_test

import (
	mockhttp "github.com/karupanerura/go-mock-http-response"
	"github.com/sharovik/devbot/internal/client"
	"github.com/sharovik/devbot/internal/container"
	"github.com/stretchr/testify/assert"
	"testing"
)

var slackClient *client.SlackClient

func init() {
	container.C = container.C.Init()
}

func MockSlackResponse(statusCode int, headers map[string]string, body []byte) {
	slackClient = &client.SlackClient{
		Client:     mockhttp.NewResponseMock(statusCode, headers, body).MakeClient(),
		BaseURL:    "__TEST_BASE_URL__",
		OAuthToken: "__TEST_TOKEN__",
	}
}

func TestSlackClient_Post(t *testing.T) {

	t.Run("Simulate the 404 error", func(t *testing.T) {
		MockSlackResponse(404, map[string]string{}, []byte{})

		response, statusCode, err := slackClient.Post("some-endpoint", []byte{})
		assert.Empty(t, response, "Response should be nil")
		assert.Error(t, err, "We should receive error")
		assert.Equal(t, 404, statusCode, "Status code should be equal 404")
	})

	t.Run("Simulate the 399 status code", func(t *testing.T) {
		MockSlackResponse(399, map[string]string{}, []byte{})

		response, statusCode, err := slackClient.Post("some-endpoint", []byte{})
		assert.Empty(t, response, "Response should be nil")
		assert.NoError(t, err, "There should be no error in response")
		assert.Equal(t, 399, statusCode, "Status code should be equal 399")
	})

	t.Run("Simulate the 200 status code", func(t *testing.T) {
		MockSlackResponse(200, map[string]string{}, []byte(`{status: true"}`))

		response, statusCode, err := slackClient.Post("some-endpoint", []byte{})
		assert.Equal(t, []byte(`{status: true"}`), response, "Response should be equal expectation")
		assert.NoError(t, err, "There should be no error in response")
		assert.Equal(t, 200, statusCode, "Status code should be equal 200")
	})
}

func TestSlackClient_Put(t *testing.T) {

	t.Run("Simulate the 404 error", func(t *testing.T) {
		MockSlackResponse(404, map[string]string{}, []byte{})

		response, statusCode, err := slackClient.Put("some-endpoint", []byte{})
		assert.Empty(t, response, "Response should be nil")
		assert.Error(t, err, "We should receive error")
		assert.Equal(t, 404, statusCode, "Status code should be equal 404")
	})

	t.Run("Simulate the 399 status code", func(t *testing.T) {
		MockSlackResponse(399, map[string]string{}, []byte{})

		response, statusCode, err := slackClient.Put("some-endpoint", []byte{})
		assert.Empty(t, response, "Response should be nil")
		assert.NoError(t, err, "There should be no error in response")
		assert.Equal(t, 399, statusCode, "Status code should be equal 399")
	})

	t.Run("Simulate the 200 status code", func(t *testing.T) {
		MockSlackResponse(200, map[string]string{}, []byte(`{status: true"}`))

		response, statusCode, err := slackClient.Put("some-endpoint", []byte{})
		assert.Equal(t, []byte(`{status: true"}`), response, "Response should be equal expectation")
		assert.NoError(t, err, "There should be no error in response")
		assert.Equal(t, 200, statusCode, "Status code should be equal 200")
	})
}

func TestSlackClient_Get(t *testing.T) {

	t.Run("Simulate the 404 error", func(t *testing.T) {
		MockSlackResponse(404, map[string]string{}, []byte{})

		response, statusCode, err := slackClient.Get("some-endpoint")
		assert.Empty(t, response, "Response should be nil")
		assert.Error(t, err, "We should receive error")
		assert.Equal(t, 404, statusCode, "Status code should be equal 404")
	})

	t.Run("Simulate the 399 status code", func(t *testing.T) {
		MockSlackResponse(399, map[string]string{}, []byte{})

		response, statusCode, err := slackClient.Get("some-endpoint")
		assert.Empty(t, response, "Response should be nil")
		assert.NoError(t, err, "There should be no error in response")
		assert.Equal(t, 399, statusCode, "Status code should be equal 399")
	})

	t.Run("Simulate the 200 status code", func(t *testing.T) {
		MockSlackResponse(200, map[string]string{}, []byte(`{status: true"}`))

		response, statusCode, err := slackClient.Get("some-endpoint")
		assert.Equal(t, []byte(`{status: true"}`), response, "Response should be equal expectation")
		assert.NoError(t, err, "There should be no error in response")
		assert.Equal(t, 200, statusCode, "Status code should be equal 200")
	})
}
