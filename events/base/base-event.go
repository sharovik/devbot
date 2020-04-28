package base

import "github.com/sharovik/devbot/internal/dto"

//Event main interface for events
type Event interface {
	//The main execution method, which will run the actual functionality for the event
	Execute(message dto.BaseChatMessage) (dto.BaseChatMessage, error)

	//The installation method, which will executes the installation parts of the event
	Install() error

	//The update method, which will update the application to use new version of this event
	Update() error
}

//Events the struct object which is used for events storing
type Events struct {
	Events map[string]Event
}
