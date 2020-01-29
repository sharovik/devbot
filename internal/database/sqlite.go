package database

import (
	"database/sql"
	"fmt"

	//Register the sqlite3 lib
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sharovik/devbot/internal/config"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/helper"
)

//SQLiteDictionary the sqlite dictionary object
type SQLiteDictionary struct {
	client *sql.DB
	Cfg    config.Config
}

//InitDatabaseConnection initialise the database connection
func (d *SQLiteDictionary) InitDatabaseConnection() error {
	if _, err := os.Stat(d.Cfg.DatabaseHost); err != nil {
		return err
	}
	db, err := sql.Open("sqlite3", d.Cfg.DatabaseHost)
	if err != nil {
		return err
	}

	d.client = db
	return nil
}

//CloseDatabaseConnection method for database connection close
func (d *SQLiteDictionary) CloseDatabaseConnection() error {
	return d.client.Close()
}

//FindAnswer used for searching of message in the database
func (d SQLiteDictionary) FindAnswer(message *dto.SlackResponseEventMessage) (dto.DictionaryMessage, error) {
	var (
		id                 int64
		answer             string
		question           string
		questionRegex      string
		questionRegexGroup string
		alias              string
	)

	err := d.client.QueryRow(`
		select
		s.id,
		q.answer,
		q.question,
		q.question_regex,
		q.question_regex_group,
		e.alias
		from questions q
		join scenarios s on q.scenario_id = s.id
		left join events e on s.event_id = e.id
		where q.question like ? order by e.id;
	`, message.Text+"%").Scan(&id, &answer, &question, &questionRegex, &questionRegexGroup, &alias)
	if err == sql.ErrNoRows {
		return dto.DictionaryMessage{}, nil
	} else if err != nil {
		return dto.DictionaryMessage{}, err
	}

	var dmAnswer = dto.DictionaryMessage{
		ScenarioID:            id,
		Answer:                answer,
		Question:              question,
		Regex:                 questionRegex,
		MainGroupIndexInRegex: questionRegexGroup,
		ReactionType:          alias,
	}

	dmAnswer.Answer = answer

	if questionRegex != "" {
		matches := helper.FindMatches(questionRegex, question)

		if len(matches) != 0 && questionRegexGroup != "" && matches[questionRegexGroup] != "" {
			dmAnswer.Answer = fmt.Sprintf(dmAnswer.Answer, matches[questionRegexGroup])
		}
	}

	return dmAnswer, nil
}

//InsertScenario used for scenario creation
func (d SQLiteDictionary) InsertScenario(name string, eventID int64) (int64, error) {
	result, err := d.client.Exec(`insert into scenarios (name, event_id) values ($1, $2)`, name, eventID)
	if err != nil {
		return 0, err
	}

	var scenarioID int64
	scenarioID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return scenarioID, nil
}

//FindScenarioByID search scenario by id
func (d SQLiteDictionary) FindScenarioByID(scenarioID int64) (int64, error) {
	err := d.client.QueryRow("select id from scenarios where id = $1", scenarioID).Scan(&scenarioID)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return scenarioID, nil
}

//GetLastScenarioID retrieve the last scenario id
func (d SQLiteDictionary) GetLastScenarioID() (int64, error) {
	var scenarioID int64
	err := d.client.QueryRow("select id from scenarios order by id desc limit 1").Scan(&scenarioID)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return scenarioID, nil
}

//FindEventByAlias search event by alias
func (d SQLiteDictionary) FindEventByAlias(eventAlias string) (int64, error) {
	var eventID int64
	err := d.client.QueryRow("select id from events where alias = $1", eventAlias).Scan(&eventID)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return eventID, nil
}

//InsertEvent used for event creation
func (d SQLiteDictionary) InsertEvent(alias string) (int64, error) {
	result, err := d.client.Exec(`insert into events (alias) values ($1)`, alias)
	if err != nil {
		return 0, err
	}

	var lastInsertID int64
	lastInsertID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

//InsertQuestion inserts the question into the database
func (d SQLiteDictionary) InsertQuestion(question string, answer string, scenarioID int64, questionRegex string, questionRegexGroup string) (int64, error) {
	result, err := d.client.Exec(`insert into questions (question, answer, scenario_id, question_regex, question_regex_group) values ($1, $2, $3, $4, $5)`,
		question,
		answer,
		scenarioID,
		questionRegex,
		questionRegexGroup,
	)
	if err != nil {
		return 0, err
	}

	var questionID int64
	questionID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return questionID, nil
}
