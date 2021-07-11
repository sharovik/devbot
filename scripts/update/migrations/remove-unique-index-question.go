package migrations

import (
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto/database_dto"
	"github.com/sharovik/orm/clients"
	"github.com/sharovik/orm/dto"
)

type RemoveUniqueIndexMigration struct {
	Client clients.BaseClientInterface
}

func (m RemoveUniqueIndexMigration) SetClient(client clients.BaseClientInterface) {
	m.Client = client
}

func (m RemoveUniqueIndexMigration) GetName() string {
	return "remove-unique-index-question"
}

func (m RemoveUniqueIndexMigration) Execute() error {
	client := container.C.Dictionary.GetNewClient()

	q := new(clients.Query).
		Alter(&database_dto.QuestionsModel).DropIndex(dto.Index{
		Name:   "questions_question_uindex",
	})
	_, err := client.Execute(q)
	if err != nil {
		return err
	}

	q = new(clients.Query).
		Alter(&database_dto.QuestionsModel).AddIndex(dto.Index{
		Name:   "questions_question_uindex",
		Target: database_dto.QuestionsModel.GetTableName(),
		Key:    "question",
		Unique: false,
	})
	_, err = client.Execute(q)
	if err != nil {
		return err
	}

	return nil
}
