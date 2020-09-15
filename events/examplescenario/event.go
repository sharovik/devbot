package examplescenario

import (
	"fmt"
	"github.com/sharovik/devbot/internal/helper"
	"github.com/sharovik/devbot/internal/service/base"
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

	//The migrations folder, which can be used for event installation or for event update
	migrationDirectoryPath = "./events/examplescenario/migrations"

	regexChannel = `(?im)(?:<#)(\w+)(?:\|)`
)

//ExmplEvent the struct for the event object. It will be used for initialisation of the event in defined-events.go file.
type ExmplSnrEvent struct {
	EventName string
}

//Event - object which is ready to use
var Event = ExmplSnrEvent{
	EventName: EventName,
}

//Execute method which is called by message processor
func (e ExmplSnrEvent) Execute(message dto.BaseChatMessage) (dto.BaseChatMessage, error) {
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

	if len(currentConversation.Variables) == 2 {
		if currentConversation.Variables[0] != "" {
			whatToWrite = currentConversation.Variables[0]
		}

		if currentConversation.Variables[1] != "" {
			whereToWrite = extractChannelName(currentConversation.Variables[1])
		}
	}

	_, _, err = container.C.MessageClient.SendMessage(dto.SlackRequestChatPostMessage{
		Channel:           whereToWrite,
		Text:              whatToWrite,
		AsUser:            true,
		Ts:                time.Now(),
		DictionaryMessage: dto.DictionaryMessage{},
		OriginalMessage:   dto.SlackResponseEventMessage{},
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
func (e ExmplSnrEvent) Install() error {
	log.Logger().Debug().
		Str("event_name", EventName).
		Str("event_version", EventVersion).
		Msg("Triggered event installation")

	if err := container.C.Dictionary.InstallEvent(
		EventName,      //We specify the event name which will be used for scenario generation
		EventVersion,   //This will be set during the event creation
		"write a message", //Actual question, which system will wait and which will trigger our event
		"What I need to write?", //Answer which will be used by the bot
		"", //Optional field. This is regular expression which can be used for question parsing.
		"", //Optional field. This is a regex group and it can be used for parsing the match group from the regexp result
	); err != nil {
		return err
	}

	scenarioID, err := container.C.Dictionary.GetLastScenarioID()
	if err != nil {
		return err
	}

	_, err = container.C.Dictionary.InsertQuestion("", "Where I need to post this message? If it's channel, the channel should be public.", scenarioID, "", "")
	if err != nil {
		return err
	}

	return nil
}

//Update for event update actions
func (e ExmplSnrEvent) Update() error {
	return container.C.Dictionary.RunMigrations(migrationDirectoryPath)
}
