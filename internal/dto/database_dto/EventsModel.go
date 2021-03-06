package database_dto

import "github.com/sharovik/orm/dto"

type EventsStruct struct {
	TableName  string
	PrimaryKey dto.ModelField
	Fields     []interface{}
}

var EventModel = dto.BaseModel{
	TableName: "events",
	Fields: []interface{}{
		dto.ModelField{
			Name:  "alias",
			Type:  dto.VarcharColumnType,
			Length: 255,
		},
		dto.ModelField{
			Name:  "installed_version",
			Type:  dto.VarcharColumnType,
			Length: 255,
		},
	},
	PrimaryKey: dto.ModelField{
		Name:          "id",
		Type:          dto.IntegerColumnType,
		AutoIncrement: true,
		IsPrimaryKey: true,
	},
}

func (m *EventsStruct) SetTableName(name string) {
	m.TableName = name
}

func (m EventsStruct) GetTableName() string {
	return m.TableName
}

func (m EventsStruct) GetColumns() []interface{} {
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

func (m *EventsStruct) AddModelField(field dto.ModelField) {
	m.Fields = append(m.GetColumns(), field)
}

func (m EventsStruct) GetField(name string) dto.ModelField {
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

func (m *EventsStruct) SetField(name string, value interface{}) {
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

func (m EventsStruct) GetPrimaryKey() dto.ModelField {
	m.PrimaryKey.IsPrimaryKey = true
	return m.PrimaryKey
}

func (m *EventsStruct) SetPrimaryKey(field dto.ModelField) {
	field.IsPrimaryKey = true
	m.PrimaryKey = field
}