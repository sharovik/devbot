package client_test

import (
	"net/http"
	"os"
	"path"
	"runtime"
	"testing"

	mockhttp "github.com/karupanerura/go-mock-http-response"
	"github.com/sharovik/devbot/internal/client"
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/test"
	"github.com/stretchr/testify/assert"
)

var slackClient *client.SlackClient

func init() {
	//We switch pointer to the root directory for control the path from which we need to generate test-data file-paths
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	_ = os.Chdir(dir)

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

func TestSlackClient_GetConversationsList_Bad(t *testing.T) {
	badStatusCases := map[int]string{
		http.StatusBadGateway:          "Bad gateway",
		http.StatusInternalServerError: "Internal error",
		http.StatusNotFound:            "Not found",
		http.StatusForbidden:           "Forbidden",
		http.StatusBadRequest:          "Bad request",
	}

	for status, errorType := range badStatusCases {
		MockSlackResponse(status, map[string]string{}, []byte(`{"ok": false, "error": "`+errorType+`"}`))

		response, statusCode, err := slackClient.GetConversationsList()
		assert.Error(t, err)
		assert.Empty(t, response)
		assert.Equal(t, status, statusCode)
	}
}

func TestSlackClient_GetConversationsList_Ok(t *testing.T) {
	MockSlackResponse(http.StatusOK, map[string]string{}, test.FileToBytes(t, "./test/testdata/slack/conversations.list.ok.json"))

	response, statusCode, err := slackClient.GetConversationsList()
	assert.NoError(t, err)
	assert.NotEmpty(t, response)
	assert.NotEmpty(t, response.Channels)
	assert.Empty(t, response.Error)
	assert.Equal(t, true, response.Ok)
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestSlackClient_GetUsersList_Bad(t *testing.T) {
	badStatusCases := map[int]string{
		http.StatusBadGateway:          "Bad gateway",
		http.StatusInternalServerError: "Internal error",
		http.StatusNotFound:            "Not found",
		http.StatusForbidden:           "Forbidden",
		http.StatusBadRequest:          "Bad request",
	}

	for status, errorType := range badStatusCases {
		MockSlackResponse(status, map[string]string{}, []byte(`{"ok": false, "error": "`+errorType+`"}`))

		response, statusCode, err := slackClient.GetUsersList()
		assert.Error(t, err)
		assert.Empty(t, response)
		assert.Equal(t, status, statusCode)
	}
}

func TestSlackClient_GetUsersList_Ok(t *testing.T) {
	MockSlackResponse(http.StatusOK, map[string]string{}, test.FileToBytes(t, "./test/testdata/slack/users.list.ok.json"))

	response, statusCode, err := slackClient.GetUsersList()
	assert.NoError(t, err)
	assert.NotEmpty(t, response)
	assert.NotEmpty(t, response.Members)
	assert.Empty(t, response.Error)
	assert.Equal(t, true, response.Ok)
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestSlackClient_SendMessage_Bad(t *testing.T) {
	badStatusCases := map[int]string{
		http.StatusBadGateway:          "Bad gateway",
		http.StatusInternalServerError: "Internal error",
		http.StatusNotFound:            "Not found",
		http.StatusForbidden:           "Forbidden",
		http.StatusBadRequest:          "Bad request",
	}

	for status, errorType := range badStatusCases {
		MockSlackResponse(status, map[string]string{}, []byte(`{"ok": false, "error": "`+errorType+`"}`))

		response, statusCode, err := slackClient.SendMessage(dto.SlackRequestChatPostMessage{})
		assert.Error(t, err)
		assert.Empty(t, response)
		assert.Equal(t, status, statusCode)
	}
}

func TestSlackClient_SendMessage_Ok(t *testing.T) {
	MockSlackResponse(http.StatusOK, map[string]string{}, test.FileToBytes(t, "./test/testdata/slack/users.list.ok.json"))

	response, statusCode, err := slackClient.SendMessage(dto.SlackRequestChatPostMessage{})
	assert.NoError(t, err)
	assert.NotEmpty(t, response)
	assert.Empty(t, response.Error)
	assert.Equal(t, true, response.Ok)
	assert.Equal(t, http.StatusOK, statusCode)
}
