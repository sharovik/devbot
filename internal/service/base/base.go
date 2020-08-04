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
	ScenarioID int64
	ScenarioQuestionID int64
	LastQuestion dto.BaseChatMessage
}

//CurrentConversations contains the list of current open conversations, where each key is a channel ID and the value is the last channel question
var CurrentConversations = map[string]Conversation{}

//GetCurrentConversations returns the list of current open conversations
func GetCurrentConversations() map[string]Conversation {
	return CurrentConversations
}

//AddConversation method adds the new conversation to the list of open conversations. This will be used for scenarios build
func AddConversation(channel string, message dto.BaseChatMessage) {
	CurrentConversations[channel] = Conversation{
		ScenarioID:         message.DictionaryMessage.ScenarioID,
		ScenarioQuestionID: 0,
		LastQuestion:       message,
	}
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