package databasedto

import "github.com/sharovik/orm/dto"

//EventTriggerHistoryStruct the struct for event history table
type EventTriggerHistoryStruct struct {
	dto.BaseModel
}

//EventTriggerHistoryModel the actual object of event history table
var EventTriggerHistoryModel = New(
	"events_triggers_history",
	[]interface{}{
		dto.ModelField{
			Name: "event_id",
			Type: dto.IntegerColumnType,
		},
		dto.ModelField{
			Name: "scenario_id",
			Type: dto.IntegerColumnType,
		},
		dto.ModelField{
			Name:   "user",
			Type:   dto.VarcharColumnType,
			Length: 255,
		},
		dto.ModelField{
			Name:   "channel",
			Type:   dto.VarcharColumnType,
			Length: 255,
		},
		dto.ModelField{
			Name:   "command",
			Type:   dto.VarcharColumnType,
			Length: 255,
		},
		dto.ModelField{
			Name:       "variables",
			Type:       dto.VarcharColumnType,
			Length:     255,
			IsNullable: true,
		},
		dto.ModelField{
			Name:       "last_question_id",
			Type:       dto.IntegerColumnType,
			Length:     11,
			IsNullable: true,
		},
		dto.ModelField{
			Name:   "created",
			Type:   dto.IntegerColumnType,
			Length: 11,
		},
	},
	dto.ModelField{
		Name:          "id",
		Type:          dto.IntegerColumnType,
		AutoIncrement: true,
		IsPrimaryKey:  true,
	},
	&EventTriggerHistoryStruct{},
)
