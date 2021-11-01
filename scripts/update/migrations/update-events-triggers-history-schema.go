package migrations

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto/databasedto"
	"github.com/sharovik/orm/clients"
	"github.com/sharovik/orm/dto"
	"github.com/sharovik/orm/query"
)

type UpdateEventsTriggersHistoryMigration struct {
	Client clients.BaseClientInterface
}

func (m UpdateEventsTriggersHistoryMigration) SetClient(client clients.BaseClientInterface) {
	m.Client = client
}

func (m UpdateEventsTriggersHistoryMigration) GetName() string {
	return "update-events-triggers-history-schema"
}

func (m UpdateEventsTriggersHistoryMigration) Execute() error {
	client := container.C.Dictionary.GetNewClient()

	q := new(clients.Query).Drop(&databasedto.EventTriggerHistoryModel)
	if _, err := client.Execute(q); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to drop %s table", databasedto.EventTriggerHistoryModel.GetTableName()))
	}

	//Create events table
	q = new(clients.Query).
		Create(&databasedto.EventTriggerHistoryModel).
		AddIndex(dto.Index{
			Name:   "user_id_index",
			Target: databasedto.EventTriggerHistoryModel.GetTableName(),
			Key:    "user",
			Unique: false,
		}).
		AddIndex(dto.Index{
			Name:   "channel_index",
			Target: databasedto.EventTriggerHistoryModel.GetTableName(),
			Key:    "channel",
			Unique: false,
		}).
		AddIndex(dto.Index{
			Name:   "created_index",
			Target: databasedto.EventTriggerHistoryModel.GetTableName(),
			Key:    "created",
			Unique: false,
		}).
		AddForeignKey(dto.ForeignKey{
			Name: "event_id",
			Target: query.Reference{
				Table: databasedto.EventModel.GetTableName(),
				Key:   "id",
			},
			With: query.Reference{
				Table: databasedto.EventTriggerHistoryModel.GetTableName(),
				Key:   "event_id",
			},
			OnDelete: dto.CascadeAction,
			OnUpdate: dto.NoActionAction,
		}).
		AddForeignKey(dto.ForeignKey{
			Name: "scenario_id",
			Target: query.Reference{
				Table: databasedto.ScenariosModel.GetTableName(),
				Key:   "id",
			},
			With: query.Reference{
				Table: databasedto.EventTriggerHistoryModel.GetTableName(),
				Key:   "scenario_id",
			},
			OnDelete: dto.CascadeAction,
			OnUpdate: dto.NoActionAction,
		})
	if _, err := client.Execute(q); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create %s table", databasedto.EventTriggerHistoryModel.GetTableName()))
	}

	return nil
}
