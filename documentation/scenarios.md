# Scenarios
This feature can help with attributes guessing for your custom event. 
Imagine user triggered your event, but unfortunately from the message bot cannot parse all the attributes. In this case, instead of throwing of error message, you can trigger the scenario for your event, which will ask a questions and set needed variables for your custom event.

## Demo
### With tagging of the bot-user
![demo-tagging](images/scenario-demo-tagging.gif)

### Without tagging of the bot-user
![non-demo-tagging](images/scenario-demo-without-tagging.gif)

### How to stop active scenario
To stop the scenario, please use the following phrases:
- `stop!`
- `stop scenario!`
Once bot receives some of these phrases, he will try to stop the active scenario in the current channel, where the message posted.

## Database
Before describing of the code base, let's check the database schema and see how on the database level the scenario looks like.
First, let's have a look on the database schema

![database-schema](images/database-structure.png)

As you can see we have the scenario_id in the questions table. This field will be used for grouping of the questions.

On the screenshot below you can see that not each question row has the filled `question attribute`. This is because initially `question` attribute of `questions` table is used as the **question from the user** and `answer` field is used as **answer from the bot**.
And in case of the scenario, the bot **asks user** and **we don't expect the question** from the user.
 
![scenario-list-questions](images/scenario-list-questions.png)

The query which is used here

```sql
select q.id, q.question, q.answer, e.alias
from questions q
join scenarios s on q.scenario_id = s.id
join events e on s.event_id = e.id
where s.id = "{SCENARIO_ID}"
order by q.id asc
```

Each scenario must have:
1. at least 2 questions
2. a connected event with the alias defined(otherwise the custom event will not be triggered)
3. only first question of scenario should have the filled `question` attribute and all next questions should have that field as empty string

## Conversations
Each trigger of scenario opens the conversation for the channel from where the message was received. That means, once the bod started the scenario, you will not be able to ask him other questions, because he is expecting the answers for the open scenario conversation.

Below you can see how the example of how the scenario processing looks like

![scenario-message-processing](images/scenario-message-processing.png)

## Installation of scenario
Each scenario with multiple questions and answers should have not less than 2 questions, otherwise it **will not be handled** as scenario, but as a simple question.
So for installation you will need to create the initial question and answer first and then, based on created scenario ID, connect to it one or more questions.

Here you can see the example of Install method for your custom event, where we install the scenario

```go
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
```

As you can see, in the `database.NewEventScenario` struct you can define multiple questions, which will be connected to a one scenario.
