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
		Create(&databasedto.EventTriggerHistoryModel).
		IfNotExists().
		AddIndex(dto.Index{
			Name:   "user_index",
			Target: databasedto.EventTriggerHistoryModel.GetTableName(),
			Key:    "user",
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
		})
	if _, err := client.Execute(q); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to create %s table", databasedto.EventTriggerHistoryModel.GetTableName()))
	}

	return nil
}
