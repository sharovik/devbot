package analiser

import (
	"github.com/sharovik/devbot/internal/service/message/conversation"
	_time "github.com/sharovik/devbot/internal/service/time"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/helper"
	"github.com/sharovik/devbot/internal/log"
)

// Message the message object, from which we will generate the dto.DictionaryMessage
type Message struct {
	Channel string
	User    string
	Text    string
}

// GetDmAnswer retrieves the Dictionary Message Answer
func GetDmAnswer(message Message) (dmAnswer dto.DictionaryMessage, err error) {
	//Now we need to check if there was already opened conversation for this channel
	//If so, then we need to get the Answer from this scenario
	openConversation := conversation.GetConversation(message.Channel)

	//If that was a stop word, we need to cancel the conversation
	IsScenarioStopTriggered := conversation.IsScenarioStopTriggered(message.Text)
	if IsScenarioStopTriggered {
		dmAnswer = dto.DictionaryMessage{
			ScenarioID:            0,
			EventID:               0,
			Answer:                "Ok, no more questions!",
			QuestionID:            0,
			Question:              message.Text,
			Regex:                 "",
			MainGroupIndexInRegex: "",
			ReactionType:          "text",
		}
		conversation.FinaliseConversation(message.Channel)

		return dmAnswer, nil
	}

	//If there was a scenario triggered for this conversation, we trigger the scenario handling logic
	if openConversation.ScenarioID != 0 {
		setAnswerToVariable(message.Text, &openConversation)

		return generateDmForConversation(message, openConversation)
	}

	dmAnswer, err = container.C.Dictionary.FindAnswer(message.Text)
	if err != nil {
		return dto.DictionaryMessage{}, err
	}

	questions, err := getVariableQuestionsByScenarioID(dmAnswer.ScenarioID)
	//if we don't have an error here, then we can proceed with the questions preparing for scenarios
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to get the list of question by the scenarioID")
	}

	isHelpAnswerTriggered, err := helper.HelpMessageShouldBeTriggered(message.Text)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to parse the help template")
	}

	if dmAnswer.ReactionType != "" && isHelpAnswerTriggered {
		dmAnswer.Answer = container.C.DefinedEvents[dmAnswer.ReactionType].Help()
		dmAnswer.IsHelpTriggered = true

		return dmAnswer, nil
	}

	//If the questions amount is more than 1, we need to start the conversation algorithm
	if len(questions) > 1 && !isHelpAnswerTriggered {
		scenario := database.EventScenario{}
		SetScenarioQuestions(&scenario, questions)
		conversation.AddConversation(scenario, dto.BaseChatMessage{
			Channel:           message.Channel,
			Text:              message.Text,
			AsUser:            false,
			Ts:                _time.Service.Now(),
			DictionaryMessage: dmAnswer,
			OriginalMessage: dto.BaseOriginalMessage{
				Text: message.Text,
				User: message.User,
			},
		})

		dmAnswer.Answer = scenario.GetUnAnsweredQuestion()
	}

	return dmAnswer, nil
}

func SetScenarioQuestions(scenario *database.EventScenario, questions []database.QuestionObject) {
	for _, q := range questions {
		scenario.Questions = append(scenario.Questions, database.Question{
			Question: q.Question,
			Answer:   q.Answer,
		})

		if q.IsVariable {
			scenario.RequiredVariables = append(scenario.RequiredVariables, database.ScenarioVariable{
				Question: q.Answer,
			})
		}
	}
}

func getVariableQuestionsByScenarioID(scenarioID int64) (result []database.QuestionObject, err error) {
	questions, err := container.C.Dictionary.GetQuestionsByScenarioID(scenarioID, true)
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to get the list of questions by the scenarioID")
		return result, err
	}

	for _, q := range questions {
		if q.Question != "" {
			continue
		}

		result = append(result, q)
	}

	return result, err
}

func setAnswerToVariable(answer string, openConversation *conversation.Conversation) {
	if len(openConversation.Scenario.RequiredVariables) == 0 {
		return
	}

	for i, variable := range openConversation.Scenario.RequiredVariables {
		if variable.Value != "" || openConversation.LastQuestion.Text != variable.Question {
			continue
		}

		openConversation.Scenario.RequiredVariables[i].Value = answer
		return
	}
}

func generateDmForConversation(message Message, openConversation conversation.Conversation) (dto.DictionaryMessage, error) {
	if len(openConversation.Scenario.RequiredVariables) == 0 {
		return dto.DictionaryMessage{}, nil
	}

	for _, variable := range openConversation.Scenario.RequiredVariables {
		if variable.Value != "" {
			continue
		}

		dmAnswer := dto.DictionaryMessage{
			ScenarioID:   openConversation.ScenarioID,
			EventID:      openConversation.EventID,
			Answer:       variable.Question,
			ReactionType: openConversation.ReactionType,
		}

		return dmAnswer, nil
	}

	conversation.MarkAsReadyEventToBeExecuted(message.Channel)

	return dto.DictionaryMessage{
		ScenarioID:   openConversation.ScenarioID,
		EventID:      openConversation.ScenarioID,
		Answer:       "Ok",
		ReactionType: openConversation.ReactionType,
	}, nil
}
