package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/log"
)

// SlackClient client for message api calls
type SlackClient struct {
	BaseMessageClient
	Client     *http.Client
	BaseURL    string
	OAuthToken string
}

// AttachFileTo method for attachment file send to specific channel
func (client SlackClient) AttachFileTo(channel string, pathToFile string, filename string) ([]byte, int, error) {
	log.Logger().StartMessage("Slack attachment request")

	var buf bytes.Buffer

	writer := multipart.NewWriter(&buf)

	log.Logger().Debug().
		Str("channel", channel).
		Str("file_path", pathToFile).
		Str("filename", filename).
		Msg("Received parameters")

	fieldWriter, err := writer.CreateFormField("channels")
	if err != nil {
		return nil, 0, err
	}

	if _, err := fieldWriter.Write([]byte(channel)); err != nil {
		return nil, 0, err
	}

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, 0, err
	}

	file, err := os.Open(pathToFile)
	if err != nil {
		return nil, 0, err
	}

	if _, err = io.Copy(part, file); err != nil {
		return nil, 0, err
	}

	writer.Close()
	file.Close()

	return client.HTTPClient.Post("/files.upload", &buf, map[string]string{})
}

// SendMessage method for post message send through simple API request
func (client SlackClient) SendMessage(message dto.BaseChatMessage) (resp dto.BaseResponseInterface, status int, err error) {
	log.Logger().Debug().Interface("message", message).Msg("Start chat.postMessage")
	byteStr, err := json.Marshal(dto.SlackRequestChatPostMessage{
		Channel:           message.Channel,
		Text:              message.Text,
		AsUser:            message.AsUser,
		ThreadTS:          message.ThreadTS,
		Ts:                message.Ts,
		DictionaryMessage: dto.DictionaryMessage{},
		OriginalMessage:   dto.SlackResponseEventMessage{},
	})
	if err != nil {
		return resp, 0, err
	}

	response, statusCode, err := client.HTTPClient.Post("/chat.postMessage", byteStr, map[string]string{})
	if err != nil {
		log.Logger().AddError(err).
			RawJSON("response", response).
			Int("status_code", statusCode).
			Msg("Failed send message")
		return resp, statusCode, err
	}

	var dtoResponse dto.SlackResponseChatPostMessage
	dtoResponse.SetByteResponse(response)
	if err := json.Unmarshal(response, &dtoResponse); err != nil {
		return resp, statusCode, err
	}

	if !dtoResponse.Ok {
		return &dtoResponse, statusCode, errors.New(dtoResponse.Error)
	}

	log.Logger().Debug().Interface("message", message).Msg("Finish chat.postMessage")
	dtoResponse.SetByteResponse(response)

	return &dtoResponse, statusCode, nil
}

// GetConversationsList method which returns the conversations list of current workspace
func (client SlackClient) GetConversationsList() (dto.SlackResponseConversationsList, int, error) {
	response, statusCode, err := client.HTTPClient.Get("/conversations.list", map[string]string{})
	if err != nil {
		return dto.SlackResponseConversationsList{}, statusCode, err
	}

	var dtoResponse dto.SlackResponseConversationsList
	if err := json.Unmarshal(response, &dtoResponse); err != nil {
		return dto.SlackResponseConversationsList{}, statusCode, err
	}

	if !dtoResponse.Ok {
		return dtoResponse, statusCode, errors.New(dtoResponse.Error)
	}

	return dtoResponse, statusCode, nil
}

// GetUsersList method which returns the users list of current workspace
func (client SlackClient) GetUsersList() (dto.SlackResponseUsersList, int, error) {
	response, statusCode, err := client.HTTPClient.Get("/users.list", map[string]string{})
	if err != nil {
		return dto.SlackResponseUsersList{}, statusCode, err
	}

	var dtoResponse dto.SlackResponseUsersList
	if err := json.Unmarshal(response, &dtoResponse); err != nil {
		return dto.SlackResponseUsersList{}, statusCode, err
	}

	if !dtoResponse.Ok {
		return dtoResponse, statusCode, errors.New(dtoResponse.Error)
	}

	return dtoResponse, statusCode, nil
}

// GetUsersListPaged method which returns the users list of current workspace using selected cursor
func (client SlackClient) GetUsersListPaged(cursor string) (result dto.SlackResponseUsersList, err error) {
	response, _, err := client.HTTPClient.Get("/users.list", map[string]string{
		"cursor": cursor,
	})
	if err != nil {
		return dto.SlackResponseUsersList{}, err
	}

	var dtoResponse dto.SlackResponseUsersList
	if err := json.Unmarshal(response, &dtoResponse); err != nil {
		return dto.SlackResponseUsersList{}, err
	}

	if !dtoResponse.Ok {
		return dtoResponse, errors.New(dtoResponse.Error)
	}

	return dtoResponse, nil
}
