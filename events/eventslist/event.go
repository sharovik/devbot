package eventslist

import (
	"fmt"

	"github.com/sharovik/devbot/internal/log"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
)

const (
	//EventName the name of the event
	EventName = "eventslist"

	//EventVersion the version of the event
	EventVersion = "1.0.0"

	//The migrations folder, which can be used for event installation or for event update
	migrationDirectoryPath = "./events/eventslist/migrations"
)

//EListEvent the struct for the event object. It will be used for initialisation of the event in defined-events.go file.
type EListEvent struct {
	EventName string
}

//Event - object which is ready to use
var Event = EListEvent{
	EventName: EventName,
}

//Execute method which is called by message processor
func (e EListEvent) Execute(message dto.BaseChatMessage) (dto.BaseChatMessage, error) {
	//This answer will be show once the event get triggered.
	//Leave message.Text empty, once you need to not show the message, once this event get triggered.
	c := container.C.Dictionary.GetClient()

	rows, err := c.Query(`
	select e.id, e.alias, q.question
	from events e
	join scenarios s on e.id = s.event_id
	join questions q on s.id = q.scenario_id
	`)
	if err != nil {
		message.Text = fmt.Sprintf("Hmm. I tried to get the list of the available events and I failed. Here is the error: ```%s```", err)
		return message, err
	}

	var (
		id       int64
		alias    string
		question string
	)

	message.Text = "Here is the list of the events:"
	for rows.Next() {
		if err := rows.Scan(&id, &alias, &question); err != nil {
			message.Text += fmt.Sprintf("Oops!I tried to prepare the report for you and I failed. Here is the error: ```%s```", err)
			return message, err
		}

		message.Text += fmt.Sprintf("\n#%d event: `%s`. Try to ask `%s`", id, alias, question)
	}

	return message, nil
}

//Install method for installation of event
func (e EListEvent) Install() error {
	log.Logger().Debug().
		Str("event_name", EventName).
		Str("event_version", EventVersion).
		Msg("Triggered event installation")

	return container.C.Dictionary.InstallEvent(
		EventName,     //We specify the event name which will be used for scenario generation
		EventVersion,  //This will be set during the event creation
		"events list", //Actual question, which system will wait and which will trigger our event
		"Just a sec, I will prepare the list for you.", //Answer which will be used by the bot
		"(?i)events list", //Optional field. This is regular expression which can be used for question parsing.
		"", //Optional field. This is a regex group and it can be used for parsing the match group from the regexp result
	)
}

//Update for event update actions
func (e EListEvent) Update() error {
	return container.C.Dictionary.RunMigrations(migrationDirectoryPath)
}
