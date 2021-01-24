# Migrations
There is a possibility to use migrations functionality in the project. This can help you to update your custom event features, play with the database of the devbot.

1. Using plain sql
For those who want to execute simple plain sql queries.

2. Using available migration logic
For those who want to use complex logic in their migrations. Such as variables or loops and etc.

### How to use
Inside of container there is service.MigrationService injected. So you can use available functionality from that service in your event.

`SetMigration` - method for preparing of your migration for execution. If you didn't set your migration using this method, then your migration will not be executed.
`RunMigrations` - method will run all the migrations which were prepared for execution.

Each migration should implement the `database.BaseMigrationInterface`.