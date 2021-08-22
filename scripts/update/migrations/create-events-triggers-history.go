package migrations

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sharovik/devbot/internal/container"
	"github.com/sharovik/devbot/internal/dto/database_dto"
	"github.com/sharovik/orm/clients"
	"github.com/sharovik/orm/dto"
	"github.com/sharovik/orm/query"
)

type EventsTriggersHistoryMigration struct {
	Client clients.BaseClientInterface
}

func (m EventsTriggersHistoryMigration) SetClient(client clients.BaseClientInterface) {
	m.Client = client
}

func (m EventsTriggersHistoryMigration) GetName() string {
	return "create-events-triggers-history"
}

func (m EventsTriggersHistoryMigration) Execute() error {
	client := container.C.Dictionary.GetNewClient()

	//Create events table
	q := new(clients.Query).
		Create(&database_dto.EventTriggerHistoryModel).
		AddIndex(dto.Index{
			Name:   "user_id_index",
			Target: database_dto.EventTriggerHistoryModel.GetTableName(),
			Key:    "user_id",
			Unique: false,
		}).
		AddForeignKey(dto.ForeignKey{
			Name: "event_id",
			Target: query.Reference{
				Table: database_dto.EventModel.GetTableName(),
				Key:   "id",
			},
			With: query.Reference{
				Table: database_dto.EventTriggerHistoryModel.GetTableName(),
				Key:   "event_id",
			},
			OnDelete: dto.CascadeAction,
			OnUpdate: dto.NoActionAction,
		})
	if _, err := client.Execute(q); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create %s table", database_dto.EventTriggerHistoryModel.GetTableName()))
	}

	return nil
}
