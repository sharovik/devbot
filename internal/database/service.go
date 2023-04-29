package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/dto/databasedto"
	"github.com/sharovik/devbot/internal/helper"
	"github.com/sharovik/orm/clients"
	cdto "github.com/sharovik/orm/dto"
	cquery "github.com/sharovik/orm/query"
)

// Dictionary the sqlite dictionary object
type Dictionary struct {
	db clients.BaseClientInterface
}

// GetDBClient method returns the client connection
func (d *Dictionary) GetDBClient() clients.BaseClientInterface {
	return d.db
}

// InitDatabaseConnection initialise the database connection
func (d *Dictionary) InitDatabaseConnection(cfg clients.DatabaseConfig) error {
	var err error
	d.db, err = clients.InitClient(cfg)

	return err
}

// CloseDatabaseConnection method for database connection close
func (d *Dictionary) CloseDatabaseConnection() error {
	return d.db.Disconnect()
}

// FindAnswer used for searching of message in the database
func (d *Dictionary) FindAnswer(message string) (dto.DictionaryMessage, error) {
	var (
		dmAnswer dto.DictionaryMessage
		regexID  int64
		err      error
	)

	//We do that because it can be that we can parse this question by available regex. If so, it will help main query to find the answer for this message
	regexID, err = d.parsedByAvailableRegex(message)
	if err != nil {
		return dto.DictionaryMessage{}, err
	}

	dmAnswer, err = d.answerByQuestionString(message, regexID)
	if err != nil {
		return dto.DictionaryMessage{}, err
	}

	//Finally we parse data by using selected regex in our question
	if dmAnswer.Regex != "" {
		matches := helper.FindMatches(dmAnswer.Regex, message)

		if len(matches) != 0 && dmAnswer.MainGroupIndexInRegex != "" && matches[dmAnswer.MainGroupIndexInRegex] != "" {
			dmAnswer.Answer = fmt.Sprintf(dmAnswer.Answer, matches[dmAnswer.MainGroupIndexInRegex])
		}
	}

	return dmAnswer, nil
}

func (d *Dictionary) parsedByAvailableRegex(question string) (int64, error) {
	availableRegex, err := d.GetAllRegex()
	if err != nil {
		return int64(0), err
	}

	for regexID, regex := range availableRegex {
		matches := helper.FindMatches(regex, question)
		if len(matches) != 0 {
			return regexID, nil
		}
	}

	return 0, nil
}

// answerByQuestionString method retrieves the answer data by selected question string
func (d *Dictionary) answerByQuestionString(questionText string, regexID int64) (dto.DictionaryMessage, error) {
	query := new(clients.Query).
		Select([]interface{}{
			"scenarios.id",
			"scenarios.event_id",
			"questions.id as question_id",
			"questions.answer",
			"questions.question",
			"questions_regex.regex as question_regex",
			"questions_regex.regex_group as question_regex_group",
			"events.alias",
		}).From(&cdto.BaseModel{TableName: "questions"}).
		Join(cquery.Join{
			Target:    cquery.Reference{Table: "scenarios", Key: "id"},
			With:      cquery.Reference{Table: "questions", Key: "scenario_id"},
			Condition: "=",
			Type:      cquery.InnerJoinType,
		}).
		Join(cquery.Join{
			Target:    cquery.Reference{Table: "questions_regex", Key: "id"},
			With:      cquery.Reference{Table: "questions", Key: "regex_id"},
			Condition: "=",
			Type:      cquery.LeftJoinType,
		}).
		Join(cquery.Join{
			Target:    cquery.Reference{Table: "events", Key: "id"},
			With:      cquery.Reference{Table: "scenarios", Key: "event_id"},
			Condition: "=",
			Type:      cquery.LeftJoinType,
		})

	if regexID != 0 {
		query.Where(cquery.Where{
			First:    "questions.regex_id",
			Operator: "=",
			Second: cquery.Bind{
				Field: "regex_id",
				Value: regexID,
			},
		})
	} else {
		query.Where(cquery.Where{
			First:    "questions.question",
			Operator: "LIKE",
			Second: cquery.Bind{
				Field: "question",
				Value: questionText + "%",
			},
		})
	}

	query.
		OrderBy("questions.id", cquery.OrderDirectionAsc).
		Limit(cquery.Limit{From: 0, To: 1})

	res, err := d.db.Execute(query)

	if err == sql.ErrNoRows {
		return dto.DictionaryMessage{}, nil
	} else if err != nil {
		return dto.DictionaryMessage{}, err
	}

	if len(res.Items()) == 0 {
		return dto.DictionaryMessage{}, nil
	}

	//We take first item and use it as the result
	item := res.Items()[0]

	var (
		r  string
		rg string
	)

	if item.GetField("question_regex").Value != nil {
		r = item.GetField("question_regex").Value.(string)
	}

	if item.GetField("question_regex_group").Value != nil {
		rg = item.GetField("question_regex_group").Value.(string)
	}

	return dto.DictionaryMessage{
		ScenarioID:            int64(item.GetField("id").Value.(int)),
		EventID:               int64(item.GetField("event_id").Value.(int)),
		Answer:                item.GetField("answer").Value.(string),
		QuestionID:            int64(item.GetField("question_id").Value.(int)),
		Question:              item.GetField("question").Value.(string),
		Regex:                 r,
		MainGroupIndexInRegex: rg,
		ReactionType:          item.GetField("alias").Value.(string),
	}, nil
}

