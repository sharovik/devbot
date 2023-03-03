package database

import (
	"database/sql"
	"strings"

	"github.com/sharovik/orm/clients"
	cdto "github.com/sharovik/orm/dto"
	cquery "github.com/sharovik/orm/query"

	"github.com/sharovik/devbot/internal/dto"
)

//BaseDatabaseInterface interface for base database client
type BaseDatabaseInterface interface {
	InitDatabaseConnection(cfg clients.DatabaseConfig) error
	GetDBClient() clients.BaseClientInterface
	CloseDatabaseConnection() error
	FindAnswer(message string) (dto.DictionaryMessage, error)
	InsertQuestion(question string, answer string, scenarioID int64, questionRegex string, questionRegexGroup string, isVariable bool) (int64, error)
	InsertScenario(name string, eventID int64) (int64, error)
	FindScenarioByID(scenarioID int64) (int64, error)
	GetLastScenarioID() (int64, error)
	FindEventByAlias(eventAlias string) (int64, error)
	FindEventBy(eventAlias string, version string) (int64, error)
	InsertEvent(alias string, version string) (int64, error)
	FindRegex(regex string) (int64, error)
	InsertQuestionRegex(questionRegex string, questionRegexGroup string) (int64, error)
	GetAllRegex() (map[int64]string, error)
	GetQuestionsByScenarioID(scenarioID int64, isVariable bool) (result []QuestionObject, err error)

	//RunMigrations Should be used for custom event migrations loading
	RunMigrations(path string) error
	IsMigrationAlreadyExecuted(name string) (bool, error)
	MarkMigrationExecuted(name string) error

	//InstallEvent Should be used for your custom event installation. This will create a new event row in the database if previously this row wasn't
	//exists and insert new scenario for specified question and answer
	InstallEvent(eventName string, eventVersion string, question string, answer string, questionRegex string, questionRegexGroup string) error

	//InstallNewEventScenario the method will be used for the new better way of scenario installation
	InstallNewEventScenario(scenario EventScenario) error
}

//QuestionObject used for proper data mapping from questions table
type QuestionObject struct {
	ID           int64
	Question     string
	Answer       string
	ReactionType string
	IsVariable   bool
}

//ScenarioVariable the object for scenario variable
type ScenarioVariable struct {
	Name     string
	Value    string
	Question string
}

//EventScenario the object can be used for the new event scenario installation
type EventScenario struct {
	//EventName name of the event, to which we need connect the scenario
	EventName string

	//ScenarioName name of the scenario
	ScenarioName string

	//EventID the id of existing event in the database
	EventID int64

	//EventVersion version of the event
	EventVersion string

	//ID id of scenario
	ID int64

	//Questions scenario questions list
	Questions []Question

	//RequiredVariables required variables list we expecting for this scenario
	RequiredVariables []ScenarioVariable
}

func (e *EventScenario) VariablesToString() string {
	var result []string
	for _, variable := range e.RequiredVariables {
		result = append(result, variable.Value)
	}

	return strings.Join(result, ";")
}

//Question the scenario question
type Question struct {
	Question      string
	Answer        string
	QuestionRegex string
	QuestionGroup string
}

//GetUnAnsweredQuestion retrieves unanswered question from the list of questions of the scenario
func (e *EventScenario) GetUnAnsweredQuestion() string {
	for _, variable := range e.RequiredVariables {
		if variable.Value != "" {
			continue
		}

		return variable.Question
	}

	return ""
}

func installNewEventScenario(d BaseDatabaseInterface, scenario EventScenario) error {
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

	name := scenario.ScenarioName
	if name == "" {
		name = scenario.EventName
	}

	scenarioID, err := d.InsertScenario(name, eventID)
	if err != nil {
		return err
	}

	for _, q := range scenario.Questions {
		_, err = d.InsertQuestion(q.Question, q.Answer, scenarioID, q.QuestionRegex, q.QuestionRegex, false)
		if err != nil {
			return err
		}
	}

	for _, v := range scenario.RequiredVariables {
		_, err = d.InsertQuestion("", v.Question, scenarioID, "", "", true)
		if err != nil {
			return err
		}
	}

	return nil
}

func getQuestionsByScenarioID(d BaseDatabaseInterface, scenarioID int64, isVariable bool) (result []QuestionObject, err error) {
	query := new(clients.Query).
		Select([]interface{}{
			"questions.id",
			"questions.question",
			"questions.answer",
			"questions.is_variable",
			"events.alias",
		}).
		From(&cdto.BaseModel{TableName: "questions"}).
		Join(cquery.Join{
			Target:    cquery.Reference{Table: "scenarios", Key: "id"},
			With:      cquery.Reference{Table: "questions", Key: "scenario_id"},
			Condition: "=",
			Type:      cquery.InnerJoinType,
		}).
		Join(cquery.Join{
			Target:    cquery.Reference{Table: "events", Key: "id"},
			With:      cquery.Reference{Table: "scenarios", Key: "event_id"},
			Condition: "=",
			Type:      cquery.InnerJoinType,
		}).
		Where(cquery.Where{
			First:    "scenarios.id",
			Operator: "=",
			Second:   cquery.Bind{Field: "scenarios.id", Value: scenarioID},
		}).
		OrderBy("questions.id", cquery.OrderDirectionAsc)

	if isVariable {
		query.Where(cquery.Where{
			First:    "questions.is_variable",
			Operator: "=",
			Second:   cquery.Bind{Field: "is_variable", Value: true},
		})
	}

	res, err := d.GetDBClient().Execute(query)
	if err == sql.ErrNoRows {
		return result, nil
	} else if err != nil {
		return result, err
	}

	for _, item := range res.Items() {
		isVar := false
		if item.GetField("is_variable").Value.(int) == 1 {
			isVar = true
		}

		result = append(result, QuestionObject{
			ID:           int64(item.GetField("id").Value.(int)),
			Question:     item.GetField("question").Value.(string),
			Answer:       item.GetField("answer").Value.(string),
			ReactionType: item.GetField("alias").Value.(string),
			IsVariable:   isVar,
		})
	}

	return result, nil
}
