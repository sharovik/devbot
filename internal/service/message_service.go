package service

import (
	"fmt"
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/log"
	"regexp"
)

const (
	eventTypeMessage = "message"

	defaultAnswer = "Sorry I don't know what to answer :("
)

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

	m, dmAnswer, err := analyseMessage(message)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to analyse received")
		return err
	}

	if err := answerToMessage(m); err != nil {
		log.Logger().AddError(err).Msg("Failed to sent answer message")
		return err
	}

	if dmAnswer.ReactionType != "" {
		//@todo: Add reaction action here
		log.Logger().Debug().Str("reaction_type", dmAnswer.ReactionType).Msg("Reaction type executed")
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

func analyseMessage(message *dto.SlackResponseEventMessage) (dto.SlackRequestChatPostMessage, dto.DictionaryMessage, error) {
	var (
		responseMessage dto.SlackRequestChatPostMessage
		err             error
		dmAnswer dto.DictionaryMessage
	)

	dmAnswer = findDictionaryMessageType(message)

	responseMessage, err = prepareAnswer(message, dmAnswer)
	if err != nil {
		return dto.SlackRequestChatPostMessage{}, dto.DictionaryMessage{}, err
	}

	return responseMessage, dmAnswer, nil
}

func findDictionaryMessageType(message *dto.SlackResponseEventMessage) dto.DictionaryMessage {
	var dmAnswer dto.DictionaryMessage
	for index, dm := range getMessageDictionary(message) {
		re := regexp.MustCompile(dm.Question)

		matches := re.FindStringSubmatch(message.Text)
		if matches != nil {
			log.Logger().Debug().
				Int("index", index).
				Str("found_matches", dm.Question).
				Str("selected_answer", dm.Answer).
				Str("selected_type_of_action", dm.ReactionType).
				Interface("matches", matches).
				Msg("Selected answer")

			dmAnswer = dm
			if dm.MainGroupIndexInRegex != 0 {
				dmAnswer.Answer = fmt.Sprintf(dm.Answer, matches[dm.MainGroupIndexInRegex])
			}

			return dmAnswer
		}
	}

	return dmAnswer
}

func getMessageDictionary(message *dto.SlackResponseEventMessage) []dto.DictionaryMessage {
	if message.Files != nil {
		log.Logger().Debug().Str("dictionary", "file_message_dictionary").Msg("Selected dictionary")
		return container.C.Dictionary.FileMessageDictionary
	}

	log.Logger().Debug().Str("dictionary", "text_message_dictionary").Msg("Selected dictionary")
	return container.C.Dictionary.TextMessageDictionary
}

func prepareAnswer(message *dto.SlackResponseEventMessage, dm dto.DictionaryMessage) (dto.SlackRequestChatPostMessage, error) {
	log.Logger().StartMessage("Answer prepare")

	var answer = defaultAnswer
	if dm.Answer != "" {
		answer = dm.Answer
	}

	responseMessage := dto.SlackRequestChatPostMessage{
		Channel: message.Channel,
		Text:    answer,
		AsUser:  true,
	}

	log.Logger().FinishMessage("Answer prepare")
	return responseMessage, nil
}
