package repeatevent

import (
	"fmt"
	"github.com/sharovik/devbot/internal/service/message/conversation"
	"github.com/sharovik/devbot/internal/service/schedule"
	"strings"
	"time"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/dto/databasedto"
	"github.com/sharovik/devbot/internal/log"
	"github.com/sharovik/orm/clients"
	cdto "github.com/sharovik/orm/dto"
	"github.com/sharovik/orm/query"
)

const (
	//EventName the name of the event
	EventName = "repeatevent"

	//EventVersion the version of the event
	EventVersion = "1.0.0"

	helpMessage = "Ask me `repeat` and I will try repeat again the event I previously executed by your command. Eg: You ask me to run something and then you need to rerun in again. You write repeat and I repeat the event."
)

// EventStruct the struct for the event object. It will be used for initialisation of the event in defined-events.go file.
type EventStruct struct {
}

// Event - object which is ready to use
var (
	Event = EventStruct{}
)

func (e EventStruct) Help() string {
	return helpMessage
}

func (e EventStruct) Alias() string {
	return EventName
}

// Execute method which is called by message processor
func (e EventStruct) Execute(message dto.BaseChatMessage) (dto.BaseChatMessage, error) {
	executedEvent, err := lastExecutedEvent(message)
	if err != nil {
		message.Text = "Failed to fetch the last executed events for that channel."
		return message, err
	}

	if executedEvent.GetField("id").Value == nil {
		message.Text = "It looks like you didn't executed any event yet at this channel."
		return message, nil
	}

	answer, err := triggerScenario(executedEvent)
	if err != nil {
		message.Text = fmt.Sprintf("Failed to execute the event.\n```%s```", err)
		return message, err
	}

	message.Text = answer.Text
	return message, nil
}

func getRepeatEventID() (int64, error) {
	eventID, err := container.C.Dictionary.FindEventByAlias(EventName)
	if err != nil {
		return 0, err
	}

	return eventID, nil
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
				Question:      "repeat",
				Answer:        "one sec",
				QuestionRegex: "(?i)repeat",
				QuestionGroup: "",
			},
		},
	})
}

func lastExecutedEvent(message dto.BaseChatMessage) (cdto.ModelInterface, error) {
	currentEventID, err := getRepeatEventID()
	if err != nil {
		return nil, err
	}

	item, err := container.C.Dictionary.GetDBClient().Execute(new(clients.Query).
		Select([]interface{}{}).From(databasedto.EventTriggerHistoryModel).
		Where(query.Where{
			First:    "channel",
			Operator: "=",
			Second: query.Bind{
				Field: "channel",
				Value: message.Channel,
			},
		}).Where(query.Where{
		First:    "user",
		Operator: "=",
		Second: query.Bind{
			Field: "user",
			Value: message.OriginalMessage.User,
		},
	}).Where(query.Where{
		First:    "event_id",
		Operator: "<>",
		Second: query.Bind{
			Field: "event_id",
			Value: currentEventID,
		},
	}).OrderBy("id", query.OrderDirectionDesc).Limit(query.Limit{
		From: 0,
		To:   1,
	}))

	if err != nil {
		return nil, err
	}

	if len(item.Items()) == 0 {
		return nil, nil
	}

	model := databasedto.EventTriggerHistoryModel

	for _, field := range item.Items()[0].GetColumns() {
		switch f := field.(type) {
		case cdto.ModelField:
			model.AddModelField(cdto.ModelField{
				Name:  f.Name,
				Value: f.Value,
			})
			break
		default:
			continue
		}

	}

	return model, nil
}

func getEventAliasByID(eventID int) (string, error) {
	item, err := container.C.Dictionary.GetDBClient().Execute(new(clients.Query).
		Select([]interface{}{"alias"}).From(databasedto.EventModel).Where(query.Where{
		First:    "id",
		Operator: "=",
		Second: query.Bind{
			Field: "id",
			Value: eventID,
		},
	}))
	if err != nil {
		return "", err
	}

	if len(item.Items()) == 0 {
		return "", fmt.Errorf("Failed to find the event alias for selected ID")
	}

	return item.Items()[0].GetField("alias").Value.(string), nil
}

func triggerScenario(item cdto.ModelInterface) (dto.BaseChatMessage, error) {
	eventAlias, err := getEventAliasByID(item.GetField("event_id").Value.(int))
	if err != nil {
		return dto.BaseChatMessage{}, err
	}

	var (
		variables []string
		channel   = item.GetField("channel").Value.(string)
	)

	if "" != item.GetField("variables").Value.(string) {
		scenario := database.EventScenario{
			RequiredVariables: nil,
		}

		variables = strings.Split(item.GetField("variables").Value.(string), schedule.VariablesDelimiter)
		for _, variable := range variables {
			scenario.RequiredVariables = append(scenario.RequiredVariables, database.ScenarioVariable{
				Value: variable,
			})
		}

		dmAnswer := dto.DictionaryMessage{
			ScenarioID:   int64(item.GetField("scenario_id").Value.(int)),
			Question:     item.GetField("command").Value.(string),
			QuestionID:   int64(item.GetField("last_question_id").Value.(int)),
			EventID:      int64(item.GetField("event_id").Value.(int)),
			ReactionType: eventAlias,
		}

		conversation.AddConversation(scenario, dto.BaseChatMessage{
			Channel:           channel,
			Text:              item.GetField("command").Value.(string),
			AsUser:            false,
			Ts:                time.Now(),
			DictionaryMessage: dmAnswer,
			OriginalMessage:   dto.BaseOriginalMessage{},
		})

		return container.C.DefinedEvents[eventAlias].Execute(conversation.GetConversation(channel).LastQuestion)
	}

	return container.C.DefinedEvents[eventAlias].Execute(dto.BaseChatMessage{
		Channel: channel,
		Text:    item.GetField("command").Value.(string),
		AsUser:  false,
		Ts:      time.Time{},
		DictionaryMessage: dto.DictionaryMessage{
			ScenarioID: int64(item.GetField("scenario_id").Value.(int)),
			Question:   "",
			QuestionID: int64(item.GetField("last_question_id").Value.(int)),
			EventID:    int64(item.GetField("event_id").Value.(int)),
		},
		OriginalMessage: dto.BaseOriginalMessage{
			Text:  item.GetField("command").Value.(string),
			User:  item.GetField("user").Value.(string),
			Files: nil,
		},
	})
}

// Update for event update actions
func (e EventStruct) Update() error {
	return nil
}
