package slack

import (
	"fmt"
	"github.com/sharovik/devbot/internal/service/base"
	"github.com/sharovik/devbot/internal/service/history"
	"regexp"
	"time"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/log"
)

var messagesReceived = map[string]dto.BaseChatMessage{}

func getPreparedAnswer(channel string) dto.BaseChatMessage {
	return messagesReceived[channel]
}

func answerToMessage(m dto.BaseChatMessage) error {
	response, statusCode, err := container.C.MessageClient.SendMessageV2(m)
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

func readyToAnswer(message dto.BaseChatMessage) {
	messagesReceived[message.Channel] = message
}

//SendAnswerForReceivedMessage method which sends the answer by specific message
func SendAnswerForReceivedMessage(message dto.BaseChatMessage) error {
	if err := answerToMessage(message); err != nil {
		log.Logger().AddError(err).Msg("Failed to sent answer message")
		return err
	}

	messageExpired(message)
	return nil
}

func messageExpired(message dto.BaseChatMessage) {
	delete(messagesReceived, message.Channel)
}

//IsChannelID method checks the ID of received message and if it is a channel ID, then the TRUE will be returned
func IsChannelID(ID string) (isChannelMessage bool, err error) {
	regex, err := regexp.Compile(`(?i)(^C(\w+))`)
	if err != nil {
		return
	}

	matches := regex.FindStringSubmatch(ID)
	if len(matches) == 0 {
		return
	}

	return true, nil
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

func prepareAnswer(message *dto.SlackResponseEventMessage, dm dto.DictionaryMessage) (dto.BaseChatMessage, error) {
	log.Logger().StartMessage("Answer prepare")

	//If we don't have any answer, we trigger the event for unknown question
	if dm.Answer == "" {
		if !container.C.Config.LearningEnabled {
			return triggerUnknownAnswerScenario(message)
		}

		//@todo trigger learn scenario
	}

	//We will trigger the event history save in case when we don't have open conversation
	//or when we do have open conversation, but it is time to trigger the event execution
	//so we can store all variables
	if 0 == base.GetConversation(message.Channel).ScenarioID || base.GetConversation(message.Channel).EventReadyToBeExecuted {
		history.RememberEventExecution(dto.BaseChatMessage{
			Channel:           message.Channel,
			Text:              message.Text,
			AsUser:            true,
			Ts:                time.Now(),
			DictionaryMessage: dm,
			OriginalMessage: dto.BaseOriginalMessage{
				Text:  message.Text,
				User:  message.User,
				Files: message.Files,
			},
		})
	}

	responseMessage := dto.BaseChatMessage{
		Channel:         message.Channel,
		Text:            dm.Answer,
		ThreadTS:        message.ThreadTS,
		AsUser:          true,
		Ts:              time.Now(),
		OriginalMessage: message.ToBaseOriginalMessage(),
	}

	log.Logger().FinishMessage("Answer prepare")
	return responseMessage, nil
}

func triggerUnknownAnswerScenario(message *dto.SlackResponseEventMessage) (answer dto.BaseChatMessage, err error) {
	message.Text = fmt.Sprintf("similar questions %s", message.Text)
	return dto.BaseChatMessage{
		Channel:         message.Channel,
		Text:            "Hmmm",
		AsUser:          true,
		ThreadTS:        message.ThreadTS,
		Ts:              time.Now(),
		OriginalMessage: message.ToBaseOriginalMessage(),
		DictionaryMessage: dto.DictionaryMessage{
			ReactionType: "unknownquestion",
		},
	}, nil
}

//TriggerAnswer triggers an answer for received message
func TriggerAnswer(channel string) error {
	answerMessage := getPreparedAnswer(channel)

	if err := SendAnswerForReceivedMessage(answerMessage); err != nil {
		log.Logger().AddError(err).Msg("Failed to send prepared answers")
		return err
	}

	if answerMessage.DictionaryMessage.ReactionType == "" || container.C.DefinedEvents[answerMessage.DictionaryMessage.ReactionType] == nil {
		log.Logger().Warn().
			Interface("answer", answerMessage).
			Msg("Reaction type wasn't found")
		return nil
	}

	activeConversation := base.GetConversation(channel)
	if activeConversation.ScenarioID != int64(0) && !activeConversation.EventReadyToBeExecuted {
		log.Logger().Info().
			Interface("conversation", activeConversation).
			Msg("This conversation isn't finished yet, so event cannot be executed.")
		return nil
	}

	go func() {
		msg := dto.BaseChatMessage{
			Channel:           answerMessage.Channel,
			Text:              answerMessage.Text,
			AsUser:            answerMessage.AsUser,
			Ts:                answerMessage.Ts,
			DictionaryMessage: answerMessage.DictionaryMessage,
			OriginalMessage: dto.BaseOriginalMessage{
				Text:  answerMessage.OriginalMessage.Text,
				User:  answerMessage.OriginalMessage.User,
				Files: answerMessage.OriginalMessage.Files,
			},
		}
		answer, err := container.C.DefinedEvents[answerMessage.DictionaryMessage.ReactionType].Execute(msg)
		if err != nil {
			log.Logger().AddError(err).Msg("Failed to execute the event")
		}

		if answer.Text != "" {
			answerMessage.Text = answer.Text
			if err := SendAnswerForReceivedMessage(answerMessage); err != nil {
				log.Logger().AddError(err).Msg("Failed to send post-answer for selected event")
			}
		}
	}()

	return nil
}
