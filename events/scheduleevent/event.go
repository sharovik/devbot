package scheduleevent

import (
	"fmt"
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/helper"
	"github.com/sharovik/devbot/internal/log"
	"github.com/sharovik/devbot/internal/service"
	"github.com/sharovik/devbot/internal/service/message"
	"github.com/sharovik/devbot/internal/service/message/conversation"
	"github.com/sharovik/devbot/internal/service/schedule"
	"strings"
)

const (
	//EventName the name of the event
	EventName = "scheduleevent"

	//EventVersion the version of the event
	EventVersion = "1.0.0"

	supportedTimeFormats = "`YYYY-mm-dd HH:ii`; `DD days`; `HH hours`; `ii minutes`"

	helpMessage = "Ask me `schedule event {event_name} {time_string}` and I will schedule the event. I may ask you the requested event questions. \nThe use-case scenario: you triggered the event, and you want to repeat it, but in a few hours. In that case, ask me: `schedule event {event_name} in 2 hours`. The event will be executed in 2 hours from current time" +
		"\nSupported time formats: " + supportedTimeFormats + ".\nMake sure target event is configured OR does have the scenario questions."

	questionTime       = "When we need to trigger this event? Supported formats: " + supportedTimeFormats
	questionEventAlias = "Which event I need to execute?\nPlease, provide the event alias(use events list command to get it)."
)

// EventStruct the struct for the event object. It will be used for initialisation of the event in defined-events.go file.
type EventStruct struct {
}

type requestedScenario struct {
	Scenario  database.EventScenario
	ExecuteAt schedule.ExecuteAt
}

// Event - object which is ready to use
var (
	Event              = EventStruct{}
	requestedScenarios = map[string]requestedScenario{}
)

func (e EventStruct) Help() string {
	return helpMessage
}

func (e EventStruct) Alias() string {
	return EventName
}

// Execute method which is called by message processor
func (e EventStruct) Execute(message dto.BaseChatMessage) (dto.BaseChatMessage, error) {
	initScenarioService()

	//We schedule the scenario
	if requestedScenarios[message.Channel].Scenario.ID != 0 {
		rScenario := requestedScenarios[message.Channel]
		executeAt := rScenario.ExecuteAt
		if executeAt.IsEmpty() {
			message.Text = "Failed to schedule scenario"
			message.Text += "\nReason: No time specified. Don't know when to trigger."

			return message, nil
		}

		delete(requestedScenarios, message.Channel)

		if err := scheduleRequestedScenario(rScenario, message, executeAt); err != nil {
			message.Text = "Failed to schedule scenario"
			message.Text += fmt.Sprintf("\nReason: %s", err.Error())

			return message, err
		}

		message.Text = fmt.Sprintf("Scenario `%s` scheduled.", rScenario.Scenario.EventName)

		return message, nil
	}

	if !isAllVariablesDefined(message) {
		scenarioID, err := getMainScenarioID()
		if err != nil {
			message.Text = "Failed to trigger the main questions for the schedule scenario"
			return message, err
		}

		if err = askScenarioQuestions(scenarioID, message.Channel); err != nil {
			message.Text = "Failed to trigger the main questions for the schedule scenario"
			return message, err
		}

		message.Text = ""

		return message, nil
	}

	eventType := getReactionType(message)
	scheduleTime := getScheduleTime(message)
	//if event type is defined, we ask questions from that event to collect the answers and use them during schedule.
	if eventType != "" {
		if err := askEventQuestions(eventType, message.Channel); err != nil {
			message.Text = "I cannot ask you the questions from that event. Please, try again."
			return message, err
		}

		rScenario := requestedScenarios[message.Channel]
		if rScenario.ExecuteAt.IsEmpty() {
			rScenario.ExecuteAt = scheduleTime
			requestedScenarios[message.Channel] = rScenario
		}

		message.Text = ""

		return message, nil
	}

	message.Text = ""

	return message, nil
}

func isAllVariablesDefined(message dto.BaseChatMessage) bool {
	reactionType := getReactionType(message)
	scheduleTime := getScheduleTime(message)
	if reactionType == "" || scheduleTime.IsEmpty() {
		return false
	}

	return true
}

