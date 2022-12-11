package base

import (
	"github.com/sharovik/devbot/internal/database"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/stretchr/testify/assert"
)

func init() {
	//We switch pointer to the root directory for control the path from which we need to generate test-data file-paths
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../../")
	_ = os.Chdir(dir)

	cnt, err := container.Init()
	if err != nil {
		panic(err)
	}
	container.C = cnt
}

func TestGetCurrentConversations(t *testing.T) {
	currentConversations["_test_channel_"] = Conversation{
		ScenarioID:         0,
		ScenarioQuestionID: 0,
		LastQuestion: dto.BaseChatMessage{
			Channel:           "_test_channel_",
			Text:              "Testing",
			AsUser:            false,
			Ts:                time.Time{},
			DictionaryMessage: dto.DictionaryMessage{},
			OriginalMessage:   dto.BaseOriginalMessage{},
		},
	}

	currentConversations["_test_channel2_"] = Conversation{
		ScenarioID:         0,
		ScenarioQuestionID: 0,
		LastQuestion: dto.BaseChatMessage{
			Channel:           "_test_channel2_",
			Text:              "Testing",
			AsUser:            false,
			Ts:                time.Time{},
			DictionaryMessage: dto.DictionaryMessage{},
			OriginalMessage:   dto.BaseOriginalMessage{},
		},
	}

	list := GetCurrentConversations()
	assert.NotEmpty(t, list)
	assert.NotEmpty(t, list["_test_channel_"])
	assert.NotEmpty(t, list["_test_channel2_"])
}

func TestAddConversation(t *testing.T) {
	scenario := database.EventScenario{}
	AddConversation(scenario, dto.BaseChatMessage{
		Channel:           "_test_channel_",
		Text:              "Testing",
		AsUser:            false,
		Ts:                time.Time{},
		DictionaryMessage: dto.DictionaryMessage{},
		OriginalMessage:   dto.BaseOriginalMessage{},
	})

	list := GetCurrentConversations()
	assert.NotEmpty(t, list)
	assert.NotEmpty(t, list["_test_channel_"])
}

func TestCleanUpExpiredMessages(t *testing.T) {
	now := time.Now()
	scenario := database.EventScenario{}

	AddConversation(scenario, dto.BaseChatMessage{
		Channel:           "_test_channel_",
		Text:              "Testing",
		AsUser:            false,
		Ts:                now.Add(-time.Second * 600),
		DictionaryMessage: dto.DictionaryMessage{},
		OriginalMessage:   dto.BaseOriginalMessage{},
	})

	AddConversation(scenario, dto.BaseChatMessage{
		Channel:           "_test_channel2_",
		Text:              "Testing",
		AsUser:            false,
		Ts:                now.Add(time.Second * 600),
		DictionaryMessage: dto.DictionaryMessage{},
		OriginalMessage:   dto.BaseOriginalMessage{},
	})

	CleanUpExpiredMessages()

	assert.NotEmpty(t, currentConversations)
	assert.Equal(t, 1, len(currentConversations))
	assert.NotEmpty(t, currentConversations["_test_channel2_"])
}

func TestGetConversation(t *testing.T) {
	now := time.Now()
	currentConversations = map[string]Conversation{}
	scenario := database.EventScenario{}

	AddConversation(scenario, dto.BaseChatMessage{
		Channel:           "_test_channel_",
		Text:              "Testing",
		AsUser:            false,
		Ts:                now.Add(time.Second * 600),
		DictionaryMessage: dto.DictionaryMessage{},
		OriginalMessage:   dto.BaseOriginalMessage{},
	})

	conversation := GetConversation("_test_channel_")

	assert.NotEmpty(t, conversation)
	assert.Equal(t, "Testing", conversation.LastQuestion.Text)

	conversation = GetConversation("_test_channel2_")
	assert.Empty(t, conversation)
}

func TestDeleteConversation(t *testing.T) {
	scenario := database.EventScenario{}
	AddConversation(scenario, dto.BaseChatMessage{
		Channel:           "_test_channel_",
		Text:              "Testing",
		AsUser:            false,
		Ts:                time.Time{},
		DictionaryMessage: dto.DictionaryMessage{},
		OriginalMessage:   dto.BaseOriginalMessage{},
	})

	assert.NotEmpty(t, currentConversations["_test_channel_"])
	FinaliseConversation("_test_channel_")
	assert.Empty(t, currentConversations["_test_channel_"])
}
