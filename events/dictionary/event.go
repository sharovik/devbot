package dictionary

import (
	"fmt"
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/helper"
	"html"
	"strconv"
)

const (
	//EventName the name of the event
	EventName = "dictionary"

	//Regex for catching of the information from the received message
	regexScenarioIdAttribute         = "(?im)((?:scenario id:) (?P<scenario_id>.+))"
	regexScenarioNameAttribute       = "(?im)((?:scenario name:) (?P<scenario_name>.+))"
	regexQuestionAttribute           = "(?im)((?:question:) (?P<question>.+))"
	regexQuestionRegexAttribute      = "(?im)((?:question regex:) (?P<question_regex>.+))"
	regexQuestionRegexGroupAttribute = "(?im)((?:question regex group:) (?P<question_regex_group>.+))"
	regexAnswerAttribute             = "(?im)((?:answer:) (?P<answer>.+))"
	regexEventAliasAttribute         = "(?im)((?:event alias:) (?P<event_alias>.+))"
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

//ThemerEvent the struct for the event object
type ThemerEvent struct {
	EventName string
}

//Event - object which is ready to use
var Event = ThemerEvent{
	EventName: EventName,
}

//Execute method which is called by message processor
func (e ThemerEvent) Execute(message dto.SlackRequestChatPostMessage) (dto.SlackRequestChatPostMessage, error) {
	var answerMessage = message

	if err := parseAttributes(html.UnescapeString(message.OriginalMessage.Text)); err != nil {
		answerMessage.Text = "Error received during the attributes parsing: " + err.Error()
		return answerMessage, err
	}

	//We get the event id for selected event alias. The eventId we will use for scenarioId and question inserting
	eventId, err := container.C.Dictionary.FindEventByAlias(eventAlias)
	if err != nil {
		panic(err)
	}

	//If we received empty event id, it means that for that event-alias we don't have any row created. We need to create it now
	if eventId == 0 {
		eventId, err = container.C.Dictionary.InsertEvent(eventAlias)
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

	answerMessage.Text = fmt.Sprintf("I added this information to the dictionary.\nQuestionID: %d\nQuestion: %s\nAnswer: %s\nScenarioID: %d\nRegex: %s\nRegex group: %s", questionId, question, answer, scenarioId, questionRegex, questionRegexGroup)
	return answerMessage, nil
}

func parseAttributes(text string) error {
	var (
		err                 error
		_scenarioId         string
		_scenarioName       string
		_question           string
		_questionRegex      string
		_questionRegexGroup string
		_answer             string
		_eventAlias         string
	)

	_scenarioId = parseAttribute(text, regexScenarioIdAttribute, "scenario_id")
	_scenarioName = parseAttribute(text, regexScenarioNameAttribute, "scenario_name")
	_question = parseAttribute(text, regexQuestionAttribute, "question")
	_questionRegex = parseAttribute(text, regexQuestionRegexAttribute, "question_regex")
	_questionRegexGroup = parseAttribute(text, regexQuestionRegexGroupAttribute, "question_regex_group")
	_answer = parseAttribute(text, regexAnswerAttribute, "answer")
	_eventAlias = parseAttribute(text, regexEventAliasAttribute, "event_alias")

	if _question == "" {
		return fmt.Errorf("Question cannot be empty. ")
	}

	if _answer == "" {
		return fmt.Errorf("Question cannot be empty. ")
	}

	if _eventAlias == "" {
		return fmt.Errorf("Question cannot be empty. ")
	}

	if _scenarioId == "" {
		scenarioId = int64(0)
	} else {
		scenarioId, err = strconv.ParseInt(_scenarioId, 10, 64)
		if err != nil {
			return err
		}
	}

	scenarioName = _scenarioName
	question = _question
	questionRegex = _questionRegex
	questionRegexGroup = _questionRegexGroup
	answer = _answer
	eventAlias = _eventAlias

	return nil
}

func parseAttribute(text string, regex string, group string) string {
	matches := helper.FindMatches(regex, text)

	if len(matches) != 0 && group != "" && matches[group] != "" {
		return matches[group]
	}

	return ""
}
