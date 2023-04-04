CREATE INDEX messages_idx ON theredshirts_message.message (lobby_id, player_id, number DESC);
CREATE INDEX messages_time_idx ON theredshirts_message.message (send_time);