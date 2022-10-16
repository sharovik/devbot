package slack

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/sharovik/devbot/internal/client"
	"github.com/sharovik/devbot/internal/config"
	"github.com/sharovik/devbot/internal/service/analiser"
	"github.com/sharovik/devbot/internal/service/base"
	"net/http"
	"strings"
	"time"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/log"
	"golang.org/x/net/websocket"
)

const slackAPIOrigin = "https://slack.com/api"

//S slack service object
var S base.ServiceInterface

func InitService() {
	S = EventsApiService{}
}

//EventsApiService struct of slack service
type EventsApiService struct {
}

func (EventsApiService) fetchMainChannelID() error {
	availableChannels, statusCode, err := container.C.MessageClient.GetConversationsList()
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

func (EventsApiService) fetchBotUserID() error {
	availableUsers, statusCode, err := container.C.MessageClient.GetUsersList()
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

//BeforeWSConnectionStart runs methods before the WS connection start
func (s EventsApiService) BeforeWSConnectionStart() error {
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
func (s EventsApiService) InitWebSocketReceiver() error {
	if err := s.BeforeWSConnectionStart(); err != nil {
		log.Logger().AddError(err).Msg("Failed to prepare service for WS connection")
		return err
	}

	ws, statusCode, err := s.wsConnect()
	if err != nil {
		log.Logger().AddError(err).Int("status_code", statusCode).Msg("Failed connect to the websocket")
		return err
	}

	var (
		event interface{}
	)

	for {
		var message dto.SlackResponseEventApiMessage

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

		if message.Type == "hello" {
			continue
		}

		log.Logger().Debug().
			RawJSON("message_body", str).
			Str("envelope_id", message.EnvelopeID).
			Msg("Received event message")

		if err := acknowledge(ws, message.EnvelopeID); err != nil {
			log.Logger().AddError(err).
				RawJSON("message_body", str).
				Str("envelope_id", message.EnvelopeID).
				Msg("Failed to acknowledge the message")

			return err
		}

		if !isValidMessage(MessageAttributes{
			Type:    message.Payload.Event.Type,
			Channel: message.Payload.Event.Channel,
			Text:    message.Payload.Event.Text,
			User:    message.Payload.Event.User,
			BotID:   message.Payload.Event.BotID,
		}) {
			continue
		}

		if err := s.ProcessMessage(&message); err != nil {
			log.Logger().AddError(err).Interface("message_object", &message).Msg("Can't check or answer to the message")
		}

		base.CleanUpExpiredMessages()
	}
}

func acknowledge(ws *websocket.Conn, envelopeId string) error {
	res := dto.SlackRequestAcknowledge{
		EnvelopeID: envelopeId,
	}

	log.Logger().Debug().
		Str("envelope_id", envelopeId).
		Msg("Acknowledge event message")

	return websocket.JSON.Send(ws, res)
}

//ProcessMessage processes the message from the WS connection
func (s EventsApiService) ProcessMessage(msg interface{}) error {
	message := msg.(*dto.SlackResponseEventApiMessage)
	log.Logger().Debug().
		Str("type", message.Type).
		Str("text", message.Payload.Event.Text).
		Str("team", message.Payload.Event.Team).
		Str("ts", message.Payload.Event.Ts).
		Str("user", message.Payload.Event.User).
		Str("channel", message.Payload.Event.Channel).
		Msg("Message received")

	//We need to trim the message before all checks
	message.Payload.Event.Text = strings.TrimSpace(message.Payload.Event.Text)

	dmAnswer, err := analiser.GetDmAnswer(analiser.Message{
		Channel: message.Payload.Event.Channel,
		User:    message.Payload.Event.User,
		Text:    message.Payload.Event.Text,
	})

	m, err := prepareAnswer(&dto.SlackResponseEventMessage{
		Channel:      message.Payload.Event.Channel,
		ClientMsgID:  message.Payload.Event.ClientMsgID,
		DisplayAsBot: false,
		EventTs:      message.Payload.Event.EventTs,
		ThreadTS:     message.Payload.Event.ThreadTS,
		Files:        nil,
		Team:         message.Payload.Event.Team,
		Text:         message.Payload.Event.Text,
		Ts:           message.Payload.Event.Ts,
		Type:         message.Payload.Event.Type,
		User:         message.Payload.Event.User,
	}, dmAnswer)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to analyse received message")
		return err
	}

	emptyDmMessage := dto.DictionaryMessage{}
	if dmAnswer == emptyDmMessage {
		log.Logger().Debug().
			Str("type", message.Payload.Event.Type).
			Str("text", message.Payload.Event.Text).
			Str("team", message.Payload.Event.Team).
			Str("ts", message.Payload.Event.Ts).
			Str("user", message.Payload.Event.User).
			Str("channel", message.Payload.Event.Channel).
			Msg("No answer found for the received message")
	} else {
		//We put a dictionary message into our message object,
		// so later we can identify what kind of reaction will be executed
		m.DictionaryMessage = dmAnswer
	}

	//We need to put this message into our small queue,
	//because we need to make sure if we received our notification.
	readyToAnswer(m)

	if err := TriggerAnswer(message.Payload.Event.Channel); err != nil {
		log.Logger().AddError(err).Msg("Failed trigger the answer")
		return err
	}

	openConversation := base.GetConversation(message.Payload.Event.Channel)
	if openConversation.ScenarioID != 0 {
		isChannelMessage, err := IsChannelID(message.Payload.Event.Channel)
		if err != nil {
			log.Logger().AddError(err).Msg("Error received during channel ID check")
			return err
		}

		if isChannelMessage {
			return TriggerAnswer(message.Payload.Event.Channel)
		}
	}

	refreshPreparedMessages()
	return nil
}

func getWSClient() client.SlackClient {
	netTransport := &http.Transport{
		TLSHandshakeTimeout: time.Duration(container.C.Config.HTTPClient.TLSHandshakeTimeout) * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: container.C.Config.HTTPClient.InsecureSkipVerify,
		},
	}

	httpClient := http.Client{
		Timeout:   time.Duration(container.C.Config.HTTPClient.RequestTimeout) * time.Second,
		Transport: netTransport,
	}

	c := client.HTTPClient{
		Client: &httpClient,
	}

	c.SetOauthToken(container.C.Config.SlackConfig.OAuthToken)
	c.SetBaseURL(container.C.Config.SlackConfig.BaseURL)

	sc := client.SlackClient{}
	sc.HttpClient = &c

	return sc
}

//wsConnect method for receiving of websocket URL which we will use for our connection
func (EventsApiService) wsConnect() (*websocket.Conn, int, error) {
	response, statusCode, err := getWSClient().HttpClient.Post("/apps.connections.open", []byte{}, map[string]string{})
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

	ws, err := websocket.Dial(dtoResponse.URL, "", slackAPIOrigin)

	return ws, statusCode, nil
}
