package databasedto

import "github.com/sharovik/orm/dto"

// MigrationStruct struct for MigrationModel object
type MigrationStruct struct {
	dto.BaseModel
}

// MigrationModel model for migrations table
var MigrationModel = New(
	"migration",
	[]interface{}{
		dto.ModelField{
			Name:   "version",
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
	&MigrationStruct{},
)
