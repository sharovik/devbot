package schedule

import (
	"github.com/sharovik/devbot/internal/config"
	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/dto/databasedto"
	"github.com/sharovik/devbot/internal/log"
	"github.com/sharovik/orm/clients"
)

//Service schedule service struct
type Service struct {
	Config config.Config
	db     clients.BaseClientInterface
}

//Item the item struct for schedule object
type Item struct {
	Author       string
	Channel      string
	ScenarioID   int64
	EventID      int64
	ReactionType string
	Scenario     database.EventScenario
	ExecuteAt    string //cron-tab syntax
}

//Run runs the schedule service in goroutine
func (s Service) Run() (err error) {
	//schedules, err :=
	go func() {
		log.Logger().Debug().Msg("Start schedule service")

		for {
			
		}
	}()

	return err
}

func (s Service) getSchedules() (items []Item, err error) {
	q := new(clients.Query).
		Select(databasedto.SchedulesModel.GetColumns()).
		From(databasedto.SchedulesModel)

	result, err := s.db.Execute(q)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to retrieve schedule list")
		return nil, err
	}

	for _, item := range result.Items() {
		items = append(items, Item{
			Author:       item.GetField("author").Value.(string),
			Channel:      item.GetField("channel").Value.(string),
			ScenarioID:   item.GetField("scenario_id").Value.(int64),
			EventID:      item.GetField("event_id").Value.(int64),
			ReactionType: item.GetField("reaction_type").Value.(string),
			Scenario:     database.EventScenario{},
			ExecuteAt:    "",
		})
	}

	return nil, err
}
