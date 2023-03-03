package scheduleevent

import (
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/dto/databasedto"
	"github.com/sharovik/devbot/internal/service/analiser"
	"github.com/sharovik/orm/clients"
	"github.com/sharovik/orm/query"
)

type ScenarioService struct {
	Dictionary database.BaseDatabaseInterface
}

var scenarioService ScenarioService

func initScenarioService() {
	scenarioService = ScenarioService{
		Dictionary: container.C.Dictionary,
	}
}

func (s ScenarioService) prepareEventScenario(eventID int64, reactionType string) (scenario database.EventScenario, err error) {
	//We are getting scenario
	q := new(clients.Query).Select(databasedto.ScenariosModel.GetColumns()).
		From(databasedto.ScenariosModel).
		Where(query.Where{
			First:    "event_id",
			Operator: "=",
			Second: query.Bind{
				Field: "event_id",
				Value: eventID,
			},
		}).OrderBy("id", query.OrderDirectionDesc).Limit(query.Limit{
		From: 0,
		To:   1,
	})
	res, err := s.Dictionary.GetDBClient().Execute(q)
	if err != nil {
		return scenario, err
	}

	scenario.ID = int64(res.Items()[0].GetField("id").Value.(int))
	questions, err := s.Dictionary.GetQuestionsByScenarioID(scenario.ID, true)
	if err != nil {
		return scenario, err
	}

	scenario.EventName = reactionType
	scenario.EventID = eventID
	analiser.SetScenarioQuestions(&scenario, questions)

	return scenario, nil
}

func (s ScenarioService) prepareScenario(scenarioID int64, reactionType string) (scenario database.EventScenario, err error) {
	scenario.ID = scenarioID
	questions, err := s.Dictionary.GetQuestionsByScenarioID(scenario.ID, true)
	if err != nil {
		return scenario, err
	}

	eventID, err := s.Dictionary.FindEventByAlias(reactionType)
	if err != nil {
		return scenario, err
	}

	scenario.EventName = reactionType
	scenario.EventID = eventID
	analiser.SetScenarioQuestions(&scenario, questions)

	return scenario, nil
}
