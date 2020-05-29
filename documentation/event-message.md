# Event Message
Here you can find information about the event-message object and how it can be used in you custom-event.

## The structure
Each event receives in the `Execute` method the object type of [dto.BaseChatMessage](../internal/dto/base_chat_message.go). While working with this message you might need to use the following attributes:
- `DictionaryMessage` - this is a generated answer, which bot found in the database. 
- `Text` - text of answer which can be modified by your custom event 
- `OriginalMessage` - this is a copy of original message, which should not be changed and should be used as source of truth.

More details you can find in the [dto.BaseChatMessage](../internal/dto/base_chat_message.go)
