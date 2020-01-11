##Dictionary
This is the functionality by which the DevBot understand what kind of the event he need to trigger for your message. Currently it is a simplified version, where we go through simple json file and search by regexp the matches and if match was found, bot run the logic for that message.

## Table of contents
- [The dictionary database](#the-dictionary-database)
- [The dictionary object structure](#the-dictionary-object-structure)
- [Dictionary object generator tool](#dictionary-object-generator-tool)
-- [Available options](#available-options)
-- [Example of generator tool execution](#example-of-generator-tool-execution)

## The dictionary database
Currently I used the simple way of storing data for dictionary messages - simple json file. You can find it here `internal/dictionary/slack_dictionary.json`.
There is 2 groups of messages `text_message_dictionary` and `file_message_dictionary`. Each of these groups contains the objects which have questions and answers.
Below you can see the example of the structure of the object which have question and answer:
```json
{"question":"(?i)Hey","answer":"Yo","main_group_index_in_regex":1,"reaction_type":"help"}
```

## The dictionary object structure
Below you can see the description for each field of the question and answer object
* **question** - regexp for question. Example: `Foo`
* **answer** - regexp for answer. Example: `Bar`
* **main_group_index_in_regex** - the number of the group in regexp which should be taken. Example: `(?:Foo (Bar))`, from that regexp we need to take the group which contains the `Bar`. In that regexp this group will be first, so we put in `main_group_index_in_regex` the `1` value 
* **reaction_type** - the name of the event which should be triggered after the message answer prepare. Example: see [events documentation](events.md)

## Dictionary object generator tool
There is a tool which you can use for creation of new questions and answers for our dictionary. You can find this tool here `scripts/dictionary-loader/dictionary-loader`

### Available options:
* selectedDictionary - existing dictionary which you can find in internal/selectedDictionary folder
* type - Type of selectedDictionary. By default will be `text_message_dictionary`
* question - a question. It can be static or can be regex
* answer - the answer
* groupIndex - Group index in regex. This will get by selected index, group in your string regex and try to use it for answers and actions
* reactionAction - Type of reaction, which should be used for this answer. If it's empty, then only text message reaction will be executed

### Example of generator tool execution
``` 
./scripts/dictionary-loader/dictionary-loader --question=Foo --answer=Bar --reactionAction=help
```
This command will insert new question for `text_message_dictionary` in the file `internal/dictionary/slack_dictionary.json`