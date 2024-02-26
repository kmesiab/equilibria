-- +goose Up
-- This section is executed when the migration is applied.

CREATE TABLE conversations
(
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    -- 'id' is a unique identifier for each conversation record.

    user_id BIGINT NOT NULL,
    -- 'user_id' links the conversation to a user, usually the one who initiated the conversation.

    start_time DATETIME NOT NULL,
    -- 'start_time' records when the conversation started.

    end_time DATETIME,
    -- 'end_time' records when the conversation ended.

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    -- 'created_at' records the date and time when the conversation record was created.

    deleted_at DATETIME,
    -- 'deleted_at' marks the timestamp when the conversation record was deleted, useful for soft deletes.

    updated_at DATETIME,
    -- 'updated_at' marks the timestamp when the conversation record was last updated.

    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +goose Down
-- This section is executed when the migration is rolled back.

DROP TABLE IF EXISTS conversations;
-- This command removes the 'conversations' table if it exists.
