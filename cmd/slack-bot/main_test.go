package main

import (
	"errors"
	"testing"

	"github.com/sharovik/devbot/internal/service/slack"
	"github.com/sharovik/devbot/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestMain_Retries(t *testing.T) {
	slack.S = mock.MockedSlackService

	mock.ErrorInitWebSocketReceiver = map[int]error{
		0: errors.New("First retry "),
		1: errors.New("Second retry "),
		2: errors.New("Third retry "),
		3: errors.New("Fourth retry "),
		4: errors.New("Fifth retry "),
		5: errors.New("Sixth retry "),
	}

	err := run()
	assert.Error(t, err)
	assert.Equal(t, "Fifth retry ", err.Error())
}
