package cancelscenario

import (
	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/helper"
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

	regexChannel = `(?im)(?:[<#@]|(?:&lt;))(\w+)(?:[|>]|(?:&gt;))`
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
	channel := extractChannelName(message.OriginalMessage.Text)
	if channel == "" {
		message.Text = "Please specify the channel name."
		return message, nil
	}

	//This answer will be show once the event get triggered.
	//Leave message.Text empty, once you need to not show the message, once this event get triggered.
	message.Text = "Done."
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
				Question:      "stop conversation",
				Answer:        "Ok, will do it now.",
				QuestionRegex: "(?i)(stop conversation)",
				QuestionGroup: "",
			},
		},
	})
}

// Update for event update actions
func (e EventStruct) Update() error {
	return nil
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
