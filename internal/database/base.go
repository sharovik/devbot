package database

import "github.com/sharovik/devbot/internal/dto"

//ConnectionSQLite the sqlite database connection type
const ConnectionSQLite = "sqlite"

//BaseDatabaseInterface interface for base database client
type BaseDatabaseInterface interface {
	InitDatabaseConnection() error
	CloseDatabaseConnection() error
	FindAnswer(message *dto.SlackResponseEventMessage) (dto.DictionaryMessage, error)
	InsertQuestion(question string, answer string, scenarioID int64, questionRegex string, questionRegexGroup string) (int64, error)
}
