package eventslist

import (
	"github.com/sharovik/devbot/internal/container"
)

//InsertHelpMessageMigration the migration name
type InsertHelpMessageMigration struct {
}

//GetName the name in the database
func (m InsertHelpMessageMigration) GetName() string {
	return "help-message-introducing.sql" //I leave this name to be backwards compatible with the old version of migrations flow
}

//Execute the body of the migration
func (m InsertHelpMessageMigration) Execute() error {
	eventID, err := container.C.Dictionary.FindEventByAlias(EventName)
	if err != nil {
		return err
	}

	scenarioID, err := container.C.Dictionary.InsertScenario("Scenario help", eventID)
	if err != nil {
		return err
	}

	_, err = container.C.Dictionary.InsertQuestion("Help", "If you want to see the list of my functions, please try ask me the following question `events list`. This will printout all possible phrases what currently I can understand.", scenarioID, "(?i)^(help)$", "")
	if err != nil {
		return err
	}

	return nil
}
