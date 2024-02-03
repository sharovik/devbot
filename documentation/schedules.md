# Schedules
In order to schedule execution of your events, you might trigger `schedule.S.Schedule` method of `schedule` package. You can do that either via migrations OR during the event installation/update

## How to
First, you need to prepare scenario, with your event name to make sure we execute the right one at the end. PrepareEventScenario method generates the scenario object based on the eventId and reaction type
```go
scenario, err := service.PrepareEventScenario(eventID, EventName)
if err != nil {
    //@todo: handle error
}
```
Secondary, you need to prepare a scheduled time. You can use the following format:
1. `in 1 hour and 2 minutes`
2. `1 hour`
3. `schedule event examplescenario every 1 minute`
4. `23 minutes` OR `in 20 minutes`
5. `2022-12-18 11:22`
6. `in 1 day`
7. `repeat 1 days and at 9:30`
8. `Sunday at 10:00`
9. `every monday at 9:10`
```go
scheduleTime, err := new(schedule.ExecuteAt).FromString(text)
if err != nil {
    //@todo: handle error
}
```
Optionally, you can define the `variables` attribute, where you specify variables for your scheduled event. The values should have the same order as it is defined in the target event. 
As delimiter, you need to use `schedule.VariablesDelimiter`. Example:
```
testValue1;testValue2
```
And finally, you generate the `schedule.Item` object, with the selected event and schedule time. After that, you call `schedule.S.Schedule(item)` method to schedule your event.
```go
var item = schedule.Item{
    Author:       "AUTHOR_ID",
    Channel:      "CHANNEL_ID",
    ScenarioID:   scenario.ID,
    EventID:      scenario.EventID,
    ReactionType: scenario.EventName,
    Variables:    strings.Join(variables, schedule.VariablesDelimiter),
    Scenario:     scenario,
    ExecuteAt:    scheduleTime,
    IsRepeatable: false,
}
```
## Example of usage inside the event
Below you can see the example of implementation inside the custom event `Execute` method, where we are receiving `message` object, which contain the received chat-message
```go
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
