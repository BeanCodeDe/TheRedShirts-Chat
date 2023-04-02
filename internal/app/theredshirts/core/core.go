package core

import (
	"errors"
	"fmt"
	"time"

	"github.com/BeanCodeDe/TheRedShirts-Chat/internal/app/theredshirts/db"
	"github.com/google/uuid"
)

type (

	//Facade
	CoreFacade struct {
		db db.DB
	}

	Core interface {
		//Message
		CreateMessage(playerId uuid.UUID, message *Message) error
		GetMessages(playerId uuid.UUID, lobbyId uuid.UUID, number int) ([]*Message, error)

		//Player
		CreatePlayer(player *Player) error
		DeletePlayer(playerId uuid.UUID) error
		GetPlayer(playerId uuid.UUID) (*Player, error)
	}

	//Objects
	Message struct {
		ID       uuid.UUID
		SendTime time.Time
		LobbyId  uuid.UUID
		Number   int
		Topic    string
		Message  map[string]interface{}
	}

	Player struct {
		ID          uuid.UUID
		LobbyId     uuid.UUID
		LastRefresh time.Time
	}
)

var (
	ErrWrongLobbyPassword = errors.New("wrong password")
)

func NewCore() (Core, error) {
	db, err := db.NewConnection()
	if err != nil {
		return nil, fmt.Errorf("error while initializing database: %v", err)
	}
	core := &CoreFacade{db: db}
	return core, nil
}
