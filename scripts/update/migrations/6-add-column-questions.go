package migrations

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto/databasedto"
	"github.com/sharovik/orm/clients"
	"github.com/sharovik/orm/dto"
)

type AddColumnQuestions struct {
	Client clients.BaseClientInterface
}

func (m AddColumnQuestions) SetClient(client clients.BaseClientInterface) {
	m.Client = client
}

func (m AddColumnQuestions) GetName() string {
	return "6-add-column-questions"
}

func (m AddColumnQuestions) Execute() error {
	client := container.C.Dictionary.GetDBClient()

	q := new(clients.Query).Drop(databasedto.EventTriggerHistoryModel)
	if _, err := client.Execute(q); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to drop %s table", databasedto.EventTriggerHistoryModel.GetTableName()))
	}

	//Create events table
	q = new(clients.Query).
		Alter(databasedto.QuestionsModel).
		AddColumn(dto.ModelField{
			Name:    "is_variable",
			Type:    dto.BooleanColumnType,
			Value:   nil,
			Default: false,
		}).
		AddIndex(dto.Index{
			Name:   "is_variable_index",
			Target: databasedto.QuestionsModel.GetTableName(),
			Key:    "is_variable",
			Unique: false,
		})
	if _, err := client.Execute(q); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create %s table", databasedto.EventTriggerHistoryModel.GetTableName()))
	}

	return nil
}
