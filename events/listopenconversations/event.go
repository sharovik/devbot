package listopenconversations

import (
	"fmt"

	"github.com/sharovik/devbot/internal/service/message/conversation"

	"github.com/sharovik/devbot/internal/database"
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
	currentConversations := conversation.GetCurrentConversations()
	if len(currentConversations) == 0 {
		message.Text = "There is no open conversations."
		return message, nil
	}

	message.Text = "Here is the list:"
	for _, conv := range conversation.GetCurrentConversations() {
		message.Text += "\n-------"
		message.Text += fmt.Sprintf("\nScenario #%d was triggered in <@%s> chat", conv.ScenarioID, conv.LastQuestion.Channel)
		if len(conv.Scenario.RequiredVariables) == 0 {
			message.Text += "\nAnd there is no answers received yet for that scenario."
		} else {
			message.Text += "\nWith the next filled answers:"
			for _, variable := range conv.Scenario.RequiredVariables {
				message.Text += fmt.Sprintf("\n* %s: `%s`", variable.Question, variable.Value)
			}
		}
		message.Text += "\n-------"
	}

	return message, nil
}

// Install method for installation of event
func (e EventStruct) Install() error {
	log.Logger().Debug().
		Str("event_name", EventName).
		Str("event_version", EventVersion).
		Msg("Triggered event installation")

	return container.C.Dictionary.InstallNewEventScenario(database.EventScenario{
		EventName:    EventName,
		EventVersion: EventVersion,
		Questions: []database.Question{
			{
				Question:      "show open conversations",
				Answer:        "Give me a sec.",
				QuestionRegex: "(?i)(show open conversations)",
				QuestionGroup: "",
			},
		},
	})
}

// Update for event update actions
func (e EventStruct) Update() error {
	return nil
}
