package container

import (
	"github.com/sharovik/devbot/internal/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMain_Init(t *testing.T) {
	c := C.Init()
	assert.IsType(t, Main{}, c)

	assert.Equal(t, true, c.Config.IsInitialised())
	assert.Equal(t, true, log.IsInitialized())
}
