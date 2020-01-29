package slack

import (
	"time"

	"github.com/sharovik/devbot/events"
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/log"
)

const (
	eventTypeMessage             = "message"
	eventTypeDesktopNotification = "desktop_notification"
	eventTypeFileShared          = "file_shared"

	defaultAnswer = "Sorry, I don't have answer for that :("
)

var (
	messagesReceived     = map[string]dto.SlackRequestChatPostMessage{}
	acceptedMessageTypes = map[string]string{
		eventTypeMessage:             eventTypeMessage,
		eventTypeDesktopNotification: eventTypeDesktopNotification,
		eventTypeFileShared:          eventTypeFileShared,
	}
)

func isValidMessage(message *dto.SlackResponseEventMessage) bool {
	if message.Type == eventTypeDesktopNotification {
		if messagesReceived[message.Channel].Channel == "" {
			log.Logger().Debug().Str("type", message.Type).Msg("We received desktop notification, but answer wasn't prepared before.")
			return false
		}

		return true
	}

	if acceptedMessageTypes[message.Type] == "" {
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

	switch message.Type {
	case eventTypeDesktopNotification:
		if !isAnswerWasPrepared(message) {
			log.Logger().Warn().
				Interface("message_object", message).
				Msg("Answer wasn't prepared")
			return nil
		}

		answerMessage := getPreparedAnswer(message)

		if err := SendAnswerForReceivedMessage(answerMessage); err != nil {
			log.Logger().AddError(err).Msg("Failed to send prepared answers")
			return err
		}

		if answerMessage.DictionaryMessage.ReactionType == "" || events.DefinedEvents.Events[answerMessage.DictionaryMessage.ReactionType] == nil {
			log.Logger().Warn().
				Interface("answer", answerMessage).
				Msg("Reaction type wasn't found")
			return nil
		}

		go func() {
			answer, err := events.DefinedEvents.Events[answerMessage.DictionaryMessage.ReactionType].Execute(answerMessage)
			if err != nil {
				log.Logger().AddError(err).Msg("Failed to execute the event")
			}

			if answer.Text != "" {
				if err := SendAnswerForReceivedMessage(answer); err != nil {
					log.Logger().AddError(err).Msg("Failed to send post-answer for selected event")
				}
			}
		}()
	default:
		m, dmAnswer, err := analyseMessage(message)
		if err != nil {
			log.Logger().AddError(err).Msg("Failed to analyse received message")
			return err
		}

		emptyDmMessage := dto.DictionaryMessage{}
		if dmAnswer == emptyDmMessage {
			log.Logger().Debug().
				Str("type", message.Type).
				Str("text", message.Text).
				Str("team", message.Team).
				Str("source_team", message.SourceTeam).
				Str("ts", message.Ts).
				Str("user", message.User).
				Str("channel", message.Channel).
				Msg("No answer found for the received message")
		}

		//We put a dictionary message into our message object,
		// so later we can identify what kind of reaction will be executed
		m.DictionaryMessage = dmAnswer

		//We need to put this message into our small queue,
		//because we need to make sure if we received our notification.
		readyToAnswer(m)
	}

	refreshPreparedMessages()
	return nil
}

func getPreparedAnswer(message *dto.SlackResponseEventMessage) dto.SlackRequestChatPostMessage {
	return messagesReceived[message.Channel]
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

	log.Logger().Info().Interface("message", m).Msg("Message sent")
	return nil
}

func analyseMessage(message *dto.SlackResponseEventMessage) (dto.SlackRequestChatPostMessage, dto.DictionaryMessage, error) {
	var (
		responseMessage dto.SlackRequestChatPostMessage
		err             error
		dmAnswer        dto.DictionaryMessage
	)

	dmAnswer, err = container.C.Dictionary.FindAnswer(message)
	if err != nil {
		return dto.SlackRequestChatPostMessage{}, dto.DictionaryMessage{}, err
	}

	responseMessage, err = prepareAnswer(message, dmAnswer)
	if err != nil {
		return dto.SlackRequestChatPostMessage{}, dto.DictionaryMessage{}, err
	}

	return responseMessage, dmAnswer, nil
}

func readyToAnswer(message dto.SlackRequestChatPostMessage) {
	messagesReceived[message.Channel] = message
}

func isAnswerWasPrepared(message *dto.SlackResponseEventMessage) bool {
	return messagesReceived[message.Channel].Channel != ""
}

//SendAnswerForReceivedMessage method which sends the answer by specific message
func SendAnswerForReceivedMessage(message dto.SlackRequestChatPostMessage) error {
	if err := answerToMessage(message); err != nil {
		log.Logger().AddError(err).Msg("Failed to sent answer message")
		return err
	}

	messageExpired(message)
	return nil
}

func messageExpired(message dto.SlackRequestChatPostMessage) {
	delete(messagesReceived, message.Channel)
}

func refreshPreparedMessages() {
	log.Logger().Debug().
		Interface("answers_prepared", messagesReceived).
		Msg("Trigger refresh messages")

	for channelID, msg := range messagesReceived {
		if time.Since(msg.Ts).Seconds() >= 1 {
			log.Logger().Debug().
				Time("message_ts", msg.Ts).
				Interface("message_object", msg).
				Msg("Message timestamp expired")
			delete(messagesReceived, channelID)
		}
	}
}

func prepareAnswer(message *dto.SlackResponseEventMessage, dm dto.DictionaryMessage) (dto.SlackRequestChatPostMessage, error) {
	log.Logger().StartMessage("Answer prepare")

	var answer = defaultAnswer
	if dm.Answer != "" {
		answer = dm.Answer
	}

	responseMessage := dto.SlackRequestChatPostMessage{
		Channel:         message.Channel,
		Text:            answer,
		AsUser:          true,
		Ts:              time.Now(),
		OriginalMessage: *message,
	}

	log.Logger().FinishMessage("Answer prepare")
	return responseMessage, nil
}
