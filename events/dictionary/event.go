package dictionary

import (
	"fmt"
	"html"
	"strconv"

	"github.com/sharovik/devbot/internal/log"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/helper"
)

const (
	//EventName the name of the event
	EventName = "dictionary"

	//EventVersion the version of the event
	EventVersion = "1.0.0"

	helpMessage = "Please use the following template:```" + `
New answer
Scenario id: SCENARIO_ID_PUT_HERE (optional)
Question: QUESTION_PUT_HERE (required)
Question regex: QUESTION_REGEX_PUT_HERE (optional)
Question regex group: QUESTION_REGEX_GROUP_PUT_HERE (optional)
Answer: ANSWER_PUT_HERE (required)
Event alias: EVENT_ALIAS_PUT_HERE (optional, by default it will be used as text message)
	` + "```"

	//Regex for catching of the information from the received message
	regexScenarioIDAttribute         = "(?im)((?:scenario id:) (?P<scenario_id>.+))"
	regexScenarioNameAttribute       = "(?im)((?:scenario name:) (?P<scenario_name>.+))"
	regexQuestionAttribute           = "(?im)((?:question:) (?P<question>.+))"
	regexQuestionRegexAttribute      = "(?im)((?:question regex:) (?P<question_regex>.+))"
	regexQuestionRegexGroupAttribute = "(?im)((?:question regex group:) (?P<question_regex_group>.+))"
	regexAnswerAttribute             = "(?im)((?:answer:) (?P<answer>.+))"
	regexEventAliasAttribute         = "(?im)((?:event alias:) (?P<event_alias>.+))"
	defaultEventAlias                = "text"
)

var (
	scenarioID         int64
	scenarioName       string
	question           string
	questionRegex      string
	questionRegexGroup string
	answer             string
	eventAlias         string
)

//DctnrEvent the struct for the event object
type DctnrEvent struct {
	EventName string
}

//Event - object which is ready to use
var Event = DctnrEvent{
	EventName: EventName,
}

//Execute method which is called by message processor
func (e DctnrEvent) Execute(message dto.BaseChatMessage) (dto.BaseChatMessage, error) {
	var answerMessage = message

	isHelpAnswerTriggered, err := helper.HelpMessageShouldBeTriggered(answerMessage.OriginalMessage.Text)
	if err != nil {
		log.Logger().Warn().Err(err).Msg("Something went wrong with help message parsing")
	}

	if isHelpAnswerTriggered {
		answerMessage.Text = helpMessage
		return answerMessage, nil
	}

	if err := parseAttributes(html.UnescapeString(message.OriginalMessage.Text)); err != nil {
		answerMessage.Text = "Error received during the attributes parsing: " + err.Error()
		return answerMessage, err
	}

	//We get the event id for selected event alias. The eventID we will use for scenarioID and question inserting
	eventID, err := container.C.Dictionary.FindEventByAlias(eventAlias)
	if err != nil {
		panic(err)
	}

	//If we received empty event id, it means that for that event-alias we don't have any row created. We need to create it now
	if eventID == 0 {
		eventID, err = container.C.Dictionary.InsertEvent(eventAlias, EventVersion)
		if err != nil {
			panic(err)
		}
	}

	//Now we need to do the similar procedure for the scenarioID
	scenarioID, err = container.C.Dictionary.FindScenarioByID(scenarioID)
	if err != nil {
		panic(err)
	}

	//If the scenarioID is 0 it means that scenarioID is not created. We need to create it now
	if scenarioID == 0 {
		lastScenarioID, err := container.C.Dictionary.GetLastScenarioID()
		if err != nil {
			panic(err)
		}

		if scenarioName == "" {
			scenarioName = fmt.Sprintf("Scenario #%d", lastScenarioID+1)
		}

		scenarioID, err = container.C.Dictionary.InsertScenario(scenarioName, eventID)
		if err != nil {
			panic(err)
		}
	}

	//In that step we have valid scenarioID and eventID. It means that we can proceed with question creation
	var questionID int64
	questionID, err = container.C.Dictionary.InsertQuestion(question, answer, scenarioID, questionRegex, questionRegexGroup)
	if err != nil {
		panic(err)
	}

	answerMessage.Text = fmt.Sprintf("I added this information to the dictionary.\nQuestionID: %d\nQuestion: %s\nAnswer: %s\nScenarioID: %d\nRegex: %s\nRegex group: %s", questionID, question, answer, scenarioID, questionRegex, questionRegexGroup)
	return answerMessage, nil
}

//Install method for installation of event
func (e DctnrEvent) Install() error {
	log.Logger().Debug().
		Str("event_name", EventName).
		Str("event_version", EventVersion).
		Msg("Triggered event installation")

	return container.C.Dictionary.InstallEvent(
		EventName,           //We specify the event name which will be used for scenario generation
		EventVersion,        //This will be set during the event creation
		"New answer", //Actual question, which system will wait and which will trigger our event
		"Ok, will do it now.",
		"(?i)(New answer)", //Optional field. This is regular expression which can be used for question parsing.
		"",                        //Optional field. This is a regex group and it can be used for parsing the match group from the regexp result
	)
}

//Update for event update actions
func (e DctnrEvent) Update() error {
	return nil
}

func parseAttributes(text string) error {
	var (
		err                 error
		_scenarioID         string
		_scenarioName       string
		_question           string
		_questionRegex      string
		_questionRegexGroup string
		_answer             string
		_eventAlias         string
	)

	_scenarioID = parseAttribute(text, regexScenarioIDAttribute, "scenario_id")
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
		return fmt.Errorf("Answer cannot be empty. ")
	}

	if _eventAlias == "" {
		eventAlias = defaultEventAlias
	}

	if _scenarioID == "" {
		scenarioID = int64(0)
	} else {
		scenarioID, err = strconv.ParseInt(_scenarioID, 10, 64)
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
