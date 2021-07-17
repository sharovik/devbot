package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	err := Init(Config{})

	assert.NoError(t, err)
	assert.Equal(t, true, loggerInstance.initialized)

	Refresh()
}

func TestLogger(t *testing.T) {
	err := Init(Config{})

	assert.NoError(t, err)

	assert.IsType(t, &LoggerInstance{}, Logger())
	Refresh()
}
