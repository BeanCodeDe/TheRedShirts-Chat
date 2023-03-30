package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/BeanCodeDe/TheRedShirts-Chat/internal/app/theredshirts/core"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

const chat_root_path = "/chat"
const message_path = "/message"
const player_path = "/player"
const message_id_param = "messageId"
const number_id_param = "number"
const lobby_id_param = "lobbyId"
const player_id_param = "playerId"

type (
	MessageCreate struct {
		PlayerId uuid.UUID `json:"player_id" validate:"required"`
		Message  string    `json:"message" validate:"required"`
	}

	Message struct {
		ID         uuid.UUID `json:"id"`
		SendTime   time.Time `json:"send_time"`
		PlayerName string    `json:"player_name"`
		Number     int       `json:"number"`
		Message    string    `json:"message"`
	}

	PlayerCreate struct {
		Name string `json:"name" validate:"required"`
		Team string `json:"team" validate:"required"`
	}
)

func initChatInterface(group *echo.Group, api *EchoApi) {
	group.POST("/:"+lobby_id_param+message_path, api.createMessageId)
	group.PUT("/:"+lobby_id_param+message_path+"/:"+message_id_param, api.createMessage)
	group.PUT("/:"+lobby_id_param+player_path+"/:"+player_id_param, api.addPlayer)
	group.DELETE("/:"+lobby_id_param+player_path+"/:"+player_id_param, api.deletePlayer)
	group.GET("/:"+lobby_id_param+message_path+"/:"+number_id_param, api.getMessages)
}

func (api *EchoApi) createMessageId(context echo.Context) error {
	logger := context.Get(logger_key).(*log.Entry)
	logger.Debug("Create message Id")
	return context.String(http.StatusCreated, uuid.NewString())
}

func (api *EchoApi) createMessage(context echo.Context) error {
	logger := context.Get(logger_key).(*log.Entry)
	logger.Debug("Create message")

	message, lobbyId, messageId, err := bindMessageCreationDTO(context)

	if err != nil {
		logger.Warnf("Error while binding message: %v", err)
		return echo.ErrBadRequest
	}

	coreMessage := mapMessageCreateToMessage(messageId, lobbyId, message)
	err = api.core.CreateMessage(message.PlayerId, coreMessage)

	if err != nil {
		logger.Warnf("Error while creating message: %v", err)
		return echo.ErrInternalServerError
	}

	return context.NoContent(http.StatusCreated)
}

func (api *EchoApi) addPlayer(context echo.Context) error {
	logger := context.Get(logger_key).(*log.Entry)
	logger.Debug("Create player")

	player, lobbyId, playerId, err := bindPlayerCreationDTO(context)
	if err != nil {
		logger.Warnf("Error while binding player: %v", err)
		return echo.ErrBadRequest
	}

	corePlayer := mapPlayerCreateToPlayer(lobbyId, playerId, player)
	err = api.core.CreatePlayer(corePlayer)

	if err != nil {
		logger.Warnf("Error while creating player: %v", err)
		return echo.ErrInternalServerError
	}

	return context.NoContent(http.StatusCreated)
}

func (api *EchoApi) deletePlayer(context echo.Context) error {
	logger := context.Get(logger_key).(*log.Entry)
	logger.Debug("delete player")

	playerId, err := getParamPlayerId(context)
	if err != nil {
		logger.Warnf("Error while binding playerId: %v", err)
		return echo.ErrBadRequest
	}

	err = api.core.DeletePlayer(playerId)

	if err != nil {
		logger.Warnf("Error while deleting player: %v", err)
		return echo.ErrInternalServerError
	}

	return context.NoContent(http.StatusNoContent)
}

func (api *EchoApi) getMessages(context echo.Context) error {
	logger := context.Get(logger_key).(*log.Entry)
	logger.Debug("Get all messages")

	lobbyId, err := getLobbyId(context)
	if err != nil {
		logger.Warnf("Error while binding lobbyId: %v", err)
		return echo.ErrBadRequest
	}

	playerId, err := getQueryPlayerId(context)
	if err != nil {
		logger.Warnf("Error while binding playerId: %v", err)
		return echo.ErrBadRequest
	}

	number, err := getNumber(context)
	if err != nil {
		logger.Warnf("Error while binding number: %v", err)
		return echo.ErrBadRequest
	}

	messages, err := api.core.GetMessages(playerId, lobbyId, number)
	if err != nil {
		logger.Warnf("Error while loading messages: %v", err)
		return echo.ErrInternalServerError
	}
	return context.JSON(http.StatusOK, mapToMessages(messages))
}

