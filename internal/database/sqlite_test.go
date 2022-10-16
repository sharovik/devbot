package database

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sharovik/devbot/internal/config"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/stretchr/testify/assert"
)

const testSQLiteDatabasePath = "./test/testdata/database/devbot.sqlite"

var (
	cfg             config.Config
	dictionary      SQLiteDictionary
	availableTables = map[string]string{}
)

func init() {
	//We switch pointer to the root directory for control the path from which we need to generate test-data file-paths
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	_ = os.Chdir(dir)
}

func TestSQLiteDictionary_InitDatabaseConnection(t *testing.T) {
	initDatabase()
	cfg.DatabaseHost = "./wrong_path"
	dictionary.Cfg = cfg

	err := dictionary.InitDatabaseConnection()

	assert.Error(t, err)
	assert.Empty(t, dictionary.client.GetClient())

	cfg.DatabaseHost = testSQLiteDatabasePath
	dictionary.Cfg = cfg

	err = dictionary.InitDatabaseConnection()
	assert.NoError(t, err)

	defer dictionary.client.Disconnect()

	checkIfDataCanBeReturned(t)
	dropTestTables(t)
	removeDatabase()
}

func TestSQLiteDictionary_CloseDatabaseConnection(t *testing.T) {
	initDatabase()
	cfg.DatabaseHost = testSQLiteDatabasePath
	dictionary.Cfg = cfg

	err := dictionary.InitDatabaseConnection()
	assert.NoError(t, err)

	checkIfDataCanBeReturned(t)

	dropTestTables(t)
	err = dictionary.CloseDatabaseConnection()
	assert.NoError(t, err)
	removeDatabase()
}

func checkIfDataCanBeReturned(t *testing.T) {
	initDatabase()
	sqlStmt := `
	drop table if exists foo;
	create table foo (id integer not null primary key, name text);
	`
	_, err := dictionary.client.GetClient().Exec(sqlStmt)
	assert.NoError(t, err)

	_, err = dictionary.client.GetClient().Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	assert.NoError(t, err)

	rows, err := dictionary.client.GetClient().Query("select id, name from foo where id = 1")
	assert.NoError(t, err)
	defer rows.Close()

	var id int
	var name string

	for rows.Next() {
		err = rows.Scan(&id, &name)
		assert.NoError(t, err)
	}

	assert.Equal(t, 1, id)
	assert.Equal(t, "foo", name)

	err = rows.Err()
	assert.NoError(t, err)

	sqlStmt = `drop table if exists foo;`
	_, err = dictionary.client.GetClient().Exec(sqlStmt)
	assert.NoError(t, err)
}

func TestSQLiteDictionary_FindAnswer(t *testing.T) {
	initDatabase()
	cfg.DatabaseHost = testSQLiteDatabasePath
	dictionary.Cfg = cfg

	err := dictionary.InitDatabaseConnection()
	assert.NoError(t, err)
	upTestTables(t)
	insertTestData(t)

	goodCases := map[string]string{
		"Hello":              "Hello",
		"Hello world":        "Hello",
		"Say hello to John":  "Hello John",
		"Say hello to Pavel": "Hello Pavel",
	}

	for question, answer := range goodCases {
		var dmAnswer dto.DictionaryMessage
		dmAnswer, err = dictionary.FindAnswer(question)
		assert.NoError(t, err)
		assert.NotEmpty(t, dmAnswer)
		assert.Equal(t, answer, dmAnswer.Answer, question)
		assert.NotEmpty(t, dmAnswer.Question)
		dmAnswer = dto.DictionaryMessage{}
	}

	badCases := map[string]string{
		"Test": "",
	}

	for question := range badCases {
		var dmAnswer dto.DictionaryMessage
		dmAnswer, err = dictionary.FindAnswer(question)
		assert.Empty(t, dmAnswer, question)
	}

	dropTestTables(t)
	dictionary.CloseDatabaseConnection()
	removeDatabase()
}

