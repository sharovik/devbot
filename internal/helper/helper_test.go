package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileToBytes(t *testing.T) {
	t.Run("File doesn't exists", func(t *testing.T) {
		var obj struct {
			ID   string
			Name string
		}

		assert.Empty(t, obj)

		bytes, err := FileToBytes("wrong/path/file.json")
		assert.Error(t, err)
		assert.Empty(t, bytes)
	})

	t.Run("File exists", func(t *testing.T) {
		var obj struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}

		assert.Empty(t, obj)

		bytes, err := FileToBytes("../../test/testdata/helper/file.json")
		assert.NoError(t, err)
		assert.NotEmpty(t, bytes)
		assert.Equal(t, []byte(`{"id":1,"name":"John"}`), bytes)
	})
}

func TestHelpMessageShouldBeTriggered(t *testing.T) {
	var text = "this is test"
	isHelpExists, err := HelpMessageShouldBeTriggered(text)
	assert.False(t, isHelpExists)
	assert.NoError(t, err)

	text = "this is test --help"
	isHelpExists, err = HelpMessageShouldBeTriggered(text)
	assert.True(t, isHelpExists)
	assert.NoError(t, err)
}
