# Add new answer event
This event will insert into our database a new answer based on your request.

## Installation guide
To install it please run 
``` 
make build-installation-script && scripts/install/run --event_alias=dictionary
```

## Usage
Write in PM or tag the bot user with this message
```
New answer
Scenario id: SCENARIO_ID_PUT_HERE (optional)
Question: QUESTION_PUT_HERE (required)
Question regex: QUESTION_REGEX_PUT_HERE (optional)
Question regex group: QUESTION_REGEX_GROUP_PUT_HERE (optional)
Answer: ANSWER_PUT_HERE (required)
Event alias: EVENT_ALIAS_PUT_HERE (optional, by default it will be used as text message)
```
As successful result of event execution a new answer will in inserted into questions and answers tables 
