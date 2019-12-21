package main

import (
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/log"
	"github.com/sharovik/devbot/internal/service/slack"
)

func init() {
	container.C = container.C.Init()
}

func main() {
	log.Logger().StartMessage("SlackBot")

	if err := slack.InitWebSocketReceiver(); err != nil {
		log.Logger().AddError(err).Msg("Error received")
	}

	log.Logger().FinishMessage("SlackBot")
}
