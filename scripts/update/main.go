package main

import (
	"flag"
	"fmt"
	"github.com/sharovik/devbot/events"
	"github.com/sharovik/devbot/internal/container"
	"os"
	"path"
	"runtime"
)

const descriptionEventAlias = "The event alias for which Update method will be called"

func init() {
	//We switch pointer to the root directory for control the path from which we need to generate test-data file-paths
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	_ = os.Chdir(dir)
}

func main()  {
	eventAlias := parseEventAlias()
	if eventAlias == "" {
		fmt.Println("Event alias cannot be empty")
		return
	}

	container.C = container.C.Init()
	if events.DefinedEvents.Events[eventAlias] == nil {
		fmt.Println("Event is not defined in the defined-events")
		return
	}

	if err := events.DefinedEvents.Events[eventAlias].Update(); err != nil {
		fmt.Println("Failed to update the event. Error:" + err.Error())
	}

	fmt.Println("Done")
}

func parseEventAlias() string {
	eventAlias := flag.String("event_alias", "", descriptionEventAlias)
	flag.Parse()

	fmt.Println("eventAlias:" + fmt.Sprintf("%s", *eventAlias))
	return *eventAlias
}