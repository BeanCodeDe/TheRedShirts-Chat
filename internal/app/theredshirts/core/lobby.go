package core

import (
	"errors"
	"fmt"

	"github.com/BeanCodeDe/TheRedShirts-Lobby/internal/app/theredshirts/db"
	"github.com/google/uuid"
)

func (core CoreFacade) CreateLobby(lobby *Lobby) error {
	dbLobby := mapToDBLobby(lobby)

	if err := core.db.CreateLobby(dbLobby); err != nil {
		if !errors.Is(err, db.ErrLobbyAlreadyExists) {
			return fmt.Errorf("error while creating lobby: %v", err)
		}
		foundLobby, err := core.db.GetLobbyById(lobby.ID)
		if err != nil {
			return fmt.Errorf("something went wrong while checking if lobby [%v] is already created: %v", lobby.ID, err)
		}

		if lobby.Name != foundLobby.Name || lobby.Password != foundLobby.Password {
			return fmt.Errorf("request of lobby [%v] doesn't match lobby from database [%v]", lobby, foundLobby)
		}

	}

	if err := core.JoinLobby(&Join{PlayerId: lobby.Owner.ID, LobbyID: lobby.ID, Name: lobby.Owner.Name, Password: lobby.Password}); err != nil {
		return err
	}

	return nil
}

func (core CoreFacade) UpdateLobby(lobby *Lobby) error {
	dbLobby, err := core.db.GetLobbyById(lobby.ID)
	if err != nil {
		return fmt.Errorf("something went wrong while loading lobby [%v] from database: %v", lobby.ID, err)
	}

	dbLobby.Name = lobby.Name
	dbLobby.Difficulty = lobby.Difficulty
	dbLobby.Password = lobby.Password

	if err := core.db.UpdateLobby(dbLobby); err != nil {
		if err != nil {
			return fmt.Errorf("something went wrong while updating lobby [%v]: %v", lobby.ID, err)
		}
	}
	return nil
}

func (core CoreFacade) DeleteLobby(lobbyId uuid.UUID) error {

	if err := core.db.DeleteAllPlayerInLobby(lobbyId); err != nil {
		return fmt.Errorf("an error accourd while deleting all players from lobby [%v]: %v", lobbyId, err)
	}

	if err := core.db.DeleteLobby(lobbyId); err != nil {
		return fmt.Errorf("an error accourd while deleting lobby [%v]: %v", lobbyId, err)
	}
	return nil
}

func (core CoreFacade) GetLobby(lobbyId uuid.UUID) (*Lobby, error) {
	lobby, err := core.db.GetLobbyById(lobbyId)
	if err != nil {
		return nil, fmt.Errorf("something went wrong while loading lobby [%v] from database: %v", lobby.ID, err)
	}

	players, err := core.db.GetAllPlayersInLobby(lobbyId)
	if err != nil {
		return nil, fmt.Errorf("something went wrong while loading players of lobby [%v] from database: %v", lobby.ID, err)
	}

	owner, err := core.db.GetPlayerById(lobby.Owner)
	if err != nil {
		return nil, fmt.Errorf("something went wrong while loading owner [%v] of lobby [%v] from database: %v", lobby.Owner, lobby.ID, err)
	}

	return mapToLobby(lobby, mapToPlayer(owner), mapToPlayers(players)), nil
}

func (core CoreFacade) GetLobbies() ([]*Lobby, error) {
	lobbies, err := core.db.GetAllLobbies()
	if err != nil {
		return nil, fmt.Errorf("something went wrong while loading all lobbies from database: %v", err)
	}

	coreLobbies := make([]*Lobby, len(lobbies))
	for index, lobby := range lobbies {

		players, err := core.db.GetAllPlayersInLobby(lobby.ID)
		if err != nil {
			return nil, fmt.Errorf("something went wrong while loading players of lobby [%v] from database: %v", lobby.ID, err)
		}

		owner, err := core.db.GetPlayerById(lobby.Owner)
		if err != nil {
			return nil, fmt.Errorf("something went wrong while loading owner [%v] of lobby [%v] from database: %v", lobby.Owner, lobby.ID, err)
		}
		coreLobbies[index] = mapToLobby(lobby, mapToPlayer(owner), mapToPlayers(players))
	}

	return coreLobbies, nil
}

func (core CoreFacade) JoinLobby(join *Join) error {
	if err := core.LeaveLobby(join.PlayerId); err != nil {
		return err
	}

	lobby, err := core.db.GetLobbyById(join.LobbyID)
	if err != nil {
		return fmt.Errorf("something went wrong while loading lobby %v from database: %v", join.LobbyID, err)
	}

	if lobby.Password != join.Password {
		return ErrWrongLobbyPassword
	}

	if err := core.db.CreatePlayer(&db.Player{ID: join.PlayerId, Name: join.Name, LobbyId: join.LobbyID}); err != nil {
		return fmt.Errorf("something went wrong while creating player %v from database: %v", join.PlayerId, err)
	}

	return nil
}

func (core CoreFacade) LeaveLobby(playerId uuid.UUID) error {
	if err := core.db.DeletePlayer(playerId); err != nil {
		return fmt.Errorf("something went wrong while deleting player %v from database: %v", playerId, err)
	}
	return nil
}
