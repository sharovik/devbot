package examplescenario

import (
	"fmt"
	"regexp"
	"time"

	"github.com/sharovik/devbot/internal/service/message/conversation"

	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/helper"

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

// EventStruct the struct for the event object. It will be used for initialisation of the event in defined-events.go file.
type EventStruct struct {
}

// Event - object which is ready to use
var Event = EventStruct{}

// Help retrieves the help message
func (e EventStruct) Help() string {
	return helpMessage
}

// Alias retrieves the event alias
func (e EventStruct) Alias() string {
	return EventName
}

// Execute method which is called by message processor
func (e EventStruct) Execute(message dto.BaseChatMessage) (dto.BaseChatMessage, error) {
	currentConversation := conversation.GetConversation(message.Channel)

	whatToWrite := ""
	whereToWrite := ""

	//If we don't have all variables for our conversation, that means, we didn't receive answers for all questions of our scenario
	if len(currentConversation.Scenario.RequiredVariables) != 2 {
		message.Text = "Not all variables are defined."

		return message, nil
	}

	if currentConversation.Scenario.RequiredVariables[0].Value != "" {
		whatToWrite = removeCurrentUserFromTheMessage(currentConversation.Scenario.RequiredVariables[0].Value)
	}

	if currentConversation.Scenario.RequiredVariables[1].Value != "" {
		whereToWrite = extractChannelName(currentConversation.Scenario.RequiredVariables[1].Value)

		if whereToWrite == "" {
			message.Text = "Something went wrong and I can't parse properly the channel name."
			return message, nil
		}
	}

	_, _, err := container.C.MessageClient.SendMessage(dto.BaseChatMessage{
		Channel:           whereToWrite,
		Text:              whatToWrite,
		AsUser:            true,
		Ts:                time.Now(),
		DictionaryMessage: dto.DictionaryMessage{},
		OriginalMessage:   dto.BaseOriginalMessage{},
	})

	if err != nil {
		log.Logger().AddError(err).Msg("Failed to send post-answer for selected event")
		message.Text = fmt.Sprintf("Failed to send the message to the channel %s.\nReason: ```%s```", currentConversation.Scenario.RequiredVariables[1].Value, err.Error())
		return message, err
	}

	//This answer will be show once the event get triggered.
	//Leave message.Text empty, once you need to not show the message, once this event get triggered.
	message.Text = "Done"

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

// Install method for installation of event
func (e EventStruct) Install() error {
	log.Logger().Debug().
		Str("event_name", EventName).
		Str("event_version", EventVersion).
		Msg("Triggered event installation")

	if err := container.C.Dictionary.InstallNewEventScenario(database.EventScenario{
		EventName:    EventName,
		EventVersion: EventVersion,
		Questions: []database.Question{
			{
				Question:      "write a message",
				QuestionRegex: "(?i)(write a message)",
			},
			{
				Question:      "write message",
				QuestionRegex: "(?i)(write message)",
			},
		},
		RequiredVariables: []database.ScenarioVariable{
			{
				Question: stepMessage,
			},
			{
				Question: stepChannel,
			},
		},
	}); err != nil {
		return err
	}

	return nil
}

// Update for event update actions
func (e EventStruct) Update() error {
	return nil
}

func removeCurrentUserFromTheMessage(message string) string {
	regexString := fmt.Sprintf("(?im)(<@%s>)", container.C.Config.MessagesAPIConfig.BotUserID)
	re := regexp.MustCompile(regexString)
	result := re.ReplaceAllString(message, ``)

	return result
}
