package db

import (
	"errors"
	"strings"
	"time"

	"github.com/BeanCodeDe/TheRedShirts-Chat/internal/app/theredshirts/util"
	"github.com/google/uuid"
)

type (
	Message struct {
		ID       uuid.UUID              `db:"id"`
		SendTime time.Time              `db:"send_time"`
		LobbyId  uuid.UUID              `db:"lobby_id"`
		Number   int                    `db:"number"`
		Topic    string                 `db:"topic"`
		Message  map[string]interface{} `db:"message"`
	}

	Player struct {
		ID          uuid.UUID `db:"id"`
		LobbyId     uuid.UUID `db:"lobby_id"`
		LastRefresh time.Time `db:"last_refresh"`
	}

	DB interface {
		Close()
		StartTransaction() (DBTx, error)
	}

	DBTx interface {
		HandleTransaction(err error)
		//Message
		CreateMessage(message *Message) error
		GetMessages(lobbyId uuid.UUID, number int) ([]*Message, error)

		//Player
		CreatePlayer(player *Player) error
		UpdatePlayerLastRefresh(id uuid.UUID, lastRefresh time.Time) error
		DeletePlayer(id uuid.UUID) error
		GetPlayer(id uuid.UUID) (*Player, error)
	}
)

const (
	schema_name = "theredshirts_chat"
)

func NewConnection() (DB, error) {
	switch db := strings.ToLower(util.GetEnvWithFallback("DATABASE", "postgresql")); db {
	case "postgresql":
		return newPostgresConnection()
	default:
		return nil, errors.New("no configuration for %s found")
	}
}
