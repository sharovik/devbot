package database

import (
	"database/sql"
	"github.com/sharovik/devbot/internal/config"
	"github.com/sharovik/devbot/internal/dto"
	"os"
)

type SQLiteDictionary struct {
	client sql.DB
	cfg    config.Config
}

//InitDatabaseConnection initialise the database connection
func (d SQLiteDictionary) InitDatabaseConnection() (*sql.DB, error) {
	if _, err := os.Stat(d.cfg.DatabaseHost); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", d.cfg.DatabaseHost)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (d SQLiteDictionary) FindInDictionary(message *dto.SlackResponseEventMessage) dto.DictionaryMessage {
	return dto.DictionaryMessage{}
}
