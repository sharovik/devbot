package databasedto

import "github.com/sharovik/orm/dto"

//EventTriggerHistoryStruct the struct for event history table
type EventTriggerHistoryStruct struct {
	TableName  string
	PrimaryKey dto.ModelField
	Fields     []interface{}
}

//EventTriggerHistoryModel the actual object of event history table
var EventTriggerHistoryModel = dto.BaseModel{
	TableName: "events_triggers_history",
	Fields: []interface{}{
		dto.ModelField{
			Name:  "event_id",
			Type:  dto.IntegerColumnType,
		},
		dto.ModelField{
			Name:  "scenario_id",
			Type:  dto.IntegerColumnType,
		},
		dto.ModelField{
			Name:  "user",
			Type:  dto.VarcharColumnType,
			Length: 255,
		},
		dto.ModelField{
			Name:  "channel",
			Type:  dto.VarcharColumnType,
			Length: 255,
		},
		dto.ModelField{
			Name:  "command",
			Type:  dto.VarcharColumnType,
			Length: 255,
		},
		dto.ModelField{
			Name:  "variables",
			Type:  dto.VarcharColumnType,
			Length: 255,
			IsNullable: true,
		},
		dto.ModelField{
			Name:  "last_question_id",
			Type:  dto.IntegerColumnType,
			Length: 11,
			IsNullable: true,
		},
		dto.ModelField{
			Name:  "created",
			Type:  dto.IntegerColumnType,
			Length: 11,
		},
	},
	PrimaryKey: dto.ModelField{
		Name:          "id",
		Type:          dto.IntegerColumnType,
		AutoIncrement: true,
		IsPrimaryKey: true,
	},
}

func (m *EventTriggerHistoryStruct) SetTableName(name string) {
	m.TableName = name
}

func (m EventTriggerHistoryStruct) GetTableName() string {
	return m.TableName
}

func (m EventTriggerHistoryStruct) GetColumns() []interface{} {
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

func (m *EventTriggerHistoryStruct) AddModelField(field dto.ModelField) {
	m.Fields = append(m.GetColumns(), field)
}

func (m EventTriggerHistoryStruct) GetField(name string) dto.ModelField {
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

func (m *EventTriggerHistoryStruct) SetField(name string, value interface{}) {
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

func (m EventTriggerHistoryStruct) GetPrimaryKey() dto.ModelField {
	m.PrimaryKey.IsPrimaryKey = true
	return m.PrimaryKey
}

func (m *EventTriggerHistoryStruct) SetPrimaryKey(field dto.ModelField) {
	field.IsPrimaryKey = true
	m.PrimaryKey = field
}

func (m *EventTriggerHistoryStruct) RemoveModelField(field string) {
	var columns []interface{}
	for _, f := range m.Fields {
		switch v := f.(type) {
		case dto.ModelField:
			if field == v.Name {
				continue
			}
		}

		columns = append(columns, f)
	}

	m.Fields = columns
}