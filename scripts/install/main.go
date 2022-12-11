package main

import (
	"flag"
	"os"

	"github.com/sharovik/devbot/internal/dto/databasedto"
	"github.com/sharovik/devbot/internal/service/definedevents"
	"github.com/sharovik/orm/clients"
	"github.com/sharovik/orm/dto"

	"github.com/sharovik/devbot/internal/config"
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/log"
	im "github.com/sharovik/devbot/scripts/install/database"
)

const descriptionEventAlias = "The event alias for which Install method will be called"
const envFilePath = "./.env"

var (
	cfg = config.Config{}
	m   = []database.BaseMigrationInterface{
		im.InstallMigration{},
	}
)

func init() {
	configuration, err := config.Init()
	if err != nil {
		panic(err)
	}

	cfg = configuration
	_ = log.Init(cfg.LogConfig)
}

func main() {
	if err := run(); err != nil {
		if err := container.C.Dictionary.CloseDatabaseConnection(); err != nil {
			log.Logger().AddError(err).Msg("Failed to close connection")
		}
	}
}

func triggerMigrations() error {
	//Create migrations table
	q := new(clients.Query).
		Create(databasedto.MigrationModel).
		IfNotExists().
		AddIndex(dto.Index{
			Name:   "migration_version_uindex",
			Target: databasedto.MigrationModel.GetTableName(),
			Key:    "version",
			Unique: true,
		})
	if _, err := container.C.Dictionary.GetDBClient().Execute(q); err != nil {
		log.Logger().AddError(err).Msg("Failed to create migration table. Already exists?")
	}

	for _, migration := range m {
		container.C.MigrationService.SetMigration(migration)
	}

	if err := container.C.MigrationService.RunMigrations(); err != nil {
		return err
	}

	return nil
}

func run() error {
	if err := checkIfEnvFilesExists(); err != nil {
		log.Logger().AddError(err).Msg("Failed check the .env file step")
		return err
	}

	if err := checkIfDatabaseExists(); err != nil {
		log.Logger().AddError(err).Msg("Database check error")
		return err
	}

	eventAlias := parseEventAlias()
	cnt, err := container.Init()
	if err != nil {
		log.Logger().AddError(err).Msg("Failed to init container")
		return err
	}

	container.C = cnt

	definedevents.InitializeDefinedEvents()

	if err = triggerMigrations(); err != nil {
		log.Logger().AddError(err).Msg("Database check error")
		return err
	}

	if eventAlias != "" {
		log.Logger().Info().Msg("Event alias cannot be empty")

		if container.C.DefinedEvents[eventAlias] == nil {
			log.Logger().Info().Msg("Event is not defined in the defined-events")
			return nil
		}

		eventID, err := container.C.Dictionary.FindEventByAlias(eventAlias)
		if err != nil {
			log.Logger().AddError(err).Msg("Failed to check if event exists")
			return err
		}

		if eventID != int64(0) {
			log.Logger().Info().Msg("Event is already installed")
			return nil
		}

		if err := container.C.DefinedEvents[eventAlias].Install(); err != nil {
			log.Logger().AddError(err).Str("event_name", eventAlias).Msg("Failed to install the event.")
		}

		if err := container.C.Dictionary.CloseDatabaseConnection(); err != nil {
			log.Logger().AddError(err).Msg("Failed to close connection")
		}

		log.Logger().Info().Msg("Done")
		return nil
	}

	log.Logger().Debug().Msg("Trying to install all defined events if it's possible")
	for eventAlias, event := range container.C.DefinedEvents {

		eventID, err := container.C.Dictionary.FindEventByAlias(eventAlias)
		if err != nil {
			log.Logger().AddError(err).Msg("Failed to check if event exists")
			return err
		}

		if eventID != int64(0) {
			log.Logger().Info().Int64("event_id", eventID).Str("event_alias", eventAlias).Msg("Event is already installed")
			continue
		}

		if err := event.Install(); err != nil {
			log.Logger().AddError(err).Str("event_alias", eventAlias).Msg("Failed to install the event.")
		}
	}

	log.Logger().Info().Msg("Done")
	return nil
}

func checkIfDatabaseExists() error {
	log.Logger().AppendGlobalContext(map[string]interface{}{
		"database_connection": cfg.Database.Type,
		"database_host":       cfg.Database.Host,
	})
	log.Logger().Info().Msg("Check if the database exists")
	switch cfg.Database.Type {
	case clients.DatabaseTypeSqlite:
		_, err := os.Stat(cfg.Database.Host)
		if err == nil {
			log.Logger().Info().Msg("Database file already exists")
			return nil
		}

		log.Logger().Info().Msg("Creating the database file")

		_, err = os.Create(cfg.Database.Host)
		if err != nil {
			log.Logger().AddError(err).Msg("Failed to create database file")
			return err
		}
	default:
		log.Logger().Info().Msg("No action for selected type of database")
	}

	return nil
}

func checkIfEnvFilesExists() error {
	if _, err := os.Stat(envFilePath); err != nil {
		log.Logger().AddError(err).
			Str("path", envFilePath).
			Msg("The .env file does not exists in selected path")

		return err
	}

	return nil
}

func parseEventAlias() string {
	eventAlias := flag.String("event_alias", "", descriptionEventAlias)
	flag.Parse()

	return *eventAlias
}
