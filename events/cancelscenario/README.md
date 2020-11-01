# Stop the scenario event
This event will stop graceful execution of scenario for specified channel.

## Installation guide
To install it please run 
``` 
make build-installation-script && scripts/install/run --event_alias=cancelscenario
```

## Usage
Write in PM or tag the bot user with this message
```
stop conversation #channel-name|@username|<CHANNELID>
```
As successful result of event execution the scenario execution in the selected channel will be interrupted. 
