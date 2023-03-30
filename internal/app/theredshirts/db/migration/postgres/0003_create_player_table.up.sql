CREATE TABLE theredshirts_chat.player (
    id uuid PRIMARY KEY NOT NULL,
    lobby_id uuid NOT NULL,
    name varchar NOT NULL,
    team varchar NOT NULL
);