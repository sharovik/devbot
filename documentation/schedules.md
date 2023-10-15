# Schedules
In order to schedule your events you might need to trigger `schedule.S.Schedule` method.

This method receives a `schedule.Item` model as an argument. Below you can find an example:
```go
//We prepare the scenario, with our event name, to make sure we execute the right at the end.
//By doing this, we are rewriting the scenario event alias to make sure, after all questions were asked, we point to our schedule event back
scenario, err := service.PrepareEventScenario(eventID, EventName)
if err != nil {
    //@todo: handle error
}

scheduleTime, err := new(schedule.ExecuteAt).FromString(text)
if err != nil {
    //@todo: handle error
}

var item = schedule.Item{
    Author:       message.OriginalMessage.User,
    Channel:      message.Channel,
    ScenarioID:   scenario.ID,
    EventID:      scenario.EventID,
    ReactionType: scenario.EventName,
    Variables:    strings.Join(variables, schedule.VariablesDelimiter),
    Scenario:     scenario,
    ExecuteAt:    scheduleTime,
    IsRepeatable: false,
}
schedule.S.Schedule(item)
```

