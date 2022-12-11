package history

import (
	"strings"
	"time"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/dto/databasedto"
	"github.com/sharovik/devbot/internal/log"
	"github.com/sharovik/devbot/internal/service/base"
	"github.com/sharovik/orm/clients"
	cdto "github.com/sharovik/orm/dto"
)

const variablesSeparator = ";"

//RememberEventExecution method for store the history of the event execution
func RememberEventExecution(msg dto.BaseChatMessage) {
	command := msg.Text
	if base.GetConversation(msg.Channel).Question != "" {
		command = base.GetConversation(msg.Channel).Question
	}

	var variables []string
	for _, variable := range base.GetConversation(msg.Channel).Scenario.RequiredVariables {
		variables = append(variables, variable.Value)
	}

	item := databasedto.EventTriggerHistoryModel
	item.AddModelField(cdto.ModelField{
		Name:  "event_id",
		Value: msg.DictionaryMessage.EventID,
	})
	item.AddModelField(cdto.ModelField{
		Name:  "scenario_id",
		Value: msg.DictionaryMessage.ScenarioID,
	})
	item.AddModelField(cdto.ModelField{
		Name:  "user",
		Value: msg.OriginalMessage.User,
	})
	item.AddModelField(cdto.ModelField{
		Name:  "channel",
		Value: msg.Channel,
	})
	item.AddModelField(cdto.ModelField{
		Name:  "command",
		Value: command,
	})
	item.AddModelField(cdto.ModelField{
		Name:  "variables",
		Value: strings.Join(variables, variablesSeparator),
	})
	item.AddModelField(cdto.ModelField{
		Name:  "last_question_id",
		Value: base.GetConversation(msg.Channel).LastQuestion.DictionaryMessage.QuestionID,
	})
	item.AddModelField(cdto.ModelField{
		Name:  "created",
		Value: time.Now().Unix(),
	})

	c := container.C.Dictionary.GetDBClient()

	if _, err := c.Execute(new(clients.Query).Insert(item)); err != nil {
		log.Logger().AddError(err).Msg("Failed to insert a log entry into events history table")
	}
}
