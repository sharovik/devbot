package database

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto/database_dto"
	"github.com/sharovik/orm/clients"
	"github.com/sharovik/orm/dto"
	cquery "github.com/sharovik/orm/query"
)

type InstallMigration struct {
	Client clients.BaseClientInterface
}

func (m InstallMigration) SetClient(client clients.BaseClientInterface) {
	m.Client = client
}

func (m InstallMigration) GetName() string {
	return "install"
}

func (m InstallMigration) Execute() error {
	client := container.C.Dictionary.GetNewClient()

	q := new(clients.Query).Select([]interface{}{"*"}).From(&database_dto.MigrationModel)
	if _, err := client.Execute(q); err == nil {
		//We already have triggered the schema setup
		return nil
	}

	if err := createSchema(client); err != nil {
		return err
	}

	return nil
}

func createSchema(client clients.BaseClientInterface) error {
	//Create events table
	q := new(clients.Query).
		Create(&database_dto.EventModel).
		AddIndex(dto.Index{
			Name:   "events_name_uindex",
			Target: database_dto.EventModel.GetTableName(),
			Key:    "alias",
			Unique: true,
		})
	if _, err := client.Execute(q); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create %s table", database_dto.EventModel.GetTableName()))
	}

	//Create migrations table
	q = new(clients.Query).
		Create(&database_dto.MigrationModel).
		AddIndex(dto.Index{
			Name:   "migration_version_uindex",
			Target: database_dto.MigrationModel.GetTableName(),
			Key:    "version",
			Unique: true,
		})
	if _, err := client.Execute(q); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create %s table", database_dto.MigrationModel.GetTableName()))
	}

	//Create questions_regex table
	q = new(clients.Query).
		Create(&database_dto.QuestionsRegexModel).
		AddIndex(dto.Index{
			Name:   "question_regex_regex_uindex",
			Target: database_dto.QuestionsRegexModel.GetTableName(),
			Key:    "regex",
			Unique: true,
		})
	if _, err := client.Execute(q); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create %s table", database_dto.QuestionsRegexModel.GetTableName()))
	}

	//Create scenarios table
	q = new(clients.Query).
		Create(&database_dto.ScenariosModel).
		AddIndex(dto.Index{
			Name:   "scenarios_name_uindex",
			Target: database_dto.ScenariosModel.GetTableName(),
			Key:    "name",
			Unique: true,
		}).
		AddForeignKey(dto.ForeignKey{
			Name: "event_id",
			Target: cquery.Reference{
				Table: database_dto.EventModel.GetTableName(),
				Key:   "id",
			},
			With: cquery.Reference{
				Table: database_dto.ScenariosModel.GetTableName(),
				Key:   "event_id",
			},
			OnDelete: dto.SetNullAction,
			OnUpdate: dto.NoActionAction,
		})
	if _, err := client.Execute(q); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create %s table", database_dto.ScenariosModel.GetTableName()))
	}

	//Create questions table
	q = new(clients.Query).
		Create(&database_dto.QuestionsModel).
		AddIndex(dto.Index{
			Name:   "questions_question_index",
			Target: database_dto.QuestionsModel.GetTableName(),
			Key:    "question",
			Unique: false,
		}).
		AddForeignKey(dto.ForeignKey{
			Name: "scenario_id",
			Target: cquery.Reference{
				Table: database_dto.ScenariosModel.GetTableName(),
				Key:   "id",
			},
			With: cquery.Reference{
				Table: database_dto.QuestionsModel.GetTableName(),
				Key:   "scenario_id",
			},
			OnDelete: dto.CascadeAction,
			OnUpdate: dto.NoActionAction,
		}).
		AddForeignKey(dto.ForeignKey{
			Name: "regex_id",
			Target: cquery.Reference{
				Table: database_dto.QuestionsRegexModel.GetTableName(),
				Key:   "id",
			},
			With: cquery.Reference{
				Table: database_dto.QuestionsModel.GetTableName(),
				Key:   "regex_id",
			},
			OnDelete: dto.SetNullAction,
			OnUpdate: dto.NoActionAction,
		})
	if _, err := client.Execute(q); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create %s table", database_dto.QuestionsModel.GetTableName()))
	}

	return nil
}