// InsertScenario used for scenario creation
func (d *Dictionary) InsertScenario(name string, eventID int64) (int64, error) {
	var model = cdto.BaseModel{
		TableName: "scenarios",
		Fields: []interface{}{
			cdto.ModelField{
				Name:  "name",
				Value: name,
			},
			cdto.ModelField{
				Name:  "event_id",
				Value: eventID,
			},
		},
	}

	res, err := d.db.Execute(new(clients.Query).Insert(&model))
	if err != nil {
		return 0, err
	}

	return res.LastInsertID(), nil
}

// FindScenarioByID search scenario by id
func (d *Dictionary) FindScenarioByID(scenarioID int64) (int64, error) {
	query := new(clients.Query).
		Select([]interface{}{"id"}).
		From(&cdto.BaseModel{TableName: "scenarios"}).
		Where(cquery.Where{
			First:    "id",
			Operator: "=",
			Second: cquery.Bind{
				Field: "id",
				Value: scenarioID,
			},
		})
	res, err := d.db.Execute(query)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	if len(res.Items()) == 0 {
		return 0, nil
	}

	return scenarioID, nil
}

// GetLastScenarioID retrieve the last scenario id
func (d *Dictionary) GetLastScenarioID() (int64, error) {
	query := new(clients.Query).
		Select([]interface{}{cdto.ModelField{
			Name: "id",
			Type: cdto.IntegerColumnType,
		}}).
		From(databasedto.ScenariosModel).
		OrderBy("id", cquery.OrderDirectionDesc).
		Limit(cquery.Limit{From: 0, To: 1})
	res, err := d.db.Execute(query)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	if len(res.Items()) == 0 {
		return 0, nil
	}

	item := res.Items()[0]
	return int64(item.GetField("id").Value.(int)), nil
}

// FindEventByAlias search event by alias
func (d *Dictionary) FindEventByAlias(eventAlias string) (int64, error) {
	query := new(clients.Query).
		Select([]interface{}{"id"}).
		From(&cdto.BaseModel{TableName: "events"}).
		Where(cquery.Where{
			First:    "alias",
			Operator: "=",
			Second: cquery.Bind{
				Field: "alias",
				Value: eventAlias,
			},
		})
	res, err := d.db.Execute(query)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	if len(res.Items()) == 0 {
		return 0, nil
	}

	item := res.Items()[0]
	return int64(item.GetField("id").Value.(int)), nil
}