func scheduleRequestedScenario(rScenario requestedScenario, message dto.BaseChatMessage, scheduleTime schedule.ExecuteAt) error {
	var variables []string

	for _, value := range conversation.GetConversation(message.Channel).Scenario.RequiredVariables {
		variables = append(variables, value.Value)
	}

	item := schedule.Item{
		Author:       message.OriginalMessage.User,
		Channel:      message.Channel,
		ScenarioID:   rScenario.Scenario.ID,
		EventID:      rScenario.Scenario.EventID,
		ReactionType: rScenario.Scenario.EventName,
		Variables:    strings.Join(variables, schedule.VariablesDelimiter),
		Scenario:     rScenario.Scenario,
		ExecuteAt:    scheduleTime,
		IsRepeatable: rScenario.ExecuteAt.IsRepeatable,
	}

	return schedule.S.Schedule(item)
}

func askEventQuestions(eventType string, channel string) error {
	eventID, err := container.C.Dictionary.FindEventByAlias(eventType)
	if err != nil {
		return err
	}

	//We prepare the scenario, with our event name, to make sure we execute the right at the end
	scenario, err := service.PrepareEventScenario(eventID, EventName)
	if err != nil {
		return err
	}

	if err = message.TriggerScenario(channel, scenario, false); err != nil {
		return err
	}

	//We change back the event type of scenario to the original one, to make sure we schedule the right event
	scenario.EventName = eventType

	requestedScenarios[channel] = requestedScenario{
		Scenario: scenario,
	}

	return nil
}

func askScenarioQuestions(scenarioID int64, channel string) error {
	//We prepare the scenario, with our event name, to make sure we execute the right at the end
	scenario, err := service.PrepareScenario(scenarioID, EventName)
	if err != nil {
		return err
	}

	if err = message.TriggerScenario(channel, scenario, false); err != nil {
		return err
	}

	return nil
}

func getReactionType(message dto.BaseChatMessage) (eventType string) {
	conv := conversation.GetConversation(message.Channel)

	//If we already have opened conversation, we will try to get the answer from the required variables
	if conv.Scenario.ID != int64(0) {
		for _, variable := range conv.Scenario.RequiredVariables {
			if questionEventAlias == variable.Question {
				return variable.Value
			}
		}
	}

	res := helper.FindMatches("(?im)event ([a-zA-Z_]+)", message.OriginalMessage.Text)

	return res["1"]
}

func getScheduleTime(message dto.BaseChatMessage) schedule.ExecuteAt {
	conv := conversation.GetConversation(message.Channel)

	text := message.OriginalMessage.Text
	//If we already have opened conversation, we will try to get the answer from the required variables
	if conv.Scenario.ID != int64(0) {
		for _, variable := range conv.Scenario.RequiredVariables {
			if questionTime == variable.Question {
				text = variable.Value
				break
			}
		}
	}

	r, err := new(schedule.ExecuteAt).FromString(text)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to parse time string")
		return schedule.ExecuteAt{}
	}

	return r
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
				Question:      "schedule event",
				Answer:        "One sec",
				QuestionRegex: "(?i)(schedule event)",
				QuestionGroup: "",
			},
		},
	}); err != nil {
		return err
	}

	return container.C.Dictionary.InstallNewEventScenario(database.EventScenario{
		EventName:    EventName,
		ScenarioName: "schedule_event_scenario",
		EventVersion: EventVersion,
		Questions: []database.Question{
			{
				Question:      "schedule_event_scenario",
				Answer:        "Ok",
				QuestionRegex: "(?i)(schedule_event_scenario)",
				QuestionGroup: "",
			},
		},
		RequiredVariables: []database.ScenarioVariable{
			{
				Name:     "time",
				Value:    "",
				Question: questionTime,
			},
			{
				Name:     "event",
				Value:    "",
				Question: questionEventAlias,
			},
		},
	})
}

// Update for event update actions
func (e EventStruct) Update() error {
	return nil
}

func getMainScenarioID() (int64, error) {
	res, err := container.C.Dictionary.FindAnswer("schedule_event_scenario")
	if err != nil {
		return 0, err
	}

	return res.ScenarioID, nil
}
