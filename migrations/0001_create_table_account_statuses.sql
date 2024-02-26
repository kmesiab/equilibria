-- +goose Up
CREATE TABLE account_statuses
(
    id   BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64) UNIQUE NOT NULL
);

INSERT INTO account_statuses (name)
VALUES ('Pending Activation'),
       ('Active'),
       ('Suspended'),
       ('Expired');
-- +goose Down
DROP TABLE IF EXISTS account_statuses;
