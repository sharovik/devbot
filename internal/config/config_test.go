package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_GetEnvironment(t *testing.T) {
	var cfg Config

	cfg.appEnv = "testing"
	assert.Equal(t, "testing", cfg.GetAppEnv())
}

func TestInit(t *testing.T) {
	c, err := Init()
	assert.NoError(t, err)
	assert.Equal(t, true, c.initialised)
}
