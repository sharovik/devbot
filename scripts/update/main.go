package main

import (
	"flag"
	"github.com/sharovik/devbot/internal/service/schedule"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/log"
	"github.com/sharovik/devbot/internal/service/definedevents"
	"github.com/sharovik/devbot/scripts/update/migrations"
)

const (
	descriptionEventAlias = "The event alias for which Update method will be called"
)

var (
	m = []database.BaseMigrationInterface{
		migrations.ExampleMigration{},
		migrations.EventsTriggersHistoryMigration{},
		migrations.UpdateEventsTriggersHistoryMigration{},
	}
)

func main() {
	if err := run(); err != nil {
		if err := container.C.Dictionary.CloseDatabaseConnection(); err != nil {
			log.Logger().AddError(err).Msg("Failed to close connection")
		}
	}
}

func run() error {
	cnt, err := container.Init()
	if err != nil {
		return err
	}

	container.C = cnt

	definedevents.InitializeDefinedEvents()
	schedule.InitS(container.C.Config, container.C.Dictionary.GetDBClient(), container.C.DefinedEvents)

	if err := runMigrations(); err != nil {
		log.Logger().AddError(err).Msg("Failed to run migrations")
		return err
	}

	eventAlias := parseEventAlias()

	if eventAlias != "" {
		if container.C.DefinedEvents[eventAlias] == nil {
			log.Logger().Info().Msg("Event is not defined in the defined-events")
			return nil
		}

		if err := container.C.DefinedEvents[eventAlias].Update(); err != nil {
			log.Logger().Info().Msg("Failed to update the event. Error:" + err.Error())
			return err
		}

		log.Logger().Info().Msg("Done")
		return nil
	}

	log.Logger().Debug().Msg("Trying to update all defined events if it's possible")
	for eventAlias, event := range container.C.DefinedEvents {
		if err := event.Update(); err != nil {
			log.Logger().AddError(err).Str("event_alias", eventAlias).Msg("Failed to update the event.")
		}
	}

	log.Logger().Info().Msg("Done")
	return nil
}

func runMigrations() error {
	for _, migration := range m {
		container.C.MigrationService.SetMigration(migration)
	}

	if err := container.C.MigrationService.RunMigrations(); err != nil {
		return err
	}

	return nil
}

func parseEventAlias() string {
	eventAlias := flag.String("event_alias", "", descriptionEventAlias)
	flag.Parse()

	log.Logger().Info().Str("event_alias", *eventAlias).Msg("Parsed event alias")
	return *eventAlias
}
