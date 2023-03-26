package event

import "github.com/sharovik/devbot/internal/dto"

//DefinedEventInterface the interface for events
type DefinedEventInterface interface {
	//Help returns the help message string
	Help() string

	Alias() string

	//Execute The main execution method, which will run the actual functionality for the event
	Execute(message dto.BaseChatMessage) (dto.BaseChatMessage, error)

	//Install The installation method, which will executes the installation parts of the event
	Install() error

	//Update The update method, which will update the application to use new version of this event
	Update() error
}
