package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sharovik/devbot/internal/service/definedevents"
	"os"
	"path"
	"runtime"

	"github.com/sharovik/devbot/internal/container"
)

var (
	scenarioId         int64
	scenarioName       string
	question           string
	questionRegex      string
	questionRegexGroup string
	answer             string
	eventAlias         string
)

const (
	defaultVersion = "1.0.0"

	//Description constants
	descriptionScenarioIdAttr         = "Scenario id, to which we need to attach this question. If 0 then new scenarioId will be created for this question"
	descriptionScenarioNameAttr       = "Scenario name"
	descriptionQuestionAttr           = "the question. It can be static or can be regex"
	descriptionQuestionRegexAttr      = "Will be used for identifying of the specific information from the question. Ex: from string 'Release master branch', by regex we can take the branch name."
	descriptionQuestionRegexGroupAttr = "Group which will be taken from selected regex"
	descriptionAnswerAttr             = "The answer for selected question"
	descriptionEventAliasAttr         = "The event alias. If alias doesn't exists in the database, then it will be created and used for this question"
)

func init() {
	//We switch pointer to the root directory for control the path from which we need to generate test-data file-paths
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	_ = os.Chdir(dir)

	container.C = container.C.Init()
	definedevents.InitializeDefinedEvents()
}

func main() {
	//We parse the args and printout them
	parseArgs()

	if err := validateArgs(); err != nil {
		panic(err)
	}

	//We get the event id for selected event alias. The eventId we will use for scenarioId and question inserting
	eventId, err := container.C.Dictionary.FindEventByAlias(eventAlias)
	if err != nil {
		panic(err)
	}

	//If we received empty event id, it means that for that event-alias we don't have any row created. We need to create it now
	if eventId == 0 {
		eventId, err = container.C.Dictionary.InsertEvent(eventAlias, defaultVersion)
		if err != nil {
			panic(err)
		}
	}

	//Now we need to do the similar procedure for the scenarioId
	scenarioId, err = container.C.Dictionary.FindScenarioByID(scenarioId)
	if err != nil {
		panic(err)
	}

	//If the scenarioId is 0 it means that scenarioId is not created. We need to create it now
	if scenarioId == 0 {
		lastScenarioId, err := container.C.Dictionary.GetLastScenarioID()
		if err != nil {
			panic(err)
		}

		if scenarioName == "" {
			scenarioName = fmt.Sprintf("Scenario #%d", lastScenarioId+1)
		}

		scenarioId, err = container.C.Dictionary.InsertScenario(scenarioName, eventId)
		if err != nil {
			panic(err)
		}
	}

	//In that step we have valid scenarioId and eventId. It means that we can proceed with question creation
	var questionId int64
	questionId, err = container.C.Dictionary.InsertQuestion(question, answer, scenarioId, questionRegex, questionRegexGroup)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("The question #%d was created!", questionId))
	fmt.Println(fmt.Sprintf("Now you can ask bot your question: %s", question))
}

func validateArgs() error {
	if answer == "" {
		return errors.New(fmt.Sprintf("Answer cannot be empty"))
	}

	if eventAlias == "" {
		return errors.New(fmt.Sprintf("Event alias cannot be empty"))
	}

	if question == "" {
		return errors.New(fmt.Sprintf("Question cannot be empty"))
	}

	return nil
}

func parseArgs() {
	_scenarioId := flag.Int64("scenario_id", 0, descriptionScenarioIdAttr)
	_scenarioName := flag.String("scenario_name", "", descriptionScenarioNameAttr)
	_question := flag.String("question", "Hello world", descriptionQuestionAttr)
	_questionRegex := flag.String("question_regex", "", descriptionQuestionRegexAttr)
	_questionRegexGroup := flag.String("question_regex_group", "", descriptionQuestionRegexGroupAttr)
	_answer := flag.String("answer", "Hey mate", descriptionAnswerAttr)
	_eventAlias := flag.String("event_alias", "", descriptionEventAliasAttr)

	flag.Parse()

	//Retrieve the value to the global vars
	scenarioId = *_scenarioId
	scenarioName = *_scenarioName
	question = *_question
	questionRegex = *_questionRegex
	questionRegexGroup = *_questionRegexGroup
	answer = *_answer
	eventAlias = *_eventAlias

	fmt.Println("scenarioId:" + fmt.Sprintf("%d", scenarioId))
	fmt.Println("scenario_name:" + fmt.Sprintf("%s", scenarioName))
	fmt.Println("question:" + fmt.Sprintf("%s", question))
	fmt.Println("questionRegex:" + fmt.Sprintf("%s", questionRegex))
	fmt.Println("questionRegexGroup:" + fmt.Sprintf("%s", questionRegexGroup))
	fmt.Println("answer:" + fmt.Sprintf("%s", answer))
	fmt.Println("eventAlias:" + fmt.Sprintf("%s", eventAlias))
}
