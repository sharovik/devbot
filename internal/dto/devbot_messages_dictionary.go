package dto

//DictionaryMessage child struct of DevBotMessageDictionary object
type DictionaryMessage struct {
	ScenarioID            int64
	Question              string
	QuestionID            int64
	Regex                 string
	Answer                string
	MainGroupIndexInRegex string
	ReactionType          string
}
