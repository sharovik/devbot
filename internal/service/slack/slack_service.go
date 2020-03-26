package slack

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/sharovik/devbot/internal/config"
	"github.com/sharovik/devbot/internal/service/base"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/log"
	"golang.org/x/net/websocket"
)

//Service struct of slack service
type Service struct {
}

//S slack service object
var S base.ServiceInterface = Service{}

func (Service) fetchMainChannelID() error {
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
		if err := container.C.Config.SetToEnv(config.SlackEnvMainChannelID, generalChannel.ID, true); err != nil {
			log.Logger().AddError(err).Str("channel_id", generalChannel.ID).Msg("Failed to save slackEnvMainChannelID in .env file")
			return err
		}

		container.C.Config.SlackConfig.MainChannelID = generalChannel.ID
	}

	return nil
}

func (Service) fetchBotUserID() error {
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
		if err := container.C.Config.SetToEnv(config.SlackEnvUserID, botMember.ID, true); err != nil {
			log.Logger().AddError(err).Str("user_id", botMember.ID).Msg("Failed to save slackEnvMainChannelID in .env file")
			return err
		}

		container.C.Config.SlackConfig.BotUserID = botMember.ID
	}

	return nil
}

func (s Service) beforeWSConnectionStart() error {
	if container.C.Config.SlackConfig.MainChannelID == "" {
		log.Logger().Info().Msg("Main channel ID wasn't specified. Trying to fetch main channel from API")
		if err := s.fetchMainChannelID(); err != nil {
			log.Logger().AddError(err).Msg("Failed to fetch channels")
			return err
		}
	}

	if container.C.Config.SlackConfig.BotUserID == "" {
		log.Logger().Info().Msg("Bot user ID wasn't specified. Trying to fetch user ID from API")
		if err := s.fetchBotUserID(); err != nil {
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
func (s Service) InitWebSocketReceiver() error {
	if err := s.beforeWSConnectionStart(); err != nil {
		log.Logger().AddError(err).Msg("Failed to prepare service for WS connection")
		return err
	}

	ws, statusCode, err := s.wsConnect()
	if err != nil {
		log.Logger().AddError(err).Int("status_code", statusCode).Msg("Failed connect to the websocket")
		return err
	}

	var (
		event   interface{}
		message dto.SlackResponseEventMessage
	)

	for {
		//Receive message
		err = websocket.JSON.Receive(ws, &event)
		if err != nil {
			log.Logger().AddError(err).Msg("Something went wrong with message receiving from EventsAPI")
			return err
		}

		str, _ := json.Marshal(&event)
		if strings.Contains(string(str), `"channel":{"created"`) || strings.Contains(string(str), `"type":"user_change"`) {
			log.Logger().Warn().RawJSON("message_body", str).Msg("Received unsupported type of message")
			continue
		}

		if err := json.Unmarshal(str, &message); err != nil {
			log.Logger().AddError(err).
				RawJSON("message_body", str).
				Msg("Something went wrong with message parsing")
			return err
		}

		if !isValidMessage(&message) {
			continue
		}

		if err := processMessage(&message); err != nil {
			log.Logger().AddError(err).Interface("message_object", &message).Msg("Can't check or answer to the message")
		}
	}
}

//wsConnect method for receiving of websocket URL which we will use for our connection
func (Service) wsConnect() (*websocket.Conn, int, error) {
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
