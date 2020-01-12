package main

import (
	"time"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/log"
	"github.com/sharovik/devbot/internal/service/slack"
)

func init() {
	container.C = container.C.Init()
}

const maximumRetries = 4

var (
	numberOfRetries = 0
	lastRetry       = time.Now()
)

func run() error {
	for {
		if err := slack.S.InitWebSocketReceiver(); err != nil {
			log.Logger().AddError(err).Msg("Error received during application run")

			if numberOfRetries >= maximumRetries {
				return err
			}

			numberOfRetries++
			lastRetry = time.Now()

			log.Logger().AppendGlobalContext(map[string]interface{}{
				"number_retries": numberOfRetries,
				"last_retry":     lastRetry,
			})

			log.Logger().Debug().Msg("Triggered retry")
			continue
		}
	}
}

func main() {
	log.Logger().AppendGlobalContext(map[string]interface{}{
		"number_retries":  numberOfRetries,
		"maximum_retries": maximumRetries,
		"started":         lastRetry,
	})

	log.Logger().StartMessage("SlackBot")
	if err := run(); err != nil {
		log.Logger().AddError(err).Msg("Application was interrupted by an error")
	}
	log.Logger().FinishMessage("SlackBot")
}
