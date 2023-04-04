package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/BeanCodeDe/TheRedShirts-Message/internal/app/theredshirts/core"
	"github.com/BeanCodeDe/TheRedShirts-Message/internal/app/theredshirts/util"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const message_root_path = "/message"
const message_path = "/msg"
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
		PlayerId uuid.UUID              `json:"player_id"`
		SendTime time.Time              `json:"send_time"`
		Number   int                    `json:"number"`
		Topic    string                 `json:"topic"`
		Message  map[string]interface{} `json:"message"`
	}
)

func initChatInterface(group *echo.Group, api *EchoApi) {
	group.POST("/:"+lobby_id_param+message_path, api.createMessageId)
	group.PUT("/:"+lobby_id_param+message_path+"/:"+message_id_param, api.createMessage)
	group.GET("/:"+lobby_id_param+message_path+"/:"+number_id_param, api.getMessages)
}

func (api *EchoApi) createMessageId(context echo.Context) error {
	customContext := context.Get(context_key).(*util.Context)
	logger := customContext.Logger
	logger.Debug("Create message Id")
	return context.String(http.StatusCreated, uuid.NewString())
}

func (api *EchoApi) createMessage(context echo.Context) error {
	customContext := context.Get(context_key).(*util.Context)
	logger := customContext.Logger
	logger.Debug("Create message")

	message, err := bindMessageCreationDTO(context)
	if err != nil {
		logger.Warnf("Error while binding message: %v", err)
		return echo.ErrBadRequest
	}
	playerId, err := getHeaderPlayerId(context)
	if err != nil {
		logger.Warnf("Error while binding playerId: %v", err)
		return echo.ErrBadRequest
	}

	coreMessage := mapMessageCreateToMessage(message)
	err = api.core.CreateMessage(customContext, playerId, coreMessage)

	if err != nil {
		logger.Warnf("Error while creating message: %v", err)
		return echo.ErrInternalServerError
	}

	return context.NoContent(http.StatusCreated)
}

func (api *EchoApi) getMessages(context echo.Context) error {
	customContext := context.Get(context_key).(*util.Context)
	logger := customContext.Logger
	logger.Debug("Get all messages")

	message, err := bindMessageGet(context)
	if err != nil {
		logger.Warnf("Error while binding get message: %v", err)
		return echo.ErrBadRequest
	}

	playerId, err := getHeaderPlayerId(context)
	if err != nil {
		logger.Warnf("Error while binding playerId: %v", err)
		return echo.ErrBadRequest
	}

	messages, err := api.core.GetMessages(customContext, playerId, message.LobbyId, message.Number)
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

func getHeaderPlayerId(context echo.Context) (uuid.UUID, error) {
	playerId, err := uuid.Parse(context.Request().Header.Get(player_id_param))
	if err != nil {
		return uuid.Nil, fmt.Errorf("error while binding playerId: %v", err)
	}
	return playerId, nil
}

func mapMessageCreateToMessage(message *MessageCreate) *core.Message {
	return &core.Message{ID: message.ID, SendTime: time.Now(), LobbyId: message.LobbyId, Topic: message.Topic, Message: message.Message}
}

func mapToMessages(coreMessages []*core.Message) []*Message {
	messages := make([]*Message, len(coreMessages))
	for index, message := range coreMessages {
		messages[index] = mapToMessage(message)
	}
	return messages
}

func mapToMessage(message *core.Message) *Message {
	return &Message{ID: message.ID, PlayerId: message.PlayerId, SendTime: message.SendTime, Number: message.Number, Topic: message.Topic, Message: message.Message}
}
