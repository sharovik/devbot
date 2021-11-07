package definedevents

import (
	"github.com/sharovik/devbot/events"
	"github.com/sharovik/devbot/internal/container"
)

func InitializeDefinedEvents() {
	container.C.DefinedEvents = events.DefinedEvents
}