func bindMessageCreationDTO(context echo.Context) (message *MessageCreate, lobbyId uuid.UUID, messageId uuid.UUID, err error) {
	message = new(MessageCreate)
	if err := context.Bind(message); err != nil {
		return nil, uuid.Nil, uuid.Nil, fmt.Errorf("could not bind message, %v", err)
	}
	if err := context.Validate(message); err != nil {
		return nil, uuid.Nil, uuid.Nil, fmt.Errorf("could not validate message, %v", err)
	}
	lobbyId, err = getLobbyId(context)
	if err != nil {
		return nil, uuid.Nil, uuid.Nil, err
	}

	messageId, err = getMessageId(context)
	if err != nil {
		return nil, uuid.Nil, uuid.Nil, err
	}

	return message, lobbyId, messageId, nil
}

func bindPlayerCreationDTO(context echo.Context) (player *PlayerCreate, lobbyId uuid.UUID, playerId uuid.UUID, err error) {
	player = new(PlayerCreate)
	if err := context.Bind(player); err != nil {
		return nil, uuid.Nil, uuid.Nil, fmt.Errorf("could not bind player, %v", err)
	}
	if err := context.Validate(player); err != nil {
		return nil, uuid.Nil, uuid.Nil, fmt.Errorf("could not validate player, %v", err)
	}
	lobbyId, err = getLobbyId(context)
	if err != nil {
		return nil, uuid.Nil, uuid.Nil, err
	}

	playerId, err = getParamPlayerId(context)
	if err != nil {
		return nil, uuid.Nil, uuid.Nil, err
	}

	return player, lobbyId, playerId, nil
}

func getLobbyId(context echo.Context) (uuid.UUID, error) {
	lobbyId, err := uuid.Parse(context.Param(lobby_id_param))
	if err != nil {
		return uuid.Nil, fmt.Errorf("error while binding lobbbyId: %v", err)
	}
	return lobbyId, nil
}

func getMessageId(context echo.Context) (uuid.UUID, error) {
	messageId, err := uuid.Parse(context.Param(message_id_param))
	if err != nil {
		return uuid.Nil, fmt.Errorf("error while binding messageId: %v", err)
	}
	return messageId, nil
}

func getNumber(context echo.Context) (int, error) {
	number, err := strconv.Atoi(context.Param(number_id_param))
	if err != nil {
		return -1, fmt.Errorf("error while binding number: %v", err)
	}
	return number, nil
}

func getQueryPlayerId(context echo.Context) (uuid.UUID, error) {
	playerId, err := uuid.Parse(context.QueryParam(player_id_param))
	if err != nil {
		return uuid.Nil, fmt.Errorf("error while binding playerId: %v", err)
	}
	return playerId, nil
}

func getParamPlayerId(context echo.Context) (uuid.UUID, error) {
	playerId, err := uuid.Parse(context.Param(player_id_param))
	if err != nil {
		return uuid.Nil, fmt.Errorf("error while binding playerId: %v", err)
	}
	return playerId, nil
}

func mapMessageCreateToMessage(messageId uuid.UUID, lobbyId uuid.UUID, message *MessageCreate) *core.Message {
	return &core.Message{ID: messageId, SendTime: time.Now(), LobbyId: lobbyId, Message: message.Message}
}

func mapToMessages(coreMessages []*core.Message) []*Message {
	messages := make([]*Message, len(coreMessages))
	for index, message := range coreMessages {
		messages[index] = mapToMessage(message)
	}
	return messages
}

func mapToMessage(message *core.Message) *Message {
	return &Message{ID: message.ID, SendTime: message.SendTime, PlayerName: message.PlayerName, Number: message.Number, Message: message.Message}
}

func mapPlayerCreateToPlayer(lobbyId uuid.UUID, playerId uuid.UUID, player *PlayerCreate) *core.Player {
	return &core.Player{ID: playerId, LobbyId: lobbyId, Name: player.Name, Team: player.Team}
}
