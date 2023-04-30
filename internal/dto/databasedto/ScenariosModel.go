package databasedto

import "github.com/sharovik/orm/dto"

// ScenariosStruct struct for scenarios object
type ScenariosStruct struct {
	dto.BaseModel
}

// ScenariosModel model for scenarios table
var ScenariosModel = New(
	"scenarios",
	[]interface{}{
		dto.ModelField{
			Name:   "name",
			Type:   dto.VarcharColumnType,
			Length: 255,
		},
		dto.ModelField{
			Name:       "event_id",
			Type:       dto.IntegerColumnType,
			IsNullable: true,
		},
	},
	dto.ModelField{
		Name:          "id",
		Type:          dto.IntegerColumnType,
		AutoIncrement: true,
		IsPrimaryKey:  true,
	},
	&ScenariosStruct{},
)
