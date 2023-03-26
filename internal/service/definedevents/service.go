package definedevents

import (
	"github.com/sharovik/devbot/events"
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto/event"
)

// InitializeDefinedEvents method initialize the defined events from the events.DefinedEvents configuration
func InitializeDefinedEvents() {
	container.C.DefinedEvents = map[string]event.DefinedEventInterface{}

	for _, item := range events.DefinedEvents {
		container.C.DefinedEvents[item.Alias()] = item
	}
}
