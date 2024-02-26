-- +goose Up
-- This section is executed when the migration is applied.

CREATE TABLE messages
(
    id                BIGINT AUTO_INCREMENT PRIMARY KEY,
    reference_id      VARCHAR(255),
    conversation_id   BIGINT,

    from_user_id      BIGINT             NOT NULL,
    to_user_id        BIGINT             NOT NULL,
    body              TEXT               NOT NULL,
    message_type_id   INT                NOT NULL,
    sent_at           DATETIME DEFAULT NULL,
    received_at       DATETIME DEFAULT NULL,
    message_status_id INT      DEFAULT 1 NOT NULL,

    created_at        DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at        DATETIME DEFAULT NULL,

    FOREIGN KEY (conversation_id) REFERENCES conversations (id),
    FOREIGN KEY (message_type_id) REFERENCES message_types (id),
    FOREIGN KEY (to_user_id) REFERENCES users (id),
    FOREIGN KEY (from_user_id) REFERENCES users (id),
    FOREIGN KEY (message_status_id) REFERENCES message_statuses (id),

    INDEX idx_reference_id (reference_id, deleted_at)
);

-- +goose Down
-- This section is executed when the migration is rolled back.

DROP TABLE IF EXISTS messages;
