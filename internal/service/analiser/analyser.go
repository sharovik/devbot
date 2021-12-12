package analiser

import (
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/helper"
	"github.com/sharovik/devbot/internal/log"
	"github.com/sharovik/devbot/internal/service/base"
	"time"
)

type PreparedAnswer struct {
	Channel  string
	dmAnswer dto.DictionaryMessage
	Text     string
}

type Message struct {
	Channel string
	User    string
	Text    string
}

func GetDmAnswer(message Message) (dto.DictionaryMessage, error) {
	var (
		err      error
		dmAnswer dto.DictionaryMessage
	)

	//Now we need to check if there was already opened conversation for this channel
	//If so, then we need to get the Answer from this scenario
	openConversation := base.GetConversation(message.Channel)

	IsScenarioStopTriggered := base.IsScenarioStopTriggered(message.Text)
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
		base.DeleteConversation(message.Channel)

		return dmAnswer, nil
	}

	if openConversation.ScenarioID != 0 {
		return generateDmForConversation(message, openConversation)
	}

	dmAnswer, err = container.C.Dictionary.FindAnswer(message.Text)
	if err != nil {
		return dto.DictionaryMessage{}, err
	}

	questions, err := container.C.Dictionary.GetQuestionsByScenarioID(dmAnswer.ScenarioID)
	//if we don't have an error here, then we can proceed with the questions preparing for scenarios
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to get the list of question by the scenarioID")
	}

	isHelpAnswerTriggered, err := helper.HelpMessageShouldBeTriggered(message.Text)
	//If the questions amount is more than 1, we need to start the conversation algorithm
	if len(questions) > 1 && !isHelpAnswerTriggered {
		base.AddConversation(message.Channel, dmAnswer.QuestionID, dto.BaseChatMessage{
			Channel:           message.Channel,
			Text:              message.Text,
			AsUser:            false,
			Ts:                time.Now(),
			DictionaryMessage: dmAnswer,
			OriginalMessage: dto.BaseOriginalMessage{
				Text: message.Text,
				User: message.User,
			},
		}, "")
	}

	return dmAnswer, nil
}

func generateDmForConversation(message Message, openConversation base.Conversation) (dto.DictionaryMessage, error) {
	questions, err := container.C.Dictionary.GetQuestionsByScenarioID(openConversation.ScenarioID)
	//if we don't have an error here, then we can proceed with the questions preparing for scenarios
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to get the list of question by the scenarioID")
		return dto.DictionaryMessage{}, err
	}

	if len(questions) == 0 {
		return dto.DictionaryMessage{}, nil
	}

	//We do questions check only in case of multiple questions attached to the scenario.
	//In other cases we do as we did before
	scenarioNextQuestion := getNextQuestion(openConversation, questions)

	//if we have 0 as ID, then that means we didn't found next question, so we can try to execute the event and have a look what will be
	if scenarioNextQuestion.ID == int64(0) {
		//We mark the scenario, to be executed
		openConversation = base.MarkAsReadyEventToBeExecuted(openConversation)

		//We get the last question object from the scenario and we use it as the answer
		dmAnswer := dto.DictionaryMessage{
			ScenarioID:            openConversation.ScenarioID,
			EventID:               openConversation.ScenarioID,
			Question:              scenarioNextQuestion.Question,
			QuestionID:            scenarioNextQuestion.ID,
			Regex:                 "",
			Answer:                "Ok",
			MainGroupIndexInRegex: "",
			ReactionType:          openConversation.ReactionType,
		}

		//We add the last message to the variables
		openConversation.Variables = append(openConversation.Variables, message.Text)
		base.CurrentConversations[message.Channel] = openConversation

		return dmAnswer, nil
	}

	//In that case we have the last questionID so that means, we need to use the last question here.
	dmAnswer := dto.DictionaryMessage{
		ScenarioID:            openConversation.ScenarioID,
		EventID:               openConversation.EventID,
		Question:              scenarioNextQuestion.Question,
		QuestionID:            scenarioNextQuestion.ID,
		Regex:                 "",
		Answer:                scenarioNextQuestion.Answer,
		MainGroupIndexInRegex: "",
		ReactionType:          scenarioNextQuestion.ReactionType,
	}

	//We also add the new state for this conversation
	base.AddConversation(message.Channel, dmAnswer.QuestionID, dto.BaseChatMessage{
		Channel:           message.Channel,
		Text:              message.Text,
		AsUser:            false,
		Ts:                time.Now(),
		DictionaryMessage: dmAnswer,
		OriginalMessage: dto.BaseOriginalMessage{
			Text: message.Text,
			User: message.User,
		},
	}, message.Text)

	return dmAnswer, nil
}

func triggerUnknownAnswerScenario(message Message) (answer PreparedAnswer, err error) {
	//todo: trigger the unknown answer event execution

	return PreparedAnswer{}, nil
}

func getNextQuestion(openConversation base.Conversation, questions []database.QuestionObject) database.QuestionObject {
	shouldAskNewQuestion := false
	lastQuestion := database.QuestionObject{}
	lastQuestionID := int64(0)
	for _, question := range questions {

		lastQuestionID = question.ID
		if openConversation.ScenarioQuestionID == question.ID {
			shouldAskNewQuestion = true
			continue
		}

		if shouldAskNewQuestion {
			lastQuestion = question
			break
		}

		//In that case we always store here the latest question object
		lastQuestion = question
	}

	//This means there was the last question already triggered
	if lastQuestionID == openConversation.ScenarioQuestionID {
		lastQuestion = database.QuestionObject{}
	}

	return lastQuestion
}
