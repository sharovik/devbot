package service

import (
	"database/sql"

	"github.com/sharovik/devbot/internal/log"

	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/dto/databasedto"
	"github.com/sharovik/devbot/internal/service/analiser"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/orm/clients"
	cdto "github.com/sharovik/orm/dto"
	cquery "github.com/sharovik/orm/query"
)

// GenerateDMAnswerForScenarioStep method generates DM object for selected scenario step
func GenerateDMAnswerForScenarioStep(step string) (dto.DictionaryMessage, error) {
	query := new(clients.Query).
		Select([]interface{}{
			"scenarios.id",
			"scenarios.event_id",
			"questions.id as question_id",
			"questions.answer",
			"questions.question",
			"questions_regex.regex as question_regex",
			"questions_regex.regex_group as question_regex_group",
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
				Table: "questions_regex",
				Key:   "id",
			},
			With: cquery.Reference{
				Table: "questions",
				Key:   "regex_id",
			},
			Condition: "=",
			Type:      cquery.LeftJoinType,
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
		Operator: "=",
		Second: cquery.Bind{
			Field: "answer",
			Value: step,
		},
	})

	res, err := container.C.Dictionary.GetDBClient().Execute(query)

	if err == sql.ErrNoRows {
		return dto.DictionaryMessage{}, nil
	} else if err != nil {
		return dto.DictionaryMessage{}, err
	}

	if len(res.Items()) == 0 {
		return dto.DictionaryMessage{}, nil
	}

	//We take first item and use it as the result
	item := res.Items()[0]

	var (
		r  string
		rg string
	)

	if item.GetField("question_regex").Value != nil {
		r = item.GetField("question_regex").Value.(string)
	}

	if item.GetField("question_regex_group").Value != nil {
		rg = item.GetField("question_regex_group").Value.(string)
	}

	return dto.DictionaryMessage{
		ScenarioID:            int64(item.GetField("id").Value.(int)),
		EventID:               int64(item.GetField("event_id").Value.(int)),
		Answer:                item.GetField("answer").Value.(string),
		QuestionID:            int64(item.GetField("question_id").Value.(int)),
		Question:              item.GetField("question").Value.(string),
		Regex:                 r,
		MainGroupIndexInRegex: rg,
		ReactionType:          item.GetField("alias").Value.(string),
	}, nil
}

// PrepareScenario based on scenarioID and reaction type, database.EventScenario will be generated
func PrepareScenario(scenarioID int64, reactionType string) (scenario database.EventScenario, err error) {
	scenario.ID = scenarioID
	questions, err := container.C.Dictionary.GetQuestionsByScenarioID(scenario.ID, true)
	if err != nil {
		return scenario, err
	}

	eventID, err := container.C.Dictionary.FindEventByAlias(reactionType)
	if err != nil {
		return scenario, err
	}

	scenario.EventName = reactionType
	scenario.EventID = eventID
	analiser.SetScenarioQuestions(&scenario, questions)

	return scenario, nil
}

// PrepareEventScenario
// Method generates the scenario object based on the eventId and reaction type received
// eventID - is the ID of event
// reactionType
func PrepareEventScenario(eventID int64, reactionType string) (scenario database.EventScenario, err error) {
	//We are getting scenario
	q := new(clients.Query).Select(databasedto.ScenariosModel.GetColumns()).
		From(databasedto.ScenariosModel).
		Where(cquery.Where{
			First:    "event_id",
			Operator: "=",
			Second: cquery.Bind{
				Field: "event_id",
				Value: eventID,
			},
		}).OrderBy("id", cquery.OrderDirectionAsc).Limit(cquery.Limit{
		From: 0,
		To:   1,
	})
	res, err := container.C.Dictionary.GetDBClient().Execute(q)
	if err != nil {
		return scenario, err
	}

	if len(res.Items()) == 0 {
		log.Logger().Warn().
			Int64("event_id", eventID).
			Str("reaction_type", reactionType).
			Msg("No scenarios found")
		return scenario, nil
	}

	scenario.ID = int64(res.Items()[0].GetField("id").Value.(int))
	variables, err := container.C.Dictionary.GetQuestionsByScenarioID(scenario.ID, true)
	if err != nil {
		return scenario, err
	}
	analiser.SetScenarioQuestions(&scenario, variables)

	scenario.EventName = reactionType
	scenario.EventID = eventID

	return scenario, nil
}
