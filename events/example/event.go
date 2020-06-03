package example

import (
	"fmt"
	"github.com/sharovik/devbot/internal/helper"

	"github.com/sharovik/devbot/internal/log"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
)

const (
	//EventName the name of the event
	EventName = "example"

	//EventVersion the version of the event
	EventVersion = "1.0.1"

	helpMessage = "Ask me `who are you?` and you will see the answer."

	//The migrations folder, which can be used for event installation or for event update
	migrationDirectoryPath = "./events/example/migrations"
)

//ExmplEvent the struct for the event object. It will be used for initialisation of the event in defined-events.go file.
type ExmplEvent struct {
	EventName string
}

//Event - object which is ready to use
var Event = ExmplEvent{
	EventName: EventName,
}

//Execute method which is called by message processor
func (e ExmplEvent) Execute(message dto.BaseChatMessage) (dto.BaseChatMessage, error) {
	isHelpAnswerTriggered, err := helper.HelpMessageShouldBeTriggered(message.OriginalMessage.Text)
	if err != nil {
		log.Logger().Warn().Err(err).Msg("Something went wrong with help message parsing")
	}

	if isHelpAnswerTriggered {
		message.Text = helpMessage
		return message, nil
	}

	//This answer will be show once the event get triggered.
	//Leave message.Text empty, once you need to not show the message, once this event get triggered.
	message.Text = "This is an example of the answer."
	return message, nil
}

//Install method for installation of event
func (e ExmplEvent) Install() error {
	log.Logger().Debug().
		Str("event_name", EventName).
		Str("event_version", EventVersion).
		Msg("Triggered event installation")

	return container.C.Dictionary.InstallEvent(
		EventName,      //We specify the event name which will be used for scenario generation
		EventVersion,   //This will be set during the event creation
		"who are you?", //Actual question, which system will wait and which will trigger our event
		fmt.Sprintf("Hello, my name is %s", container.C.Config.SlackConfig.BotName), //Answer which will be used by the bot
		"(?i)who are you?", //Optional field. This is regular expression which can be used for question parsing.
		"", //Optional field. This is a regex group and it can be used for parsing the match group from the regexp result
	)
}

//Update for event update actions
func (e ExmplEvent) Update() error {
	return container.C.Dictionary.RunMigrations(migrationDirectoryPath)
}
