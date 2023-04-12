package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

const (
	message_table_name                  = "message"
	create_message_sql                  = "INSERT INTO %s.%s(id, send_time, lobby_id, player_id, topic, message) VALUES($1, $2, $3, $4, $5, $6)"
	select_messages_by_lobby_and_number = "SELECT id, send_time, lobby_id, player_id, number, topic, message FROM %s.%s WHERE lobby_id = $1 AND player_id != $2 AND number > $3"
	select_first_messages_of_player     = "SELECT id, send_time, lobby_id, player_id, number, topic, message FROM %s.%s WHERE lobby_id = $1 AND player_id != $2 AND number > (SELECT number FROM %s.%s WHERE lobby_id = $1 AND player_id = $2 AND topic = 'PLAYER_JOINS_LOBBY' ORDER BY number ASC LIMIT 1)"
	delete_messages_by_older_then       = "DELETE FROM %s.%s WHERE send_time < $1"
)

var (
	ErrMessageAlreadyExists = errors.New("message already exists")
)

func (tx *postgresTransaction) CreateMessage(message *Message) error {
	if _, err := tx.tx.Exec(context.Background(), fmt.Sprintf(create_message_sql, schema_name, message_table_name), message.ID, message.SendTime, message.LobbyId, message.PlayerId, message.Topic, message.Message); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return ErrMessageAlreadyExists
			}
		}

		return fmt.Errorf("unknown error when inserting message: %v", err)
	}
	return nil
}

func (tx *postgresTransaction) GetMessages(lobbyId uuid.UUID, toIgnoreplayerId uuid.UUID, number int) ([]*Message, error) {
	var messages []*Message
	if err := pgxscan.Select(context.Background(), tx.tx, &messages, fmt.Sprintf(select_messages_by_lobby_and_number, schema_name, message_table_name), lobbyId, toIgnoreplayerId, number); err != nil {
		return nil, fmt.Errorf("error while selecting all messages: %v", err)
	}

	return messages, nil
}

func (tx *postgresTransaction) GetMessagesFirstRequest(lobbyId uuid.UUID, toIgnoreplayerId uuid.UUID) ([]*Message, error) {
	var messages []*Message
	if err := pgxscan.Select(context.Background(), tx.tx, &messages, fmt.Sprintf(select_first_messages_of_player, schema_name, message_table_name, schema_name, message_table_name), lobbyId, toIgnoreplayerId); err != nil {
		return nil, fmt.Errorf("error while selecting first messages: %v", err)
	}

	return messages, nil
}

func (tx *postgresTransaction) DeleteMessages(time time.Time) error {
	if _, err := tx.tx.Exec(context.Background(), fmt.Sprintf(delete_messages_by_older_then, schema_name, message_table_name), time); err != nil {
		return fmt.Errorf("unknown error when deliting messages: %v", err)
	}
	return nil
}
