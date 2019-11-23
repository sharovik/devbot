package log

import (
	"testing"

	"github.com/sharovik/devbot/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	cfg := new(config.Config)
	err := Init(Config(cfg))

	assert.NoError(t, err)
	assert.Equal(t, true, loggerInstance.initialized)

	Refresh()
}

func TestLogger(t *testing.T) {
	cfg := new(config.Config)
	err := Init(Config(cfg))

	assert.NoError(t, err)

	assert.IsType(t, &LoggerInstance{}, Logger())
	Refresh()
}
