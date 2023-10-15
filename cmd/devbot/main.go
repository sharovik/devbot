package main

import (
	_time "github.com/sharovik/devbot/internal/service/time"
	"time"

	"github.com/sharovik/devbot/internal/service/schedule"

	"github.com/sharovik/devbot/internal/service/definedevents"

	"github.com/sharovik/devbot/internal/config"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/log"
	"github.com/sharovik/devbot/internal/service/message"
)

func init() {
	cnt, err := container.Init()
	if err != nil {
		panic(err)
	}

	container.C = cnt
	definedevents.InitializeDefinedEvents()
	message.InitService()
	schedule.InitS(container.C.Config, container.C.Dictionary.GetDBClient(), container.C.DefinedEvents)
}

const (
	maximumRetries      = 4
	delayBetweenRetries = time.Second * 600 //10 minutes
)

var (
	numberOfRetries = 0
	lastRetry       = time.Now()
)

func run() error {
	if err := schedule.S.Run(); err != nil {
		log.Logger().AddError(err).Msg("Failed to run schedule service")

		return err
	}

	for {
		if err := message.S.InitWebSocketReceiver(); err != nil {
			log.Logger().AddError(err).Msg("Error received during application run")

			if numberOfRetries >= maximumRetries {
				return err
			}

			currentTime := _time.Service.Now()

			//We set to 0 number of retries if there were no any retries since 30 minutes
			elapsed := time.Duration(currentTime.Sub(lastRetry).Nanoseconds())
			if elapsed > delayBetweenRetries {
				numberOfRetries = 0
			}

			numberOfRetries++
			lastRetry = _time.Service.Now()

			log.Logger().AppendGlobalContext(map[string]interface{}{
				"number_retries": numberOfRetries,
				"last_retry":     lastRetry,
			})

			log.Logger().Debug().Msg("Triggered retry")
			if container.C.Config.GetAppEnv() != config.EnvironmentTesting {
				time.Sleep(time.Duration(numberOfRetries) * time.Minute)
			}

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

	log.Logger().StartMessage("DevBot")
	if err := run(); err != nil {
		log.Logger().AddError(err).Msg("Application was interrupted by an error")
	}

	container.C.Terminate()
	log.Logger().FinishMessage("DevBot")
}