// FindEventBy search event by alias and version
func (d *Dictionary) FindEventBy(eventAlias string, version string) (int64, error) {
	query := new(clients.Query).
		Select([]interface{}{"id"}).
		From(&cdto.BaseModel{TableName: "events"}).
		Where(cquery.Where{
			First:    "alias",
			Operator: "=",
			Second: cquery.Bind{
				Field: "alias",
				Value: eventAlias,
			},
		}).
		Where(cquery.Where{
			First:    "installed_version",
			Operator: "=",
			Second: cquery.Bind{
				Field: "installed_version",
				Value: version,
			},
			Type: cquery.WhereOrType,
		})
	res, err := d.db.Execute(query)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	if len(res.Items()) == 0 {
		return 0, nil
	}

	item := res.Items()[0]
	return int64(item.GetField("id").Value.(int)), nil
}

// InsertEvent used for event creation
func (d *Dictionary) InsertEvent(alias string, version string) (int64, error) {
	var model = cdto.BaseModel{
		TableName: "events",
		Fields: []interface{}{
			cdto.ModelField{
				Name:  "alias",
				Value: alias,
			},
			cdto.ModelField{
				Name:  "installed_version",
				Value: version,
			},
		},
	}

	res, err := d.db.Execute(new(clients.Query).Insert(&model))
	if err != nil {
		return 0, err
	}

	return res.LastInsertID(), nil
}

// InsertQuestion inserts the question into the database
func (d *Dictionary) InsertQuestion(question string, answer string, scenarioID int64, questionRegex string, questionRegexGroup string, isVariable bool) (int64, error) {
	var (
		regexID int64
		err     error
	)
	if questionRegex != "" {
		//We need to find the existing regexID
		regexID, err = d.FindRegex(questionRegex)
		if err != nil {
			return 0, err
		}

		//If we don't have this regex in our database, then we need to add it
		if regexID == 0 {
			regexID, err = d.InsertQuestionRegex(questionRegex, questionRegexGroup)
			if err != nil {
				return 0, err
			}
		}
	}

	var model = cdto.BaseModel{
		TableName: "questions",
		Fields: []interface{}{
			cdto.ModelField{
				Name:  "question",
				Value: question,
			},
			cdto.ModelField{
				Name:  "answer",
				Value: answer,
			},
			cdto.ModelField{
				Name:  "scenario_id",
				Value: scenarioID,
			},
			cdto.ModelField{
				Name:  "is_variable",
				Value: isVariable,
			},
		},
	}

	if regexID != 0 {
		model.AddModelField(cdto.ModelField{
			Name:  "regex_id",
			Value: regexID,
		})
	}

	res, err := d.db.Execute(new(clients.Query).Insert(&model))
	if err != nil {
		return 0, err
	}

	return res.LastInsertID(), nil
}

// FindRegex search regex by regex string
func (d *Dictionary) FindRegex(regex string) (int64, error) {
	query := new(clients.Query).
		Select([]interface{}{"id"}).
		From(&cdto.BaseModel{TableName: "questions_regex"}).
		Where(cquery.Where{
			First:    "regex",
			Operator: "=",
			Second: cquery.Bind{
				Field: "regex",
				Value: regex,
			},
		})
	res, err := d.db.Execute(query)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	if len(res.Items()) == 0 {
		return 0, nil
	}

	item := res.Items()[0]
	return int64(item.GetField("id").Value.(int)), nil
}

// InsertQuestionRegex method insert the regex and returns the regexId. This regex can be connected to the multiple questions
func (d *Dictionary) InsertQuestionRegex(questionRegex string, questionRegexGroup string) (int64, error) {
	var model = cdto.BaseModel{
		TableName: "questions_regex",
		Fields: []interface{}{
			cdto.ModelField{
				Name:  "regex",
				Value: questionRegex,
			},
			cdto.ModelField{
				Name:  "regex_group",
				Value: questionRegexGroup,
			},
		},
	}

	res, err := d.db.Execute(new(clients.Query).Insert(&model))
	if err != nil {
		return 0, err
	}

	return res.LastInsertID(), nil
}

