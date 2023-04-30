package message

import (
	"fmt"
	"time"

	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/service/message/conversation"

	"github.com/sharovik/devbot/internal/service/history"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/log"
)

var messagesReceived = map[string]dto.BaseChatMessage{}

func answerToMessage(m dto.BaseChatMessage) error {
	response, statusCode, err := container.C.MessageClient.SendMessage(m)
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

// SendAnswerForReceivedMessage method which sends the answer by specific message
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

		//@todo: trigger learn scenario
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

// TriggerScenario triggers the scenario for selected channel
func TriggerScenario(channel string, scenario database.EventScenario, shouldRemember bool) error {
	dmAnswer := dto.DictionaryMessage{
		ScenarioID:   scenario.ID,
		EventID:      scenario.EventID,
		ReactionType: scenario.EventName,
	}

	conversation.AddConversation(scenario, dto.BaseChatMessage{
		Channel:           channel,
		AsUser:            true,
		Text:              scenario.GetUnAnsweredQuestion(),
		Ts:                time.Now(),
		DictionaryMessage: dmAnswer,
		OriginalMessage: dto.BaseOriginalMessage{
			Text: scenario.GetUnAnsweredQuestion(),
		},
	})

	if err := TriggerAnswer(channel, conversation.GetConversation(channel).LastQuestion, shouldRemember); err != nil {
		return err
	}

	return nil
}

// TriggerAnswer triggers an answer for received message
func TriggerAnswer(channel string, answerMessage dto.BaseChatMessage, shouldRemember bool) error {
	if answerMessage.Text != "" {
		if err := SendAnswerForReceivedMessage(answerMessage); err != nil {
			log.Logger().AddError(err).Msg("Failed to send prepared answers")
			conversation.FinaliseConversation(channel)

			return err
		}
	}

	if answerMessage.DictionaryMessage.IsHelpTriggered {
		return nil
	}

	conversation.SetLastQuestion(answerMessage)

	if answerMessage.DictionaryMessage.ReactionType == "" || container.C.DefinedEvents[answerMessage.DictionaryMessage.ReactionType] == nil {
		log.Logger().Warn().
			Interface("answer", answerMessage).
			Msg("Reaction type wasn't found")
		return nil
	}

	activeConversation := conversation.GetConversation(channel)
	if activeConversation.ScenarioID != int64(0) && !activeConversation.EventReadyToBeExecuted {
		log.Logger().Info().
			Interface("conversation", activeConversation).
			Msg("This conversation isn't finished yet, so event cannot be executed.")
		return nil
	}

	go func() {
		answer, err := container.C.DefinedEvents[answerMessage.DictionaryMessage.ReactionType].Execute(answerMessage)
		if err != nil {
			log.Logger().AddError(err).Msg("Failed to execute the event")
			conversation.FinaliseConversation(channel)
		}

		if answer.Text != "" {
			answerMessage.Text = answer.Text
			if err = SendAnswerForReceivedMessage(answerMessage); err != nil {
				log.Logger().AddError(err).Msg("Failed to send post-answer for selected event")
			}
		}

		//We will trigger the event history save in case when we don't have open conversation
		//or when we do have open conversation, but it is time to trigger the event execution
		//so, we can store all variables
		if shouldRemember && (conversation.GetConversation(answerMessage.Channel).ScenarioID == 0 || conversation.GetConversation(answerMessage.Channel).EventReadyToBeExecuted) {
			history.RememberEventExecution(answerMessage)
		}

		if conversation.GetConversation(answerMessage.Channel).EventReadyToBeExecuted {
			conversation.FinaliseConversation(channel)
		}
	}()

	return nil
}
