package migrations

import (
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto/database_dto"
	"github.com/sharovik/orm/clients"
	"github.com/sharovik/orm/query"
)

type ExampleMigration struct {
	Client clients.BaseClientInterface
}

func (m ExampleMigration) SetClient(client clients.BaseClientInterface) {
	m.Client = client
}

func (m ExampleMigration) GetName() string {
	return "example-migration"
}

func (m ExampleMigration) Execute() error {
	client := container.C.Dictionary.GetNewClient()

	q := new(clients.Query).
		Select(database_dto.MigrationModel.GetColumns()).
		From(&database_dto.MigrationModel).
		Where(query.Where{
			First:    "1",
			Operator: "=",
			Second:   "1",
		})
	_, err := client.Execute(q)
	if err != nil {
		return err
	}

	return nil
}
