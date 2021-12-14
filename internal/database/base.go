package database

import (
	"database/sql"
	"github.com/sharovik/orm/clients"

	"github.com/sharovik/devbot/internal/dto"
)

const (
	//ConnectionSQLite the sqlite database connection type
	ConnectionSQLite = "sqlite"
	//ConnectionMySQL the sqlite database connection type
	ConnectionMySQL = "mysql"
)

type Message struct {
	Channel string
	User string
	Text string
}

//BaseDatabaseInterface interface for base database client
type BaseDatabaseInterface interface {
	InitDatabaseConnection() error
	GetClient() *sql.DB
	GetNewClient() clients.BaseClientInterface
	CloseDatabaseConnection() error
	FindAnswer(message string) (dto.DictionaryMessage, error)
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

	//InstallNewEventScenario the method will be used for the new better way of scenario installation
	InstallNewEventScenario(scenario NewEventScenario) error
}

//QuestionObject used for proper data mapping from questions table
type QuestionObject struct {
	ID           int64
	Question     string
	Answer       string
	ReactionType string
}

//NewEventScenario the object can be used for the new event scenario installation
type NewEventScenario struct {
	EventName string
	EventVersion string
	Questions []Question
}

//Question the scenario question
type Question struct {
	Question string
	Answer string
	QuestionRegex string
	QuestionGroup string
}

func installNewEventScenario(d BaseDatabaseInterface, scenario NewEventScenario) error {
	eventID, err := d.FindEventByAlias(scenario.EventName)
	if err != nil {
		return err
	}

	if eventID == 0 {
		eventID, err = d.InsertEvent(scenario.EventName, scenario.EventVersion)
		if err != nil {
			return err
		}
	}

	scenarioID, err := d.InsertScenario(scenario.EventName, eventID)
	if err != nil {
		return err
	}

	for _, q := range scenario.Questions {
		_, err = d.InsertQuestion(q.Question, q.Answer, scenarioID, q.QuestionRegex, q.QuestionRegex)
		if err != nil {
			return err
		}
	}

	return nil
}