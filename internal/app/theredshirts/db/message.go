package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

const (
	message_table_name                  = "message"
	create_message_sql                  = "INSERT INTO %s.%s(id, send_time, player_name, lobby_id, number, message) VALUES($1, $2, $3, $4, $5, $6)"
	select_messages_by_lobby_and_number = "SELECT id, send_time, player_name, lobby_id, number, message FROM %s.%s WHERE lobby_id = $1 AND number = $2"
)

var (
	ErrMessageAlreadyExists = errors.New("message already exists")
)

func (tx *postgresTransaction) CreateMessage(message *Message) error {
	if _, err := tx.tx.Exec(context.Background(), fmt.Sprintf(create_message_sql, schema_name, message_table_name), message.ID, message.SendTime, message.PlayerName, message.LobbyId, message.Number, message.Message); err != nil {
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

func (tx *postgresTransaction) GetMessages(lobbyId uuid.UUID, number int) ([]*Message, error) {
	var messages []*Message
	if err := pgxscan.Select(context.Background(), tx.tx, &messages, fmt.Sprintf(select_messages_by_lobby_and_number, schema_name, message_table_name), lobbyId, number); err != nil {
		return nil, fmt.Errorf("error while selecting all messages: %v", err)
	}

	return messages, nil
}