func TestSQLiteDictionary_FindEventByAlias(t *testing.T) {
	initDatabase()
	cfg.DatabaseHost = testSQLiteDatabasePath
	dictionary.Cfg = cfg

	err := dictionary.InitDatabaseConnection()
	assert.NoError(t, err)
	upTestTables(t)
	insertTestData(t)

	var eventID int64
	eventID, err = dictionary.FindEventByAlias("test")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), eventID)

	eventID, err = dictionary.FindEventByAlias("hello")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), eventID)

	dropTestTables(t)
	dictionary.CloseDatabaseConnection()
	removeDatabase()
}

func TestSQLiteDictionary_FindScenarioById(t *testing.T) {
	initDatabase()
	cfg.DatabaseHost = testSQLiteDatabasePath
	dictionary.Cfg = cfg

	err := dictionary.InitDatabaseConnection()
	assert.NoError(t, err)
	upTestTables(t)
	insertTestData(t)

	var scenarioID int64
	scenarioID, err = dictionary.FindScenarioByID(int64(3))
	assert.NoError(t, err)
	assert.Equal(t, int64(0), scenarioID)

	scenarioID, err = dictionary.FindScenarioByID(int64(1))
	assert.NoError(t, err)
	assert.Equal(t, int64(1), scenarioID)

	dropTestTables(t)
	dictionary.CloseDatabaseConnection()
	removeDatabase()
}

func TestSQLiteDictionary_InsertScenario(t *testing.T) {
	initDatabase()
	cfg.DatabaseHost = testSQLiteDatabasePath
	dictionary.Cfg = cfg

	err := dictionary.InitDatabaseConnection()
	assert.NoError(t, err)
	upTestTables(t)
	insertTestData(t)

	var scenarioID int64
	scenarioID, err = dictionary.InsertScenario("test", int64(1))
	assert.NoError(t, err)
	assert.Equal(t, int64(2), scenarioID)

	scenarioID, err = dictionary.InsertScenario("test", int64(1))
	assert.Error(t, err)
	assert.Equal(t, "UNIQUE constraint failed: scenarios.name", err.Error())

	dropTestTables(t)
	dictionary.CloseDatabaseConnection()
	removeDatabase()
}

func TestSQLiteDictionary_InsertEvent(t *testing.T) {
	initDatabase()
	cfg.DatabaseHost = testSQLiteDatabasePath
	dictionary.Cfg = cfg

	err := dictionary.InitDatabaseConnection()
	assert.NoError(t, err)
	upTestTables(t)
	insertTestData(t)

	var eventID int64
	eventID, err = dictionary.InsertEvent("test", "1.0.0")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), eventID)

	eventID, err = dictionary.InsertEvent("test", "1.0.0")
	assert.Error(t, err)
	assert.Equal(t, "UNIQUE constraint failed: events.alias", err.Error())

	dropTestTables(t)
	dictionary.CloseDatabaseConnection()
	removeDatabase()
}

func TestSQLiteDictionary_InsertQuestion(t *testing.T) {
	initDatabase()
	cfg.DatabaseHost = testSQLiteDatabasePath
	dictionary.Cfg = cfg

	err := dictionary.InitDatabaseConnection()
	assert.NoError(t, err)
	upTestTables(t)
	insertTestData(t)

	var questionID int64
	questionID, err = dictionary.InsertQuestion(
		"Hello bot",
		"Yo",
		int64(1),
		"",
		"",
	)
	assert.NoError(t, err)
	assert.Equal(t, int64(4), questionID)

	questionID, err = dictionary.InsertQuestion(
		"Hello bot",
		"Yo",
		int64(1),
		"",
		"",
	)
	assert.Error(t, err)
	assert.Equal(t, "UNIQUE constraint failed: questions.question", err.Error())

	dropTestTables(t)
	dictionary.CloseDatabaseConnection()
	removeDatabase()
}

