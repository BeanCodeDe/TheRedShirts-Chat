package db

import (
	"errors"
	"strings"
	"time"

	"github.com/BeanCodeDe/TheRedShirts-Message/internal/app/theredshirts/util"
	"github.com/google/uuid"
)

type (
	Message struct {
		ID       uuid.UUID              `db:"id"`
		SendTime time.Time              `db:"send_time"`
		LobbyId  uuid.UUID              `db:"lobby_id"`
		PlayerId uuid.UUID              `db:"player_id"`
		Number   int                    `db:"number"`
		Topic    string                 `db:"topic"`
		Message  map[string]interface{} `db:"message"`
	}

	DB interface {
		Close()
		StartTransaction() (DBTx, error)
	}

	DBTx interface {
		Commit() error
		Rollback() error
		//Message
		CreateMessage(message *Message) error
		GetMessages(lobbyId uuid.UUID, toIgnoreplayerId uuid.UUID, number int) ([]*Message, error)
		GetMessagesFirstRequest(lobbyId uuid.UUID, toIgnoreplayerId uuid.UUID) ([]*Message, error)
		DeleteMessages(time time.Time) error
	}
)

const (
	schema_name = "theredshirts_message"
)

func NewConnection() (DB, error) {
	switch db := strings.ToLower(util.GetEnvWithFallback("DATABASE", "postgresql")); db {
	case "postgresql":
		return newPostgresConnection()
	default:
		return nil, errors.New("no configuration for %s found")
	}
}
