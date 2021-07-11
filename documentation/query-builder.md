# The query builder
Here you can find the instructions how to use the query builder of the internal dictionary.

## Usage
Here you can find the basic example of the query builder usage:

```go
c := container.C.Dictionary.GetNewClient()

q := new(clients.Query).
    Select([]interface{}{"events.id", "events.alias", "questions.question"}).
    From(&database_dto.EventModel).
    Join(query.Join{
        Target: query.Reference{
            Table: database_dto.ScenariosModel.GetTableName(),
            Key:   "event_id",
        },
        With: query.Reference{
            Table: database_dto.EventModel.GetTableName(),
            Key:   database_dto.EventModel.GetPrimaryKey().Name,
        },
        Condition: "=",
    }).
    Join(query.Join{
        Target: query.Reference{
            Table: database_dto.QuestionsModel.GetTableName(),
            Key:   "scenario_id",
        },
        With: query.Reference{
            Table: database_dto.ScenariosModel.GetTableName(),
            Key:   database_dto.ScenariosModel.GetPrimaryKey().Name,
        },
        Condition: "=",
    }).
    Where(query.Where{
        First:    "questions.question",
        Operator: "<>",
        Second:   "''",
    })
res, err := c.Execute(q)
if err != nil {
    //do something with error
}
```
The example of usage in the event you can find here: [events/eventslist](../events/eventslist)`.

Another examples of the database query-builder and it's supported functionality [you can find here](https://github.com/sharovik/orm).
