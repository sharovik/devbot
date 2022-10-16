package client_test

import (
	mockhttp "github.com/karupanerura/go-mock-http-response"
	"net/http"
	"os"
	"path"
	"runtime"
	"testing"

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
	slackClient = &client.SlackClient{}
	slackClient.HttpClient = &client.HTTPClient{Client: mockhttp.NewResponseMock(statusCode, headers, body).MakeClient()}
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

		response, statusCode, err := slackClient.SendMessageV2(dto.BaseChatMessage{})
		assert.Error(t, err)
		assert.Empty(t, response)
		assert.Equal(t, status, statusCode)
	}
}

func TestSlackClient_SendMessage_Ok(t *testing.T) {
	MockSlackResponse(http.StatusOK, map[string]string{}, test.FileToBytes(t, "./test/testdata/slack/users.list.ok.json"))

	response, statusCode, err := slackClient.SendMessageV2(dto.BaseChatMessage{})

	assert.NoError(t, err)
	assert.NotEmpty(t, response)
	assert.Equal(t, http.StatusOK, statusCode)
}
