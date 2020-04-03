package main

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/sharovik/devbot/events"
	"github.com/sharovik/devbot/internal/config"
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/log"
)

const descriptionEventAlias = "The event alias for which Install method will be called"
const databaseInstallationDataSQLitePath = "./scripts/install/database/sqlite.sql"
const envExampleFilePath = "./.env.example"
const envFilePath = "./.env"

var cfg = config.Config{}

func init() {
	//We switch pointer to the root directory for control the path from which we need to generate test-data file-paths
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	_ = os.Chdir(dir)

	_ = log.Init(cfg)
}

func main() {
	if err := checkIfEnvFilesExists(); err != nil {
		log.Logger().AddError(err).Msg("Failed to check the .env file")
		return
	}

	cfg = config.Init()
	if err := checkIfDatabaseExists(); err != nil {
		log.Logger().AddError(err).Msg("Database check error")
		return
	}

	eventAlias := parseEventAlias()
	container.C = container.C.Init()

	if eventAlias != "" {
		log.Logger().Info().Msg("Event alias cannot be empty")

		if events.DefinedEvents.Events[eventAlias] == nil {
			log.Logger().Info().Msg("Event is not defined in the defined-events")
			return
		}

		if err := events.DefinedEvents.Events[eventAlias].Install(); err != nil {
			log.Logger().AddError(err).Str("event_name", eventAlias).Msg("Failed to install the event.")
		}

		log.Logger().Info().Msg("Done")
		return
	}

	for eventAlias, event := range events.DefinedEvents.Events {
		if err := event.Install(); err != nil {
			log.Logger().AddError(err).Str("event_alias", eventAlias).Msg("Failed to install the event.")
		}
	}

	log.Logger().Info().Msg("Done")
}

func checkIfDatabaseExists() error {
	log.Logger().Info().Msg("Check if the database exists")
	switch cfg.DatabaseConnection {
	case database.ConnectionSQLite:
		if _, err := os.Stat(cfg.DatabaseHost); err != nil {
			log.Logger().Info().Msg("We will try to create the database file")

			_, err := os.Stat(databaseInstallationDataSQLitePath)
			if err != nil {
				log.Logger().AddError(err).Msg("Failed to find installation data")
				return err
			}

			_, err = os.Create(cfg.DatabaseHost)
			if err != nil {
				log.Logger().AddError(err).Msg("Failed to create database file")
				return err
			}

			db, err := sql.Open("sqlite3", cfg.DatabaseHost)
			if err != nil {
				log.Logger().AddError(err).Msg("Failed to open connection")
				return err
			}

			installationData, err := ioutil.ReadFile(databaseInstallationDataSQLitePath)
			if err != nil {
				log.Logger().AddError(err).Msg("Failed to open installation file")
				return err
			}

			_, err = db.Exec(string(installationData))
			if err != nil {
				log.Logger().AddError(err).Msg("Failed to install the data")
				return err
			}
		}
	default:
		log.Logger().Warn().Msg("Unfortunately, current version supports only SQLite connection")
	}

	return nil
}

func checkIfEnvFilesExists() error {
	if _, err := os.Stat(envFilePath); err != nil {
		log.Logger().AddError(err).Msg("We will create the .env file from example file")

		if _, err := os.Stat(envExampleFilePath); err != nil {
			log.Logger().AddError(err).Msg("Failed to find example file")

			return err
		}

		envData, err := ioutil.ReadFile(envExampleFilePath)
		if err != nil {
			log.Logger().AddError(err).Msg("Failed to read .env.example file")
			return err
		}

		file, err := os.Create(envFilePath)
		if err != nil {
			log.Logger().AddError(err).Msg("Failed to create .env file")
			return err
		}

		_, err = file.Write(envData)
		if err != nil {
			log.Logger().AddError(err).Msg("Failed write into .env file")
			return err
		}
	}
	return nil
}

func parseEventAlias() string {
	eventAlias := flag.String("event_alias", "", descriptionEventAlias)
	flag.Parse()

	return *eventAlias
}
