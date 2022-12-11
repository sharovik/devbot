package eventslist

import (
	"fmt"

	"github.com/sharovik/devbot/internal/database"
	"github.com/sharovik/devbot/internal/dto/databasedto"
	"github.com/sharovik/orm/clients"
	"github.com/sharovik/orm/query"

	"github.com/sharovik/devbot/internal/log"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
)

const (
	//EventName the name of the event
	EventName = "eventslist"

	//EventVersion the version of the event
	EventVersion = "1.0.1"
)

//EListEvent the struct for the event object. It will be used for initialisation of the event in defined-events.go file.
type EListEvent struct {
	EventName string
}

//Event - object which is ready to use
var (
	Event = EListEvent{
		EventName: EventName,
	}
	m = []database.BaseMigrationInterface{
		InsertHelpMessageMigration{},
	}
)

//Execute method which is called by message processor
func (e EListEvent) Execute(message dto.BaseChatMessage) (dto.BaseChatMessage, error) {
	//This answer will be show once the event get triggered.
	//Leave message.Text empty, once you need to not show the message, once this event get triggered.
	c := container.C.Dictionary.GetDBClient()

	q := new(clients.Query).
		Select([]interface{}{"events.id", "events.alias", "questions.question"}).
		From(databasedto.EventModel).
		Join(query.Join{
			Target: query.Reference{
				Table: databasedto.ScenariosModel.GetTableName(),
				Key:   "event_id",
			},
			With: query.Reference{
				Table: databasedto.EventModel.GetTableName(),
				Key:   databasedto.EventModel.GetPrimaryKey().Name,
			},
			Condition: "=",
		}).
		Join(query.Join{
			Target: query.Reference{
				Table: databasedto.QuestionsModel.GetTableName(),
				Key:   "scenario_id",
			},
			With: query.Reference{
				Table: databasedto.ScenariosModel.GetTableName(),
				Key:   databasedto.ScenariosModel.GetPrimaryKey().Name,
			},
			Condition: "=",
		}).
		Where(query.Where{
			First:    "questions.question",
			Operator: "<>",
			Second:   "''",
		})
	res, err := c.Execute(q)
	if err != nil {
		message.Text = fmt.Sprintf("Hmm. I tried to get the list of the available events and I failed. Here is the error: ```%s```", err)
		return message, err
	}

	message.Text = "Here is the list of the possible phrases which your can use:"
	for _, item := range res.Items() {
		message.Text += fmt.Sprintf("\n#%d event: `%s`. Try to ask `%s`. Also you could try to ask `%s --help`",
			item.GetField("id").Value,
			item.GetField("alias").Value,
			item.GetField("question").Value,
			item.GetField("question").Value,
		)
	}

	return message, nil
}

//Install method for installation of event
func (e EListEvent) Install() error {
	log.Logger().Debug().
		Str("event_name", EventName).
		Str("event_version", EventVersion).
		Msg("Triggered event installation")

	return container.C.Dictionary.InstallNewEventScenario(database.EventScenario{
		EventName:    EventName,
		EventVersion: EventVersion,
		Questions: []database.Question{
			{
				Question:      "events list",
				Answer:        "Just a sec, I will prepare the list for you.",
				QuestionRegex: "(?i)events list",
				QuestionGroup: "",
			},
		},
	})
}

//Update for event update actions
func (e EListEvent) Update() error {
	for _, migration := range m {
		container.C.MigrationService.SetMigration(migration)
	}

	return container.C.MigrationService.RunMigrations()
}
