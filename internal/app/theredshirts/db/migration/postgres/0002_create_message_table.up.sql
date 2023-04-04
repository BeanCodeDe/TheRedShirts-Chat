CREATE TABLE theredshirts_message.message (
    id uuid PRIMARY KEY NOT NULL,
    send_time timestamp NOT NULL,
    lobby_id uuid NOT NULL,
    player_id uuid NOT NULL,
    number SERIAL NOT NULL,
    topic varchar NOT NULL,
    message json NOT NULL
);