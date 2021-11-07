# Events
This feature will help you to improve the skills of your bot. With it you are able to create your own event for your custom message.

## Table of contents
- [Good to know](#good-to-know-for-event-setup)
- [Prerequisites](#prerequisites)
- [The event diagram](#the-event-diagram)
- [Event setup](#event-setup)
- [Example](#example)

## Good to know for event setup
- `defined-events.go.dist` file example which contains the way of how to load your custom events
- [the event message object](event-message.md) - there you can find more details about event message object

## Prerequisites
* run `cp defined-events.go.dist defined-events.go` to create the file where you will define your events

## The event diagram
![base-event-diagram](images/base-event-scheme.png)

## Event setup
* create your event directory in `events` directory. Ex: `events/mybrandnewevent`
* create in your new directory file with name `event.go`. There is no black magic inside the naming, we just introduce the structured way of how to define the event files.
* create the logic for your new event struct object and make sure this logic is compatible with the interface `container.DefinedEvent`
* add your object to the "map" of the events `events.DefinedEvents` in the `defined-events.go` file 
```go
var DefinedEvents = map[string]container.DefinedEvent{
    //...
	mybrandnewevent.EventName: mybrandnewevent.Event,
}
```

## Example
Now let's take a deep dive into custom event creation and create a custom event.

### Prepare event folder
You need to specify the folder name without spaces. There should not be dashes or underscores in the package name, so please make sure you use a proper naming.
Let's call it `example` and create this example folder in the `events` folder.

### The code
There is a [Base DefinedEvent interface](../internal/container/container.go) which defines the structure of each event. So, let's create our custom event using this base interface definition:
```go
package example

import (
	"fmt"

    "github.com/sharovik/devbot/internal/helper"
    "github.com/sharovik/devbot/internal/log"

	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto"
)

const (
	//EventName the name of the event
	EventName = "example"

	//EventVersion the version of the event
	EventVersion = "1.0.1"

    helpMessage = "Ask me `who are you?` and you will see the answer."
)

//ExmplEvent the struct for the event object. It will be used for initialisation of the event in defined-events.go file.
type ExmplEvent struct {
	EventName string
}

//Event - object which is ready to use
var Event = ExmplEvent{
	EventName: EventName,
}

//Execute method which is called by message processor
func (e ExmplEvent) Execute(message dto.BaseChatMessage) (dto.BaseChatMessage, error) {
    isHelpAnswerTriggered, err := helper.HelpMessageShouldBeTriggered(message.OriginalMessage.Text)
    if err != nil {
        log.Logger().Warn().Err(err).Msg("Something went wrong with help message parsing")
    }

    if isHelpAnswerTriggered {
        message.Text = helpMessage
        return message, nil
    }

    //This answer will be show once the event get triggered.
    //Leave message.Text empty, once you need to not show the message, once this event get triggered.
	message.Text = "This is an example of the answer."
	return message, nil
}

//Install method for installation of event
func (e ExmplEvent) Install() error {
	log.Logger().Debug().
		Str("event_name", EventName).
		Str("event_version", EventVersion).
		Msg("Triggered event installation")

	return container.C.Dictionary.InstallNewEventScenario(database.NewEventScenario{
        EventName:    EventName,
        EventVersion: EventVersion,
        Questions:    []database.Question{
            {
                Question:      "who are you?",
                Answer:        fmt.Sprintf("Hello, my name is %s", container.C.Config.SlackConfig.BotName),
                QuestionRegex: "(?i)who are you?",
                QuestionGroup: "",
            },
        },
    })
}

//Update for event update actions
func (e ExmplEvent) Update() error {
	for _, migration := range m {
		container.C.MigrationService.SetMigration(migration)
	}

	return container.C.MigrationService.RunMigrations()
}
```

### Result
#### With empty text message in Execute method
![empty text message](images/example-event-text-empty.png)

#### With the filled text message in Execute method
![with text message](images/example-event-with-text.png)

### Source code
You can find the source code of the event in [events/example](https://github.com/sharovik/devbot/tree/master/events/example) folder