package databasedto

import "github.com/sharovik/orm/dto"

//EventsStruct struct for EventModel object
type EventsStruct struct {
	dto.BaseModel
}

//EventModel model for events table
var EventModel = New(
	"events",
	[]interface{}{
		dto.ModelField{
			Name:   "alias",
			Type:   dto.VarcharColumnType,
			Length: 255,
		},
		dto.ModelField{
			Name:   "installed_version",
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
	&EventsStruct{},
)
