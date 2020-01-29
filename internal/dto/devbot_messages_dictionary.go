package dto

//DictionaryMessage child struct of DevBotMessageDictionary object
type DictionaryMessage struct {
	ScenarioID            int64
	Question              string
	Regex                 string
	Answer                string
	MainGroupIndexInRegex string
	ReactionType          string
}
