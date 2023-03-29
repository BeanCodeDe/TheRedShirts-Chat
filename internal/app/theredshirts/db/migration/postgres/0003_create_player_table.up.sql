CREATE TABLE theredshirts_chat.player (
    id uuid PRIMARY KEY NOT NULL,
    lobby_id uuid NOT NULL,
    name varchar NOT NULL,
    team varchar NOT NULL,
);


	Player struct {
		ID      uuid.UUID `db:"id"`
		LobbyId uuid.UUID `db:"lobby_id"`
		Name    string    `db:"name"`
		Team    string    `db:"team"`
	}