// GetAllRegex method retrieves all available regexs
func (d *Dictionary) GetAllRegex() (res map[int64]string, err error) {
	rows, err := d.db.Execute(new(clients.Query).Select(databasedto.QuestionsRegexModel.GetColumns()).From(&cdto.BaseModel{TableName: "questions_regex"}))
	if err == sql.ErrNoRows {
		return res, nil
	} else if err != nil {
		return res, err
	}

	res = map[int64]string{}
	if len(rows.Items()) == 0 {
		return nil, nil
	}

	for _, item := range rows.Items() {
		res[int64(item.GetField("id").Value.(int))] = item.GetField("regex").Value.(string)
	}

	return res, nil
}

// RunMigrations method for migrations load from specified path
func (d *Dictionary) RunMigrations(pathToFiles string) error {
	if _, err := os.Stat(pathToFiles); os.IsNotExist(err) {
		return nil
	}

	var files = map[string]string{}
	err := filepath.Walk(pathToFiles, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files[info.Name()] = path
		}

		return nil
	})
	if err != nil {
		return err
	}

	var db = d.GetDBClient().GetClient()
	for file, filePath := range files {
		migrationData, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		isMigrationAlreadyExecuted, err := d.IsMigrationAlreadyExecuted(file)
		if err != nil {
			return err
		}

		if isMigrationAlreadyExecuted {
			continue
		}

		_, err = db.Exec(string(migrationData))
		if err != nil {
			return err
		}

		if err := d.MarkMigrationExecuted(file); err != nil {
			return err
		}
	}

	return nil
}

// IsMigrationAlreadyExecuted checks if the migration name was already executed
func (d *Dictionary) IsMigrationAlreadyExecuted(version string) (executed bool, err error) {
	query := new(clients.Query).
		Select([]interface{}{"id"}).
		From(&cdto.BaseModel{TableName: "migration"}).
		Where(cquery.Where{
			First:    "version",
			Operator: "=",
			Second: cquery.Bind{
				Field: "version",
				Value: version,
			},
		})
	rows, err := d.db.Execute(query)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	if len(rows.Items()) == 0 {
		return false, nil
	}

	return true, err
}

// MarkMigrationExecuted marks the selected migration version as executed
func (d *Dictionary) MarkMigrationExecuted(version string) (err error) {
	var model = cdto.BaseModel{
		TableName: "migration",
		Fields: []interface{}{
			cdto.ModelField{
				Name:  "version",
				Value: version,
			},
		},
	}

	_, err = d.db.Execute(new(clients.Query).Insert(&model))
	return
}

// InstallEvent method installs the event(if it wasn't installed before) and creates the scenario for selected event with selected question and answer
// Deprecated: please use InstallNewEventScenario instead
func (d *Dictionary) InstallEvent(eventName string, eventVersion string, question string, answer string, questionRegex string, questionRegexGroup string) error {
	eventID, err := d.FindEventByAlias(eventName)
	if err != nil {
		return err
	}

	if eventID != 0 {
		return nil
	}

	eventID, err = d.InsertEvent(eventName, eventVersion)
	if err != nil {
		return err
	}

	scenarioID, err := d.InsertScenario(eventName, eventID)
	if err != nil {
		return err
	}

	_, err = d.InsertQuestion(question, answer, scenarioID, questionRegex, questionRegexGroup, false)
	if err != nil {
		return err
	}

	return nil
}

// GetQuestionsByScenarioID method retrieves all available questions and answers for selected scenarioID
func (d *Dictionary) GetQuestionsByScenarioID(scenarioID int64, isVariable bool) (result []QuestionObject, err error) {
	return getQuestionsByScenarioID(d, scenarioID, isVariable)
}

// InstallNewEventScenario the method for installing of the new event scenario
func (d *Dictionary) InstallNewEventScenario(scenario EventScenario) error {
	return installNewEventScenario(d, scenario)
}
