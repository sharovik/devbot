package service

import (
	"database/sql"
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/service/base"
	"time"
)

//GenerateDMAnswerForScenarioStep method generates DM object for selected scenario step
func GenerateDMAnswerForScenarioStep(step string) (dto.DictionaryMessage, error) {
	var (
		id                 int64
		answer             string
		questionID         int64
		question           string
		questionRegex      sql.NullString
		questionRegexGroup sql.NullString
		alias              string
		err                error
	)

	query := `
		select
		s.id,
		q.id as question_id,
		q.answer,
		q.question,
		qr.regex as question_regex,
		qr.regex_group as question_regex_group,
		e.alias
		from questions q
		join scenarios s on q.scenario_id = s.id
		left join questions_regex qr on qr.id = q.regex_id
		left join events e on s.event_id = e.id
		where q.answer = ?
`
	err = container.C.Dictionary.GetClient().QueryRow(query, step).Scan(&id, &questionID, &answer, &question, &questionRegex, &questionRegexGroup, &alias)
	if err == sql.ErrNoRows {
		return dto.DictionaryMessage{}, nil
	} else if err != nil {
		return dto.DictionaryMessage{}, err
	}

	return dto.DictionaryMessage{
		ScenarioID:            id,
		Answer:                answer,
		QuestionID:            questionID,
		Question:              question,
		Regex:                 questionRegex.String,
		MainGroupIndexInRegex: questionRegexGroup.String,
		ReactionType:          alias,
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
