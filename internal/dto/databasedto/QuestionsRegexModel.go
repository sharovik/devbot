package databasedto

import "github.com/sharovik/orm/dto"

// QuestionsRegexStruct struct for QuestionsRegexModel
type QuestionsRegexStruct struct {
	dto.BaseModel
}

// QuestionsRegexModel model for questions_regex table
var QuestionsRegexModel = New(
	"questions_regex",
	[]interface{}{
		dto.ModelField{
			Name:   "regex",
			Type:   dto.VarcharColumnType,
			Length: 255,
		},
		dto.ModelField{
			Name:   "regex_group",
			Type:   dto.VarcharColumnType,
			Length: 255,
		},
	},
	dto.ModelField{
		Name:          "id",
		Type:          dto.IntegerColumnType,
		AutoIncrement: true,
		IsPrimaryKey:  true,
	},
	&QuestionsRegexStruct{},
)
