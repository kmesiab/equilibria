-- +goose Up
-- This section is executed when the migration is applied.

CREATE TABLE facts
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    -- 'id' is a unique identifier for each conversation record.

    user_id    BIGINT NOT NULL,
    -- 'user_id' links the fact to a user.

    conversation_id BIGINT NOT NULL,
    -- 'message_id' links the fact to the message that the fact was derived from

    body       TEXT,
    -- 'body' is the body text of the fact

    reasoning  TEXT,
    -- 'reasoning' is the reasoning behind why this is a valid fact

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    -- 'created_at' records the date and time when the fact record was created.

    deleted_at DATETIME DEFAULT NULL,
    -- 'deleted_at' marks the timestamp when the fact record was deleted, useful for soft deletes.

    updated_at DATETIME DEFAULT NULL,
    -- 'updated_at' marks the timestamp when the fact record was last updated.

    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (conversation_id) REFERENCES conversations (id),

    INDEX (user_id, deleted_at)

);

-- +goose Down
-- This section is executed when the migration is rolled back.

DROP TABLE IF EXISTS facts;
-- This command removes the 'conversations' table if it exists.
