-- +goose Up
-- This section is executed when the migration is applied.

CREATE TABLE nrclex
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id    BIGINT NOT NULL,
    message_id BIGINT NOT NULL,
    anger          decimal(10, 2) NOT NULL DEFAULT 0,
    anticipation   decimal(10, 2) NOT NULL DEFAULT 0,
    disgust        decimal(10, 2) NOT NULL DEFAULT 0,
    fear           decimal(10, 2) NOT NULL DEFAULT 0,
    trust          decimal(10, 2) NOT NULL DEFAULT 0,
    joy            decimal(10, 2) NOT NULL DEFAULT 0,
    negative       decimal(10, 2) NOT NULL DEFAULT 0,
    positive       decimal(10, 2) NOT NULL DEFAULT 0,
    sadness        decimal(10, 2) NOT NULL DEFAULT 0,
    surprise       decimal(10, 2) NOT NULL DEFAULT 0,
    vader_compound decimal(10, 2) NOT NULL DEFAULT 0,
    vader_neg      decimal(10, 2) NOT NULL DEFAULT 0,
    vader_neu      decimal(10, 2) NOT NULL DEFAULT 0,
    vader_pos      decimal(10, 2) NOT NULL DEFAULT 0,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (message_id) REFERENCES messages (id)

);

-- +goose Down
-- This section is executed when the migration is rolled back.

DROP TABLE IF EXISTS nrclex;
-- This command removes the 'transactions' table if it exists.
