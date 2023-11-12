package conversation

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	_time "github.com/sharovik/devbot/internal/service/time"

	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/dto"
)

// Conversation the conversation object which contains the information about the scenario selected for the conversation and the last question asked by the customer.
type Conversation struct {
	ScenarioID             int64
	EventID                int64
	ScenarioQuestionID     int64
	Question               string
	Channel                string
	EventReadyToBeExecuted bool
	LastQuestion           dto.BaseChatMessage
	ReactionType           string
	Scenario               database.EventScenario
}

// currentConversations contains the list of current open conversations, where each key is a channel ID and the value is the last channel question
var currentConversations = map[string]Conversation{}

const openConversationTimeout = time.Second * 600

// GetCurrentConversations returns the list of current open conversations
func GetCurrentConversations() map[string]Conversation {
	return currentConversations
}

// AddConversation add the new conversation to the list of open conversations. This will be used for scenarios build
func AddConversation(scenario database.EventScenario, message dto.BaseChatMessage) {
	conversation := Conversation{
		ScenarioID:         message.DictionaryMessage.ScenarioID,
		EventID:            message.DictionaryMessage.EventID,
		Question:           message.DictionaryMessage.Question,
		ScenarioQuestionID: message.DictionaryMessage.QuestionID,
		LastQuestion:       message,
		ReactionType:       message.DictionaryMessage.ReactionType,
		Scenario:           scenario,
	}

	currentConversations[message.Channel] = conversation
}

// SetLastQuestion sets the last question to the current conversation
func SetLastQuestion(message dto.BaseChatMessage) {
	conv := GetConversation(message.Channel)
	conv.LastQuestion = message

	currentConversations[message.Channel] = conv
}

// GetConversation method retrieve the conversation for selected channel
func GetConversation(channel string) Conversation {
	if conversation, ok := currentConversations[channel]; ok {
		return conversation
	}

	return Conversation{}
}

// MarkAsReadyEventToBeExecuted method set the conversation event ready to be executed
func MarkAsReadyEventToBeExecuted(channel string) {
	conversation := currentConversations[channel]
	conversation.EventReadyToBeExecuted = true

	currentConversations[channel] = conversation
}

// CleanUpExpiredMessages removes the messages from the CleanUpExpiredMessages map object, which are expired
func CleanUpExpiredMessages() {
	currentTime := _time.Service.Now()

	for channel, conversation := range GetCurrentConversations() {
		elapsed := time.Duration(currentTime.Sub(conversation.LastQuestion.Ts).Nanoseconds())
		if elapsed >= openConversationTimeout {
			FinaliseConversation(channel)
		}
	}
}

// FinaliseConversation method delete the conversation for selected channel
func FinaliseConversation(channel string) {
	if _, ok := currentConversations[channel]; !ok {
		return
	}

	delete(currentConversations, channel)
}

// getStopScenarioWords method returns the stop words, which will be used for identification if we need to stop the scenario.
func getStopScenarioWords() []string {
	var stopPhrases []string

	for _, text := range []string{
		"stop!",
		"stop scenario!",
		"exit",
		"stop",
		"cancel",
	} {
		modifiedText := fmt.Sprintf("(%s)", text)
		stopPhrases = append(stopPhrases, modifiedText)
	}

	return stopPhrases
}

// IsScenarioStopTriggered method checks if the scenario stop action was triggered
func IsScenarioStopTriggered(text string) bool {
	regexStr := fmt.Sprintf("(?i)%s", strings.Join(getStopScenarioWords(), "|"))
	regex, err := regexp.Compile(regexStr)
	if err != nil {
		return false
	}

	matches := regex.FindStringSubmatch(text)

	return len(matches) != 0
}
