package listopenconversations

import (
	"fmt"
	"github.com/sharovik/devbot/internal/helper"
	"github.com/sharovik/devbot/internal/service/base"

	"github.com/sharovik/devbot/internal/log"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
)

const (
	//EventName the name of the event
	EventName = "listopenconversations"

	//EventVersion the version of the event
	EventVersion = "1.0.0"

	helpMessage = "Write me ```show open conversations``` and I will show you the list of open conversations."

	//The migrations folder, which can be used for event installation or for event update
	migrationDirectoryPath = "./events/listopenconversations/migrations"
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

	currentConversations := base.GetCurrentConversations()
	if len(currentConversations) == 0 {
		message.Text = "There is no open conversations."
		return message, nil
	}

	message.Text = "Here is the list:"
	for _, conv := range base.GetCurrentConversations() {
		message.Text += "\n-------"
		message.Text += fmt.Sprintf("\nScenario #%d was triggered in <@%s> chat", conv.ScenarioID, conv.LastQuestion.Channel)
		if len(conv.Variables) == 0 {
			message.Text += fmt.Sprintf("\nAnd there is no answers received yet for that scenario.")
		} else {
			message.Text += "\nWith the next filled answers:"
			for i, variable := range conv.Variables {
				id := i + 1
				message.Text += fmt.Sprintf("\n* #%d: %s", id, variable)
			}
		}
		message.Text += "\n-------"
	}

	return message, nil
}

//Install method for installation of event
func (e EventStruct) Install() error {
	log.Logger().Debug().
		Str("event_name", EventName).
		Str("event_version", EventVersion).
		Msg("Triggered event installation")

	return container.C.Dictionary.InstallEvent(
		EventName,                 //We specify the event name which will be used for scenario generation
		EventVersion,              //This will be set during the event creation
		"show open conversations", //Actual question, which system will wait and which will trigger our event
		"Give me a sec.",
		"(?i)(show open conversations)", //Optional field. This is regular expression which can be used for question parsing.
		"",                              //Optional field. This is a regex group and it can be used for parsing the match group from the regexp result
	)
}

//Update for event update actions
func (e EventStruct) Update() error {
	return nil
}
