package databasedto

import "github.com/sharovik/orm/dto"

//QuestionsStruct struct for QuestionsModel object
type QuestionsStruct struct {
	dto.BaseModel
}

//QuestionsModel model for questions table
var QuestionsModel = New(
	"questions",
	[]interface{}{
		dto.ModelField{
			Name:    "question",
			Type:    dto.VarcharColumnType,
			Length:  255,
			Default: "",
		},
		dto.ModelField{
			Name:    "answer",
			Type:    dto.VarcharColumnType,
			Length:  255,
			Default: "",
		},
		dto.ModelField{
			Name: "scenario_id",
			Type: dto.IntegerColumnType,
		},
		dto.ModelField{
			Name:       "regex_id",
			Type:       dto.IntegerColumnType,
			IsNullable: true,
		},
		dto.ModelField{
			Name:       "is_variable",
			Type:       dto.BooleanColumnType,
			IsNullable: true,
		},
	},
	dto.ModelField{
		Name:          "id",
		Type:          dto.IntegerColumnType,
		AutoIncrement: true,
		IsPrimaryKey:  true,
	},
	&QuestionsStruct{},
)
