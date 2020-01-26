package database

import "github.com/sharovik/devbot/internal/dto"

type BaseDatabaseInterface interface {
	InitDatabaseConnection() error
	FindInDictionary(message *dto.SlackResponseEventMessage) dto.DictionaryMessage
}
