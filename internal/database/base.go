package database

import (
	"database/sql"

	"github.com/sharovik/devbot/internal/dto"
)

//ConnectionSQLite the sqlite database connection type
const ConnectionSQLite = "sqlite"

//BaseDatabaseInterface interface for base database client
type BaseDatabaseInterface interface {
	InitSQLiteDatabaseConnection() error
	GetClient() *sql.DB
	CloseDatabaseConnection() error
	FindAnswer(message *dto.SlackResponseEventMessage) (dto.DictionaryMessage, error)
	InsertQuestion(question string, answer string, scenarioID int64, questionRegex string, questionRegexGroup string) (int64, error)
	InsertScenario(name string, eventID int64) (int64, error)
	FindScenarioByID(scenarioID int64) (int64, error)
	GetLastScenarioID() (int64, error)
	FindEventByAlias(eventAlias string) (int64, error)
	FindEventBy(eventAlias string, version string) (int64, error)
	InsertEvent(alias string, version string) (int64, error)
	FindRegex(regex string) (int64, error)
	InsertQuestionRegex(questionRegex string, questionRegexGroup string) (int64, error)
	GetAllRegex() (map[int64]string, error)
	GetQuestionsByScenarioID(scenarioID int64) (result []QuestionObject, err error)

	//Should be used for custom event migrations loading
	RunMigrations(path string) error
	IsMigrationAlreadyExecuted(name string) (bool, error)
	MarkMigrationExecuted(name string) error

	//Should be used for your custom event installation. This will create a new event row in the database if previously this row wasn't
	//exists and insert new scenario for specified question and answer
	InstallEvent(eventName string, eventVersion string, question string, answer string, questionRegex string, questionRegexGroup string) error
}

//QuestionObject used for proper data mapping from questions table
type QuestionObject struct {
	ID int64
	Question string
	Answer string
	ReactionType string
}
