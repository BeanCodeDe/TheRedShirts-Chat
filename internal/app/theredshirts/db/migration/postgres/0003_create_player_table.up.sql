CREATE TABLE theredshirts_message.player (
    id uuid PRIMARY KEY NOT NULL,
    lobby_id uuid NOT NULL,
    last_refresh timestamp NOT NULL
);