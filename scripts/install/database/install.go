package database

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto/databasedto"
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
	client := container.C.Dictionary.GetDBClient()

	if err := createSchema(client); err != nil {
		return err
	}

	return nil
}

func createSchema(client clients.BaseClientInterface) error {
	//Create events table
	q := new(clients.Query).
		Create(databasedto.EventModel).
		IfNotExists().
		AddIndex(dto.Index{
			Name:   "events_name_uindex",
			Target: databasedto.EventModel.GetTableName(),
			Key:    "alias",
			Unique: true,
		})
	if _, err := client.Execute(q); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create %s table", databasedto.EventModel.GetTableName()))
	}

	//Create questions_regex table
	q = new(clients.Query).
		Create(databasedto.QuestionsRegexModel).
		IfNotExists().
		AddIndex(dto.Index{
			Name:   "question_regex_regex_uindex",
			Target: databasedto.QuestionsRegexModel.GetTableName(),
			Key:    "regex",
			Unique: true,
		})
	if _, err := client.Execute(q); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create %s table", databasedto.QuestionsRegexModel.GetTableName()))
	}

	//Create scenarios table
	q = new(clients.Query).
		Create(databasedto.ScenariosModel).
		IfNotExists().
		AddIndex(dto.Index{
			Name:   "scenarios_name_uindex",
			Target: databasedto.ScenariosModel.GetTableName(),
			Key:    "name",
			Unique: true,
		}).
		AddForeignKey(dto.ForeignKey{
			Name: "event_id",
			Target: cquery.Reference{
				Table: databasedto.EventModel.GetTableName(),
				Key:   "id",
			},
			With: cquery.Reference{
				Table: databasedto.ScenariosModel.GetTableName(),
				Key:   "event_id",
			},
			OnDelete: dto.SetNullAction,
			OnUpdate: dto.NoActionAction,
		})
	if _, err := client.Execute(q); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create %s table", databasedto.ScenariosModel.GetTableName()))
	}

	//Create questions table
	q = new(clients.Query).
		Create(databasedto.QuestionsModel).
		IfNotExists().
		AddIndex(dto.Index{
			Name:   "questions_question_index",
			Target: databasedto.QuestionsModel.GetTableName(),
			Key:    "question",
			Unique: false,
		}).
		AddForeignKey(dto.ForeignKey{
			Name: "scenario_id",
			Target: cquery.Reference{
				Table: databasedto.ScenariosModel.GetTableName(),
				Key:   "id",
			},
			With: cquery.Reference{
				Table: databasedto.QuestionsModel.GetTableName(),
				Key:   "scenario_id",
			},
			OnDelete: dto.CascadeAction,
			OnUpdate: dto.NoActionAction,
		}).
		AddForeignKey(dto.ForeignKey{
			Name: "regex_id",
			Target: cquery.Reference{
				Table: databasedto.QuestionsRegexModel.GetTableName(),
				Key:   "id",
			},
			With: cquery.Reference{
				Table: databasedto.QuestionsModel.GetTableName(),
				Key:   "regex_id",
			},
			OnDelete: dto.SetNullAction,
			OnUpdate: dto.NoActionAction,
		})
	if _, err := client.Execute(q); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create %s table", databasedto.QuestionsModel.GetTableName()))
	}

	//Create questions table
	q = new(clients.Query).
		Create(databasedto.SchedulesModel).
		IfNotExists().
		AddIndex(dto.Index{
			Name:   "execute_at_index",
			Target: databasedto.SchedulesModel.GetTableName(),
			Key:    "execute_at",
			Unique: false,
		}).
		AddForeignKey(dto.ForeignKey{
			Name: "scenario_id",
			Target: cquery.Reference{
				Table: databasedto.ScenariosModel.GetTableName(),
				Key:   "id",
			},
			With: cquery.Reference{
				Table: databasedto.SchedulesModel.GetTableName(),
				Key:   "scenario_id",
			},
			OnDelete: dto.SetNullAction,
			OnUpdate: dto.NoActionAction,
		}).
		AddForeignKey(dto.ForeignKey{
			Name: "event_id",
			Target: cquery.Reference{
				Table: databasedto.EventModel.GetTableName(),
				Key:   "id",
			},
			With: cquery.Reference{
				Table: databasedto.SchedulesModel.GetTableName(),
				Key:   "event_id",
			},
			OnDelete: dto.SetNullAction,
			OnUpdate: dto.NoActionAction,
		})
	if _, err := client.Execute(q); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create %s table", databasedto.SchedulesModel.GetTableName()))
	}

	return nil
}
