package core

import (
	"errors"
	"fmt"
	"time"

	"github.com/BeanCodeDe/TheRedShirts-Chat/internal/app/theredshirts/db"
	"github.com/google/uuid"
)

func (core CoreFacade) CreatePlayer(player *Player) error {
	tx, err := core.db.StartTransaction()
	defer tx.HandleTransaction(err)
	if err != nil {
		return fmt.Errorf("something went wrong while creating transaction: %v", err)
	}
	err = core.createPlayer(tx, player)
	return err
}

func (core CoreFacade) createPlayer(tx db.DBTx, player *Player) error {
	if err := tx.CreatePlayer(mapToDBPlayer(player)); err != nil {
		if !errors.Is(err, db.ErrPlayerAlreadyExists) {
			return fmt.Errorf("error while creating player: %v", err)
		}
	}

	message := &Message{ID: uuid.New(), SendTime: time.Now(), PlayerName: "Lobby", LobbyId: player.LobbyId, Message: fmt.Sprintf("Player %s joint team %s", player.Name, player.Team)}
	if err := tx.CreateMessage(mapToDBMessage(message)); err != nil {
		if !errors.Is(err, db.ErrMessageAlreadyExists) {
			return fmt.Errorf("error while creating message: %v", err)
		}
	}

	return nil
}

func (core CoreFacade) DeletePlayer(playerId uuid.UUID) error {
	tx, err := core.db.StartTransaction()
	defer tx.HandleTransaction(err)
	if err != nil {
		return fmt.Errorf("something went wrong while creating transaction: %v", err)
	}
	err = core.deletePlayer(tx, playerId)
	return err
}

func (core CoreFacade) deletePlayer(tx db.DBTx, playerId uuid.UUID) error {

	player, err := core.getPlayer(tx, playerId)
	if err != nil {
		return err
	}

	message := &Message{ID: uuid.New(), SendTime: time.Now(), PlayerName: "Lobby", LobbyId: player.LobbyId, Message: fmt.Sprintf("Player %s left lobby", player.Name)}
	if err := tx.CreateMessage(mapToDBMessage(message)); err != nil {
		if !errors.Is(err, db.ErrMessageAlreadyExists) {
			return fmt.Errorf("error while creating message: %v", err)
		}
	}

	if err := tx.DeletePlayer(playerId); err != nil {
		return fmt.Errorf("an error accourd while deleting player [%v]: %v", playerId, err)
	}

	return nil
}

func (core CoreFacade) GetPlayer(playerId uuid.UUID) (*Player, error) {
	tx, err := core.db.StartTransaction()
	defer tx.HandleTransaction(err)

	player, err := core.getPlayer(tx, playerId)
	return player, err
}

func (core CoreFacade) getPlayer(tx db.DBTx, playerId uuid.UUID) (*Player, error) {
	player, err := tx.GetPlayer(playerId)
	if err != nil {
		return nil, fmt.Errorf("something went wrong while loading player [%v] from database: %v", playerId, err)
	}

	return mapToPlayer(player), nil
}

func mapToPlayer(player *db.Player) *Player {
	return &Player{ID: player.ID, LobbyId: player.LobbyId, Name: player.Name, Team: player.Team}
}

func mapToDBPlayer(player *Player) *db.Player {
	return &db.Player{ID: player.ID, LobbyId: player.LobbyId, Name: player.Name, Team: player.Team}
}
