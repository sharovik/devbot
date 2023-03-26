package unknownquestion

import (
	"database/sql"
	"fmt"
	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/service"
	"github.com/sharovik/devbot/internal/service/message"
	"github.com/sharovik/devbot/internal/service/message/conversation"
	"github.com/sharovik/orm/clients"
	cdto "github.com/sharovik/orm/dto"
	cquery "github.com/sharovik/orm/query"
	"regexp"
	"strings"

	"github.com/sharovik/devbot/internal/log"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
)

const (
	//EventName the name of the event
	EventName = "unknownquestion"

	//EventVersion the version of the event
	EventVersion = "1.0.0"

	helpMessage = "Ask me `similar questions QUESTION_STRING` and I will try to find the similar events and questions in my memory."

	scenarioShouldITriggerThis = `Should I trigger this event?`
)

// EventStruct the struct for the event object. It will be used for initialisation of the event in defined-events.go file.
type EventStruct struct {
}

var (
	// Event - object which is ready to use
	Event                 = EventStruct{}
	selectedConversations = map[string]string{}
)

func (e EventStruct) Help() string {
	return helpMessage
}

func (e EventStruct) Alias() string {
	return EventName
}

// Execute method which is called by message processor
func (e EventStruct) Execute(message dto.BaseChatMessage) (dto.BaseChatMessage, error) {
	conv := conversation.GetConversation(message.Channel)
	if len(conv.Scenario.RequiredVariables) > 0 && selectedConversations[message.Channel] != "" {
		selected := selectedConversations[message.Channel]
		delete(selectedConversations, message.Channel)

		if getAnswer(message) {
			return triggerSelectedScenarioQuestion(message, selected)
		}

		message.Text = "I've got a negative answer from you. I will not trigger this scenario. Perhaps, you may use `events list` command, to see all events."

		return message, nil
	}

	textStr, err := extractRequestString(message.OriginalMessage.Text)
	if err != nil {
		message.Text = "Failed to parse the answer"

		return message, nil
	}

	words := strings.Split(textStr, " ")

	var foundResults, cleanItems = findResults(words)

	switch len(foundResults) {
	case 0:
		message.Text = "Unfortunately, I don't understand what you mean. Please, have a look on current available `events list`.\nJust type `events list` and I will show you everything what I can."
		return message, nil
	case 1:
		return triggerPotentialScenarioQuestion(message, cleanItems[0])
	}

	message.Text = "Maybe you mean:\n"
	for _, str := range foundResults {
		message.Text += str
	}

	//This answer will be show once the event get triggered.
	//Leave message.Text empty, once you need to not show the message, once this event get triggered.
	message.Text += "\n\nOR you also can write `events list` and I will show you all available events."
	return message, nil
}

func getAnswer(message dto.BaseChatMessage) (result bool) {
	conv := conversation.GetConversation(message.Channel)

	//If we already have opened conversation, we will try to get the answer from the required variables
	if conv.Scenario.ID != int64(0) {
		for _, variable := range conv.Scenario.RequiredVariables {
			if "" != variable.Value {
				answer := strings.ToLower(variable.Value)
				switch answer {
				case "yes":
					return true
				default:
					return false
				}
			}
		}
	}

	return false
}

func triggerSelectedScenarioQuestion(msg dto.BaseChatMessage, eventAlias string) (dto.BaseChatMessage, error) {
	eventID, err := container.C.Dictionary.FindEventByAlias(eventAlias)
	if err != nil {
		msg.Text = "Failed to prepare event for triggering. Try again later. Sorry."

		return msg, err
	}

	//We prepare the scenario, with our event name, to make sure we execute the right at the end
	scenario, err := service.PrepareEventScenario(eventID, eventAlias)
	if err != nil {
		msg.Text = "Failed to prepare event-scenario for triggering. Try again later. Sorry."

		return msg, err
	}

	if len(scenario.RequiredVariables) == 0 || len(scenario.Questions) == 0 {
		scenario.ID = 0
	}

	if err = message.TriggerScenario(msg.Channel, scenario, false); err != nil {
		msg.Text = "Failed to trigger selected scenario. Try again later. Sorry."

		return msg, err
	}

	msg.Text = ""

	return msg, nil
}

func triggerPotentialScenarioQuestion(msg dto.BaseChatMessage, item cdto.ModelInterface) (dto.BaseChatMessage, error) {
	scenarioID, err := getDoubleCheckEventScenarioID()
	if err != nil {
		msg.Text = "Failed to trigger the main questions for the schedule scenario"
		return msg, err
	}

	//We prepare the scenario, with our event name, to make sure we execute the right at the end
	scenario, err := service.PrepareScenario(scenarioID, EventName)
	if err != nil {
		msg.Text = "Failed to get the scenario"

		return msg, err
	}

	scenario.RequiredVariables = []database.ScenarioVariable{
		{
			Question: fmt.Sprintf("Do you want me to trigger the `%s` event (later you can ask `%s --help` for more details)?\nPlease, answer yes or no", item.GetField("alias").Value.(string), item.GetField("question").Value.(string)),
		},
	}

	selectedConversations[msg.Channel] = item.GetField("alias").Value.(string)

	if err = message.TriggerScenario(msg.Channel, scenario, false); err != nil {
		msg.Text = "Failed to ask scenario questions"

		return msg, err
	}

	msg.Text = ""

	return msg, nil
}

