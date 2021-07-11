package database_dto

import "github.com/sharovik/orm/dto"

type ScenariosStruct struct {
	TableName  string
	PrimaryKey dto.ModelField
	Fields     []interface{}
}

var ScenariosModel = dto.BaseModel{
	TableName: "scenarios",
	Fields: []interface{}{
		dto.ModelField{
			Name:  "name",
			Type:  dto.VarcharColumnType,
			Length: 255,
		},
		dto.ModelField{
			Name:  "event_id",
			Type:  dto.IntegerColumnType,
			IsNullable: true,
		},
	},
	PrimaryKey: dto.ModelField{
		Name:          "id",
		Type:          dto.IntegerColumnType,
		AutoIncrement: true,
		IsPrimaryKey: true,
	},
}

func (m *ScenariosStruct) SetTableName(name string) {
	m.TableName = name
}

func (m ScenariosStruct) GetTableName() string {
	return m.TableName
}

func (m ScenariosStruct) GetColumns() []interface{} {
	var columns []interface{}

	if m.GetPrimaryKey() != (dto.ModelField{IsPrimaryKey: true}) {
		columns = append(columns, m.GetPrimaryKey())
	}

	if len(m.Fields) == 0 {
		return nil
	}

	for _, field := range m.Fields {
		columns = append(columns, field)
	}
	return columns
}

func (m *ScenariosStruct) AddModelField(field dto.ModelField) {
	m.Fields = append(m.GetColumns(), field)
}

func (m ScenariosStruct) GetField(name string) dto.ModelField {
	for _, field := range m.GetColumns() {
		switch v := field.(type) {
		case dto.ModelField:
			if v.Name == name {
				return v
			}
		}
	}

	return dto.ModelField{}
}

func (m *ScenariosStruct) SetField(name string, value interface{}) {
	var columns []interface{}
	for _, field := range m.GetColumns() {
		switch v := field.(type) {
		case dto.ModelField:
			if m.GetPrimaryKey() == v {
				continue
			}

			if v.Name == name {
				v.Value = value
			}
			columns = append(columns, v)
		}
	}

	m.Fields = columns
}

func (m ScenariosStruct) GetPrimaryKey() dto.ModelField {
	m.PrimaryKey.IsPrimaryKey = true
	return m.PrimaryKey
}

func (m *ScenariosStruct) SetPrimaryKey(field dto.ModelField) {
	field.IsPrimaryKey = true
	m.PrimaryKey = field
}