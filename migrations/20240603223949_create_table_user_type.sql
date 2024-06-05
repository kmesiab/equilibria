-- +goose Up
-- This section is executed when the migration is applied.

CREATE TABLE user_types
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY,
    -- 'id' is a unique identifier for each conversation record.

    name       varchar(128) NOT NULL,
    -- 'name' is the string name of the user type

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    -- 'created_at' records the date and time when the fact record was created.

    deleted_at DATETIME DEFAULT NULL,
    -- 'deleted_at' marks the timestamp when the fact record was deleted, useful for soft deletes.

    updated_at DATETIME DEFAULT NULL
    -- 'updated_at' marks the timestamp when the fact record was last updated.
);

-- Insert a new record into the 'user_types' table
INSERT INTO user_types (name)
VALUES ('Patient'),
       ('Therapist');

-- +goose Down
-- This section is executed when the migration is rolled back.

DROP TABLE IF EXISTS user_types;
-- This command removes the 'conversations' table if it exists.
