package examplescenario

import (
	"fmt"
	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/helper"
	"github.com/sharovik/devbot/internal/service/base"
	"regexp"
	"time"

	"github.com/sharovik/devbot/internal/log"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
)

const (
	//EventName the name of the event
	EventName = "examplescenario"

	//EventVersion the version of the event
	EventVersion = "1.0.0"

	helpMessage = "Ask me `write a message` and you will see the answer."

	regexChannel = `(?im)(?:[<#@]|(?:&lt;))(\w+)(?:[|>]|(?:&gt;))`

	stepMessage = "What I need to write?"
	stepChannel = "Where I need to post this message? If it's channel, the channel should be public."
)

//EventStruct the struct for the event object. It will be used for initialisation of the event in defined-events.go file.
type EventStruct struct {
	EventName string
}

//Event - object which is ready to use
var Event = EventStruct{
	EventName: EventName,
}

//Execute method which is called by message processor
func (e EventStruct) Execute(message dto.BaseChatMessage) (dto.BaseChatMessage, error) {
	isHelpAnswerTriggered, err := helper.HelpMessageShouldBeTriggered(message.OriginalMessage.Text)
	if err != nil {
		log.Logger().Warn().Err(err).Msg("Something went wrong with help message parsing")
	}

	if isHelpAnswerTriggered {
		message.Text = helpMessage
		return message, nil
	}

	currentConversation := base.GetConversation(message.Channel)

	whatToWrite := ""
	whereToWrite := ""

	//If we don't have all variables for our conversation, that means, we didn't received answers for all questions of our scenario
	if len(currentConversation.Variables) != 2 {
		message.Text = "Not all variables are defined."

		//We remove this conversation from the memory, because it is expired. You must do this, otherwise bot will think that this conversation is still opened.
		base.DeleteConversation(message.Channel)
		return message, nil
	}

	if currentConversation.Variables[0] != "" {
		whatToWrite = removeCurrentUserFromTheMessage(currentConversation.Variables[0])
	}

	if currentConversation.Variables[1] != "" {
		whereToWrite = extractChannelName(currentConversation.Variables[1])

		if whereToWrite == "" {
			message.Text = "Something went wrong and I can't parse properly the channel name."
			base.DeleteConversation(message.Channel)
			return message, nil
		}
	}

	_, _, err = container.C.MessageClient.SendMessageV2(dto.BaseChatMessage{
		Channel:           whereToWrite,
		Text:              whatToWrite,
		AsUser:            true,
		Ts:                time.Now(),
		DictionaryMessage: dto.DictionaryMessage{},
		OriginalMessage:   dto.BaseOriginalMessage{},
	})

	if err != nil {
		log.Logger().AddError(err).Msg("Failed to send post-answer for selected event")
		message.Text = fmt.Sprintf("Failed to send the message to the channel %s.\nReason: ```%s```", currentConversation.Variables[1], err.Error())
		return message, err
	}

	//This answer will be show once the event get triggered.
	//Leave message.Text empty, once you need to not show the message, once this event get triggered.
	message.Text = "Done"

	//We remove this conversation from the memory, because it is expired
	base.DeleteConversation(message.Channel)
	return message, nil
}

func extractChannelName(text string) string {
	matches := helper.FindMatches(regexChannel, text)
	if len(matches) == 0 {
		return ""
	}

	if matches["1"] == "" {
		return ""
	}

	return matches["1"]
}

//Install method for installation of event
func (e EventStruct) Install() error {
	log.Logger().Debug().
		Str("event_name", EventName).
		Str("event_version", EventVersion).
		Msg("Triggered event installation")

	if err := container.C.Dictionary.InstallNewEventScenario(database.NewEventScenario{
		EventName:    EventName,
		EventVersion: EventVersion,
		Questions: []database.Question{
			{
				Question:      "write a message",
				Answer:        stepMessage,
				QuestionRegex: "(?i)(write a message)",
				QuestionGroup: "",
			},
			{
				Question:      "",
				Answer:        stepChannel,
				QuestionRegex: "(?i)(write a message)",
				QuestionGroup: "",
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

//Update for event update actions
func (e EventStruct) Update() error {
	return nil
}

func removeCurrentUserFromTheMessage(message string) string {
	regexString := fmt.Sprintf("(?im)(<@%s>)", container.C.Config.SlackConfig.BotUserID)
	re := regexp.MustCompile(regexString)
	result := re.ReplaceAllString(message, ``)

	return result
}
