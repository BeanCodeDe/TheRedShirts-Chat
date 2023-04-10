package core

import (
	"errors"
	"fmt"

	"github.com/BeanCodeDe/TheRedShirts-Message/internal/app/theredshirts/db"
	"github.com/BeanCodeDe/TheRedShirts-Message/internal/app/theredshirts/util"
	"github.com/google/uuid"
)

func (core CoreFacade) CreateMessage(context *util.Context, message *Message) error {
	context.Logger.Debugf("Create Message: %+v", *message)
	tx, err := core.db.StartTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err != nil {
		return fmt.Errorf("something went wrong while creating transaction: %v", err)
	}
	if err := core.createMessage(context, tx, message); err != nil {
		return err
	}
	return tx.Commit()
}

func (core CoreFacade) createMessage(context *util.Context, tx db.DBTx, message *Message) error {
	if message.PlayerId != core.lobbyPlayerId {
		player, err := core.lobbyAdapter.GetPlayer(context, message.PlayerId)
		if err != nil {
			return fmt.Errorf("error while getting player %v: %v", message.PlayerId, err)
		}

		if player.LobbyId != message.LobbyId {
			return fmt.Errorf("error player %v from lobby %v is not authorised to write in lobby %v", message.PlayerId, player.LobbyId, message.LobbyId)
		}
	}

	if err := tx.CreateMessage(mapToDBMessage(message)); err != nil {
		if !errors.Is(err, db.ErrMessageAlreadyExists) {
			return fmt.Errorf("error while creating message: %v", err)
		}
	}
	return nil
}

func (core CoreFacade) GetMessages(context *util.Context, playerId uuid.UUID, lobbyId uuid.UUID, number int) ([]*Message, error) {
	tx, err := core.db.StartTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	messages, err := core.getMessages(context, tx, playerId, lobbyId, number)
	if err != nil {
		return nil, err
	}
	return messages, tx.Commit()
}

func (core CoreFacade) getMessages(context *util.Context, tx db.DBTx, playerId uuid.UUID, lobbyId uuid.UUID, number int) ([]*Message, error) {
	player, err := core.lobbyAdapter.GetPlayer(context, playerId)
	if err != nil {
		return nil, fmt.Errorf("error while getting player %v: %v", playerId, err)
	}

	if player.LobbyId != lobbyId {
		return nil, fmt.Errorf("error player %v from lobby %v is not authorised to load messages from lobby %v", playerId, player.LobbyId, lobbyId)
	}
	var messages []*db.Message
	if number != -1 {
		messages, err = tx.GetMessages(lobbyId, playerId, number)
	} else {
		messages, err = tx.GetMessagesFirstRequest(lobbyId, playerId)
	}
	if err != nil {
		return nil, fmt.Errorf("something went wrong while loading messages in lobby [%v] from database: %v", lobbyId, err)
	}

	if err := core.lobbyAdapter.UpdatePlayerLastRefresh(context, playerId); err != nil {
		return nil, fmt.Errorf("error while updating player %v: %v", playerId, err)
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
	return &Message{ID: message.ID, SendTime: message.SendTime, LobbyId: message.LobbyId, PlayerId: message.PlayerId, Number: message.Number, Topic: message.Topic, Message: message.Message}
}

func mapToDBMessage(message *Message) *db.Message {
	return &db.Message{ID: message.ID, SendTime: message.SendTime, LobbyId: message.LobbyId, PlayerId: message.PlayerId, Number: message.Number, Topic: message.Topic, Message: message.Message}
}
