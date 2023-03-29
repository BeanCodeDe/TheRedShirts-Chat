CREATE TABLE theredshirts_chat.message (
    id uuid PRIMARY KEY NOT NULL,
    send_time timestamp NOT NULL,
    player_name varchar NOT NULL,
    lobby_id uuid NOT NULL,
    number number NOT NULL,
    message varchar NOT NULL
);