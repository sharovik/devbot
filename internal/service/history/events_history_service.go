package history

import (
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/dto/databasedto"
	"github.com/sharovik/devbot/internal/log"
	"github.com/sharovik/devbot/internal/service/base"
	"github.com/sharovik/orm/clients"
	cdto "github.com/sharovik/orm/dto"
	"strings"
	"time"
)

const VariablesSeparator = ";"

//RememberEventExecution method for store the history of the event execution
func RememberEventExecution(msg dto.BaseChatMessage) {
	command := msg.Text
	if base.GetConversation(msg.Channel).Question != "" {
		command = base.GetConversation(msg.Channel).Question
	}

	item := databasedto.EventTriggerHistoryStruct{
		TableName: databasedto.EventTriggerHistoryModel.GetTableName(),
		Fields: []interface{}{
			cdto.ModelField{
				Name:  "event_id",
				Value: msg.DictionaryMessage.EventID,
			},
			cdto.ModelField{
				Name:  "scenario_id",
				Value: msg.DictionaryMessage.ScenarioID,
			},
			cdto.ModelField{
				Name:  "user",
				Value: msg.OriginalMessage.User,
			},
			cdto.ModelField{
				Name:  "channel",
				Value: msg.Channel,
			},
			cdto.ModelField{
				Name:  "command",
				Value: command,
			},
			cdto.ModelField{
				Name:  "variables",
				Value: strings.Join(base.GetConversation(msg.Channel).Variables, VariablesSeparator),
			},
			cdto.ModelField{
				Name:  "last_question_id",
				Value: base.GetConversation(msg.Channel).LastQuestion.DictionaryMessage.QuestionID,
			},
			cdto.ModelField{
				Name:  "created",
				Value: time.Now().Unix(),
			},
		},
	}

	c := container.C.Dictionary.GetNewClient()

	if _, err := c.Execute(new(clients.Query).Insert(&item)); err != nil {
		log.Logger().AddError(err).Msg("Failed to insert a log entry into events history table")
	}
}