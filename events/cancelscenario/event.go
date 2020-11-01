package cancelscenario

import (
	"github.com/sharovik/devbot/internal/helper"
	"github.com/sharovik/devbot/internal/service/base"

	"github.com/sharovik/devbot/internal/log"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
)

const (
	//EventName the name of the event
	EventName = "cancelscenario"

	//EventVersion the version of the event
	EventVersion = "1.0.0"

	helpMessage = "Write me ```stop conversation #channel-name|@username``` and I will stop any conversation which is started for it."

	//The migrations folder, which can be used for event installation or for event update
	migrationDirectoryPath = "./events/cancelscenario/migrations"

	regexChannel = `(?im)(?:[<#@]|(?:&lt;))(\w+)(?:[|>]|(?:&gt;))`
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

	channel := extractChannelName(message.OriginalMessage.Text)
	if channel == "" {
		message.Text = "Please specify the channel name."
		return message, nil
	}

	base.DeleteConversation(channel)

	//This answer will be show once the event get triggered.
	//Leave message.Text empty, once you need to not show the message, once this event get triggered.
	message.Text = "Done."
	return message, nil
}

//Install method for installation of event
func (e EventStruct) Install() error {
	log.Logger().Debug().
		Str("event_name", EventName).
		Str("event_version", EventVersion).
		Msg("Triggered event installation")

	return container.C.Dictionary.InstallEvent(
		EventName,      //We specify the event name which will be used for scenario generation
		EventVersion,   //This will be set during the event creation
		"stop conversation", //Actual question, which system will wait and which will trigger our event
		"Ok, will do it now.",
		"(?i)(stop conversation)", //Optional field. This is regular expression which can be used for question parsing.
		"",                 //Optional field. This is a regex group and it can be used for parsing the match group from the regexp result
	)
}

//Update for event update actions
func (e EventStruct) Update() error {
	return container.C.Dictionary.RunMigrations(migrationDirectoryPath)
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
