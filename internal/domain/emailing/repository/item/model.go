package item

import (
	"encoding/json"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/neurochar/backend/internal/domain/emailing/entity"
	"github.com/neurochar/backend/internal/infra/emailing"
)

// DBModel - database model
type DBModel struct {
	ID          uuid.UUID       `db:"id"`
	MessageData json.RawMessage `db:"message_data"`
	RequestIP   net.IP          `db:"request_ip"`
	SentAt      *time.Time      `db:"sent_at"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (db *DBModel) ToEntity() *entity.Item {
	messageData := emailing.Message{}
	err := json.Unmarshal(db.MessageData, &messageData)
	if err != nil {
		panic(err)
	}

	return &entity.Item{
		ID:          db.ID,
		MessageData: messageData,
		RequestIP:   db.RequestIP,
		SentAt:      db.SentAt,

		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
	}
}

func mapEntityToDBModel(entity *entity.Item) *DBModel {
	messageData, err := json.Marshal(entity.MessageData)
	if err != nil {
		panic(err)
	}

	return &DBModel{
		ID:          entity.ID,
		MessageData: messageData,
		RequestIP:   entity.RequestIP,
		SentAt:      entity.SentAt,

		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}
