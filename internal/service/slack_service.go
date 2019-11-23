package service

import (
	"encoding/json"
	"errors"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/log"
	"golang.org/x/net/websocket"
)

const (
	slackEnvUserID        = "SLACK_USER_ID"
	slackEnvMainChannelID = "SLACK_MAIN_CHANNEL_ID"

	eventTypeMessage = "message"
)

func fetchMainChannelID() error {
	availableChannels, statusCode, err := container.C.SlackClient.GetConversationsList()
	if err != nil {
		log.Logger().AddError(err).Int("status_code", statusCode).Msg("Failed conversations list fetching")
		return err
	}

	var generalChannel dto.Channel
	for _, channel := range availableChannels.Channels {
		if channel.Name == container.C.Config.SlackConfig.MainChannelAlias {
			generalChannel = channel
			break
		}
	}

	if container.C.Config.SlackConfig.MainChannelID == "" {
		if err := container.C.Config.SetToEnv(slackEnvMainChannelID, generalChannel.ID, true); err != nil {
			log.Logger().AddError(err).Str("channel_id", generalChannel.ID).Msg("Failed to save slackEnvMainChannelID in .env file")
			return err
		}

		container.C.Config.SlackConfig.MainChannelID = generalChannel.ID
	}

	return nil
}

func fetchBotUserID() error {
	availableUsers, statusCode, err := container.C.SlackClient.GetUsersList()
	if err != nil {
		log.Logger().AddError(err).Int("status_code", statusCode).Msg("Failed conversations list fetching")
		return err
	}

	var botMember dto.SlackMember
	for _, member := range availableUsers.Members {
		if member.Profile.RealName == container.C.Config.SlackConfig.BotName {
			botMember = member
			break
		}
	}

	if container.C.Config.SlackConfig.BotUserID == "" {
		if err := container.C.Config.SetToEnv(slackEnvUserID, botMember.ID, true); err != nil {
			log.Logger().AddError(err).Str("user_id", botMember.ID).Msg("Failed to save slackEnvMainChannelID in .env file")
			return err
		}

		container.C.Config.SlackConfig.BotUserID = botMember.ID
	}

	return nil
}

func beforeWSConnectionStart() error {
	if container.C.Config.SlackConfig.MainChannelID == "" {
		log.Logger().Info().Msg("Main channel ID wasn't specified. Trying to fetch main channel from API")
		if err := fetchMainChannelID(); err != nil {
			log.Logger().AddError(err).Msg("Failed to fetch channels")
			return err
		}
	}

	if container.C.Config.SlackConfig.BotUserID == "" {
		log.Logger().Info().Msg("Bot user ID wasn't specified. Trying to fetch user ID from API")
		if err := fetchBotUserID(); err != nil {
			log.Logger().AddError(err).Msg("Failed to fetch user ID")
			return err
		}
	}

	log.Logger().AppendGlobalContext(map[string]interface{}{
		"main_channel_id":    container.C.Config.SlackConfig.MainChannelID,
		"main_channel_alias": container.C.Config.SlackConfig.MainChannelAlias,
		"bot_user_id":        container.C.Config.SlackConfig.BotUserID,
		"bot_user_name":      container.C.Config.SlackConfig.BotName,
	})

	return nil
}

//InitWebSocketReceiver method for initialization of websocket receiver
func InitWebSocketReceiver() error {
	if err := beforeWSConnectionStart(); err != nil {
		log.Logger().AddError(err).Msg("Failed to prepare service for WS connection")
		return err
	}

	ws, statusCode, err := wsConnect()
	if err != nil {
		log.Logger().AddError(err).Int("status_code", statusCode).Msg("Failed send message")
		return err
	}

	var (
		message dto.SlackResponseEventMessage
	)

	for {
		//Receive message
		err = websocket.JSON.Receive(ws, &message)
		if err != nil {
			log.Logger().AddError(err).Msg("Something went wrong with message receiving from EventsAPI")
		}

		if !isValidMessage(&message) {
			continue
		}

		if err := processMessage(&message); err != nil {
			log.Logger().AddError(err).Interface("object", &message).Msg("Can't check or answer to the message")
		}
	}
}

//wsConnect method for receiving of websocket URL which we will use for our connection
func wsConnect() (*websocket.Conn, int, error) {
	response, statusCode, err := container.C.SlackClient.Get("/rtm.connect")
	if err != nil {
		log.Logger().AddError(err).RawJSON("response", response).Int("status_code", statusCode).Msg("Failed send message")
		return &websocket.Conn{}, statusCode, err
	}

	var dtoResponse dto.SlackResponseRTMConnect
	if err := json.Unmarshal(response, &dtoResponse); err != nil {
		return &websocket.Conn{}, statusCode, err
	}

	if !dtoResponse.Ok {
		return &websocket.Conn{}, statusCode, errors.New(dtoResponse.Error)
	}

	ws, err := websocket.Dial(dtoResponse.URL, "", "https://api.slack.com/")

	return ws, statusCode, nil
}

func isValidMessage(message *dto.SlackResponseEventMessage) bool {
	if message.Type != eventTypeMessage {
		log.Logger().Debug().Str("type", message.Type).Msg("Skip message check for this message type")
		return false
	}

	if message.Channel == "" {
		log.Logger().Debug().Msg("Message channel cannot be empty")
		return false
	}

	if message.User == container.C.Config.SlackConfig.BotUserID {
		log.Logger().Debug().Msg("This message is from our bot user")
		return false
	}

	return true
}

func processMessage(message *dto.SlackResponseEventMessage) error {
	log.Logger().Debug().
		Str("type", message.Type).
		Str("text", message.Text).
		Str("team", message.Team).
		Str("source_team", message.SourceTeam).
		Str("ts", message.Ts).
		Str("user", message.User).
		Str("channel", message.Channel).
		Msg("Message received")

	m := dto.SlackRequestChatPostMessage{
		Channel: message.Channel,
		Text:    "Text",
		AsUser:  true,
	}

	if err := answerToMessage(m); err != nil {
		log.Logger().AddError(err).Msg("Failed to sent answer message")
	}

	log.Logger().Info().Interface("message", m).Msg("Message sent")
	return nil
}

func answerToMessage(m dto.SlackRequestChatPostMessage) error {
	response, statusCode, err := container.C.SlackClient.SendMessage(m)
	if err != nil {
		log.Logger().AddError(err).
			Interface("response", response).
			Interface("status", statusCode).
			Msg("Failed to sent answer message")
		return err
	}

	return nil
}
