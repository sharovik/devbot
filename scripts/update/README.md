# Update script
This script should be used for post-installation actions. For example, when event should install or insert custom data into the database. 

### Before run
Make sure that you added your event into `events/defined-events.go` file. 

### How to use
run `make build-update-script && scripts/update/run --event_alias={YOUR_ALIAS_NAME}`