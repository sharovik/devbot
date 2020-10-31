package base

import (
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