func getDoubleCheckEventScenarioID() (int64, error) {
	res, err := container.C.Dictionary.FindAnswer("potential_event")
	if err != nil {
		return 0, err
	}

	return res.ScenarioID, nil
}

func implodeDatabaseResults(items []cdto.ModelInterface) string {
	var result string
	for _, item := range items {
		result += fmt.Sprintf("`%s` to trigger event `%s`? Then ask `%s --help` for more details;\n", item.GetField("question").Value.(string), item.GetField("alias").Value.(string), item.GetField("question").Value.(string))
	}

	return result
}

func findResults(words []string) (foundResults []string, cleanItems []cdto.ModelInterface) {
	var processedEvents = map[string]string{}

	for _, word := range words {
		if word == "" {
			continue
		}

		wordRes, err := findPotentialEvents(word)
		if err != nil {
			log.Logger().AddError(err).
				Str("word", word).
				Msg("Failed to find potential event for word")
			continue
		}

		//We remove duplicates
		for _, item := range wordRes {
			if processedEvents[item.GetField("alias").Value.(string)] != "" {
				continue
			}

			cleanItems = append(cleanItems, item)
			processedEvents[item.GetField("alias").Value.(string)] = item.GetField("alias").Value.(string)
		}

		if len(wordRes) > 0 {
			foundResults = append(foundResults, implodeDatabaseResults(cleanItems))
		}
	}

	return
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
				Question:      "similar questions",
				Answer:        "",
				QuestionRegex: "(?im)similar questions",
				QuestionGroup: "",
			},
		},
	}); err != nil {
		return err
	}

	return container.C.Dictionary.InstallNewEventScenario(database.EventScenario{
		EventName:    EventName,
		ScenarioName: "potential_event",
		EventVersion: EventVersion,
		Questions: []database.Question{
			{
				Question:      "potential_event",
				Answer:        "Ok",
				QuestionRegex: "(?i)(potential_question)",
				QuestionGroup: "",
			},
		},
		RequiredVariables: []database.ScenarioVariable{
			{
				Name:     "answer",
				Value:    "",
				Question: scenarioShouldITriggerThis,
			},
		},
	})
}

// Update for event update actions
func (e EventStruct) Update() error {
	return nil
}

func extractRequestString(text string) (result string, err error) {
	re, err := regexp.Compile(`similar questions (.+)`)
	if err != nil {
		return "", err
	}

	matches := re.FindStringSubmatch(text)
	if len(matches) != 2 {
		return "", nil
	}

	if matches[1] == "" {
		return "", nil
	}

	return matches[1], err
}

func findPotentialEvents(text string) (result []cdto.ModelInterface, err error) {
	query := new(clients.Query).
		Select([]interface{}{
			"questions.answer",
			"questions.question",
			"events.alias",
		}).
		From(&cdto.BaseModel{TableName: "questions"}).
		Join(cquery.Join{
			Target: cquery.Reference{
				Table: "scenarios",
				Key:   "id",
			},
			With: cquery.Reference{
				Table: "questions",
				Key:   "scenario_id",
			},
			Condition: "=",
			Type:      cquery.InnerJoinType,
		}).
		Join(cquery.Join{
			Target: cquery.Reference{
				Table: "events",
				Key:   "id",
			},
			With: cquery.Reference{
				Table: "scenarios",
				Key:   "event_id",
			},
			Condition: "=",
			Type:      cquery.LeftJoinType,
		}).Where(cquery.Where{
		First:    "questions.answer",
		Operator: "LIKE",
		Type:     cquery.WhereOrType,
		Second:   `"%` + text + `%"`,
	}).Where(cquery.Where{
		First:    "questions.question",
		Operator: "LIKE",
		Type:     cquery.WhereOrType,
		Second:   `"%` + text + `%"`,
	}).Where(cquery.Where{
		First:    "events.alias",
		Operator: "LIKE",
		Type:     cquery.WhereOrType,
		Second:   `"%` + text + `%"`,
	}).GroupBy("events.id")

	res, err := container.C.Dictionary.GetDBClient().Execute(query)
	if err == sql.ErrNoRows {
		return result, nil
	} else if err != nil {
		return result, err
	}

	if len(res.Items()) == 0 {
		return result, nil
	}

	return res.Items(), nil
}
