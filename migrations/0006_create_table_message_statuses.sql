-- +goose Up
CREATE TABLE message_statuses
(
    id   INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64) UNIQUE NOT NULL
);

INSERT INTO message_statuses (name)
VALUES ('Pending'),
       ('Sent'),
       ('Received'),
       ('Delivered'),
       ('Canceled'),
       ('Failed'),
       ('Accepted'),
       ('Queued'),
       ('Receiving'),
       ('Read'),
       ('Sending'),
       ('Unknown');

-- +goose Down
DROP TABLE IF EXISTS message_statuses;
