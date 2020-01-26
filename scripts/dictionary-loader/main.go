package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/sharovik/devbot/internal/dto"
	"github.com/sharovik/devbot/internal/helper"
)

func main() {

	var (
		currentDictionary dto.DevBotMessageDictionary

		selectedDictionary    = flag.String("selectedDictionary", "slack_dictionary", "existing selectedDictionary which you can find in internal/selectedDictionary folder")
		question              = flag.String("question", "", "a question. It can be static or can be regex")
		answer                = flag.String("answer", "", "the answer")
		mainGroupIndexInRegex = flag.Int("groupIndex", 0, "Group index in regex. This will get by selected index, group in your string regex and try to use it for answers and actions")
		reactionAction        = flag.String("reactionAction", "", "Type of reaction, which should be used for this answer. If it's empty, then only text message reaction will be executed")
	)

	flag.Parse()

	var pathToDictionary = fmt.Sprintf("./internal/dictionary/%s.json", *selectedDictionary)

	fmt.Println("path:" + fmt.Sprintf("%s", pathToDictionary))
	fmt.Println("question:" + fmt.Sprintf("%s", *question))
	fmt.Println("answer:" + fmt.Sprintf("%s", *answer))
	fmt.Println("groupIndex:" + fmt.Sprintf("%d", *mainGroupIndexInRegex))
	fmt.Println("reactionAction:" + fmt.Sprintf("%s", *reactionAction))

	bytes, err := helper.FileToBytes(pathToDictionary)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(bytes, &currentDictionary); err != nil {
		panic(err)
	}

	if *answer == "" {
		panic("Answer cannot be empty")
	}

	msg := dto.DictionaryMessage{
		Question:              *question,
		Answer:                *answer,
		MainGroupIndexInRegex: *mainGroupIndexInRegex,
		ReactionType:          *reactionAction,
	}
	currentDictionary = addToDictionary(msg, currentDictionary)
	if err := helper.ObjectToFile(pathToDictionary, currentDictionary); err != nil {
		panic(err)
	}
}

func addToDictionary(newMessage dto.DictionaryMessage, dictionary dto.DevBotMessageDictionary) dto.DevBotMessageDictionary {
	dictionary.Messages = append(dictionary.Messages, newMessage)
	return dictionary
}
