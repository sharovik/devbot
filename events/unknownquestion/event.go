package unknownquestion

import (
	"database/sql"
	"fmt"
	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/helper"
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

	textStr, err := extractRequestString(message.OriginalMessage.Text)
	words := strings.Split(textStr, " ")

	var foundResults = findResults(words)

	if len(foundResults) == 0 {
		message.Text = "Unfortunately, I don't understand what you mean. Please, have a look on current available `events list`.\nJust type `events list` and I will show you everything what I can."
		return message, nil
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

func implodeDatabaseResults(items []cdto.ModelInterface) string {
	var result string
	for _, item := range items {
		result += fmt.Sprintf("`%s` to trigger event `%s`? Then ask `%s --help` for more details;\n", item.GetField("question").Value.(string), item.GetField("alias").Value.(string), item.GetField("question").Value.(string))
	}

	return result
}

func findResults(words []string) []string {
	var (
		foundResults       []string
		processedEvents = map[string]string{}
	)
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

		var cleanItems []cdto.ModelInterface
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

	return foundResults
}

//Install method for installation of event
func (e EventStruct) Install() error {
	log.Logger().Debug().
		Str("event_name", EventName).
		Str("event_version", EventVersion).
		Msg("Triggered event installation")

	return container.C.Dictionary.InstallNewEventScenario(database.NewEventScenario{
		EventName:    EventName,
		EventVersion: EventVersion,
		Questions:    []database.Question{
			{
				Question:      "similar questions",
				Answer:        "",
				QuestionRegex: "(?im)similar questions",
				QuestionGroup: "",
			},
		},
	})
}

//Update for event update actions
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
		Type: cquery.WhereOrType,
		Second: `"%`+text+`%"`,
	}).Where(cquery.Where{
		First:    "questions.question",
		Operator: "LIKE",
		Type: cquery.WhereOrType,
		Second: `"%`+text+`%"`,
	}).Where(cquery.Where{
		First:    "events.alias",
		Operator: "LIKE",
		Type: cquery.WhereOrType,
		Second: `"%`+text+`%"`,
	}).GroupBy("events.id")

	res, err := container.C.Dictionary.GetNewClient().Execute(query)
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