func TestSQLiteDictionary_GetQuestionsByScenarioID(t *testing.T) {
	initDatabase()
	cfg.DatabaseHost = testSQLiteDatabasePath
	dictionary.Cfg = cfg

	err := dictionary.InitDatabaseConnection()
	assert.NoError(t, err)
	upTestTables(t)
	insertTestData(t)

	var questionID int64
	questionID, err = dictionary.InsertQuestion(
		"Hello bot",
		"Yo",
		int64(1),
		"",
		"",
	)
	assert.NoError(t, err)
	assert.Equal(t, int64(4), questionID)

	questionID, err = dictionary.InsertQuestion(
		"",
		"test",
		int64(1),
		"",
		"",
	)

	assert.NoError(t, err)

	questions, err := dictionary.GetQuestionsByScenarioID(int64(1))
	assert.NoError(t, err)
	assert.Equal(t, 5, len(questions))

	dropTestTables(t)
	dictionary.CloseDatabaseConnection()
	removeDatabase()
}

func TestSQLiteDictionary_GetLastScenarioID(t *testing.T) {
	initDatabase()
	cfg.DatabaseHost = testSQLiteDatabasePath
	dictionary.Cfg = cfg

	err := dictionary.InitDatabaseConnection()
	assert.NoError(t, err)
	upTestTables(t)
	insertTestData(t)

	var scenarioID int64
	scenarioID, err = dictionary.GetLastScenarioID()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), scenarioID)

	scenarioID, err = dictionary.InsertScenario("test", int64(1))
	assert.NoError(t, err)
	assert.Equal(t, int64(2), scenarioID)

	scenarioID, err = dictionary.GetLastScenarioID()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), scenarioID)

	dropTestTables(t)
	dictionary.CloseDatabaseConnection()
	removeDatabase()
}

func TestSQLiteDictionary_GetAllRegex(t *testing.T) {
	initDatabase()
	cfg.DatabaseHost = testSQLiteDatabasePath
	dictionary.Cfg = cfg

	err := dictionary.InitDatabaseConnection()
	assert.NoError(t, err)
	upTestTables(t)

	var result = map[int64]string{}
	result, err = dictionary.GetAllRegex()
	assert.NoError(t, err)
	assert.Empty(t, result)

	insertTestData(t)

	result, err = dictionary.GetAllRegex()
	assert.NoError(t, err)
	assert.NotEmpty(t, result)

	expected := map[int64]string{
		int64(1): `(?i)hello`,
		int64(2): `(?i)Say hello to (?P<name>\w+)`,
	}

	assert.Equal(t, expected, result)

	dropTestTables(t)
	dictionary.CloseDatabaseConnection()
	removeDatabase()
}

func TestSQLiteDictionary_FindRegex(t *testing.T) {
	initDatabase()
	cfg.DatabaseHost = testSQLiteDatabasePath
	dictionary.Cfg = cfg

	err := dictionary.InitDatabaseConnection()
	assert.NoError(t, err)
	upTestTables(t)
	insertTestData(t)

	var regexID int64
	regexID, err = dictionary.FindRegex("test")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), regexID)

	regexID, err = dictionary.FindRegex("(?i)hello")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), regexID)

	dropTestTables(t)
	dictionary.CloseDatabaseConnection()
	removeDatabase()
}

func TestSQLiteDictionary_FindScenarioByID(t *testing.T) {
	initDatabase()
	cfg.DatabaseHost = testSQLiteDatabasePath
	dictionary.Cfg = cfg

	err := dictionary.InitDatabaseConnection()
	assert.NoError(t, err)
	upTestTables(t)
	insertTestData(t)

	var scenarioID int64
	scenarioID, err = dictionary.FindScenarioByID(int64(11))
	assert.NoError(t, err)
	assert.Equal(t, int64(0), scenarioID)

	scenarioID, err = dictionary.FindScenarioByID(int64(1))
	assert.NoError(t, err)
	assert.Equal(t, int64(1), scenarioID)

	dropTestTables(t)
	dictionary.CloseDatabaseConnection()
	removeDatabase()
}

