# Migrations
There is a possibility to use migrations functionality in the project. This can help you to update your custom event features, play with the database of the devbot.
Currently, there is only one way how to trigger the migrations for your event or for project itself - using the `container.C.MigrationService` service.

### How to use
Inside of container there is service.MigrationService injected. So you can use available functionality from that service in your event by simply calling `container.C.MigrationService.SetMigration(migration)`.

`SetMigration` - method for scheduling of your migration for execution. As attribute, it receives an object type of `database.BaseMigrationInterface`
`RunMigrations` - method will run all the migrations which were prepared for execution.

See the example in `Update` method of `events/eventslist/event.go` event.