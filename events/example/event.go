package example

import (
	"fmt"

	"github.com/sharovik/devbot/internal/database"

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
)

// EventStruct the struct for the event object. It will be used for initialisation of the event in defined-events.go file.
type EventStruct struct {
}

// Event - object which is ready to use
var Event = EventStruct{}

func (e EventStruct) Help() string {
	return helpMessage
}

func (e EventStruct) Alias() string {
	return EventName
}

// Execute method which is called by message processor
func (e EventStruct) Execute(message dto.BaseChatMessage) (dto.BaseChatMessage, error) {
	//This answer will be show once the event get triggered.
	//Leave message.Text empty, once you need to not show the message, once this event get triggered.
	message.Text = "This is an example of the answer."
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
				Question:      "who are you?",
				Answer:        fmt.Sprintf("Hello, my name is %s", container.C.Config.MessagesAPIConfig.BotName),
				QuestionRegex: "(?i)who are you?",
				QuestionGroup: "",
			},
		},
	})
}

// Update for event update actions
func (e EventStruct) Update() error {
	return nil
}
