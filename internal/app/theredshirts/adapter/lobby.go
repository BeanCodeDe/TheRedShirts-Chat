package adapter

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BeanCodeDe/TheRedShirts-Message/internal/app/theredshirts/util"
	"github.com/google/uuid"
)

type (
	LobbyAdapter struct {
		ServerUrl string
	}
	SimplePlayer struct {
		ID      uuid.UUID `json:"id" `
		Name    string    `json:"name" `
		LobbyId uuid.UUID `json:"lobby_id"`
	}
)

const (
	lobby_get_player_path     = "%s/player/%s"
	lobby_refresh_player_path = "%s/player/%s/last-refresh"
	correlation_id            = "X-Correlation-ID"
	content_typ_value         = "application/json; charset=utf-8"
	content_typ               = "Content-Type"
)

func NewLobbyAdapter() *LobbyAdapter {
	serverUrl := util.GetEnvWithFallback("CHAT_SERVER_URL", "http://theredshirts-lobby:1203")
	return &LobbyAdapter{ServerUrl: serverUrl}
}

func (adapter *LobbyAdapter) GetPlayer(context *util.Context, playerId uuid.UUID) (*SimplePlayer, error) {
	response, err := adapter.sendGetPlayer(context, playerId)
	if err != nil {
		return nil, fmt.Errorf("error while getting player: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status of response while getting player: %v", response.StatusCode)
	}

	var simplePlayer SimplePlayer
	err = json.NewDecoder(response.Body).Decode(&simplePlayer)
	if err != nil {
		return nil, fmt.Errorf("something went wrong while parsing player: %v", err)
	}
	return &simplePlayer, nil
}

func (adapter *LobbyAdapter) UpdatePlayerLastRefresh(context *util.Context, playerId uuid.UUID) error {
	response, err := adapter.sendUpdatePlayerLastRefresh(context, playerId)
	if err != nil {
		return fmt.Errorf("error while updating player: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("wrong status of response while updating player: %v", response.StatusCode)
	}

	return nil
}

func (adapter *LobbyAdapter) sendGetPlayer(context *util.Context, playerId uuid.UUID) (*http.Response, error) {
	client := &http.Client{}

	path := fmt.Sprintf(lobby_get_player_path, adapter.ServerUrl, playerId)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("request to get player could not be build: %v", err)
	}

	req.Header.Set(correlation_id, context.CorrelationId)
	req.Header.Set(content_typ, content_typ_value)
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("request to get player not possible: %v", err)
	}
	return resp, nil
}

func (adapter *LobbyAdapter) sendUpdatePlayerLastRefresh(context *util.Context, playerId uuid.UUID) (*http.Response, error) {
	client := &http.Client{}

	path := fmt.Sprintf(lobby_refresh_player_path, adapter.ServerUrl, playerId)
	req, err := http.NewRequest(http.MethodPatch, path, nil)
	if err != nil {
		return nil, fmt.Errorf("request to update last refesh for player chat could not be build: %v", err)
	}

	req.Header.Set(correlation_id, context.CorrelationId)
	req.Header.Set(content_typ, content_typ_value)
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("request to update last refesh for player not possible: %v", err)
	}
	return resp, nil
}
