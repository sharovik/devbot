package base

import (
	"database/sql"
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"time"
)

//ServiceInterface base interface for messages APIs services
type ServiceInterface interface {
	InitWebSocketReceiver() error
	BeforeWSConnectionStart() error
	ProcessMessage(message interface{}) error
}

//Conversation the conversation object which contains the information about the scenario selected for the conversation and the last question asked by the customer.
type Conversation struct {
	ScenarioID             int64
	ScenarioQuestionID     int64
	EventReadyToBeExecuted bool
	LastQuestion           dto.BaseChatMessage
	ReactionType           string
	//This is map for custom variables of conversation, which can be used in custom events
	Variables []string
}

//CurrentConversations contains the list of current open conversations, where each key is a channel ID and the value is the last channel question
var CurrentConversations = map[string]Conversation{}

//GetCurrentConversations returns the list of current open conversations
func GetCurrentConversations() map[string]Conversation {
	return CurrentConversations
}

//AddConversation method adds the new conversation to the list of open conversations. This will be used for scenarios build
func AddConversation(channel string, questionID int64, message dto.BaseChatMessage, variable string) {
	updatedConversation := Conversation{}
	if CurrentConversations[channel].ScenarioID != int64(0) {
		updatedConversation = CurrentConversations[channel]
		updatedConversation.ScenarioQuestionID = questionID
		updatedConversation.LastQuestion = message
		updatedConversation.ReactionType = message.DictionaryMessage.ReactionType
	} else {
		updatedConversation = Conversation{
			ScenarioID:         message.DictionaryMessage.ScenarioID,
			ScenarioQuestionID: questionID,
			LastQuestion:       message,
			ReactionType:       message.DictionaryMessage.ReactionType,
		}
	}

	updatedConversation = AddConversationVariable(updatedConversation, variable)

	CurrentConversations[channel] = updatedConversation
}

//GetConversation method retrieve the conversation for selected channel
func GetConversation(channel string) Conversation {
	if conversation, ok := CurrentConversations[channel]; ok {
		return conversation
	}

	return Conversation{}
}

//MarkAsReadyEventToBeExecuted method set the conversation event ready to be executed
func MarkAsReadyEventToBeExecuted(conversation Conversation) Conversation {
	conversation.EventReadyToBeExecuted = true
	return conversation
}

//CleanUpExpiredMessages removes the messages from the CleanUpExpiredMessages map object, which are expired
func CleanUpExpiredMessages() {
	currentTime := time.Now()

	for channel, conversation := range GetCurrentConversations() {
		elapsed := time.Duration(currentTime.Sub(conversation.LastQuestion.Ts).Nanoseconds())
		if elapsed >= container.C.Config.OpenConversationTimeout {
			DeleteConversation(channel)
		}
	}
}

//DeleteConversation method delete the conversation for selected channel
func DeleteConversation(channel string) {
	delete(CurrentConversations, channel)
}

//AddConversationVariable method add the variable to the variables of the selected conversation
func AddConversationVariable(conversation Conversation, variable string) Conversation {
	if variable != "" {
		//We save the current list of variables and add the new value there
		conversation.Variables = append(conversation.Variables, variable)
	}

	return conversation
}

func generateDMAnswerForScenarioStep(step string) (dto.DictionaryMessage, error) {
	var (
		id                 int64
		answer             string
		questionID           int64
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

//RunScenario method initialize the scenario for selected step
//scenarioFirstStepString - it's the answer string of your scenario, by this string we will try to find your scenario data.
func RunScenario(scenarioFirstStepString string, message dto.BaseChatMessage) error {
	dmAnswer, err := generateDMAnswerForScenarioStep(scenarioFirstStepString)
	if err != nil {
		return err
	}

	AddConversation(message.Channel, dmAnswer.QuestionID, dto.BaseChatMessage{
		Channel:           message.Channel,
		Text:              scenarioFirstStepString,
		AsUser:            false,
		Ts:                time.Now(),
		DictionaryMessage: dmAnswer,
		OriginalMessage:   dto.BaseOriginalMessage{
			Text:  message.OriginalMessage.Text,
			User:  message.OriginalMessage.User,
			Files: message.OriginalMessage.Files,
		},
	}, "")

	return nil
}