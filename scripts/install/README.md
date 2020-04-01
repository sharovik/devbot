# Installation script
This script should be used for new events installation for your custom DevBot. Otherwise your custom event will not be catched by the system properly

### Before run
Make sure that you added your event into `events/defined-events.go` file. 

### How to use
run `make build-installation-script && scripts/install/run --event_alias={YOUR_ALIAS_NAME}`