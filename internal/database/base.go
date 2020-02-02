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
	InsertScenario(name string, eventID int64) (int64, error)
	FindScenarioByID(scenarioID int64) (int64, error)
	GetLastScenarioID() (int64, error)
	FindEventByAlias(eventAlias string) (int64, error)
	InsertEvent(alias string) (int64, error)
	FindRegex(regex string) (int64, error)
	InsertQuestionRegex(questionRegex string, questionRegexGroup string) (int64, error)
	GetAllRegex() (map[int64]string, error)
}
