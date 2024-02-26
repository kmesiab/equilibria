-- +goose Up
CREATE TABLE users
(
    id                BIGINT AUTO_INCREMENT PRIMARY KEY,
    phone_number      VARCHAR(20) UNIQUE NOT NULL,
    phone_verified    BOOLEAN            NOT NULL DEFAULT FALSE,
    firstname         VARCHAR(100),
    lastname          VARCHAR(100),
    password          VARCHAR(1024),
    email             VARCHAR(100),
    account_status_id BIGINT             NOT NULL DEFAULT 1,
    created_at        DATETIME           not null DEFAULT CURRENT_TIMESTAMP,
    deleted_at        DATETIME                    DEFAULT null,
    updated_at        DATETIME                    DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (account_status_id) REFERENCES account_statuses (id)
);

INSERT INTO users (phone_number, phone_verified, firstname, lastname, email, account_status_id)
VALUES ('+18333595081', TRUE, 'System', 'User', '-@-', 2),
       ('+13607102634', TRUE, 'Deanne', 'Doucette', 'deannedoucette@gmail.com', 2),
       ('+17072468797', TRUE, 'Dennis', 'Christo', 'bigplans777@gmail.com', 2),
       ('+12533243071', TRUE, 'Kevin', 'Mesiab', 'kmesiab+equilibria_sms@gmail.com', 2);

-- +goose Down
DROP TABLE IF EXISTS users;

