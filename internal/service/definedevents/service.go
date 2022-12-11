package definedevents

import (
	"github.com/sharovik/devbot/events"
	"github.com/sharovik/devbot/internal/container"
)

//InitializeDefinedEvents method initialize the defined events from the events.DefinedEvents configuration
func InitializeDefinedEvents() {
	container.C.DefinedEvents = events.DefinedEvents
}
