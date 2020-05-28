package main

import (
	"database/sql"
	"flag"
	"github.com/sharovik/devbot/events"
	"github.com/sharovik/devbot/internal/config"
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/log"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	descriptionEventAlias  = "The event alias for which Update method will be called"
	migrationDirectoryPath = "./scripts/update/migrations"
)

var (
	cfg = config.Config{}
	db  *sql.DB
)

func init() {
	_ = log.Init(cfg)
	cfg = config.Init()

	var err error
	switch cfg.DatabaseConnection {
	case database.ConnectionSQLite:
		db, err = sql.Open("sqlite3", cfg.DatabaseHost)
		if err != nil {
			log.Logger().AddError(err).Msg("Failed to open connection")
		}
	default:
		db, err = sql.Open(cfg.DatabaseConnection, cfg.DatabaseHost)
		if err != nil {
			log.Logger().AddError(err).Msg("Failed to open connection")
		}
	}
}

func main() {
	if err := runMigrations(); err != nil {
		log.Logger().AddError(err).Msg("Failed to run migrations")
		return
	}

	eventAlias := parseEventAlias()
	container.C = container.C.Init()

	if eventAlias == "" {
		log.Logger().Debug().Msg("Trying to install all defined events if it's possible")
		for eventAlias, event := range events.DefinedEvents.Events {
			if err := event.Update(); err != nil {
				log.Logger().AddError(err).Str("event_alias", eventAlias).Msg("Failed to update the event.")
			}
		}

		log.Logger().Info().Msg("Done")
		return
	}

	container.C = container.C.Init()
	if events.DefinedEvents.Events[eventAlias] == nil {
		log.Logger().Info().Msg("Event is not defined in the defined-events")
		return
	}

	if err := events.DefinedEvents.Events[eventAlias].Update(); err != nil {
		log.Logger().Info().Msg("Failed to update the event. Error:" + err.Error())
	}

	if err := db.Close(); err != nil {
		log.Logger().AddError(err).Msg("Failed to close the connection")
	}

	log.Logger().Info().Msg("Done")
}

func runMigrations() error {
	if _, err := os.Stat(migrationDirectoryPath); err != nil {
		log.Logger().AddError(err).Str("migration_directory", migrationDirectoryPath).Msg("The migration directory was not found.")
		return err
	}

	var files = map[string]string{}
	err := filepath.Walk(migrationDirectoryPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files[info.Name()] = path
		}

		return nil
	})
	if err != nil {
		log.Logger().AddError(err).Msg("Could not extract files for selected directory")
	}

	for file, filePath := range files {
		migrationData, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Logger().AddError(err).Msg("Failed to open installation file")
			return err
		}

		var id int64
		err = db.QueryRow("select id from migration where version = $1", file).Scan(&id)
		if err == sql.ErrNoRows {
			_, err := db.Exec(string(migrationData))
			if err != nil {
				log.Logger().AddError(err).Str("version", file).Msg("Failed to execute migration")
				return err
			}

			_, err = db.Exec("insert into migration (version) values ($1)", file)
			if err != nil {
				log.Logger().AddError(err).Str("version", file).Msg("Failed to store version into migration database")
				return err
			}

		} else if err != nil {
			log.Logger().AddError(err).Str("version", file).Msg("Failed to check if version already exists in the database")
			return err
		}
	}

	return nil
}

func parseEventAlias() string {
	eventAlias := flag.String("event_alias", "", descriptionEventAlias)
	flag.Parse()

	log.Logger().Info().Str("event_alias", *eventAlias).Msg("Parsed event alias")
	return *eventAlias
}
