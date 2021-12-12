package databasedto

import "github.com/sharovik/orm/dto"

type QuestionsRegexStruct struct {
	TableName  string
	PrimaryKey dto.ModelField
	Fields     []interface{}
}

var QuestionsRegexModel = dto.BaseModel{
	TableName: "questions_regex",
	Fields: []interface{}{
		dto.ModelField{
			Name:  "regex",
			Type:  dto.VarcharColumnType,
			Length: 255,
		},
		dto.ModelField{
			Name:  "regex_group",
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

func (m *QuestionsRegexStruct) SetTableName(name string) {
	m.TableName = name
}

func (m QuestionsRegexStruct) GetTableName() string {
	return m.TableName
}

func (m QuestionsRegexStruct) GetColumns() []interface{} {
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

func (m *QuestionsRegexStruct) AddModelField(field dto.ModelField) {
	m.Fields = append(m.GetColumns(), field)
}

func (m QuestionsRegexStruct) GetField(name string) dto.ModelField {
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

func (m *QuestionsRegexStruct) SetField(name string, value interface{}) {
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

func (m QuestionsRegexStruct) GetPrimaryKey() dto.ModelField {
	m.PrimaryKey.IsPrimaryKey = true
	return m.PrimaryKey
}

func (m *QuestionsRegexStruct) SetPrimaryKey(field dto.ModelField) {
	field.IsPrimaryKey = true
	m.PrimaryKey = field
}

func (m *QuestionsRegexStruct) RemoveModelField(field string) {
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