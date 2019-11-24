package dto

//DictionaryMessage child struct of DevBotMessageDictionary object
type DictionaryMessage struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	MainGroupIndexInRegex int `json:"main_group_index_in_regex"`
	ReactionType string `json:"reaction_type"`
}

//DevBotMessageDictionary main dictionary of DevBot
type DevBotMessageDictionary struct {
	TextMessageDictionary []DictionaryMessage `json:"text_message_dictionary"`
	FileMessageDictionary []DictionaryMessage `json:"file_message_dictionary"`
}
