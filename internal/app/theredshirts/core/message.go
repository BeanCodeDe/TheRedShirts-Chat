package core

import (
	"errors"
	"fmt"
	"time"

	"github.com/BeanCodeDe/TheRedShirts-Chat/internal/app/theredshirts/db"
	"github.com/google/uuid"
)

func (core CoreFacade) CreateMessage(playerId uuid.UUID, message *Message) error {
	tx, err := core.db.StartTransaction()
	defer tx.HandleTransaction(err)
	if err != nil {
		return fmt.Errorf("something went wrong while creating transaction: %v", err)
	}
	err = core.createMessage(tx, playerId, message)
	return err
}

func (core CoreFacade) createMessage(tx db.DBTx, playerId uuid.UUID, message *Message) error {
	player, err := tx.GetPlayer(playerId)
	if err != nil {
		return fmt.Errorf("error while getting player %v: %v", playerId, err)
	}

	if player.LobbyId != message.LobbyId {
		return fmt.Errorf("error player %v from lobby %v is not authorised to write in lobby %v", playerId, player.LobbyId, message.LobbyId)
	}

	if err := tx.CreateMessage(mapToDBMessage(message)); err != nil {
		if !errors.Is(err, db.ErrMessageAlreadyExists) {
			return fmt.Errorf("error while creating message: %v", err)
		}
	}
	return nil
}

func (core CoreFacade) GetMessages(playerId uuid.UUID, lobbyId uuid.UUID, number int) ([]*Message, error) {
	tx, err := core.db.StartTransaction()
	defer tx.HandleTransaction(err)

	messages, err := core.getMessages(tx, playerId, lobbyId, number)
	return messages, err
}

func (core CoreFacade) getMessages(tx db.DBTx, playerId uuid.UUID, lobbyId uuid.UUID, number int) ([]*Message, error) {
	player, err := tx.GetPlayer(playerId)
	if err != nil {
		return nil, fmt.Errorf("error while getting player %v: %v", playerId, err)
	}

	if player.LobbyId != lobbyId {
		return nil, fmt.Errorf("error player %v from lobby %v is not authorised to write in lobby %v", playerId, player.LobbyId, lobbyId)
	}
	messages, err := tx.GetMessages(lobbyId, number)
	if err != nil {
		return nil, fmt.Errorf("something went wrong while loading messages in lobby [%v] from database: %v", lobbyId, err)
	}

	if err := core.updatePlayerLastRefresh(tx, playerId); err != nil {
		return nil, err
	}

	return mapToMessages(messages), nil
}

func mapToMessages(dbMessages []*db.Message) []*Message {
	messages := make([]*Message, len(dbMessages))
	for index, message := range dbMessages {
		messages[index] = mapToMessage(message)
	}
	return messages
}

func mapToMessage(message *db.Message) *Message {
	return &Message{ID: message.ID, SendTime: message.SendTime, LobbyId: message.LobbyId, Number: message.Number, Topic: message.Topic, Message: message.Message}
}

func mapToDBMessage(message *Message) *db.Message {
	return &db.Message{ID: message.ID, SendTime: message.SendTime, LobbyId: message.LobbyId, Number: message.Number, Topic: message.Topic, Message: message.Message}
}

func getPlayerJoinsMessage(lobbyId uuid.UUID) *Message {
	basicMessage := getBasicMessage(lobbyId)
	basicMessage.Topic = "PLAYER_JOIN"
	return basicMessage
}

func getPlayerLeavesMessage(lobbyId uuid.UUID) *Message {
	basicMessage := getBasicMessage(lobbyId)
	basicMessage.Topic = "PLAYER_LEAVES"
	return basicMessage
}

func getBasicMessage(lobbyId uuid.UUID) *Message {
	return &Message{ID: uuid.New(), SendTime: time.Now(), LobbyId: lobbyId}
}
