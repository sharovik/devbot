package service

import (
	"database/sql"
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/service/base"
	"github.com/sharovik/orm/clients"
	cdto "github.com/sharovik/orm/dto"
	cquery "github.com/sharovik/orm/query"
	"time"
)

//GenerateDMAnswerForScenarioStep method generates DM object for selected scenario step
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

	res, err := container.C.Dictionary.GetNewClient().Execute(query)

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

//RunScenario runs the selected scenario step
func RunScenario(answer dto.BaseChatMessage, step string) error {
	dmAnswer, err := GenerateDMAnswerForScenarioStep(step)
	if err != nil {
		return err
	}

	base.AddConversation(answer.Channel, dmAnswer.QuestionID, dto.BaseChatMessage{
		Channel:           answer.Channel,
		Text:              step,
		AsUser:            false,
		Ts:                time.Now(),
		DictionaryMessage: dmAnswer,
		OriginalMessage: dto.BaseOriginalMessage{
			Text:  answer.OriginalMessage.Text,
			User:  answer.OriginalMessage.User,
			Files: answer.OriginalMessage.Files,
		},
	}, "")

	return nil
}
