package databasedto

import "github.com/sharovik/orm/dto"

//SchedulesStruct the struct for schedules model
type SchedulesStruct struct {
	dto.BaseModel
}

//SchedulesModel the model for schedules table
var SchedulesModel = New(
	"schedules",
	[]interface{}{
		dto.ModelField{
			Name:       "author",
			Type:       dto.VarcharColumnType,
			Length:     255,
			IsNullable: true,
		},
		dto.ModelField{
			Name:   "channel",
			Type:   dto.VarcharColumnType,
			Length: 255,
		},
		dto.ModelField{
			Name: "event_id",
			Type: dto.IntegerColumnType,
		},
		dto.ModelField{
			Name: "scenario_id",
			Type: dto.IntegerColumnType,
		},
		dto.ModelField{
			Name: "variables",
			Type: dto.VarcharColumnType,
		},
		dto.ModelField{
			Name:   "reaction_type",
			Type:   dto.VarcharColumnType,
			Length: 255,
		},
		dto.ModelField{
			Name:   "execute_at",
			Type:   dto.VarcharColumnType,
			Length: 255,
		},
		dto.ModelField{
			Name:    "is_repeatable",
			Type:    dto.BooleanColumnType,
			Default: false,
		},
	},
	dto.ModelField{
		Name:          "id",
		Type:          dto.IntegerColumnType,
		AutoIncrement: true,
		IsPrimaryKey:  true,
	},
	&SchedulesStruct{},
)
