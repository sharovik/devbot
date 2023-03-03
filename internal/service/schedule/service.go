package schedule

import (
	"fmt"
	"github.com/sharovik/devbot/internal/config"
	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/dto/databasedto"
	"github.com/sharovik/devbot/internal/dto/event"
	"github.com/sharovik/devbot/internal/log"
	"github.com/sharovik/devbot/internal/service/message/conversation"
	"github.com/sharovik/orm/clients"
	cdto "github.com/sharovik/orm/dto"
	"github.com/sharovik/orm/query"
	"strings"
	"time"
)

// Service schedule service struct
type Service struct {
	Config        config.Config
	DB            clients.BaseClientInterface
	DefinedEvents map[string]event.DefinedEventInterface
}

const (
	//VariablesDelimiter global variables
	VariablesDelimiter = ";"
)

var (
	S            Service
	toBeExecuted = map[string][]Item{}
)

// Item the item struct for schedule object
type Item struct {
	ID           int
	Author       string
	Channel      string
	ScenarioID   int64
	EventID      int64
	ReactionType string
	Variables    string //; separated
	Scenario     database.EventScenario
	ExecuteAt    ExecuteAt
	//IsRepeatable if it is set to true, that means we want to repeat it
	IsRepeatable bool
}

func InitS(cfg config.Config, db clients.BaseClientInterface, definedEvents map[string]event.DefinedEventInterface) {
	S = Service{
		Config:        cfg,
		DB:            db,
		DefinedEvents: definedEvents,
	}
}

// Run runs the schedule service in goroutine
func (s *Service) Run() (err error) {
	log.Logger().Debug().Msg("Start schedule service")
	go func() {
		for {
			time.Sleep(time.Second)
			s.triggerEvents()
		}
	}()

	log.Logger().Debug().Msg("Finished schedule service")

	return err
}

func alreadyExists(timeStr string, item Item) bool {
	if len(toBeExecuted[timeStr]) == 0 {
		return false
	}

	//Check if the entry already exists. If so, false will be returned
	for _, entry := range toBeExecuted[timeStr] {
		if generateItemID(entry) == generateItemID(item) {
			//it's already exists
			return true
		}
	}

	return false
}

func (s *Service) triggerEvents() {
	now := time.Now()
	for _, item := range s.getSchedules() {
		if !item.IsRepeatable && now.After(item.ExecuteAt.getDatetime()) {
			if alreadyExists(item.ExecuteAt.getDatetime().Format(timeFormat), item) {
				continue
			}

			toBeExecuted[now.Format(timeFormat)] = append(toBeExecuted[now.Format(timeFormat)], item)
			continue
		}

		if alreadyExists(item.ExecuteAt.getDatetime().Format(timeFormat), item) {
			continue
		}

		toBeExecuted[item.ExecuteAt.getDatetime().Format(timeFormat)] = append(toBeExecuted[item.ExecuteAt.getDatetime().Format(timeFormat)], item)
	}

	for _, item := range toBeExecuted[now.Format(timeFormat)] {
		s.trigger(item)
	}

	delete(toBeExecuted, now.Format(timeFormat))
}

func (s *Service) trigger(item Item) {
	scenario := database.EventScenario{
		ID:        item.ScenarioID,
		EventName: item.ReactionType,
		EventID:   item.EventID,
	}

	for _, variable := range strings.Split(item.Variables, VariablesDelimiter) {
		scenario.RequiredVariables = append(scenario.RequiredVariables, database.ScenarioVariable{
			Value: variable,
		})
	}

	conversation.AddConversation(scenario, dto.BaseChatMessage{
		Channel: item.Channel,
		AsUser:  true,
		Ts:      time.Now(),
		DictionaryMessage: dto.DictionaryMessage{
			ScenarioID:   item.ScenarioID,
			EventID:      item.EventID,
			ReactionType: item.ReactionType,
		},
		OriginalMessage: dto.BaseOriginalMessage{},
	})

	go func() {
		if _, err := s.DefinedEvents[item.ReactionType].Execute(conversation.GetConversation(item.Channel).LastQuestion); err != nil {
			log.Logger().AddError(err).Msg("Failed to execute event")
		}

		conversation.FinaliseConversation(item.Channel)
	}()

	if !item.IsRepeatable {
		q := new(clients.Query).Delete().From(databasedto.SchedulesModel).Where(query.Where{
			First:    "id",
			Operator: "=",
			Second: query.Bind{
				Field: "id",
				Value: item.ID,
			},
		})
		if _, err := s.DB.Execute(q); err != nil {
			log.Logger().
				AddError(err).
				Int("item_id", item.ID).
				Msg("Failed to delete scheduled item from the database")
		}
	}
}

func (s *Service) Schedule(item Item) (err error) {
	model := databasedto.SchedulesModel
	model.AddModelField(cdto.ModelField{
		Name:  "author",
		Value: item.Author,
	})
	model.AddModelField(cdto.ModelField{
		Name:  "channel",
		Value: item.Channel,
	})
	model.AddModelField(cdto.ModelField{
		Name:  "scenario_id",
		Value: item.ScenarioID,
	})
	model.AddModelField(cdto.ModelField{
		Name:  "event_id",
		Value: item.ScenarioID,
	})
	model.AddModelField(cdto.ModelField{
		Name:  "reaction_type",
		Value: item.ReactionType,
	})
	model.AddModelField(cdto.ModelField{
		Name:  "is_repeatable",
		Value: item.IsRepeatable,
	})
	model.AddModelField(cdto.ModelField{
		Name:  "execute_at",
		Value: item.ExecuteAt.toString(),
	})
	model.AddModelField(cdto.ModelField{
		Name:  "variables",
		Value: item.Variables,
	})
	q := new(clients.Query).
		Insert(model)

	_, err = s.DB.Execute(q)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to insert data into database")
		return err
	}

	return nil
}

func (s *Service) getSchedules() (items []Item) {
	q := new(clients.Query).
		Select(databasedto.SchedulesModel.GetColumns()).
		From(databasedto.SchedulesModel)

	result, err := s.DB.Execute(q)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to retrieve schedule list")
		return nil
	}

	for _, item := range result.Items() {
		executeAt, err := new(ExecuteAt).FromString(item.GetField("execute_at").Value.(string))
		if err != nil {
			log.Logger().AddError(err).Msg("Failed to parse execute_at")
			return nil
		}

		isRepeatable := false
		if item.GetField("is_repeatable").Value.(int) == 1 {
			isRepeatable = true
			executeAt.IsRepeatable = isRepeatable
		}

		items = append(items, Item{
			ID:           item.GetField("id").Value.(int),
			Author:       item.GetField("author").Value.(string),
			Channel:      item.GetField("channel").Value.(string),
			ScenarioID:   int64(item.GetField("scenario_id").Value.(int)),
			EventID:      int64(item.GetField("event_id").Value.(int)),
			ReactionType: item.GetField("reaction_type").Value.(string),
			Scenario:     database.EventScenario{},
			ExecuteAt:    executeAt,
			IsRepeatable: isRepeatable,
			Variables:    item.GetField("variables").Value.(string),
		})
	}

	return items
}

func generateItemID(item Item) string {
	return fmt.Sprintf("%d-%s-%s", item.ID, item.Channel, item.ReactionType)
}
