package main

import (
	"os"
	"path"
	"runtime"
	"time"

	"github.com/sharovik/devbot/internal/config"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/log"
	"github.com/sharovik/devbot/internal/service/slack"
)

func init() {
	//We switch pointer to the root directory for control the path from which we need to generate test-data file-paths
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	_ = os.Chdir(dir)

	container.C = container.C.Init()
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
	for {
		if err := slack.S.InitWebSocketReceiver(); err != nil {
			log.Logger().AddError(err).Msg("Error received during application run")

			if numberOfRetries >= maximumRetries {
				return err
			}

			currentTime := time.Now()

			//We set to 0 number of retries if there were no any retries since 30 minutes
			elapsed := time.Duration(currentTime.Sub(lastRetry).Nanoseconds())
			if elapsed > delayBetweenRetries {
				numberOfRetries = 0
			}

			numberOfRetries++
			lastRetry = time.Now()

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

	log.Logger().StartMessage("SlackBot")
	if err := run(); err != nil {
		log.Logger().AddError(err).Msg("Application was interrupted by an error")
	}

	container.C.Terminate()
	log.Logger().FinishMessage("SlackBot")
}
