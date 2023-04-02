package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/BeanCodeDe/TheRedShirts-Chat/internal/app/theredshirts/core"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

const message_root_path = "/message"
const message_path = "/msg"
const player_path = "/player"
const message_id_param = "messageId"
const number_id_param = "number"
const lobby_id_param = "lobbyId"
const player_id_param = "playerId"

type (
	MessageCreate struct {
		ID      uuid.UUID              `param:"messageId" validate:"required"`
		LobbyId uuid.UUID              `param:"lobbyId" validate:"required"`
		Topic   string                 `json:"topic" validate:"required"`
		Message map[string]interface{} `json:"message"`
	}

	MessageGet struct {
		LobbyId uuid.UUID `param:"lobbyId" validate:"required"`
		Number  int       `param:"number" validate:"required"`
	}

	Message struct {
		ID       uuid.UUID              `json:"id"`
		SendTime time.Time              `json:"send_time"`
		Number   int                    `json:"number"`
		Topic    string                 `json:"topic"`
		Message  map[string]interface{} `json:"message"`
	}

	PlayerCreate struct {
		ID      uuid.UUID `param:"playerId" validate:"required"`
		LobbyId uuid.UUID `param:"lobbyId" validate:"required"`
	}

	PlayerDelete struct {
		ID      uuid.UUID `param:"playerId" validate:"required"`
		LobbyId uuid.UUID `param:"lobbyId" validate:"required"`
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

	message, err := bindMessageCreationDTO(context)
	if err != nil {
		logger.Warnf("Error while binding message: %v", err)
		return echo.ErrBadRequest
	}
	playerId, err := getQueryPlayerId(context)
	if err != nil {
		logger.Warnf("Error while binding playerId: %v", err)
		return echo.ErrBadRequest
	}

	coreMessage := mapMessageCreateToMessage(message)
	err = api.core.CreateMessage(playerId, coreMessage)

	if err != nil {
		logger.Warnf("Error while creating message: %v", err)
		return echo.ErrInternalServerError
	}

	return context.NoContent(http.StatusCreated)
}

func (api *EchoApi) addPlayer(context echo.Context) error {
	logger := context.Get(logger_key).(*log.Entry)
	logger.Debug("Create player")

	player, err := bindPlayerCreationDTO(context)
	if err != nil {
		logger.Warnf("Error while binding player: %v", err)
		return echo.ErrBadRequest
	}

	corePlayer := mapPlayerCreateToPlayer(player)
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

	player, err := bindPlayerDelete(context)
	if err != nil {
		logger.Warnf("Error while binding playerId: %v", err)
		return echo.ErrBadRequest
	}

	err = api.core.DeletePlayer(player.ID)
	if err != nil {
		logger.Warnf("Error while deleting player: %v", err)
		return echo.ErrInternalServerError
	}

	return context.NoContent(http.StatusNoContent)
}

func (api *EchoApi) getMessages(context echo.Context) error {
	logger := context.Get(logger_key).(*log.Entry)
	logger.Debug("Get all messages")

	message, err := bindMessageGet(context)
	if err != nil {
		logger.Warnf("Error while binding get message: %v", err)
		return echo.ErrBadRequest
	}

	playerId, err := getQueryPlayerId(context)
	if err != nil {
		logger.Warnf("Error while binding playerId: %v", err)
		return echo.ErrBadRequest
	}

	messages, err := api.core.GetMessages(playerId, message.LobbyId, message.Number)
	if err != nil {
		logger.Warnf("Error while loading messages: %v", err)
		return echo.ErrInternalServerError
	}
	return context.JSON(http.StatusOK, mapToMessages(messages))
}

func bindMessageCreationDTO(context echo.Context) (message *MessageCreate, err error) {
	message = new(MessageCreate)
	if err := context.Bind(message); err != nil {
		return nil, fmt.Errorf("could not bind message, %v", err)
	}
	if err := context.Validate(message); err != nil {
		return nil, fmt.Errorf("could not validate message, %v", err)
	}

	return message, nil
}

func bindMessageGet(context echo.Context) (message *MessageGet, err error) {
	message = new(MessageGet)
	if err := context.Bind(message); err != nil {
		return nil, fmt.Errorf("could not bind message, %v", err)
	}
	if err := context.Validate(message); err != nil {
		return nil, fmt.Errorf("could not validate message, %v", err)
	}

	return message, nil
}

func bindPlayerCreationDTO(context echo.Context) (player *PlayerCreate, err error) {
	player = new(PlayerCreate)
	if err := context.Bind(player); err != nil {
		return nil, fmt.Errorf("could not bind player, %v", err)
	}
	if err := context.Validate(player); err != nil {
		return nil, fmt.Errorf("could not validate player, %v", err)
	}

	return player, nil
}

func bindPlayerDelete(context echo.Context) (player *PlayerDelete, err error) {
	player = new(PlayerDelete)
	if err := context.Bind(player); err != nil {
		return nil, fmt.Errorf("could not bind player, %v", err)
	}
	if err := context.Validate(player); err != nil {
		return nil, fmt.Errorf("could not validate player, %v", err)
	}

	return player, nil
}

func getQueryPlayerId(context echo.Context) (uuid.UUID, error) {
	playerId, err := uuid.Parse(context.QueryParam(player_id_param))
	if err != nil {
		return uuid.Nil, fmt.Errorf("error while binding playerId: %v", err)
	}
	return playerId, nil
}

func mapMessageCreateToMessage(message *MessageCreate) *core.Message {
	return &core.Message{ID: message.ID, SendTime: time.Now(), LobbyId: message.LobbyId, Message: message.Message}
}

func mapToMessages(coreMessages []*core.Message) []*Message {
	messages := make([]*Message, len(coreMessages))
	for index, message := range coreMessages {
		messages[index] = mapToMessage(message)
	}
	return messages
}

func mapToMessage(message *core.Message) *Message {
	return &Message{ID: message.ID, SendTime: message.SendTime, Number: message.Number, Topic: message.Topic, Message: message.Message}
}

func mapPlayerCreateToPlayer(player *PlayerCreate) *core.Player {
	return &core.Player{ID: player.ID, LobbyId: player.LobbyId, LastRefresh: time.Now()}
}