func TestSQLiteDictionary_InsertQuestionRegex(t *testing.T) {
	initDatabase()
	cfg.DatabaseHost = testSQLiteDatabasePath
	dictionary.Cfg = cfg

	err := dictionary.InitDatabaseConnection()
	assert.NoError(t, err)
	upTestTables(t)
	insertTestData(t)

	var insertedID int64
	insertedID, err = dictionary.InsertQuestionRegex("test", "")
	assert.NoError(t, err)
	assert.NotEmpty(t, insertedID)

	insertedID, err = dictionary.InsertQuestionRegex("test", "")
	assert.Error(t, err)
	assert.Equal(t, "UNIQUE constraint failed: questions_regex.regex", err.Error())

	insertedID, err = dictionary.InsertQuestionRegex("test2", "test_group")
	assert.NoError(t, err)
	assert.NotEmpty(t, insertedID)

	dropTestTables(t)
	dictionary.CloseDatabaseConnection()
	removeDatabase()
}

func upTestTables(t *testing.T) {
	availableTables["events"] = `
	drop table if exists events;
	create table events
	(
		id integer
			constraint events_pk
				primary key autoincrement,
		alias varchar not null,
		installed_version varchar
	);
	
	create unique index events_name_uindex
		on events (alias);
	`
	availableTables["questions"] = `
	drop table if exists questions;
	create table questions
	(
		id          integer
			constraint questions_pk
				primary key autoincrement,
		question    varchar not null,
		answer      varchar not null,
		scenario_id int     not null
			references scenarios
				on delete cascade,
		regex_id    integer
							references questions_regex
								on delete set null
	);
	
	create unique index questions_question_uindex
		on questions (question);
	`
	availableTables["scenarios"] = `
	drop table if exists scenarios;
	create table scenarios
	(
		id       integer
			constraint scenarios_pk
				primary key autoincrement,
		name     varchar not null,
		event_id integer not null
			references events
				on delete set null
	);
	
	create unique index scenarios_name_uindex
		on scenarios (name);
	`
	availableTables["questions_regex"] = `
	drop table if exists questions_regex;
	create table questions_regex
	(
		id          integer not null
			constraint question_regex_pk
				primary key autoincrement,
		regex       varchar not null,
		regex_group varchar
	);
	
	create unique index question_regex_regex_uindex
		on questions_regex (regex);
	`

	for _, query := range availableTables {
		_, err := dictionary.client.GetClient().Exec(query)
		if err != nil {
			assert.NoError(t, err)
		}
	}
}

func insertTestData(t *testing.T) {
	_, err := dictionary.client.GetClient().Exec(`
	INSERT INTO events (id, alias) VALUES (1, 'hello');
	INSERT INTO questions (id, question, answer, scenario_id, regex_id) VALUES (1, 'Hello', 'Hello', 1, 1);
	INSERT INTO questions (id, question, answer, scenario_id, regex_id) VALUES (2, 'Say hello to John', 'Hello %s', 1, 2);
	INSERT INTO questions (id, question, answer, scenario_id, regex_id) VALUES (3, 'Hello world', 'Hello', 1, 0);
	INSERT INTO questions_regex (id, regex, regex_group) VALUES (1, '(?i)hello', '');
	INSERT INTO questions_regex (id, regex, regex_group) VALUES (2, '(?i)Say hello to (?P<name>\w+)', 'name');
	INSERT INTO scenarios (id, name, event_id) VALUES (1, 'Scenario #1', 1);
	`)
	if err != nil {
		assert.NoError(t, err)
	}
}

func dropTestTables(t *testing.T) {
	for tableName := range availableTables {
		_, err := dictionary.client.GetClient().Exec("drop table if exists " + tableName)
		if err != nil {
			assert.NoError(t, err)
		}
	}
}

func initDatabase() {
	_, err := os.Create(testSQLiteDatabasePath)
	if err != nil {
		fmt.Println("Failed to create database file: " + err.Error())
	}
}

func removeDatabase() {
	err := os.Remove(testSQLiteDatabasePath)
	if err != nil {
		fmt.Println("Failed to remove database file: " + err.Error())
	}
}
