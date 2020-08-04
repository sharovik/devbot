package base

import (
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"runtime"
	"testing"
	"time"
)

func init() {
	//We switch pointer to the root directory for control the path from which we need to generate test-data file-paths
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../../")
	_ = os.Chdir(dir)

	container.C = container.C.Init()
}

func TestGetCurrentConversations(t *testing.T) {
	CurrentConversations["_test_channel_"] = Conversation{
		ScenarioID:         0,
		ScenarioQuestionID: 0,
		LastQuestion:       dto.BaseChatMessage{
			Channel:           "_test_channel_",
			Text:              "Testing",
			AsUser:            false,
			Ts:                time.Time{},
			DictionaryMessage: dto.DictionaryMessage{},
			OriginalMessage:   dto.BaseOriginalMessage{},
		},
	}

	CurrentConversations["_test_channel2_"] = Conversation{
		ScenarioID:         0,
		ScenarioQuestionID: 0,
		LastQuestion:       dto.BaseChatMessage{
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
	AddConversation("_test_channel_", dto.BaseChatMessage{
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

	AddConversation("_test_channel_", dto.BaseChatMessage{
		Channel:           "_test_channel_",
		Text:              "Testing",
		AsUser:            false,
		Ts:                now.Add(-time.Second * 600),
		DictionaryMessage: dto.DictionaryMessage{},
		OriginalMessage:   dto.BaseOriginalMessage{},
	})

	AddConversation("_test_channel2_", dto.BaseChatMessage{
		Channel:           "_test_channel2_",
		Text:              "Testing",
		AsUser:            false,
		Ts:                now.Add(time.Second * 600),
		DictionaryMessage: dto.DictionaryMessage{},
		OriginalMessage:   dto.BaseOriginalMessage{},
	})

	CleanUpExpiredMessages()

	assert.NotEmpty(t, CurrentConversations)
	assert.Equal(t, 1, len(CurrentConversations))
	assert.NotEmpty(t, CurrentConversations["_test_channel2_"])
}

func TestDeleteConversation(t *testing.T) {
	AddConversation("_test_channel_", dto.BaseChatMessage{
		Channel:           "_test_channel_",
		Text:              "Testing",
		AsUser:            false,
		Ts:                time.Time{},
		DictionaryMessage: dto.DictionaryMessage{},
		OriginalMessage:   dto.BaseOriginalMessage{},
	})

	assert.NotEmpty(t, CurrentConversations["_test_channel_"])
	DeleteConversation("_test_channel_")
	assert.Empty(t, CurrentConversations["_test_channel_"])
}