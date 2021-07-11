package slack

import (
	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/helper"
	"github.com/sharovik/devbot/internal/service/base"
	"regexp"
	"time"

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

	if isGlobalAlertTriggered(message.Text) {
		log.Logger().Debug().Msg("The global alert is triggered. Skipping.")
		return false
	}

	if message.User == container.C.Config.SlackConfig.BotUserID {
		log.Logger().Debug().Msg("This message is from our bot user")
		return false
	}

	return true
}

func getPreparedAnswer(message *dto.SlackResponseEventMessage) dto.SlackRequestChatPostMessage {
	return messagesReceived[message.Channel]
}

func isGlobalAlertTriggered(text string) bool {
	re, err := regexp.Compile(`(?i)(\<\!(here|channel)\>)`)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to parse global alert text part")
		return false
	}

	return re.MatchString(text)
}

func answerToMessage(m dto.SlackRequestChatPostMessage) error {
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

func analyseMessage(message *dto.SlackResponseEventMessage) (dto.SlackRequestChatPostMessage, dto.DictionaryMessage, error) {
	var (
		responseMessage dto.SlackRequestChatPostMessage
		err             error
		dmAnswer        dto.DictionaryMessage
	)

	//Now we need to check if there was already opened conversation for this channel
	//If so, then we need to get the Answer from this scenario
	openConversation := base.GetConversation(message.Channel)

	IsScenarioStopTriggered := base.IsScenarioStopTriggered(message.Text)
	if openConversation.ScenarioID != 0 && !IsScenarioStopTriggered {
		questions, err := container.C.Dictionary.GetQuestionsByScenarioID(openConversation.ScenarioID)
		//if we don't have an error here, then we can proceed with the questions preparing for scenarios
		if err != nil {
			log.Logger().AddError(err).Msg("Failed to get the list of question by the scenarioID")
		}

		//We do questions check only in case of multiple questions attached to the scenario.
		//In other cases we do as we did before
		if len(questions) > 1 && err == nil {
			scenarioNextQuestion := extractLastQuestionID(openConversation, questions)

			//if we have 0 as ID, then that means we didn't found next question, so we can try to execute the event and have a look what will be
			if scenarioNextQuestion.ID == int64(0) {
				//We mark the scenario, to be executed
				openConversation = base.MarkAsReadyEventToBeExecuted(openConversation)

				//We get the last question object from the scenario and we use it as the answer
				dmAnswer = dto.DictionaryMessage{
					ScenarioID:            openConversation.ScenarioID,
					Question:              scenarioNextQuestion.Question,
					QuestionID:            scenarioNextQuestion.ID,
					Regex:                 "",
					Answer:                "Ok",
					MainGroupIndexInRegex: "",
					ReactionType:          openConversation.ReactionType,
				}

				//We add the last message to the variables
				openConversation.Variables = append(openConversation.Variables, message.Text)
				base.CurrentConversations[message.Channel] = openConversation

				//In that case we have the last questionID so that means, we need to use the last question here.
			} else {
				dmAnswer = dto.DictionaryMessage{
					ScenarioID:            openConversation.ScenarioID,
					Question:              scenarioNextQuestion.Question,
					QuestionID:            scenarioNextQuestion.ID,
					Regex:                 "",
					Answer:                scenarioNextQuestion.Answer,
					MainGroupIndexInRegex: "",
					ReactionType:          scenarioNextQuestion.ReactionType,
				}

				//We also add the new state for this conversation
				base.AddConversation(message.Channel, dmAnswer.QuestionID, dto.BaseChatMessage{
					Channel:           message.Channel,
					Text:              message.Text,
					AsUser:            false,
					Ts:                time.Now(),
					DictionaryMessage: dmAnswer,
					OriginalMessage: dto.BaseOriginalMessage{
						Text:  message.Text,
						User:  message.User,
						Files: message.Files,
					},
				}, message.Text)
			}

			responseMessage, err = prepareAnswer(message, dmAnswer)
			if err != nil {
				return dto.SlackRequestChatPostMessage{}, dto.DictionaryMessage{}, err
			}

			return responseMessage, dmAnswer, nil
		}
	}

	if IsScenarioStopTriggered {
		dmAnswer = dto.DictionaryMessage{
			ScenarioID:            0,
			Answer:                "Ok, no more questions!",
			QuestionID:            0,
			Question:              message.Text,
			Regex:                 "",
			MainGroupIndexInRegex: "",
			ReactionType:          "text",
		}
		base.DeleteConversation(message.Channel)
	} else {
		dmAnswer, err = container.C.Dictionary.FindAnswer(message)
		if err != nil {
			return dto.SlackRequestChatPostMessage{}, dto.DictionaryMessage{}, err
		}
	}

	questions, err := container.C.Dictionary.GetQuestionsByScenarioID(dmAnswer.ScenarioID)
	//if we don't have an error here, then we can proceed with the questions preparing for scenarios
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to get the list of question by the scenarioID")
	}

	isHelpAnswerTriggered, err := helper.HelpMessageShouldBeTriggered(message.Text)
	//If the questions amount is more than 1, we need to start the conversation algorithm
	if len(questions) > 1 && !isHelpAnswerTriggered {
		base.AddConversation(message.Channel, dmAnswer.QuestionID, dto.BaseChatMessage{
			Channel:           message.Channel,
			Text:              message.Text,
			AsUser:            false,
			Ts:                time.Now(),
			DictionaryMessage: dmAnswer,
			OriginalMessage: dto.BaseOriginalMessage{
				Text:  message.Text,
				User:  message.User,
				Files: message.Files,
			},
		}, "")
	}

	responseMessage, err = prepareAnswer(message, dmAnswer)
	if err != nil {
		return dto.SlackRequestChatPostMessage{}, dto.DictionaryMessage{}, err
	}

	return responseMessage, dmAnswer, nil
}

func extractLastQuestionID(openConversation base.Conversation, questions []database.QuestionObject) database.QuestionObject {
	shouldAskNewQuestion := false
	lastQuestion := database.QuestionObject{}
	lastQuestionID := int64(0)
	for _, question := range questions {

		lastQuestionID = question.ID
		if openConversation.ScenarioQuestionID == question.ID {
			shouldAskNewQuestion = true
			continue
		}

		if shouldAskNewQuestion {
			lastQuestion = question
			break
		}

		//In that case we always store here the latest question object
		lastQuestion = question
	}

	//This means there was the last question already triggered
	if lastQuestionID == openConversation.ScenarioQuestionID {
		lastQuestion = database.QuestionObject{}
	}

	return lastQuestion
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
