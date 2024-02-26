-- +goose Up
CREATE TABLE message_types
(
    id                   INT AUTO_INCREMENT PRIMARY KEY,
    name                 VARCHAR(64) UNIQUE NOT NULL,
    bill_rate_in_credits FLOAT                 NOT NULL DEFAULT 1
);

INSERT INTO message_types (name, bill_rate_in_credits)
VALUES ('Push Notification', 1),
       ('SMS', 1),
       ('Horoscope', 0),
       ('Email', 0),
       ('Facebook Messenger', 1),
       ('Twitter', 1),
       ('WhatsApp', 1),
       ('Telegram', 1),
       ('WeChat', 1),
       ('Instagram', 1),
       ('SnapChat', 1);

-- +goose Down
DROP TABLE IF EXISTS message_types;
